package logger

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logFile *os.File

func InitLogger() error {
	_ = os.MkdirAll("./storage/logs", 0755)
	file, err := os.OpenFile(filepath.Join("./storage/logs", "transactions.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	logFile = file
	return nil
}

func LogTransaction(data any) {
	if logFile == nil {
		log.Println("log file not initialized")
		return
	}
	entry := struct {
		Timestamp string      `json:"timestamp"`
		Data      interface{} `json:"data"`
	}{
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      data,
	}
	json.NewEncoder(logFile).Encode(entry)
}
