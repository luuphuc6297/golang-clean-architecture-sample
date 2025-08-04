# Clean Architecture API

A RESTful API server built with Clean Architecture principles using Golang and Gin framework, featuring authentication and authorization.

## Project Structure

```
clean-architecture-api/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
├── internal/
│   ├── domain/                     # Business logic layer
│   │   ├── entities/               # Domain entities
│   │   ├── constants/              # Domain constants
│   │   ├── errors/                 # Domain errors
│   │   └── repositories/           # Repository interfaces
│   ├── usecase/                    # Use cases (business logic)
│   ├── delivery/                   # Delivery layer
│   │   ├── http/                   # HTTP handlers
│   │   └── middleware/             # Middleware
│   └── infrastructure/             # Infrastructure layer
│       ├── database/               # Database connection
│       ├── auth/                   # Authentication service
│       └── repository/             # Repository implementations
├── pkg/                           # Shared packages
│   └── logger/                    # Logging utilities
├── scripts/                       # Build and test scripts
├── go.mod                         # Go modules
└── README.md                      # Documentation
```

## Features

- ✅ Clean Architecture pattern
- ✅ JWT Authentication & Authorization
- ✅ Role-based access control
- ✅ RESTful API endpoints
- ✅ Database integration (PostgreSQL/SQLite)
- ✅ Structured logging
- ✅ Input validation
- ✅ Error handling
- ✅ Pagination support
- ✅ Optimized build scripts
- ✅ Comprehensive testing

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 12+ (or SQLite for development)

### Installation

1. **Clone repository**
```bash
git clone <repository-url>
cd clean-architecture-api
```

2. **Install dependencies**
```bash
go mod tidy
```

3. **Configure environment**
```bash
cp env.example .env
# Edit .env with your database settings
```

4. **Run with PostgreSQL**
```bash
make run
```

5. **Run with SQLite (no Docker required)**
```bash
make run-sqlite
```

6. **Run with hot reload**
```bash
make dev
```

## API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Register new user | ❌ |
| POST | `/api/v1/auth/login` | User login | ❌ |
| POST | `/api/v1/auth/refresh` | Refresh token | ❌ |

### Users (Admin only)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/users` | List users | ✅ (Admin) |
| GET | `/api/v1/users/:id` | Get user by ID | ✅ (Admin) |
| PUT | `/api/v1/users/:id` | Update user | ✅ (Admin) |
| DELETE | `/api/v1/users/:id` | Delete user | ✅ (Admin) |

### Products

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/products` | List products | ❌ |
| GET | `/api/v1/products/:id` | Get product by ID | ❌ |
| GET | `/api/v1/products/category/:category` | Get products by category | ❌ |
| POST | `/api/v1/products` | Create product | ✅ |
| PUT | `/api/v1/products/:id` | Update product | ✅ |
| DELETE | `/api/v1/products/:id` | Delete product | ✅ |

## Development

### Available Commands

```bash
# Build application
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make format

# Check code quality
make lint

# Test API endpoints
make test-api

# Clean build artifacts
make clean
```

### Code Quality

The project uses golangci-lint for code quality checks:

```bash
# Run linter
golangci-lint run

# Run specific linters
golangci-lint run --enable=gocritic,gocyclo
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/usecase
```

## Configuration

### Environment Variables

```bash
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=clean_architecture_api

# JWT
JWT_SECRET_KEY=your-secret-key-change-in-production
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=7d
```

### Database Setup

#### PostgreSQL
```sql
CREATE DATABASE clean_architecture_api;
```

#### SQLite
No setup required - database file will be created automatically.

## Architecture

### Clean Architecture Layers

1. **Domain Layer**: Business entities and rules
2. **Use Case Layer**: Application business logic
3. **Delivery Layer**: HTTP handlers and middleware
4. **Infrastructure Layer**: Database, external services

### Key Principles

- **Dependency Inversion**: High-level modules don't depend on low-level modules
- **Single Responsibility**: Each component has one reason to change
- **Open/Closed**: Open for extension, closed for modification
- **Interface Segregation**: Clients depend only on interfaces they use

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## License

This project is licensed under the MIT License. 