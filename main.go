package main

import (
  "encoding/json"
	"fmt"
  "io"
	"net/http"
)

type LocationResponse struct {
  Response struct {
    Location []Location `json:"location"`
  } `json:"response`
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

//step.1
func helloHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
  }
  fmt.Fprintln(w, "Hello go!")
}


func addressHandler(w http.ResponseWriter, r *http.Request) {
  query := r.URL.Query()
  postalCode := query.Get("postal_code")
  url := fmt.Sprintf("https://geoapi.heartrails.com/api/json?method=searchByPostal&postal=%s", postalCode)
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    panic(err)
  }
  client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

  var locationResponse LocationResponse
  err = json.Unmarshal(body, &locationResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}

  fmt.Fprintln(w,locationResponse)
}


func main() {
	http.HandleFunc("/", helloHandler) // Step.1
  http.HandleFunc("/address", addressHandler)

  fmt.Println("Server is running on :8080...")
  err := http.ListenAndServe(":8080", nil)
  if err!= nil{
    panic(err)
  }
}