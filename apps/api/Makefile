.PHONY: help run build test clean migrate

help:
	@echo "Available targets:"
	@echo "  make run      - Run the server"
	@echo "  make build    - Build the server binary"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make migrate  - Run database migrations"

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
	rm -rf uploads/chunks/*

migrate:
	mysql -u root -p < sql/init.sql

install-deps:
	go mod download
	go mod tidy
