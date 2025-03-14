package models

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Group struct {
	ID        string        `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string        `gorm:"type:varchar(255);not null" json:"name"`
	Members   []GroupMember `gorm:"foreignKey:GroupID"` // Relasi ke member
	CreatedAt time.Time     `gorm:"autoCreateTime" json:"createdAt"`
}

func AddGroup(db *gorm.DB, name string) string {
	group := Group{
		ID:   uuid.NewString(), // Generate UUID baru
		Name: name,
	}

	if err := db.Create(&group).Error; err != nil {
		log.Println("Gagal menambahkan user ke grup:", err)
		return ""
	}

	return group.ID
}
