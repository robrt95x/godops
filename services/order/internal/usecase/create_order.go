package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/robrt95x/godops/services/order/internal/entity"
	"github.com/robrt95x/godops/services/order/internal/errors"
	"github.com/robrt95x/godops/services/order/internal/repository"
	"github.com/sirupsen/logrus"
)

type CreateOrderCase struct {
	repository repository.OrderRepository
	logger     *logrus.Logger
}

func NewCreateOrderCase(repository repository.OrderRepository, logger *logrus.Logger) *CreateOrderCase {
	return &CreateOrderCase{
		repository: repository,
		logger:     logger,
	}
}

func (uc *CreateOrderCase) Execute(userID string, items []entity.OrderItem) (*entity.Order, error) {
	logEntry := uc.logger.WithFields(logrus.Fields{
		"use_case":    "CreateOrder",
		"user_id":     userID,
		"items_count": len(items),
	})
	
	logEntry.Debug("Starting create order use case")
	
	// Validate input
	if userID == "" {
		logEntry.Warning("Create order failed: missing user ID")
		return nil, errors.ErrValidationMissingUserID
	}
	
	if len(items) == 0 {
		logEntry.Warning("Create order failed: no items provided")
		return nil, errors.ErrValidationEmptyItems
	}
	
	// Validate items
	var total float64
	for i, item := range items {
		if item.ProductID == "" {
			logEntry.WithField("item_index", i).Warning("Create order failed: missing product ID")
			return nil, errors.ErrValidationMissingProductID
		}
		if item.Quantity <= 0 {
			logEntry.WithFields(logrus.Fields{
				"item_index": i,
				"quantity":   item.Quantity,
			}).Warning("Create order failed: invalid quantity")
			return nil, errors.ErrValidationInvalidQuantity
		}
		if item.Price <= 0 {
			logEntry.WithFields(logrus.Fields{
				"item_index": i,
				"price":      item.Price,
			}).Warning("Create order failed: invalid price")
			return nil, errors.ErrValidationInvalidPrice
		}
		total += item.Price * float64(item.Quantity)
	}

	orderID := uuid.NewString()
	order := &entity.Order{
		ID:        orderID,
		UserID:    userID,
		Items:     items,
		Status:    entity.Pending,
		Total:     total,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	logEntry = logEntry.WithFields(logrus.Fields{
		"order_id": orderID,
		"total":    total,
	})

	err := uc.repository.Save(order)
	if err != nil {
		logEntry.WithError(err).Error("Failed to save order to repository")
		return nil, errors.ErrDatabaseQuery
	}
	
	logEntry.Info("Order created successfully")
	return order, nil
}
