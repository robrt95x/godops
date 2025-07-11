package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"
const RequestIDContextKey = "request_id"

// RequestID middleware generates a unique request ID for each request
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request ID already exists in header
		requestID := r.Header.Get(RequestIDHeader)
		
		// Generate new request ID if not provided
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// Add request ID to response header
		w.Header().Set(RequestIDHeader, requestID)
		
		// Add request ID to request header for downstream processing
		r.Header.Set(RequestIDHeader, requestID)
		
		// Add request ID to context
		ctx := context.WithValue(r.Context(), RequestIDContextKey, requestID)
		r = r.WithContext(ctx)
		
		next.ServeHTTP(w, r)
	})
}

// GetRequestID extracts request ID from context
func GetRequestID(r *http.Request) string {
	if requestID, ok := r.Context().Value(RequestIDContextKey).(string); ok {
		return requestID
	}
	return ""
}

// GetRequestIDFromContext extracts request ID from context
func GetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDContextKey).(string); ok {
		return requestID
	}
	return ""
}
