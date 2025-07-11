# Order Service

A microservice for managing orders with support for multiple storage backends and RESTful API endpoints.

## Features

- **RESTful API**: Create and retrieve orders via HTTP endpoints
- **Multiple Storage Backends**: Memory (for testing) and PostgreSQL (for production)
- **Configuration Management**: Environment-based configuration with .env support
- **Clean Architecture**: Separation of concerns with use cases, repositories, and handlers
- **Chi Router**: Modern HTTP router with middleware support
- **Unit Testing**: Comprehensive tests with memory repository for fast testing

## Architecture

```
cmd/
├── main.go                 # Application entry point

internal/
├── config/
│   └── config.go          # Configuration management
├── delivery/
│   └── http/
│       └── handler.go     # HTTP handlers
├── entity/
│   └── order.go          # Domain entities
├── infra/
│   ├── factory.go        # Repository factory
│   ├── memory/
│   │   └── order_repository.go  # In-memory implementation
│   └── postgres/
│       └── order_repository.go  # PostgreSQL implementation
├── repository/
│   └── order_repository.go      # Repository interface
└── usecase/
    ├── create_order.go          # Create order use case
    ├── get_order_by_id.go       # Get order by ID use case
    └── get_order_by_id_test.go  # Unit tests
```

## API Endpoints

### Create Order
```http
POST /orders
Content-Type: application/json

{
  "user_id": "user123",
  "items": [
    {
      "product_id": "product1",
      "quantity": 2,
      "price": 29.99
    }
  ]
}
```

### Get Order by ID
```http
GET /orders/{id}
```

## Configuration

The service supports environment-based configuration via `.env` files:

### Environment Variables

| Variable | Description | Default | Options |
|----------|-------------|---------|---------|
| `STORAGE_TYPE` | Storage backend | `postgres` | `memory`, `postgres` |
| `DB_HOST` | Database host | `localhost` | - |
| `DB_PORT` | Database port | `5432` | - |
| `DB_USER` | Database user | `user` | - |
| `DB_PASSWORD` | Database password | `pass` | - |
| `DB_NAME` | Database name | `godops` | - |
| `DB_SSLMODE` | SSL mode | `disable` | - |
| `SERVER_PORT` | Server port | `8080` | - |
| `LOG_LEVEL` | Log level | `info` | - |
| `APP_ENV` | Environment | `development` | `development`, `production`, `test` |

## Quick Start

### Development (Memory Storage)
```bash
# Copy development configuration
cp .env.development .env

# Build and run
go build -o order-service ./cmd/main.go
./order-service
```

### Production (PostgreSQL)
```bash
# Copy and configure production settings
cp .env.example .env
# Edit .env with your database credentials

# Build and run
go build -o order-service ./cmd/main.go
./order-service
```

## Testing

### Run Unit Tests
```bash
go test ./internal/usecase/... -v
```

### Test API with Memory Storage
```bash
# Start server with memory storage
cp .env.development .env
./order-service

# Create an order
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "items": [
      {
        "product_id": "product1",
        "quantity": 2,
        "price": 29.99
      }
    ]
  }'

# Get order by ID (replace {id} with actual order ID from create response)
curl http://localhost:8080/orders/{id}
```

## Dependencies

- **Chi Router**: HTTP router and middleware
- **godotenv**: Environment configuration
- **PostgreSQL Driver**: Database connectivity
- **UUID**: Unique identifier generation

## Development

### Adding New Features
1. Define new use cases in `internal/usecase/`
2. Add repository methods if needed in `internal/repository/`
3. Implement in both memory and postgres repositories
4. Add HTTP handlers in `internal/delivery/http/`
5. Update routes in `cmd/main.go`
6. Write unit tests

### Storage Backends
- **Memory**: Perfect for unit testing, development, and CI/CD pipelines
- **PostgreSQL**: Production-ready with ACID compliance and persistence

The factory pattern makes it easy to add new storage backends by implementing the `OrderRepository` interface.
