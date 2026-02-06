// Command pattern-preview generates SVG preview files for all pattern templates.
// Usage: go run ./cmd/pattern-preview
// Output: pattern-preview/ directory with SVG files
package main

import (
	"fmt"
	"os"

	"pact/pkg/pact"
)

func main() {
	outDir := "pattern-preview"

	if err := pact.GeneratePatternPreviews(pact.PatternPreviewConfig{
		OutputDir: outDir,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate pattern previews: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated pattern previews in %s/\n", outDir)
	fmt.Printf("Open %s/index.html in a browser to view all patterns\n", outDir)
}
