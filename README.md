# SSH Gate API

REST API для управления пользователями и их доступом к серверам через SSH.

## API Endpoints

### Пользователи

#### Создание пользователя
```
POST /api/users
Content-Type: application/json

{
    "username": "user1",
    "public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC6eNtGpNGwstc...."
}
```

#### Получение пользователя по ID
```
GET /api/users/{id}
```

#### Получение всех пользователей
```
GET /api/users
```

### Серверы

#### Создание сервера
```
POST /api/servers
Content-Type: application/json

{
    "ip": "192.168.1.100"
}
```

#### Получение сервера по ID
```
GET /api/servers/{id}
```

#### Получение всех серверов
```
GET /api/servers
```

### Управление доступом пользователей к серверам

#### Привязка сервера к пользователю
```
POST /api/users/{userId}/servers/{serverId}
```

При этом публичный ключ пользователя будет автоматически добавлен в файл `authorized_keys` на сервере.

#### Получение всех серверов пользователя
```
GET /api/users/{userId}/servers
```

#### Удаление привязки сервера к пользователю
```
DELETE /api/users/{userId}/servers/{serverId}
```

При этом публичный ключ пользователя будет автоматически удален из файла `authorized_keys` на сервере.

## Запуск сервера

```bash
go run main.go
```

Сервер будет доступен по адресу `http://localhost:8080`. 