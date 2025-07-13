package db

import (
	"encoding/json"
)

type Transaction struct {
	MessageID string
	UserID    int64
	Email     string
	Type      string
	TierName  string
	Amount    string
	Currency  string
	Data      interface{}
}

func SaveTransaction(d Transaction) error {
	dataJSON, _ := json.Marshal(d.Data)

	_, err := DB.Exec(`INSERT INTO transactions 
	(message_id, user_id, email, type, tier_name, amount, currency, data)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		d.MessageID, d.UserID, d.Email, d.Type, d.TierName, d.Amount, d.Currency, string(dataJSON))
	return err
}

func EnqueueGRPC(userID int64, amount string) error {
	_, err := DB.Exec(`INSERT INTO grpc_queue (user_id, amount) VALUES (?, ?)`, userID, amount)
	return err
}

func GetPendingGRPC() ([]struct {
	ID     int64
	UserID int64
	Amount string
}, error) {
	rows, err := DB.Query(`SELECT id, user_id, amount FROM grpc_queue WHERE processed_at IS NULL AND attempts < 5`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []struct {
		ID     int64
		UserID int64
		Amount string
	}
	for rows.Next() {
		var item struct {
			ID     int64
			UserID int64
			Amount string
		}
		rows.Scan(&item.ID, &item.UserID, &item.Amount)
		list = append(list, item)
	}
	return list, nil
}

func MarkGRPCProcessed(id int64) {
	DB.Exec(`UPDATE grpc_queue SET processed_at = NOW(), attempts = attempts + 1 WHERE id = ?`, id)
}

func MarkGRPCError(id int64, err string) {
	DB.Exec(`UPDATE grpc_queue SET attempts = attempts + 1, last_error = ? WHERE id = ?`, err, id)
}
