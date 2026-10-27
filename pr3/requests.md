# Тестовые запросы к API

Сервер запускается командой `make run` (или `PORT=8080 ./bin/server`) и
слушает `http://localhost:8080`.

## 1. Проверка состояния

```bash
curl -i http://localhost:8080/health
```

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{"status":"ok"}
```

## 2. Создание задачи

```bash
curl -i -X POST http://localhost:8080/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"Купить молоко"}'
```

```
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8

{"id":1,"title":"Купить молоко","done":false,
 "created_at":"2026-10-23T11:02:02.424928623Z",
 "updated_at":"2026-10-23T11:02:02.424928623Z"}
```

## 3. Список задач

```bash
curl -i http://localhost:8080/tasks
```

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

[{"id":1,"title":"Купить молоко","done":false,...},
 {"id":2,"title":"Позвонить маме","done":false,...},
 {"id":3,"title":"Купить хлеб","done":false,...}]
```

## 4. Фильтр по подстроке

```bash
curl -i "http://localhost:8080/tasks?q=купить"
```

```
HTTP/1.1 200 OK

[{"id":1,"title":"Купить молоко",...},
 {"id":3,"title":"Купить хлеб",...}]
```

## 5. Получение задачи по идентификатору

```bash
curl -i http://localhost:8080/tasks/1
```

```
HTTP/1.1 200 OK

{"id":1,"title":"Купить молоко","done":false,...}
```

## 6. Отметка о выполнении (PATCH)

```bash
curl -i -X PATCH http://localhost:8080/tasks/1 \
  -H 'Content-Type: application/json' \
  -d '{"done":true}'
```

```
HTTP/1.1 200 OK

{"id":1,"title":"Купить молоко","done":true,...}
```

## 7. Переименование задачи (PATCH)

```bash
curl -i -X PATCH http://localhost:8080/tasks/2 \
  -H 'Content-Type: application/json' \
  -d '{"title":"Позвонить маме вечером"}'
```

```
HTTP/1.1 200 OK

{"id":2,"title":"Позвонить маме вечером","done":false,...}
```

## 8. Удаление задачи

```bash
curl -i -X DELETE http://localhost:8080/tasks/3
```

```
HTTP/1.1 204 No Content
```

## 9. Ошибки

```bash
curl -i -X POST http://localhost:8080/tasks \
  -H 'Content-Type: application/json' -d '{"title":"ab"}'
```

```
HTTP/1.1 422 Unprocessable Entity

{"error":"title length must be between 3 and 140 characters"}
```

```bash
curl -i http://localhost:8080/tasks/9999
```

```
HTTP/1.1 404 Not Found

{"error":"task not found"}
```

```bash
curl -i http://localhost:8080/tasks/abc
```

```
HTTP/1.1 400 Bad Request

{"error":"invalid id"}
```

```bash
curl -i -X POST http://localhost:8080/tasks \
  -H 'Content-Type: application/json' -d '{"title":'
```

```
HTTP/1.1 400 Bad Request

{"error":"invalid json: unexpected EOF"}
```

```bash
curl -i -X PUT http://localhost:8080/tasks
```

```
HTTP/1.1 405 Method Not Allowed
```

## 10. CORS

```bash
curl -i -X OPTIONS http://localhost:8080/tasks \
  -H 'Origin: http://localhost:3000'
```

```
HTTP/1.1 204 No Content
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PATCH, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type
```
