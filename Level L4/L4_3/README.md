# Сервис Календарь

## Описание используемых библиотек
1. Gin - web framework(запуск сервера)

## Структура проекта
1. cmd/calendar/main.go - точка входа в программу
2. internal/entity - структура передаваемых данных
3. internal/hanler - обработка ручек
4. internal/middleware - логирование запросов и время их обработки
5. internal/repository - база данных
6. internal/usecase - бизнес логика
7. internal/logger - бизнес логика
8. tests/ - папка с тестами

## Запуск сервиса
1. Переходим в папку cmd/calendar
2. Запускаем main файл: `go run main.go`
3. Отправляем запросы на ручки:
   POST /create_event — создание нового события;
   POST /update_event — обновление существующего;
   POST /delete_event — удаление;
   GET /events_for_day — получить все события на день;
   GET /events_for_week — события на неделю;
   GET /events_for_month — события на месяц.
4. Запусть тест бизнес логики: переходим в папку /tests и прописываем `go test -v` или `go test`.

## Примеры запросов и ответов

### POST /create_event
**Создание нового события**

**Запрос:**
```bash
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "name_event": "Встреча с командой",
    "data_event": "2025-12-25T14:00:00Z",
    "text": "Обсуждение планов на следующий квартал",
    "remind_at": "2025-12-25T13:30:00Z"
  }'
```

**Ответ (успех):**
```json
{
  "Created Event success": {
    "user_id": "user123",
    "name_event": "Встреча с командой",
    "data_event": "2025-12-25T14:00:00Z",
    "text": "Обсуждение планов на следующий квартал",
    "remind_at": "2025-12-25T13:30:00Z"
  }
}
```

### POST /update_event
**Обновление существующего события**

**Запрос:**
```bash
curl -X POST http://localhost:8080/update_event \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "name_event": "Встреча с командой (обновлено)",
    "data_event": "2025-12-25T15:00:00Z",
    "text": "Обсуждение планов на следующий квартал - изменено время",
    "remind_at": "2025-12-25T14:30:00Z"
  }'
```

**Ответ (успех):**
```json
{
  "Updated Event success": {
    "user_id": "user123",
    "name_event": "Встреча с командой (обновлено)",
    "data_event": "2025-12-25T15:00:00Z",
    "text": "Обсуждение планов на следующий квартал - изменено время",
    "remind_at": "2025-12-25T14:30:00Z"
  }
}
```

### POST /delete_event
**Удаление события**

**Запрос:**
```bash
curl -X POST http://localhost:8080/delete_event \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "name_event": "Встреча с командой",
    "data_event": "2025-12-25T14:00:00Z",
    "text": "",
    "remind_at": "2025-12-25T13:30:00Z"
  }'
```

**Ответ (успех):**
```json
{
  "Deleted Event success": "2024-12-25T14:00:00Z"
}
```

### GET /events_for_day/:user_id/:date
**Получение событий на день**

**Запрос:**
```bash
curl -X GET http://localhost:8080/events_for_day/user123/2025-12-25
```

**Ответ (успех):**
```json
{
  "event": [
    {
      "user_id": "user123",
      "name_event": "Встреча с командой",
      "data_event": "2025-12-25T14:00:00Z",
      "text": "Обсуждение планов на следующий квартал",
      "remind_at": "2025-12-25T13:30:00Z"
    }
  ]
}
```


### GET /events_for_week/:user_id/:date
**Получение событий на неделю**

**Запрос:**
```bash
curl -X GET http://localhost:8080/events_for_week/user123/2024-12-25
```

**Ответ (успех):**
```json
{
  "events": [
    {
      "user_id": "user123",
      "name_event": "Встреча с командой",
      "data_event": "2025-12-25T14:00:00Z",
      "text": "Обсуждение планов на следующий квартал",
      "remind_at": "2025-12-25T13:30:00Z"
    },
    {
      "user_id": "user123",
      "name_event": "Презентация проекта",
      "data_event": "2025-12-27T10:00:00Z",
      "text": "Демонстрация нового функционала",
      "remind_at": "2025-12-27T09:30:00Z"
    }
  ]
}
```

### GET /events_for_month/:user_id/:date
**Получение событий на месяц**

**Запрос:**
```bash
curl -X GET http://localhost:8080/events_for_month/user123/2024-12-25
```

**Ответ (успех):**
```json
{
  "event": [
    {
      "user_id": "user123",
      "name_event": "Встреча с командой",
      "data_event": "2025-12-25T14:00:00Z",
      "text": "Обсуждение планов на следующий квартал",
      "remind_at": "2025-12-25T13:30:00Z"
    },
    {
      "user_id": "user123",
      "name_event": "Новогодний корпоратив",
      "data_event": "2025-12-31T18:00:00Z",
      "text": "Празднование Нового года",
      "remind_at": "2025-12-31T17:00:00Z"
    }
  ]
}
```
