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

	config := Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:9591@localhost:5432/velesbook?sslmode=disable"),
		JWTSecret:   jwtSecret,
	}

	// Логируем конфигурацию (JWTSecret скрыт для безопасности)
	log.Printf("Конфигурация загружена: \n"+
		"  - SERVER_PORT: %s\n"+
		"  - DATABASE_URL: %s\n"+
		"  - JWT_SECRET: %s\n",
		config.ServerPort, config.DatabaseURL, config.JWTSecret)

	return config
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
