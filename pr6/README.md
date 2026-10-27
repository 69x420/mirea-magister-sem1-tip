# Практическое занятие 6 - GORM: модели, миграции и связи

**Выполнил:** Зайцев А. А., группа ЭФМО-02-26.

REST-сервис заметок на GORM и chi. Модели `User`, `Note`, `Tag` описывают
связи 1:N и M:N, схема создаётся через `AutoMigrate`, связанные данные
подгружаются через `Preload`. В качестве СУБД используется SQLite (файл
`pz6_gorm.db`), поэтому отдельный сервер не нужен.

## Требования

- Go 1.21+
- `CGO_ENABLED=1` и компилятор C (драйвер `gorm.io/driver/sqlite`)
- консольный клиент `sqlite3` для просмотра схемы

Для PostgreSQL достаточно заменить драйвер на `gorm.io/driver/postgres` и
задать `DB_DSN`, например:

```bash
export DB_DSN='host=127.0.0.1 user=postgres password=postgres dbname=pz6_gorm port=5432 sslmode=disable'
```

## Запуск

```bash
cd ~/mirea/magister/semester1/tid/pr6
make build
DB_DSN="file:pz6_gorm.db?_journal_mode=WAL&_foreign_keys=1" PORT=8080 ./bin/server
make test
make schema     # схема users/notes/tags/note_tags
make clean      # удалить бинарник и файл базы
```

Переменные окружения:

- `DB_DSN` - строка подключения, по умолчанию
  `file:pz6_gorm.db?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=1`;
- `PORT` - порт HTTP-сервера, по умолчанию `8080`.

## Маршруты

| Метод и путь | Назначение | Коды |
| :--- | :--- | :--- |
| `GET /health` | проверка состояния | 200 |
| `POST /users` | создать пользователя | 201, 400, 409 |
| `GET /users` | список пользователей с `Preload("Notes")` | 200 |
| `POST /notes` | создать заметку с тегами | 201, 400, 404 |
| `GET /notes` | список заметок с автором и тегами | 200 |
| `GET /notes/{id}` | одна заметка с автором и тегами | 200, 400, 404 |

Ошибки возвращаются в виде `{"error":"..."}`.

## Модели

| Модель | Поля и теги | Связи |
| :--- | :--- | :--- |
| `User` | `Name size:100`, `Email size:200;uniqueIndex` | `Notes []Note` (1:N) |
| `Note` | `Title size:200`, `Content type:text`, `UserID not null;index` | `User User`, `Tags []Tag many2many:note_tags` |
| `Tag` | `Name size:50;uniqueIndex` | `Notes []Note many2many:note_tags` |

## Структура проекта

Дерево построено командой `eza -T --git-ignore`:

```text
pr6/
|-- cmd/server/main.go
|-- go.mod
|-- go.sum
|-- images/
|-- internal/db/db.go
|-- internal/httpapi/handlers.go
|-- internal/httpapi/handlers_test.go
|-- internal/httpapi/responses.go
|-- internal/httpapi/router.go
|-- internal/models/models.go
|-- Makefile
|-- pkg/middleware/cors.go
|-- pkg/middleware/logger.go
|-- README.md
|-- report.md
`-- requests.md
```

Примеры запросов собраны в `requests.md`.
