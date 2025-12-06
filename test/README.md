# Testing Guide

This directory contains all tests for the civix-backend project, including unit tests and integration tests.

## Test Structure

```
test/
├── integration/          # Integration tests with real database and Redis
│   ├── api_test.go      # API endpoint integration tests
│   ├── setup.go         # Test database setup helpers
│   └── testcontainers.go # Container management for tests
└── README.md            # This file
```

Unit tests are located alongside the code they test:
- `internal/services/*_test.go` - Service layer unit tests
- `internal/repository/*_test.go` - Repository layer unit tests (if any)

## Running Tests

### Prerequisites

- Docker must be running (for integration tests with testcontainers)
- Go 1.23 or higher

### Run All Tests

```bash
make test
```

This runs both unit tests and integration tests.

### Run Unit Tests Only

```bash
make test-unit
```

Unit tests are fast and don't require Docker. They use mocks for dependencies.

### Run Integration Tests Only

```bash
make test-integration
```

Integration tests use testcontainers to spin up real PostgreSQL and Redis instances. The containers are automatically created and destroyed.

### Run Tests with Coverage

```bash
make test-coverage
```

This generates a coverage report in `coverage.html` that you can open in a browser.

### Clean Test Artifacts

```bash
make test-clean
```

## Test Technologies

- **testify** - Assertion library and test suite framework
- **testcontainers-go** - Provides real PostgreSQL and Redis containers for integration tests
- **testcontainers-go/modules/postgres** - PostgreSQL-specific testcontainer module
- **testcontainers-go/modules/redis** - Redis-specific testcontainer module

## Writing Tests

### Unit Tests

Unit tests should:
- Use mocks for external dependencies
- Be fast (< 100ms per test)
- Test individual functions in isolation
- Be located in `*_test.go` files next to the code they test

Example:
```go
func TestAuthService_Login_Success(t *testing.T) {
    mockUserRepo := new(MockUserRepository)
    // ... setup mocks
    authService := NewAuthService(mockUserRepo, ...)
    // ... test assertions
}
```

### Integration Tests

Integration tests should:
- Use real database and Redis via testcontainers
- Test full request/response cycles
- Be located in `test/integration/`
- Clean up data between tests

Example:
```go
func (suite *APITestSuite) TestHealthCheck() {
    req, _ := http.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)
    assert.Equal(suite.T(), http.StatusOK, w.Code)
}
```

## Testcontainers

Integration tests automatically manage Docker containers for:

- **PostgreSQL** - Using `postgis/postgis:15-3.3` image
- **Redis** - Using `redis:7-alpine` image

Containers are:
- Started before each test suite
- Cleaned up after tests complete
- Isolated per test run
- Automatically configured with correct ports and credentials

No manual Docker setup required!

## Continuous Integration

These tests are designed to run in CI/CD environments. The Makefile commands work in any environment with Docker installed.

## Troubleshooting

### "Cannot connect to Docker daemon"
Ensure Docker is running on your system.

### "Port already in use"
Testcontainers uses random ports, but if you have conflicts, stop other containers:
```bash
docker ps
docker stop <container-id>
```

### Integration tests are slow
This is normal. Integration tests spin up real databases and can take 10-30 seconds to start.

### Tests fail with database errors
Make sure migrations are up to date in the `migrations/` directory. Integration tests run all migrations automatically.
