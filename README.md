# Clean Architecture API

Một RESTful API server được xây dựng theo Clean Architecture với Golang và Gin framework, bao gồm authentication và authorization.

## Cấu trúc dự án

```
clean-architecture-api/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
├── internal/
│   ├── domain/                     # Business logic layer
│   │   ├── entities/               # Domain entities
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
├── configs/                       # Configuration files
├── go.mod                         # Go modules
└── README.md                      # Documentation
```

## Tính năng

- ✅ Clean Architecture pattern
- ✅ JWT Authentication & Authorization
- ✅ Role-based access control
- ✅ RESTful API endpoints
- ✅ Database integration (PostgreSQL)
- ✅ Structured logging
- ✅ Input validation
- ✅ Error handling
- ✅ Pagination support

## Cài đặt

### Yêu cầu

- Go 1.21+
- PostgreSQL 12+

### Bước 1: Clone repository

```bash
git clone <repository-url>
cd clean-architecture-api
```

### Bước 2: Cài đặt dependencies

```bash
go mod tidy
```

### Bước 3: Cấu hình database

1. Tạo database PostgreSQL:
```sql
CREATE DATABASE clean_architecture_api;
```

2. Tạo file `.env` từ `env.example`:
```bash
cp env.example .env
```

3. Cập nhật thông tin database trong file `.env`

### Bước 4: Chạy ứng dụng

```bash
go run cmd/server/main.go
```

Server sẽ chạy tại `http://localhost:8080`

## API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Đăng ký user mới | ❌ |
| POST | `/api/v1/auth/login` | Đăng nhập | ❌ |
| POST | `/api/v1/auth/refresh` | Refresh token | ❌ |

### Users (Admin only)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/users` | Lấy danh sách users | ✅ (Admin) |
| GET | `/api/v1/users/:id` | Lấy thông tin user | ✅ (Admin) |
| PUT | `/api/v1/users/:id` | Cập nhật user | ✅ (Admin) |
| DELETE | `/api/v1/users/:id` | Xóa user | ✅ (Admin) |

### Products

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/products` | Lấy danh sách products | ❌ |
| GET | `/api/v1/products/:id` | Lấy thông tin product | ❌ |
| GET | `/api/v1/products/category/:category` | Lấy products theo category | ❌ |
| POST | `/api/v1/products` | Tạo product mới | ✅ |
| PUT | `/api/v1/products/:id` | Cập nhật product | ✅ |
| DELETE | `/api/v1/products/:id` | Xóa product | ✅ |

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |

## Authentication

### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "message": "Login successful",
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 1703123456
  }
}
```

### Sử dụng token

```bash
curl -X GET http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Roles

- `user`: Người dùng thông thường
- `admin`: Quản trị viên (có quyền truy cập tất cả endpoints)

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | 8080 |
| `ENV` | Environment | development |
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_USER` | Database user | postgres |
| `DB_PASSWORD` | Database password | password |
| `DB_NAME` | Database name | clean_architecture_api |
| `JWT_SECRET_KEY` | JWT secret key | your-secret-key-change-in-production |

## Development

### Chạy tests

```bash
go test ./...
```

### Format code

```bash
go fmt ./...
```

### Lint code

```bash
golangci-lint run
```

## Production Deployment

1. Cập nhật `JWT_SECRET_KEY` với một giá trị bảo mật
2. Cấu hình database production
3. Sử dụng reverse proxy (nginx)
4. Cấu hình SSL/TLS
5. Monitoring và logging

## License

