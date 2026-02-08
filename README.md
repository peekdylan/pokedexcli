# Pokedex CLI 🎮

A feature-rich command-line Pokedex application built in Go that uses the [PokeAPI](https://pokeapi.co/) to explore Pokemon locations, catch Pokemon, and battle!

[![LinkedIn](https://img.shields.io/badge/LinkedIn-Connect-blue?style=flat&logo=linkedin)](https://www.linkedin.com/in/dylan-p-3b9297248)
![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## ✨ Features

- 🗺️ **Explore Pokemon Locations** - Navigate through the Pokemon world
- 🔍 **Search for Pokemon** - Find Pokemon in different areas
- ⚡ **Catch Pokemon** - Probability-based catching with progress bars
- 🎨 **ASCII Art** - Beautiful ASCII representations of Pokemon
- ⚔️ **Battle System** - Fight wild Pokemon with your party
- 👥 **Party System** - Build a team of up to 6 Pokemon
- 💾 **Persistent Storage** - Auto-save your progress
- 🎨 **Colored Output** - Beautiful terminal colors
- ⬆️ **Command History** - Use arrow keys to navigate previous commands
- 🧪 **Unit Tested** - Reliable and well-tested code

## 📦 Installation
```bash
# Clone the repository
git clone https://github.com/peekdylan/pokedexcli.git
cd pokedexcli

# Install dependencies
go mod download

# Build the application
go build

# Run it!
./pokedexcli
```

## 🎮 Usage

### Available Commands

| Command | Description |
|---------|-------------|
| `help` | Display available commands |
| `map` | Show next 20 location areas |
| `mapb` | Show previous 20 location areas |
| `explore <location>` | List Pokemon in a specific area |
| `catch <pokemon>` | Attempt to catch a Pokemon |
| `inspect <pokemon>` | View details of a caught Pokemon |
| `pokedex` | List all caught Pokemon |
| `party` | View your party of Pokemon |
| `addparty <pokemon>` | Add a Pokemon to your party |
| `battle <pokemon>` | Battle a wild Pokemon |
| `save` | Manually save your progress |
| `exit` | Save and exit the application |

### Example Session
```
Pokedex > map
📍 Location Areas:
  • canalave-city-area
  • eterna-city-area
  ...

Pokedex > explore canalave-city-area
🔍 Exploring canalave-city-area...

✨ Found Pokemon:
  • tentacool
  • tentacruel
  ...

Pokedex > catch pikachu
🎯 Throwing a Pokeball at pikachu...
[========================================] 100%
✓ pikachu was caught!

    ⢀⣠⣤⣤⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
  ⢀⣾⡿⠋⠀⠀⠈⠙⢷⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
  ⢸⡏⠀⠀⠀⠀⠀⠀⠀⢻⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
  ⢸⡇⠀⣀⣀⠀⠀⣀⣀⠀⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
  ⠘⣷⡀⠛⠛⠀⠀⠛⠛⢀⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
  ⠀⠘⢿⣦⣀⣀⣀⣀⣤⡾⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
  ⠀⠀⠀⠉⠛⠛⠛⠛⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀

   You may now inspect it with the inspect command.

Pokedex > addparty pikachu
✓ Added pikachu to your party!

Pokedex > battle charmander
╔════════════════════════════════════════╗
║          POKEMON BATTLE!             ║
╚════════════════════════════════════════╝

pikachu VS charmander
...
```

## 🏗️ Project Structure
```
pokedexcli/
├── main.go              # Main application and command handlers
├── battle.go            # Battle system logic
├── storage.go           # Save/load functionality
├── ascii.go             # ASCII art for Pokemon
├── repl.go              # Input cleaning utilities
├── repl_test.go         # Tests for input cleaning
├── pokedex_save.json    # Auto-generated save file
└── internal/
    ├── pokeapi/         # PokeAPI client
    │   └── pokeapi.go
    └── pokecache/       # Caching system
        ├── cache.go
        └── cache_test.go
```

## 🛠️ Technical Details

### Dependencies

- **[fatih/color](https://github.com/fatih/color)** - Colored terminal output
- **[schollz/progressbar](https://github.com/schollz/progressbar)** - Progress bars
- **[eiannone/keyboard](https://github.com/eiannone/keyboard)** - Keyboard input handling

### Architecture

- **Language**: Go 1.23+
- **API**: [PokeAPI v2](https://pokeapi.co/docs/v2)
- **Concurrency**: Thread-safe caching with mutex locks
- **Storage**: JSON-based persistence
- **Testing**: Unit tests with table-driven test patterns

## 🧪 Running Tests
```bash
go test ./...
```

## 🎯 Learning Outcomes

This project demonstrates:

- ✅ Building CLI applications with Go
- ✅ Making HTTP requests and parsing JSON
- ✅ Implementing caching strategies
- ✅ Thread-safe concurrent programming
- ✅ File I/O and data persistence
- ✅ Test-driven development
- ✅ Clean architecture patterns
- ✅ Terminal UI/UX design
- ✅ Game mechanics and probability systems

## 📝 License

MIT License - feel free to use this project for learning!

## 🙏 Acknowledgments

- Built as part of [Boot.dev](https://boot.dev)'s curriculum
- Pokemon data from [PokeAPI](https://pokeapi.co/)
- ASCII art inspired by Pokemon sprites

---

**Gotta Code 'Em All!** 🚀
