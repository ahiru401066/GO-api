package main

import (
  "encoding/json"
	"fmt"
  "io"
  "math"
	"net/http"
  "strconv"
)

type GeoApiResponse struct {
  Response struct {
    Location []Location `json:"location"`
  } `json:"response"`
}

type AddressResponse struct {
  PostalCode string `json:"postal_code"`
  HitCount int `json:"hit_count"`
  Address string `json:"address"`
  FromTokyoDistance float64 `json:"tokyo_sta_distance"`
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

func getHitCount(locations []Location)int {
  return len(locations)
}

func getCommonAddress(locations []Location)string {
  if len(locations) == 0 {
    return ""
  }
  pref := locations[0].Prefecture
  city := locations[0].City
  townPrefix := locations[0].Town

  for _, loc := range locations[1:] {
    if loc.Prefecture != pref {
      pref = ""
    }
    if loc.City != city {
      city = ""
    }
    townPrefix = getCommonPrefix(townPrefix, loc.Town)
  }

  result := pref + city
  if townPrefix != "" {
    result += townPrefix
  }
  return result
}

func getCommonPrefix(a, b string) string {
  arune := []rune(a)
  brune := []rune(b)
  minLen := len(arune)
  if len(brune) < len(arune) {
    minLen = len(b)
  } 
  i := 0
  for i < minLen && arune[i] == brune[i] {
    i++
  }
  return string(arune[:i])
}

func getFromTokyoStation(locations []Location) float64 {
  var maxDistance float64 = 0
  for _, loc := range locations {
    x, _ := strconv.ParseFloat(loc.X, 64)
    y, _ := strconv.ParseFloat(loc.Y, 64)
    dis := calculateDistanceFromTokyoStation(x, y)
    if dis > maxDistance {
      maxDistance = dis
    }
  }
  return maxDistance
}

func calculateDistanceFromTokyoStation(x, y float64) float64 {
  const (  
    R = 6371.0
    xt = 139.7673068
    yt = 35.6809591
  )
  distance := (math.Pi / 180) * R * math.Sqrt(
    math.Pow((x - xt) * math.Cos(math.Pi*(y + yt)/360),2) + math.Pow(y - yt, 2),
  )
  return RoundToDecimal(distance)
}

func RoundToDecimal(f float64) float64 {
  return math.Round(f*10) / 10
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

  var geoApiResponse GeoApiResponse
  err = json.Unmarshal(body, &geoApiResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}

  hitCount := getHitCount(geoApiResponse.Response.Location)
  commonAddress := getCommonAddress(geoApiResponse.Response.Location)
  DistanceFromTokyo := getFromTokyoStation(geoApiResponse.Response.Location)

  addressResponse := AddressResponse{
    PostalCode: postalCode,
    HitCount: hitCount,
    Address: commonAddress,
    FromTokyoDistance: DistanceFromTokyo,
  }

  w.Header().Set("Content-Type", "application/json")
  jsonAddressResponse, err := json.Marshal(addressResponse)
  if err != nil {
    http.Error(w, "Error marshaling response", http.StatusInternalServerError)
    return 
  }
  // fmt.Fprintln(w, string(jsonAddressResponse))
  w.Write(jsonAddressResponse)
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