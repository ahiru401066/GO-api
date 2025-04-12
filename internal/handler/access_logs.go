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

	accessLogs, err := db.GetAccessLogs()
	if err != nil {
		http.Error(w, "Failed to get access logs", http.StatusInternalServerError)
		return		
	}

  jsonResponse, err := json.Marshal(accessLogs)
  if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
    return    
  }
	w.Header().Set("Content-Type", "application/json")
  w.Write(jsonResponse)
}