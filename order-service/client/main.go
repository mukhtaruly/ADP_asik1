package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "order-service/pkg/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := os.Getenv("ORDER_SERVICE_GRPC_ADDR")
	if addr == "" {
		addr = "localhost:50052"
	}

	orderID := os.Getenv("ORDER_ID")
	if orderID == "" {
		orderID = "your-id"
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)
	stream, err := client.SubscribeToOrderUpdates(context.Background(), &pb.OrderRequest{OrderId: orderID})
	if err != nil {
		log.Fatal(err)
	}

	for {
		update, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Status changed: %s at %v\n", update.GetStatus(), update.GetUpdatedAt().AsTime())
	}
}
