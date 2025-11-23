# Распределённый grep

## Запуск
### 1. Запуск воркера (worker)

Воркер обрабатывает запросы от лидера. Запустите в отдельном терминале:

```bash
go run ./cmd -Mode=worker -Addr=localhost:8082
```

Или несколько воркеров:
```bash
# Терминал 1
go run ./cmd -Mode=worker -Addr=localhost:8082

# Терминал 2
go run ./cmd -Mode=worker -Addr=localhost:8083
```

### 2. Запуск лидера (leader)

Лидер читает файл, делит его на части и распределяет между воркерами:

```bash
go run ./cmd -Mode=leader -Peers="localhost:8082" -Quorum=1 "pattern" file.txt
```

Если несколько воркеров
```bash
go run ./cmd -Mode=leader -Peers="localhost:8082,localhost:8083" -Quorum=2 "pattern" file.txt
```

**Параметры:**
- `-Mode=leader` - режим лидера (по умолчанию)
- `-Peers="localhost:8082,localhost:8083"` - адреса воркеров через запятую
- `-Quorum=2` - минимальное количество успешных ответов
- `pattern` - шаблон для поиска
- `file.txt` - файл для поиска (или `-` для stdin)

**Grep флаги:**
- `-n` - выводить номера строк
- `-c` - только количество совпадений
- `-i` - игнорировать регистр
- `-v` - инвертировать (строки БЕЗ паттерна)
- `-A N` - вывести N строк после совпадения
- `-B N` - вывести N строк до совпадения
- `-C N` - вывести N строк контекста

## Примеры

### Базовый поиск

```bash
# Воркер
go run ./cmd -Mode=worker -Addr=localhost:8082

# Лидер
go run ./cmd -Mode=leader -Peers="localhost:8082" -Quorum=1 "pattern" file.txt
```

**Вывод:**
```
pattern matching line
another line with pattern
```

### С номерами строк

```bash
go run ./cmd -Mode=leader -Peers="localhost:8082" -Quorum=1 -n  "pattern" file.txt    
```

**Вывод:**
```
3:pattern matching line
4:another line with pattern
```

### Несколько воркеров

```bash
# Воркер 1
go run ./cmd -Mode=worker -Addr=localhost:8082

# Воркер 2
go run ./cmd -Mode=worker -Addr=localhost:8083

# Лидер (требует кворум 2)
go run ./cmd -Mode=leader -Peers="localhost:8082,localhost:8083" -Quorum=2 -n -i  "pattern" file.txt
```

**Вывод:**
```
3:pattern matching line
4:another line with pattern
6:Pattern in different case
```

## Структура проекта

```
cmd/
  └── cli_utility.go      # Точка входа
internal/
  ├── app/                # Запуск сервера/воркера
  ├── service/            # Бизнес-логика распределённого поиска
  ├── handler/            # HTTP обработчики
  ├── search/             # Логика поиска (grep флаги)
  ├── cli/                # Парсинг флагов
  └── entity/             # Структуры данных
```

