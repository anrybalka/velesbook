package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Структура пользователя
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password,omitempty" gorm:"not null"`
}

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", registerUser(db))
		auth.POST("/login", loginUser(db))
	}
}

// POST /auth/register
//
//	{
//	  "email": "user@example.com",
//	  "password": "securepassword"
//	}
func registerUser(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		var user User

		// Парсим JSON-запрос в структуру пользователя
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
			return
		}

		// Проверяем, существует ли уже пользователь с таким email
		var existingUser User
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким email уже существует"})
			return
		}

		// Хэшируем пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хэшировании пароля"})
			return
		}

		// Сохраняем пользователя в базе данных
		user.Password = string(hashedPassword)
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить пользователя"})
			return
		}

		// Возвращаем успешный ответ без пароля
		c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно зарегистрирован", "user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		}})
	}
}

// POST /auth/login
//
//	{
//	  "email": "user@example.com",
//	  "password": "securepassword"
//	}
func loginUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		// Парсим JSON-запрос
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
			return
		}

		// Ищем пользователя в базе по email
		var user User
		if err := db.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
			return
		}

		// Проверяем пароль
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
			return
		}

		// Успешный вход (здесь можно добавить JWT-токен в будущем)
		c.JSON(http.StatusOK, gin.H{"message": "Вход выполнен успешно", "user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		}})
	}
}
