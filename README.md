# Subscription Service

REST API сервис для агрегации данных об онлайн подписках пользователей.

## Стек

- **Go 1.25.6**
- **chi** — HTTP роутер
- **GORM** — ORM для работы с базой данных
- **PostgreSQL** — база данных
- **goose** — миграции
- **uber/fx** — dependency injection
- **logrus** — логирование
- **Docker Compose** — запуск сервиса

## Запуск

### Через Docker Compose
```bash
docker compose up --build
```

Сервис будет доступен на `http://localhost:8040`

### Локально

1. Заполни `.env` файл
2. Подними PostgreSQL
3. Запусти сервис:
```bash
go run ./cmd/main.go
```

## Переменные окружения

| Переменная   | Описание             | Default       |
|--------------|----------------------|---------------|
| APP_PORT     | Порт сервиса         | 8040          |
| DB_HOST      | Хост базы данных     | localhost     |
| DB_PORT      | Порт базы данных     | 5432          |
| DB_USER      | Пользователь БД      | postgres      |
| DB_PASSWORD  | Пароль БД            | postgres      |
| DB_NAME      | Название БД          | subscriptions |
| DB_SSLMODE   | SSL режим            | disable       |

## API

| Метод  | Путь                          | Описание                      |
|--------|-------------------------------|-------------------------------|
| POST   | /api/v1/subscriptions         | Создать подписку              |
| GET    | /api/v1/subscriptions         | Получить список подписок      |
| GET    | /api/v1/subscriptions/{id}    | Получить подписку по ID       |
| PUT    | /api/v1/subscriptions/{id}    | Обновить подписку             |
| DELETE | /api/v1/subscriptions/{id}    | Удалить подписку              |
| GET    | /api/v1/subscriptions/total   | Суммарная стоимость за период |

### Создание подписки
```
POST /api/v1/subscriptions
```
```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025"
}
```

### Обновление подписки
```
PUT /api/v1/subscriptions/{id}
```
```json
{
  "price": 500,
  "end_date": "12-2025"
}
```

### Суммарная стоимость за период
```
GET /api/v1/subscriptions/total?from=01-2025&to=12-2025
GET /api/v1/subscriptions/total?from=01-2025&to=12-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba
GET /api/v1/subscriptions/total?from=01-2025&to=12-2025&service_name=Yandex Plus
```
```json
{
  "total": 4800
}
```

### Список подписок с фильтрацией
```
GET /api/v1/subscriptions
GET /api/v1/subscriptions?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba
GET /api/v1/subscriptions?service_name=Yandex Plus
```

## Структура проекта
```
.
├── cmd/
│   └── main.go
├── internal/
│   ├── api/
│   │   ├── dto.go
│   │   ├── handler.go
│   │   ├── middleware.go
│   │   └── router.go
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── postgres.go
│   ├── di/
│   │   └── app.go
│   ├── model/
│   │   └── subscription.go
│   ├── repository/
│   │   ├── subscriptionImpl.go
│   │   └── subscription.go
│   └── service/
│       ├── subscriptionImpl.go
│       └── subscription.go
├── migrations/
│   ├── 001_init.up.sql
│   └── 001_init.down.sql
├── docs/
│   └── swagger.yaml
├── .env
├── Dockerfile
├── docker-compose.yml
└── README.md
```