package main

import (
	// "bytes"
	"encoding/json"
	"fmt"
	// "io"
	"net/http"
)

type RequestData struct {
	Method string `json:"method"`
	Postal string `json:"postal"`
}

func AddressHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
  postalCode := query.Get("postal_code")
	requestData := RequestData {
		Method: "searchByPostal",
		Postal: postalCode,
	}

	jsonData, _ := json.Marshal(requestData)

	fmt.Fprint(w, string(jsonData))
}



func main() {
	http.HandleFunc("/address", AddressHandler)

	fmt.Println("Server is running on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
