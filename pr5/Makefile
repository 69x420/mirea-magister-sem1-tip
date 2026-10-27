APP ?= app
DB  ?= todo.db

.PHONY: run build test fmt vet tidy db show schema clean

run:
	go run .

build:
	go build -o bin/$(APP) .

test:
	go test ./... -count=1

fmt:
	gofmt -l -w .

vet:
	go vet ./...

tidy:
	go mod tidy

db:
	sqlite3 $(DB) < schema.sql

show:
	sqlite3 -header -column $(DB) "SELECT id, title, done, created_at FROM tasks ORDER BY id;"

schema:
	sqlite3 $(DB) ".schema tasks"

clean:
	rm -rf bin $(DB) $(DB)-wal $(DB)-shm
