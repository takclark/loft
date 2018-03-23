package main

import (
	"github.com/takclark/loft/server"
)

func main() {
	photoStreamServer := server.NewPhotoStreamServer()
	photoStreamServer.Start()
}
