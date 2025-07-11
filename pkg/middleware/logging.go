package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// Logging middleware logs HTTP requests with structured logging
func Logging(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Wrap response writer to capture status code
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     0,
			}
			
			// Get request ID from context
			requestID := GetRequestID(r)
			
			// Create log entry with request details
			logEntry := logger.WithFields(logrus.Fields{
				"request_id":   requestID,
				"method":       r.Method,
				"path":         r.URL.Path,
				"query":        r.URL.RawQuery,
				"remote_addr":  r.RemoteAddr,
				"user_agent":   r.Header.Get("User-Agent"),
				"content_type": r.Header.Get("Content-Type"),
			})
			
			// Log request start
			logEntry.Info("Request started")
			
			// Process request
			next.ServeHTTP(wrapped, r)
			
			// Calculate duration
			duration := time.Since(start)
			
			// Create response log entry
			responseEntry := logEntry.WithFields(logrus.Fields{
				"status_code":  wrapped.statusCode,
				"duration_ms":  duration.Milliseconds(),
				"response_size": wrapped.written,
			})
			
			// Log with appropriate level based on status code
			switch {
			case wrapped.statusCode >= 500:
				responseEntry.Error("Request completed with server error")
			case wrapped.statusCode >= 400:
				responseEntry.Warning("Request completed with client error")
			case wrapped.statusCode >= 300:
				responseEntry.Info("Request completed with redirect")
			default:
				responseEntry.Info("Request completed successfully")
			}
		})
	}
}

// LoggingWithSkipPaths creates a logging middleware that skips certain paths
func LoggingWithSkipPaths(logger *logrus.Logger, skipPaths []string) func(next http.Handler) http.Handler {
	skipMap := make(map[string]bool)
	for _, path := range skipPaths {
		skipMap[path] = true
	}
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip logging for specified paths (e.g., health checks)
			if skipMap[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}
			
			// Use regular logging middleware
			Logging(logger)(next).ServeHTTP(w, r)
		})
	}
}

// ErrorLogging middleware logs panics and recovers
func ErrorLogging(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					requestID := GetRequestID(r)
					
					logger.WithFields(logrus.Fields{
						"request_id": requestID,
						"method":     r.Method,
						"path":       r.URL.Path,
						"panic":      err,
					}).Error("Panic recovered in HTTP handler")
					
					// Return 500 error
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			
			next.ServeHTTP(w, r)
		})
	}
}
