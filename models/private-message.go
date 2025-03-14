package models

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Chat menyimpan informasi chat (bisa pribadi atau grup)
type Chat struct {
	ID        string    `gorm:"type:char(36);primaryKey"`
	Name      string    `gorm:"null" json:"Name"`
	GroupID   *string   `gorm:"type:char(36);index"` // NULL jika chat pribadi
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Group     *Group    `gorm:"foreignKey:GroupID"` // Relasi opsional ke grup
	Messages  []Message `gorm:"foreignKey:ChatID"`  // Relasi ke pesan
}

// Message menyimpan pesan dalam chat
type Message struct {
	ID         string `gorm:"primaryKey;type:char(36)" json:"id"`
	ChatID     string `gorm:"not null" json:"chat_id"`
	SenderID   string `gorm:"not null" json:"sender_id"`
	ReceiverID string `gorm:"not null" json:"receiver_id"`
	Content    string `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time
}

// Message menyimpan pesan dalam chat
type MessageGrup struct {
	ID         string `gorm:"primaryKey;type:char(36)" json:"id"`
	ChatID     string `gorm:"not null" json:"chat_id"`
	SenderID   string `gorm:"not null" json:"sender_id"`
	UserSender User   `gorm:"foreignKey:SenderID" json:"userSender"`
	Content    string `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time
}

type MessageWithUserInfo struct {
	Message
	SenderID     string `json:"sender_id"`
	SenderName   string `json:"sender_name"`
	ReceiverID   string `json:"receiver_id"`
	ReceiverName string `json:"receiver_name"`
}

func GetExistingChatID(db *gorm.DB, senderID string, receiverID string) (string, error) {
	var chatID string

	// Cek apakah ada chat antara sender_id dan receiver_id, atau sebaliknya
	err := db.Raw(`
        SELECT chat_id FROM messages 
        WHERE (sender_id = ? AND receiver_id = ?) 
        OR (sender_id = ? AND receiver_id = ?) 
        LIMIT 1
    `, senderID, receiverID, receiverID, senderID).Scan(&chatID).Error

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

func CreateNewChat(db *gorm.DB) (string, error) {
	newChat := Chat{
		ID:        uuid.NewString(),
		Name:      "chat Private",
		GroupID:   nil,
		CreatedAt: time.Now(),
	}

	if err := db.Create(&newChat).Error; err != nil {
		log.Println("Gagal membuat chat:", err)
		return "", err
	}

	log.Println("Chat baru dibuat, ID:", newChat.ID)
	return newChat.ID, nil
}

// Simpan Pesan ke Database
func SaveMessageToDB(db *gorm.DB, senderID, chatID, receiverID, content string) error {
	// Cek dulu chat ID yang sudah ada atau buat baru jika perlu
	// Simpan pesan ke database
	message := Message{
		ID:         uuid.NewString(),
		ChatID:     chatID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  time.Now(),
	}

	if err := db.Create(&message).Error; err != nil {
		log.Println("Gagal menyimpan pesan ke database:", err)
		return err
	}

	log.Println("Pesan berhasil disimpan di chat:", chatID)
	return nil
}

// GetMessagesByUserID retrieves messages for a given user and includes sender/receiver information
func GetMessagesByUserID(db *gorm.DB, userID string) ([]MessageWithUserInfo, error) {
	var messages []MessageWithUserInfo

	// Perform the query to get messages where user is either sender or receiver
	err := db.Table("messages").
		Select("messages.*, sender.id  AS sender_name, receiver.id AS receiver_id").
		Joins("LEFT JOIN users AS sender ON sender.id = messages.sender_id").
		Joins("LEFT JOIN users AS receiver ON receiver.id = messages.receiver_id").
		Where("messages.sender_id = ? OR messages.receiver_id = ?", userID, userID).
		Order("messages.created_at ASC").
		Scan(&messages).Error

	return messages, err
}
