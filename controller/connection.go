package controller

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "root:5431sabi@tcp(localhost:3306)/chat_db?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}

	// Auto-migrate tabel
	// database.AutoMigrate(
	// 	&models.User{},
	// 	&models.Chat{},
	// 	&models.ChatMember{},
	// 	&models.Message{})

	DB = database
	fmt.Println("Database Connected!")
}

func GetDB() *gorm.DB {
	return DB
}
