package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"websocket-chat/helper"
	"websocket-chat/models"

	"github.com/gorilla/websocket"
)

// controller selected user yang akan di masukkan ke group
func (s *Server) UpdateSelectedUserToGrupController(w http.ResponseWriter, r *http.Request) {
	selectIdUser := helper.SelectedIdUser{}

	if err := json.NewDecoder(r.Body).Decode(&selectIdUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Cek apakah id_user kosong
	if len(selectIdUser.IDUser) == 0 {
		http.Error(w, "id_user tidak boleh kosong", http.StatusBadRequest)
		return
	}

	errSelectedUser := models.UpdateStatusSelectedUserToGroup(s.DB, selectIdUser.IDUser)

	if errSelectedUser != nil {
		http.Error(w, "error ya", http.StatusBadRequest)
		return
	}

	// Response sukses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Berhasil memperbarui status pengguna",
	})

}

// tambahkan ke grup yang sudah di selected true
func (s *Server) AddGrupMemberUsersController(w http.ResponseWriter, r *http.Request) {
	group := models.Group{}
	var idGroup string

	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if group.ID != "" {
		idGroup = group.ID
		fmt.Println("ini menambahkan user dengan grup yang sudah di buat", idGroup)
	} else if group.Name != "" {
		idGroup = models.AddGroup(s.DB, group.Name)
		fmt.Println("ini id group yang baru saja dibuat: ", idGroup)
	}

	// Tambahkan user ke grup yang sudah dipastikan ada
	err := models.CreatedUserSelectedToGrupmember(s.DB, idGroup)
	if err != nil {
		http.Error(w, "Gagal menambahkan user ke grup", http.StatusBadRequest)
		return
	}

	// Response sukses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Berhasil membuat grup",
	})

}

// WebSocket upgraderGroup
var upgraderGroup = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Struktur data gabungan
var muGourp sync.Mutex

// Channel yang menampung GroupMessage
var broadcastGroup = make(chan models.GroupMessage)
var groups = make(map[string]map[*websocket.Conn]string)

// Channel untuk broadcast pesan grup
func (s *Server) HandleConnectionsGrupController(w http.ResponseWriter, r *http.Request) {
	ws, err := upgraderGroup.Upgrade(w, r, nil)

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

	var groupID *string

	// Ambil nilai dari query parameter
	groupIDStr := r.URL.Query().Get("grup_id")

	// Ambil user_id dan grup_id dari query params
	userID := r.URL.Query().Get("user_id")

	// Jika kosong, biarkan nil, jika ada isi, gunakan nilai string-nya
	if groupIDStr != "" {
		groupID = &groupIDStr
	}

	if userID == "" || groupID == nil {
		log.Println("user_id atau grup_id kosong!")
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Parameter tidak lengkap"))
		return
	}

	// Cek apakah user adalah anggota grup
	var userGroupMember models.GroupMember
	err = s.DB.Where("group_id = ? AND user_id = ?", groupID, userID).
		First(&userGroupMember).Error
	if err != nil {
		log.Println("User bukan anggota grup:", err)
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Bukan anggota grup"))
		return
	}

	// Masukkan user ke daftar koneksi WebSocket
	muGourp.Lock()
	if groups[*groupID] == nil {
		groups[*groupID] = make(map[*websocket.Conn]string)
	}
	groups[*groupID][ws] = userID
	muGourp.Unlock()

	log.Println("User", userID, "bergabung di grup", groupID)

	// Kirim riwayat chat grup ke user baru
	chatHistory, err := models.GetChatHistoryGrup(s.DB, *groupID)
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

		chatID, err := models.GetExistingChatIDByIdSender(s.DB, msgGroup.SenderID)
		if err != nil {
			log.Println("Gagal mendapatkan chat ID:", err)
			return
		}

		if chatID == "" {
			// Chat belum ada, buat chat baru
			chatID, err = models.CreateNewChatGroup(s.DB, groupID)
			if err != nil {
				log.Println("Gagal membuat chat baru:", err)
				return
			}
		}

		// Simpan pesan ke database
		errMsGroups := models.SaveMessageToDB(s.DB, msgGroup.SenderID, chatID, "", msgGroup.Content)
		if errMsGroups != nil {
			log.Println("Gagal menyimpan pesan:", errMsGroups)
			continue
		}

		// Kirim pesan ke semua anggota grup
		broadcastGroup <- models.GroupMessage{Msg: msgGroup, Chat: models.Chat{GroupID: groupID}}
	}
}

// Fungsi untuk menangani broadcast pesan grup
func (s *Server) HandleMessagesGrupModel() {
	for {
		groupMsg := <-broadcastGroup

		msg := groupMsg.Msg
		groupID := groupMsg.Chat.GroupID // Ambil ID grup dari chat

		// Kirim pesan ke semua anggota grup
		if groupMembers, ok := groups[*groupID]; ok {
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
