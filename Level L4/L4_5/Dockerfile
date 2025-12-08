FROM golang:1.24.5-alpine

WORKDIR /app

# Копируем модули и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY api_optimization .

# Собираем бинарник из папки cmd
RUN go build -o main ./cmd/api_optimization

EXPOSE 8081

CMD ["./main"]
