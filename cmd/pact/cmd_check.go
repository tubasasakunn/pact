package main

import (
	"fmt"
	"strings"

	"pact/pkg/pact"
)

func cmdCheck(args []string) error {
	showMissing := false
	var patterns []string

	for _, arg := range args {
		if arg == "--missing" || arg == "-m" {
			showMissing = true
		} else if !strings.HasPrefix(arg, "-") {
			patterns = append(patterns, arg)
		}
	}

	// Default to current directory
	if len(patterns) == 0 {
		patterns = []string{"."}
	}

	// Find all pact files
	pactFiles := expandFiles(patterns)

	if len(pactFiles) == 0 {
		if showMissing {
			fmt.Println("No .pact files found")
			return fmt.Errorf("no files found")
		}
		fmt.Println("No .pact files found")
		return nil
	}

	client := pact.New()
	hasErrors := false
	components := make(map[string]bool)
	dependencies := make(map[string][]string)

	// Parse all files and collect components/dependencies
	for _, file := range pactFiles {
		spec, err := client.ParseFile(file)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", file, err)
			hasErrors = true
			continue
		}

		// Collect components
		if spec.Component != nil {
			components[spec.Component.Name] = true
			for _, rel := range spec.Component.Body.Relations {
				if rel.Kind == "depends_on" {
					dependencies[spec.Component.Name] = append(dependencies[spec.Component.Name], rel.Target)
				}
			}
		}
		for _, comp := range spec.Components {
			components[comp.Name] = true
			for _, rel := range comp.Body.Relations {
				if rel.Kind == "depends_on" {
					dependencies[comp.Name] = append(dependencies[comp.Name], rel.Target)
				}
			}
		}
	}

	// Check for missing dependencies
	if showMissing {
		missing := make(map[string]bool)
		for _, deps := range dependencies {
			for _, dep := range deps {
				if !components[dep] {
					missing[dep] = true
				}
			}
		}

		if len(missing) > 0 {
			fmt.Println("Missing components:")
			for name := range missing {
				fmt.Printf("  - %s\n", name)
			}
			return fmt.Errorf("missing components found")
		}
		fmt.Println("No missing components")
	}

	if hasErrors {
		return fmt.Errorf("check failed")
	}

	fmt.Printf("Checked %d files, found %d components\n", len(pactFiles), len(components))
	return nil
}
