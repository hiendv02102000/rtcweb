package handler

import (
	"chat-service/dto"
	db "chat-service/pkg/infrastucture/database"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	guuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var AllRooms RoomMap

func SendMessage(c *gin.Context) {
	roomId := c.Query("room_id")
	msg := dto.MessageChat{}
	err := c.ShouldBind(msg)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	if msg.SendAt == nil {
		timeNow := time.Now()
		msg.SendAt = &timeNow
	}

	score := msg.SendAt.Unix() / 10000
	res, err := db.RedisPool.ZAdd(c, "room-chat-"+roomId, &redis.Z{Member: msg, Score: float64(score)}).Result()
	if err != nil || res <= 0 {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  "send message failure",
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	room, ok := AllRooms.Get(roomId)
	if !ok {
		room = AllRooms.CreateRoom(roomId)
	}
	var wg sync.WaitGroup
	for _, chanel := range room.Chanels {
		wg.Add(1)
		go func(w *sync.WaitGroup, ch chan dto.MessageChat, ms dto.MessageChat) {
			defer wg.Done()
			ch <- ms
		}(&wg, chanel, msg)
	}
	go func(w *sync.WaitGroup) {
		time.Sleep(10 * 60 * time.Second)
		for {
			wg.Done()
		}
	}(&wg)
	wg.Wait()
	data := dto.BaseResponse{
		Status: http.StatusOK,
		Result: "success",
	}
	c.JSON(http.StatusOK, data)
}

func GetMessage(c *gin.Context) {
	roomId := c.Query("room_id")
	startS := c.Query("start")
	sizeS := c.Query("size")
	start, errP := strconv.Atoi(startS)
	size, errS := strconv.Atoi(sizeS)
	if errP != nil || errS != nil {

		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  "page size is int",
		}
		c.JSON(http.StatusBadRequest, data)
		return

	}

	res, err := db.RedisPool.ZRevRange(c, "room-chat-"+roomId, int64(start), int64(size+start)).Result()
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	dataMess := []dto.MessageChat{}
	for _, msgString := range res {
		mess := dto.MessageChat{}
		errC := json.Unmarshal([]byte(msgString), &mess)
		if errC != nil {
			continue
		}
		dataMess = append(dataMess, mess)
	}
	data := dto.BaseResponse{
		Status: http.StatusOK,
		Result: dataMess,
	}
	c.JSON(http.StatusOK, data)
}
func JoinRoomChat(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	roomId := c.Query("room_id")
	room, ok := AllRooms.Get(roomId)
	if !ok {
		room = AllRooms.CreateRoom(roomId)
	}
	channelId := guuid.New().String()
	AllRooms.InsertIntoRoom(roomId, channelId)
	defer AllRooms.LeaveRoom(roomId, channelId)
	channel := room.Chanels[roomId]
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
	for {
		msg := <-channel
		err := ws.WriteJSON(msg)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Read Error: "+err.Error())
			ws.Close()
			return
		}
	}
	data := dto.BaseResponse{
		Status: http.StatusOK,
		Result: "success",
	}
	c.JSON(http.StatusOK, data)
}
