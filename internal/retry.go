package internal

import (
	"log"
	"time"

	"kofi-gateway/grpc"
	"kofi-gateway/internal/db"
)

func StartGRPCRetryLoop() {
	go func() {
		for {
			time.Sleep(10 * time.Second)

			tasks, err := db.GetPendingGRPC()
			if err != nil {
				log.Println("Retry query error:", err)
				continue
			}

			for _, task := range tasks {
				err := grpc.MarkPaymentCompleted(
					task.UserID,
					task.OrderID,
					task.Amount,
					task.Credits,
					task.Email,
					task.Username,
					task.Provider,
				)
				if err != nil {
					log.Println("❌ Retry gRPC failed:", err)
					db.MarkGRPCError(task.ID, err.Error())
				} else {
					log.Println("✅ Retry gRPC OK")
					db.MarkGRPCProcessed(task.ID)
				}
			}
		}
	}()
}
