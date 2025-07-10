package postgres

import (
	"database/sql"
	"encoding/json"

	"github.com/robrt95x/godops/services/order/internal/entity"
)

type OrderPostgresRespository struct {
	db *sql.DB
}

func NewOrderPostgresRepository(db *sql.DB) *OrderPostgresRespository {
	return &OrderPostgresRespository{db: db}
}

func (r *OrderPostgresRespository) Save(order *entity.Order) error {
	itemsJson, _ := json.Marshal(order.Items)

	_, err := r.db.Exec(
		`INSERT INTO orders (id, user_id, items, status, coupon_code, total, shipping_address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		order.ID,
		order.UserID,
		itemsJson,
		order.Status,
		order.CouponCode,
		order.Total,
		order.ShippingAddress,
		order.CreatedAt,
		order.UpdatedAt,
	)

	return err
}

func (r *OrderPostgresRespository) FindByID(id string) (*entity.Order, error) {
	var order entity.Order
	var itemsJson []byte

	err := r.db.QueryRow(
		`SELECT id, user_id, items, status, coupon_code, total, shipping_address, created_at, updated_at 
		FROM orders WHERE id = $1`, id).Scan(
		&order.ID,
		&order.UserID,
		&itemsJson,
		&order.Status,
		&order.CouponCode,
		&order.Total,
		&order.ShippingAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(itemsJson, &order.Items); err != nil {
		return nil, err
	}

	return &order, nil
}
