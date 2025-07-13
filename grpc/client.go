// grpc/client.go
package grpc

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/digkill/kofi-gateway/proto"
	"google.golang.org/grpc"
)

func MarkPaymentCompleted(userID int64, amount string) error {
	addr := os.Getenv("CORE_SERVICE_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ConfirmPayment(ctx, &pb.PaymentRequest{
		UserId: userID,
		Amount: amount,
	})
	if err != nil {
		return err
	}

	log.Printf("gRPC response: success=%v, message=%s", resp.Success, resp.Message)
	return nil
}
