package svg

import (
	"bytes"
	"strings"
	"testing"

	"pact/internal/domain/diagram/class"
)

// =============================================================================
// RCL001-RCL013: ClassRenderer Tests
// =============================================================================

// RCL001: 空図
func TestClassRenderer_EmptyDiagram(t *testing.T) {
	diagram := &class.Diagram{}

	renderer := NewClassRenderer()
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

// RCL002: 単一ノード
func TestClassRenderer_SingleNode(t *testing.T) {
	diagram := &class.Diagram{
		Nodes: []class.Node{
			{ID: "User", Name: "User"},
		},
	}

	renderer := NewClassRenderer()
	var buf bytes.Buffer
	err := renderer.Render(diagram, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	if !strings.Contains(svg, "<rect") {
		t.Error("expected rect element for node")
	}
	if !strings.Contains(svg, "User") {
		t.Error("expected node name in output")
	}
}

// RCL003: 属性付きノード
func TestClassRenderer_NodeWithAttributes(t *testing.T) {
	diagram := &class.Diagram{
		Nodes: []class.Node{
			{
				ID:   "User",
				Name: "User",
				Attributes: []class.Attribute{
					{Name: "id", Type: "string", Visibility: class.VisibilityPrivate},
					{Name: "name", Type: "string", Visibility: class.VisibilityPublic},
				},
			},
		},
	}

	renderer := NewClassRenderer()
	var buf bytes.Buffer
	err := renderer.Render(diagram, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	// コンパートメント区切りの線
	lineCount := strings.Count(svg, "<line")
	if lineCount < 1 {
		t.Error("expected line for compartment separator")
	}
}

// RCL004: メソッド付きノード
func TestClassRenderer_NodeWithMethods(t *testing.T) {
	diagram := &class.Diagram{
		Nodes: []class.Node{
			{
				ID:   "Service",
				Name: "Service",
				Methods: []class.Method{
					{Name: "Get", Visibility: class.VisibilityPublic},
					{Name: "Set", Visibility: class.VisibilityPublic},
				},
			},
		},
	}

	renderer := NewClassRenderer()
	var buf bytes.Buffer
	err := renderer.Render(diagram, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	if !strings.Contains(svg, "Get") {
		t.Error("expected method name in output")
	}
}

// RCL005: ステレオタイプ
func TestClassRenderer_Stereotype(t *testing.T) {
	diagram := &class.Diagram{
		Nodes: []class.Node{
			{
				ID:         "API",
				Name:       "API",
				Stereotype: "interface",
			},
		},
	}

	renderer := NewClassRenderer()
	var buf bytes.Buffer
	renderer.Render(diagram, &buf)

	svg := buf.String()
	if !strings.Contains(svg, "<<interface>>") {
		t.Error("expected stereotype in output")
	}
}

// RCL006: 可視性記号
func TestClassRenderer_VisibilitySymbols(t *testing.T) {
	diagram := &class.Diagram{
		Nodes: []class.Node{
			{
				ID:   "Test",
				Name: "Test",
				Attributes: []class.Attribute{
					{Name: "pub", Type: "string", Visibility: class.VisibilityPublic},
					{Name: "priv", Type: "string", Visibility: class.VisibilityPrivate},
					{Name: "prot", Type: "string", Visibility: class.VisibilityProtected},
					{Name: "pkg", Type: "string", Visibility: class.VisibilityPackage},
				},
			},
		},
	}

	renderer := NewClassRenderer()
	var buf bytes.Buffer
	renderer.Render(diagram, &buf)

	svg := buf.String()
	if !strings.Contains(svg, "+ pub") {
		t.Error("expected + for public")
	}
	if !strings.Contains(svg, "- priv") {
		t.Error("expected - for private")
	}
	if !strings.Contains(svg, "# prot") {
		t.Error("expected # for protected")
	}
	if !strings.Contains(svg, "~ pkg") {
		t.Error("expected ~ for package")
	}
}

// RCL007: 依存エッジ
func TestClassRenderer_Edge_Dependency(t *testing.T) {
	diagram := &class.Diagram{
		Nodes: []class.Node{
			{ID: "A", Name: "A"},
			{ID: "B", Name: "B"},
		},
		Edges: []class.Edge{
			{From: "A", To: "B", Type: class.EdgeTypeDependency, LineStyle: class.LineStyleDashed},
		},
	}

	renderer := NewClassRenderer()
	var buf bytes.Buffer
	err := renderer.Render(diagram, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	if !strings.Contains(svg, "<line") {
		t.Error("expected line for edge")
	}
}

// RCL008: 継承エッジ
func TestClassRenderer_Edge_Inheritance(t *testing.T) {
	diagram := &class.Diagram{
		Nodes: []class.Node{
			{ID: "Child", Name: "Child"},
			{ID: "Parent", Name: "Parent"},
		},
		Edges: []class.Edge{
			{From: "Child", To: "Parent", Type: class.EdgeTypeInheritance, Decoration: class.DecorationTriangle},
		},
	}

	renderer := NewClassRenderer()
	var buf bytes.Buffer
	renderer.Render(diagram, &buf)

	svg := buf.String()
	// 三角形の装飾
	if !strings.Contains(svg, "<polygon") {
		t.Error("expected polygon for triangle decoration")
	}
}

// RCL010: コンポジションエッジ
func TestClassRenderer_Edge_Composition(t *testing.T) {
	diagram := &class.Diagram{
		Nodes: []class.Node{
			{ID: "Whole", Name: "Whole"},
			{ID: "Part", Name: "Part"},
		},
		Edges: []class.Edge{
			{From: "Whole", To: "Part", Type: class.EdgeTypeComposition, Decoration: class.DecorationFilledDiamond},
		},
	}

	renderer := NewClassRenderer()
	var buf bytes.Buffer
	renderer.Render(diagram, &buf)

	svg := buf.String()
	// 黒ひし形の装飾
	polygonCount := strings.Count(svg, "<polygon")
	if polygonCount < 1 {
		t.Error("expected polygon for diamond decoration")
	}
}

// RCL011: 集約エッジ
func TestClassRenderer_Edge_Aggregation(t *testing.T) {
	diagram := &class.Diagram{
		Nodes: []class.Node{
			{ID: "Container", Name: "Container"},
			{ID: "Item", Name: "Item"},
		},
		Edges: []class.Edge{
			{From: "Container", To: "Item", Type: class.EdgeTypeAggregation, Decoration: class.DecorationEmptyDiamond},
		},
	}

	renderer := NewClassRenderer()
	var buf bytes.Buffer
	renderer.Render(diagram, &buf)

	svg := buf.String()
	// 白ひし形の装飾
	if !strings.Contains(svg, "<polygon") {
		t.Error("expected polygon for empty diamond decoration")
	}
}

// RCL012: レイアウト
func TestClassRenderer_Layout(t *testing.T) {
	diagram := &class.Diagram{
		Nodes: []class.Node{
			{ID: "A", Name: "A"},
			{ID: "B", Name: "B"},
			{ID: "C", Name: "C"},
		},
	}

	renderer := NewClassRenderer()
	var buf bytes.Buffer
	renderer.Render(diagram, &buf)

	svg := buf.String()
	// 3つのノードがある
	rectCount := strings.Count(svg, "<rect")
	if rectCount < 3 {
		t.Errorf("expected at least 3 rect elements, got %d", rectCount)
	}
}
