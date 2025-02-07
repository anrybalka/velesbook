package server

import (
	"velesbook/internal/auth"
	"velesbook/internal/page"
	"velesbook/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Группируем маршруты
	auth.RegisterRoutes(router, db) // Маршруты для аутентификации
	user.RegisterRoutes(router, db) // Маршруты для работы с пользователями
	page.RegisterRoutes(router, db) // Маршруты для страниц
}
