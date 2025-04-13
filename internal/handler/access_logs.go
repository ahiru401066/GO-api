package handler

import (
	"encoding/json"
	"net/http"

	"GO-api/internal/db"
)

func AccessLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// データベースからアクセスログを取得
	accessLogs, err := db.ReadAccessLogs()
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(accessLogs)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
