package main

import (
	"bytes"
	"strings"
	"testing"
)

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
