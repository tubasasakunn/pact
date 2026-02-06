package svg

import (
	"bytes"
	"strings"
	"testing"

	"pact/internal/domain/diagram/flow"
)

// =============================================================================
// RFL001-RFL013: FlowRenderer Tests
// =============================================================================

// RFL001: 空図
func TestFlowRenderer_EmptyDiagram(t *testing.T) {
	diagram := &flow.Diagram{}

	renderer := NewFlowRenderer()
	var buf bytes.Buffer
	err := renderer.Render(diagram, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	if !strings.Contains(svg, "<svg") {
		t.Error("expected valid SVG output")
	}
}

// RFL002: 端子ノード
func TestFlowRenderer_TerminalNode(t *testing.T) {
	diagram := &flow.Diagram{
		Nodes: []flow.Node{
			{ID: "start", Label: "Start", Shape: flow.NodeShapeTerminal},
		},
	}

	renderer := NewFlowRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(diagram, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	// 角丸長方形
	if !strings.Contains(svg, "<rect") {
		t.Error("expected rect for terminal node")
	}
}

// RFL003: 処理ノード
func TestFlowRenderer_ProcessNode(t *testing.T) {
	diagram := &flow.Diagram{
		Nodes: []flow.Node{
			{ID: "process", Label: "Process", Shape: flow.NodeShapeProcess},
		},
	}

	renderer := NewFlowRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(diagram, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	if !strings.Contains(svg, "<rect") {
		t.Error("expected rect for process node")
	}
}

// RFL004: 判断ノード
func TestFlowRenderer_DecisionNode(t *testing.T) {
	diagram := &flow.Diagram{
		Nodes: []flow.Node{
			{ID: "decision", Label: "condition?", Shape: flow.NodeShapeDecision},
		},
	}

	renderer := NewFlowRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(diagram, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	// ひし形
	if !strings.Contains(svg, "<polygon") {
		t.Error("expected polygon for decision node")
	}
}

// RFL005: 入出力ノード
func TestFlowRenderer_IONode(t *testing.T) {
	diagram := &flow.Diagram{
		Nodes: []flow.Node{
			{ID: "io", Label: "Input", Shape: flow.NodeShapeIO},
		},
	}

	renderer := NewFlowRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(diagram, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	// 平行四辺形
	if !strings.Contains(svg, "<polygon") {
		t.Error("expected polygon for IO node")
	}
}

// RFL006: データベースノード
func TestFlowRenderer_DatabaseNode(t *testing.T) {
	diagram := &flow.Diagram{
		Nodes: []flow.Node{
			{ID: "db", Label: "Database", Shape: flow.NodeShapeDatabase},
		},
	}

	renderer := NewFlowRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(diagram, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	// 円柱
	if !strings.Contains(svg, "<ellipse") {
		t.Error("expected ellipse for database node")
	}
}

// RFL007: ノードラベル
func TestFlowRenderer_NodeLabel(t *testing.T) {
	diagram := &flow.Diagram{
		Nodes: []flow.Node{
			{ID: "process", Label: "Process Data", Shape: flow.NodeShapeProcess},
		},
	}

	renderer := NewFlowRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(diagram, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	if !strings.Contains(svg, "Process Data") {
		t.Error("expected label in output")
	}
}
