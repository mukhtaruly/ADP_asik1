# Assignment 2: gRPC Microservices

## Repositories

- Proto repository: `https://github.com/mukhtaruly/mukhtaruly-proto`
- Generated repository: `https://github.com/mukhtaruly/mukhtaruly-generated`

Note: after publishing `mukhtaruly-generated` as a real Go module tag like `v1.0.0`, replace this local dependency in [order-service/go.mod](/Users/nurassylmuktaruly/Downloads/as2/order-service/go.mod:1):

```bash
cd /Users/nurassylmuktaruly/Downloads/as2/order-service
go get github.com/mukhtaruly/mukhtaruly-generated@v1.0.0
```

Then remove:

```go
replace github.com/mukhtaruly/as2/payment-service => ../payment-service
```

## Run

```bash
cd /Users/nurassylmuktaruly/Downloads/as2/payment-service && go run ./cmd/main.go
cd /Users/nurassylmuktaruly/Downloads/as2/order-service && go run ./cmd/main.go
cd /Users/nurassylmuktaruly/Downloads/as2/order-service && ORDER_ID=your-order-id go run ./client/main.go
```

## Environment Variables

| Variable | Service | Default | Description |
|---|---|---|---|
| `DATABASE_URL` | `order-service` | `postgres://postgres:1234@localhost:5433/orders_db?sslmode=disable` | PostgreSQL connection string |
| `ORDER_SERVICE_GRPC_ADDR` | `order-service` | `:50052` | gRPC address for order streaming server |
| `PAYMENT_SERVICE_ADDR` | `order-service` | `localhost:50051` | Address of payment gRPC server used by order-service client |
| `PAYMENT_GRPC_ADDR` | `payment-service` | `:50051` | gRPC address for payment-service |
| `ORDER_ID` | `order-service/client` | `your-order-id` | Order ID for streaming demo client |

## Streaming Demo

1. Start the streaming client for an existing order ID.
2. In another terminal run:

```sql
UPDATE orders SET status='Shipped' WHERE id='your-order-id';
```

3. The client prints the new status immediately from the gRPC stream.

## Streaming Screenshot

Add your screenshot here before submission, for example:

`docs/streaming-client.png`
