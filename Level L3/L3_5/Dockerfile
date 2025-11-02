FROM golang:1.24.5-alpine

WORKDIR /app

# Копируем модули и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект в рабочую директорию
COPY . .

RUN go build -o main ./cmd/event_booker

EXPOSE 8081

CMD ["./main"]
