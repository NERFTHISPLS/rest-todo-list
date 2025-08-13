# REST API для Todo List

Простое REST API приложение для управления списком задач, построенное на Go с использованием Fiber и PostgreSQL.

## Описание

Это приложение предоставляет REST API для:

- Создания, чтения, обновления и удаления задач (CRUD операции)
- Управления списком задач с возможностью пометки их как выполненных
- Хранения данных в PostgreSQL базе данных
- Работы через HTTP API endpoints

## Стэк

- **Go 1.24.2**
- **Fiber v2** - для создания API
- **pgx** - для работы с базой данных
- **PostgreSQL** - реляционная база данных
- **Docker** - контейнеризация приложения

## Требования

- Go 1.24.2+
- Docker

## Переменные окружения

Создайте файл `.env` в корневой директории проекта со следующими переменными:

```env
ENV=dev

SERVER_PORT=8080
SERVER_TIMEOUT_READ=3s
SERVER_TIMEOUT_WRITE=5s
SERVER_TIMEOUT_IDLE=5s

DB_HOST=db
DB_PORT=5432
DB_USER=app_user
DB_PASSWORD=app_password
DB_NAME=app_name

```

ENV может также иметь значение ```prod```

## Запуск программы c использованием Docker

1. Убедитесь, что Docker установлен
2. Создайте файл `.env` с необходимыми переменными окружения
3. Запустите приложение:

```bash
docker compose up --build
```

Приложение будет доступно по адресу: `http://localhost:{порт_указанный_в_env}`

## API Endpoints

- `GET /tasks` - получить список всех задач
- `POST /tasks` - создать новую задачу
- `PUT /tasks/:id` - обновить задачу
- `DELETE /tasks/:id` - удалить задачу

## Структура проекта

```
rest-todo-list/
├── cmd/            # Основные файлы приложения
│   ├── main.go     # Точка входа
├── docs/           # Документация Swagger
├── internal/
│   ├── config/     # Конфигурация приложения
│   ├── database/   # База данных
│   ├── handlers/   # Обработчики HTTP запросов
│   ├── logger/     # Логирование
│   ├── models/     # Модели данных
│   ├── repository/ # Запросы к базе данных
│   ├── server/     # Сервер
├── .air.toml       # Конфигурации Air
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Разработка

Для разработки с автоматической перезагрузкой при изменении кода:

```bash
go install github.com/air-verse/air@latest
air -c .air.toml
```
