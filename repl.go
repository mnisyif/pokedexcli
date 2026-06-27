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
	callback    func() error
}

func init() {
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
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
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
