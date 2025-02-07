package page

import (
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
	ParentID  *uint     `json:"parent_id"` // Ссылка на родительскую страницу (если есть)
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
	}
}

func getAllPages(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Список страниц"})
	}
}
