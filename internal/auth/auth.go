package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Структура пользователя
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password,omitempty" gorm:"not null"`
}

// Секретный ключ для подписи JWT (можно вынести в .env)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func RegisterRoutes(group *gin.RouterGroup, db *gorm.DB) {
	group.POST("/register", registerUser(db))
	group.POST("/login", loginUser(db))
}

// Генерация JWT-токена
func generateToken(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"id":    userID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Токен действителен 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
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

		// Генерируем JWT-токен
		token, err := generateToken(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании токена"})
			return
		}

		// Возвращаем успешный ответ
		c.JSON(http.StatusOK, gin.H{
			"message": "Пользователь успешно зарегистрирован",
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
			},
			"token": token,
		})
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

		// Генерируем JWT-токен
		token, err := generateToken(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании токена"})
			return
		}

		// Возвращаем успешный ответ
		c.JSON(http.StatusOK, gin.H{
			"message": "Вход выполнен успешно",
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
			},
			"token": token,
		})

	}
}
