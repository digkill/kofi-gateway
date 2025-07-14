package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/digkill/kofi-gateway/grpc"
	"github.com/digkill/kofi-gateway/internal"
	"github.com/digkill/kofi-gateway/internal/logger"
	"github.com/digkill/kofi-gateway/internal/types"
	"github.com/digkill/kofi-gateway/internal/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func KofiWebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusInternalServerError)
		return
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	rawData := values.Get("data")
	if rawData == "" {
		http.Error(w, "no data field", http.StatusBadRequest)
		return
	}

	var payload types.KofiWebhookData
	if err := json.Unmarshal([]byte(rawData), &payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	expected := os.Getenv("KOFI_SECRET")
	if expected != "" && payload.VerificationToken != expected {
		http.Error(w, "invalid token", http.StatusForbidden)
		return
	}

	// 🧾 Лог в файл
	logger.LogTransaction(payload)

	// 🧠 Поиск пользователя
	userID := internal.LookupUserByEmail(payload.Email)
	if userID == 0 {
		log.Printf("⚠️ Не найден пользователь для email %s", payload.Email)
	}

	// 📌 Поддержка подписки: VIP и др.
	if payload.TierName == "VIP" {
		log.Printf("🌟 Новый VIP подписчик: %s", payload.Email)
	}

	// 🎯 Конвертация суммы
	amountInt, _ := strconv.Atoi(payload.Amount) // best-effort
	credits := utils.EstimateCredits(amountInt)  // например, 1 руб = 1.5 кредита

	// ✅ Отправка в CoreService
	err = grpc.MarkPaymentCompleted(
		userID,
		payload.TransactionID, // → OrderID
		int32(amountInt),
		int32(credits),
		payload.Email,
		payload.FromName,
		"kofi",
	)
	if err != nil {
		log.Printf("❌ gRPC ошибка: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}
