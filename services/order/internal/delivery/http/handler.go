package http

import (
	"encoding/json"
	"net/http"

	"github.com/robrt95x/godops/services/order/internal/entity"
	"github.com/robrt95x/godops/services/order/internal/usecase"
)

type OrderHandler struct {
	CreateUC *usecase.CreateOrderCase
}

func NewOrderHandler(createUC *usecase.CreateOrderCase) *OrderHandler {
	return &OrderHandler{
		CreateUC: createUC,
	}
}

type CreateOrderRequest struct {
	UserID string             `json:"user_id"`
	Items  []entity.OrderItem `json:"items"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	order, err := h.CreateUC.Execute(req.UserID, req.Items)
	if err != nil {
		http.Error(w, "Failed to create order: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
