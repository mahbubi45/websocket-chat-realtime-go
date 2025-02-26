package controller

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Server struct {
	DB *gorm.DB
}

func (s *Server) ConnectDatabase() {
	dsn := "root:5431sabi@tcp(localhost:3306)/chat_db?charset=utf8mb4&parseTime=True&loc=Local"
	cndatabase, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}

	// Auto-migrate tabel
	// database.AutoMigrate(
	// 	&models.User{},
	// 	&models.Chat{},
	// 	&models.ChatMember{},
	// 	&models.Message{})

	s.DB = cndatabase
	fmt.Println("Database Connected!")
}

func (s *Server) GetDB() *gorm.DB {
	return s.DB
}
