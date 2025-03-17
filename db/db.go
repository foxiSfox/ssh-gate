package db

import (
	"database/sql"
	"fmt"
	"log"

	"ssh-gate/models"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB инициализирует соединение с базой данных и создает необходимые таблицы
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения с базой данных: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Создаем таблицу пользователей
	if err := models.CreateUserTable(db); err != nil {
		log.Printf("Ошибка при создании таблицы пользователей: %v", err)
		return db, err
	}

	log.Println("База данных успешно инициализирована")
	return db, nil
}
