package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// AllRooms is the global hashmap for the server
var AllRooms RoomMap
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func CreateRoom(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	room := AllRooms.CreateRoom(0, "vv")

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Web Socket Upgrade Error")
		return
	}
	defer AllRooms.DeleteRoom(room.RoomID)
	for {
		var msg map[string]interface{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Read Error: "+err.Error())
			return
		}
		go broadcaster(msg, room)
	}
}
func broadcaster(msg map[string]interface{}, room Room) {

	for _, msgChannel := range room.Chanels {
		go func(c *chan MessageVideo) { *c <- msg }(&msgChannel)
	}

}

func JoinRoom(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	roomId := c.Query("uuid")
	if len(roomId) == 0 {
		c.JSON(http.StatusBadRequest, "uuid is required")
		return
	}
	if _, ok := AllRooms.Get(roomId); !ok {
		c.JSON(http.StatusBadRequest, "room is not exist")
		return
	}
	userChannelId := guuid.New().String()
	AllRooms.InsertIntoRoom(roomId, userChannelId)

	room, _ := AllRooms.Get(roomId)
	cMsg := room.Chanels[userChannelId]
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Web Socket Upgrade Error")
		return
	}
	for {

		msg := <-cMsg
		err := ws.WriteJSON(msg)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Read Error: "+err.Error())
			return
		}

	}
}
