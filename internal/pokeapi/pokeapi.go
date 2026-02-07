package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

const baseURL = "https://pokeapi.co/api/v2"

// LocationAreasResponse represents the API response for location areas
type LocationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// LocationArea represents detailed information about a specific location area
type LocationArea struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// Pokemon represents a Pokemon's details
type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

// Client handles API requests with caching
type Client struct {
	cache Cache
}

// Cache interface defines the caching methods
type Cache interface {
	Add(key string, val []byte)
	Get(key string) ([]byte, bool)
}

// NewClient creates a new PokeAPI client with a cache
func NewClient(cache Cache) *Client {
	return &Client{
		cache: cache,
	}
}

// GetLocationAreas fetches location areas from the PokeAPI with caching
func (c *Client) GetLocationAreas(pageURL *string) (*LocationAreasResponse, error) {
	// Use the provided URL or default to the first page
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	// Check if we have this URL in cache
	if cachedData, ok := c.cache.Get(url); ok {
		// Unmarshal cached data
		var locationsResp LocationAreasResponse
		err := json.Unmarshal(cachedData, &locationsResp)
		if err != nil {
			return nil, err
		}
		return &locationsResp, nil
	}

	// Make the HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Add to cache
	c.cache.Add(url, body)

	// Unmarshal JSON into our struct
	var locationsResp LocationAreasResponse
	err = json.Unmarshal(body, &locationsResp)
	if err != nil {
		return nil, err
	}

	return &locationsResp, nil
}

// GetLocationArea fetches details for a specific location area
func (c *Client) GetLocationArea(locationName string) (*LocationArea, error) {
	url := baseURL + "/location-area/" + locationName

	// Check if we have this URL in cache
	if cachedData, ok := c.cache.Get(url); ok {
		// Unmarshal cached data
		var locationArea LocationArea
		err := json.Unmarshal(cachedData, &locationArea)
		if err != nil {
			return nil, err
		}
		return &locationArea, nil
	}

	// Make the HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Add to cache
	c.cache.Add(url, body)

	// Unmarshal JSON into our struct
	var locationArea LocationArea
	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return nil, err
	}

	return &locationArea, nil
}

// GetPokemon fetches details for a specific Pokemon
func (c *Client) GetPokemon(pokemonName string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + pokemonName

	// Check if we have this URL in cache
	if cachedData, ok := c.cache.Get(url); ok {
		// Unmarshal cached data
		var pokemon Pokemon
		err := json.Unmarshal(cachedData, &pokemon)
		if err != nil {
			return Pokemon{}, err
		}
		return pokemon, nil
	}

	// Make the HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	// Add to cache
	c.cache.Add(url, body)

	// Unmarshal JSON into our struct
	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}
