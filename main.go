package main

import (
	"log"
	"velesbook/cmd/server"
	"velesbook/config"
	"velesbook/internal/database"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Подключаем базу данных
	db := database.InitDB(cfg.DatabaseURL)

	// Запускаем сервер
	server.Run(db, cfg.ServerPort)

	log.Println("Сервер завершил работу")
}
