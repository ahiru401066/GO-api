package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

type RequestData struct {
	Method string `json:"method"`
	Postal string `json:"postal"`
}

type Location struct {
	City string `json:"city"`
	CityKana string `json:"city_kana"`
	Town string `json:"town"`
	TownKana string `json:"town_kana"`
	X string `json:"x"`
	Y string `json:"y"`
	Prefecture string `json:"prefecture"`
	Postal string `json:"postal"`
}

type ApiResponse struct {
	Response struct {
		Location []Location `json:"location"`
	} `json:"response"`
}

func AddressHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
  postalCode := query.Get("postal_code")
	requestData := RequestData {
		Method: "searchByPostal",
		Postal: postalCode,
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		panic(err)
	}

	url := "https://geoapi.heartrails.com/api/json?method=searchByPostal"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
  if err != nil {
    panic(err)
  }
	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("%s", dump)
}



func main() {
	// http.HandleFunc("/address", AddressHandler)
	// fmt.Println("Server is running on :8080...")
	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	fmt.Println("Error starting server:", err)
	// }

	var postal_code = "5016121"
  url := fmt.Sprintf("https://geoapi.heartrails.com/api/json?method=searchByPostal&postal=%s", postal_code)
	req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    panic(err)
  }
	req.Header.Set("Content-Type", "application/json")
	
  client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}
	// fmt.Println("API Response:", apiResponse)
	fmt.Println(apiResponse.Response.Location)
	
	
}
