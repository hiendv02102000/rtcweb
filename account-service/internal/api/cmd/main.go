package main

import "api/internal/api/router"

func main() {
	r := router.NewRouter()
	r.Engine.Run(":8080")
}
