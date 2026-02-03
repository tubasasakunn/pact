package service

import (
	"bytes"
	"io"
	"testing"

	"pact/internal/application/transformer"
	"pact/internal/domain/ast"
	"pact/internal/domain/diagram/class"
	"pact/internal/domain/diagram/flow"
	"pact/internal/domain/diagram/sequence"
	"pact/internal/domain/diagram/state"
)

// Mock renderers for testing

type mockClassRenderer struct {
	called bool
}

func (m *mockClassRenderer) Render(d *class.Diagram, w io.Writer) error {
	m.called = true
	_, err := w.Write([]byte("<svg>class</svg>"))
	return err
}

type mockSequenceRenderer struct {
	called bool
}

func (m *mockSequenceRenderer) Render(d *sequence.Diagram, w io.Writer) error {
	m.called = true
	_, err := w.Write([]byte("<svg>sequence</svg>"))
	return err
}

type mockStateRenderer struct {
	called bool
}

func (m *mockStateRenderer) Render(d *state.Diagram, w io.Writer) error {
	m.called = true
	_, err := w.Write([]byte("<svg>state</svg>"))
	return err
}

type mockFlowRenderer struct {
	called bool
}

func (m *mockFlowRenderer) Render(d *flow.Diagram, w io.Writer) error {
	m.called = true
	_, err := w.Write([]byte("<svg>flow</svg>"))
	return err
}

func TestNewDiagramService(t *testing.T) {
	t.Parallel()
	svc := NewDiagramService(
		&mockClassRenderer{},
		&mockSequenceRenderer{},
		&mockStateRenderer{},
		&mockFlowRenderer{},
	)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestDiagramService_TransformClassDiagram(t *testing.T) {
	t.Parallel()
	svc := NewDiagramService(
		&mockClassRenderer{},
		&mockSequenceRenderer{},
		&mockStateRenderer{},
		&mockFlowRenderer{},
	)

	spec := &ast.SpecFile{
		Component: &ast.ComponentDecl{
			Name: "TestService",
			Body: ast.ComponentBody{},
		},
	}

	diagram, err := svc.TransformClassDiagram([]*ast.SpecFile{spec}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diagram == nil {
		t.Fatal("expected non-nil diagram")
	}
}

func TestDiagramService_TransformSequenceDiagram_NoFlowName(t *testing.T) {
	t.Parallel()
	svc := NewDiagramService(
		&mockClassRenderer{},
		&mockSequenceRenderer{},
		&mockStateRenderer{},
		&mockFlowRenderer{},
	)

	spec := &ast.SpecFile{}
	_, err := svc.TransformSequenceDiagram([]*ast.SpecFile{spec}, nil)
	if err == nil {
		t.Fatal("expected error for nil options")
	}
}

func TestDiagramService_GenerateClassDiagram(t *testing.T) {
	t.Parallel()
	cr := &mockClassRenderer{}
	svc := NewDiagramService(
		cr,
		&mockSequenceRenderer{},
		&mockStateRenderer{},
		&mockFlowRenderer{},
	)

	spec := &ast.SpecFile{
		Component: &ast.ComponentDecl{
			Name: "TestService",
			Body: ast.ComponentBody{},
		},
	}

	var buf bytes.Buffer
	err := svc.GenerateClassDiagram([]*ast.SpecFile{spec}, nil, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cr.called {
		t.Error("expected class renderer to be called")
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestDiagramService_GenerateSequenceDiagram(t *testing.T) {
	t.Parallel()
	sr := &mockSequenceRenderer{}
	svc := NewDiagramService(
		&mockClassRenderer{},
		sr,
		&mockStateRenderer{},
		&mockFlowRenderer{},
	)

	spec := &ast.SpecFile{
		Component: &ast.ComponentDecl{
			Name: "TestService",
			Body: ast.ComponentBody{
				Flows: []ast.FlowDecl{
					{
						Name:  "TestFlow",
						Steps: []ast.Step{},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := svc.GenerateSequenceDiagram(
		[]*ast.SpecFile{spec},
		&transformer.SequenceOptions{FlowName: "TestFlow"},
		&buf,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sr.called {
		t.Error("expected sequence renderer to be called")
	}
}

func TestDiagramService_GenerateStateDiagram(t *testing.T) {
	t.Parallel()
	str := &mockStateRenderer{}
	svc := NewDiagramService(
		&mockClassRenderer{},
		&mockSequenceRenderer{},
		str,
		&mockFlowRenderer{},
	)

	spec := &ast.SpecFile{
		Component: &ast.ComponentDecl{
			Name: "TestService",
			Body: ast.ComponentBody{
				States: []ast.StatesDecl{
					{
						Name:   "TestStates",
						States: []ast.StateDecl{},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := svc.GenerateStateDiagram(
		[]*ast.SpecFile{spec},
		&transformer.StateOptions{StatesName: "TestStates"},
		&buf,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !str.called {
		t.Error("expected state renderer to be called")
	}
}

func TestDiagramService_GenerateFlowchart(t *testing.T) {
	t.Parallel()
	fr := &mockFlowRenderer{}
	svc := NewDiagramService(
		&mockClassRenderer{},
		&mockSequenceRenderer{},
		&mockStateRenderer{},
		fr,
	)

	spec := &ast.SpecFile{
		Component: &ast.ComponentDecl{
			Name: "TestService",
			Body: ast.ComponentBody{
				Flows: []ast.FlowDecl{
					{
						Name:  "TestFlow",
						Steps: []ast.Step{},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := svc.GenerateFlowchart(
		[]*ast.SpecFile{spec},
		&transformer.FlowOptions{FlowName: "TestFlow"},
		&buf,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fr.called {
		t.Error("expected flow renderer to be called")
	}
}
