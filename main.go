package main

import (
	"websocket-chat/controller"
)

func main() {
	server := controller.Server{}
	server.ConnectDatabase()
	// WebSocket handler end to end

	go server.HandleMessagesGrupModel()
	go server.HandleMessagesPrivateModel()

	server.RunServer()

}
