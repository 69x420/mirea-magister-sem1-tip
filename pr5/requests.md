# Проверочные запросы и команды (ПЗ 5)

База: `todo.db` (SQLite), строка подключения по умолчанию
`file:todo.db?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=1`,
переопределяется переменной `DATABASE_URL`.

## 1. Создание таблицы

```bash
make db
```

```sql
CREATE TABLE IF NOT EXISTS tasks (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	title      TEXT NOT NULL CHECK (length(trim(title)) > 0),
	done       BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL
);
```

## 2. Схема

```bash
make schema
```

```text
CREATE TABLE tasks (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	title      TEXT NOT NULL CHECK (length(trim(title)) > 0),
	done       BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL
);
CREATE INDEX idx_tasks_done ON tasks (done);
```

## 3. Запуск приложения

```bash
make run
```

Ожидаемый вывод: подключение, три `INSERT`, список задач, выборка
невыполненных, поиск задачи по `id`, массовая вставка в транзакции.

## 4. Все задачи

```bash
make show
```

```sql
SELECT id, title, done, created_at FROM tasks ORDER BY id;
```

## 5. Выполненные задачи

```sql
SELECT id, title, done, created_at FROM tasks WHERE done = 1 ORDER BY id;
```

## 6. Задача по идентификатору

```sql
SELECT id, title, done, created_at FROM tasks WHERE id = 2;
```

## 7. Количество записей

```sql
SELECT COUNT(*) FROM tasks;
```

## 8. Очистка

```bash
make clean
```

Удаляет `bin/`, файл базы и служебные файлы WAL.
