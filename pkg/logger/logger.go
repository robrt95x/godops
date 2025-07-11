package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config holds logger configuration
type Config struct {
	Level       string
	Format      string
	Output      string
	FilePath    string
	MaxSize     int
	MaxBackups  int
	MaxAge      int
	Compress    bool
	ServiceName string
}

// Setup initializes and configures the logger
func Setup(config Config) *logrus.Logger {
	logger := logrus.New()
	
	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
		logger.Warnf("Invalid log level '%s', defaulting to INFO", config.Level)
	}
	logger.SetLevel(level)
	
	// Set log format
	switch strings.ToLower(config.Format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		// Default to JSON for production
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	}
	
	// Set output destination
	var output io.Writer
	switch strings.ToLower(config.Output) {
	case "console":
		output = os.Stdout
	case "file":
		output = setupFileOutput(config)
	case "both":
		fileOutput := setupFileOutput(config)
		output = io.MultiWriter(os.Stdout, fileOutput)
	default:
		output = os.Stdout
	}
	
	logger.SetOutput(output)
	
	// Add service name as a default field
	if config.ServiceName != "" {
		logger = logger.WithField("service", config.ServiceName).Logger
	}
	
	return logger
}

// setupFileOutput configures file output with rotation
func setupFileOutput(config Config) io.Writer {
	// Ensure log directory exists
	if config.FilePath != "" {
		dir := filepath.Dir(config.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			logrus.Errorf("Failed to create log directory: %v", err)
			return os.Stdout
		}
	}
	
	return &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize,    // MB
		MaxBackups: config.MaxBackups, // number of backups
		MaxAge:     config.MaxAge,     // days
		Compress:   config.Compress,   // compress old files
	}
}

// NewDefaultConfig returns a default logger configuration
func NewDefaultConfig() Config {
	return Config{
		Level:       "info",
		Format:      "json",
		Output:      "console",
		FilePath:    "logs/service.log",
		MaxSize:     100, // 100 MB
		MaxBackups:  5,
		MaxAge:      30, // 30 days
		Compress:    true,
		ServiceName: "service",
	}
}

// WithRequestID adds a request ID field to the logger
func WithRequestID(logger *logrus.Logger, requestID string) *logrus.Entry {
	return logger.WithField("request_id", requestID)
}

// WithUserID adds a user ID field to the logger
func WithUserID(logger *logrus.Entry, userID string) *logrus.Entry {
	return logger.WithField("user_id", userID)
}

// WithDuration adds a duration field to the logger
func WithDuration(logger *logrus.Entry, durationMs int64) *logrus.Entry {
	return logger.WithField("duration_ms", durationMs)
}

// WithHTTPFields adds common HTTP fields to the logger
func WithHTTPFields(logger *logrus.Logger, method, path string, statusCode int) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"method":      method,
		"path":        path,
		"status_code": statusCode,
	})
}
