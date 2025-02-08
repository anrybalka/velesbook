package server

import (
	"database/sql"
	"velesbook/internal/auth"
	"velesbook/internal/page"
	"velesbook/internal/user"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {

	// Группа для версии v1
	v1 := router.Group("/v1")

	// Группа аутентификации (без middleware)
	authGroup := v1.Group("/auth")
	auth.RegisterRoutes(authGroup, db) // Публичные маршруты

	// Группа защищенных маршрутов (требует токен)
	api := v1.Group("/api")
	api.Use(auth.AuthMiddleware()) // Требует аутентификации
	user.RegisterRoutes(api, db)   // Маршруты для работы с пользователями
	page.RegisterRoutes(api, db)   // Маршруты для страниц
}
