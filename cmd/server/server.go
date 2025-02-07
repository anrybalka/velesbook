package server

import (
	"log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"velesbook/internal/auth"
	"velesbook/internal/user"
	"velesbook/internal/page"
)

func Run(db *gorm.DB, port string) {
	router := gin.Default()

	// Подключаем маршруты
	auth.RegisterRoutes(router, db)
	user.RegisterRoutes(router, db)
	page.RegisterRoutes(router, db)

	log.Println("Сервер запущен на порту", port)
	router.Run(":" + port)
}
