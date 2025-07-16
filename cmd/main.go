package main

import (
	_ "fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"kofi-gateway/grpc"
	"kofi-gateway/internal"
	"kofi-gateway/internal/db"
	"kofi-gateway/internal/handlers"
	"kofi-gateway/internal/logger"
)

func main() {
	// Загрузка переменных окружения из .env (если есть)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env не найден — продолжаем с системными переменными")
	}

	// Инициализация БД
	if err := db.Init(); err != nil {
		log.Fatalf("❌ DB error: %v", err)
	}

	// Инициализация логгера
	if err := logger.InitLogger(); err != nil {
		log.Printf("⚠️ Logger error: %v", err)
	}

	// Старт retry-потока
	internal.StartGRPCRetryLoop()

	// gRPC-сервер (запускается в фоне)
	go grpc.StartGRPCServer()

	// HTTP-эндпоинты
	http.HandleFunc("/webhook/kofi", handlers.KofiWebhookHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Старт HTTP-сервера
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("🌐 HTTP server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
