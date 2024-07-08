package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	ExternalAPI string
	*Storage
	*Server
}

type Storage struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

type Server struct {
	Host    string
	Port    string
	Timeout time.Duration
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}

	timeout, err := strconv.Atoi(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		log.Panic("Error loading SERVER_TIMEOUT variable")
	}

	return &Config{
		os.Getenv("ENV"),
		os.Getenv("EXTERNAL_API_URL"),
		&Storage{
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Database: os.Getenv("DATABASE_NAME"),
		},
		&Server{
			Host:    os.Getenv("SERVER_HOST"),
			Port:    os.Getenv("SERVER_PORT"),
			Timeout: time.Duration(timeout),
		},
	}
}
