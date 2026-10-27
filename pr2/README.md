# Практическое занятие 2 - Структура Go-проекта

**Выполнил:** Зайцев А. А., группа ЭФМО-02-26.

Минимальный сервис, демонстрирующий рекомендованную структуру Go-проекта:
тонкая точка входа в `cmd/`, приватный код в `internal/`, общие хелперы в
`utils/`.

## Требования

- Go 1.21+ (`go version`)
- Git
- curl (или браузер) для проверки ответов

## Запуск

```bash
# из каталога ~/mirea/magister/semester1/tid/pr2
go run ./cmd/myapp                  # слушает :8080
PORT=8081 go run ./cmd/myapp        # порт из переменной окружения
go build -o bin/myapp ./cmd/myapp
./bin/myapp
```

## Маршруты

| Метод и путь | Ответ |
| :--- | :--- |
| `GET /` | 200, текст `Hello, Go project structure!` |
| `GET /ping` | 200, JSON `{"status":"ok","time":"<RFC3339>"}` |
| `GET /fail` | 400, JSON `{"error":"bad_request_example"}` |
| любой другой путь | 404 |

Все ответы содержат заголовок `X-Request-Id`: он берётся из запроса, а если
его нет - генерируется.

## Структура проекта

Дерево построено командой `eza -T --git-ignore`:

```text
pr2/
|-- cmd/myapp/main.go
|-- internal/app/app.go
|-- internal/app/app_test.go
|-- internal/app/handlers/ping.go
|-- utils/httpjson.go
|-- utils/logger.go
|-- images/
|-- go.mod
|-- README.md
`-- report.md
```

## Конфигурация

- `PORT` - порт HTTP-сервера, по умолчанию `8080`.
- Логи выводятся в stdout: время, адрес клиента, метод и путь (`utils.LogRequest`).
- Пакет `internal/` недоступен для импорта извне модуля по правилам Go.
