package main

import (
	"time"

	"github.com/mnisyif/pokedexcli/internal/pokeapi"
	"github.com/mnisyif/pokedexcli/internal/pokecache"
)

func main() {
	newCache, _ := pokecache.NewCache(5 * time.Second)
	newClient := &pokeapi.Client{
		Cache: newCache,
	}
	cfg := &config{
		client: newClient,
	}

	startREPL(cfg)
}
