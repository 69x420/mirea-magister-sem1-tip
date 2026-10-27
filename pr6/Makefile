DB  ?= pz6_gorm.db
PORT ?= 8080

.PHONY: run build test fmt vet tidy show schema clean

run:
	PORT=$(PORT) go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -count=1

fmt:
	gofmt -l -w .

vet:
	go vet ./...

tidy:
	go mod tidy

show:
	sqlite3 -header -column $(DB) "SELECT id, title, user_id FROM notes ORDER BY id;"

schema:
	sqlite3 $(DB) ".schema users" ".schema notes" ".schema tags" ".schema note_tags"

clean:
	rm -rf bin $(DB) $(DB)-wal $(DB)-shm
