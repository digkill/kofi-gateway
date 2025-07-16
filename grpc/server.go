package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"kofi-gateway/internal/db"
	"kofi-gateway/internal/payment"
	pb "kofi-gateway/proto"
	"log"
	"net"
)

type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
}

// ✅ Метод вызывается из Telegram-бота при создании заказа
func (s *PaymentServer) RequestPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	log.Printf("💳 RequestPayment: user_id=%d, credits=%d, amount=%d, email=%s, provider=%s",
		req.UserId, req.Credits, req.Amount, req.Email, req.Provider)

	link, err := payment.GeneratePaymentLink(req.OrderId, req.Amount, req.Email)
	if err != nil {
		log.Printf("❌ Ошибка при генерации ссылки Ko-fi: %v", err)
		return &pb.PaymentResponse{
			Status:  "error",
			Message: "Ошибка генерации ссылки оплаты",
		}, nil
	}

	// ✅ Сохраняем транзакцию в базу
	tx := db.Transaction{
		MessageID: req.OrderId,
		UserID:    req.UserId,
		Email:     req.Email,
		Type:      "payment",
		TierName:  "default",
		Amount:    fmt.Sprintf("%d", req.Amount),
		Currency:  "RUB",
		Data: map[string]interface{}{
			"provider": req.Provider,
			"credits":  req.Credits,
			"link":     link,
		},
	}

	if err := db.SaveTransaction(tx); err != nil {
		log.Printf("❌ Не удалось сохранить транзакцию: %v", err)
		return &pb.PaymentResponse{
			Status:  "error",
			Message: "Ошибка записи транзакции",
		}, nil
	}

	log.Printf("✅ Ссылка успешно сгенерирована: %s", link)

	return &pb.PaymentResponse{
		Status:      "ok",
		PaymentLink: link,
		Message:     "Ссылка на оплату успешно создана",
	}, nil
}

// ✅ Метод вызывается шлюзом после подтверждения оплаты
func (s *PaymentServer) ConfirmPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	log.Printf("✅ ConfirmPayment: user_id=%d, amount_rub=%d, credits=%d, email=%s, order_id=%s",
		req.UserId, req.Amount, req.Credits, req.Email, req.OrderId)

	// 🔎 Проверка user_id
	if req.UserId == 0 {
		return &pb.PaymentResponse{
			Status:  "error",
			Message: "Не указан user_id",
		}, nil
	}

	// 💰 Начисляем кредиты пользователю
	err := db.AddCredits(req.UserId, int(req.Credits))
	if err != nil {
		log.Printf("❌ Ошибка при начислении кредитов: %v", err)
		return &pb.PaymentResponse{
			Status:  "error",
			Message: "Не удалось начислить кредиты",
		}, nil
	}

	// 🧾 Сохраняем успешное подтверждение в транзакции
	tx := db.Transaction{
		MessageID: req.OrderId,
		UserID:    req.UserId,
		Email:     req.Email,
		Type:      "confirm",
		TierName:  "confirmed",
		Amount:    fmt.Sprintf("%d", req.Amount),
		Currency:  "RUB",
		Data: map[string]interface{}{
			"provider": req.Provider,
			"credits":  req.Credits,
			"status":   "confirmed",
		},
	}

	if err := db.SaveTransaction(tx); err != nil {
		log.Printf("❌ Ошибка при сохранении подтверждения: %v", err)
	}

	log.Printf("🎉 Кредиты начислены пользователю %d (+%d)", req.UserId, req.Credits)

	return &pb.PaymentResponse{
		Status:  "ok",
		Message: "Платёж успешно подтверждён и кредиты начислены",
	}, nil
}

// 🚀 Запуск gRPC-сервера
func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ Не удалось запустить gRPC-сервер: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, &PaymentServer{})

	log.Println("🚀 gRPC сервер слушает на :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("❌ Ошибка при запуске gRPC-сервера: %v", err)
	}
}
