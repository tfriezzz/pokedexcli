// Package pokeapi makes calls to pokeapi.co
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"

	// "encoding/json"
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

func RunArea() {
	// FetchLocationAreas()
}

func FetchLocationAreas(in io.Reader, out io.Writer, url string) (response, error) {
	api := "https://pokeapi.co/api/v2/location-area"
	if url == "" {
		url = api
	}
	res, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(out, "error creating request: %v", err)
		return response{}, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(out, "response failed: %v", err)
		return response{}, fmt.Errorf("response failed: %w", err)
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		fmt.Fprintf(out, "response failed with status code: %d and\nbody: %s", res.StatusCode, body)
		return response{}, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}

	var resp response
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Fprintf(out, "unmarshal failed: %v", err)
		return response{}, fmt.Errorf("unmarshal failed: %w", err)
	}
	return resp, nil
}
