// Package pokeapi makes calls to pokeapi.co
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tfriezzz/pokedexcli/internal/pokecache"
)

type response struct {
	Count    int       `json:"count"`
	Next     *string   `json:"next"`
	Previous *string   `json:"previous"`
	Results  []Results `json:"results"`
}
type Results struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type HTTPAPI struct {
	Cache *pokecache.Cache
}

type nameCall struct {
	// Location             Location               `json:"location"`
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}

type PokemonEncounters struct {
	Pokemon PokemonInfo `json:"pokemon"`
}
type PokemonInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (a *HTTPAPI) Get(u string) ([]byte, error) {
	if a.Cache != nil {
		if val, ok := a.Cache.Get(u); ok {
			return val, nil
		}
	}
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if a.Cache != nil {
		a.Cache.Add(u, body)
	}
	return body, nil
}

func FetchLocationAreas(get func(string) ([]byte, error), url string) (response, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	body, err := get(url)
	if err != nil {
		return response{}, fmt.Errorf("request failed: %w", err)
	}

	var resp response
	if err := json.Unmarshal(body, &resp); err != nil {
		return response{}, fmt.Errorf("unmarshal failed: %w", err)
	}
	return resp, nil
}

func FetchEncounters(get func(string) ([]byte, error), location string) (nameCall, error) {
	baseURL := "https://pokeapi.co/api/v2/location-area/"
	exploreLocation := fmt.Sprintf("%s%s", baseURL, location)
	// fmt.Printf("test_location: %s", exploreLocation)
	body, err := get(exploreLocation)
	if err != nil {
		return nameCall{}, fmt.Errorf("request failed %w", err)
	}
	var resp nameCall
	if err := json.Unmarshal(body, &resp); err != nil {
		return nameCall{}, fmt.Errorf("unmarshal failed: %w", err)
	}
	return resp, nil
}
