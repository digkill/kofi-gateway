// grpc/server.go
package grpc

import (
	"context"
	"log"
	"net"

	pb "github.com/digkill/kofi-gateway/proto"
	"google.golang.org/grpc"
)

type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *PaymentServer) ConfirmPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	log.Printf("✅ ConfirmPayment called: user_id=%d, amount=%s", req.UserId, req.Amount)

	// Тут ты можешь активировать товар, доступ и т.д.

	return &pb.PaymentResponse{
		Success: true,
		Message: "Payment confirmed",
	}, nil
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, &PaymentServer{})

	log.Println("🚀 gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
