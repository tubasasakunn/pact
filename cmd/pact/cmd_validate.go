package main

import (
	"fmt"

	"pact/pkg/pact"
)

func cmdValidate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no files specified")
	}

	// Expand file patterns
	files := expandFiles(args)

	client := pact.New()
	hasErrors := false

	for _, file := range files {
		_, err := client.ParseFile(file)
		if err != nil {
			fmt.Printf("Error in %s: %v\n", file, err)
			hasErrors = true
		} else {
			fmt.Printf("%s: OK\n", file)
		}
	}

	if hasErrors {
		return fmt.Errorf("validation failed")
	}

	fmt.Println("All files valid!")
	return nil
}
