package pokeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"

	"github.com/frankielb/pokedex/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2/location-area/"
const pokemonURL = "https://pokeapi.co/api/v2/pokemon/"

func AttemptCatch(name string) (bool, PokemonExp, error) {
	url := pokemonURL + name + "/"
	res, err := http.Get(url)
	if err != nil {

		return false, PokemonExp{}, err
	}
	defer res.Body.Close()
	var pokemonExp PokemonExp
	if err := json.NewDecoder(res.Body).Decode(&pokemonExp); err != nil {
		fmt.Printf("%v not found\n", name)
		return false, PokemonExp{}, err
	}
	//lower exp, more likely
	chance := max(100-pokemonExp.BaseExp/3, 10)
	randN := rand.IntN(100) + 1
	return randN <= chance, pokemonExp, nil
}

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
