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
		"catch": {
			name:        "catch",
			description: "Try and catch the pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Show details of pokemon in pokedex",
			callback:    commandInspect,
		},
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *pokeapi.Config, cache *pokecache.Cache, input string, pokedex *Pokedex) error
}

type Pokedex struct {
	caughtPokemon map[string]pokeapi.PokemonExp
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	config := &pokeapi.Config{}
	cache := pokecache.NewCache(5 * time.Second)
	pokedex := &Pokedex{
		caughtPokemon: make(map[string]pokeapi.PokemonExp),
	}
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
				err := command.callback(config, cache, userInput, pokedex)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		} //adding explore 2nd input logic
		if len(words) == 2 {
			commandName := words[0]
			if commandName != "explore" && commandName != "catch" && commandName != "inspect" {
				continue
			}
			userInput = words[1]
			command, exists := commandRegistry[commandName]
			if exists {
				err := command.callback(config, cache, userInput, pokedex)
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

func commandHelp(config *pokeapi.Config, cache *pokecache.Cache, userInput string, pokedex *Pokedex) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	// Iterate through commands to display them all
	for _, cmd := range commandRegistry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}
func commandExit(config *pokeapi.Config, cache *pokecache.Cache, userInput string, pokedex *Pokedex) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(config *pokeapi.Config, cache *pokecache.Cache, userInput string, pokedex *Pokedex) error {
	locations, err := pokeapi.GetLocations(config, cache)
	if err != nil {
		return err
	}
	for _, loc := range locations {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(config *pokeapi.Config, cache *pokecache.Cache, userInput string, pokedex *Pokedex) error {
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

func commandExplore(config *pokeapi.Config, cache *pokecache.Cache, userInput string, pokedex *Pokedex) error {
	pokemons, err := pokeapi.GetPokemons(userInput, cache)
	if err != nil {
		return err
	}
	for _, pok := range pokemons {
		fmt.Println(pok.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *pokeapi.Config, cache *pokecache.Cache, userInput string, pokedex *Pokedex) error {
	fmt.Printf("Throwing a Pokeball at %v...\n", userInput)
	caught, pokemonInfo, err := pokeapi.AttemptCatch(userInput)
	if err != nil {
		return err
	}
	if !caught {
		fmt.Printf("%v escaped!\n", userInput)
	} else {
		fmt.Printf("%v was caught!\n", userInput)
		pokedex.caughtPokemon[userInput] = pokemonInfo
	}
	return nil
}

func commandInspect(config *pokeapi.Config, cache *pokecache.Cache, userInput string, pokedex *Pokedex) error {
	pokemonData, exists := pokedex.caughtPokemon[userInput]
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %v\n", pokemonData.Name)
	fmt.Printf("Height: %v\n", pokemonData.Height)
	fmt.Printf("Weight: %v\n", pokemonData.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemonData.Stats {
		fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemonData.Types {
		fmt.Printf("  - %v\n", t.Type.Name)
	}
	return nil

}
