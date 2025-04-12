package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"ssh-gate/db"
	"ssh-gate/handlers"
)

func main() {
	// Инициализируем базу данных

	database, err := db.InitDB("users.db")
	if err != nil {
		log.Fatal("Ошибка инициализации базы данных:", err)
	}
	defer database.Close()

	// Получаем путь к приватному ключу из переменной окружения
	os.Setenv("SSH_PRIVATE_KEY_PATH", "id_rsa_jump_server")
	keyPath := os.Getenv("SSH_PRIVATE_KEY_PATH")
	if keyPath == "" {
		log.Fatal("Не указан путь к приватному ключу (SSH_PRIVATE_KEY_PATH)")
	}

	// Создаем обработчики
	userHandler := handlers.NewUserHandler(database)
	serverHandler := handlers.NewServerHandler(database, keyPath)

	// Создаем роутер
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll().Handler)
	// Определяем маршруты
	r.Route("/api", func(r chi.Router) {
		// Маршруты для пользователей
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/", userHandler.GetAllUsers)
			r.Get("/{id}", userHandler.GetUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})

		// Маршруты для серверов
		r.Route("/servers", func(r chi.Router) {
			r.Post("/", serverHandler.CreateServer)
			r.Get("/", serverHandler.GetAllServers)
			r.Get("/{id}", serverHandler.GetServer)
			r.Delete("/{id}", serverHandler.DeleteServer)
		})

		// Маршруты для управления доступом пользователей к серверам
		r.Route("/users/{userId}/servers", func(r chi.Router) {
			r.Get("/", serverHandler.GetUserServers)
			r.Post("/{serverId}", serverHandler.AssignServerToUser)
			r.Delete("/{serverId}", serverHandler.RemoveServerFromUser)
		})
	})

	// Обслуживаем статические файлы фронтенда
	fileServer(r, "/", http.Dir("./frontend/dist"))

	// Запускаем сервер
	log.Println("Сервер запущен на порту :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

// Функция для удобного обслуживания фронтенда и SPA-роутинга (todo временно)
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	fs := http.StripPrefix(path, http.FileServer(root))
	r.Get(path+"*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := root.Open(r.URL.Path); os.IsNotExist(err) {
			http.ServeFile(w, r, "./frontend/dist/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	})
}
