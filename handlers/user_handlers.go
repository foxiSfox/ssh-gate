package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"hello-world/models"

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

	if user.Name == "" || user.PublicKey == "" {
		http.Error(w, "Имя и публичный ключ обязательны", http.StatusBadRequest)
		return
	}

	id, err := models.AddUser(h.DB, user)
	if err != nil {
		http.Error(w, "Ошибка при добавлении пользователя: "+err.Error(), http.StatusInternalServerError)
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
