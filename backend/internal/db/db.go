package db

import (
	"fmt"
	"log"

	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	// DB CONNECTION
	cfg := config.Load()
	db, err := gorm.Open(postgres.Open(cfg.DBURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected to database!")
	// db.AutoMigrate(&entities.User{}, &entities.Post{})
	DB = db
}
