package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware проверяет JWT-токен перед доступом к защищенным маршрутам
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		tokenString := c.GetHeader("Authorization")

		// Проверяем, передан ли токен и его формат (Bearer token)
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен отсутствует или некорректен"})
			c.Abort()
			return
		}

		// Убираем "Bearer " из строки токена
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Парсим и проверяем токен
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil // Используем секретный ключ для проверки
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			c.Abort()
			return
		}

		// Сохраняем userID и email в контексте Gin
		c.Set("userID", claims["id"])
		c.Set("email", claims["email"])

		// Передаем управление дальше
		c.Next()
	}
}
