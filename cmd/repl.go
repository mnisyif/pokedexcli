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
	callback    func(*config, ...string) error
}

type config struct {
	client              *pokeapi.Client
	nextLocationURL     *string
	previousLocationURL *string
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

		command := words[0]
		args := []string{}

		if len(words) > 1 {
			args = words[1:]
		}

		cmd, ok := getCommands()[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := cmd.callback(cfg, args...)
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
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
	}
	return commands
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMapf(cfg *config, args ...string) error {
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

func commandMapb(cfg *config, args ...string) error {
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

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a location name")
	}

	name := args[0]
	pokemons, err := cfg.client.EncounterPokemons(&name)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", name)
	fmt.Println("Found Pokemon: ")
	for _, value := range pokemons.PokemonEncounters {
		fmt.Printf("  - %s\n", *value.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}

	name := args[0]
	pokemon, err := cfg.client.FetchPokemon(&name)
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Base Experience: %d\n", pokemon.BaseExperience)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Printf("Throwing ball at %s...\n", args[0])
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
