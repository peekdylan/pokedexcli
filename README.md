# Pokedex CLI

A command-line Pokedex application built in Go that uses the [PokeAPI](https://pokeapi.co/) to explore Pokemon locations and catch Pokemon.

## Features

- 🗺️ Explore Pokemon world locations
- 🔍 Search for Pokemon in different areas
- ⚡ Catch Pokemon with probability-based mechanics
- 📊 View detailed stats for caught Pokemon
- 💾 Intelligent caching system for fast performance
- 🧪 Unit tested

## Installation
```bash
git clone https://github.com/YOUR_USERNAME/pokedexcli.git
cd pokedexcli
go build
./pokedexcli
```

## Usage

### Commands

- `help` - Display available commands
- `map` - Show next 20 location areas
- `mapb` - Show previous 20 location areas
- `explore <location>` - List Pokemon in a specific area
- `catch <pokemon>` - Attempt to catch a Pokemon
- `inspect <pokemon>` - View details of a caught Pokemon
- `pokedex` - List all caught Pokemon
- `exit` - Exit the application

### Example Session
```
Pokedex > map
canalave-city-area
eterna-city-area
...

Pokedex > explore canalave-city-area
Exploring canalave-city-area...
Found Pokemon:
 - tentacool
 - tentacruel
 ...

Pokedex > catch tentacool
Throwing a Pokeball at tentacool...
tentacool was caught!

Pokedex > inspect tentacool
Name: tentacool
Height: 9
Weight: 455
Stats:
  -hp: 40
  -attack: 40
  ...
Types:
  - water
  - poison

Pokedex > pokedex
Your Pokedex:
 - tentacool
```

## Project Structure
```
pokedexcli/
├── main.go              # Main application and command handlers
├── repl.go              # Input cleaning utilities
├── repl_test.go         # Tests for input cleaning
└── internal/
    ├── pokeapi/         # PokeAPI client
    │   └── pokeapi.go
    └── pokecache/       # Caching system
        ├── cache.go
        └── cache_test.go
```

## Technical Details

- **Language**: Go 1.23+
- **API**: [PokeAPI v2](https://pokeapi.co/docs/v2)
- **Concurrency**: Thread-safe caching with mutex locks
- **Testing**: Unit tests with table-driven test patterns

## Learning Outcomes

This project was built as part of [Boot.dev](https://boot.dev)'s curriculum and demonstrates:

- Building CLI applications with Go
- Making HTTP requests and parsing JSON
- Implementing caching strategies
- Thread-safe concurrent programming
- Test-driven development
- Clean architecture patterns

## License

MIT
