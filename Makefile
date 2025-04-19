.PHONY: help
help:
@echo "Usage: make [target]"
@echo ""
@echo "Targets:"
@echo "  run              Run the application"
@echo "  docker-up        Start Docker containers"
@echo "  docker-down      Stop Docker containers"
@echo "  test            Run tests"
@echo "  deps            Download dependencies"
@echo "  clean           Clean build artifacts"
@echo "  migrate-up       Run migrations up"
@echo "  migrate-down     Run migrations down"
@echo "  frontend-install Install frontend dependencies"
@echo "  frontend-start   Start frontend development server"
@echo "  frontend-build   Build frontend assets"

.PHONY: run build test clean docker-up docker-down migrate-up migrate-down frontend-install frontend-start frontend-build

# Go commands
run:
	go run cmd/app/main.go

build:
	go build -o bin/app cmd/app/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Database commands
migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/timebride?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/timebride?sslmode=disable" down

# Frontend commands
frontend-install:
	cd web && ./install.sh

frontend-start:
	cd web && npm start

frontend-build:
	cd web && npm run build

.PHONY: deps
deps:
go mod download
