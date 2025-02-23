package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"websocket-chat/controller"
	"websocket-chat/models"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// WebSocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client map
var clients = make(map[*websocket.Conn]string)
var mu sync.Mutex

// Channel untuk broadcast pesan
var broadcast = make(chan models.Message)

// Handler WebSocket
func handleConnectionsController(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade gagal:", err)
		return
	}
	defer ws.Close()

	// Ambil user ID dari parameter
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		log.Println("User ID tidak ditemukan")
		return
	}

	mu.Lock()
	clients[ws] = userID
	mu.Unlock()

	log.Println("User terhubung:", userID)

	messages, err := models.GetMessagesByUserID(db, userID)
	if err != nil {
		log.Println("Gagal mendapatkan pesan:", err)
	}

	// Send previous messages to the client
	for _, msg := range messages {
		if err := ws.WriteJSON(msg); err != nil {
			log.Println("Gagal mengirim pesan:", err)
			break
		}
	}

	for {
		var msg models.Message
		if err := ws.ReadJSON(&msg); err != nil {
			log.Println("Kesalahan membaca pesan:", err)
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			break
		}

		if msg.ReceiverID == "" {
			log.Println("Receiver ID kosong, tidak bisa buat chat!")
			continue
		}

		msg.SenderID = userID

		chatID, err := models.GetExistingChatID(db, msg.SenderID, msg.ReceiverID)
		if err != nil {
			log.Println("Gagal mendapatkan chat ID:", err)
			return
		}

		if chatID == "" {
			// Chat belum ada, buat chat baru
			chatID, err = models.CreateNewChat(db)
			if err != nil {
				log.Println("Gagal membuat chat baru:", err)
				return
			}
		}

		if chatID == "" {
			// Kalau chat belum ada, buat baru
			chatID, err = models.CreateNewChat(db)
			if err != nil {
				log.Println("Gagal membuat chat baru:", err)
				continue
			}
		}

		msg.ChatID = chatID

		// Simpan pesan ke database
		if err := models.SaveMessageToDB(db, msg.SenderID, chatID, msg.ReceiverID, msg.Content); err != nil {
			log.Println("Gagal menyimpan pesan:", err)
			continue
		}

		// Kirim pesan ke penerima
		broadcast <- msg
	}
}

// Broadcast pesan ke penerima
func handleMessagesModel() {
	for {
		msg := <-broadcast

		mu.Lock()
		for client, userID := range clients {
			if userID == msg.ReceiverID || userID == msg.SenderID {
				if err := client.WriteJSON(msg); err != nil {
					log.Println("Kesalahan mengirim pesan:", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
		mu.Unlock()
	}
}

func main() {
	controller.ConnectDatabase()
	http.HandleFunc("/ws/chat", func(w http.ResponseWriter, r *http.Request) {
		handleConnectionsController(controller.GetDB(), w, r)
	})

	go handleMessagesModel()

	fmt.Println("WebSocket server berjalan di :7070")
	log.Fatal(http.ListenAndServe(":7070", nil))
}
