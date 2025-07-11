package memory

import (
	"database/sql"
	"sync"

	"github.com/robrt95x/godops/services/order/internal/entity"
)

type OrderMemoryRepository struct {
	orders map[string]*entity.Order
	mutex  sync.RWMutex
}

func NewOrderMemoryRepository() *OrderMemoryRepository {
	return &OrderMemoryRepository{
		orders: make(map[string]*entity.Order),
		mutex:  sync.RWMutex{},
	}
}

func (r *OrderMemoryRepository) Save(order *entity.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Create a copy of the order to avoid reference issues
	orderCopy := *order
	itemsCopy := make([]entity.OrderItem, len(order.Items))
	copy(itemsCopy, order.Items)
	orderCopy.Items = itemsCopy
	
	r.orders[order.ID] = &orderCopy
	return nil
}

func (r *OrderMemoryRepository) FindByID(id string) (*entity.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	order, exists := r.orders[id]
	if !exists {
		return nil, sql.ErrNoRows
	}
	
	// Return a copy to avoid external modifications
	orderCopy := *order
	itemsCopy := make([]entity.OrderItem, len(order.Items))
	copy(itemsCopy, order.Items)
	orderCopy.Items = itemsCopy
	
	return &orderCopy, nil
}

// Additional helper methods for testing
func (r *OrderMemoryRepository) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.orders = make(map[string]*entity.Order)
}

func (r *OrderMemoryRepository) Count() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.orders)
}

func (r *OrderMemoryRepository) GetAll() []*entity.Order {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	orders := make([]*entity.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orderCopy := *order
		itemsCopy := make([]entity.OrderItem, len(order.Items))
		copy(itemsCopy, order.Items)
		orderCopy.Items = itemsCopy
		orders = append(orders, &orderCopy)
	}
	return orders
}
