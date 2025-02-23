package models

import (
	"gorm.io/gorm"
)

// User menyimpan data user
type User struct {
	ID       string `gorm:"primaryKey;type:char(36)" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
}

func GetChatHistory(db *gorm.DB, chatID string) ([]Message, error) {
	var messages []Message

	// Ambil semua pesan yang ada di chat ini (baik yang dikirim maupun diterima)
	err := db.Where("chat_id = ?", chatID).
		Order("created_at ASC").Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}
