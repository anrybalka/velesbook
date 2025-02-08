package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pg"
)

// InitDB инициализирует подключение к базе данных PostgreSQL
func InitDB(databaseURL string) *sql.DB {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("❌ Ошибка подключения к базе данных:", err)
	}

	// Проверяем соединение
	if err = db.Ping(); err != nil {
		log.Fatal("❌ База данных недоступна:", err)
	}

	log.Println("✅ База данных подключена")
	return db
}
