package api

import (
	"fmt"
	"encoding/json"
)


type Cities struct {
    Cities []City `json:"cities"`
}

type City struct {
	Insee     string   `json:"insee"`
	Name      string   `json:"name"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
}


func GetCityCode(cityName string, client *Client) ([]City, error) {
	query := map[string]string{ "search": cityName }

	resp, err := client.APICall("location/cities", query)
	if err != nil {
		return nil, fmt.Errorf("city search failed: %w", err)
	}

	var cities Cities

	err = json.Unmarshal(resp, &cities)
	if err != nil {
		return nil, fmt.Errorf("error while reading JSON data: %w", err)
	}
	if len(cities.Cities) == 0 {
		return nil, fmt.Errorf("city not recognized. Try again")
	}

	return cities.Cities, nil
}