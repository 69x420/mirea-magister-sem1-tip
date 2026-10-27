# Проверочные запросы (ПЗ 6)

Сервер: `http://localhost:8080`. База: SQLite, файл `pz6_gorm.db`,
строка подключения задаётся переменной `DB_DSN`.

## 1. Запуск

```bash
make build
DB_DSN="file:pz6_gorm.db?_journal_mode=WAL&_foreign_keys=1" PORT=8080 ./bin/server
```

Ожидаемые строки журнала:

```text
Connected to SQLite via GORM
AutoMigrate is done
listening on :8080
```

## 2. Проверка состояния

```bash
curl -s http://localhost:8080/health
```

```json
{"status":"ok"}
```

## 3. Создание пользователя

```bash
curl -s -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com"}'
```

```json
{"id":1,"name":"Alice","email":"alice@example.com","createdAt":"2026-10-23T12:59:07Z","updatedAt":"2026-10-23T12:59:07Z"}
```

## 4. Повторный e-mail (уникальный индекс)

```bash
curl -s -i -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com"}'
```

```text
HTTP/1.1 409 Conflict
{"error":"UNIQUE constraint failed: users.email"}
```

## 5. Пустые поля

```bash
curl -s -i -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"","email":""}'
```

```text
HTTP/1.1 400 Bad Request
{"error":"name and email are required"}
```

## 6. Создание заметки с тегами

```bash
curl -s -X POST http://localhost:8080/notes \
  -H 'Content-Type: application/json' \
  -d '{"title":"Первая заметка","content":"Текст заметки","userId":1,"tags":["go","gorm"]}'
```

```text
HTTP/1.1 201 Created
```

Ответ содержит автора и теги, потому что после `Create` выполняется
`Preload("User").Preload("Tags")`.

## 7. Получение заметки по идентификатору

```bash
curl -s http://localhost:8080/notes/1
```

```json
{
  "id": 1,
  "title": "Первая заметка",
  "content": "Текст заметки",
  "userId": 1,
  "user": {"id": 1, "name": "Alice", "email": "alice@example.com"},
  "tags": [{"id": 1, "name": "go"}, {"id": 2, "name": "gorm"}]
}
```

## 8. Список заметок и пользователей

```bash
curl -s http://localhost:8080/notes
curl -s http://localhost:8080/users
```

## 9. Ошибки

```bash
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8080/notes/abc
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8080/notes/99
```

```text
400
404
```

## 10. Схема базы

```bash
make schema
```

```text
CREATE TABLE `users` (`id` integer PRIMARY KEY AUTOINCREMENT,`name` text NOT NULL,`email` text NOT NULL,`created_at` datetime,`updated_at` datetime);
CREATE UNIQUE INDEX `idx_users_email` ON `users`(`email`);
CREATE TABLE `notes` (`id` integer PRIMARY KEY AUTOINCREMENT,`title` text NOT NULL,`content` text,`user_id` integer NOT NULL,`created_at` datetime,`updated_at` datetime,CONSTRAINT `fk_users_notes` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`));
CREATE TABLE `note_tags` (`tag_id` integer,`note_id` integer,PRIMARY KEY (`tag_id`,`note_id`),CONSTRAINT `fk_note_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags`(`id`),CONSTRAINT `fk_note_tags_note` FOREIGN KEY (`note_id`) REFERENCES `notes`(`id`));
```

## 11. Сводный запрос по связям

```sql
SELECT n.id, n.title, u.name AS author, group_concat(t.name) AS tags
FROM notes n
JOIN users u ON u.id = n.user_id
LEFT JOIN note_tags nt ON nt.note_id = n.id
LEFT JOIN tags t ON t.id = nt.tag_id
GROUP BY n.id;
```

```text
id  title           author  tags
--  --------------  ------  ---------
1   Первая заметка  Alice   go,gorm
2   Вторая заметка  Alice   go,sqlite
```

## 12. Очистка

```bash
make clean
```
