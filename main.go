package main

import (
	"log"
	"velesbook/config"
	"velesbook/internal/database"
	"velesbook/cmd/server"
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
