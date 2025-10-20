package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/tfriezzz/pokedexcli/internal/pokeapi"
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
	API      api
	Location string
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
			description: "Displays 20 location areas. (Each subsequent call displays the next 20)",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the Previous 20 location areas.",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "'explore <area_name>' lists all the Pokemon in the area",
			callback:    commandExplore,
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
		if len(fields) >= 2 {
			cfg.Location = fields[1]
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
	fmt.Fprintln(out, "\nWelcome to the Pokedex!")
	fmt.Fprint(out, `Usage:

`)
	for _, cmd := range replCommands {
		fmt.Fprintf(out, "%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Fprint(out, "\n")

	return false, nil
}

func val(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func commandMap(in io.Reader, out io.Writer, cfg *config) (bool, error) {
	resp, err := pokeapi.FetchLocationAreas(cfg.API.Get, cfg.Next)
	if err != nil {
		return false, err
	}
	for _, r := range resp.Results {
		fmt.Fprintf(out, "%v\n", r.Name)
	}
	cfg.Next = val(resp.Next)
	cfg.Previous = val(resp.Previous)

	return false, nil
}

func commandMapb(in io.Reader, out io.Writer, cfg *config) (bool, error) {
	if cfg.Previous == "" {
		fmt.Fprintln(out, "you're on the first page")
		return false, nil
	}

	resp, err := pokeapi.FetchLocationAreas(cfg.API.Get, cfg.Previous)
	if err != nil {
		return false, err
	}
	fmt.Fprintln(out, "you're on the first page")

	for _, r := range resp.Results {
		fmt.Fprintf(out, "%v\n", r.Name)
	}
	cfg.Next = val(resp.Next)
	cfg.Previous = val(resp.Previous)

	return false, nil
}

func commandExplore(in io.Reader, out io.Writer, cfg *config) (bool, error) {
	resp, err := pokeapi.FetchEncounters(cfg.API.Get, cfg.Location)
	if err != nil {
		return false, err
	}

	for _, encounter := range resp.PokemonEncounters {
		fmt.Fprintf(out, "%v\n", encounter.Pokemon.Name)
	}

	return false, nil
}
