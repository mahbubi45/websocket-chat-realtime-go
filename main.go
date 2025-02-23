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

// // WebSocket Upgrader
// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// // Client map
// var clients = make(map[*websocket.Conn]string)
// var mu sync.Mutex

// // Channel untuk broadcast pesan
// var broadcast = make(chan models.Message)

// // Handler WebSocket
// func handleConnectionsController(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("WebSocket upgrade gagal:", err)
// 		return
// 	}
// 	defer ws.Close()

// 	// Ambil user ID dari parameter
// 	userID := r.URL.Query().Get("user_id")
// 	if userID == "" {
// 		log.Println("User ID tidak ditemukan")
// 		return
// 	}

// 	mu.Lock()
// 	clients[ws] = userID
// 	mu.Unlock()

// 	log.Println("User terhubung:", userID)

// 	messages, err := models.GetMessagesByUserID(db, userID)
// 	if err != nil {
// 		log.Println("Gagal mendapatkan pesan:", err)
// 	}

// 	// Send previous messages to the client
// 	for _, msg := range messages {
// 		if err := ws.WriteJSON(msg); err != nil {
// 			log.Println("Gagal mengirim pesan:", err)
// 			break
// 		}
// 	}

// 	for {
// 		var msg models.Message
// 		if err := ws.ReadJSON(&msg); err != nil {
// 			log.Println("Kesalahan membaca pesan:", err)
// 			mu.Lock()
// 			delete(clients, ws)
// 			mu.Unlock()
// 			break
// 		}

// 		if msg.ReceiverID == "" {
// 			log.Println("Receiver ID kosong, tidak bisa buat chat!")
// 			continue
// 		}

// 		msg.SenderID = userID

// 		chatID, err := models.GetExistingChatID(db, msg.SenderID, msg.ReceiverID)
// 		if err != nil {
// 			log.Println("Gagal mendapatkan chat ID:", err)
// 			return
// 		}

// 		if chatID == "" {
// 			// Chat belum ada, buat chat baru
// 			chatID, err = models.CreateNewChat(db)
// 			if err != nil {
// 				log.Println("Gagal membuat chat baru:", err)
// 				return
// 			}
// 		}

// 		if chatID == "" {
// 			// Kalau chat belum ada, buat baru
// 			chatID, err = models.CreateNewChat(db)
// 			if err != nil {
// 				log.Println("Gagal membuat chat baru:", err)
// 				continue
// 			}
// 		}

// 		msg.ChatID = chatID

// 		// Simpan pesan ke database
// 		if err := models.SaveMessageToDB(db, msg.SenderID, chatID, msg.ReceiverID, msg.Content); err != nil {
// 			log.Println("Gagal menyimpan pesan:", err)
// 			continue
// 		}

// 		// Kirim pesan ke penerima
// 		broadcast <- msg
// 	}
// }

// // Broadcast pesan ke penerima
// func handleMessagesModel() {
// 	for {
// 		msg := <-broadcast

// 		mu.Lock()
// 		for client, userID := range clients {
// 			if userID == msg.ReceiverID || userID == msg.SenderID {
// 				if err := client.WriteJSON(msg); err != nil {
// 					log.Println("Kesalahan mengirim pesan:", err)
// 					client.Close()
// 					delete(clients, client)
// 				}
// 			}
// 		}
// 		mu.Unlock()
// 	}
// }

// func main() {
// 	controller.ConnectDatabase()
// 	http.HandleFunc("/ws/chat", func(w http.ResponseWriter, r *http.Request) {
// 		handleConnectionsController(controller.GetDB(), w, r)
// 	})

// 	go handleMessagesModel()

// 	fmt.Println("WebSocket server berjalan di :7070")
// 	log.Fatal(http.ListenAndServe(":7070", nil))
// }

// WebSocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Struktur data gabungan
type GroupMessage struct {
	Msg  models.Message
	Chat models.Chat
}

var muGourp sync.Mutex

// Channel yang menampung GroupMessage
var broadcastGroup = make(chan GroupMessage)
var groups = make(map[string]map[*websocket.Conn]string)

// Channel untuk broadcast pesan grup
func HandleConnectionsGrupController(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade gagal:", err)
		return
	}
	defer func() {
		// Hapus user saat disconnect
		muGourp.Lock()
		if groupID := r.URL.Query().Get("grup_id"); groupID != "" {
			delete(groups[groupID], ws)
		}
		muGourp.Unlock()
		ws.Close()
	}()

	// Ambil user_id dan grup_id dari query params
	userID := r.URL.Query().Get("user_id")
	groupID := r.URL.Query().Get("grup_id")
	if userID == "" || groupID == "" {
		log.Println("user_id atau grup_id kosong!")
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Parameter tidak lengkap"))
		return
	}

	// Cek apakah user adalah anggota grup
	var userGroupMember models.GroupMember
	err = db.Where("group_id = ? AND user_id = ?", groupID, userID).First(&userGroupMember).Error
	if err != nil {
		log.Println("User bukan anggota grup:", err)
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Bukan anggota grup"))
		return
	}

	// Masukkan user ke daftar koneksi WebSocket
	muGourp.Lock()
	if groups[groupID] == nil {
		groups[groupID] = make(map[*websocket.Conn]string)
	}
	groups[groupID][ws] = userID
	muGourp.Unlock()

	log.Println("User", userID, "bergabung di grup", groupID)

	// Kirim riwayat chat grup ke user baru
	chatHistory, err := models.GetChatHistoryGrup(db, groupID)
	if err != nil {
		log.Println("Gagal mengambil riwayat chat:", err)
		ws.WriteJSON(map[string]string{"error": "Gagal mengambil riwayat chat"})
	} else {
		for _, msg := range chatHistory {
			if err := ws.WriteJSON(msg); err != nil {
				log.Println("Gagal mengirim pesan riwayat chat:", err)
				break
			}
		}
	}

	// Loop untuk terima pesan
	for {
		var msgGroup models.Message
		errMsGroup := ws.ReadJSON(&msgGroup)
		if errMsGroup != nil {
			log.Println("Kesalahan membaca pesan:", errMsGroup)
			break
		}
		msgGroup.SenderID = userID

		chatID, err := models.GetExistingChatIDByIdSender(db, msgGroup.SenderID)
		if err != nil {
			log.Println("Gagal mendapatkan chat ID:", err)
			return
		}

		if chatID == "" {
			// Chat belum ada, buat chat baru
			chatID, err = models.CreateNewChatGroup(db, groupID)
			if err != nil {
				log.Println("Gagal membuat chat baru:", err)
				return
			}
		}

		// Simpan pesan ke database
		errMsGroups := models.SaveMessageToDB(db, msgGroup.SenderID, chatID, "", msgGroup.Content)
		if errMsGroups != nil {
			log.Println("Gagal menyimpan pesan:", errMsGroups)
			continue
		}

		// Kirim pesan ke semua anggota grup
		broadcastGroup <- GroupMessage{Msg: msgGroup, Chat: models.Chat{GroupID: groupID}}
	}
}

// Fungsi untuk menangani broadcast pesan grup
func HandleMessagesGrupModel() {
	for {
		groupMsg := <-broadcastGroup

		msg := groupMsg.Msg
		groupID := groupMsg.Chat.GroupID // Ambil ID grup dari chat

		// Kirim pesan ke semua anggota grup
		if groupMembers, ok := groups[groupID]; ok {
			for client := range groupMembers {
				// Pastikan pesan tidak dikirim ke pengirimnya
				if groupMembers[client] != msg.SenderID {
					err := client.WriteJSON(msg)
					if err != nil {
						log.Println("Kesalahan mengirim pesan:", err)
						client.Close()               // Menutup koneksi yang gagal
						delete(groupMembers, client) // Hapus klien yang koneksinya error
					}
				}
			}
		}
	}
}

func main() {
	controller.ConnectDatabase()
	// WebSocket handler end to end

	// WebSocket handler group
	http.HandleFunc("/ws/group", func(w http.ResponseWriter, r *http.Request) {
		HandleConnectionsGrupController(controller.GetDB(), w, r)
	})

	go HandleMessagesGrupModel()

	fmt.Println("WebSocket server berjalan di :6070")
	log.Fatal(http.ListenAndServe(":6070", nil))
}
