package postgres

import (
	"database/sql"

	"order-service/internal/domain"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Save(order domain.Order) error {
	_, err := r.db.Exec(
		`INSERT INTO orders (id, customer_id, item_name, amount, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		order.ID,
		order.CustomerID,
		order.ItemName,
		order.Amount,
		order.Status,
		order.CreatedAt,
	)
	return err
}

func (r *OrderRepo) GetByID(id string) (domain.Order, error) {
	row := r.db.QueryRow(
		`SELECT id, customer_id, item_name, amount, status, created_at 
		 FROM orders WHERE id=$1`, id,
	)

	var order domain.Order
	err := row.Scan(
		&order.ID,
		&order.CustomerID,
		&order.ItemName,
		&order.Amount,
		&order.Status,
		&order.CreatedAt,
	)

	return order, err
}

func (r *OrderRepo) Update(order domain.Order) error {
	_, err := r.db.Exec(
		`UPDATE orders SET status=$1 WHERE id=$2`,
		order.Status,
		order.ID,
	)
	return err
}
func (r *OrderRepo) GetStatus(orderID string) (string, error) {
	var status string

	err := r.db.QueryRow(
		"SELECT status FROM orders WHERE id = $1",
		orderID,
	).Scan(&status)

	if err != nil {
		return "", err
	}

	return status, nil
}
