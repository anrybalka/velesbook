package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	JWTSecret   string
}

func LoadConfig() Config {

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("Ошибка: переменная окружения JWT_SECRET не установлена")
	}

	return Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:9591@localhost:5432/velesbook?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "velesbook2025"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
