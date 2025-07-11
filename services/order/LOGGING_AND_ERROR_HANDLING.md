# Professional Logging and Error Handling Implementation

This document describes the comprehensive logging and error handling system implemented for the Order Service.

## 🎯 Overview

We've implemented a production-ready logging and error handling system with the following features:

- **Structured Logging**: JSON-formatted logs with configurable levels
- **Standardized Error Responses**: Consistent API error format with error codes
- **Request Tracing**: Unique request IDs for correlation across logs
- **Log Rotation**: Automatic log file rotation with configurable retention
- **Environment Configuration**: All logging settings configurable via .env

## 📁 Architecture

### New Components Added

```
internal/
├── errors/
│   ├── catalog.go          # Error definitions and catalog
│   └── handler.go          # HTTP error response handler
├── logger/
│   └── logger.go           # Logger configuration and setup
└── middleware/
    ├── request_id.go       # Request ID generation middleware
    └── logging.go          # HTTP request logging middleware
```

## 🔧 Error Handling System

### Error Response Format

All API errors now return a standardized JSON format:

```json
{
  "error_code": "ORDER_NOT_FOUND",
  "error_message": "The requested order could not be found"
}
```

### Error Categories

**Order Errors:**
- `ORDER_NOT_FOUND` - Order doesn't exist
- `ORDER_INVALID_ID` - Invalid order ID format
- `ORDER_ALREADY_EXISTS` - Duplicate order ID

**Validation Errors:**
- `VALIDATION_MISSING_USER_ID` - User ID required
- `VALIDATION_EMPTY_ITEMS` - Order must have items
- `VALIDATION_INVALID_QUANTITY` - Invalid item quantity
- `VALIDATION_INVALID_PRICE` - Invalid item price
- `VALIDATION_MISSING_PRODUCT_ID` - Product ID required

**Database Errors:**
- `DATABASE_CONNECTION_ERROR` - Connection failed
- `DATABASE_QUERY_ERROR` - Query execution failed
- `DATABASE_TRANSACTION_ERROR` - Transaction failed

**System Errors:**
- `SYSTEM_INTERNAL_ERROR` - Generic internal error
- `SYSTEM_SERVICE_UNAVAILABLE` - Service unavailable
- `SYSTEM_TIMEOUT` - Request timeout

### HTTP Status Code Mapping

- **400 Bad Request**: Validation errors, invalid input
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource already exists
- **408 Request Timeout**: Timeout errors
- **500 Internal Server Error**: Database and system errors
- **503 Service Unavailable**: Service unavailable

## 📊 Logging System

### Log Levels

- **DEBUG**: Detailed flow information, variable values
- **INFO**: General application flow, successful operations
- **WARNING**: Unexpected situations that don't stop operation
- **ERROR**: Error conditions affecting specific operations
- **FATAL**: Critical errors causing application shutdown

### Log Format

**JSON Format (Production):**
```json
{
  "timestamp": "2025-01-10T20:07:06.084-06:00",
  "level": "INFO",
  "service": "order-service",
  "message": "Order created successfully",
  "request_id": "req-123e4567-e89b-12d3-a456-426614174000",
  "user_id": "user-456",
  "order_id": "order-789",
  "duration_ms": 45
}
```

**Text Format (Development):**
```
2025-01-10 20:07:06 INFO Order created successfully request_id=req-123... user_id=user-456
```

### Configuration Options

All logging settings are configurable via environment variables:

```env
# Log levels: DEBUG, INFO, WARNING, ERROR
LOG_LEVEL=INFO

# Log formats: json, text
LOG_FORMAT=json

# Log outputs: console, file, both
LOG_OUTPUT=console

# Log file settings (when LOG_OUTPUT=file or both)
LOG_FILE_PATH=logs/order-service.log
LOG_MAX_SIZE=100        # MB
LOG_MAX_BACKUPS=5       # number of backup files
LOG_MAX_AGE=30          # days to retain
LOG_COMPRESS=true       # compress old files
```

## 🔄 Request Flow with Logging

### 1. Request Arrives
- **Request ID Middleware**: Generates unique ID for request
- **Logging Middleware**: Logs request start with details

### 2. Handler Processing
- **Handler Logs**: Debug/Info logs with request context
- **Use Case Logs**: Business logic logging with structured fields
- **Repository Logs**: Data access logging (if needed)

### 3. Response
- **Success**: Info log with response details
- **Error**: Warning/Error log with error details
- **Logging Middleware**: Logs request completion with duration

### 4. Error Handling
- **Domain Errors**: Mapped to appropriate HTTP status and error code
- **Unexpected Errors**: Logged with full details, generic error returned to client

## 🚀 Usage Examples

### In Use Cases
```go
logEntry := uc.logger.WithFields(logrus.Fields{
    "use_case": "CreateOrder",
    "user_id":  userID,
    "items_count": len(items),
})

logEntry.Debug("Starting create order use case")

if userID == "" {
    logEntry.Warning("Create order failed: missing user ID")
    return nil, customErrors.ErrValidationMissingUserID
}

logEntry.Info("Order created successfully")
```

### In HTTP Handlers
```go
logEntry := h.Logger.WithFields(logrus.Fields{
    "handler":    "CreateOrder",
    "request_id": requestID,
})

order, err := h.CreateUC.Execute(req.UserID, req.Items)
if err != nil {
    logEntry.WithError(err).Error("Create order use case failed")
    h.ErrorHandler.HandleError(w, r, err)
    return
}

logEntry.WithField("order_id", order.ID).Info("Order created successfully")
```

## 🧪 Testing

### Unit Tests
- All tests pass with new error handling
- Memory repository used for fast testing
- Structured logging visible in test output

### Integration Testing
```bash
# Run tests with verbose logging
go test ./internal/usecase/... -v

# Test with different log levels
LOG_LEVEL=DEBUG go test ./internal/usecase/... -v
```

## 📈 Benefits

### For Development
- **Better Debugging**: Structured logs with request correlation
- **Fast Testing**: Memory repository for unit tests
- **Clear Errors**: Specific error codes and messages

### For Production
- **Observability**: JSON logs for log aggregation tools
- **Monitoring**: Structured fields for metrics and alerts
- **Troubleshooting**: Request IDs for tracing issues
- **Compliance**: Audit trails with detailed logging

### For Operations
- **Log Management**: Automatic rotation and retention
- **Performance**: Efficient logging with minimal overhead
- **Scalability**: Configurable log levels and outputs

## 🔧 Configuration Examples

### Development Environment
```env
STORAGE_TYPE=memory
LOG_LEVEL=DEBUG
LOG_FORMAT=text
LOG_OUTPUT=console
APP_ENV=development
```

### Production Environment
```env
STORAGE_TYPE=postgres
LOG_LEVEL=INFO
LOG_FORMAT=json
LOG_OUTPUT=both
LOG_FILE_PATH=/var/log/order-service/app.log
LOG_MAX_SIZE=100
LOG_MAX_BACKUPS=10
LOG_MAX_AGE=30
LOG_COMPRESS=true
APP_ENV=production
```

## 🎉 Summary

The Order Service now features enterprise-grade logging and error handling:

✅ **Structured JSON Logging** with configurable levels and outputs  
✅ **Standardized Error Responses** with consistent format and codes  
✅ **Request Tracing** with unique IDs for correlation  
✅ **Log Rotation** with automatic file management  
✅ **Environment Configuration** for different deployment scenarios  
✅ **Production Ready** with proper error handling and observability  
✅ **Developer Friendly** with clear error messages and debug logging  
✅ **Test Coverage** with comprehensive unit tests  

The system is now ready for production deployment with full observability and professional error handling.
