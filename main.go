package main

import (
	"os"
	"time"

	"github.com/tfriezzz/pokedexcli/internal/pokeapi"
	"github.com/tfriezzz/pokedexcli/internal/pokecache"
)

func main() {
	cache := pokecache.NewCache(5 * time.Second)
	// pokedex := pokeapi.NewPokedex()
	cfg := &config{
		API: &pokeapi.HTTPAPI{
			Cache: cache,
			// Pokedex: pokedex,
		},
	}
	runREPL(os.Stdin, os.Stdout, replCommands, cfg)
	// pokeapi.RunArea()
}
