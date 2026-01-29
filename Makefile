.PHONY: help dev dev:api dev:web build build:api build:web install clean

help:
	@echo "Available targets:"
	@echo "  make dev         - Start both API and Web (dev mode)"
	@echo "  make dev:api     - Start API server only"
	@echo "  make dev:web     - Start Web frontend only"
	@echo "  make build       - Build both API and Web"
	@echo "  make build:api   - Build API binary"
	@echo "  make build:web   - Build Web frontend"
	@echo "  make install     - Install all dependencies"
	@echo "  make clean       - Clean build artifacts"

dev:
	@echo "Starting API and Web..."
	pnpm dev

dev:api:
	@echo "Starting API server..."
	cd apps/api && go run cmd/server/main.go

dev:web:
	@echo "Starting Web frontend..."
	pnpm dev:web

build: build:web build:api

build:api:
	@echo "Building API..."
	cd apps/api && go build -o ../../bin/server cmd/server/main.go

build:web:
	@echo "Building Web..."
	pnpm build:web

install:
	@echo "Installing dependencies..."
	pnpm install
	cd apps/api && go mod download

clean:
	@echo "Cleaning build artifacts..."
	rm -rf apps/web/node_modules apps/web/dist apps/api/bin bin

# Database
migrate:
	mysql -u root -p < apps/api/sql/init.sql

# Development
setup: install
	@echo "Setup complete! Run 'make dev' to start development."
