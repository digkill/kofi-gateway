package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/digkill/kofi-gateway/proto"
	"google.golang.org/grpc"
)

var grpcConn pb.PaymentServiceClient

// Инициализация gRPC-клиента
func InitGRPCClient(address string) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return fmt.Errorf("не удалось подключиться к gRPC: %w", err)
	}
	grpcConn = pb.NewPaymentServiceClient(conn)
	log.Println("✅ gRPC клиент инициализирован")
	return nil
}

// Отправка запроса на оплату
func RequestPayment(
	userID int64,
	orderID string,
	amount int32,
	credits int32,
	email string,
	username string,
	provider string,
) (*pb.PaymentResponse, error) {
	if grpcConn == nil {
		return nil, fmt.Errorf("gRPC клиент не инициализирован")
	}

	req := &pb.PaymentRequest{
		UserId:   userID,
		OrderId:  orderID,
		Amount:   amount,
		Credits:  credits,
		Email:    email,
		Username: username,
		Provider: provider,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := grpcConn.RequestPayment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе оплаты: %w", err)
	}

	if resp.Status != "ok" {
		log.Printf("⚠️ Ошибка оплаты: %s", resp.Message)
	} else {
		log.Printf("✅ Оплата успешно инициирована: %s", resp.PaymentLink)
	}

	return resp, nil
}

// Подтверждение оплаты (например, из retry-очереди)
func MarkPaymentCompleted(
	userID int64,
	orderID string,
	amount int32,
	credits int32,
	email string,
	username string,
	provider string,
) error {
	if grpcConn == nil {
		return fmt.Errorf("gRPC клиент не инициализирован")
	}

	req := &pb.PaymentRequest{
		UserId:   userID,
		OrderId:  orderID,
		Amount:   amount,
		Credits:  credits,
		Email:    email,
		Username: username,
		Provider: provider,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := grpcConn.ConfirmPayment(ctx, req)
	if err != nil {
		return fmt.Errorf("ошибка при вызове ConfirmPayment: %w", err)
	}

	if resp.Status != "ok" {
		return fmt.Errorf("оплата не подтверждена: %s", resp.Message)
	}

	log.Printf("✅ Оплата подтверждена для пользователя %d: %s", userID, resp.Message)
	return nil
}
