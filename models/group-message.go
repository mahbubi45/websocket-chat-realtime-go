package models

import (
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

func SaveGroupMessage(db *gorm.DB, chatID, senderID, content string) (MessageGrup, error) {
	message := MessageGrup{
		ID:        uuid.NewString(),
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := db.Create(&message).Error; err != nil {
		return message, err
	}
	return message, nil
}
