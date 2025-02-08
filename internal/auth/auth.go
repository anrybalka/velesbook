package auth

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Структура пользователя
type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// Секретный ключ для подписи JWT (можно вынести в .env)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func RegisterRoutes(group *gin.RouterGroup, db *sql.DB) {
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
func registerUser(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		var user User

		// Парсим JSON-запрос в структуру пользователя
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Printf("❌ auth.registerUser.ShouldBindJSON: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
			return
		}

		// Проверяем, существует ли уже пользователь с таким email
		var existingUser User
		err := db.QueryRow("SELECT id, email, password FROM users WHERE email = $1", user.Email).Scan(&existingUser.ID, &existingUser.Email, &existingUser.Password)
		if err == nil {
			log.Printf("❌ auth.registerUser.db.QueryRow: Пользователь с таким email уже существует: %v", user.Email)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким email уже существует"})
			return
		}

		// Хэшируем пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("❌ auth.registerUser.bcrypt.GenerateFromPassword: Ошибка при хэшировании пароля: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хэшировании пароля"})
			return
		}

		// Сохраняем пользователя в базе данных
		_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, string(hashedPassword))
		if err != nil {
			log.Printf("❌ auth.registerUser.db.Exec: Не удалось сохранить пользователя: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить пользователя"})
			return
		}

		// Генерируем JWT-токен
		token, err := generateToken(user.ID, user.Email)
		if err != nil {
			log.Printf("❌ auth.registerUser.generateToken: Ошибка при создании токена: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании токена"})
			return
		}

		log.Printf("✅ Регистрация пользователя с ID: %v", user.ID)
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
func loginUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		// Парсим JSON-запрос
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			log.Printf("❌ auth.loginUser.ShouldBindJSON: Неверный формат данных: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
			return
		}

		// Ищем пользователя в базе по email
		var user User
		err := db.QueryRow("SELECT id, email, password FROM users WHERE email = $1", loginRequest.Email).Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			log.Printf("❌ auth.loginUser.db.QueryRow: Неверный пароль: %v", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
			return
		}

		// Проверяем пароль
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
			log.Printf("❌ auth.loginUser.bcrypt.CompareHashAndPassword: Неверный пароль: %v", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
			return
		}

		// Генерируем JWT-токен
		token, err := generateToken(user.ID, user.Email)
		if err != nil {
			log.Printf("❌ auth.loginUser.generateToken: Ошибка при создании токена: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании токена"})
			return
		}

		log.Printf("✅ Авторизация пользователя с ID: %v", user.ID)
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
