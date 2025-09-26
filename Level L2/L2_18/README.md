###Сервис Календарь###

##Описание используемых библиотек##
1. Gin - web framework(запуск сервера)
2. Zap logger - логгер используемый вместо стандартного.

##Структура проекта##
1. cmd/calendar/main.go - точка входа в программу
2. internal/entity - структура передаваемых данных
3. internal/hanler - обработка ручек
4. internal/middleware - логирование запросов и время их обработки
5. internal/repository - база данных
6. internal/usecase - бизнес логика
7. tests/ - папка с тестами

##Запуск сервиса##
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
