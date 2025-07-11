package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	pkgErrors "github.com/robrt95x/godops/pkg/errors"
	"github.com/robrt95x/godops/services/order/internal/entity"
	"github.com/robrt95x/godops/services/order/internal/errors"
	"github.com/robrt95x/godops/services/order/internal/usecase"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	CreateUC       *usecase.CreateOrderCase
	GetOrderByIDUC *usecase.GetOrderByIDCase
	ErrorHandler   *pkgErrors.HTTPErrorHandler
	Logger         *logrus.Logger
}

func NewOrderHandler(createUC *usecase.CreateOrderCase, getOrderByIDUC *usecase.GetOrderByIDCase, logger *logrus.Logger) *OrderHandler {
	errorCatalog := errors.NewOrderErrorCatalog()
	return &OrderHandler{
		CreateUC:       createUC,
		GetOrderByIDUC: getOrderByIDUC,
		ErrorHandler:   pkgErrors.NewHTTPErrorHandler(logger, errorCatalog),
		Logger:         logger,
	}
}

type CreateOrderRequest struct {
	UserID string             `json:"user_id"`
	Items  []entity.OrderItem `json:"items"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-Request-ID")
	logEntry := h.Logger.WithFields(logrus.Fields{
		"handler":    "CreateOrder",
		"request_id": requestID,
	})
	
	logEntry.Debug("Processing create order request")
	
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logEntry.WithError(err).Warning("Failed to decode request body")
		h.ErrorHandler.HandleValidationError(w, r, "Invalid request body format")
		return
	}
	
	logEntry = logEntry.WithFields(logrus.Fields{
		"user_id":     req.UserID,
		"items_count": len(req.Items),
	})
	
	order, err := h.CreateUC.Execute(req.UserID, req.Items)
	if err != nil {
		logEntry.WithError(err).Error("Create order use case failed")
		h.ErrorHandler.HandleError(w, r, err)
		return
	}
	
	logEntry.WithField("order_id", order.ID).Info("Order created successfully")
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	requestID := r.Header.Get("X-Request-ID")
	
	logEntry := h.Logger.WithFields(logrus.Fields{
		"handler":    "GetOrderByID",
		"request_id": requestID,
		"order_id":   orderID,
	})
	
	logEntry.Debug("Processing get order by ID request")
	
	order, err := h.GetOrderByIDUC.Execute(orderID)
	if err != nil {
		logEntry.WithError(err).Warning("Get order by ID use case failed")
		h.ErrorHandler.HandleError(w, r, err)
		return
	}
	
	logEntry.Info("Order retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
