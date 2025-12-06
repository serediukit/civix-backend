# Testing Documentation

## Overview

This project now has comprehensive test coverage including both unit tests and integration tests.

## What Was Added

### 1. Testing Dependencies

Added the following packages via `go get`:
- `github.com/stretchr/testify` - Assertion and mocking library
- `github.com/testcontainers/testcontainers-go` - Container orchestration for tests
- `github.com/testcontainers/testcontainers-go/modules/postgres` - PostgreSQL testcontainer
- `github.com/testcontainers/testcontainers-go/modules/redis` - Redis testcontainer

### 2. Unit Tests

**Location:** `internal/services/*_test.go`

#### Auth Service Tests (`internal/services/auth_service_test.go`)
- ✅ User registration success
- ✅ User registration with existing email
- ✅ Login success
- ✅ Login with non-existent user
- ✅ Login with incorrect password
- ✅ Logout success
- ✅ Refresh token success
- ✅ Refresh token with blacklisted token
- ✅ Refresh token with invalid token

#### User Service Tests (`internal/services/user_service_test.go`)
- ✅ Get user by ID
- ✅ Get user by email
- ✅ Get user not found
- ✅ Get user without identifier
- ✅ Update profile success
- ✅ Update profile without user in context

**Total Unit Tests:** 15 tests, all passing

### 3. Integration Tests

**Location:** `test/integration/`

#### Test Infrastructure
- `testcontainers.go` - Manages PostgreSQL and Redis containers
- `setup.go` - Database setup and migration runner
- `api_test.go` - API endpoint tests

#### API Integration Tests
- ✅ Health check endpoint
- ✅ User registration and login flow
- ✅ Login with invalid credentials
- ✅ Authenticated endpoints access
- ✅ Unauthorized access handling
- ✅ Token refresh flow

**Features:**
- Automatic container lifecycle management
- Real PostgreSQL database with PostGIS
- Real Redis instance
- Automatic migration execution
- Clean state between tests

### 4. Makefile Commands

Added the following commands to `Makefile`:

```makefile
make test              # Run all tests (unit + integration)
make test-unit         # Run unit tests only (fast, no Docker)
make test-integration  # Run integration tests (with containers)
make test-coverage     # Generate coverage report (coverage.html)
make test-clean        # Clean test artifacts and containers
```

### 5. Helper Functions

Added to `pkg/database/postgres.go`:
```go
func NewDBFromDSN(ctx context.Context, dsn string) (*Store, error)
```

Added to `internal/middleware/auth.go`:
```go
func SetUserIDInContext(ctx context.Context, userID uint64) context.Context
```

Added to `internal/server/router.go`:
```go
func SetupRouter(...) *gin.Engine  // Exported for testing
```

### 6. Documentation

- `test/README.md` - Comprehensive testing guide
- `TESTING.md` - This file

## Quick Start

### Run All Tests
```bash
make test
```
**Note:** Integration tests require Docker to be running.

### Run Only Unit Tests (Fast)
```bash
make test-unit
```
Unit tests run quickly without Docker (typically < 5 seconds).

### Run Only Integration Tests
```bash
make test-integration
```
**Requirements:** Docker must be running. Integration tests will:
- Automatically start PostgreSQL (PostGIS) container
- Automatically start Redis container
- Run all database migrations
- Execute API endpoint tests
- Clean up containers after completion

### Generate Coverage Report
```bash
make test-coverage
open coverage.html  # View in browser
```

## Important Notes

### Integration Tests
The integration tests use testcontainers which requires:
1. **Docker Desktop** (or Docker daemon) running
2. **Network access** to pull container images (first run only)
3. **~30-60 seconds** for container startup and test execution

If Docker is not available, you can still run unit tests with `make test-unit`.

### Redis Client Setup
The test infrastructure properly initializes Redis using the project's `redis.CachedStore` wrapper, ensuring consistency between tests and production code.

## Test Architecture

### Unit Tests
- Use testify mocks for dependencies
- Fast execution (< 5 seconds total)
- No external dependencies
- Located next to source code

### Integration Tests
- Use testcontainers for real database/Redis
- Test full HTTP request/response cycle
- Run migrations automatically
- Clean state between tests
- Located in `test/integration/`

## Mock Implementations

All unit tests use mock implementations for:
- `UserRepository`
- `CityRepository`
- `CacheRepository`

Mocks are defined in test files using testify/mock.

## Container Configuration

Integration tests automatically configure:

**PostgreSQL:**
- Image: `postgis/postgis:15-3.3`
- Database: `civix_test`
- User/Password: `postgres/postgres`
- Migrations: Applied automatically

**Redis:**
- Image: `redis:7-alpine`
- Configuration: Default with verbose logging

Containers use random ports to avoid conflicts and are cleaned up automatically.

## Coverage

Current test coverage focuses on:
- ✅ Authentication flow (register, login, logout, refresh)
- ✅ User management (get, update)
- ✅ API endpoints
- ✅ Middleware (auth)
- ✅ Service layer business logic

Future coverage can be added for:
- Report service and endpoints
- Repository layer
- Utility packages

## CI/CD Integration

These tests are designed for CI/CD:
- No manual setup required
- Docker-based (works in any CI environment)
- Fast unit tests for quick feedback
- Integration tests for deployment validation

Example GitHub Actions:
```yaml
- name: Run unit tests
  run: make test-unit

- name: Run integration tests
  run: make test-integration
```

## Best Practices

1. **Unit tests should:**
   - Test one function at a time
   - Use mocks for dependencies
   - Be fast (< 100ms per test)
   - Not depend on external services

2. **Integration tests should:**
   - Test real scenarios end-to-end
   - Use real database and cache
   - Clean up data between tests
   - Verify API contracts

3. **When adding new features:**
   - Write unit tests for business logic
   - Write integration tests for new endpoints
   - Update this documentation

## Troubleshooting

See `test/README.md` for detailed troubleshooting guide.

## Summary

- ✅ 15 unit tests covering core services
- ✅ 6 integration tests covering API endpoints
- ✅ Automatic container management
- ✅ Full CI/CD ready
- ✅ Comprehensive documentation
- ✅ Make commands for easy execution
