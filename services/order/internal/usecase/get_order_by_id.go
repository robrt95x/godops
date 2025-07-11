package usecase

import (
	"database/sql"
	"errors"

	"github.com/robrt95x/godops/services/order/internal/entity"
	"github.com/robrt95x/godops/services/order/internal/repository"
)

var (
	ErrOrderNotFound  = errors.New("order not found")
	ErrInvalidOrderID = errors.New("invalid order ID")
)

type GetOrderByIDCase struct {
	repository repository.OrderRepository
}

func NewGetOrderByIDCase(repository repository.OrderRepository) *GetOrderByIDCase {
	return &GetOrderByIDCase{repository: repository}
}

func (uc *GetOrderByIDCase) Execute(id string) (*entity.Order, error) {
	if id == "" {
		return nil, ErrInvalidOrderID
	}

	order, err := uc.repository.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return order, nil
}
