package main

import (
	"fmt"
	"os"
	"strings"
)

var commands map[string]cliCommand

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
		// "map": {
		// 	name: "map",
		// 	description: "Display next 20 locations",
		// 	callback: commandMap,
		// },
		// "mapb": {
		// 	name: "mapb",
		// 	description: "Display the previous 20 locations",
		// 	callback: commandMapB,
		// },
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
	fmt.Println("Usage:\n")
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
