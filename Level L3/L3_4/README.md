# 🖼️ Image Processor Service

Сервис для загрузки, обработки и управления изображениями с поддержкой resize, миниатюр и водяных знаков.

## 🚀 Возможности

- **Загрузка изображений** - поддержка JPG, PNG и тд.
- **Изменение размера** - resize с указанием ширины и высоты
- **Создание миниатюр** - автоматическое уменьшение до 500x500 пикселей
- **Водяные знаки** - добавление текстовых водяных знаков
- **Асинхронная обработка** - через Apache Kafka
- **Хранение файлов** - MinIO S3-совместимое хранилище
- **База данных** - PostgreSQL для метаданных
- **Web интерфейс** - современный фронтенд для управления

## 🛠️ Технологии

### Backend
- **Gin** - HTTP веб-фреймворк
- **Apache Kafka** - очередь сообщений для асинхронной обработки
- **PostgreSQL** - база данных для метаданных
- **MinIO** - S3-совместимое хранилище файлов
- **libvips** (через bimg) - обработка изображений

### Инфраструктура
- **Docker** - контейнеризация
- **Docker Compose** - оркестрация сервисов


## 📡 API Endpoints

### 1. Загрузка изображения
```http
POST /upload
Content-Type: multipart/form-data

Form Data:
- image: файл изображения (JPG, PNG, GIF)
```

**Пример ответа:**
```json
{
  "photo_id": "e7ee2e48-b860-42a5-a9a4-fcb16f8ac0f8"
}
```

### 2. Обработка изображения
```http
POST /process
Content-Type: application/json

# Для watermark:
{
  "photo_id": "e7ee2e48-b860-42a5-a9a4-fcb16f8ac0f8",
  "version": "watermark"
}

# Для resize/miniature (с query параметрами):
POST /process?widthPhoto=200&heightPhoto=200
Content-Type: application/json

{
  "photo_id": "e7ee2e48-b860-42a5-a9a4-fcb16f8ac0f8",
  "version": "resize",
  "width": "200",
  "height": "200"
}
```

**Пример ответа:**
```json
{
  "start_processing": "e7ee2e48-b860-42a5-a9a4-fcb16f8ac0f8"
}
```

### 3. Получение обработанного изображения
```http
GET /image/{photo_id}/{version}

# Примеры:
GET /image/e7ee2e48-b860-42a5-a9a4-fcb16f8ac0f8/watermark
GET /image/e7ee2e48-b860-42a5-a9a4-fcb16f8ac0f8/resize
GET /image/e7ee2e48-b860-42a5-a9a4-fcb16f8ac0f8/miniature
```

**Пример ответа:**
```json
{
  "img_url": "http://localhost:9000/photos/e7ee2e48-b860-42a5-a9a4-fcb16f8ac0f8_watermark.jpg"
}
```

### 4. Удаление изображения
```http
DELETE /image/delete
Content-Type: application/json

{
  "photo_id": "e7ee2e48-b860-42a5-a9a4-fcb16f8ac0f8",
  "version": "",
  "width": "",
  "height": "",
  "bucket_name": ""
}
```

**Пример ответа:**
```json
{
  "deleted_photo": "e7ee2e48-b860-42a5-a9a4-fcb16f8ac0f8"
}
```

### Быстрый старт

1. **Клонируйте репозиторий:**
```bash
git clone <repository-url>
cd L3_4
```

2. **Запустите сервисы через Docker Compose:**
```bash
docker-compose up -d
```

3. **Откройте веб-интерфейс:**
```
http://localhost:8081
```

### Разработка

1. **Установите зависимости:**
```bash
go mod download
```

2. **Настройте переменные окружения:**
```bash
cp .env.example .env
# Отредактируйте .env файл
```

3. **Запустите сервис:**
```bash
go run cmd/image_processor/main.go
```

## ⚙️ Конфигурация

### Переменные окружения

```bash
# База данных
DB_HOST=localhost
DB_PORT=5432
DB_NAME=image_processor
DB_USER=postgres
DB_PASSWORD=postgres

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
MINIO_PUBLIC_ENDPOINT=http://localhost:9000
BUCKET_NAME=photos

# Kafka
KAFKA_ADDR=kafka:9092
KAFKA_TOPIC=photos_topic
KAFKA_GROUP_ID=photo-consumer-group

# Настройки соединения
MAX_OPEN_CONNS=10
MAX_IDLE_CONNS=5
CONN_MAX_LIFETIME=30s
MAX_RETRIES=3
RETRY_DELAY=5s
```

## 📝 Примеры использования

### 1. Загрузка и обработка изображения

```bash
# Загрузка
curl -X POST -F "image=@photo.jpg" http://localhost:8081/upload

# Ответ: {"photo_id": "abc-123-def"}

# Обработка watermark
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"photo_id": "abc-123-def", "version": "watermark"}' \
  http://localhost:8081/process

# Получение результата
curl http://localhost:8081/image/abc-123-def/watermark
```

### 2. Resize изображения

```bash
# Обработка resize
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"photo_id": "abc-123-def", "version": "resize", "width": "800", "height": "600"}' \
  "http://localhost:8081/process?widthPhoto=800&heightPhoto=600"
```

### 3. Создание миниатюры

```bash
# Обработка miniature
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"photo_id": "abc-123-def", "version": "miniature", "width": "200", "height": "200"}' \
  "http://localhost:8081/process?widthPhoto=200&heightPhoto=200"
```

## 🔧 Разработка

### Структура проекта

```
L3_4/
├── cmd/image_processor/     # Точка входа
├── internal/                # Внутренние пакеты
│   ├── app/                # Инициализация приложения
│   ├── config/             # Конфигурация
│   ├── entity/             # Модели данных
│   ├── handler/            # HTTP обработчики
│   ├── kafka/              # Kafka интеграция
│   ├── repository/         # Слой данных
│   └── usecase/            # Бизнес-логика
├── migrations/             # Миграции БД
├── pkg/wbf/               # Внутренняя библиотека
├── web/                   # Frontend
└── docker-compose.yml     # Docker конфигурация
```
