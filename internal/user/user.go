package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Структура пользователя
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
}

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	user := router.Group("/users")
	{
		user.GET("/", getAllUsers(db))
	}
}

// GET /users
func getAllUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []User

		// Получаем всех пользователей из базы
		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список пользователей"})
			return
		}

		// Получаем ID текущего пользователя из контекста
		userID, exists := c.Get("userID")
		if !exists {
			userID = "неизвестный" // Если ID отсутствует
		}

		// Логируем действие
		log.Printf("Всех пользователей вывел пользователь с ID: %v", userID)

		// Возвращаем список пользователей
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Всех пользователей вывел пользователь с ID: %v", userID),
			"users":   users,
		})
	}
}
