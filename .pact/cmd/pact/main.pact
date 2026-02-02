// Pact CLI - Command Line Interface for diagram generation
// Source: cmd/pact/main.go
@version("1.0")
component PactCLI {
	// generateOptions holds options for the generate command
	type GenerateOptions {
		outputDir: string
		diagramTypes: string[]
		inputFiles: string[]
	}

	// Command handlers
	provides CLI {
		// Main entry point
		Main(args: string[])

		// Initialize a new .pactconfig file
		Init(args: string[]) -> error

		// Generate diagrams from .pact files
		Generate(args: string[]) -> error

		// Validate .pact files for syntax errors
		Validate(args: string[]) -> error

		// Check for missing component dependencies
		Check(args: string[]) -> error

		// Watch for file changes and regenerate
		Watch(args: string[]) -> error

		// Print usage information
		PrintUsage()

		// Parse command line options for generate command
		ParseGenerateOptions(args: string[]) -> GenerateOptions

		// Check if a diagram type should be generated
		ShouldGenerate(diagramTypes: string[], target: string) -> bool

		// Generate class diagram from spec
		GenerateClassDiagram(client: Client, spec: SpecFile, outputDir: string, baseName: string) -> error

		// Generate sequence diagrams from spec
		GenerateSequenceDiagrams(client: Client, spec: SpecFile, outputDir: string, baseName: string) -> error

		// Generate state diagrams from spec
		GenerateStateDiagrams(client: Client, spec: SpecFile, outputDir: string, baseName: string) -> error

		// Generate flowcharts from spec
		GenerateFlowcharts(client: Client, spec: SpecFile, outputDir: string, baseName: string) -> error

		// Get flow names from spec
		GetFlowNames(spec: SpecFile) -> string[]

		// Get state names from spec
		GetStateNames(spec: SpecFile) -> string[]
	}

	depends on Client

	// CLI command flow
	flow ProcessCommand {
		cmd = self.parseCommand(args)
		result = self.dispatchCommand(cmd, args)
		return result
	}

	// Diagram generation flow
	flow GenerateDiagrams {
		opts = self.ParseGenerateOptions(args)
		self.validateOptions(opts)
		self.processInputFiles(opts)
		return success
	}
}
