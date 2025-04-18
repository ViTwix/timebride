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

.PHONY: run
run:
go run cmd/app/main.go

.PHONY: docker-up
docker-up:
docker-compose up -d

.PHONY: docker-down
docker-down:
docker-compose down

.PHONY: test
test:
go test -v ./...

.PHONY: deps
deps:
go mod download

.PHONY: clean
clean:
rm -rf bin/
go clean -cache
