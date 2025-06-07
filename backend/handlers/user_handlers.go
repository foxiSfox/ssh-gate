package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"ssh-gate/models"
	"ssh-gate/ssh"

	"github.com/go-chi/chi/v5"
)

// UserHandler содержит обработчики для API пользователей
type UserHandler struct {
	DB *sql.DB
}

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// CreateUser обрабатывает запрос на создание нового пользователя
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Ошибка при разборе запроса: "+err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "Имя пользователя обязательно", http.StatusBadRequest)
		return
	}

	if user.PublicKey == "" {
		http.Error(w, "Публичный ключ обязателен", http.StatusBadRequest)
		return
	}

	id, err := models.AddUser(h.DB, user)
	if err != nil {
		http.Error(w, "Ошибка при добавлении пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Добавляем публичный ключ локально на jump сервер
	authorizedKeysFile := "authorized_keys"
	f, err := os.OpenFile(authorizedKeysFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Ошибка при открытии файла authorized_keys: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(user.PublicKey + "\n"); err != nil {
		http.Error(w, "Ошибка при записи ключа в файл authorized_keys: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser обрабатывает запрос на получение пользователя по ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByID(h.DB, id)
	if err != nil {
		http.Error(w, "Пользователь не найден: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetAllUsers обрабатывает запрос на получение всех пользователей
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAllUsers(h.DB)
	if err != nil {
		http.Error(w, "Ошибка при получении пользователей: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// UpdateUser обрабатывает запрос на обновление пользователя
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Ошибка при разборе запроса: "+err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "Имя пользователя обязательно", http.StatusBadRequest)
		return
	}

	if user.PublicKey == "" {
		http.Error(w, "Публичный ключ обязателен", http.StatusBadRequest)
		return
	}

	user.ID = id
	if err := models.UpdateUser(h.DB, user); err != nil {
		http.Error(w, "Ошибка при обновлении пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeleteUser обрабатывает запрос на удаление пользователя
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	/////// Удаление ключа. Начало
	user, err := models.GetUserByID(h.DB, id)
	if err != nil {
		http.Error(w, "Пользователь не найден: "+err.Error(), http.StatusNotFound)
		return
	}

	// Получаем серверы, к которым имеет доступ пользователь
	servers, err := models.GetUserServers(h.DB, id)
	if err != nil {
		http.Error(w, "Ошибка при получении серверов пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отзываем ключ с каждого сервера
	for _, server := range servers {
		sshConfig := ssh.SSHConfig{
			Host:     server.IP,
			Port:     server.Port,
			User:     server.Login,
			Password: server.Password,
		}

		_ = ssh.RemoveAuthorizedKey(sshConfig, user.PublicKey)
	}

	// Удаляем привязки серверов к пользователю в БД
	if err := models.RemoveAllServersFromUser(h.DB, id); err != nil {
		http.Error(w, "Ошибка при удалении привязок серверов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Читаем файл authorized_keys
	authorizedKeysFile := "authorized_keys"
	data, err := os.ReadFile(authorizedKeysFile)
	if err != nil {
		http.Error(w, "Ошибка при чтении файла authorized_keys: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Разделяем содержимое файла на строки и удаляем нужную строку
	lines := strings.Split(string(data), "\n")
	var updatedKeys []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" && trimmedLine != strings.TrimSpace(user.PublicKey) {
			updatedKeys = append(updatedKeys, line)
		}
	}

	// Перезаписываем файл authorized_keys
	err = os.WriteFile(authorizedKeysFile, []byte(strings.Join(updatedKeys, "\n")), 0644)
	if err != nil {
		http.Error(w, "Ошибка при записи файла authorized_keys: "+err.Error(), http.StatusInternalServerError)
		return
	}

	/////// Удаление ключа. Конец

	err = models.DeleteUser(h.DB, id)
	if err != nil {
		http.Error(w, "Ошибка при удалении пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
