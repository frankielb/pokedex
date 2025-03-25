package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/frankielb/pokedex/internal/pokeapi"
	"github.com/frankielb/pokedex/internal/pokecache"
)

var commandRegistry map[string]cliCommand

func init() {
	commandRegistry = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Show pokemon in location",
			callback:    commandExplore,
		},
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *pokeapi.Config, cache *pokecache.Cache, input string) error
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	config := &pokeapi.Config{}
	cache := pokecache.NewCache(5 * time.Second)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}
		userInput := ""
		if len(words) == 1 {
			commandName := words[0]
			command, exists := commandRegistry[commandName]
			if exists {
				err := command.callback(config, cache, userInput)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		} //adding explore 2nd input logic
		if len(words) == 2 {
			commandName := words[0]
			if commandName != "explore" {
				continue
			}
			userInput = words[1]
			command, exists := commandRegistry[commandName]
			if exists {
				err := command.callback(config, cache, userInput)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}

func cleanInput(text string) []string {

	trim := strings.TrimSpace(text)
	lower := strings.ToLower(trim)
	words := strings.Fields(lower)

	return words
}

func commandHelp(config *pokeapi.Config, cache *pokecache.Cache, userInput string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	// Iterate through commands to display them all
	for _, cmd := range commandRegistry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}
func commandExit(config *pokeapi.Config, cache *pokecache.Cache, userInput string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(config *pokeapi.Config, cache *pokecache.Cache, userInput string) error {
	locations, err := pokeapi.GetLocations(config, cache)
	if err != nil {
		return err
	}
	for _, loc := range locations {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(config *pokeapi.Config, cache *pokecache.Cache, userInput string) error {
	if config.Previous == "" {
		fmt.Println("You're on the first page")
		return nil
	}
	locations, err := pokeapi.GetPreviousLocations(config, cache)
	if err != nil {
		return err
	}
	for _, loc := range locations {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandExplore(config *pokeapi.Config, cache *pokecache.Cache, userInput string) error {
	pokemons, err := pokeapi.GetPokemons(userInput, cache)
	if err != nil {
		return err
	}
	for _, pok := range pokemons {
		fmt.Println(pok.Pokemon.Name)
	}
	return nil
}
