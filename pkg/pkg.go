package pkg

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetUserID извлекает и преобразует userID из контекста Gin
func GetUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("Неавторизованный пользователь")
	}

	// Преобразуем userID в uint (JWT может передавать как float64)
	switch v := userID.(type) {
	case uint:
		return v, nil
	case float64:
		return uint(v), nil
	default:
		return 0, fmt.Errorf("Ошибка приведения userID")
	}
}
