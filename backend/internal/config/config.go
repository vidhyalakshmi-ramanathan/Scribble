package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBURL      string
	SECRET_KEY string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		Port:       os.Getenv("PORT"),
		DBURL:      os.Getenv("DBURL"),
		SECRET_KEY: os.Getenv("SECRET_KEY"),
	}
}
