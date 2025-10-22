# Makefile
.PHONY: help

help:
	@echo "Available commands:"
	@echo "  make run       - Run the application"
	@echo "  make build     - Build the application"
	@echo "  make test      - Run tests"
	@echo "  make migrate   - Run database migrations"
	@echo "  make seed      - Seed the database"
	@echo "  make docker    - Run with Docker Compose"

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v ./...

migrate-up:
	migrate -path internal/database/migrations -database "postgresql://postgres:password@localhost:5432/bmginventory?sslmode=disable" up

migrate-down:
	migrate -path internal/database/migrations -database "postgresql://postgres:password@localhost:5432/bmginventory?sslmode=disable" down

docker:
	docker-compose up -d

docker-down:
	docker-compose down

lint:
	golangci-lint run

swagger:
	swag init -g cmd/api/main.go