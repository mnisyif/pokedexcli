// Package pokeapi provides with api calls to comment with PokeAPI servers
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mnisyif/pokedexcli/internal/pokecache"
)

type Client struct {
	Cache *pokecache.Cache
}

type PokeLocationAreas struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokemonEncounter struct {
	Pokemon struct {
		Name *string `json:"name"`
		URL  *string `json:"url"`
	} `json:"pokemon"`
}

type LocationAreaDetails struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonDetails struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
}

func FetchAndCache[T any](cache *pokecache.Cache, url string) (T, error) {
	var result T

	cached, exists := cache.Get(url)
	if exists {
		err := json.Unmarshal(cached, &result)
		return result, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, &result)

	cache.Add(url, data)

	return result, err
}

func (c *Client) FetchLocations(pageURL *string) (PokeLocationAreas, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	return FetchAndCache[PokeLocationAreas](c.Cache, url)
}

func (c *Client) EncounterPokemons(locationID *string) (LocationAreaDetails, error) {
	url := fmt.Sprintf("/%s/location-area/%s", baseURL, *locationID)

	return FetchAndCache[LocationAreaDetails](c.Cache, url)
}

func (c *Client) FetchPokemon(pokemonName *string) (PokemonDetails, error) {
	url := fmt.Sprint("/%s/pokemon/%s", baseURL, *pokemonName)

	return FetchAndCache[PokemonDetails](c.Cache, url)
}
