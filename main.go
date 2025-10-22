package main

import (
	"os"
	"time"

	"github.com/tfriezzz/pokedexcli/internal/pokeapi"
	"github.com/tfriezzz/pokedexcli/internal/pokecache"
)

func main() {
	cache := pokecache.NewCache(5 * time.Second)
	cfg := &config{
		API: &pokeapi.HTTPAPI{
			Cache:   cache,
			Pokedex: make(map[string]pokeapi.Pokemon),
		},
	}
	runREPL(os.Stdin, os.Stdout, replCommands, cfg)
}
