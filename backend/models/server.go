package models

import (
	"database/sql"
	"fmt"
)

// Server представляет модель сервера
type Server struct {
	ID int64  `json:"id"`
	IP string `json:"ip"`
}

// CreateServerTable создает таблицу серверов и связующую таблицу
func CreateServerTable(db *sql.DB) error {
	// Создаем таблицу серверов
	serverQuery := `
	CREATE TABLE IF NOT EXISTS servers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL UNIQUE
	);
	`

	// Создаем связующую таблицу
	userServerQuery := `
	CREATE TABLE IF NOT EXISTS user_servers (
		user_id INTEGER NOT NULL,
		server_id INTEGER NOT NULL,
		PRIMARY KEY (user_id, server_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE
	);
	`

	// Создаем таблицу серверов
	if _, err := db.Exec(serverQuery); err != nil {
		return fmt.Errorf("ошибка создания таблицы серверов: %w", err)
	}

	// Создаем связующую таблицу
	if _, err := db.Exec(userServerQuery); err != nil {
		return fmt.Errorf("ошибка создания связующей таблицы: %w", err)
	}

	return nil
}

// AddServer добавляет новый сервер в базу данных
func AddServer(db *sql.DB, server Server) (int64, error) {
	query := `
	INSERT INTO servers (ip)
	VALUES (?);
	`

	result, err := db.Exec(query, server.IP)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления сервера: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID: %w", err)
	}

	return id, nil
}

// GetServerByID получает сервер по ID
func GetServerByID(db *sql.DB, id int64) (Server, error) {
	query := `
	SELECT id, ip
	FROM servers
	WHERE id = ?;
	`

	var server Server
	err := db.QueryRow(query, id).Scan(&server.ID, &server.IP)
	if err != nil {
		if err == sql.ErrNoRows {
			return Server{}, fmt.Errorf("сервер с ID %d не найден", id)
		}
		return Server{}, fmt.Errorf("ошибка получения сервера: %w", err)
	}

	return server, nil
}

// GetAllServers получает все серверы
func GetAllServers(db *sql.DB) ([]Server, error) {
	query := `
	SELECT id, ip
	FROM servers;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения серверов: %w", err)
	}
	defer rows.Close()

	var servers []Server
	for rows.Next() {
		var server Server
		if err := rows.Scan(&server.ID, &server.IP); err != nil {
			return nil, fmt.Errorf("ошибка чтения данных сервера: %w", err)
		}
		servers = append(servers, server)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе строк: %w", err)
	}

	return servers, nil
}

// AssignServerToUser привязывает сервер к пользователю
func AssignServerToUser(db *sql.DB, userID, serverID int64) error {
	query := `
	INSERT INTO user_servers (user_id, server_id)
	VALUES (?, ?);
	`

	_, err := db.Exec(query, userID, serverID)
	if err != nil {
		return fmt.Errorf("ошибка привязки сервера к пользователю: %w", err)
	}

	return nil
}

// GetUserServers получает все серверы пользователя
func GetUserServers(db *sql.DB, userID int64) ([]Server, error) {
	query := `
	SELECT s.id, s.ip
	FROM servers s
	JOIN user_servers us ON s.id = us.server_id
	WHERE us.user_id = ?;
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения серверов пользователя: %w", err)
	}
	defer rows.Close()

	var servers []Server
	for rows.Next() {
		var server Server
		if err := rows.Scan(&server.ID, &server.IP); err != nil {
			return nil, fmt.Errorf("ошибка чтения данных сервера: %w", err)
		}
		servers = append(servers, server)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе строк: %w", err)
	}

	return servers, nil
}

// RemoveServerFromUser удаляет привязку сервера к пользователю
func RemoveServerFromUser(db *sql.DB, userID, serverID int64) error {
	query := `
	DELETE FROM user_servers
	WHERE user_id = ? AND server_id = ?;
	`

	result, err := db.Exec(query, userID, serverID)
	if err != nil {
		return fmt.Errorf("ошибка удаления привязки сервера: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества затронутых строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("привязка сервера к пользователю не найдена")
	}

	return nil
}
