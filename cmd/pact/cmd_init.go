package main

import (
	"fmt"
	"os"
)

func cmdInit(args []string) error {
	configPath := ".pactconfig"

	// Check if already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Println("Config file already exists, overwriting...")
	}

	config := `# Pact configuration file
src: .
output: diagrams
types:
  - class
  - sequence
  - state
  - flow
`

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	fmt.Println("Created .pactconfig")
	return nil
}
