package server

import (
	"velesbook/internal/auth"
	"velesbook/internal/page"
	"velesbook/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Группа для аутентификации (без защиты)
	authGroup := router.Group("/auth")
	auth.RegisterRoutes(authGroup, db) // Теперь маршруты аутентификации в /auth

	// Группа защищенных маршрутов (требует токен)
	api := router.Group("/api")
	api.Use(auth.AuthMiddleware()) // Добавляем middlewar
	user.RegisterRoutes(api, db)   // Маршруты для работы с пользователями
	page.RegisterRoutes(api, db)   // Маршруты для страниц
}
