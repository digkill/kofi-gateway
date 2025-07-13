package main

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
)

func main() {
	data := url.Values{}
	data.Set("data", `{
		"verification_token": "your-secret",
		"message_id": "test-message-001",
		"timestamp": "2025-07-13T21:00:00Z",
		"type": "Donation",
		"is_public": true,
		"from_name": "Test User",
		"message": "For test",
		"amount": "10.00",
		"url": "https://ko-fi.com/test",
		"email": "jo.example@example.com",
		"currency": "USD",
		"is_subscription_payment": false,
		"is_first_subscription_payment": false,
		"kofi_transaction_id": "txn-123",
		"tier_name": null,
		"shop_items": null,
		"shipping": null
	}`)

	resp, err := http.Post(
		"http://localhost:8080/webhook/kofi",
		"application/x-www-form-urlencoded",
		bytes.NewBufferString(data.Encode()),
	)
	if err != nil {
		log.Fatal("❌ Error:", err)
	}
	defer resp.Body.Close()

	log.Println("✅ Webhook test sent. Status:", resp.Status)
}
