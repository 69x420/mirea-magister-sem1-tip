package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = defaultDSN
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("openDB error: %v", err)
	}
	defer func() { _ = db.Close() }()

	repo := NewRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := applySchema(ctx, db); err != nil {
		log.Fatalf("applySchema error: %v", err)
	}

	titles := []string{"Сделать ПЗ №5", "Купить кофе", "Проверить отчёты"}
	for _, title := range titles {
		id, err := repo.CreateTask(ctx, title)
		if err != nil {
			log.Fatalf("CreateTask error: %v", err)
		}
		log.Printf("Inserted task id=%d (%s)", id, title)
	}

	tasks, err := repo.ListTasks(ctx)
	if err != nil {
		log.Fatalf("ListTasks error: %v", err)
	}

	fmt.Println("=== Tasks ===")
	printTasks(tasks)

	open, err := repo.ListDone(ctx, false)
	if err != nil {
		log.Fatalf("ListDone error: %v", err)
	}
	fmt.Println("=== Not done ===")
	printTasks(open)

	task, err := repo.FindByID(ctx, 2)
	if errors.Is(err, ErrNotFound) {
		log.Println("FindByID: task not found")
	} else if err != nil {
		log.Fatalf("FindByID error: %v", err)
	} else {
		fmt.Printf("=== Task #%d ===\ntitle=%s\ndone=%v\ncreated_at=%s\n",
			task.ID, task.Title, task.Done, task.CreatedAt.Format(time.RFC3339))
	}

	if err := repo.CreateMany(ctx, []string{"Ещё одна задача", "Последняя задача"}); err != nil {
		log.Fatalf("CreateMany error: %v", err)
	}
	log.Println("CreateMany inserted 2 tasks in one transaction")

	if err := repo.SetDone(ctx, 1, true); err != nil {
		log.Fatalf("SetDone error: %v", err)
	}
	log.Println("SetDone marked task id=1 as done")

	done, err := repo.ListDone(ctx, true)
	if err != nil {
		log.Fatalf("ListDone error: %v", err)
	}
	fmt.Println("=== Done ===")
	printTasks(done)

	total, err := repo.Count(ctx)
	if err != nil {
		log.Fatalf("Count error: %v", err)
	}

	all, err := repo.ListTasks(ctx)
	if err != nil {
		log.Fatalf("ListTasks error: %v", err)
	}
	fmt.Println("=== All tasks ===")
	printTasks(all)
	fmt.Printf("total=%d\n", total)

	if _, err := repo.FindByID(ctx, 999); errors.Is(err, ErrNotFound) {
		fmt.Println("FindByID 999 -> task not found")
	}
}

func printTasks(tasks []Task) {
	for _, t := range tasks {
		fmt.Printf("#%d | %-26s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}
}
