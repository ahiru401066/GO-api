package main

import (
  "fmt"
  "net/http"

  "GO-api/internal/handler"
)

func main() {
	http.HandleFunc("/", handler.HelloHandler)
  http.HandleFunc("/address", handler.AddressHandler)
  http.HandleFunc("/address/access_logs", handler.AccessLogsHandler)

  fmt.Println("Server is running on :8080...")
  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }
}