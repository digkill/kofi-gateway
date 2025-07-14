package db

import (
	_ "database/sql"
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

// ⚠️ Расширенная очередь gRPC
type GRPCTask struct {
	ID       int64
	UserID   int64
	OrderID  string
	Amount   int32
	Credits  int32
	Email    string
	Username string
	Provider string
}

// Сохраняем задачу на подтверждение оплаты
func EnqueueGRPC(task GRPCTask) error {
	_, err := DB.Exec(`
		INSERT INTO grpc_queue 
		(user_id, order_id, amount, credits, email, username, provider) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		task.UserID, task.OrderID, task.Amount, task.Credits, task.Email, task.Username, task.Provider)
	return err
}

// Получаем все неподтверждённые задачи
func GetPendingGRPC() ([]GRPCTask, error) {
	rows, err := DB.Query(`
		SELECT id, user_id, order_id, amount, credits, email, username, provider 
		FROM grpc_queue 
		WHERE processed_at IS NULL AND attempts < 5`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []GRPCTask
	for rows.Next() {
		var t GRPCTask
		err := rows.Scan(&t.ID, &t.UserID, &t.OrderID, &t.Amount, &t.Credits, &t.Email, &t.Username, &t.Provider)
		if err != nil {
			continue
		}
		list = append(list, t)
	}
	return list, nil
}

// Успешная обработка
func MarkGRPCProcessed(id int64) {
	DB.Exec(`UPDATE grpc_queue SET processed_at = NOW(), attempts = attempts + 1 WHERE id = ?`, id)
}

// Ошибка при повторной попытке
func MarkGRPCError(id int64, err string) {
	DB.Exec(`UPDATE grpc_queue SET attempts = attempts + 1, last_error = ? WHERE id = ?`, err, id)
}

func AddCredits(userID int64, credits int) error {
	_, err := DB.Exec(`UPDATE users SET credits = credits + ? WHERE telegram_id = ?`, credits, userID)
	return err
}
