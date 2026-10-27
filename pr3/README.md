# Практическое занятие 3 - HTTP API на net/http

**Выполнил:** Зайцев А. А., группа ЭФМО-02-26.

CRUD-сервис списка задач на стандартной библиотеке `net/http`, без сторонних
роутеров. Данные хранятся в памяти и защищены `sync.RWMutex`.

## Требования

- Go 1.22+ (нужны method-pattern'ы и `r.PathValue`)
- curl (или Postman) для проверки запросов

## Запуск

```bash
cd ~/mirea/magister/semester1/tid/pr3
make run                 # PORT=8080 по умолчанию
PORT=8081 make run       # другой порт
make build && ./bin/server
make test
```

## Маршруты

| Метод и путь | Назначение | Коды |
| :--- | :--- | :--- |
| `GET /health` | проверка состояния | 200 |
| `GET /tasks` | список задач, фильтр `?q=` | 200 |
| `POST /tasks` | создание по `{"title":"..."}` | 201, 400, 422 |
| `GET /tasks/{id}` | одна задача | 200, 400, 404 |
| `PATCH /tasks/{id}` | изменить `title` и/или `done` | 200, 400, 404, 422 |
| `DELETE /tasks/{id}` | удалить задачу | 204, 400, 404 |

Ответы и ошибки имеют вид `{"error":"..."}`, все ответы в формате JSON.
Заголовки CORS выставляются для всех маршрутов, `OPTIONS` отвечает `204`.

## Структура проекта

Дерево построено командой `eza -T --git-ignore`:

```text
pr3/
|-- cmd/server/main.go
|-- internal/api/handlers.go
|-- internal/api/handlers_test.go
|-- internal/api/middleware.go
|-- internal/api/responses.go
|-- internal/storage/memory.go
|-- Makefile
|-- requests.md
|-- images/
|-- go.mod
|-- README.md
`-- report.md
```

## Конфигурация

- `PORT` - порт HTTP-сервера, по умолчанию `8080`.
- Логи: метод, URI, код ответа, размер и длительность запроса.
- Остановка по `Ctrl+C` или `SIGTERM` выполняется через `srv.Shutdown` с
  таймаутом 5 секунд.

Примеры запросов собраны в `requests.md`.
