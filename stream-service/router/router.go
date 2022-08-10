package router

import (
	"stream-service/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
}

func (r *Router) Routes() {

	webSocket := r.Engine.Group("/ws")
	{
		webSocket.POST("/create_room", handler.CreateRoom)
		webSocket.GET("/join_room/:uuid", handler.JoinRoom)
	}

}
func NewRouter() Router {
	var r Router
	r.Engine = gin.Default()
	r.Routes()
	return r
}
