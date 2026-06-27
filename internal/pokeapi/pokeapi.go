// Package pokeapi provides with api calls to comment with PokeAPI servers
package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type PokeLocationAreas struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func FetchLocations(pageURL *string) (PokeLocationAreas, error) {
	url := baseURL
	if pageURL != nil {
		url = *pageURL
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokeLocationAreas{}, err
	}

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return PokeLocationAreas{}, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PokeLocationAreas{}, err
	}

	var locationAreas PokeLocationAreas
	err = json.Unmarshal(data, &locationAreas)
	if err != nil {
		return PokeLocationAreas{}, err
	}

	return locationAreas, nil
}
