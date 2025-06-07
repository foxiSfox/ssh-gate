# Этап 1: Сборка Frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Копируем и устанавливаем зависимости frontend
COPY frontend/package*.json ./
RUN npm install

# Копируем исходники и собираем frontend
COPY frontend/ .
RUN npm run build

# Этап 2: Сборка Backend
FROM golang:1.24.1-alpine AS backend-builder

# Установка зависимостей для работы с SQLite (если используется SQLite)
RUN apk add --no-cache gcc musl-dev sqlite sqlite-dev

WORKDIR /app/backend

# Копируем файлы backend и устанавливаем зависимости
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Копируем исходники backend
COPY backend/ .

# Компиляция backend
RUN go build -o main .

# Финальный этап: запуск обоих приложений
FROM alpine:latest

WORKDIR /app

# Установка SQLite runtime-библиотек (если используется SQLite)
RUN apk add --no-cache sqlite

# Копируем собранный фронтенд
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Копируем собранный бэкэнд
COPY --from=backend-builder /app/backend/main ./backend/main
COPY --from=backend-builder /app/backend/authorized_keys ./backend/authorized_keys
COPY --from=backend-builder /app/backend/id_rsa_jump_server ./backend/id_rsa_jump_server

# Указываем порт
EXPOSE 8080

# Запускаем backend
CMD ["sh", "-c", "cd backend && ./main"]