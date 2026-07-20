# Corporate Messenger - `GopherTalk`

Backend корпоративного мессенджера на Go.

Проект разрабатывается как pet-project с целью изучения архитектуры backend-приложений, REST API, WebSocket, PostgreSQL, Docker и принципов построения многослойных сервисов.

---

## Что реализовано

✔ Graceful Shutdown

✔ JWT

✔ slog

✔ PostgreSQL

✔ Connection Pool

✔ Многослойная архитектура

✔ Middleware

✔ Валидация

✔ Конфигурация

✔ Контексты запросов

---

## Что планируется реализовать

Планируется:

- WebSocket
- диалоги
- групповые чаты
- отправка сообщений
- редактирование сообщений
- удаление сообщений
- Docker Compose
- Swagger
- unit-тесты
- интеграционные тесты
- миграции БД

---

## Технологии

- Go 1.25.6
- PostgreSQL
- JWT
- bcrypt
- Docker
- slog
- Docker Compose

---

## Архитектура проекта

Client

    │

HTTP Request 

    │

Middleware

    │

Handlers

    │

Services

    │

Database

    │

PostgreSQL

---

## REST API

| Метод | URL           | Описание                  |
| ----- | ------------- | ------------------------- |
| POST  | /api/register | Регистрация               |
| POST  | /api/login    | Авторизация               |
| GET   | /api/me       | Информация о пользователе |

---

## Переменные окружения

| Переменная  | Описание        |
| ----------- | --------------- |
| DB_HOST     | PostgreSQL host |
| DB_PORT     | PostgreSQL port |
| DB_USER     | пользователь БД |
| DB_PASSWORD | пароль БД       |
| DB_NAME     | имя БД          |
| JWT_SECRET  | секрет          |
| SERVER_PORT | порт            |

---

## Запуск проекта

git clone [...](https://github.com/amat-andre/corporate_messenger)

cd corporate_messenger

docker compose up

go run ./cmd

---

## Структура проекта

- `cmd`
  - `main.go` — точка входа приложения
- `internal`
  - `auth` — JWT и пароль
  - `validation` - валидация входных данных
  - `db` — подключение к PostgreSQL и управление пулом соединений
  - `config` — загрузка и хранение конфигурации
  - `middlewate` — HTTP middleware
  - `handlers` — HTTP обработчики
  - `helpers` — общие HTTP-хелперы
  - `server` — создание и запуск HTTP-сервера
  - `service` - бизнес-логика
- `tests` — 
- `web` - 

- `go.mod` - зависимости проекта
- `go.sum` - lock-файл зависимостей
- `Dockerfile` — сборка docker-образа
- `README.md` — документация по проекту

---
