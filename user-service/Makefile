.PHONY: help build run test docker-up docker-down migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make build       - Build the service"
	@echo "  make run         - Run the service locally"
	@echo "  make test        - Run tests"
	@echo "  make docker-up   - Start Docker services"
	@echo "  make docker-down - Stop Docker services"
	@echo "  make migrate-up  - Run database migrations"
	@echo "  make migrate-down- Rollback migrations"

build:
	go build -o bin/user-service ./cmd/api

run:
	go run ./cmd/api -db-dsn=${FASHIONUSERS_DB_DSN} -jwt-secret=${JWT_SECRET}
migrate-up:
	migrate -path ./migrations -database "${FASHIONUSERS_DB_DSN}" up

migrate-down:
	migrate -path ./migrations -database "${FASHIONUSERS_DB_DSN}" down 1