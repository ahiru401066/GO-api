package main

import (
  "database/sql"
  "encoding/json"
	"fmt"
  "io"
	"log"
  "math"
	"net/http"
  "strconv"

	_ "github.com/go-sql-driver/mysql"
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

type AccessLog struct {
  PostalCode string `json:"postal_code"`
  RequestCount int `json:"request_count"`
}

type AccessLogsResponse struct {
  AccessLogs []AccessLog `json:"access_logs"`
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

  addLog(postalCode) //DBにアクセスログの追加
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
  w.Write(jsonAddressResponse)
}

func addLog(postalCode string) {
  dsn := "develop:password@tcp(db:3306)/mydb"
  db, err := sql.Open("mysql", dsn)
  if err != nil {
    log.Fatal("DB接続エラー：", err)
  }
  defer db.Close()

  if err := db.Ping(); err != nil {
    log.Fatal("DBに接続できません：", err)
  }
  _, err = db.Exec("INSERT INTO access_logs (postal_code) VALUES (?)", postalCode)
  if err != nil {
    log.Fatal("INSERTエラー：", err)
  }
  fmt.Println("データを挿入しました!")
}

func accessLogsHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
  }
  db, err := sql.Open("mysql", "develop:password@tcp(db:3306)/mydb")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  rows, err := db.Query(`
		SELECT postal_code, COUNT(*) AS request_count
		FROM access_logs
		GROUP BY postal_code
		ORDER BY request_count DESC
  `)

  if err != nil {
		http.Error(w, "Failed to query DB", http.StatusInternalServerError)
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

  var accessLogs []AccessLog
  for rows.Next() {
    var accessLog AccessLog
    if err := rows.Scan(&accessLog.PostalCode, &accessLog.RequestCount);  err != nil {
      http.Error(w, "Failed to read row", http.StatusInternalServerError)
			log.Println("Row scan error:", err)
			return
    }
    accessLogs = append(accessLogs, accessLog)
  }

  response := AccessLogsResponse{
    AccessLogs: accessLogs,
  }

  w.Header().Set("Content-Type", "application/json")
  jsonResponse, err := json.Marshal(response)
  if err != nil {
    http.Error(w, "Error marshaling response", http.StatusInternalServerError)
    return    
  }
  w.Write(jsonResponse)
}


func main() {
	http.HandleFunc("/", helloHandler) // Step.1
  http.HandleFunc("/address", addressHandler)
  http.HandleFunc("/address/access_logs", accessLogsHandler)

  fmt.Println("Server is running on :8080...")
  err := http.ListenAndServe(":8080", nil)
  if err!= nil{
    panic(err)
  }
}