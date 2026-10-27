# Практическое занятие 5 - база данных через database/sql

**Выполнил:** Зайцев А. А., группа ЭФМО-02-26.

Работа с реляционной базой данных из Go: подключение через `database/sql`,
пул соединений, параметризованные `INSERT` и `SELECT`, транзакции. В
качестве СУБД используется SQLite, поэтому отдельный сервер не нужен.

## Требования

- Go 1.21+
- `CGO_ENABLED=1` и компилятор C (драйвер `github.com/mattn/go-sqlite3`)
- консольный клиент `sqlite3` для проверки данных

## Запуск

```bash
cd ~/mirea/magister/semester1/tid/pr5
make clean          # удалить старую базу и артефакты
make db             # создать таблицу из schema.sql
make run            # вставка, выборка и печать задач
make show           # SELECT всех задач через консольный клиент
make test           # юнит-тесты репозитория
```

Строка подключения берётся из переменной `DATABASE_URL`; по умолчанию
используется файл `todo.db` в каталоге проекта:

```bash
DATABASE_URL="file:/tmp/other.db?_journal_mode=WAL" make run
```

## API репозитория

| Метод | Назначение | SQL |
| :--- | :--- | :--- |
| `CreateTask(ctx, title) (int, error)` | создать задачу, вернуть `id` | `INSERT ... RETURNING id` |
| `ListTasks(ctx) ([]Task, error)` | все задачи по возрастанию `id` | `SELECT ... ORDER BY id` |
| `ListDone(ctx, done) ([]Task, error)` | фильтр по статусу | `SELECT ... WHERE done = ?` |
| `FindByID(ctx, id) (*Task, error)` | одна задача или `ErrNotFound` | `SELECT ... WHERE id = ?` |
| `SetDone(ctx, id, done) error` | изменить статус | `UPDATE ... WHERE id = ?` |
| `CreateMany(ctx, titles) error` | массовая вставка в транзакции | `INSERT` в `tx` |
| `Count(ctx) (int, error)` | количество задач | `SELECT COUNT(*)` |

## Структура проекта

Дерево построено командой `eza -T --git-ignore`:

```text
pr5/
|-- db.go
|-- go.mod
|-- go.sum
|-- images/
|   |-- 01-structure.png
|   |-- 02-schema.png
|   |-- 03-run.png
|   |-- 04-select.png
|   `-- 05-tests.png
|-- main.go
|-- Makefile
|-- README.md
|-- report.md
|-- repository.go
|-- repository_test.go
|-- requests.md
`-- schema.sql
```

## Конфигурация

- `DATABASE_URL` - строка подключения, по умолчанию `todo.db` в режиме WAL.
- Пул: `SetMaxOpenConns(10)`, `SetMaxIdleConns(5)`,
  `SetConnMaxLifetime(30 * time.Minute)`.
- Проверка соединения: `PingContext` с таймаутом 2 секунды.
- Все запросы выполняются через `QueryContext`/`QueryRowContext`/`ExecContext`.

Проверочные SQL-запросы собраны в `requests.md`.
