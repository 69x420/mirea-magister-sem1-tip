.PHONY: run build test fmt vet tidy clean

PORT ?= 8080

run:
	PORT=$(PORT) go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -count=1

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf bin
