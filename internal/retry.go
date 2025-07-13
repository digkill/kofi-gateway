package internal

import (
	"github.com/digkill/kofi-gateway/grpc"
	"github.com/digkill/kofi-gateway/internal/db"
	"log"
	"time"
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
				err := grpc.MarkPaymentCompleted(task.UserID, task.Amount)
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
