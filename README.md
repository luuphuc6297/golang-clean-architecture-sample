# Clean Architecture API

A production-ready RESTful API server built with **Clean Architecture** principles using **Go** and **Gin** framework. Features comprehensive **JWT authentication**, **role-based authorization**, **monitoring**, and **multiple database support**.

## 🏗️ Project Structure

```
clean-architecture-api/
├── cmd/
│   └── server/
│       ├── main.go                 # Main entry point (PostgreSQL)
│       └── main_sqlite.go          # SQLite entry point  
├── internal/
│   ├── domain/                     # Business logic layer
│   │   ├── entities/               # Domain entities
│   │   ├── repositories/           # Repository interfaces
│   │   ├── constants/              # Application constants
│   │   ├── errors/                 # Custom error definitions
│   │   └── validators/             # Input validation
│   ├── usecase/                    # Use cases (business logic)
│   ├── delivery/                   # Delivery layer
│   │   ├── http/                   # HTTP handlers
│   │   │   ├── handlers/           # API handlers
│   │   │   └── server.go           # Server configuration
│   │   └── middleware/             # HTTP middleware
│   └── infrastructure/             # Infrastructure layer
│       ├── database/               # Database connections
│       ├── auth/                   # Authentication & authorization
│       └── repository/             # Repository implementations
├── pkg/                           # Shared packages
│   ├── logger/                    # Structured logging
│   └── newrelic/                  # New Relic monitoring
├── data/                          # Database files (SQLite)
├── scripts/                       # Utility scripts
├── docker-compose.yml             # Development environment
├── docker-compose.prod.yml        # Production environment
├── Dockerfile                     # Container image
├── Makefile                       # Build & development commands
└── sonar-project.properties       # SonarCloud configuration
```

## ✨ Features

- ✅ **Clean Architecture** pattern with clear separation of concerns
- ✅ **JWT Authentication** with access & refresh tokens
- ✅ **Policy-based Authorization** with role-based access control (RBAC)
- ✅ **Multiple Database Support** (PostgreSQL for production, SQLite/In-memory for local)
- ✅ **New Relic APM** integration for monitoring and performance tracking
- ✅ **SonarCloud** integration for code quality analysis
- ✅ **RESTful API** endpoints with proper HTTP status codes
- ✅ **Structured Logging** with configurable levels
- ✅ **Input Validation** and error handling
- ✅ **Docker Support** for containerized deployment
- ✅ **Health Check** endpoint
- ✅ **Audit Logging** for security events
- ✅ **Pagination Support** for list endpoints
- ✅ **Comprehensive Testing** with test helpers

## 🛠️ Technology Stack

- **Language**: Go 1.23+
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL (production), SQLite (development), In-memory (testing)
- **Authentication**: JWT with RS256/HS256 signing
- **ORM**: GORM v2
- **Monitoring**: New Relic APM
- **Code Quality**: SonarCloud
- **Logging**: Logrus
- **Containerization**: Docker & Docker Compose
- **Build Tool**: Make

## 🚀 Quick Start

### Prerequisites

- **Go 1.23+**
- **Docker & Docker Compose** (for production setup)
- **PostgreSQL 15+** (for production database)

### 1. Clone Repository

```bash
git clone <repository-url>
cd clean-architecture-api
```

### 2. Install Dependencies

```bash
make deps
# or
go mod tidy
```

### 3. Development Setup (Local)

#### Option A: In-Memory Database (Fastest)
```bash
# No setup required - uses in-memory SQLite
make run-memory
```

#### Option B: SQLite Database (Persistent)
```bash
# Copy SQLite environment configuration
cp env.sqlite.example .env

# Run with SQLite
make run-sqlite
```

#### Option C: PostgreSQL with Docker
```bash
# Copy PostgreSQL environment configuration
cp env.example .env

# Start PostgreSQL database
docker-compose up postgres -d

# Run application
make run
```

### 4. Production Setup

```bash
# Configure environment variables
cp env.example .env
# Edit .env with production values

# Start full production stack
docker-compose -f docker-compose.prod.yml up -d
```

The server will be available at `http://localhost:8080`

## 🗄️ Database Configuration

The application supports multiple database configurations:

### Production Environment
- **Database**: PostgreSQL 15+
- **Connection**: Via environment variables
- **Migrations**: Auto-migration on startup
- **Monitoring**: Full New Relic database monitoring

### Local Development
- **In-Memory**: SQLite in-memory database (fastest, no persistence)
- **SQLite File**: Persistent SQLite database in `./data/` directory
- **Docker PostgreSQL**: Full PostgreSQL setup via Docker Compose

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `ENV` | Environment (development/production) | development | No |
| `PORT` | Server port | 8080 | No |
| `DB_HOST` | Database host | localhost | Yes (PostgreSQL) |
| `DB_PORT` | Database port | 5432 | Yes (PostgreSQL) |
| `DB_USER` | Database user | postgres | Yes (PostgreSQL) |
| `DB_PASSWORD` | Database password | - | Yes (PostgreSQL) |
| `DB_NAME` | Database name | clean_architecture_api | Yes (PostgreSQL) |
| `SQLITE_DB_PATH` | SQLite database file path | ./data/clean_architecture_api.db | No |
| `JWT_SECRET_KEY` | JWT signing secret | - | Yes |
| `LOG_LEVEL` | Logging level | info | No |

## 📊 Monitoring & Observability

### New Relic APM Integration

The application includes comprehensive New Relic monitoring:

- **Application Performance Monitoring (APM)**
- **Database Query Monitoring**
- **Custom Metrics and Events**
- **Error Tracking and Alerting**
- **Distributed Tracing**

#### Configuration

| Variable | Description | Required |
|----------|-------------|----------|
| `NEW_RELIC_ENABLED` | Enable/disable New Relic | No |
| `NEW_RELIC_APP_NAME` | Application name in New Relic | No |
| `NEW_RELIC_LICENSE_KEY` | New Relic license key | Yes (if enabled) |

```bash
# Enable New Relic monitoring
NEW_RELIC_ENABLED=true
NEW_RELIC_APP_NAME=clean-architecture-api
NEW_RELIC_LICENSE_KEY=your-license-key
```

### SonarCloud Code Quality

The project is configured for SonarCloud analysis:

- **Code Quality Gates**
- **Security Vulnerability Detection**
- **Code Coverage Analysis**
- **Technical Debt Monitoring**
- **Duplicated Code Detection**

Configuration in `sonar-project.properties`:
```properties
sonar.projectKey=luuphuc6297_golang-clean-architecture-sample
sonar.organization=luuphuc6297
sonar.host.url=https://sonarcloud.io
```

## 🔐 Authentication & Authorization

### JWT Authentication

The API uses JWT tokens for authentication:

- **Access Token**: Short-lived (configurable expiration)
- **Refresh Token**: Long-lived for token renewal
- **Signing Algorithm**: Configurable (HS256/RS256)

### Role-Based Access Control (RBAC)

The authorization system implements a policy-based RBAC:

#### Roles
- **`admin`**: Full system access
- **`user`**: Limited access to user resources

#### Permissions System
- **Resource-based**: Permissions tied to specific resources
- **Action-based**: CRUD operations (Create, Read, Update, Delete, List)
- **Policy Engine**: Flexible policy evaluation with conditions
- **Context-aware**: IP-based, time-based, and resource ownership checks

#### Policy Examples

**Admin Policy** (Full Access):
```json
{
  "name": "admin-full-access",
  "statements": [{
    "effect": "Allow",
    "principal": "role:admin",
    "action": "*",
    "resource": "*"
  }]
}
```

**User Policy** (Limited Access):
```json
{
  "name": "user-product-access", 
  "statements": [{
    "effect": "Allow",
    "principal": "role:user",
    "action": "create|read|update|delete|list",
    "resource": "product:*"
  }]
}
```

## 🔄 API Endpoints

### Authentication
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Register new user | ❌ |
| POST | `/api/v1/auth/login` | User login | ❌ |
| POST | `/api/v1/auth/refresh` | Refresh access token | ❌ |

### Users (Admin Only)
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/users` | List all users | ✅ (Admin) |
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

### System
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check endpoint |

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 🛠️ Development Commands

The project includes a comprehensive Makefile:

```bash
# Development
make run              # Run with PostgreSQL
make run-sqlite       # Run with SQLite  
make run-memory       # Run with in-memory DB
make dev             # Run with hot reload (requires air)

# Code Quality
make lint            # Lint code
make lint-fix        # Lint and fix issues
make format          # Format code with gofmt, goimports, etc.
make format-deps     # Install formatting dependencies

# Build & Deploy
make build           # Build binary
make docker-build    # Build Docker image
make clean           # Clean build artifacts

# Dependencies
make deps            # Install/update dependencies
```

## 🐳 Docker Deployment

### Development Environment
```bash
# Start PostgreSQL only
docker-compose up postgres -d

# Start full development stack
docker-compose up -d
```

### Production Environment
```bash
# Production deployment with optimized settings
docker-compose -f docker-compose.prod.yml up -d
```

### Environment Files
- `env.example` - PostgreSQL configuration template
- `env.sqlite.example` - SQLite configuration template

## 🌐 GCP Deployment

The application is ready for Google Cloud Platform deployment:

### Cloud Run
```bash
# Build and deploy to Cloud Run
gcloud run deploy clean-architecture-api \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

### Cloud SQL (PostgreSQL)
```bash
# Create Cloud SQL instance
gcloud sql instances create clean-architecture-db \
  --database-version=POSTGRES_15 \
  --tier=db-f1-micro \
  --region=us-central1

# Create database
gcloud sql databases create clean_architecture_api \
  --instance=clean-architecture-db
```

### Required Environment Variables for GCP
```bash
DB_HOST=<cloud-sql-connection-name>
DB_PASSWORD=<cloud-sql-password>
NEW_RELIC_LICENSE_KEY=<your-license-key>
JWT_SECRET_KEY=<production-secret>
```

## 📝 API Usage Examples

### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securePassword123",
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
    "password": "securePassword123"
  }'
```

### Access Protected Endpoint
```bash
curl -X GET http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <your-access-token>"
```

### Create Product
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <your-access-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sample Product",
    "description": "A sample product description",
    "price": 29.99,
    "category": "electronics"
  }'
```

## 🔧 Configuration

### Logging Configuration
```bash
LOG_LEVEL=debug|info|warn|error
```

### Database Connection Pooling
The application automatically configures connection pooling for optimal performance:
- **Max Open Connections**: 25
- **Max Idle Connections**: 5  
- **Connection Max Lifetime**: 30 minutes

### Security Headers
All API responses include security headers:
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`

## 🚀 Performance & Scalability

- **Stateless Design**: Fully stateless for horizontal scaling
- **Connection Pooling**: Optimized database connection management
- **Caching**: In-memory policy cache for authorization
- **Pagination**: Efficient pagination for large datasets
- **Monitoring**: Full observability with New Relic APM

## 🔒 Security Features

- **JWT Token Authentication** with configurable expiration
- **Password Hashing** using bcrypt
- **Rate Limiting** (configurable)
- **Input Validation** with custom validators
- **SQL Injection Protection** via GORM ORM
- **CORS Configuration** for cross-origin requests
- **Security Headers** on all responses
- **Audit Logging** for security events

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Quality Standards
- All code must pass `make lint`
- Test coverage should be maintained above 80%
- Follow Go best practices and Clean Architecture principles
- All commits must pass SonarCloud quality gates

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:

1. Check the [Issues](../../issues) page
2. Review the [API Testing Documentation](API_TESTING_EN.md)
3. Check application logs for error details
4. Verify environment configuration

## 📚 Additional Resources

- [Clean Architecture Principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Best Practices](https://golang.org/doc/effective_go.html)
- [Gin Framework Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [New Relic Go Agent](https://docs.newrelic.com/docs/agents/go-agent/)