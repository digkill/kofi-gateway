// migrate.go
package main

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env не найден — продолжаем с системными переменными")
	}

	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("❌ MYSQL_DSN не найден в .env")
	}

	db, err := goose.OpenDBWithDriver("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Goose DB open error:", err)
	}
	defer db.Close()

	// Запуск миграций
	if err := goose.Up(db, "./internal/db/migrations"); err != nil {
		log.Fatal("❌ Goose up error:", err)
	}

	log.Println("✅ Goose миграции применены")
}
