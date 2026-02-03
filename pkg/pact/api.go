// Package pact provides the public API for the Pact diagram generation library.
package pact

import (
	"io"
	"os"

	"pact/internal/application/service"
	"pact/internal/application/transformer"
	"pact/internal/domain/ast"
	"pact/internal/domain/diagram/class"
	"pact/internal/domain/diagram/flow"
	"pact/internal/domain/diagram/sequence"
	"pact/internal/domain/diagram/state"
	"pact/internal/infrastructure/parser"
	"pact/internal/infrastructure/renderer/svg"
)

// Type aliases for public use
type (
	SpecFile      = ast.SpecFile
	ComponentDecl = ast.ComponentDecl
	ComponentBody = ast.ComponentBody
	RelationDecl  = ast.RelationDecl
	FlowDecl      = ast.FlowDecl
	StatesDecl    = ast.StatesDecl
)

// Client provides the main API for parsing and generating diagrams.
type Client struct {
	service          *service.DiagramService
	classRenderer    *svg.ClassRenderer
	sequenceRenderer *svg.SequenceRenderer
	stateRenderer    *svg.StateRenderer
	flowRenderer     *svg.FlowRenderer
}

// New creates a new Client instance with default SVG renderers.
func New() *Client {
	cr := svg.NewClassRenderer()
	sr := svg.NewSequenceRenderer()
	str := svg.NewStateRenderer()
	fr := svg.NewFlowRenderer()

	svc := service.NewDiagramService(cr, sr, str, fr)

	return &Client{
		service:          svc,
		classRenderer:    cr,
		sequenceRenderer: sr,
		stateRenderer:    str,
		flowRenderer:     fr,
	}
}

// ParseFile parses a .pact file and returns the AST.
func (c *Client) ParseFile(path string) (*ast.SpecFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return c.ParseString(string(content))
}

// ParseString parses a .pact string and returns the AST.
func (c *Client) ParseString(content string) (*ast.SpecFile, error) {
	lexer := parser.NewLexer(content)
	p := parser.NewParser(lexer)
	return p.Parse()
}

// ToClassDiagram transforms the AST to a class diagram.
func (c *Client) ToClassDiagram(spec *ast.SpecFile) (*class.Diagram, error) {
	return c.service.TransformClassDiagram([]*ast.SpecFile{spec}, nil)
}

// ToSequenceDiagram transforms the AST to a sequence diagram for the given flow.
func (c *Client) ToSequenceDiagram(spec *ast.SpecFile, flowName string) (*sequence.Diagram, error) {
	return c.service.TransformSequenceDiagram([]*ast.SpecFile{spec}, &transformer.SequenceOptions{FlowName: flowName})
}

// ToStateDiagram transforms the AST to a state diagram for the given states.
func (c *Client) ToStateDiagram(spec *ast.SpecFile, statesName string) (*state.Diagram, error) {
	return c.service.TransformStateDiagram([]*ast.SpecFile{spec}, &transformer.StateOptions{StatesName: statesName})
}

// ToFlowchart transforms the AST to a flowchart for the given flow.
func (c *Client) ToFlowchart(spec *ast.SpecFile, flowName string) (*flow.Diagram, error) {
	return c.service.TransformFlowchart([]*ast.SpecFile{spec}, &transformer.FlowOptions{FlowName: flowName})
}

// RenderClassDiagram renders a class diagram to SVG.
func (c *Client) RenderClassDiagram(diagram *class.Diagram, w io.Writer) error {
	return c.classRenderer.Render(diagram, w)
}

// RenderSequenceDiagram renders a sequence diagram to SVG.
func (c *Client) RenderSequenceDiagram(diagram *sequence.Diagram, w io.Writer) error {
	return c.sequenceRenderer.Render(diagram, w)
}

// RenderStateDiagram renders a state diagram to SVG.
func (c *Client) RenderStateDiagram(diagram *state.Diagram, w io.Writer) error {
	return c.stateRenderer.Render(diagram, w)
}

// RenderFlowchart renders a flowchart to SVG.
func (c *Client) RenderFlowchart(diagram *flow.Diagram, w io.Writer) error {
	return c.flowRenderer.Render(diagram, w)
}
