// Package service provides application-level orchestration services.
package service

import (
	"io"

	"pact/internal/application/transformer"
	"pact/internal/domain/ast"
	"pact/internal/domain/diagram/class"
	"pact/internal/domain/diagram/flow"
	"pact/internal/domain/diagram/sequence"
	"pact/internal/domain/diagram/state"
)

// Parser parses .pact source code into an AST.
type Parser interface {
	Parse() (*ast.SpecFile, error)
}

// ClassRenderer renders class diagrams.
type ClassRenderer interface {
	Render(d *class.Diagram, w io.Writer) error
}

// SequenceRenderer renders sequence diagrams.
type SequenceRenderer interface {
	Render(d *sequence.Diagram, w io.Writer) error
}

// StateRenderer renders state diagrams.
type StateRenderer interface {
	Render(d *state.Diagram, w io.Writer) error
}

// FlowRenderer renders flow diagrams.
type FlowRenderer interface {
	Render(d *flow.Diagram, w io.Writer) error
}

// DiagramService orchestrates the parse → transform → render pipeline.
type DiagramService struct {
	classRenderer    ClassRenderer
	sequenceRenderer SequenceRenderer
	stateRenderer    StateRenderer
	flowRenderer     FlowRenderer
}

// NewDiagramService creates a new DiagramService with the given renderers.
func NewDiagramService(
	classRenderer ClassRenderer,
	sequenceRenderer SequenceRenderer,
	stateRenderer StateRenderer,
	flowRenderer FlowRenderer,
) *DiagramService {
	return &DiagramService{
		classRenderer:    classRenderer,
		sequenceRenderer: sequenceRenderer,
		stateRenderer:    stateRenderer,
		flowRenderer:     flowRenderer,
	}
}

// GenerateClassDiagram transforms AST files into a class diagram and renders it.
func (s *DiagramService) GenerateClassDiagram(files []*ast.SpecFile, opts *transformer.TransformOptions, w io.Writer) error {
	tr := transformer.NewClassTransformer()
	diagram, err := tr.Transform(files, opts)
	if err != nil {
		return err
	}
	return s.classRenderer.Render(diagram, w)
}

// GenerateSequenceDiagram transforms AST files into a sequence diagram and renders it.
func (s *DiagramService) GenerateSequenceDiagram(files []*ast.SpecFile, opts *transformer.SequenceOptions, w io.Writer) error {
	tr := transformer.NewSequenceTransformer()
	diagram, err := tr.Transform(files, opts)
	if err != nil {
		return err
	}
	return s.sequenceRenderer.Render(diagram, w)
}

// GenerateStateDiagram transforms AST files into a state diagram and renders it.
func (s *DiagramService) GenerateStateDiagram(files []*ast.SpecFile, opts *transformer.StateOptions, w io.Writer) error {
	tr := transformer.NewStateTransformer()
	diagram, err := tr.Transform(files, opts)
	if err != nil {
		return err
	}
	return s.stateRenderer.Render(diagram, w)
}

// GenerateFlowchart transforms AST files into a flowchart and renders it.
func (s *DiagramService) GenerateFlowchart(files []*ast.SpecFile, opts *transformer.FlowOptions, w io.Writer) error {
	tr := transformer.NewFlowTransformer()
	diagram, err := tr.Transform(files, opts)
	if err != nil {
		return err
	}
	return s.flowRenderer.Render(diagram, w)
}

// TransformClassDiagram transforms AST files into a class diagram model.
func (s *DiagramService) TransformClassDiagram(files []*ast.SpecFile, opts *transformer.TransformOptions) (*class.Diagram, error) {
	tr := transformer.NewClassTransformer()
	return tr.Transform(files, opts)
}

// TransformSequenceDiagram transforms AST files into a sequence diagram model.
func (s *DiagramService) TransformSequenceDiagram(files []*ast.SpecFile, opts *transformer.SequenceOptions) (*sequence.Diagram, error) {
	tr := transformer.NewSequenceTransformer()
	return tr.Transform(files, opts)
}

// TransformStateDiagram transforms AST files into a state diagram model.
func (s *DiagramService) TransformStateDiagram(files []*ast.SpecFile, opts *transformer.StateOptions) (*state.Diagram, error) {
	tr := transformer.NewStateTransformer()
	return tr.Transform(files, opts)
}

// TransformFlowchart transforms AST files into a flowchart model.
func (s *DiagramService) TransformFlowchart(files []*ast.SpecFile, opts *transformer.FlowOptions) (*flow.Diagram, error) {
	tr := transformer.NewFlowTransformer()
	return tr.Transform(files, opts)
}
