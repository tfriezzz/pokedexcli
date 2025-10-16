package main

import (
	"bufio"
	"fmt"
	"io"

	// "os"
	"strings"

	pokeapi "github.com/tfriezzz/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(in io.Reader, out io.Writer, cfg *config) (done bool, err error)
}

type api interface{ Get(string) ([]byte, error) }

type config struct {
	Next     string
	Previous string
	URL      string
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
		"map": {
			name:        "map",
			description: "Displays 20 location areas. Each subsequent call displays the next 20",
			callback:    commandMap,
		},
	}
}

func runREPL(in io.Reader, out io.Writer, cmds map[string]cliCommand, cfg *config) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, "Pokedex > ")
		if !scanner.Scan() {
			return
		}
		input := strings.TrimSpace(strings.ToLower(scanner.Text()))
		fields := strings.Fields(input)
		if len(fields) == 0 {
			continue
		}
		name := fields[0]
		cmd, ok := cmds[name]
		if !ok {
			fmt.Fprintln(out, "Unknown command")
			continue
		}
		done, err := cmd.callback(in, out, cfg)
		if err != nil {
			fmt.Fprint(out, "error:", err)
		}
		if done {
			return
		}
	}
}

func commandExit(in io.Reader, out io.Writer, cfg *config) (bool, error) {
	fmt.Fprint(out, "Closing the Pokedex... Goodbye!")

	return true, nil
}

func commandHelp(in io.Reader, out io.Writer, cfg *config) (bool, error) {
	fmt.Fprintln(out, "Welcome to the Pokedex!")
	fmt.Fprint(out, `Usage:

`)
	for _, cmd := range replCommands {
		fmt.Fprintf(out, "%s: %s\n", cmd.name, cmd.description)
	}

	return false, nil
}

func val(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func commandMap(in io.Reader, out io.Writer, cfg *config) (bool, error) {
	resp, err := pokeapi.FetchLocationAreas(in, out, cfg.Next)
	if err != nil {
		return false, err
	}
	results := resp.Results
	for _, r := range results {
		fmt.Fprintf(out, "%v\n", r.Name)
	}
	cfg.Next = val(resp.Next)
	cfg.Previous = val(resp.Previous)

	return false, nil
}
