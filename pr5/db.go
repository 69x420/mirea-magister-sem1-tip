package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const defaultDSN = "file:todo.db?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=1"

const schema = `
CREATE TABLE IF NOT EXISTS tasks (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	title      TEXT NOT NULL CHECK (length(trim(title)) > 0),
	done       BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tasks_done ON tasks (done);
`

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	log.Println("Connected to SQLite")
	return db, nil
}

func applySchema(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return err
	}
	log.Println("Schema is ready")
	return nil
}
