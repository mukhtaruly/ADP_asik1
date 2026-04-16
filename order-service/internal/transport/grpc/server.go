package grpc

import (
	"database/sql"
	"time"

	"order-service/internal/repository"
	pb "order-service/pkg/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	repo repository.OrderRepository
}

func NewOrderServer(repo repository.OrderRepository) *OrderServer {
	return &OrderServer{repo: repo}
}

func (s *OrderServer) SubscribeToOrderUpdates(req *pb.OrderRequest, stream pb.OrderService_SubscribeToOrderUpdatesServer) error {
	if req.GetOrderId() == "" {
		return status.Error(codes.InvalidArgument, "order_id is required")
	}

	orderID := req.GetOrderId()
	order, err := s.repo.GetByID(orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return status.Error(codes.NotFound, "order not found")
		}
		return status.Error(codes.Internal, err.Error())
	}

	sendUpdate := func(statusValue string) error {
		return stream.Send(&pb.OrderStatusUpdate{
			OrderId:   orderID,
			Status:    statusValue,
			UpdatedAt: timestamppb.Now(),
		})
	}

	finalStatus := "Pending"
	switch {
	case order.Amount < 500:
		finalStatus = "Paid"
	case order.Amount <= 5000:
		finalStatus = "Pending"
	default:
		finalStatus = "Failed"
	}

	statuses := []string{"Pending", "Processing", finalStatus}

	var lastStatus string

	for i, statusValue := range statuses {
		if i > 0 {
			time.Sleep(1 * time.Second)
		}

		select {
		case <-stream.Context().Done():
			return nil
		default:
		}

		if statusValue == lastStatus {
			continue
		}

		order.Status = statusValue
		if err := s.repo.Update(order); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		currentStatus, err := s.repo.GetStatus(orderID)
		if err != nil {
			if err == sql.ErrNoRows {
				return status.Error(codes.NotFound, "order not found")
			}
			return status.Error(codes.Internal, err.Error())
		}

		if currentStatus != lastStatus {
			if err := sendUpdate(currentStatus); err != nil {
				return nil
			}
			lastStatus = currentStatus
		}
	}

	return nil
}
