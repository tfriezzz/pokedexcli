package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"math/rand"
	"strings"

	"github.com/tfriezzz/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(in io.Reader, out io.Writer, cfg *config) (done bool, err error)
}

type api interface {
	Get(string) ([]byte, error)
	AddToPokedex(pokeapi.Pokemon)
}

type config struct {
	Next           string
	Previous       string
	API            *pokeapi.HTTPAPI
	secondArgument string
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
		"catch": {
			name:        "catch",
			description: "'catch <pokemon_name>' attempts to catch the Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "'inspect <pokemon_name>' inspect the selected Pokemon in you Pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "lists all the Pokemon in you Pokedex",
			callback:    commandPokedex,
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
			cfg.secondArgument = fields[1]
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
	resp, err := pokeapi.FetchEncounters(cfg.API.Get, cfg.secondArgument)
	if err != nil {
		return false, err
	}

	for _, encounter := range resp.PokemonEncounters {
		fmt.Fprintf(out, "%v\n", encounter.Pokemon.Name)
	}

	return false, nil
}

func commandCatch(in io.Reader, out io.Writer, cfg *config) (bool, error) {
	resp, err := pokeapi.FetchPokemon(cfg.API.Get, cfg.secondArgument)
	if err != nil {
		return false, err
	}
	didCatch := func() bool {
		baseChance := 90.0
		difficulty := 0.2 * float64(resp.BaseExperience)
		maxDifficulty := 95.0
		catchChance := baseChance - difficulty
		result := int(math.Min(catchChance, maxDifficulty))

		return rand.Intn(100) < result
	}
	fmt.Fprintf(out, "Throwing a Pokeball at %s...\n", resp.Name)
	if !didCatch() {
		fmt.Fprintf(out, "%s escaped!\n", resp.Name)
		return false, nil
	}
	cfg.API.AddToPokedex(resp)
	fmt.Fprintf(out, "%s was caught!\n", resp.Name)
	fmt.Fprint(out, "You may now inspect it with the inspect command\n.")
	return false, nil
}

func commandInspect(in io.Reader, out io.Writer, cfg *config) (bool, error) {
	// pokemon := cfg.API.Pokedex[cfg.secondArgument]
	pokemon, ok := cfg.API.Pokedex[cfg.secondArgument]
	if !ok {
		fmt.Fprint(out, "you have not caught that pokemon\n")
		return false, nil
	}
	fmt.Fprintf(out, "Name: %s\n", pokemon.Name)
	fmt.Fprintf(out, "Height: %d\n", pokemon.Height)
	fmt.Fprintf(out, "Weight: %d\n", pokemon.Weight)
	fmt.Fprint(out, "Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Fprintf(out, " -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Fprint(out, "Types:\n")
	for _, t := range pokemon.Types {
		fmt.Fprintf(out, " -%s\n", t.Type.Name)
	}
	return false, nil
}

func commandPokedex(in io.Reader, out io.Writer, cfg *config) (bool, error) {
	pokemonList := cfg.API.Pokedex
	if len(pokemonList) == 0 {
		fmt.Fprint(out, "Try to catch some Pokemon first\n")
		return false, nil
	}
	fmt.Fprint(out, "Your Pokedex:\n")
	for p := range pokemonList {
		fmt.Fprintf(out, "%s\n", p)
	}
	return false, nil
}
