package pokeapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/frankielb/pokedex/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2/location-area/"

func GetPokemons(name string, cache *pokecache.Cache) ([]PokemonEncounter, error) {
	url := baseURL + name + "/"

	var reader io.Reader
	var body []byte

	stored, exists := cache.Get(url)
	if exists { //wether cached or not the data structure is the same
		reader = bytes.NewReader(stored)
	} else { //adds the non-cached to cache
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		cache.Add(url, body)
		reader = bytes.NewReader(body)
	}
	var namedResponse NamedLocationResponse
	if err := json.NewDecoder(reader).Decode(&namedResponse); err != nil {
		return nil, err
	}
	return namedResponse.PokemonEncounters, nil
}

func GetLocations(config *Config, cache *pokecache.Cache) ([]LocationArea, error) {
	url := baseURL

	if config.Next != "" {
		url = config.Next
	}

	var reader io.Reader
	var body []byte
	stored, exists := cache.Get(url)
	if exists { //wether cached or not the data structure is the same
		reader = bytes.NewReader(stored)
	} else { //adds the non-cached to cache
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		cache.Add(url, body)
		reader = bytes.NewReader(body)
	}

	var locationResp LocationAreaResponse
	if err := json.NewDecoder(reader).Decode(&locationResp); err != nil {
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

func GetPreviousLocations(config *Config, cache *pokecache.Cache) ([]LocationArea, error) {
	if config.Previous == "" {
		return nil, nil
	}

	var reader io.Reader
	var body []byte
	stored, exists := cache.Get(config.Previous)
	if exists { //wether cached or not the data structure is the same
		reader = bytes.NewReader(stored)
	} else { //adds the non-cached to cache
		res, err := http.Get(config.Previous)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		cache.Add(config.Previous, body)
		reader = bytes.NewReader(body)
	}

	var locationResp LocationAreaResponse
	if err := json.NewDecoder(reader).Decode(&locationResp); err != nil {
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
