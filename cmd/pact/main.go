package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pact/pkg/pact"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "init":
		if err := cmdInit(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "generate":
		if err := cmdGenerate(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "validate":
		if err := cmdValidate(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "check":
		if err := cmdCheck(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "watch":
		if err := cmdWatch(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "version", "-v", "--version":
		fmt.Printf("pact version %s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`pact - Pact DSL diagram generator

Usage:
  pact <command> [options] [arguments]

Commands:
  init        Initialize a new .pactconfig file
  generate    Generate diagrams from .pact files
  validate    Validate .pact files
  check       Check for missing components
  watch       Watch for file changes and regenerate
  version     Show version information
  help        Show this help message

Examples:
  pact init
  pact generate service.pact
  pact generate -o output/ -t class service.pact
  pact validate *.pact
  pact check --missing`)
}

// =============================================================================
// init command
// =============================================================================

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

// =============================================================================
// generate command
// =============================================================================

type generateOptions struct {
	output     string
	types      []string
	files      []string
}

func parseGenerateOptions(args []string) (*generateOptions, error) {
	opts := &generateOptions{
		output: ".",
		types:  []string{"all"},
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-o" || arg == "--output":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("missing value for %s", arg)
			}
			i++
			opts.output = args[i]
		case arg == "-t" || arg == "--type":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("missing value for %s", arg)
			}
			i++
			opts.types = strings.Split(args[i], ",")
		case strings.HasPrefix(arg, "-"):
			return nil, fmt.Errorf("unknown option: %s", arg)
		default:
			opts.files = append(opts.files, arg)
		}
	}

	if len(opts.files) == 0 {
		return nil, fmt.Errorf("no input files specified")
	}

	return opts, nil
}

func cmdGenerate(args []string) error {
	opts, err := parseGenerateOptions(args)
	if err != nil {
		return err
	}

	// Create output directory if needed
	if opts.output != "." {
		if err := os.MkdirAll(opts.output, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Expand file patterns
	var files []string
	for _, pattern := range opts.files {
		// Check if it's a directory
		info, err := os.Stat(pattern)
		if err == nil && info.IsDir() {
			matches, _ := filepath.Glob(filepath.Join(pattern, "*.pact"))
			files = append(files, matches...)
		} else {
			matches, _ := filepath.Glob(pattern)
			if len(matches) == 0 {
				// Check if file exists directly
				if _, err := os.Stat(pattern); err != nil {
					return fmt.Errorf("file not found: %s", pattern)
				}
				files = append(files, pattern)
			} else {
				files = append(files, matches...)
			}
		}
	}

	if len(files) == 0 {
		return fmt.Errorf("no .pact files found")
	}

	client := pact.New()

	for _, file := range files {
		fmt.Printf("Processing %s...\n", file)

		spec, err := client.ParseFile(file)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", file, err)
		}

		baseName := strings.TrimSuffix(filepath.Base(file), ".pact")

		// Generate class diagram
		if shouldGenerate(opts.types, "class") {
			if err := generateClassDiagram(client, spec, opts.output, baseName); err != nil {
				fmt.Printf("  Warning: class diagram: %v\n", err)
			} else {
				fmt.Printf("  Generated %s_class.svg\n", baseName)
			}
		}

		// Generate sequence diagrams
		if shouldGenerate(opts.types, "sequence") {
			if err := generateSequenceDiagrams(client, spec, opts.output, baseName); err != nil {
				fmt.Printf("  Warning: sequence diagram: %v\n", err)
			}
		}

		// Generate state diagrams
		if shouldGenerate(opts.types, "state") {
			if err := generateStateDiagrams(client, spec, opts.output, baseName); err != nil {
				fmt.Printf("  Warning: state diagram: %v\n", err)
			}
		}

		// Generate flowcharts
		if shouldGenerate(opts.types, "flow") {
			if err := generateFlowcharts(client, spec, opts.output, baseName); err != nil {
				fmt.Printf("  Warning: flowchart: %v\n", err)
			}
		}
	}

	fmt.Println("Done!")
	return nil
}

func shouldGenerate(types []string, target string) bool {
	for _, t := range types {
		if t == "all" || t == target {
			return true
		}
	}
	return false
}

func generateClassDiagram(client *pact.Client, spec *pact.SpecFile, output, baseName string) error {
	diagram, err := client.ToClassDiagram(spec)
	if err != nil {
		return err
	}

	outPath := filepath.Join(output, baseName+"_class.svg")
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return client.RenderClassDiagram(diagram, f)
}

func generateSequenceDiagrams(client *pact.Client, spec *pact.SpecFile, output, baseName string) error {
	// Find all flows
	flows := getFlowNames(spec)
	if len(flows) == 0 {
		return nil
	}

	for _, flowName := range flows {
		diagram, err := client.ToSequenceDiagram(spec, flowName)
		if err != nil {
			continue
		}

		outPath := filepath.Join(output, baseName+"_sequence_"+flowName+".svg")
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}

		if err := client.RenderSequenceDiagram(diagram, f); err != nil {
			f.Close()
			return err
		}
		f.Close()
		fmt.Printf("  Generated %s_sequence_%s.svg\n", baseName, flowName)
	}
	return nil
}

func generateStateDiagrams(client *pact.Client, spec *pact.SpecFile, output, baseName string) error {
	// Find all states
	stateNames := getStateNames(spec)
	if len(stateNames) == 0 {
		return nil
	}

	for _, stateName := range stateNames {
		diagram, err := client.ToStateDiagram(spec, stateName)
		if err != nil {
			continue
		}

		outPath := filepath.Join(output, baseName+"_state_"+stateName+".svg")
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}

		if err := client.RenderStateDiagram(diagram, f); err != nil {
			f.Close()
			return err
		}
		f.Close()
		fmt.Printf("  Generated %s_state_%s.svg\n", baseName, stateName)
	}
	return nil
}

func generateFlowcharts(client *pact.Client, spec *pact.SpecFile, output, baseName string) error {
	// Find all flows
	flows := getFlowNames(spec)
	if len(flows) == 0 {
		return nil
	}

	for _, flowName := range flows {
		diagram, err := client.ToFlowchart(spec, flowName)
		if err != nil {
			continue
		}

		outPath := filepath.Join(output, baseName+"_flow_"+flowName+".svg")
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}

		if err := client.RenderFlowchart(diagram, f); err != nil {
			f.Close()
			return err
		}
		f.Close()
		fmt.Printf("  Generated %s_flow_%s.svg\n", baseName, flowName)
	}
	return nil
}

func getFlowNames(spec *pact.SpecFile) []string {
	var names []string

	if spec.Component != nil {
		for _, flow := range spec.Component.Body.Flows {
			names = append(names, flow.Name)
		}
	}

	for _, comp := range spec.Components {
		for _, flow := range comp.Body.Flows {
			names = append(names, flow.Name)
		}
	}

	return names
}

func getStateNames(spec *pact.SpecFile) []string {
	var names []string

	if spec.Component != nil {
		for _, states := range spec.Component.Body.States {
			names = append(names, states.Name)
		}
	}

	for _, comp := range spec.Components {
		for _, states := range comp.Body.States {
			names = append(names, states.Name)
		}
	}

	return names
}

// =============================================================================
// validate command
// =============================================================================

func cmdValidate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no files specified")
	}

	// Expand file patterns
	var files []string
	for _, pattern := range args {
		info, err := os.Stat(pattern)
		if err == nil && info.IsDir() {
			matches, _ := filepath.Glob(filepath.Join(pattern, "*.pact"))
			files = append(files, matches...)
		} else {
			matches, _ := filepath.Glob(pattern)
			if len(matches) == 0 {
				files = append(files, pattern)
			} else {
				files = append(files, matches...)
			}
		}
	}

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

// =============================================================================
// check command
// =============================================================================

func cmdCheck(args []string) error {
	showMissing := false
	var files []string

	for _, arg := range args {
		if arg == "--missing" || arg == "-m" {
			showMissing = true
		} else if !strings.HasPrefix(arg, "-") {
			files = append(files, arg)
		}
	}

	// Default to current directory
	if len(files) == 0 {
		files = []string{"."}
	}

	// Find all pact files
	var pactFiles []string
	for _, pattern := range files {
		info, err := os.Stat(pattern)
		if err == nil && info.IsDir() {
			matches, _ := filepath.Glob(filepath.Join(pattern, "*.pact"))
			pactFiles = append(pactFiles, matches...)
		} else {
			matches, _ := filepath.Glob(pattern)
			pactFiles = append(pactFiles, matches...)
		}
	}

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

// =============================================================================
// watch command
// =============================================================================

func cmdWatch(args []string) error {
	fmt.Println("Watch mode is not implemented yet")
	fmt.Println("Use a file watcher like entr or fswatch with 'pact generate'")
	return nil
}
