package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func now() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}

var ErrNotFound = errors.New("task not found")

type Task struct {
	ID        int
	Title     string
	Done      bool
	CreatedAt time.Time
}

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo { return &Repo{DB: db} }

func (r *Repo) CreateTask(ctx context.Context, title string) (int, error) {
	const q = `INSERT INTO tasks (title, done, created_at) VALUES (?, FALSE, ?) RETURNING id;`

	var id int
	if err := r.DB.QueryRowContext(ctx, q, title, now()).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) ListTasks(ctx context.Context) ([]Task, error) {
	const q = `SELECT id, title, done, created_at FROM tasks ORDER BY id;`

	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Task, 0, 8)
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repo) ListDone(ctx context.Context, done bool) ([]Task, error) {
	const q = `SELECT id, title, done, created_at FROM tasks WHERE done = ? ORDER BY id;`

	rows, err := r.DB.QueryContext(ctx, q, done)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Task, 0, 8)
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repo) FindByID(ctx context.Context, id int) (*Task, error) {
	const q = `SELECT id, title, done, created_at FROM tasks WHERE id = ?;`

	var t Task
	err := r.DB.QueryRowContext(ctx, q, id).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repo) SetDone(ctx context.Context, id int, done bool) error {
	const q = `UPDATE tasks SET done = ? WHERE id = ?;`

	res, err := r.DB.ExecContext(ctx, q, done, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repo) CreateMany(ctx context.Context, titles []string) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	const q = `INSERT INTO tasks (title, done, created_at) VALUES (?, FALSE, ?);`

	for i, title := range titles {
		if _, err := tx.ExecContext(ctx, q, title, now()); err != nil {
			return fmt.Errorf("insert %d (%q): %w", i+1, title, err)
		}
	}

	return tx.Commit()
}

func (r *Repo) Count(ctx context.Context) (int, error) {
	const q = `SELECT COUNT(*) FROM tasks;`

	var n int
	if err := r.DB.QueryRowContext(ctx, q).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}
