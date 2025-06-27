# Base Code Go Gin Clean

A clean architecture template for Go web applications using Gin framework.

## Features

- Clean Architecture implementation
- Gin web framework
- PostgreSQL database with Bun ORM
- Environment variable configuration
- Structured logging
- Graceful shutdown
- Health check endpoint
- User management (CRUD)

## Getting Started

### Prerequisites

- Go 1.20 or higher
- PostgreSQL 13 or higher
- Make (optional)

### Installation

1. Clone the repository
2. Copy `.env.example` to `.env` and update the values
3. Install dependencies:
   ```bash
   go mod download
   ```

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
PORT=8080
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=yourdb
DB_SSLMODE=disable
```

### Running the Application

```bash
# Run the application
go run main.go

# Or using make
make run
```

## API Endpoints

### Users

- `GET /api/v1/users/:id` - Get user by ID

### Health Check

- `GET /health` - Health check endpoint
- `GET /api/v1/ping` - Ping endpoint

## Project Structure

```
.
├── cmd/                  # Main applications for this project
├── internal/             # Private application and library code
│   ├── config/          # Configuration
│   ├── domain/          # Enterprise business rules
│   ├── handler/         # HTTP handlers
│   ├── repository/      # Data access layer
│   ├── server/          # HTTP server configuration
│   └── service/         # Business logic
├── pkg/                 # Library code that's ok to use by external applications
└── scripts/             # Build and deployment scripts
```

## License

MIT
