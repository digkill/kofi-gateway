package db

import (
	"database/sql"
	_ "fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() error {
	dsn := os.Getenv("MYSQL_DSN") // example: user:pass@tcp(localhost:3306)/kofi
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return DB.Ping()
}
