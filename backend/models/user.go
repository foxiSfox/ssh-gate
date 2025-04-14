package models

import (
	"database/sql"
	"fmt"
)

// User представляет модель пользователя
type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}

// CreateUserTable создает таблицу пользователей, если она не существует
func CreateUserTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		public_key TEXT NOT NULL
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы пользователей: %w", err)
	}

	return nil
}

// AddUser добавляет нового пользователя в базу данных
func AddUser(db *sql.DB, user User) (int64, error) {
	query := `
	INSERT INTO users (username, public_key)
	VALUES (?, ?);
	`

	result, err := db.Exec(query, user.Username, user.PublicKey)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления пользователя: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID: %w", err)
	}

	return id, nil
}

// GetUserByID получает пользователя по ID
func GetUserByID(db *sql.DB, id int64) (*User, error) {
	user := &User{}
	query := `
	SELECT id, username, public_key
	FROM users
	WHERE id = ?;
	`

	err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.PublicKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь с ID %d не найден", id)
		}
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return user, nil
}

// GetAllUsers получает всех пользователей
func GetAllUsers(db *sql.DB) ([]User, error) {
	query := `
	SELECT id, username, public_key
	FROM users;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователей: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.PublicKey); err != nil {
			return nil, fmt.Errorf("ошибка чтения данных пользователя: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе строк: %w", err)
	}

	return users, nil
}

// DeleteUser удаление пользователя
func DeleteUser(db *sql.DB, id int64) error {
	query := `
	DELETE FROM users
	WHERE id = ?;
	`

	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления пользователя: %w", err)
	}

	return nil
}
