package main

import (
	// pokeapi "github.com/tfriezzz/pokedexcli/internal/pokeapi"
	"os"
)

func main() {
	cfg := &config{}
	runREPL(os.Stdin, os.Stdout, replCommands, cfg)
	// pokeapi.RunArea()
}
