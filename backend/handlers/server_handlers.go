package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"ssh-gate/models"
	"ssh-gate/ssh"

	"github.com/go-chi/chi/v5"
)

// ServerHandler содержит обработчики для API серверов
type ServerHandler struct {
	DB *sql.DB
}

// NewServerHandler создает новый экземпляр ServerHandler
func NewServerHandler(db *sql.DB) *ServerHandler {
	return &ServerHandler{DB: db}
}

// CreateServer обрабатывает запрос на создание нового сервера
func (h *ServerHandler) CreateServer(w http.ResponseWriter, r *http.Request) {
	var server models.Server
	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		http.Error(w, "Ошибка при разборе запроса: "+err.Error(), http.StatusBadRequest)
		return
	}

	if server.IP == "" || server.Login == "" || server.Password == "" {
		http.Error(w, "IP, логин и пароль обязательны", http.StatusBadRequest)
		return
	}

	if server.Port == 0 {
		server.Port = 22
	}

	id, err := models.AddServer(h.DB, server)
	if err != nil {
		http.Error(w, "Ошибка при добавлении сервера: "+err.Error(), http.StatusInternalServerError)
		return
	}

	server.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(server)
}

// GetServer обрабатывает запрос на получение сервера по ID
func (h *ServerHandler) GetServer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	server, err := models.GetServerByID(h.DB, id)
	if err != nil {
		http.Error(w, "Сервер не найден: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(server)
}

// GetAllServers обрабатывает запрос на получение всех серверов
func (h *ServerHandler) GetAllServers(w http.ResponseWriter, r *http.Request) {
	servers, err := models.GetAllServers(h.DB)
	if err != nil {
		http.Error(w, "Ошибка при получении серверов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servers)
}

// UpdateServer обрабатывает запрос на обновление сервера
func (h *ServerHandler) UpdateServer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	var server models.Server
	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		http.Error(w, "Ошибка при разборе запроса: "+err.Error(), http.StatusBadRequest)
		return
	}

	if server.IP == "" || server.Login == "" || server.Password == "" {
		http.Error(w, "IP, логин и пароль обязательны", http.StatusBadRequest)
		return
	}

	if server.Port == 0 {
		server.Port = 22
	}

	server.ID = id
	if err := models.UpdateServer(h.DB, server); err != nil {
		http.Error(w, "Ошибка при обновлении сервера: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(server)
}

// AssignServerToUser обрабатывает запрос на привязку сервера к пользователю
func (h *ServerHandler) AssignServerToUser(w http.ResponseWriter, r *http.Request) {

	userIDStr := chi.URLParam(r, "userId")
	serverIDStr := chi.URLParam(r, "serverId")

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID пользователя", http.StatusBadRequest)
		return
	}

	serverID, err := strconv.ParseInt(serverIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID сервера", http.StatusBadRequest)
		return
	}

	// Получаем информацию о пользователе
	user, err := models.GetUserByID(h.DB, userID)
	if err != nil {
		http.Error(w, "Пользователь не найден: "+err.Error(), http.StatusNotFound)
		return
	}

	// Получаем информацию о сервере
	server, err := models.GetServerByID(h.DB, serverID)
	if err != nil {
		http.Error(w, "Сервер не найден: "+err.Error(), http.StatusNotFound)
		return
	}

	// Проверяем корректность публичного ключа
	if err := ssh.ValidatePublicKey(user.PublicKey); err != nil {
		http.Error(w, "Неверный формат публичного ключа: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Создаем конфигурацию для SSH-подключения
	sshConfig := ssh.SSHConfig{
		Host:     server.IP,
		Port:     server.Port,
		User:     server.Login,
		Password: server.Password,
	}

	// Добавляем публичный ключ на сервер, к которому надо получить доступ пользователю
	if err := ssh.AddAuthorizedKey(sshConfig, user.PublicKey); err != nil {
		http.Error(w, "Ошибка при добавлении ключа на сервер: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Привязываем сервер к пользователю в базе данных
	err = models.AssignServerToUser(h.DB, userID, serverID)
	if err != nil {
		// Если не удалось привязать сервер к пользователю, удаляем ключ с сервера
		_ = ssh.RemoveAuthorizedKey(sshConfig, user.PublicKey)
		http.Error(w, "Ошибка при привязке сервера: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetUserServers обрабатывает запрос на получение всех серверов пользователя
func (h *ServerHandler) GetUserServers(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID пользователя", http.StatusBadRequest)
		return
	}

	servers, err := models.GetUserServers(h.DB, userID)
	if err != nil {
		http.Error(w, "Ошибка при получении серверов пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servers)
}

// RemoveServerFromUser обрабатывает запрос на удаление привязки сервера к пользователю
func (h *ServerHandler) RemoveServerFromUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	serverIDStr := chi.URLParam(r, "serverId")

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID пользователя", http.StatusBadRequest)
		return
	}

	serverID, err := strconv.ParseInt(serverIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID сервера", http.StatusBadRequest)
		return
	}

	// Получаем информацию о пользователе
	user, err := models.GetUserByID(h.DB, userID)
	if err != nil {
		http.Error(w, "Пользователь не найден: "+err.Error(), http.StatusNotFound)
		return
	}

	// Получаем информацию о сервере
	server, err := models.GetServerByID(h.DB, serverID)
	if err != nil {
		http.Error(w, "Сервер не найден: "+err.Error(), http.StatusNotFound)
		return
	}

	// Создаем конфигурацию для SSH-подключения
	sshConfig := ssh.SSHConfig{
		Host:     server.IP,
		Port:     server.Port,
		User:     server.Login,
		Password: server.Password,
	}

	// Удаляем публичный ключ с сервера
	if err := ssh.RemoveAuthorizedKey(sshConfig, user.PublicKey); err != nil {
		http.Error(w, "Ошибка при удалении ключа с сервера: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Удаляем привязку сервера к пользователю в базе данных
	err = models.RemoveServerFromUser(h.DB, userID, serverID)
	if err != nil {
		http.Error(w, "Ошибка при удалении привязки сервера: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteServer обрабатывает запрос на удаление сервера по ID
func (h *ServerHandler) DeleteServer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	// Получаем информацию о сервере
	server, err := models.GetServerByID(h.DB, id)
	if err != nil {
		http.Error(w, "Сервер не найден: "+err.Error(), http.StatusNotFound)
		return
	}

	// Получаем пользователей, имеющих доступ к серверу
	users, err := models.GetServerUsers(h.DB, id)
	if err != nil {
		http.Error(w, "Ошибка при получении пользователей сервера: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Конфигурация для SSH-подключения к удаляемому серверу
	sshConfig := ssh.SSHConfig{
		Host:     server.IP,
		Port:     server.Port,
		User:     server.Login,
		Password: server.Password,
	}

	// Отзываем ключи у всех пользователей
	for _, user := range users {
		if err := ssh.RemoveAuthorizedKey(sshConfig, user.PublicKey); err != nil {
			http.Error(w, "Ошибка при удалении ключа с сервера: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Удаляем записи из user_servers после успешного отзыва ключей
	if err := models.RemoveAllUsersFromServer(h.DB, id); err != nil {
		http.Error(w, "Ошибка при удалении привязок пользователей: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Удаляем сам сервер
	if err := models.DeleteServer(h.DB, id); err != nil {
		http.Error(w, "Ошибка при удалении сервера: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
