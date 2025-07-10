package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/robrt95x/godops/services/order/internal/entity"
	"github.com/robrt95x/godops/services/order/internal/repository"
)

type CreateOrderCase struct {
	repository repository.OrderRepository
}

func NewCreateOrderCase(repository repository.OrderRepository) *CreateOrderCase {
	return &CreateOrderCase{repository: repository}
}

func (uc *CreateOrderCase) Execute(userID string, items []entity.OrderItem) (*entity.Order, error) {
	order := &entity.Order{
		ID:        uuid.NewString(),
		UserID:    userID,
		Items:     items,
		Status:    entity.Pending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := uc.repository.Save(order)
	if err != nil {
		return nil, err
	}
	return order, nil
}
