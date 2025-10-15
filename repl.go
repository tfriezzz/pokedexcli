package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func runREPL() {
	userInput()
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	text = strings.Trim(text, " ")
	splitStrings := strings.Split(text, " ")

	return splitStrings
}

func userInput() string {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		input = strings.ToLower(input)
		input = strings.Trim(input, " ")
		firstWord := strings.Split(input, " ")
		fmt.Printf("Your command was: %s\n", firstWord[0])

	}
}
