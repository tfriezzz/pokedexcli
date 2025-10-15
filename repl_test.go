package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		}, {
			input:    "test case",
			expected: []string{"test", "case"},
		}, {
			input:    "bulbasaur charmander squirtle",
			expected: []string{"bulbasaur", "charmander", "squirtle"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("length mismatch: got %d, want %d", len(actual), len(c.expected))
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("word %d mismatch: got %q, want %q", i, word, expectedWord)
			}
		}
	}
}
