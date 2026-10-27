package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func newTestRepo(t *testing.T) *Repo {
	t.Helper()

	db, err := openDB("file::memory:?cache=shared&_foreign_keys=1")
	if err != nil {
		t.Fatalf("openDB: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if _, err := db.Exec(`DROP TABLE IF EXISTS tasks`); err != nil {
		t.Fatalf("drop: %v", err)
	}

	ctx := context.Background()
	if err := applySchema(ctx, db); err != nil {
		t.Fatalf("applySchema: %v", err)
	}
	return NewRepo(db)
}

func TestCreateTask(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	id, err := repo.CreateTask(ctx, "Первая задача")
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if id != 1 {
		t.Fatalf("id = %d, want 1", id)
	}

	task, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if task.Title != "Первая задача" {
		t.Errorf("title = %q", task.Title)
	}
	if task.Done {
		t.Error("new task must be not done")
	}
	if task.CreatedAt.IsZero() {
		t.Error("created_at must be set")
	}
}

func TestCreateTaskEmptyTitle(t *testing.T) {
	repo := newTestRepo(t)

	if _, err := repo.CreateTask(context.Background(), "   "); err == nil {
		t.Fatal("expected CHECK constraint error for blank title")
	}
}

func TestListTasks(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	for _, title := range []string{"Первая", "Вторая", "Третья"} {
		if _, err := repo.CreateTask(ctx, title); err != nil {
			t.Fatalf("CreateTask(%q): %v", title, err)
		}
	}

	tasks, err := repo.ListTasks(ctx)
	if err != nil {
		t.Fatalf("ListTasks: %v", err)
	}
	if len(tasks) != 3 {
		t.Fatalf("len = %d, want 3", len(tasks))
	}
	for i, task := range tasks {
		if task.ID != i+1 {
			t.Errorf("tasks[%d].ID = %d, want %d", i, task.ID, i+1)
		}
	}
}

func TestListDone(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	for _, title := range []string{"Первая", "Вторая", "Третья"} {
		if _, err := repo.CreateTask(ctx, title); err != nil {
			t.Fatalf("CreateTask: %v", err)
		}
	}
	if err := repo.SetDone(ctx, 2, true); err != nil {
		t.Fatalf("SetDone: %v", err)
	}

	done, err := repo.ListDone(ctx, true)
	if err != nil {
		t.Fatalf("ListDone(true): %v", err)
	}
	if len(done) != 1 || done[0].ID != 2 {
		t.Fatalf("done = %+v, want only id=2", done)
	}

	todo, err := repo.ListDone(ctx, false)
	if err != nil {
		t.Fatalf("ListDone(false): %v", err)
	}
	if len(todo) != 2 || todo[0].ID != 1 || todo[1].ID != 3 {
		t.Fatalf("not done = %+v, want ids 1 and 3", todo)
	}
}

func TestFindByID(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	id, err := repo.CreateTask(ctx, "Найти меня")
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	task, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if task.ID != id || task.Title != "Найти меня" {
		t.Errorf("task = %+v", task)
	}

	if _, err := repo.FindByID(ctx, 999); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestSetDone(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	id, err := repo.CreateTask(ctx, "Отметить")
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	if err := repo.SetDone(ctx, id, true); err != nil {
		t.Fatalf("SetDone: %v", err)
	}

	task, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if !task.Done {
		t.Error("task must be done")
	}

	if err := repo.SetDone(ctx, 999, true); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestCreateMany(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	titles := []string{"Массовая 1", "Массовая 2", "Массовая 3"}
	if err := repo.CreateMany(ctx, titles); err != nil {
		t.Fatalf("CreateMany: %v", err)
	}

	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if count != len(titles) {
		t.Fatalf("count = %d, want %d", count, len(titles))
	}
}

func TestCreateManyRollback(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	err := repo.CreateMany(ctx, []string{"Хорошая", "", "Тоже хорошая"})
	if err == nil {
		t.Fatal("expected error for blank title")
	}

	count, countErr := repo.Count(ctx)
	if countErr != nil {
		t.Fatalf("Count: %v", countErr)
	}
	if count != 0 {
		t.Fatalf("count = %d, want 0 (transaction must roll back)", count)
	}
}

func TestContextTimeout(t *testing.T) {
	repo := newTestRepo(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()

	if _, err := repo.ListTasks(ctx); err == nil {
		t.Fatal("expected context deadline error")
	}
}
