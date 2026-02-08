package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/peekdylan/pokedexcli/internal/pokeapi"
	"github.com/peekdylan/pokedexcli/internal/pokecache"
)

const saveFile = "pokedex_save.json"

type SaveData struct {
	CaughtPokemon map[string]pokeapi.Pokemon `json:"caught_pokemon"`
	Party         []pokeapi.Pokemon          `json:"party"`
	Timestamp     time.Time                  `json:"timestamp"`
}

func saveProgress(cfg *config) error {
	saveData := SaveData{
		CaughtPokemon: cfg.caughtPokemon,
		Party:         cfg.party,
		Timestamp:     time.Now(),
	}

	data, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %w", err)
	}

	err = os.WriteFile(saveFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

func loadProgress() *config {
	data, err := os.ReadFile(saveFile)
	if err != nil {
		// No save file exists
		return nil
	}

	var saveData SaveData
	err = json.Unmarshal(data, &saveData)
	if err != nil {
		fmt.Printf("Warning: Could not load save file: %v\n", err)
		return nil
	}

	// Create new config with loaded data
	cache := pokecache.NewCache(5 * time.Minute)
	client := pokeapi.NewClient(cache)

	return &config{
		pokeapiClient: client,
		caughtPokemon: saveData.CaughtPokemon,
		party:         saveData.Party,
	}
}
