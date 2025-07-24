package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/agusbasari29/GoAuraBill/config" // Ganti dengan path modul Anda
	"github.com/agusbasari29/GoAuraBill/internal/model" // Ganti dengan path modul Anda
)

var DB *gorm.DB

func main() {
	// 1. Muat Konfigurasi
	cfg, err := config.LoadConfig(".") // Memuat dari direktori saat ini
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// 2. Hubungkan ke Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database connection established")

	log.Println("Running Migrations...")
	err = DB.AutoMigrate(
		&model.User{},
		&model.Router{},
		&model.ServiceProfile{},
		&model.Voucher{},
		&model.Transaction{},
	)

	// 3. Siapkan Server Gin
	router := gin.Default()

	// Rute pengujian sederhana
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 4. Jalankan Server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}