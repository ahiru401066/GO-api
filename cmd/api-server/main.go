package main

import (
  "fmt"
  "log"
  "net/http"

  "GO-api/internal/handler"
  "GO-api/internal/db"
)

func main() {
  if err := db.Init(); err != nil {
    log.Fatalf("DB接続失敗: %v", err)
  }
  defer db.DB.Close()

  mux := http.NewServeMux()
	mux.HandleFunc("/", handler.HelloHandler)
  mux.HandleFunc("/address", handler.AddressHandler)
  mux.HandleFunc("/address/access_logs", handler.AccessLogsHandler)

  fmt.Println("Server is running on :8080...")
  if err := http.ListenAndServe(":8080", mux); err != nil {
    log.Fatal(err)
  }
}