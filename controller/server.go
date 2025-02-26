package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Server struct {
	DB *gorm.DB
	R  *mux.Router
}

func (s *Server) ConnectDatabase() {
	dsn := "root:5431sabi@tcp(localhost:3306)/chat_db?charset=utf8mb4&parseTime=True&loc=Local"
	cndatabase, err := gorm.Open(mysql.Open(dsn),
		&gorm.Config{})
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
	s.R = mux.NewRouter()
	fmt.Println("Database Connected!")

	// Ping database
	sqldb, err := s.DB.DB()
	if err != nil {
		log.Fatal("Gagal mendapatkan koneksi database:", err)
	}

	if err := sqldb.Ping(); err != nil {
		log.Fatal("Database tidak merespons:", err)
	}

	fmt.Println("Database Connected and Ping Successful!")
}

func (s *Server) RunServer() {
	s.Routes()
	fmt.Println("WebSocket server berjalan di :6070")
	log.Fatal(http.ListenAndServe(":6070", s.R))
}
