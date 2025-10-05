## Уведомления: отложенная отправка

Сервис для создания, хранения и отложенной отправки уведомлений с использованием PostgreSQL, Redis и RabbitMQ. Имеет HTTP API и фонового воркера, а также пример консюмера.

### Архитектура
- `internal/app` — инициализация приложения: HTTP-сервер, миграции, подключения Redis/RabbitMQ, запуск воркера.
- `internal/handler` — HTTP-обработчики.
- `internal/usecase` — бизнес-логика (работа с БД и кэшем, выбор ближайшего уведомления).
- `internal/repository` — доступ к PostgreSQL.
- `internal/cache` — кэширование уведомлений в Redis.
- `internal/worker` — фоновые задачи: ожидание времени и публикация в RabbitMQ.
- `internal/rabbit` — обертка publisher над RabbitMQ.
- `internal/migrations` — создание таблиц и индексов.
- `cmd/delayed_notifier` — основной сервис (HTTP + воркер).
- `cmd/notify_consumer` — пример консюмера для очереди `notify_queue`.
- `internal/pkg/wbf/*` — утилиты.

### Зависимости (сторонние библиотеки)
- Gin (`github.com/gin-gonic/gin`) — HTTP сервер.
- Zerolog (`github.com/rs/zerolog`) — логирование.
- Go-Redis v8 (`github.com/go-redis/redis/v8`) — клиент Redis.
- amqp091-go (`github.com/rabbitmq/amqp091-go`) — клиент RabbitMQ.
- Viper (`github.com/spf13/viper`) — конфигурации (во внутренних утилитах).
- godotenv (`github.com/joho/godotenv`) — загрузка `.env`.
- lib/pq (`github.com/lib/pq`) — драйвер PostgreSQL.

Устанавливаются автоматически через `go mod download`.

### Переменные окружения (.env)
- `DB_NAME` — имя БД (например, `postgres`).
- `DB_USER` — пользователь БД.
- `DB_PASSWORD` — пароль БД.
- `DB_HOST` — хост БД (docker: `postgres`, локально: `localhost`).
- `DB_PORT` — порт БД (обычно `5432`).
- `TIMEZONE` — часовой пояс (например, `Europe/Moscow`).
- `REDIS_ADDR` — адрес Redis (docker: `redis:6379`, локально: `localhost:6379`).
- `RABBIT_URL` — `amqp://user:pass@host:5672/` (docker: `amqp://guest:guest@rabbitmq:5672/`).

### HTTP API
Базовый адрес: `http://localhost:8081`

- POST `/notify`
  - Создать уведомление.
  - Тело JSON:
    ```json
    {
      "user_id": "123",
      "title": "Заголовок",
      "body": "Текст",
      "status": false,
      "sending_date": "2025-10-05T12:00:00Z"
    }
    ```
  - Ответ `201 Created` с данными уведомления (включая `notify_id`).

- GET `/notify/:notifyID`
  - Проверить статус уведомления.
  - Ответы: `Notification sent` | `Notification not sent` | `Notification not found`.

- DELETE `/notify/:notifyID`
  - Удалить уведомление (из БД и кэша), разбудить воркер.
  - Ответ `200 OK`.

- GET `/notifications/:userID`
  - Получить все уведомления пользователя.

### Очереди
- Паблишер публикует сообщения в очередь `notify_queue` (RabbitMQ).
- Пример консюмера — `cmd/notify_consumer` — читает и печатает сообщения из `notify_queue`.

### Запуск

Docker Compose (все зависимости):
```bash
docker-compose build --no-cache
docker-compose up -d
docker logs -f <notification_service_container>
```