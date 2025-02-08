package page

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	"velesbook/pkg"

	"github.com/gin-gonic/gin"
)

// Структура страницы
type Page struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    uint      `json:"user_id"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB) {
	page := router.Group("/pages")
	{
		page.GET("/", getAllPages(db))
		page.GET("/my", getMyPages(db))
		page.POST("/create", сreatePage(db))
		// DELETE /delete/{id} // удалить страницу с id
		// GET /{id} // получить страницу с id
	}
}

func getMyPages(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pages []Page

		// Получаем userID через функцию
		userIDUint, err := pkg.GetUserID(c)
		if err != nil {
			log.Printf("Ошибка: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		log.Printf("🔍 Получаем страницы для userID: %d", userIDUint)

		// Запрос для получения страниц текущего пользователя
		rows, err := db.Query("SELECT id, title, content, user_id, parent_id, created_at, updated_at FROM pages WHERE user_id = ?", userIDUint)
		if err != nil {
			log.Printf("Ошибка: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список страниц текущего пользователя"})
			return
		}
		defer rows.Close()

		// Чтение всех страниц из результата запроса
		for rows.Next() {
			var page Page
			if err := rows.Scan(&page.ID, &page.Title, &page.Content, &page.UserID, &page.ParentID, &page.CreatedAt, &page.UpdatedAt); err != nil {
				log.Printf("Ошибка при чтении данных: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении данных страниц"})
				return
			}
			pages = append(pages, page)
		}

		// Проверка на ошибки при чтении строк
		if err := rows.Err(); err != nil {
			log.Printf("Ошибка при обработке данных: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке данных страниц"})
			return
		}

		log.Printf("Свои страницы вывел пользователь с ID: %v", userIDUint)
		// Возвращаем список страниц
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Свои страницы вывел пользователь с ID: %v", userIDUint),
			"pages":   pages,
		})
	}
}

// POST /pages/create
//
//	{
//	    "title": "Название страницы",
//	    "content": "Текст страницы",
//	    "parent_id": null
//	}
func сreatePage(db *sql.DB) gin.HandlerFunc {
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

		// Получаем userID через функцию
		userIDUint, err := pkg.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Запрос для создания страницы
		result, err := db.Exec("INSERT INTO pages (title, content, user_id, parent_id) VALUES (?, ?, ?, ?)", input.Title, input.Content, userIDUint, input.ParentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать страницу"})
			return
		}

		// Получаем ID новой страницы
		pageID, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить ID созданной страницы"})
			return
		}

		// Создаем структуру страницы
		page := Page{
			ID:       uint(pageID),
			Title:    input.Title,
			Content:  input.Content,
			UserID:   userIDUint,
			ParentID: input.ParentID, // Может быть nil, если родителя нет
		}

		log.Printf("Создал страницу пользователь с ID: %v", userIDUint)
		// Возвращаем успешный ответ
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Пользователь с ID %d создал страницу", userIDUint),
			"page":    page,
		})
	}
}

// GET /pages
func getAllPages(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pages []Page

		// Запрос для получения всех страниц
		rows, err := db.Query("SELECT id, title, content, user_id, parent_id, created_at, updated_at FROM pages")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список страниц"})
			return
		}
		defer rows.Close()

		// Чтение всех страниц из результата запроса
		for rows.Next() {
			var page Page
			if err := rows.Scan(&page.ID, &page.Title, &page.Content, &page.UserID, &page.ParentID, &page.CreatedAt, &page.UpdatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении данных страниц"})
				return
			}
			pages = append(pages, page)
		}

		// Проверка на ошибки при чтении строк
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке данных страниц"})
			return
		}

		// Получаем userID через функцию
		userIDUint, err := pkg.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Всех пользователей вывел пользователь с ID: %v", userIDUint)

		// Возвращаем список страниц
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Все страницы вывел пользователь с ID: %v", userIDUint),
			"pages":   pages,
		})
	}
}
