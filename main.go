package main

import (
	"os"

	"github.com/tfriezzz/pokedexcli/internal/pokeapi"
)

func main() {
	cfg := &config{
		API: pokeapi.HTTPAPI{},
	}
	runREPL(os.Stdin, os.Stdout, replCommands, cfg)
	// pokeapi.RunArea()
}
