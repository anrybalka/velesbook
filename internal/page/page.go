package page

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Структура страницы
type Page struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content"`
	UserID    uint      `json:"user_id"`   // Привязка к пользователю
	ParentID  *uint     `json:"parent_id"` // Используем указатель, чтобы поддерживать NULL
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Для GORM-связей:
	Parent   *Page  `json:"parent" gorm:"foreignKey:ParentID"`
	Children []Page `json:"children" gorm:"foreignKey:ParentID"`
}

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	page := router.Group("/pages")
	{
		page.GET("/", getAllPages(db))
		page.POST("/create", сreatePage(db))
	}
}

// POST /pages/create
//
//	{
//	    "title": "Название страницы",
//	    "content": "Текст страницы",
//	    "parent_id": null
//	}
func сreatePage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Title    string `json:"title" binding:"required"`
			Content  string `json:"content"`
			ParentID *uint  `json:"parent_id"` // Может быть nil, если родителя нет
		}
		// Парсим JSON-запрос
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
			return
		}

		// Получаем ID текущего пользователя из контекста
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизованный пользователь"})
			return
		}

		// Приводим userID к uint
		userIDUint, ok := userID.(uint)
		if !ok {
			userIDFloat, ok := userID.(float64) // JWT может передавать числа как float64
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка приведения userID"})
				return
			}
			userIDUint = uint(userIDFloat)
		}

		// Создаем новую страницу
		page := Page{
			Title:    input.Title,
			Content:  input.Content,
			UserID:   userIDUint,
			ParentID: input.ParentID, // Может быть nil, если родителя нет
		}

		// Сохраняем страницу в базе данных
		if err := db.Create(&page).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать страницу"})
			return
		}

		log.Printf("Создал страницу пользователь с ID: %v", userID)
		// Возвращаем успешный ответ
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Пользователь с ID %d создал страницу", userIDUint),
			"page":    page,
		})
	}
}

// GET /pages
func getAllPages(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pages []Page

		// Получаем всех пользователей из базы
		if err := db.Find(&pages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список пользователей"})
			return
		}

		// Получаем ID текущего пользователя из контекста
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизованный пользователь"})
			return
		}

		log.Printf("Всех пользователей вывел пользователь с ID: %v", userID)

		// Возвращаем список пользователей
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Все страницы вывел пользователь с ID: %v", userID),
			"pages":   pages,
		})
	}
}
