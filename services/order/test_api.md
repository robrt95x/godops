# Order Service API Testing

## Configuration

The service now supports multiple storage backends via environment configuration:

### Storage Options
- **Memory Storage**: Fast, in-memory storage perfect for testing and development
- **PostgreSQL**: Production-ready persistent storage

### Environment Configuration
Copy `.env.example` to `.env` and configure:

```bash
# For testing with memory storage
cp .env.development .env

# For production with PostgreSQL
cp .env.example .env
# Then edit .env with your database credentials
```

## Endpoints

### 1. Create Order
```bash
POST /orders
Content-Type: application/json

{
  "user_id": "user123",
  "items": [
    {
      "product_id": "product1",
      "quantity": 2,
      "price": 29.99
    },
    {
      "product_id": "product2", 
      "quantity": 1,
      "price": 15.50
    }
  ]
}
```

### 2. Get Order by ID
```bash
GET /orders/{id}
```

## Example Usage

### Quick Testing with Memory Storage

1. Use memory storage (no database required):
```bash
cp .env.development .env
./order-service
```

2. Create an order:
```bash
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
```

3. Get the order by ID (use the ID returned from step 2):
```bash
curl http://localhost:8080/orders/{order-id}
```

### Production with PostgreSQL

1. Configure PostgreSQL:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

2. Start the server:
```bash
./order-service
```

## Response Codes

- **200 OK**: Order found and returned successfully
- **400 Bad Request**: Invalid order ID format
- **404 Not Found**: Order not found
- **500 Internal Server Error**: Database or server error

## Features Added

- ✅ Chi router integration with middleware (Logger, Recoverer)
- ✅ RESTful routing structure (`/orders` for POST, `/orders/{id}` for GET)
- ✅ Clean architecture with separate use case for GetOrderByID
- ✅ Proper error handling with appropriate HTTP status codes
- ✅ URL parameter extraction using Chi's URLParam
