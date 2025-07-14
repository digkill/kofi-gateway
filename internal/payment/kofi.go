package payment

import (
	"fmt"
	"os"
)

type KofiPaymentRequest struct {
	Title       string `json:"title"`
	Amount      int32  `json:"amount"`
	Currency    string `json:"currency"`
	RedirectURL string `json:"redirect_url"`
	Email       string `json:"email,omitempty"`
	ExternalID  string `json:"external_id,omitempty"`
}

type KofiPaymentResponse struct {
	Success     bool   `json:"success"`
	PaymentLink string `json:"payment_link"`
	Message     string `json:"message,omitempty"`
}

// GeneratePaymentLink формирует ссылку на оплату через Ko-fi
func GeneratePaymentLink(orderID string, amount int32, email string) (string, error) {
	username := os.Getenv("KOFI_USERNAME") // например: digkill
	if username == "" {
		return "", fmt.Errorf("KOFI_USERNAME не установлен")
	}

	baseURL := fmt.Sprintf("https://ko-fi.com/%s/%d", username, amount)

	//	params := url.Values{}
	//	params.Set("amount", fmt.Sprintf("%d", amount))
	//	if email != "" {
	//		params.Set("email", email)
	//	}
	//	params.Set("message", "Оплата заказа "+orderID)

	//	fullURL := baseURL + "?" + params.Encode()
	return baseURL, nil
}
