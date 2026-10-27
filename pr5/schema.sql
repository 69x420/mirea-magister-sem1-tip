CREATE TABLE IF NOT EXISTS tasks (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	title      TEXT NOT NULL CHECK (length(trim(title)) > 0),
	done       BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tasks_done ON tasks (done);
