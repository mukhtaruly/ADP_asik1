package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"

	grpcTransport "payment-service/internal/transport/grpc"
	pb "payment-service/pkg/payment"
)

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("Method: %s | Duration: %v | Error: %v", info.FullMethod, time.Since(start), err)
	return resp, err
}

func main() {
	// ---------- gRPC SERVER ----------
	grpcAddr := os.Getenv("PAYMENT_GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = ":50051"
	}

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	pb.RegisterPaymentServiceServer(grpcServer, &grpcTransport.PaymentServer{})

	go func() {
		log.Println("gRPC Payment Service running on", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("failed to serve gRPC:", err)
		}
	}()

	// ---------- HTTP SERVER (можешь оставить или убрать) ----------
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "payment service ok",
		})
	})

	log.Println("HTTP Payment Service running on :8082")
	r.Run(":8082")
}
