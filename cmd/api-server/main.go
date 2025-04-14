package main

import (
	"fmt"
	"log"
	"net/http"

	"GO-api/internal/db"
	"GO-api/internal/handler"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer db.DB.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/address/access_logs", handler.AccessLogsHandler)
	mux.HandleFunc("/address", handler.AddressHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		handler.HelloHandler(w, r)
	})

	fmt.Println("Server is running on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
