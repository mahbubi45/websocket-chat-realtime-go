package main

import (
	"fmt"
	"log"
	"net/http"
	"websocket-chat/controller"
)

func main() {
	server := controller.Server{}
	server.ConnectDatabase()
	// WebSocket handler end to end

	go server.HandleMessagesGrupModel()
	go server.HandleMessagesPrivateModel()

	// WebSocket handler group
	http.HandleFunc("/ws/group", func(w http.ResponseWriter, r *http.Request) {
		server.HandleConnectionsGrupController(server.GetDB(), w, r)
	})

	http.HandleFunc("/ws/private", func(w http.ResponseWriter, r *http.Request) {
		server.HandleConnectionsPrivateMessageController(server.GetDB(), w, r)
	})

	fmt.Println("WebSocket server berjalan di :6070")
	log.Fatal(http.ListenAndServe(":6070", nil))
}
