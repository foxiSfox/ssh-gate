package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"

	"hello-world/db"
	"hello-world/handlers"
)

const dbPath = "./users.db"

func main() {
	// Инициализация базы данных
	database, err := db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer database.Close()

	// Создание роутера
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Обработчики пользователей
	userHandler := handlers.NewUserHandler(database)

	// Маршруты
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Hello, World!"}`))
	})

	// API для пользователей
	r.Route("/api/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser) // Создание пользователя
		r.Get("/", userHandler.GetAllUsers) // Получение всех пользователей
		r.Get("/{id}", userHandler.GetUser) // Получение пользователя по ID
	})

	// Запуск сервера
	port := ":8080"
	fmt.Printf("Сервер запущен на http://localhost%s\n", port)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
		os.Exit(1)
	}
}
