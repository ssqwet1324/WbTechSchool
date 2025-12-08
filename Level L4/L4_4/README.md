# Mem GC exporter

Сервис для автоматического просмотра через метрики **prometheus** данных о **памяти** и **GC**.

## Технологии
- **Go 1.25.3** - язык программирования
- **Gin** - веб-фреймворк для HTTP API
- **Prometheus** - инструмент для метрик
- **Grafana** - инструмент для просмотра метрик в виде графиков
- **Pprof** - инструмент для профилирования

## Установка зависимостей
```bash
go mod download
```
После
```bash
go mod tidy
```

### Быстрый старт с Docker
1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd L4_4
```

2. Запустите проект:
```bash
docker-compose up -d
```

1. Сервис будет доступен по адресу: `http://localhost:8080`
Основная ручка /metrics: `http://localhost:8080/metrics`
_Данная ручка автоматически вызывается каждые 5s для обновления актуальных данных, изменить можно в файле `/pkg/prometheus/prometheus.yml`_

2. Графики можно посмотреть по адресу `http://localhost:3000`. Пароль и логин те, которые указаны в `docker-compose.yml`(admin, admin)
Заходим в Dashboards, выбираем `Create Dashboard`, далее `Add Visualisation`, справа выбираем графана, далее `Configure a new data source`,
тут выбираем `prometheus`, вписываем `http://prometheus:9090`(или ваш хост для prometheus), сохраняем.
Далее переходим снова в Dashboards, выбираем тоже самое и там сразу вылезет prometheus, выбираем его, выбираем метрику с припиской mem
например: `mem_gc_exporter_allocations`, дале ставим `Label filters` type = runtime и нажимаем Run Queries.

3. В Prometheus можно перейти по адресу `http://localhost:9090`

4. Остановить проект:
```bash
docker-compose down
```

## Структура проекта
1. cmd/mem_gc_exporter/main.go - точка входа в программу
2. internal/entity - структура передаваемых данных
3. internal/hanler - обработка ручек
4. internal/middleware - логирование запросов и время их обработки
5. internal/usecase - бизнес логика
6. pkg/metrics - метрики для prometheus


