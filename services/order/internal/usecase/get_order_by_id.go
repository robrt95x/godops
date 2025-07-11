package usecase

import (
	"database/sql"

	"github.com/robrt95x/godops/services/order/internal/entity"
	"github.com/robrt95x/godops/services/order/internal/errors"
	"github.com/robrt95x/godops/services/order/internal/repository"
	"github.com/sirupsen/logrus"
)

type GetOrderByIDCase struct {
	repository repository.OrderRepository
	logger     *logrus.Logger
}

func NewGetOrderByIDCase(repository repository.OrderRepository, logger *logrus.Logger) *GetOrderByIDCase {
	return &GetOrderByIDCase{
		repository: repository,
		logger:     logger,
	}
}

func (uc *GetOrderByIDCase) Execute(id string) (*entity.Order, error) {
	logEntry := uc.logger.WithFields(logrus.Fields{
		"use_case": "GetOrderByID",
		"order_id": id,
	})
	
	logEntry.Debug("Starting get order by ID use case")
	
	if id == "" {
		logEntry.Warning("Invalid order ID: empty string provided")
		return nil, errors.ErrOrderInvalidID
	}

	order, err := uc.repository.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			logEntry.Info("Order not found")
			return nil, errors.ErrOrderNotFound
		}
		logEntry.WithError(err).Error("Failed to retrieve order from repository")
		return nil, errors.ErrDatabaseQuery
	}

	logEntry.Info("Order retrieved successfully")
	return order, nil
}
