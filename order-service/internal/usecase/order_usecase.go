package usecase

import (
	"context"
	"errors"
	"time"

	"order-service/internal/domain"
	"order-service/internal/repository"

	"github.com/google/uuid"
)

type OrderUsecase struct {
	paymentClient *PaymentClient
	repo          repository.OrderRepository
}

func NewOrderUsecase(pc *PaymentClient, repo repository.OrderRepository) *OrderUsecase {
	return &OrderUsecase{
		paymentClient: pc,
		repo:          repo,
	}
}

func (u *OrderUsecase) CreateOrder(customerID, item string, amount int64) domain.Order {
	order := domain.Order{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		ItemName:   item,
		Amount:     amount,
		Status:     "Pending",
		CreatedAt:  time.Now(),
	}

	if u.paymentClient == nil {
		order.Status = "Failed"
	} else {
		paymentResp, err := u.paymentClient.ProcessPayment(order.ID, order.Amount)
		if err != nil || paymentResp == nil {
			order.Status = "Failed"
		} else {
			order.Status = paymentResp.GetMessage()
		}
	}

	// 🔥 сохраняем в БД
	_ = u.repo.Save(order)

	return order
}

func (u *OrderUsecase) GetOrder(id string) (domain.Order, error) {
	return u.repo.GetByID(id)
}

func (u *OrderUsecase) CancelOrder(id string) error {
	order, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}

	if order.Status != "Pending" {
		return errors.New("cannot cancel non-pending order")
	}

	order.Status = "Cancelled"
	return u.repo.Update(order)
}

func (u *OrderUsecase) SubscribeToOrderUpdates(ctx context.Context, orderID string, onUpdate func(domain.Order) error) error {
	order, err := u.repo.GetByID(orderID)
	if err != nil {
		return err
	}

	lastStatus := order.Status
	if err := onUpdate(order); err != nil {
		return err
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			currentOrder, err := u.repo.GetByID(orderID)
			if err != nil {
				return err
			}

			if currentOrder.Status == lastStatus {
				continue
			}

			lastStatus = currentOrder.Status
			if err := onUpdate(currentOrder); err != nil {
				return err
			}
		}
	}
}
