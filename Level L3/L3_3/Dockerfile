FROM golang:1.25.3-alpine

WORKDIR /app

# Копируем модули и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник из папки cmd/comment_tree
RUN go build -o main ./cmd/comment_tree

EXPOSE 8081

CMD ["./main"]
