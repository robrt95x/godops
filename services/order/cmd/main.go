package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	pkgLogger "github.com/robrt95x/godops/pkg/logger"
	pkgMiddleware "github.com/robrt95x/godops/pkg/middleware"
	"github.com/robrt95x/godops/services/order/internal/config"
	httpDelivery "github.com/robrt95x/godops/services/order/internal/delivery/http"
	"github.com/robrt95x/godops/services/order/internal/infra"
	"github.com/robrt95x/godops/services/order/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
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
		ServiceName: "order-service",
	}
	appLogger := pkgLogger.Setup(loggerConfig)
	
	appLogger.WithField("storage_type", cfg.StorageType).Info("Starting order service")

	// Create repository using factory
	factory := infra.NewRepositoryFactory(cfg)
	repo, err := factory.CreateOrderRepository()
	if err != nil {
		appLogger.WithError(err).Fatal("Failed to create repository")
	}

	// Create use cases
	createUC := usecase.NewCreateOrderCase(repo, appLogger)
	getOrderByIDUC := usecase.NewGetOrderByIDCase(repo, appLogger)
	handler := httpDelivery.NewOrderHandler(createUC, getOrderByIDUC, appLogger)

	// Setup router with middleware
	r := chi.NewRouter()
	
	// Add custom middleware
	r.Use(pkgMiddleware.RequestID)
	r.Use(pkgMiddleware.Logging(appLogger))
	r.Use(pkgMiddleware.ErrorLogging(appLogger))
	r.Use(middleware.Recoverer)

	r.Route("/orders", func(r chi.Router) {
		r.Post("/", handler.CreateOrder)
		r.Get("/{id}", handler.GetOrderByID)
	})

	appLogger.WithField("port", cfg.ServerPort).Info("Starting HTTP server")
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		appLogger.WithError(err).Fatal("HTTP server failed")
	}
}
