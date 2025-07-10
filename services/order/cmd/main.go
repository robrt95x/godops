package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	httpDelivery "github.com/robrt95x/godops/services/order/internal/delivery/http"
	"github.com/robrt95x/godops/services/order/internal/infra/postgres"
	"github.com/robrt95x/godops/services/order/internal/usecase"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/godops?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewOrderPostgresRepository(db)
	createUC := usecase.NewCreateOrderCase(repo)
	handler := httpDelivery.NewOrderHandler(createUC)

	http.HandleFunc("/orders", handler.CreateOrder)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
