package infra

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/robrt95x/godops/services/order/internal/config"
	"github.com/robrt95x/godops/services/order/internal/infra/memory"
	"github.com/robrt95x/godops/services/order/internal/infra/postgres"
	"github.com/robrt95x/godops/services/order/internal/repository"
)

type RepositoryFactory struct {
	config *config.Config
}

func NewRepositoryFactory(config *config.Config) *RepositoryFactory {
	return &RepositoryFactory{config: config}
}

func (f *RepositoryFactory) CreateOrderRepository() (repository.OrderRepository, error) {
	switch {
	case f.config.IsMemoryStorage():
		log.Println("Using in-memory storage for orders")
		return memory.NewOrderMemoryRepository(), nil
		
	case f.config.IsPostgresStorage():
		log.Println("Using PostgreSQL storage for orders")
		db, err := f.createPostgresConnection()
		if err != nil {
			return nil, fmt.Errorf("failed to create postgres connection: %w", err)
		}
		return postgres.NewOrderPostgresRepository(db), nil
		
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", f.config.StorageType)
	}
}

func (f *RepositoryFactory) createPostgresConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", f.config.GetDatabaseURL())
	if err != nil {
		return nil, err
	}
	
	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	log.Printf("Connected to PostgreSQL at %s:%s", f.config.DBHost, f.config.DBPort)
	return db, nil
}
