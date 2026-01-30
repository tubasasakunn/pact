// Package pact provides the public API for the Pact diagram generation library.
package pact

import (
	"io"
	"strings"

	"pact/internal/application/transformer"
	"pact/internal/domain/ast"
	"pact/internal/domain/diagram/class"
	"pact/internal/domain/diagram/flow"
	"pact/internal/domain/diagram/sequence"
	"pact/internal/domain/diagram/state"
	"pact/internal/infrastructure/parser"
	"pact/internal/infrastructure/renderer/svg"
)

// Client provides the main API for parsing and generating diagrams.
type Client struct{}

// New creates a new Client instance.
func New() *Client {
	return &Client{}
}

// ParseFile parses a .pact file and returns the AST.
func (c *Client) ParseFile(path string) (*ast.SpecFile, error) {
	// Implementation: read file and parse
	return nil, nil
}

// ParseString parses a .pact string and returns the AST.
func (c *Client) ParseString(content string) (*ast.SpecFile, error) {
	lexer := parser.NewLexer(strings.NewReader(content))
	p := parser.NewParser(lexer)
	return p.Parse()
}

// ToClassDiagram transforms the AST to a class diagram.
func (c *Client) ToClassDiagram(spec *ast.SpecFile) (*class.Diagram, error) {
	tr := transformer.NewClassTransformer()
	return tr.Transform(spec)
}

// ToSequenceDiagram transforms the AST to a sequence diagram for the given flow.
func (c *Client) ToSequenceDiagram(spec *ast.SpecFile, flowName string) (*sequence.Diagram, error) {
	tr := transformer.NewSequenceTransformer()
	return tr.Transform(spec, flowName)
}

// ToStateDiagram transforms the AST to a state diagram for the given states.
func (c *Client) ToStateDiagram(spec *ast.SpecFile, statesName string) (*state.Diagram, error) {
	tr := transformer.NewStateTransformer()
	return tr.Transform(spec, statesName)
}

// ToFlowchart transforms the AST to a flowchart for the given flow.
func (c *Client) ToFlowchart(spec *ast.SpecFile, flowName string) (*flow.Diagram, error) {
	tr := transformer.NewFlowTransformer()
	return tr.Transform(spec, flowName)
}

// RenderClassDiagram renders a class diagram to SVG.
func (c *Client) RenderClassDiagram(diagram *class.Diagram, w io.Writer) error {
	renderer := svg.NewClassRenderer()
	return renderer.Render(diagram, w)
}

// RenderSequenceDiagram renders a sequence diagram to SVG.
func (c *Client) RenderSequenceDiagram(diagram *sequence.Diagram, w io.Writer) error {
	renderer := svg.NewSequenceRenderer()
	return renderer.Render(diagram, w)
}

// RenderStateDiagram renders a state diagram to SVG.
func (c *Client) RenderStateDiagram(diagram *state.Diagram, w io.Writer) error {
	renderer := svg.NewStateRenderer()
	return renderer.Render(diagram, w)
}

// RenderFlowchart renders a flowchart to SVG.
func (c *Client) RenderFlowchart(diagram *flow.Diagram, w io.Writer) error {
	renderer := svg.NewFlowRenderer()
	return renderer.Render(diagram, w)
}
