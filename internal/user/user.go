package user

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"velesbook/pkg"

	"github.com/gin-gonic/gin"
)

// Структура пользователя
type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB) {
	user := router.Group("/users")
	{
		user.GET("/", getAllUsers(db))
	}
}

// GET /users
func getAllUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []User

		// Запрос для получения всех пользователей из базы
		rows, err := db.Query("SELECT id, email FROM users")
		if err != nil {
			log.Printf("❌ user.getAllUsers.db.Query error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список пользователей"})
			return
		}
		defer rows.Close()

		// Чтение всех пользователей из результата запроса
		for rows.Next() {
			var user User
			if err := rows.Scan(&user.ID, &user.Email); err != nil {
				log.Printf("❌ user.getAllUsers.rows.Next error: %v", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении данных пользователей"})
				return
			}
			users = append(users, user)
		}

		// Проверка на наличие ошибок после чтения строк
		if err := rows.Err(); err != nil {
			log.Printf("❌ user.getAllUsers.rows.Err error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке данных пользователей"})
			return
		}

		// Получаем userID через функцию
		userIDUint, err := pkg.GetUserID(c)
		if err != nil {
			log.Printf("❌ user.getAllUsers.GetUserID error: %v", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Логируем действие
		log.Printf("✅ Всех пользователей вывел пользователь с ID: %v", userIDUint)
		// Возвращаем список пользователей
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Всех пользователей вывел пользователь с ID: %v", userIDUint),
			"users":   users,
		})
	}
}
