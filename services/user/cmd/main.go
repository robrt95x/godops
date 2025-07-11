package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robrt95x/godops/pkg/logger"
	"github.com/robrt95x/godops/pkg/middleware"
	"github.com/robrt95x/godops/services/user/internal/adapter/repository"
	userHttp "github.com/robrt95x/godops/services/user/internal/adapter/http"
	"github.com/robrt95x/godops/services/user/internal/application/usecase"
	"github.com/robrt95x/godops/services/user/internal/config"
	"github.com/robrt95x/godops/services/user/internal/domain/service"
)

func main() {
	// Initialize logger
	log := logger.New()
	
	// Initialize config
	cfg := config.New()
	
	// Initialize repository
	userRepo := repository.NewMemoryUserRepository()
	
	// Initialize domain service
	userService := service.NewUserService(userRepo)
	
	// Initialize use cases
	createUserUseCase := usecase.NewCreateUserUseCase(userService)
	getUserUseCase := usecase.NewGetUserUseCase(userService)
	
	// Initialize HTTP handler
	userHandler := userHttp.NewUserHandler(createUserUseCase, getUserUseCase)
	
	// Setup routes
	r := mux.NewRouter()
	
	// Apply middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging(log))
	
	// User routes
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	
	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	log.Infof("User service starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Errorf("Failed to start server: %v", err)
	}
}
