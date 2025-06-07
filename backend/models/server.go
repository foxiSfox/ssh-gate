package models

import (
	"database/sql"
	"fmt"
)

// Server представляет модель сервера
type Server struct {
	ID       int64  `json:"id"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// CreateServerTable создает таблицу серверов и связующую таблицу
func CreateServerTable(db *sql.DB) error {
	// Создаем таблицу серверов
	serverQuery := `
        CREATE TABLE IF NOT EXISTS servers (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                ip TEXT NOT NULL UNIQUE,
                port INTEGER NOT NULL DEFAULT 22,
                login TEXT NOT NULL,
                password TEXT NOT NULL
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
        INSERT INTO servers (ip, port, login, password)
        VALUES (?, ?, ?, ?);
        `

	result, err := db.Exec(query, server.IP, server.Port, server.Login, server.Password)
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
        SELECT id, ip, port, login, password
        FROM servers
        WHERE id = ?;
        `

	var server Server
	err := db.QueryRow(query, id).Scan(&server.ID, &server.IP, &server.Port, &server.Login, &server.Password)
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
        SELECT id, ip, port, login, password
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
		if err := rows.Scan(&server.ID, &server.IP, &server.Port, &server.Login, &server.Password); err != nil {
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
        SELECT s.id, s.ip, s.port, s.login, s.password
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
		if err := rows.Scan(&server.ID, &server.IP, &server.Port, &server.Login, &server.Password); err != nil {
			return nil, fmt.Errorf("ошибка чтения данных сервера: %w", err)
		}
		servers = append(servers, server)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе строк: %w", err)
	}

	return servers, nil
}

// GetServerUsers получает всех пользователей, имеющих доступ к серверу
func GetServerUsers(db *sql.DB, serverID int64) ([]User, error) {
	query := `
        SELECT u.id, u.username, u.public_key
        FROM users u
        JOIN user_servers us ON u.id = us.user_id
        WHERE us.server_id = ?;
        `

	rows, err := db.Query(query, serverID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователей сервера: %w", err)
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

// RemoveAllUsersFromServer удаляет все привязки пользователей к серверу
func RemoveAllUsersFromServer(db *sql.DB, serverID int64) error {
	query := `
        DELETE FROM user_servers
        WHERE server_id = ?;
        `

	if _, err := db.Exec(query, serverID); err != nil {
		return fmt.Errorf("ошибка удаления привязок пользователей к серверу: %w", err)
	}

	return nil
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

// DeleteServer удаляет сервер по ID
func DeleteServer(db *sql.DB, id int64) error {
	query := `
	DELETE FROM servers
	WHERE id = ?;
	`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления сервера: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества затронутых строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("сервер с ID %d не найден", id)
	}

	return nil
}
