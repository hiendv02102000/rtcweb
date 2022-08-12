package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"stream-service/dto"
	"stream-service/pkg/share/middleware"
	"stream-service/pkg/share/utils"
	"strings"
	"time"

	db "stream-service/pkg/infrastucture/database"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
)

var AllRooms RoomMap

const (
	rtcpPLIInterval = time.Second * 3
)

var m webrtc.MediaEngine = webrtc.MediaEngine{}
var api *webrtc.API
var peerConnectionConfig webrtc.Configuration

func InitHandler() {
	AllRooms.Init()
	m.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))
	api = webrtc.NewAPI(webrtc.WithMediaEngine(m))
	peerConnectionConfig = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
}
func StartStream(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	req := dto.StartRoomRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	clientToken := c.GetHeader("Authorization")
	extractedToken := strings.Split(clientToken, "Bearer ")
	clientToken = strings.TrimSpace(extractedToken[1])
	d, err := utils.SendRequest("POST", utils.HOST_ACCOUNT_SERVICE+"/api/room/start_room", clientToken+"-MyRoomKey", req)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	res := dto.BaseResponse{}
	err = json.Unmarshal(d, &res)
	if err != nil || res.Status != 200 {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  "start room failure",
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}

	defer utils.SendRequest("POST", utils.HOST_ACCOUNT_SERVICE+"/api/room/end_room", clientToken+"-MyRoomKey", nil)
	roomId := res.Result.(map[string]interface{})["id"].(string)
	AllRooms.CreateRoom(roomId)
	defer AllRooms.DeleteRoom(roomId)
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
	if err != nil {
		c.JSON(http.StatusBadRequest, "Web Socket Upgrade Error")
		return
	}
	room, _ := AllRooms.Get(roomId)
	for {
		var session webrtc.SessionDescription
		err := ws.ReadJSON(&session)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Read Error: "+err.Error())
			ws.Close()
			return
		}

		peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
		if err != nil {
			fmt.Println(err)
		}
		createTrack1(peerConnection, room)
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			fmt.Println(err)
		}

		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			fmt.Println(err)
		}
		err = ws.WriteJSON(answer)
		if err != nil {
			c.JSON(http.StatusBadRequest, "write Error: "+err.Error())
			ws.Close()
			return
		}
	}
}
func JoinStream(c *gin.Context) {
	roomId := c.Query("room_id")
	ex, errR := db.RedisPool.Exists(c, "room-user-"+roomId).Result()
	if ex <= 0 || errR != redis.Nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  "room is not exist",
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	_, errR = db.RedisPool.SAdd(c, "room-user-"+roomId, middleware.GetUserFromContext(c)).Result()
	if errR != redis.Nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  "join room fail",
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	defer db.RedisPool.SRem(c, "room-user-"+roomId, middleware.GetUserFromContext(c))
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
	if err != nil {
		c.JSON(http.StatusBadRequest, "Web Socket Upgrade Error")
		return
	}
	room, _ := AllRooms.Get(roomId)
	for {
		var session webrtc.SessionDescription
		err := ws.ReadJSON(&session)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Read Error: "+err.Error())
			ws.Close()
			return
		}

		peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
		if err != nil {
			fmt.Println(err)
		}
		err = peerConnection.SetRemoteDescription(session)
		if err != nil {
			fmt.Println(err)
		}
		recieveTrack1(peerConnection, room)
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			fmt.Println(err)
		}

		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			fmt.Println(err)
		}
		err = ws.WriteJSON(answer)
		if err != nil {
			c.JSON(http.StatusBadRequest, "write Error: "+err.Error())
			ws.Close()
			return
		}
	}
}
func createTrack1(peerConnection *webrtc.PeerConnection, room Room) {

	if _, err := peerConnection.AddTransceiver(webrtc.RTPCodecTypeVideo); err != nil {
		log.Fatal(err)
	}

	peerConnection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {

		go func() {
			ticker := time.NewTicker(rtcpPLIInterval)
			for range ticker.C {
				if rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()
		localTrack, newTrackErr := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
		if newTrackErr != nil {
			log.Fatal(newTrackErr)
		}

		// localTrackChan := make(chan *webrtc.Track, 1)
		// localTrackChan <- localTrack

		room.Chanels = localTrack

		rtpBuf := make([]byte, 1400)
		for {
			i, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				log.Fatal(readErr)
			}
			if _, err := localTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
				log.Fatal(err)
			}
		}
	})

}
func recieveTrack1(peerConnection *webrtc.PeerConnection,
	room Room) {
	if room.Chanels == nil {
		room.Chanels = new(webrtc.Track)
	}
	localTrack := room.Chanels
	peerConnection.AddTrack(localTrack)
}
