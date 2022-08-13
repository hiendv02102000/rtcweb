package router

import (
	"chat-service/handler"
	"chat-service/pkg/share/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
}

func (r *Router) Routes() {
	handler.AllRooms.Init()
	wsChat := r.Engine.Group("/chat_service")
	{
		wsChat.Use(middleware.AuthMiddleware(), middleware.AuthUserBanned())
		wsChat.POST("/send_message", handler.SendMessage)
		wsChat.GET("/join_room", handler.JoinRoomChat)
		wsChat.GET("/get_message", handler.GetMessage)
	} //

}
func NewRouter() Router {
	var r Router
	r.Engine = gin.Default()
	r.Routes()
	return r
}
