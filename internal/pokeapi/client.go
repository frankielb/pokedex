package pokeapi

import (
	"encoding/json"
	"net/http"
)

const baseURL = "https://pokeapi.co/api/v2/location-area/"

func GetLocations(config *Config) ([]LocationArea, error) {
	url := baseURL

	if config.Next != "" {
		url = config.Next
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var locationResp LocationAreaResponse
	if err := json.NewDecoder(res.Body).Decode(&locationResp); err != nil {
		return nil, err
	}

	if locationResp.Previous != nil {
		config.Previous = *locationResp.Previous
	} else {
		config.Previous = ""
	}

	if locationResp.Next != nil {
		config.Next = *locationResp.Next
	} else {
		config.Next = ""
	}

	return locationResp.Results, nil
}

func GetPreviousLocations(config *Config) ([]LocationArea, error) {
	if config.Previous == "" {
		return nil, nil
	}
	res, err := http.Get(config.Previous)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var locationResp LocationAreaResponse
	if err := json.NewDecoder(res.Body).Decode(&locationResp); err != nil {
		return nil, err
	}

	if locationResp.Previous != nil {
		config.Previous = *locationResp.Previous
	} else {
		config.Previous = ""
	}

	if locationResp.Next != nil {
		config.Next = *locationResp.Next
	} else {
		config.Next = ""
	}
	return locationResp.Results, nil
}
