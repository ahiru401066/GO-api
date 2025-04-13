package geoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"GO-api/internal/model"
)

func FetchLocations(postalCode string) ([]model.Location, error) {
	//リクエストの作成
	url := fmt.Sprintf("https://geoapi.heartrails.com/api/json?method=searchByPostal&postal=%s", postalCode)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	//　リクエストの送信とレスポンスの受信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received no-ok status")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geoApiResponse model.GeoApiResponse
	if err := json.Unmarshal(body, &geoApiResponse); err != nil {
		return nil, err
	}
	return geoApiResponse.Response.Location, nil
}
