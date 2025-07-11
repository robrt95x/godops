package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/robrt95x/godops/services/order/internal/config"
	httpDelivery "github.com/robrt95x/godops/services/order/internal/delivery/http"
	"github.com/robrt95x/godops/services/order/internal/infra"
	"github.com/robrt95x/godops/services/order/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Starting order service with %s storage", cfg.StorageType)

	// Create repository using factory
	factory := infra.NewRepositoryFactory(cfg)
	repo, err := factory.CreateOrderRepository()
	if err != nil {
		log.Fatal("Failed to create repository:", err)
	}

	// Create use cases
	createUC := usecase.NewCreateOrderCase(repo)
	getOrderByIDUC := usecase.NewGetOrderByIDCase(repo)
	handler := httpDelivery.NewOrderHandler(createUC, getOrderByIDUC)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/orders", func(r chi.Router) {
		r.Post("/", handler.CreateOrder)
		r.Get("/{id}", handler.GetOrderByID)
	})

	log.Printf("Starting server on :%s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
