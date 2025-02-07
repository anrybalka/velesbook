package user

import (
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

// Функция регистрации нового пользователя
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	user := router.Group("/users")
	{
		user.GET("/", getAllUsers(db))
	}
}

func getAllUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []User

		// Получаем всех пользователей из базы
		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список пользователей"})
			return
		}

		// Возвращаем список пользователей
		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}
