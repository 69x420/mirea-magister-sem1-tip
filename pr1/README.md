# Практическое занятие 1 - Установка и настройка окружения Go

**Выполнил:** Зайцев А. А., группа ЭФМО-02-26.

Минимальный HTTP-сервис на стандартной библиотеке `net/http` с внешней
зависимостью для генерации UUID.

## Требования

- Go 1.21+ (`go version`)
- Git
- curl (или браузер) для проверки ответов

## Запуск

```bash
# из каталога ~/mirea/magister/semester1/tid/pr1
go mod download
go run ./cmd/server                  # слушает :8080
APP_PORT=8081 go run ./cmd/server    # порт из переменной окружения
```

Сборка бинарника:

```bash
go build -o bin/helloapi ./cmd/server        # Linux/macOS
GOOS=windows GOARCH=amd64 go build -o bin/helloapi.exe ./cmd/server
./bin/helloapi
```

## Маршруты

| Метод и путь | Ответ |
| :--- | :--- |
| `GET /hello` | 200, текст `Hello, world!` |
| `GET /user` | 200, JSON `{"id":"<uuid>","name":"Gopher"}` и `Content-Type: application/json` |
| `GET /health` | 200, JSON `{"status":"ok","time":"<RFC3339>"}` |

Примеры запросов:

```bash
curl -i http://localhost:8080/hello
curl -i http://localhost:8080/user
curl -i http://localhost:8080/health
```

## Структура проекта

Дерево построено командой `eza -T`:

```text
pr1/
|-- cmd/server/
|   |-- main.go          # роутер, обработчики, middleware логирования
|   `-- main_test.go     # тесты маршрутов через httptest
|-- images/              # скриншоты для отчёта
|-- go.mod
|-- go.sum
|-- README.md
`-- report.md
```

## Конфигурация

- `APP_PORT` - порт HTTP-сервера, по умолчанию `8080`.
- Каждый запрос логируется в stdout: `GET /hello 12us`.
- Внешняя зависимость `github.com/google/uuid v1.6.0` зафиксирована в
  `go.mod` и `go.sum`, UUID реально генерируется в обработчике `/user`.
