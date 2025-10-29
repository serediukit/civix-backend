# Civix Backend

A Go REST API with JWT authentication, built with Gin, PostgreSQL, and Redis.

## Features

- User registration and authentication with JWT
- Password hashing with bcrypt
- Token refresh mechanism
- Protected routes with JWT middleware
- Redis for token blacklisting
- Database migrations with GORM
- Environment-based configuration
- Graceful shutdown
- Health check endpoint

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 13 or higher
- Redis 6 or higher

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/serediukit/civix-backend.git
   cd civix-backend
   ```

2. Copy the example environment file and update the values:
   ```bash
   cp .env.example .env
   ```

3. Update the `.env` file with your database and Redis credentials.

4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login with email and password
- `POST /api/v1/auth/logout` - Logout (invalidate token)
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/auth/me` - Get current user profile

### Users

- `GET /api/v1/users/me` - Get current user profile
- `PUT /api/v1/users/me` - Update current user profile
- `PUT /api/v1/users/me/password` - Change password
- `DELETE /api/v1/users/me` - Delete current user account

### Health Check

- `GET /health` - Health check endpoint

## Project Structure

```
.
├── cmd/                  # Application entry points
│   └── api/              # Main API server
│       └── main.go       # Application entry point
├── internal/             # Private application code
│   ├── config/           # Configuration
│   ├── controller/       # HTTP controllers
│   ├── middleware/       # HTTP middleware
│   ├── model/            # Data models
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   └── util/             # Utility functions
├── pkg/                  # Public library code
│   └── database/         # Database connections
├── .env                  # Environment variables
├── .gitignore           # Git ignore file
├── go.mod               # Go module file
└── README.md            # This file
```

## Running Tests

```bash
go test -v ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
