package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"ssh-gate/models"

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

	if server.IP == "" {
		http.Error(w, "IP адрес обязателен", http.StatusBadRequest)
		return
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

	err = models.AssignServerToUser(h.DB, userID, serverID)
	if err != nil {
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

	err = models.RemoveServerFromUser(h.DB, userID, serverID)
	if err != nil {
		http.Error(w, "Ошибка при удалении привязки сервера: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
