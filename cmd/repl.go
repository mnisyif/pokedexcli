package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mnisyif/pokedexcli/internal/pokeapi"
)

var commands map[string]cliCommand

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	client              *pokeapi.Client
	nextLocationURL     *string
	previousLocationURL *string
	locationArea        *string
}

func startREPL(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}

		cmd, ok := getCommands()[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if len(words) >= 2 {
			cfg.locationArea = &words[1]
		}

		err := cmd.callback(cfg)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getCommands() map[string]cliCommand {
	commands = map[string]cliCommand{
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
			description: "Display next 20 locations",
			callback:    commandMapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "List all pokemons in this area",
			callback:    commandExplore,
		},
	}
	return commands
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMapf(cfg *config) error {
	locations, err := cfg.client.FetchLocations(cfg.nextLocationURL)
	if err != nil {
		return err
	}

	cfg.nextLocationURL = locations.Next
	cfg.previousLocationURL = locations.Previous

	for _, location := range locations.Results {
		fmt.Println(location.Name)
		// fmt.Println(location.Name)
	}
	return nil
}

func commandMapb(cfg *config) error {
	if cfg.previousLocationURL == nil {
		return errors.New("you're on the first page")
	}

	locations, err := cfg.client.FetchLocations(cfg.previousLocationURL)
	if err != nil {
		return err
	}

	cfg.nextLocationURL = locations.Next
	cfg.previousLocationURL = locations.Previous

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(cfg *config) error {
	pokemons, err := cfg.client.EncounterPokemons(cfg.locationArea)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", *cfg.locationArea)
	for _, value := range pokemons.PokemonEncounters {
		fmt.Printf("  - %s\n", *value.Pokemon.Name)
	}

	return nil
}

func cleanInput(text string) []string {
	var result []string
	sep := " "
	i := strings.Index(text, sep)

	for i > -1 {
		word := strings.ToLower(text[:i])
		if word != "" {
			result = append(result, strings.ToLower(text[:i]))
		}
		text = text[i+len(sep):]
		i = strings.Index(text, sep)
	}

	if text != "" {
		result = append(result, strings.ToLower(text))
	}

	return result
}
