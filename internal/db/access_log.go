package db

import (
	"database/sql"
	"fmt"
	"log"

	"GO-api/internal/model"
	_ "github.com/go-sql-driver/mysql"
)

var dsn = "develop:password@tcp(db:3306)/mydb"

func AddLog(postalCode string) {
  db, err := sql.Open("mysql", dsn)
  if err != nil {
    log.Fatal("DB接続エラー：", err)
  }
  defer db.Close()

  _, err = db.Exec("INSERT INTO access_logs (postal_code) VALUES (?)", postalCode)
  if err != nil {
    log.Fatal("INSERTエラー：", err)
  }
  fmt.Println("データを挿入しました!")
}

func GetAccessLogs() (model.AccessLogsResponse, error) {
	db, err := sql.Open("mysql", dsn)
  if err != nil {
    return model.AccessLogsResponse{}, err
  }
  defer db.Close()

  rows, err := db.Query(`
		SELECT postal_code, COUNT(*) AS request_count
		FROM access_logs
		GROUP BY postal_code
		ORDER BY request_count DESC
  `)

  if err != nil {
		return model.AccessLogsResponse{}, err
	}
	defer rows.Close()

  var accessLogs []model.AccessLog
  for rows.Next() {
    var accessLog model.AccessLog
    if err := rows.Scan(&accessLog.PostalCode, &accessLog.RequestCount);  err != nil {
			return model.AccessLogsResponse{}, err
    }
    accessLogs = append(accessLogs, accessLog)
  }
	return model.AccessLogsResponse{AccessLogs: accessLogs}, nil
}