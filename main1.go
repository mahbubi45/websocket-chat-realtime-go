package main

import (
	"fmt"
	"log"
	"net/http"
	"websocket-chat/controller"
)

// Channel untuk broadcast pesan grup
// var broadcastGroup = make(chan models.MessageGrup)

// Handler WebSocket untuk grup
// Channel untuk broadcast pesan grup
// var (
// 	groups = make(map[string]map[*websocket.Conn]string) // 🔥 Simpan koneksi per grup
// )

// // Handler WebSocket untuk grup
// func handleConnectionsGrupController(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("WebSocket upgrade gagal:", err)
// 		return
// 	}
// 	defer ws.Close()

// 	// 🔥 Ambil user_id dan grup_id dari query params
// 	userID := r.URL.Query().Get("user_id")
// 	groupID := r.URL.Query().Get("grup_id") // Sesuai request lo
// 	if userID == "" || groupID == "" {
// 		log.Println("❌ user_id atau grup_id kosong!")
// 		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Parameter tidak lengkap"))
// 		return
// 	}

// 	// 🔥 Cek apakah user adalah anggota grup
// 	var chatMember models.ChatMember
// 	err = db.Where("chat_id = ? AND user_id = ?", groupID, userID).First(&chatMember).Error
// 	if err != nil {
// 		log.Println("❌ User bukan anggota grup:", err)
// 		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Bukan anggota grup"))
// 		return
// 	}

// 	// 🔥 Masukkan user ke daftar koneksi WebSocket
// 	mu.Lock()
// 	if groups[groupID] == nil {
// 		groups[groupID] = make(map[*websocket.Conn]string)
// 	}
// 	groups[groupID][ws] = userID
// 	mu.Unlock()

// 	log.Println("✅ User", userID, "bergabung di grup", groupID)

// 	// 🔥 Kirim riwayat chat grup ke user baru
// 	chatHistory, err := models.GetChatHistoryGrup(db, groupID)
// 	if err != nil {
// 		log.Println("⚠️ Gagal mengambil riwayat chat:", err)
// 	} else {
// 		for _, msg := range chatHistory {
// 			ws.WriteJSON(msg)
// 		}
// 	}

// 	// Loop untuk terima pesan
// 	for {
// 		var msgGroup models.MessageGrup
// 		err := ws.ReadJSON(&msgGroup)
// 		if err != nil {
// 			log.Println("❌ Kesalahan membaca pesan:", err)
// 			break
// 		}

// 		// 🔥 Simpan pesan ke database
// 		savedMsg, err := models.SaveGroupMessage(db, msgGroup.ChatID, msgGroup.SenderID, msgGroup.Content)
// 		if err != nil {
// 			log.Println("❌ Gagal menyimpan pesan:", err)
// 			continue
// 		}

// 		// 🔥 Kirim ke channel broadcast
// 		broadcastGroup <- savedMsg
// 	}

// 	// Hapus user saat disconnect
// 	mu.Lock()
// 	delete(groups[groupID], ws)
// 	mu.Unlock()
// }

// func HandleMessagesGrupModel() {
// 	for {
// 		msg := <-broadcastGroup

// 		mu.Lock()
// 		// Pastikan grup ada sebelum broadcast
// 		if groupMembers, ok := groups[msg.ChatID]; ok {
// 			for client, userID := range groupMembers {
// 				if userID != msg.SenderID { // Kirim ke semua anggota kecuali pengirim
// 					err := client.WriteJSON(msg)
// 					if err != nil {
// 						log.Println("Kesalahan mengirim pesan:", err)
// 						client.Close()
// 						delete(groupMembers, client)
// 					}
// 				}
// 			}
// 		}
// 		mu.Unlock()
// 	}
// }

func main1() {
	controller.ConnectDatabase()
	// WebSocket handler end to end

	// WebSocket handler group
	// http.HandleFunc("/ws-group", func(w http.ResponseWriter, r *http.Request) {
	// 	handleConnectionsGrupController(controller.GetDB(), w, r)
	// })

	// go HandleMessagesGrupModel()

	fmt.Println("WebSocket server berjalan di :7070")
	log.Fatal(http.ListenAndServe(":7070", nil))
}
