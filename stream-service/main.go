package main

import "stream-service/router"

func main() {
	r := router.NewRouter()
	r.Engine.Run(":80")
}
