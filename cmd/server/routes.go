package server

import (
	"velesbook/internal/auth"
	"velesbook/internal/page"
	"velesbook/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Маршруты для аутентификации (без защиты)
	auth.RegisterRoutes(router, db) // Маршруты для аутентификации

	// Группа защищенных маршрутов (требует токен)
	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware()) // Добавляем middleware
	user.RegisterRoutes(protected, db)   // Маршруты для работы с пользователями
	page.RegisterRoutes(protected, db)   // Маршруты для страниц
}
