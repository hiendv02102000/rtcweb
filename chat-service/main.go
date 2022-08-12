package main

import "chat-service/router"

func main() {
	r := router.NewRouter()
	r.Engine.Run(":80")
}
