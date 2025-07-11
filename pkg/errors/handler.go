package errors

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// ErrorInfo represents error information for API responses
type ErrorInfo struct {
	Code    string `json:"error_code"`
	Message string `json:"error_message"`
}

// ErrorCatalog interface that each service should implement
type ErrorCatalog interface {
	GetErrorInfo(err error) ErrorInfo
	IsValidationError(err error) bool
	IsDatabaseError(err error) bool
}

// HTTPErrorHandler handles HTTP error responses with standardized format
type HTTPErrorHandler struct {
	logger  *logrus.Logger
	catalog ErrorCatalog
}

// NewHTTPErrorHandler creates a new HTTP error handler
func NewHTTPErrorHandler(logger *logrus.Logger, catalog ErrorCatalog) *HTTPErrorHandler {
	return &HTTPErrorHandler{
		logger:  logger,
		catalog: catalog,
	}
}

// HandleError processes an error and sends appropriate HTTP response
func (h *HTTPErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	errorInfo := h.catalog.GetErrorInfo(err)
	statusCode := h.getHTTPStatusCode(err)
	
	// Log the error with context
	logEntry := h.logger.WithFields(logrus.Fields{
		"error_code":    errorInfo.Code,
		"error_message": errorInfo.Message,
		"status_code":   statusCode,
		"method":        r.Method,
		"path":          r.URL.Path,
		"user_agent":    r.Header.Get("User-Agent"),
	})
	
	// Add request ID if available
	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		logEntry = logEntry.WithField("request_id", requestID)
	}
	
	// Log with appropriate level based on error type
	switch {
	case h.catalog.IsValidationError(err):
		logEntry.Warning("Validation error occurred")
	case h.catalog.IsDatabaseError(err):
		logEntry.Error("Database error occurred")
	case statusCode >= 500:
		logEntry.Error("Internal server error occurred")
	default:
		logEntry.Info("Request completed with error")
	}
	
	// Send standardized error response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if encodeErr := json.NewEncoder(w).Encode(errorInfo); encodeErr != nil {
		h.logger.WithFields(logrus.Fields{
			"original_error": err.Error(),
			"encode_error":   encodeErr.Error(),
		}).Error("Failed to encode error response")
	}
}

// getHTTPStatusCode maps domain errors to HTTP status codes
// This uses a generic approach that services can override if needed
func (h *HTTPErrorHandler) getHTTPStatusCode(err error) int {
	// Check if it's a validation error
	if h.catalog.IsValidationError(err) {
		return http.StatusBadRequest
	}
	
	// Check if it's a database error
	if h.catalog.IsDatabaseError(err) {
		return http.StatusInternalServerError
	}
	
	// Get error info to determine status code based on error code pattern
	errorInfo := h.catalog.GetErrorInfo(err)
	
	switch {
	case contains(errorInfo.Code, "NOT_FOUND"):
		return http.StatusNotFound
	case contains(errorInfo.Code, "ALREADY_EXISTS"):
		return http.StatusConflict
	case contains(errorInfo.Code, "TIMEOUT"):
		return http.StatusRequestTimeout
	case contains(errorInfo.Code, "SERVICE_UNAVAILABLE"):
		return http.StatusServiceUnavailable
	case contains(errorInfo.Code, "VALIDATION_"):
		return http.StatusBadRequest
	case contains(errorInfo.Code, "DATABASE_"):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// HandleValidationError is a convenience method for validation errors
func (h *HTTPErrorHandler) HandleValidationError(w http.ResponseWriter, r *http.Request, message string) {
	if message != "" {
		// Create a custom validation error with specific message
		errorInfo := ErrorInfo{
			Code:    "VALIDATION_INVALID_REQUEST",
			Message: message,
		}
		
		logEntry := h.logger.WithFields(logrus.Fields{
			"error_code":    errorInfo.Code,
			"error_message": errorInfo.Message,
			"status_code":   http.StatusBadRequest,
			"method":        r.Method,
			"path":          r.URL.Path,
		})
		
		if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
			logEntry = logEntry.WithField("request_id", requestID)
		}
		
		logEntry.Warning("Validation error occurred")
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorInfo)
		return
	}
	
	// Fallback to generic validation error
	errorInfo := ErrorInfo{
		Code:    "VALIDATION_INVALID_REQUEST",
		Message: "Invalid request format",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errorInfo)
}

// HandleInternalError is a convenience method for internal server errors
func (h *HTTPErrorHandler) HandleInternalError(w http.ResponseWriter, r *http.Request, err error) {
	// Log the original error for debugging
	logEntry := h.logger.WithFields(logrus.Fields{
		"original_error": err.Error(),
		"method":         r.Method,
		"path":           r.URL.Path,
	})
	
	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		logEntry = logEntry.WithField("request_id", requestID)
	}
	
	logEntry.Error("Internal server error occurred")
	
	// Return generic internal error to client
	errorInfo := ErrorInfo{
		Code:    "SYSTEM_INTERNAL_ERROR",
		Message: "An internal error occurred",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errorInfo)
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		   len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
		   (len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
