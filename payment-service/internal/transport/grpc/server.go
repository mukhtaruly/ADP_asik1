package grpc

import (
	"context"

	pb "payment-service/pkg/payment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *PaymentServer) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	if req.GetOrderId() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	if req.GetAmount() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid amount")
	}

	message := "Pending"
	success := false

	switch {
	case req.GetAmount() < 500:
		message = "Paid"
		success = true
	case req.GetAmount() <= 5000:
		message = "Pending"
	default:
		message = "Failed"
	}

	return &pb.PaymentResponse{
		Success: success,
		Message: message,
	}, nil
}
