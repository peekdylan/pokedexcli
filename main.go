package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/peekdylan/pokedexcli/internal/pokeapi"
	"github.com/peekdylan/pokedexcli/internal/pokecache"
)

type config struct {
	pokeapiClient       *pokeapi.Client
	nextLocationURL     *string
	previousLocationURL *string
	caughtPokemon       map[string]pokeapi.Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
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
			description: "Displays the names of 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "View details about a caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all caught pokemon",
			callback:    commandPokedex,
		},
	}
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()

	return nil
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(cfg *config, args []string) error {
	// Fetch location areas using the client
	resp, err := cfg.pokeapiClient.GetLocationAreas(cfg.nextLocationURL)
	if err != nil {
		return err
	}

	// Update the config with new URLs
	cfg.nextLocationURL = resp.Next
	cfg.previousLocationURL = resp.Previous

	// Print all location names
	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(cfg *config, args []string) error {
	// Check if we're on the first page
	if cfg.previousLocationURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	// Fetch previous location areas using the client
	resp, err := cfg.pokeapiClient.GetLocationAreas(cfg.previousLocationURL)
	if err != nil {
		return err
	}

	// Update the config with new URLs
	cfg.nextLocationURL = resp.Next
	cfg.previousLocationURL = resp.Previous

	// Print all location names
	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a location area name")
	}

	locationName := args[0]
	fmt.Printf("Exploring %s...\n", locationName)

	// Get location area details
	locationArea, err := cfg.pokeapiClient.GetLocationArea(locationName)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a pokemon name")
	}

	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	// Get pokemon details
	pokemon, err := cfg.pokeapiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	// Calculate catch probability based on base experience
	// Higher base experience = harder to catch
	const maxBaseExperience = 300
	catchChance := 1.0 - (float64(pokemon.BaseExperience) / float64(maxBaseExperience))

	// Generate random number between 0 and 1
	randValue := rand.Float64()

	if randValue > catchChance {
		fmt.Printf("%s escaped!\n", pokemonName)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemonName)
	fmt.Println("You may now inspect it with the inspect command.")

	// Add to caught pokemon
	cfg.caughtPokemon[pokemonName] = pokemon

	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a pokemon name")
	}

	pokemonName := args[0]

	// Check if the pokemon has been caught
	pokemon, ok := cfg.caughtPokemon[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	// Display pokemon details
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, args []string) error {
	fmt.Println("Your Pokedex:")

	if len(cfg.caughtPokemon) == 0 {
		fmt.Println("  You haven't caught any pokemon yet!")
		return nil
	}

	for _, pokemon := range cfg.caughtPokemon {
		fmt.Printf(" - %s\n", pokemon.Name)
	}

	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()

	// Create cache and API client
	cache := pokecache.NewCache(5 * time.Minute)
	client := pokeapi.NewClient(cache)

	cfg := &config{
		pokeapiClient: client,
		caughtPokemon: make(map[string]pokeapi.Pokemon),
	}

	for {
		// Print the prompt
		fmt.Print("Pokedex > ")

		// Wait for user input
		scanner.Scan()

		// Get the text from the scanner
		text := scanner.Text()

		// Clean the input
		words := cleanInput(text)

		// If there are no words, continue to next iteration
		if len(words) == 0 {
			continue
		}

		// Get the command name (first word)
		commandName := words[0]

		// Get the arguments (remaining words)
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}

		// Look up the command in the registry
		command, exists := commands[commandName]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		// Execute the command's callback with config and args
		err := command.callback(cfg, args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
