package main

import (
	"bytes"
	"strings"
	"testing"
)

// Command test suite

func TestRunREPL(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		assert func(t *testing.T, got string)
	}{
		{
			name:  "help then exit",
			input: "help\nexit\n",
			assert: func(t *testing.T, got string) {
				if !strings.Contains(got, "Pokedex > ") {
					t.Fatal("missing prompt")
				}
				if !strings.Contains(got, "Welcome to the Pokedex!") {
					t.Fatal("missing help header")
				}
				if !strings.Contains(got, "help: ") {
					t.Fatal("missing help command")
				}
				if !strings.Contains(got, "exit: ") {
					t.Fatal("missing exit command")
				}
			},
		},
		{
			name:  "unknown then exit",
			input: "nonsense\nexit\n",
			assert: func(t *testing.T, got string) {
				if !strings.Contains(got, "Unknown command") {
					t.Fatal("expected unknown command message")
				}
			},
		},
		{
			name:  "blank input ignored",
			input: "\nexit\n",
			assert: func(t *testing.T, got string) {
				// Should show two prompts, no unknown message between.
				if strings.Count(got, "Pokedex > ") < 2 {
					t.Fatal("expected second prompt after blank line")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := strings.NewReader(tc.input)
			var out bytes.Buffer
			cfg := &config{}
			runREPL(in, &out, replCommands, cfg)
			tc.assert(t, out.String())
		})
	}
}

// Callback tests

type mockAPI struct {
	calls []string
	body  []byte
	err   error
}

func (f *mockAPI) Get(url string) ([]byte, error) {
	f.calls = append(f.calls, url)
	return f.body, f.err
}

func TestMap_FirstPage(t *testing.T) {
	page1 := []byte(`{"count":1234,"next":"https://pokeapi.co/api/v2/location-area?offset=2&limit=2","previous":null,"results":[{"name":"canalave-city-area","url":"x"},{"name":"eterna-city-area","url":"y"}]}`)
	f := &mockAPI{body: page1}

	in := strings.NewReader("map\nexit\n")
	var out bytes.Buffer
	cfg := &config{API: f}

	runREPL(in, &out, replCommands, cfg)

	got := out.String()
	if !strings.Contains(got, "canalave-city-area") || !strings.Contains(got, "eterna-city-area") {
		t.Fatalf("expected names in output, got: \n%s", got)
	}
	if cfg.Next == "" {
		t.Fatal("Next not updated")
	}
	if cfg.Previous != "" {
		t.Fatal("Previous should be nil on first page")
	}
}

func TestMapBack_FirstPage(t *testing.T) {
	f := &mockAPI{}

	in := strings.NewReader("mapb\nexit\n")
	var out bytes.Buffer
	cfg := &config{API: f}

	runREPL(in, &out, replCommands, cfg)

	got := out.String()
	if !strings.Contains(got, "you're on the first page") {
		t.Fatalf("expected: 'you're on the first page', got: %s", got)
	}
}

func TestMapBack_GoesBack(t *testing.T) {
	prevURL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=2"

	pagePrev := []byte(`{
		"count": 1234,
		"next": "https://pokeapi.co/api/v2/location-area?offset=2&limit=2",
		"previous": null,
		"results": [
			{"name":"canalave-city-area","url":"x"},
			{"name":"eterna-city-area","url":"y"}
		]
	}`)

	f := &mockAPI{body: pagePrev}

	in := strings.NewReader("mapb\nexit\n")
	var out bytes.Buffer
	cfg := &config{
		API:      f,
		Next:     "https://pokeapi.co/api/v2/location-area?offset=2&limit=2",
		Previous: "https://pokeapi.co/api/v2/location-area?offset=0&limit=2",
	}

	runREPL(in, &out, replCommands, cfg)

	got := out.String()
	if !strings.Contains(got, "canalave-city-area") {
		t.Fatalf("expected names in output, got %s", got)
	}
	if len(f.calls) != 1 {
		t.Fatalf("expected 1 API call, got %d: %v", len(f.calls), f.calls)
	}
	if f.calls[0] != prevURL {
		t.Fatalf("expected call to %q, got %q", prevURL, f.calls[0])
	}
	if cfg.Next == "" {
		t.Fatal("Next not updated after mapb")
	}
	if cfg.Previous != "" {
		t.Fatal("Previous should be empty on the first page after going back")
	}
}
