package db

import (
	"fmt"
	"log"

	"GO-api/internal/model"
	_ "github.com/go-sql-driver/mysql"
)

func AddLog(postalCode string) {
  _, err := DB.Exec("INSERT INTO access_logs (postal_code) VALUES (?)", postalCode)
  if err != nil {
    log.Fatal("INSERTエラー：", err)
  }
  fmt.Println("データを挿入しました!")
}

func GetAccessLogs() (model.AccessLogsResponse, error) {
  rows, err := DB.Query(`
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