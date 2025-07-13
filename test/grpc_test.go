package grpc_test

import (
	"context"
	"github.com/digkill/kofi-gateway/grpc"
	"testing"

	pb "github.com/digkill/kofi-gateway/proto"
)

func TestConfirmPayment(t *testing.T) {
	server := &grpc.PaymentServer{}
	resp, err := server.ConfirmPayment(context.Background(), &pb.PaymentRequest{
		UserId: 111,
		Amount: "15.50",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !resp.Success {
		t.Errorf("expected success=true, got false")
	}
}
