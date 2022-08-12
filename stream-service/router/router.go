package router

import (
	"stream-service/handler"
	"stream-service/pkg/share/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
}

func (r *Router) Routes() {
	handler.InitHandler()
	wsChat := r.Engine.Group("/stream_service")
	{
		wsChat.POST("/start_stream", handler.StartStream)
		wsChat.Use(middleware.AuthMiddleware())
		wsChat.POST("/join_stream", handler.JoinStream)
	} //
	//
}
func NewRouter() Router {
	var r Router
	r.Engine = gin.Default()
	r.Routes()
	return r
}
