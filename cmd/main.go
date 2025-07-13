package main

import (
	"fmt"
	"github.com/digkill/kofi-gateway/grpc"
	"github.com/digkill/kofi-gateway/internal"
	"github.com/digkill/kofi-gateway/internal/db"
	"github.com/digkill/kofi-gateway/internal/handlers"
	"github.com/digkill/kofi-gateway/internal/logger"
	"github.com/joho/godotenv"
	"log"

	"net/http"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env не найден — продолжаем с системными переменными")
	}

	if err := db.Init(); err != nil {
		log.Fatal("❌ DB error:", err)
	}
	err := logger.InitLogger()
	if err != nil {
		fmt.Println("Error init logger:", err)
	}
	internal.StartGRPCRetryLoop()

	http.HandleFunc("/webhook/kofi", handlers.KofiWebhookHandler)
	go grpc.StartGRPCServer()

	log.Println("🌐 HTTP server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

}
