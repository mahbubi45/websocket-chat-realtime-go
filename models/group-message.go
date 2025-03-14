package models

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupMessage struct {
	Msg  Message
	Chat Chat
}

type GroupMember struct {
	ID       string    `gorm:"type:char(36);primaryKey"`
	GroupID  string    `gorm:"type:char(36);not null"`
	UserID   string    `gorm:"type:char(36);not null"`
	JoinedAt time.Time `gorm:"autoCreateTime"`
	Group    Group     `gorm:"foreignKey:GroupID"`
}

// buat update user statusnya di select apa ngak untuk ke grup
func UpdateStatusSelectedUserToGroup(db *gorm.DB, idUser []string) error {
	var users User

	if len(idUser) == 0 {
		return fmt.Errorf("id user tidak boleh kosong")
	}

	return db.Model(&users).
		Where("id IN ?", idUser).
		Update("selected", true).Error

}

func UpdateStatusUnSelectedUserToGroup(db *gorm.DB) error {
	var users User

	return db.Model(&users).
		Where("selected = ? ", true).
		Update("selected", false).Error

}

// buat grup dari user yang sudah di selected
func CreatedUserSelectedToGrupmember(db *gorm.DB, idGroup string) error {
	// buat penampung user yang di selected
	var selecteduser []User
	// Ambil semua user yang selected = true
	if err := db.Where("selected = ?", true).
		Find(&selecteduser).Error; err != nil {
		log.Println("Gagal mengambil data user:", err)
		return err
	}

	if len(selecteduser) == 0 {
		return fmt.Errorf("tidak ada user yang dipilih")
	}

	// Looping setiap user dan masukkan ke GroupMember
	var groupMembers []GroupMember //penampung groupMember
	for _, users := range selecteduser {
		groupMembers = append(groupMembers, GroupMember{
			ID:       uuid.NewString(), // Generate UUID baru untuk setiap entri
			GroupID:  idGroup,
			UserID:   users.ID, // Pastikan ambil UserID dari loop
			JoinedAt: time.Now(),
		})
	}

	if err := db.Create(&groupMembers).Error; err != nil {
		log.Println("Gagal menambahkan user ke grup:", err)
		return err
	}

	//ketika sudah di tambahkan ke grup maka auto unselected lagi
	UpdateStatusUnSelectedUserToGroup(db)

	return nil
}

func GetExistingChatIDByIdSender(db *gorm.DB, senderID string) (string, error) {
	var chatID string

	// Cek apakah ada chat antara sender_id dan receiver_id, atau sebaliknya
	err := db.Raw(`
        SELECT chat_id FROM messages 
        WHERE (sender_id = ?) 
        LIMIT 1
    `, senderID).Scan(&chatID).Error

	if err == nil && chatID != "" {
		log.Println("Chat sudah ada dari messages, pakai ID:", chatID)
		return chatID, nil
	}

	// Kalau error bukan "record not found", return error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println("Gagal mencari chat:", err)
		return "", err
	}

	return "", nil
}

func CreateNewChatGroup(db *gorm.DB, group_id *string) (string, error) {
	newChat := Chat{
		ID:        uuid.NewString(),
		Name:      "chat-group",
		GroupID:   group_id,
		CreatedAt: time.Now(),
	}

	if err := db.Create(&newChat).Error; err != nil {
		log.Println("Gagal membuat chat:", err)
		return "", err
	}

	log.Println("Chat baru dibuat, ID:", newChat.ID)
	return newChat.ID, nil
}

// Fungsi untuk mengambil riwayat chat grup berdasarkan groupID
func GetChatHistoryGrup(db *gorm.DB, groupID string) ([]Message, error) {
	// Ambil semua chat yang terkait dengan groupID
	var chats []Chat
	if err := db.Where("group_id = ?", groupID).
		Find(&chats).Error; err != nil {
		return nil, err
	}

	var messages []Message
	// Ambil semua pesan terkait dengan ChatID dari setiap chat
	for _, chat := range chats {
		var msgs []Message
		if err := db.Where("chat_id = ?", chat.ID).
			Order("created_at ASC").
			Find(&msgs).Error; err != nil {
			return nil, err
		}
		messages = append(messages, msgs...)
	}
	return messages, nil
}
