// Pact Public API - Main entry point for the Pact library
// Source: pkg/pact/api.go
@version("1.0")
component Client {
	// Type aliases for public use (re-exported from ast package)
	type SpecFile {
		path: string
		imports: Import[]
		singleComponent: ComponentDecl?
		components: ComponentDecl[]
	}

	type ComponentDecl {
		name: string
		body: ComponentBody
	}

	type ComponentBody {
		types: TypeDecl[]
		relations: RelationDecl[]
		interfaces: InterfaceDecl[]
		statesBlocks: StatesDecl[]
		flows: FlowDecl[]
	}

	type RelationDecl {
		kind: string
		target: string
	}

	type FlowDecl {
		name: string
		statements: Statement[]
	}

	type StatesDecl {
		name: string
		statesBlocks: StateDecl[]
		transitions: TransitionDecl[]
	}

	depends on Parser
	depends on Lexer
	depends on ClassTransformer
	depends on SequenceTransformer
	depends on StateTransformer
	depends on FlowTransformer
	depends on ClassRenderer
	depends on SequenceRenderer
	depends on StateRenderer
	depends on FlowRenderer

	// Public API for parsing and diagram generation
	provides PactAPI {
		// Parse a .pact file from disk
		ParseFile(filePath: string) -> SpecFile

		// Parse a .pact string directly
		ParseString(content: string) -> SpecFile

		// Transform AST to class diagram
		ToClassDiagram(spec: SpecFile) -> ClassDiagram

		// Transform AST to sequence diagram for a specific flow
		ToSequenceDiagram(spec: SpecFile, flowName: string) -> SequenceDiagram

		// Transform AST to state diagram for specific states
		ToStateDiagram(spec: SpecFile, statesName: string) -> StateDiagram

		// Transform AST to flowchart for a specific flow
		ToFlowchart(spec: SpecFile, flowName: string) -> FlowDiagram

		// Render class diagram to SVG
		RenderClassDiagram(diagram: ClassDiagram, writer: Writer) -> error

		// Render sequence diagram to SVG
		RenderSequenceDiagram(diagram: SequenceDiagram, writer: Writer) -> error

		// Render state diagram to SVG
		RenderStateDiagram(diagram: StateDiagram, writer: Writer) -> error

		// Render flowchart to SVG
		RenderFlowchart(diagram: FlowDiagram, writer: Writer) -> error
	}

	// Parse file flow
	flow ParseFileFlow {
		content = self.readFile(filePath)
		if content == null {
			throw FileNotFoundError
		}
		spec = self.ParseString(content)
		return spec
	}

	// Parse string flow
	flow ParseStringFlow {
		lexer = Lexer.New(content)
		parser = Parser.New(lexer)
		spec = parser.Parse()
		if spec.hasErrors {
			throw ParseError
		}
		return spec
	}

	// Generate class diagram flow
	flow GenerateClassDiagram {
		transformer = ClassTransformer.New()
		diagram = transformer.Transform(spec)
		if diagram == null {
			throw TransformError
		}
		renderer = ClassRenderer.New()
		result = renderer.Render(diagram, writer)
		return result
	}
}
