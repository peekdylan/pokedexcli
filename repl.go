package main

import (
	"strings"
)

func cleanInput(text string) []string {
	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(text)

	// Convert to lowercase
	lowered := strings.ToLower(trimmed)

	// Split on whitespace
	words := strings.Fields(lowered)

	return words
}
