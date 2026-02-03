package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pact/pkg/pact"
)

type generateOptions struct {
	output string
	types  []string
	files  []string
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
	files := expandFiles(opts.files)

	// Verify files exist
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			return fmt.Errorf("file not found: %s", f)
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
