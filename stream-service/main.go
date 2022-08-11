package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
)

const (
	rtcpPLIInterval = time.Second * 3
)

// Sdp represent session description protocol describe media communication sessions
type Sdp struct {
	Sdp string
}

func main() {

	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)
	router := gin.Default()

	// sender to channel of track
	peerConnectionMap := make(map[string]chan *webrtc.Track)

	m := webrtc.MediaEngine{}

	// Setup the codecs you want to use.
	// Only support VP8(video compression), this makes our proxying code simpler
	m.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))

	api := webrtc.NewAPI(webrtc.WithMediaEngine(m))

	peerConnectionConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	router.POST("/webrtc/sdp/m/:meetingId/c/:userID/p/:peerId/s/:isSender", func(c *gin.Context) {
		isSender, _ := strconv.ParseBool(c.Param("isSender"))
		userID := c.Param("userID")
		peerID := c.Param("peerId")

		var session Sdp
		if err := c.ShouldBindJSON(&session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		offer := webrtc.SessionDescription{}
		Decode(session.Sdp, &offer)

		// Create a new RTCPeerConnection
		// this is the gist of webrtc, generates and process SDP
		peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
		if err != nil {
			log.Fatal(err)
		}
		if !isSender {
			recieveTrack(peerConnection, peerConnectionMap, peerID)
		} else {
			createTrack(peerConnection, peerConnectionMap, userID)
		}
		// Set the SessionDescription of remote peer
		peerConnection.SetRemoteDescription(offer)

		// Create answer
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			log.Fatal(err)
		}

		// Sets the LocalDescription, and starts our UDP listeners
		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, Sdp{Sdp: Encode(answer)})
	})

	router.Run(":8080")
}

// user is the caller of the method
// if user connects before peer: create channel and keep listening till track is added
// if peer connects before user: channel would have been created by peer and track can be added by getting the channel from cache
func recieveTrack(peerConnection *webrtc.PeerConnection,
	peerConnectionMap map[string]chan *webrtc.Track,
	peerID string) {
	if _, ok := peerConnectionMap[peerID]; !ok {
		peerConnectionMap[peerID] = make(chan *webrtc.Track, 1)
	}
	localTrack := <-peerConnectionMap[peerID]
	peerConnection.AddTrack(localTrack)
}

// user is the caller of the method
// if user connects before peer: since user is first, user will create the channel and track and will pass the track to the channel
// if peer connects before user: since peer came already, he created the channel and is listning and waiting for me to create and pass track
func createTrack(peerConnection *webrtc.PeerConnection,
	peerConnectionMap map[string]chan *webrtc.Track,
	currentUserID string) {

	if _, err := peerConnection.AddTransceiver(webrtc.RTPCodecTypeVideo); err != nil {
		log.Fatal(err)
	}

	// Set a handler for when a new remote track starts, this just distributes all our packets
	// to connected peers
	peerConnection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		// This can be less wasteful by processing incoming RTCP events, then we would emit a NACK/PLI when a viewer requests it
		go func() {
			ticker := time.NewTicker(rtcpPLIInterval)
			for range ticker.C {
				if rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()

		// Create a local track, all our SFU clients will be fed via this track
		// main track of the broadcaster
		localTrack, newTrackErr := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
		if newTrackErr != nil {
			log.Fatal(newTrackErr)
		}

		// the channel that will have the local track that is used by the sender
		// the localTrack needs to be fed to the reciever
		localTrackChan := make(chan *webrtc.Track, 1)
		localTrackChan <- localTrack
		if existingChan, ok := peerConnectionMap[currentUserID]; ok {
			// feed the exsiting track from user with this track
			existingChan <- localTrack
		} else {
			peerConnectionMap[currentUserID] = localTrackChan
		}

		rtpBuf := make([]byte, 1400)
		for { // for publisher only
			i, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				log.Fatal(readErr)
			}

			// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
			if _, err := localTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
				log.Fatal(err)
			}
		}
	})

}

// Allows compressing offer/answer to bypass terminal input limits.
const compress = false

// MustReadStdin blocks until input is received from stdin
func MustReadStdin() string {
	r := bufio.NewReader(os.Stdin)

	var in string
	for {
		var err error
		in, err = r.ReadString('\n')
		if err != io.EOF {
			if err != nil {
				panic(err)
			}
		}
		in = strings.TrimSpace(in)
		if len(in) > 0 {
			break
		}
	}

	fmt.Println("")

	return in
}

// Encode encodes the input in base64
// It can optionally zip the input before encoding
func Encode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	if compress {
		b = zip(b)
	}

	return base64.StdEncoding.EncodeToString(b)
}

// Decode decodes the input from base64
// It can optionally unzip the input after decoding
func Decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	if compress {
		b = unzip(b)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}

func zip(in []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)
	if err != nil {
		panic(err)
	}
	err = gz.Flush()
	if err != nil {
		panic(err)
	}
	err = gz.Close()
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func unzip(in []byte) []byte {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		panic(err)
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return res
}
