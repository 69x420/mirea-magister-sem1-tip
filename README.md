# Технологии индустриального программирования

Практические работы 1-6. **Зайцев А. А., ЭФМО-02-26.**

| N | Работа | Реализовано | Отчёт |
| :-: | :--- | :--- | :--- |
| 1 | [HTTP-сервер и UUID](pr1) | `/hello`, `/user`, UUID, `/health`, `APP_PORT`, лог запросов | [report.md](pr1/report.md), [report.docx](pr1/report.docx) (ГОСТ 7.32) |
| 2 | [Структура Go-проекта](pr2) | `cmd`/`internal`/`utils`, `/ping`, `X-Request-Id`, JSON-ошибки | [report.md](pr2/report.md), [report.docx](pr2/report.docx) (ГОСТ 7.32) |
| 3 | [HTTP API на net/http](pr3) | CRUD, фильтр `q`, валидация, CORS, `PATCH`/`DELETE`, graceful shutdown | [report.md](pr3/report.md), [report.docx](pr3/report.docx) (ГОСТ 7.32) |
| 4 | [CRUD-сервис на chi](pr4) | `chi.Router`, CRUD, пагинация, фильтр `done`, поиск, CORS, middleware логов | [report.md](pr4/report.md), [report.docx](pr4/report.docx) (ГОСТ 7.32) |
| 5 | [SQLite через `database/sql`](pr5) | пул соединений, `INSERT`/`SELECT`, `context`, фильтр `done`, поиск по `id`, транзакция `CreateMany` | [report.md](pr5/report.md), [report.docx](pr5/report.docx) (ГОСТ 7.32) |
| 6 | [GORM: модели, миграции и связи](pr6) | `User`/`Note`/`Tag`, `AutoMigrate`, 1:N и M:N, `Preload`, `FirstOrCreate`, `409` на дубликат | [report.md](pr6/report.md), [report.docx](pr6/report.docx) (ГОСТ 7.32) |

Каждая работа - отдельный Go-модуль. Команды запуска, переменные окружения и
проверочные запросы приведены в `README.md` соответствующей работы.
