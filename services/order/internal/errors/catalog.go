package errors

import (
	"errors"
	
	pkgErrors "github.com/robrt95x/godops/pkg/errors"
)

// Error codes for standardized API responses
const (
	// Order related errors
	OrderNotFound     = "ORDER_NOT_FOUND"
	OrderInvalidID    = "ORDER_INVALID_ID"
	OrderAlreadyExists = "ORDER_ALREADY_EXISTS"
	
	// Validation errors
	ValidationMissingUserID    = "VALIDATION_MISSING_USER_ID"
	ValidationEmptyItems       = "VALIDATION_EMPTY_ITEMS"
	ValidationInvalidQuantity  = "VALIDATION_INVALID_QUANTITY"
	ValidationInvalidPrice     = "VALIDATION_INVALID_PRICE"
	ValidationMissingProductID = "VALIDATION_MISSING_PRODUCT_ID"
	ValidationInvalidRequest   = "VALIDATION_INVALID_REQUEST"
	
	// Database errors
	DatabaseConnectionError = "DATABASE_CONNECTION_ERROR"
	DatabaseQueryError      = "DATABASE_QUERY_ERROR"
	DatabaseTransactionError = "DATABASE_TRANSACTION_ERROR"
	
	// System errors
	SystemInternalError = "SYSTEM_INTERNAL_ERROR"
	SystemServiceUnavailable = "SYSTEM_SERVICE_UNAVAILABLE"
	SystemTimeout = "SYSTEM_TIMEOUT"
)

// Domain errors that map to error codes
var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrOrderInvalidID    = errors.New("invalid order ID")
	ErrOrderAlreadyExists = errors.New("order already exists")
	
	ErrValidationMissingUserID    = errors.New("user ID is required")
	ErrValidationEmptyItems       = errors.New("order must contain at least one item")
	ErrValidationInvalidQuantity  = errors.New("item quantity must be greater than zero")
	ErrValidationInvalidPrice     = errors.New("item price must be greater than zero")
	ErrValidationMissingProductID = errors.New("product ID is required for all items")
	ErrValidationInvalidRequest   = errors.New("invalid request format")
	
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrDatabaseQuery      = errors.New("database query failed")
	ErrDatabaseTransaction = errors.New("database transaction failed")
	
	ErrSystemInternal = errors.New("internal system error")
	ErrSystemServiceUnavailable = errors.New("service temporarily unavailable")
	ErrSystemTimeout = errors.New("request timeout")
)

// ErrorInfo represents error information for API responses
type ErrorInfo struct {
	Code    string `json:"error_code"`
	Message string `json:"error_message"`
}

// ErrorCatalog maps domain errors to API error responses
var ErrorCatalog = map[error]ErrorInfo{
	ErrOrderNotFound:     {OrderNotFound, "The requested order could not be found"},
	ErrOrderInvalidID:    {OrderInvalidID, "Invalid order ID format"},
	ErrOrderAlreadyExists: {OrderAlreadyExists, "Order with this ID already exists"},
	
	ErrValidationMissingUserID:    {ValidationMissingUserID, "User ID is required"},
	ErrValidationEmptyItems:       {ValidationEmptyItems, "Order must contain at least one item"},
	ErrValidationInvalidQuantity:  {ValidationInvalidQuantity, "Item quantity must be greater than zero"},
	ErrValidationInvalidPrice:     {ValidationInvalidPrice, "Item price must be greater than zero"},
	ErrValidationMissingProductID: {ValidationMissingProductID, "Product ID is required for all items"},
	ErrValidationInvalidRequest:   {ValidationInvalidRequest, "Invalid request format"},
	
	ErrDatabaseConnection:  {DatabaseConnectionError, "Database connection failed"},
	ErrDatabaseQuery:       {DatabaseQueryError, "Database query failed"},
	ErrDatabaseTransaction: {DatabaseTransactionError, "Database transaction failed"},
	
	ErrSystemInternal:           {SystemInternalError, "An internal error occurred"},
	ErrSystemServiceUnavailable: {SystemServiceUnavailable, "Service is temporarily unavailable"},
	ErrSystemTimeout:            {SystemTimeout, "Request timeout"},
}

// GetErrorInfo returns the ErrorInfo for a given error
func GetErrorInfo(err error) ErrorInfo {
	if info, exists := ErrorCatalog[err]; exists {
		return info
	}
	// Default error for unknown errors
	return ErrorInfo{
		Code:    SystemInternalError,
		Message: "An unexpected error occurred",
	}
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	switch err {
	case ErrValidationMissingUserID, ErrValidationEmptyItems, ErrValidationInvalidQuantity,
		 ErrValidationInvalidPrice, ErrValidationMissingProductID, ErrValidationInvalidRequest:
		return true
	default:
		return false
	}
}

// IsDatabaseError checks if the error is a database error
func IsDatabaseError(err error) bool {
	switch err {
	case ErrDatabaseConnection, ErrDatabaseQuery, ErrDatabaseTransaction:
		return true
	default:
		return false
	}
}

// OrderErrorCatalog implements the pkgErrors.ErrorCatalog interface
type OrderErrorCatalog struct{}

// NewOrderErrorCatalog creates a new OrderErrorCatalog
func NewOrderErrorCatalog() *OrderErrorCatalog {
	return &OrderErrorCatalog{}
}

// GetErrorInfo returns the ErrorInfo for a given error
func (c *OrderErrorCatalog) GetErrorInfo(err error) pkgErrors.ErrorInfo {
	if info, exists := ErrorCatalog[err]; exists {
		return pkgErrors.ErrorInfo{
			Code:    info.Code,
			Message: info.Message,
		}
	}
	// Default error for unknown errors
	return pkgErrors.ErrorInfo{
		Code:    SystemInternalError,
		Message: "An unexpected error occurred",
	}
}

// IsValidationError checks if the error is a validation error
func (c *OrderErrorCatalog) IsValidationError(err error) bool {
	return IsValidationError(err)
}

// IsDatabaseError checks if the error is a database error
func (c *OrderErrorCatalog) IsDatabaseError(err error) bool {
	return IsDatabaseError(err)
}
