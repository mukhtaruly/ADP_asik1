package usecase

import (
	"context"
	"errors"
	"os"
	"time"

	pb "payment-service/pkg/payment"

	"google.golang.org/grpc"
)

type PaymentClient struct {
	client pb.PaymentServiceClient
	conn   *grpc.ClientConn
}

func NewPaymentClient() *PaymentClient {
	addr := os.Getenv("PAYMENT_SERVICE_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}

	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil
	}

	return &PaymentClient{
		client: pb.NewPaymentServiceClient(conn),
		conn:   conn,
	}
}

func (p *PaymentClient) ProcessPayment(orderID string, amount int64) (*pb.PaymentResponse, error) {
	if p == nil {
		return nil, errors.New("payment client is not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return p.client.ProcessPayment(ctx, &pb.PaymentRequest{
		OrderId: orderID,
		Amount:  float64(amount),
	})
}
