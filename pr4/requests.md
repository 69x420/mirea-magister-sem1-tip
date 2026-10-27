# Тестовые запросы к API

Сервер запускается командой `make run` (или `PORT=8080 ./bin/server`) и
слушает `http://localhost:8080`. Все маршруты находятся под префиксом `/api`.

## 1. Проверка состояния

```bash
curl -s http://localhost:8080/health
```

```json
{"status":"ok"}
```

## 2. Создание задачи

```bash
curl -s -X POST http://localhost:8080/api/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"Выучить chi"}'
```

```json
{"id":1,"title":"Выучить chi","done":false,
 "created_at":"2026-10-23T11:43:48.100730046Z",
 "updated_at":"2026-10-23T11:43:48.100730046Z"}
```

## 3. Список задач

```bash
curl -s http://localhost:8080/api/tasks
```

```json
{"items":[{"id":1,"title":"Выучить chi","done":false,...},
          {"id":2,"title":"Написать отчёт","done":false,...},
          {"id":3,"title":"Проверить тесты","done":false,...}],
 "total":3,"page":1,"limit":10}
```

## 4. Пагинация

```bash
curl -s "http://localhost:8080/api/tasks?page=2&limit=2"
```

```json
{"items":[{"id":3,"title":"Проверить тесты",...},
          {"id":4,"title":"Сделать ПЗ 4",...}],
 "total":5,"page":2,"limit":2}
```

## 5. Фильтр по статусу

```bash
curl -s "http://localhost:8080/api/tasks?done=true"
```

```json
{"items":[{"id":2,"title":"Написать отчёт","done":true,...}],
 "total":1,"page":1,"limit":10}
```

## 6. Поиск по подстроке

```bash
curl -s "http://localhost:8080/api/tasks?q=chi"
```

```json
{"items":[{"id":1,"title":"Выучить chi","done":false,...}],
 "total":1,"page":1,"limit":10}
```

## 7. Получение задачи по идентификатору

```bash
curl -s http://localhost:8080/api/tasks/1
```

```json
{"id":1,"title":"Выучить chi","done":false,...}
```

## 8. Полное обновление (PUT)

```bash
curl -s -X PUT http://localhost:8080/api/tasks/1 \
  -H 'Content-Type: application/json' \
  -d '{"title":"Выучить chi глубже","done":true}'
```

```json
{"id":1,"title":"Выучить chi глубже","done":true,...}
```

## 9. Удаление

```bash
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/api/tasks/5
```

```
204
```

## 10. Ошибки

```bash
curl -s http://localhost:8080/api/tasks/abc
```

```json
{"error":"invalid id"}
```

```bash
curl -s http://localhost:8080/api/tasks/999
```

```json
{"error":"task not found"}
```

```bash
curl -s -X POST http://localhost:8080/api/tasks \
  -H 'Content-Type: application/json' -d '{"title":"ab"}'
```

```json
{"error":"title length must be between 3 and 100 characters"}
```

```bash
curl -s "http://localhost:8080/api/tasks?limit=0"
```

```json
{"error":"limit must be between 1 and 100"}
```

```bash
curl -s -o /dev/null -w "%{http_code}\n" -X PATCH http://localhost:8080/api/tasks
```

```
405
```

## 11. CORS preflight

```bash
curl -s -i -X OPTIONS http://localhost:8080/api/tasks
```

```
HTTP/1.1 204 No Content
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```
