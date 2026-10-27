# Практическое занятие 4 - CRUD-сервис на chi

**Выполнил:** Зайцев А. А., группа ЭФМО-02-26.

REST-сервис "Список задач" на роутере `github.com/go-chi/chi/v5`. Данные
хранятся в памяти и защищены `sync.RWMutex`.

## Требования

- Go 1.22+
- curl (или Postman) для проверки запросов

## Запуск

```bash
cd ~/mirea/magister/semester1/tid/pr4
make run                 # PORT=8080 по умолчанию
PORT=8081 make run       # другой порт
make build && ./bin/server
make test
```

## Маршруты

| Метод и путь | Назначение | Коды |
| :--- | :--- | :--- |
| `GET /health` | проверка состояния | 200 |
| `GET /api/tasks` | список задач, `?q=`, `?done=`, `?page=`, `?limit=` | 200, 400 |
| `POST /api/tasks` | создание по `{"title":"..."}` | 201, 400, 422 |
| `GET /api/tasks/{id}` | одна задача | 200, 400, 404 |
| `PUT /api/tasks/{id}` | заменить `title` и/или `done` | 200, 400, 404, 422 |
| `DELETE /api/tasks/{id}` | удалить задачу | 204, 400, 404 |

Список возвращается в конверте `{"items":[...],"total":N,"page":P,"limit":L}`;
значения по умолчанию - `page=1`, `limit=10`, максимум `limit=100`. Ошибки
имеют вид `{"error":"..."}`, все ответы в формате JSON. Заголовки CORS
выставляются для всех маршрутов, `OPTIONS` отвечает `204`.

## Структура проекта

Дерево построено командой `eza -T --git-ignore`:

```text
pr4/
|-- go.mod
|-- go.sum
|-- images/
|-- internal/task/handler.go
|-- internal/task/handler_test.go
|-- internal/task/model.go
|-- internal/task/repo.go
|-- main.go
|-- Makefile
|-- pkg/middleware/cors.go
|-- pkg/middleware/logger.go
|-- README.md
|-- report.md
|-- requests.md
```

## Конфигурация

- `PORT` - порт HTTP-сервера, по умолчанию `8080`.
- Логи: метод, URI, код ответа, размер и длительность запроса.
- Остановка по `Ctrl+C` или `SIGTERM` выполняется через `srv.Shutdown` с
  таймаутом 5 секунд.

Примеры запросов собраны в `requests.md`.
