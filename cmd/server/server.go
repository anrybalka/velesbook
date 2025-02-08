package server

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

func Run(db *sql.DB, port string) {
	router := gin.Default()

	// Настройка маршрутов через отдельную функцию в routes.go
	SetupRoutes(router, db)

	log.Println("Сервер запущен на порту", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
