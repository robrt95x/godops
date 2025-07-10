package repository

import "github.com/robrt95x/godops/services/order/internal/entity"

type OrderRepository interface {
	Save(order *entity.Order) error
	FindByID(id string) (*entity.Order, error)
}
