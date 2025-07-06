# Base Code Go Gin Clean

A clean architecture template for Go web applications using Gin framework and Bun ORM.

## 🚀 Features

- **Clean Architecture** - Well-structured project following clean architecture principles
- **Gin Web Framework** - Fast and efficient HTTP web framework
- **Bun ORM** - SQL-first Golang ORM with PostgreSQL support
- **Database Migrations** - Built-in migration system using goose
- **Environment Configuration** - Easy environment-based configuration
- **Structured Logging** - JSON-formatted logs with different log levels
- **Graceful Shutdown** - Proper handling of server shutdown
- **Health Check** - Built-in health check endpoint
- **Authentication** - Secure user registration and login with bcrypt password hashing
- **User Management** - Basic user CRUD operations
- **API Documentation** - Swagger/OpenAPI documentation
- **Docker Support** - Easy containerization

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/base-code-go-gin-clean.git
   cd base-code-go-gin-clean
   ```

2. Copy environment file and update the values:

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

### ⚙️ Configuration

Update the `.env` file with your configuration:

```env
# Server Configuration
PORT=8080
ENVIRONMENT=development  # development, staging, production

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=yourdb
DB_SSLMODE=disable  # disable, require, verify-full

# Optional: Uptrace APM (set ENABLED=true to enable)
TRACING_ENABLED=false
SERVICE_NAME=base-code-go-gin-clean
SERVICE_VERSION=1.0.0
UPTRACE_DSN=
```

### 🏃 Running the Application

#### Using Make (recommended):

```bash
# Start the development server (with auto-reload)
make dev

# Run database migrations
make migrate

# Run database seeders
make seed

# Build the application
make build
```

#### Using Go commands:

```bash
# Run the application
go run main.go

# Run database migrations
cd cmd/migrate && go run .

# Run database seeders
cd cmd/seed && go run .
```

### 🐳 Using Docker

```bash
# Build and start containers
docker-compose up --build

# Run migrations in the container
docker-compose exec app make migrate

# Run seeders in the container
docker-compose exec app make seed
```

## 📚 API Documentation

Once the application is running, you can access:

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health
- **Ping Endpoint**: http://localhost:8080/api/v1/ping

## API Endpoints

### Health Check

- `GET /health` - Check if the server is running

### Authentication

#### Register a new user
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

### Users

- `GET /api/v1/users/:id` - Get user by ID

## 📂 Project Structure

```
.
# Main application entry points
├── cmd/                  
│   ├── checkdb/         # Database connection checker
│   ├── migrate/         # Database migration tool
│   └── seed/            # Database seeder

# Application code (private)
├── internal/            
│   ├── config/          # Configuration management
│   ├── domain/          # Core business models and interfaces
│   ├── handler/         # HTTP request handlers
│   ├── pkg/             # Internal shared packages
│   ├── repository/      # Data access layer
│   ├── routes/          # Route definitions
│   ├── seeders/         # Database seeders
│   ├── server/          # Server configuration and setup
│   └── service/         # Business logic layer

# Public packages
└── pkg/
    ├── logger/          # Structured logging
    └── middleware/      # HTTP middleware components
```

## License

MIT
