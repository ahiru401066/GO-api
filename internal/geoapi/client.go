package geoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"GO-api/internal/model"
)

func FetchLocations(postalCode string) ([]model.Location, error) {
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

	var geoApiResponse model.GeoApiResponse
  if err := json.Unmarshal(body, &geoApiResponse); err != nil {
    return nil, err
  }
	return geoApiResponse.Response.Location, nil
}