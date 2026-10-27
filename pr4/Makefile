.PHONY: run build test fmt vet tidy clean

PORT ?= 8080

run:
	PORT=$(PORT) go run .

build:
	go build -o bin/server .

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
