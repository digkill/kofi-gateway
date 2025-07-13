package types

type KofiWebhookData struct {
	VerificationToken          string      `json:"verification_token"`
	MessageID                  string      `json:"message_id"`
	Timestamp                  string      `json:"timestamp"`
	Type                       string      `json:"type"`
	IsPublic                   bool        `json:"is_public"`
	FromName                   string      `json:"from_name"`
	Message                    string      `json:"message"`
	Amount                     string      `json:"amount"`
	URL                        string      `json:"url"`
	Email                      string      `json:"email"`
	Currency                   string      `json:"currency"`
	IsSubscriptionPayment      bool        `json:"is_subscription_payment"`
	IsFirstSubscriptionPayment bool        `json:"is_first_subscription_payment"`
	KofiTransactionID          string      `json:"kofi_transaction_id"`
	TierName                   string      `json:"tier_name"`
	ShopItems                  interface{} `json:"shop_items"`
	Shipping                   interface{} `json:"shipping"`
}
