package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "develop:password@tcp(db:3306)/mydb"
  db, err := sql.Open("mysql", dsn)
  if err != nil {
    log.Fatal("DB接続エラー：", err)
  }
  defer db.Close()

  if err := db.Ping(); err != nil {
    log.Fatal("DBに接続できません：", err)
  }

  postalCode := "1234567"
  _, err = db.Exec("INSERT INTO access_logs (postal_code) VALUES (?)", postalCode)
  if err != nil {
    log.Fatal("INSERTエラー：", err)
  }
  fmt.Println("データを挿入しました!")
}