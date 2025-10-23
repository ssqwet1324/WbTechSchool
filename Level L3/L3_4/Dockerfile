FROM golang:1.24.5-alpine

# Устанавливаем зависимости для bimg
RUN apk add --no-cache build-base vips-dev

# Создаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник
RUN go build -o main ./cmd/image_processor

EXPOSE 8081

CMD ["./main"]
