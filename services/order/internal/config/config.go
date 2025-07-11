package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Storage Configuration
	StorageType string `env:"STORAGE_TYPE" default:"postgres"`
	
	// Database Configuration
	DBHost     string `env:"DB_HOST" default:"localhost"`
	DBPort     string `env:"DB_PORT" default:"5432"`
	DBUser     string `env:"DB_USER" default:"user"`
	DBPassword string `env:"DB_PASSWORD" default:"pass"`
	DBName     string `env:"DB_NAME" default:"godops"`
	DBSSLMode  string `env:"DB_SSLMODE" default:"disable"`
	
	// Server Configuration
	ServerPort string `env:"SERVER_PORT" default:"8080"`
	
	// Logging Configuration
	LogLevel       string `env:"LOG_LEVEL" default:"info"`
	LogFormat      string `env:"LOG_FORMAT" default:"json"`
	LogOutput      string `env:"LOG_OUTPUT" default:"console"`
	LogFilePath    string `env:"LOG_FILE_PATH" default:"logs/order-service.log"`
	LogMaxSize     int    `env:"LOG_MAX_SIZE" default:"100"`
	LogMaxBackups  int    `env:"LOG_MAX_BACKUPS" default:"5"`
	LogMaxAge      int    `env:"LOG_MAX_AGE" default:"30"`
	LogCompress    bool   `env:"LOG_COMPRESS" default:"true"`
	
	// Environment
	AppEnv string `env:"APP_ENV" default:"development"`
}

func Load() *Config {
	// Try to load .env file (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables and defaults")
	}
	
	config := &Config{
		StorageType:    getEnv("STORAGE_TYPE", "postgres"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "user"),
		DBPassword:     getEnv("DB_PASSWORD", "pass"),
		DBName:         getEnv("DB_NAME", "godops"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		LogFormat:      getEnv("LOG_FORMAT", "json"),
		LogOutput:      getEnv("LOG_OUTPUT", "console"),
		LogFilePath:    getEnv("LOG_FILE_PATH", "logs/order-service.log"),
		LogMaxSize:     getEnvInt("LOG_MAX_SIZE", 100),
		LogMaxBackups:  getEnvInt("LOG_MAX_BACKUPS", 5),
		LogMaxAge:      getEnvInt("LOG_MAX_AGE", 30),
		LogCompress:    getEnvBool("LOG_COMPRESS", true),
		AppEnv:         getEnv("APP_ENV", "development"),
	}
	
	return config
}

func (c *Config) GetDatabaseURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=" + c.DBSSLMode
}

func (c *Config) IsMemoryStorage() bool {
	return c.StorageType == "memory"
}

func (c *Config) IsPostgresStorage() bool {
	return c.StorageType == "postgres"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
