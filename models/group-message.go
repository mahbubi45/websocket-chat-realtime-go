package models

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Group struct {
	ID        string        `gorm:"type:char(36);primaryKey"`
	Name      string        `gorm:"type:varchar(255);not null"`
	Members   []GroupMember `gorm:"foreignKey:GroupID"` // Relasi ke member
	CreatedAt time.Time     `gorm:"autoCreateTime"`
}

type GroupMember struct {
	ID       string    `gorm:"type:char(36);primaryKey"`
	GroupID  string    `gorm:"type:char(36);not null"`
	UserID   string    `gorm:"type:char(36);not null"`
	JoinedAt time.Time `gorm:"autoCreateTime"`
	Group    Group     `gorm:"foreignKey:GroupID"`
}

func CreateNewChatGroup(db *gorm.DB, group_id string) (string, error) {
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
