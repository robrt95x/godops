# Shared Package (pkg)

This directory contains shared components that can be used across all services in the monorepo. The shared components provide consistent logging, error handling, and middleware functionality.

## 🎯 Overview

The `pkg` directory contains reusable components that promote consistency and reduce code duplication across services:

- **Logger**: Structured logging with configurable levels, formats, and outputs
- **Error Handling**: Generic HTTP error handler that works with service-specific error catalogs
- **Middleware**: Request ID generation and HTTP request logging middleware

## 📁 Structure

```
pkg/
├── go.mod                    # Module dependencies
├── errors/
│   └── handler.go           # Generic HTTP error handler
├── logger/
│   └── logger.go            # Logger configuration and setup
└── middleware/
    ├── request_id.go        # Request ID generation middleware
    └── logging.go           # HTTP request logging middleware
```

## 🔧 Components

### Logger (`pkg/logger`)

Provides structured logging with configurable options:

```go
import pkgLogger "github.com/robrt95x/godops/pkg/logger"

// Setup logger
config := pkgLogger.Config{
    Level:       "info",
    Format:      "json",
    Output:      "console",
    FilePath:    "logs/service.log",
    MaxSize:     100,
    MaxBackups:  5,
    MaxAge:      30,
    Compress:    true,
    ServiceName: "my-service",
}
logger := pkgLogger.Setup(config)
```

**Features:**
- JSON and text formatting
- Console, file, or both outputs
- Automatic log rotation with Lumberjack
- Configurable retention policies
- Service name tagging

### Error Handler (`pkg/errors`)

Generic HTTP error handler that works with service-specific error catalogs:

```go
import pkgErrors "github.com/robrt95x/godops/pkg/errors"

// Service must implement ErrorCatalog interface
type MyErrorCatalog struct{}

func (c *MyErrorCatalog) GetErrorInfo(err error) pkgErrors.ErrorInfo { /* ... */ }
func (c *MyErrorCatalog) IsValidationError(err error) bool { /* ... */ }
func (c *MyErrorCatalog) IsDatabaseError(err error) bool { /* ... */ }

// Create error handler
catalog := &MyErrorCatalog{}
errorHandler := pkgErrors.NewHTTPErrorHandler(logger, catalog)

// Use in HTTP handlers
errorHandler.HandleError(w, r, err)
```

**Features:**
- Standardized JSON error responses
- Automatic HTTP status code mapping
- Structured error logging
- Service-specific error catalog integration

### Middleware (`pkg/middleware`)

HTTP middleware for request processing:

```go
import pkgMiddleware "github.com/robrt95x/godops/pkg/middleware"

// Setup router with middleware
r := chi.NewRouter()
r.Use(pkgMiddleware.RequestID)
r.Use(pkgMiddleware.Logging(logger))
r.Use(pkgMiddleware.ErrorLogging(logger))
```

**Features:**
- **RequestID**: Generates unique request IDs for tracing
- **Logging**: Structured HTTP request/response logging
- **ErrorLogging**: Panic recovery with logging

## 🚀 Usage in Services

### 1. Add Dependency

Add the pkg dependency to your service's `go.mod`:

```go
require (
    github.com/robrt95x/godops/pkg v0.0.0-00010101000000-000000000000
)

replace github.com/robrt95x/godops/pkg => ../../pkg
```

### 2. Implement Error Catalog

Each service should implement the `ErrorCatalog` interface:

```go
package errors

import pkgErrors "github.com/robrt95x/godops/pkg/errors"

// Service-specific error catalog
type ServiceErrorCatalog struct{}

func NewServiceErrorCatalog() *ServiceErrorCatalog {
    return &ServiceErrorCatalog{}
}

func (c *ServiceErrorCatalog) GetErrorInfo(err error) pkgErrors.ErrorInfo {
    // Map service errors to ErrorInfo
    if info, exists := ServiceErrorCatalog[err]; exists {
        return pkgErrors.ErrorInfo{
            Code:    info.Code,
            Message: info.Message,
        }
    }
    return pkgErrors.ErrorInfo{
        Code:    "SYSTEM_INTERNAL_ERROR",
        Message: "An unexpected error occurred",
    }
}

func (c *ServiceErrorCatalog) IsValidationError(err error) bool {
    // Check if error is validation-related
}

func (c *ServiceErrorCatalog) IsDatabaseError(err error) bool {
    // Check if error is database-related
}
```

### 3. Setup in Main

Configure shared components in your service's main function:

```go
package main

import (
    pkgLogger "github.com/robrt95x/godops/pkg/logger"
    pkgMiddleware "github.com/robrt95x/godops/pkg/middleware"
    "your-service/internal/errors"
)

func main() {
    // Setup logger
    loggerConfig := pkgLogger.Config{
        Level:       cfg.LogLevel,
        Format:      cfg.LogFormat,
        Output:      cfg.LogOutput,
        FilePath:    cfg.LogFilePath,
        MaxSize:     cfg.LogMaxSize,
        MaxBackups:  cfg.LogMaxBackups,
        MaxAge:      cfg.LogMaxAge,
        Compress:    cfg.LogCompress,
        ServiceName: "your-service",
    }
    logger := pkgLogger.Setup(loggerConfig)
    
    // Setup error handling
    errorCatalog := errors.NewServiceErrorCatalog()
    
    // Create handlers with shared components
    handler := NewHandler(useCase, logger, errorCatalog)
    
    // Setup middleware
    r := chi.NewRouter()
    r.Use(pkgMiddleware.RequestID)
    r.Use(pkgMiddleware.Logging(logger))
    r.Use(pkgMiddleware.ErrorLogging(logger))
}
```

### 4. Use in Handlers

Use shared error handling in HTTP handlers:

```go
package http

import (
    pkgErrors "github.com/robrt95x/godops/pkg/errors"
    "github.com/sirupsen/logrus"
)

type Handler struct {
    useCase      UseCase
    errorHandler *pkgErrors.HTTPErrorHandler
    logger       *logrus.Logger
}

func NewHandler(useCase UseCase, logger *logrus.Logger, catalog pkgErrors.ErrorCatalog) *Handler {
    return &Handler{
        useCase:      useCase,
        errorHandler: pkgErrors.NewHTTPErrorHandler(logger, catalog),
        logger:       logger,
    }
}

func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
    result, err := h.useCase.Execute(data)
    if err != nil {
        h.errorHandler.HandleError(w, r, err)
        return
    }
    
    // Success response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
```

## 📋 Configuration

### Environment Variables

Services should support these logging configuration variables:

```env
# Log levels: DEBUG, INFO, WARNING, ERROR
LOG_LEVEL=INFO

# Log formats: json, text
LOG_FORMAT=json

# Log outputs: console, file, both
LOG_OUTPUT=console

# Log file settings (when LOG_OUTPUT=file or both)
LOG_FILE_PATH=logs/service.log
LOG_MAX_SIZE=100        # MB
LOG_MAX_BACKUPS=5       # number of backup files
LOG_MAX_AGE=30          # days to retain
LOG_COMPRESS=true       # compress old files
```

## 🎯 Benefits

### Consistency
- Standardized logging format across all services
- Consistent error response format
- Uniform request tracing

### Maintainability
- Single source of truth for shared functionality
- Easy to update logging/error handling across all services
- Reduced code duplication

### Observability
- Structured logs for easy parsing and analysis
- Request correlation with unique IDs
- Standardized error codes for monitoring

### Developer Experience
- Simple integration with existing services
- Clear interfaces and documentation
- Flexible configuration options

## 🔄 Migration Guide

To migrate an existing service to use shared components:

1. **Add pkg dependency** to service's `go.mod`
2. **Implement ErrorCatalog interface** for service-specific errors
3. **Update imports** to use pkg components
4. **Replace service-specific** logging/middleware with pkg versions
5. **Update configuration** to support new logging options
6. **Remove old** service-specific logging/middleware files
7. **Test thoroughly** to ensure functionality is preserved

## 📝 Example Services

See the `services/order` directory for a complete example of how to integrate the shared components.

## 🤝 Contributing

When adding new shared functionality:

1. Ensure it's truly generic and reusable across services
2. Maintain backward compatibility when possible
3. Update documentation and examples
4. Test with multiple services
5. Follow established patterns and interfaces

The shared package should remain focused on truly common functionality that benefits all services in the monorepo.
