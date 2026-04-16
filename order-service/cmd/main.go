package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"order-service/internal/repository/postgres"
	grpcTransport "order-service/internal/transport/grpc"
	httpHandler "order-service/internal/transport/http"
	"order-service/internal/usecase"
	orderpb "order-service/pkg/order"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:1234@localhost:5433/orders_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewOrderRepo(db)
	paymentClient := usecase.NewPaymentClient()
	orderUsecase := usecase.NewOrderUsecase(paymentClient, repo)

	grpcAddr := os.Getenv("ORDER_SERVICE_GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = ":50052"
	}

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	orderGrpcServer := grpcTransport.NewOrderServer(repo)
	orderpb.RegisterOrderServiceServer(grpcServer, orderGrpcServer)

	go func() {
		log.Println("Order gRPC Service running on", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	r := gin.Default()
	handler := httpHandler.NewHandler(orderUsecase)

	r.POST("/orders", handler.CreateOrder)
	r.GET("/orders/:id", handler.GetOrder)
	r.PATCH("/orders/:id/cancel", handler.CancelOrder)

	log.Println("Order Service running on :8081")
	r.Run(":8081")
}
