package server

import (
	"log"
	"velesbook/internal/auth"
	"velesbook/internal/page"
	"velesbook/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Run(db *gorm.DB, port string) {
	router := gin.Default()

	// Маршруты для аутентификации (без защиты)
	auth.RegisterRoutes(router, db)

	// Создаем защищенную группу маршрутов
	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware())

	// Подключаем защищенные маршруты (передаем protected, который является *gin.RouterGroup)
	user.RegisterRoutes(protected, db)
	page.RegisterRoutes(protected, db)

	log.Println("Сервер запущен на порту", port)
	router.Run(":" + port)
}
