FROM golang:1.25.3-alpine

WORKDIR /app

# Копируем модули и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект в рабочую директорию
COPY . .

RUN go build -o main ./cmd/sales_tracker

EXPOSE 8081

CMD ["./main"]
