package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/peekdylan/pokedexcli/internal/pokeapi"
	"github.com/peekdylan/pokedexcli/internal/pokecache"
	"github.com/schollz/progressbar/v3"
)

type config struct {
	pokeapiClient       *pokeapi.Client
	nextLocationURL     *string
	previousLocationURL *string
	caughtPokemon       map[string]pokeapi.Pokemon
	party               []pokeapi.Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

var (
	cyan    = color.New(color.FgCyan).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	bold    = color.New(color.Bold).SprintFunc()
)

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
		"party": {
			name:        "party",
			description: "View your party of pokemon",
			callback:    commandParty,
		},
		"addparty": {
			name:        "addparty",
			description: "Add a pokemon to your party",
			callback:    commandAddParty,
		},
		"battle": {
			name:        "battle",
			description: "Battle a wild pokemon",
			callback:    commandBattle,
		},
		"save": {
			name:        "save",
			description: "Save your progress",
			callback:    commandSave,
		},
	}
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println(cyan("\n╔════════════════════════════════════════╗"))
	fmt.Println(cyan("║") + bold("     Welcome to the Pokedex CLI!      ") + cyan("║"))
	fmt.Println(cyan("╚════════════════════════════════════════╝"))
	fmt.Println(yellow("\nUsage:"))
	fmt.Println()

	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf("  %s - %s\n", cyan(cmd.name), cmd.description)
	}
	fmt.Println()

	return nil
}

func commandExit(cfg *config, args []string) error {
	fmt.Println(yellow("\nAuto-saving your progress..."))
	err := saveProgress(cfg)
	if err != nil {
		fmt.Printf(red("Warning: Could not save progress: %v\n"), err)
	} else {
		fmt.Println(green("✓ Progress saved!"))
	}

	fmt.Println(cyan("\nClosing the Pokedex... Goodbye!"))
	os.Exit(0)
	return nil
}

func commandMap(cfg *config, args []string) error {
	resp, err := cfg.pokeapiClient.GetLocationAreas(cfg.nextLocationURL)
	if err != nil {
		return err
	}

	cfg.nextLocationURL = resp.Next
	cfg.previousLocationURL = resp.Previous

	fmt.Println(cyan("\n📍 Location Areas:"))
	for _, location := range resp.Results {
		fmt.Printf("  • %s\n", green(location.Name))
	}
	fmt.Println()

	return nil
}

func commandMapb(cfg *config, args []string) error {
	if cfg.previousLocationURL == nil {
		fmt.Println(yellow("⚠ You're on the first page"))
		return nil
	}

	resp, err := cfg.pokeapiClient.GetLocationAreas(cfg.previousLocationURL)
	if err != nil {
		return err
	}

	cfg.nextLocationURL = resp.Next
	cfg.previousLocationURL = resp.Previous

	fmt.Println(cyan("\n📍 Location Areas:"))
	for _, location := range resp.Results {
		fmt.Printf("  • %s\n", green(location.Name))
	}
	fmt.Println()

	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println(red("❌ You must provide a location area name"))
		return nil
	}

	locationName := args[0]
	fmt.Printf(cyan("\n🔍 Exploring %s...\n"), bold(locationName))

	locationArea, err := cfg.pokeapiClient.GetLocationArea(locationName)
	if err != nil {
		return err
	}

	if len(locationArea.PokemonEncounters) == 0 {
		fmt.Println(yellow("  No Pokemon found in this area."))
		return nil
	}

	fmt.Println(green("\n✨ Found Pokemon:"))
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Printf("  • %s\n", magenta(encounter.Pokemon.Name))
	}
	fmt.Println()

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println(red("❌ You must provide a pokemon name"))
		return nil
	}

	pokemonName := args[0]

	if _, exists := cfg.caughtPokemon[pokemonName]; exists {
		fmt.Printf(yellow("\n⚠ You already caught %s!\n"), pokemonName)
		return nil
	}

	fmt.Printf(cyan("\n🎯 Throwing a Pokeball at %s...\n"), bold(pokemonName))

	pokemon, err := cfg.pokeapiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	bar := progressbar.NewOptions(100,
		progressbar.OptionSetDescription("Catching"),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(20 * time.Millisecond)
	}
	fmt.Println()

	const maxBaseExperience = 300
	catchChance := 1.0 - (float64(pokemon.BaseExperience) / float64(maxBaseExperience))
	randValue := rand.Float64()

	if randValue > catchChance {
		fmt.Printf(red("💔 %s escaped!\n"), bold(pokemonName))
		fmt.Println(yellow("   Try again!"))
		return nil
	}

	fmt.Printf(green("✓ %s was caught!\n"), bold(pokemonName))
	displayPokemonArt(pokemon.Name)
	fmt.Println(cyan("   You may now inspect it with the inspect command."))

	cfg.caughtPokemon[pokemonName] = pokemon

	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println(red("❌ You must provide a pokemon name"))
		return nil
	}

	pokemonName := args[0]

	pokemon, ok := cfg.caughtPokemon[pokemonName]
	if !ok {
		fmt.Println(red("\n❌ You have not caught that pokemon"))
		return nil
	}

	fmt.Println(cyan("\n╔════════════════════════════════════════╗"))
	fmt.Printf(cyan("║")+"  %s"+cyan("║\n"), centerText(bold(pokemon.Name), 36))
	fmt.Println(cyan("╚════════════════════════════════════════╝"))

	displayPokemonArt(pokemon.Name)

	fmt.Printf("\n%s %d\n", yellow("Height:"), pokemon.Height)
	fmt.Printf("%s %d\n", yellow("Weight:"), pokemon.Weight)

	fmt.Println(yellow("\nStats:"))
	for _, stat := range pokemon.Stats {
		fmt.Printf("  %s: %s\n",
			cyan(stat.Stat.Name),
			green(fmt.Sprintf("%d", stat.BaseStat)))
	}

	fmt.Println(yellow("\nTypes:"))
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  • %s\n", magenta(typeInfo.Type.Name))
	}
	fmt.Println()

	return nil
}

func commandPokedex(cfg *config, args []string) error {
	fmt.Println(cyan("\n╔════════════════════════════════════════╗"))
	fmt.Println(cyan("║") + bold("           Your Pokedex               ") + cyan("║"))
	fmt.Println(cyan("╚════════════════════════════════════════╝"))

	if len(cfg.caughtPokemon) == 0 {
		fmt.Println(yellow("\n  You haven't caught any pokemon yet!"))
		fmt.Println(cyan("  Use 'explore' to find pokemon, then 'catch' them!\n"))
		return nil
	}

	fmt.Printf(green("\n  Total Caught: %d\n\n"), len(cfg.caughtPokemon))
	for _, pokemon := range cfg.caughtPokemon {
		fmt.Printf("  • %s (Lvl %d)\n", magenta(pokemon.Name), pokemon.BaseExperience/10)
	}
	fmt.Println()

	return nil
}

func commandParty(cfg *config, args []string) error {
	fmt.Println(cyan("\n╔════════════════════════════════════════╗"))
	fmt.Println(cyan("║") + bold("           Your Party                 ") + cyan("║"))
	fmt.Println(cyan("╚════════════════════════════════════════╝\n"))

	if len(cfg.party) == 0 {
		fmt.Println(yellow("  Your party is empty!"))
		fmt.Println(cyan("  Use 'addparty <pokemon>' to add pokemon to your party.\n"))
		return nil
	}

	for i, pokemon := range cfg.party {
		fmt.Printf("  %d. %s (HP: %d, ATK: %d, DEF: %d)\n",
			i+1,
			bold(pokemon.Name),
			getPokemonStat(pokemon, "hp"),
			getPokemonStat(pokemon, "attack"),
			getPokemonStat(pokemon, "defense"))
	}
	fmt.Println()

	return nil
}

func commandAddParty(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println(red("❌ You must provide a pokemon name"))
		return nil
	}

	if len(cfg.party) >= 6 {
		fmt.Println(red("\n❌ Your party is full! (Max 6 Pokemon)"))
		return nil
	}

	pokemonName := args[0]
	pokemon, ok := cfg.caughtPokemon[pokemonName]
	if !ok {
		fmt.Println(red("\n❌ You haven't caught that pokemon yet!"))
		return nil
	}

	for _, p := range cfg.party {
		if p.Name == pokemonName {
			fmt.Printf(yellow("\n⚠ %s is already in your party!\n"), pokemonName)
			return nil
		}
	}

	cfg.party = append(cfg.party, pokemon)
	fmt.Printf(green("\n✓ Added %s to your party!\n"), bold(pokemonName))

	return nil
}

func commandBattle(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println(red("❌ You must provide a pokemon name to battle"))
		return nil
	}

	if len(cfg.party) == 0 {
		fmt.Println(red("\n❌ You need pokemon in your party to battle!"))
		fmt.Println(cyan("   Use 'addparty <pokemon>' to add pokemon to your party."))
		return nil
	}

	wildPokemonName := args[0]
	wildPokemon, err := cfg.pokeapiClient.GetPokemon(wildPokemonName)
	if err != nil {
		return err
	}

	return startBattle(cfg, cfg.party[0], wildPokemon)
}

func commandSave(cfg *config, args []string) error {
	fmt.Println(yellow("\n💾 Saving your progress..."))
	err := saveProgress(cfg)
	if err != nil {
		fmt.Printf(red("❌ Failed to save: %v\n"), err)
		return nil
	}
	fmt.Println(green("✓ Progress saved successfully!\n"))
	return nil
}

func main() {
	cfg := loadProgress()
	if cfg == nil {
		cache := pokecache.NewCache(5 * time.Minute)
		client := pokeapi.NewClient(cache)

		cfg = &config{
			pokeapiClient: client,
			caughtPokemon: make(map[string]pokeapi.Pokemon),
			party:         []pokeapi.Pokemon{},
		}

		fmt.Println(cyan("\n╔════════════════════════════════════════╗"))
		fmt.Println(cyan("║") + bold("     Welcome to Pokedex CLI!          ") + cyan("║"))
		fmt.Println(cyan("╚════════════════════════════════════════╝"))
		fmt.Println(yellow("\n  Starting new adventure..."))
		fmt.Println(cyan("  Type 'help' for available commands\n"))
	} else {
		fmt.Println(green("\n✓ Loaded saved progress!"))
		fmt.Printf(cyan("  Caught Pokemon: %d\n"), len(cfg.caughtPokemon))
		fmt.Printf(cyan("  Party Size: %d\n\n"), len(cfg.party))
	}

	commands := getCommands()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(bold("Pokedex > "))

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf(red("Error reading input: %v\n"), err)
			continue
		}

		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}

		command, exists := commands[commandName]
		if !exists {
			fmt.Println(red("❌ Unknown command:"), commandName)
			fmt.Println(cyan("   Type 'help' for available commands"))
			continue
		}

		err = command.callback(cfg, args)
		if err != nil {
			fmt.Printf(red("Error: %v\n"), err)
		}
	}
}

func centerText(text string, width int) string {
	textLen := len(text)
	if textLen >= width {
		return text
	}
	padding := (width - textLen) / 2
	return fmt.Sprintf("%*s%s%*s", padding, "", text, padding, "")
}

func getPokemonStat(pokemon pokeapi.Pokemon, statName string) int {
	for _, stat := range pokemon.Stats {
		if stat.Stat.Name == statName {
			return stat.BaseStat
		}
	}
	return 0
}
