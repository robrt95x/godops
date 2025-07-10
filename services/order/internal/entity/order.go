package entity

import "time"

type OrderStatus string

const (
	Pending   OrderStatus = "PENDING"
	Completed OrderStatus = "COMPLETED"
	Cancelled OrderStatus = "CANCELLED"
)

func (s OrderStatus) IsPending() bool {
	return s == Pending
}

func (s OrderStatus) IsCompleted() bool {
	return s == Completed
}

func (s OrderStatus) IsCancelled() bool {
	return s == Cancelled
}

type Order struct {
	ID        string
	UserID    string
	Items     []OrderItem
	Status    OrderStatus
	CouponCode string
	Total     float64
	ShippingAddress string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}
