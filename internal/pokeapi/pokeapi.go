// Package pokeapi makes calls to pokeapi.co
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type HTTPAPI struct{}

func (HTTPAPI) Get(u string) ([]byte, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
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
