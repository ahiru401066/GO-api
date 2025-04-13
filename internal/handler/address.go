package handler 

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"GO-api/internal/db"
	"GO-api/internal/geoapi"
	"GO-api/internal/model"
)

func validatePostalCode(r *http.Request) (string, error) {
	query := r.URL.Query()
  postalCode := query.Get("postal_code")
	if !isValidPostalCode(postalCode) {
		return "", fmt.Errorf("Bad Request")
	}
	return postalCode, nil
}

func isValidPostalCode(postalCode string) bool {
	re := regexp.MustCompile(`\d{7}$`)
	return re.MatchString(postalCode)
}

func AddressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return 
	}

	// 郵便番号のバリテーション
  postalCode, err := validatePostalCode(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	// DBにアクセスログの追加
	if err := db.CreateLog(postalCode); err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	// 外部apiを利用し住所一覧の取得
	locations, err := geoapi.FetchLocations(postalCode)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return 
	}

	addressResponse := model.AddressResponse{
		PostalCode: postalCode,
		HitCount: model.GetHitCount(locations),
		Address: model.GetCommonAddress(locations),
		FromTokyoDistance: model.GetFromTokyoStation(locations),
	}

  jsonAddressResponse, err := json.Marshal(addressResponse)
  if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
    return 
  }
	w.Header().Set("Content-Type", "application/json")
  w.Write(jsonAddressResponse)
}