package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/robrt95x/godops/services/order/internal/entity"
	"github.com/robrt95x/godops/services/order/internal/usecase"
)

type OrderHandler struct {
	CreateUC      *usecase.CreateOrderCase
	GetOrderByIDUC *usecase.GetOrderByIDCase
}

func NewOrderHandler(createUC *usecase.CreateOrderCase, getOrderByIDUC *usecase.GetOrderByIDCase) *OrderHandler {
	return &OrderHandler{
		CreateUC:      createUC,
		GetOrderByIDUC: getOrderByIDUC,
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

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	
	order, err := h.GetOrderByIDUC.Execute(orderID)
	if err != nil {
		switch err {
		case usecase.ErrOrderNotFound:
			http.Error(w, "Order not found", http.StatusNotFound)
		case usecase.ErrInvalidOrderID:
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
		default:
			http.Error(w, "Failed to get order: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
