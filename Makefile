.PHONY: run restart test test-unit test-integration test-coverage test-clean

run:
	docker-compose up -d

restart:
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# Run all tests
test: test-unit test-integration

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@go test -v -race -short ./internal/services/... ./internal/repository/... ./pkg/...

# Run integration tests with testcontainers
test-integration:
	@echo "Running integration tests..."
	@echo "Starting test containers (PostgreSQL and Redis)..."
	@go test -v -race ./test/integration/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean test artifacts
test-clean:
	@echo "Cleaning test artifacts..."
	@rm -f coverage.out coverage.html
	@docker system prune -f