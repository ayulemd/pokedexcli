package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ayulemd/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config) error
}

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	words := strings.Fields(lowered)

	return words
}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	config := &pokeapi.Config{
		Next:     "https://pokeapi.co/api/v2/location-area",
		Previous: "",
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()

		words := cleanInput(text)
		if len(words) > 0 {
			command, exists := getCommands()[words[0]]
			if !exists {
				fmt.Println("Unknown command")
				continue
			}

			err := command.callback(config)
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
		}
	}
}

func commandExit(config *pokeapi.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(config *pokeapi.Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}

func commandMap(config *pokeapi.Config) error {
	if config.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}

	locationArea, err := pokeapi.GetLocationArea(config.Next)
	if err != nil {
		return err
	}

	pokeapi.DisplayLocationAreas(locationArea, config)

	return nil
}

func commandMapB(config *pokeapi.Config) error {
	if config.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	locationArea, err := pokeapi.GetLocationArea(config.Previous)
	if err != nil {
		return err
	}

	pokeapi.DisplayLocationAreas(locationArea, config)

	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 areas",
			callback:    commandMapB,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}
