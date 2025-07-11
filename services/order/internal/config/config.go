package config

import (
	"log"
	"os"

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
	LogLevel   string `env:"LOG_LEVEL" default:"info"`
	
	// Environment
	AppEnv string `env:"APP_ENV" default:"development"`
}

func Load() *Config {
	// Try to load .env file (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables and defaults")
	}
	
	config := &Config{
		StorageType: getEnv("STORAGE_TYPE", "postgres"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "user"),
		DBPassword:  getEnv("DB_PASSWORD", "pass"),
		DBName:      getEnv("DB_NAME", "godops"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		AppEnv:      getEnv("APP_ENV", "development"),
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
