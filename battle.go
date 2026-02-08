package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/peekdylan/pokedexcli/internal/pokeapi"
	"github.com/schollz/progressbar/v3"
)

func startBattle(cfg *config, playerPokemon, wildPokemon pokeapi.Pokemon) error {
	fmt.Println(cyan("\n╔════════════════════════════════════════╗"))
	fmt.Println(cyan("║") + bold("          POKEMON BATTLE!             ") + cyan("║"))
	fmt.Println(cyan("╚════════════════════════════════════════╝\n"))

	fmt.Printf("%s VS %s\n\n",
		green(bold(playerPokemon.Name)),
		red(bold(wildPokemon.Name)))

	playerHP := getPokemonStat(playerPokemon, "hp")
	wildHP := getPokemonStat(wildPokemon, "hp")

	playerMaxHP := playerHP
	wildMaxHP := wildHP

	round := 1

	for playerHP > 0 && wildHP > 0 {
		fmt.Printf(yellow("--- Round %d ---\n"), round)

		// Show HP bars
		fmt.Printf("%s HP: ", green(playerPokemon.Name))
		showHPBar(playerHP, playerMaxHP)

		fmt.Printf("%s HP: ", red(wildPokemon.Name))
		showHPBar(wildHP, wildMaxHP)
		fmt.Println()

		// Player attacks
		playerAttack := getPokemonStat(playerPokemon, "attack")
		wildDefense := getPokemonStat(wildPokemon, "defense")
		damage := calculateDamage(playerAttack, wildDefense)

		fmt.Printf("%s attacks for %s damage!\n",
			green(playerPokemon.Name),
			bold(fmt.Sprintf("%d", damage)))

		time.Sleep(1 * time.Second)
		wildHP -= damage

		if wildHP <= 0 {
			fmt.Printf(green("\n🎉 %s won the battle!\n\n"), bold(playerPokemon.Name))

			// Chance to catch the wild pokemon
			fmt.Println(cyan("💫 The wild pokemon is weakened!"))
			fmt.Println(yellow("Would you like to try catching it? (y/n)"))

			var response string
			fmt.Scanln(&response)

			if response == "y" || response == "yes" {
				return commandCatch(cfg, []string{wildPokemon.Name})
			}
			return nil
		}

		// Wild pokemon attacks
		wildAttack := getPokemonStat(wildPokemon, "attack")
		playerDefense := getPokemonStat(playerPokemon, "defense")
		damage = calculateDamage(wildAttack, playerDefense)

		fmt.Printf("%s attacks for %s damage!\n",
			red(wildPokemon.Name),
			bold(fmt.Sprintf("%d", damage)))

		time.Sleep(1 * time.Second)
		playerHP -= damage

		if playerHP <= 0 {
			fmt.Printf(red("\n💔 %s fainted! You lost the battle.\n\n"), playerPokemon.Name)
			return nil
		}

		fmt.Println()
		round++
		time.Sleep(1 * time.Second)
	}

	return nil
}

func calculateDamage(attack, defense int) int {
	baseDamage := attack - (defense / 2)
	if baseDamage < 1 {
		baseDamage = 1
	}
	// Add some randomness
	variance := rand.Intn(5) + 1
	return baseDamage + variance
}

func showHPBar(currentHP, maxHP int) {
	barWidth := 20
	filledWidth := (currentHP * barWidth) / maxHP
	if filledWidth < 0 {
		filledWidth = 0
	}

	bar := progressbar.NewOptions(maxHP,
		progressbar.OptionSetWidth(barWidth),
		progressbar.OptionShowCount(),
		progressbar.OptionSetDescription(fmt.Sprintf("%3d/%3d", currentHP, maxHP)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "█",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	bar.Set(currentHP)
	fmt.Println()
}
