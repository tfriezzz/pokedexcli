package main

import (
	"strings"
)

func cleanInput(text string) []string {
	// empty := []string{text}
	text = strings.ToLower(text)
	text = strings.Trim(text, " ")
	splitStrings := strings.Split(text, " ")

	return splitStrings
}

func main() {
	cleanInput("")
}
