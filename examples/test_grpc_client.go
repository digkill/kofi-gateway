package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "kofi-gateway/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	client := pb.NewPaymentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ConfirmPayment(ctx, &pb.PaymentRequest{
		UserId: 12345,
		Amount: "5.00",
	})
	if err != nil {
		log.Fatal("❌ Error:", err)
	}

	log.Printf("✅ gRPC Response: success=%v, message=%s", resp.Success, resp.Message)
}
