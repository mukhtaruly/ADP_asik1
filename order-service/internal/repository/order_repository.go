package repository

import "order-service/internal/domain"

type OrderRepository interface {
	Save(order domain.Order) error
	GetByID(id string) (domain.Order, error)
	GetStatus(orderID string) (string, error)
	Update(order domain.Order) error
}
