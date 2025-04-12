package handler 

import (
	"encoding/json"
	"net/http"

	"GO-api/internal/db"
	"GO-api/internal/geoapi"
	"GO-api/internal/model"
)

func AddressHandler(w http.ResponseWriter, r *http.Request) {
  query := r.URL.Query()
  postalCode := query.Get("postal_code")

  db.AddLog(postalCode) //DBにアクセスログの追加

	locations, err := geoapi.FetchLocations(postalCode)
	if err != nil {
		http.Error(w, "Failed to fetch location", http.StatusInternalServerError)
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
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
    return 
  }
	w.Header().Set("Content-Type", "application/json")
  w.Write(jsonAddressResponse)
}