package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var replCommands map[string]cliCommand

func init() {
	replCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}
}

func runREPL() {
	userInput()
}

// func cleanInput(text string) []string {
// 	text = strings.ToLower(text)
// 	text = strings.Trim(text, " ")
// 	splitStrings := strings.Split(text, " ")
//
// 	return splitStrings
// }

func userInput() string {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		input = strings.ToLower(input)
		input = strings.Trim(input, " ")
		firstWord := strings.Split(input, " ")
		switch firstWord[0] {
		case "exit":
			replCommands["exit"].callback()
		case "help":
			replCommands["help"].callback()
		default:
			fmt.Println("Unknown command")
		}

	}
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print(`Usage:

`)
	for _, cmd := range replCommands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}
