package transformer

import (
	"testing"

	"pact/internal/domain/ast"
	"pact/internal/domain/diagram/class"
)

// =============================================================================
// TC001-TC017: ClassTransformer Tests
// =============================================================================

// TC001: 空コンポーネント
func TestClassTransformer_EmptyComponent(t *testing.T) {
	files := []*ast.SpecFile{
		{
			Component: &ast.ComponentDecl{
				Name: "Foo",
				Body: ast.ComponentBody{},
			},
		},
	}

	transformer := NewClassTransformer()
	diagram, err := transformer.Transform(files, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(diagram.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(diagram.Nodes))
	}
	if len(diagram.Edges) != 0 {
		t.Errorf("expected 0 edges, got %d", len(diagram.Edges))
	}
}

// TC002: 型がノードになる
func TestClassTransformer_ComponentWithType(t *testing.T) {
	files := []*ast.SpecFile{
		{
			Component: &ast.ComponentDecl{
				Name: "Foo",
				Body: ast.ComponentBody{
					Types: []ast.TypeDecl{
						{Name: "User", Kind: ast.TypeKindStruct},
					},
				},
			},
		},
	}

	transformer := NewClassTransformer()
	diagram, err := transformer.Transform(files, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(diagram.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(diagram.Nodes))
	}
}

// TC003: メソッドがコンパートメントに
func TestClassTransformer_ComponentWithMethods(t *testing.T) {
	files := []*ast.SpecFile{
		{
			Component: &ast.ComponentDecl{
				Name: "Foo",
				Body: ast.ComponentBody{
					Provides: []ast.InterfaceDecl{
						{
							Name: "API",
							Methods: []ast.MethodDecl{
								{Name: "Get"},
								{Name: "Set"},
							},
						},
					},
				},
			},
		},
	}

	transformer := NewClassTransformer()
	diagram, err := transformer.Transform(files, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	componentNode := diagram.Nodes[0]
	if len(componentNode.Methods) != 2 {
		t.Errorf("expected 2 methods, got %d", len(componentNode.Methods))
	}
}

// TC004: フィールドが属性に
func TestClassTransformer_TypeFields(t *testing.T) {
	files := []*ast.SpecFile{
		{
			Component: &ast.ComponentDecl{
				Name: "Foo",
				Body: ast.ComponentBody{
					Types: []ast.TypeDecl{
						{
							Name: "User",
							Kind: ast.TypeKindStruct,
							Fields: []ast.FieldDecl{
								{Name: "id", Type: ast.TypeExpr{Name: "string"}},
								{Name: "name", Type: ast.TypeExpr{Name: "string"}},
							},
						},
					},
				},
			},
		},
	}

	transformer := NewClassTransformer()
	diagram, err := transformer.Transform(files, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// コンポーネント + 型 = 2ノード
	var typeNode *class.Node
	for i := range diagram.Nodes {
		if diagram.Nodes[i].Name == "User" {
			typeNode = &diagram.Nodes[i]
			break
		}
	}
	if typeNode == nil {
		t.Fatal("expected User node")
	}
	if len(typeNode.Attributes) != 2 {
		t.Errorf("expected 2 attributes, got %d", len(typeNode.Attributes))
	}
}

// TC005: 可視性マッピング
func TestClassTransformer_FieldVisibility(t *testing.T) {
	tests := []struct {
		input    ast.Visibility
		expected class.Visibility
	}{
		{ast.VisibilityPublic, class.VisibilityPublic},
		{ast.VisibilityPrivate, class.VisibilityPrivate},
		{ast.VisibilityProtected, class.VisibilityProtected},
		{ast.VisibilityPackage, class.VisibilityPackage},
	}

	for _, tt := range tests {
		files := []*ast.SpecFile{
			{
				Component: &ast.ComponentDecl{
					Name: "Foo",
					Body: ast.ComponentBody{
						Types: []ast.TypeDecl{
							{
								Name: "Bar",
								Kind: ast.TypeKindStruct,
								Fields: []ast.FieldDecl{
									{Name: "field", Type: ast.TypeExpr{Name: "string"}, Visibility: tt.input},
								},
							},
						},
					},
				},
			},
		}

		transformer := NewClassTransformer()
		diagram, _ := transformer.Transform(files, nil)

		var typeNode *class.Node
		for i := range diagram.Nodes {
			if diagram.Nodes[i].Name == "Bar" {
				typeNode = &diagram.Nodes[i]
				break
			}
		}

		if typeNode.Attributes[0].Visibility != tt.expected {
			t.Errorf("expected visibility %q, got %q", tt.expected, typeNode.Attributes[0].Visibility)
		}
	}
}

// TC006: enum
func TestClassTransformer_Enum(t *testing.T) {
	files := []*ast.SpecFile{
		{
			Component: &ast.ComponentDecl{
				Name: "Foo",
				Body: ast.ComponentBody{
					Types: []ast.TypeDecl{
						{Name: "Status", Kind: ast.TypeKindEnum, Values: []string{"A", "B"}},
					},
				},
			},
		},
	}

	transformer := NewClassTransformer()
	diagram, err := transformer.Transform(files, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var enumNode *class.Node
	for i := range diagram.Nodes {
		if diagram.Nodes[i].Name == "Status" {
			enumNode = &diagram.Nodes[i]
			break
		}
	}
	if enumNode == nil {
		t.Fatal("expected Status node")
	}
	if enumNode.Stereotype != "enum" {
		t.Errorf("expected stereotype 'enum', got %q", enumNode.Stereotype)
	}
}

// TC007: 依存エッジ
func TestClassTransformer_DependsOn(t *testing.T) {
	files := []*ast.SpecFile{
		{
			Component: &ast.ComponentDecl{
				Name: "Foo",
				Body: ast.ComponentBody{
					Relations: []ast.RelationDecl{
						{Kind: ast.RelationDependsOn, Target: "Bar"},
					},
				},
			},
		},
	}

	transformer := NewClassTransformer()
	diagram, err := transformer.Transform(files, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(diagram.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(diagram.Edges))
	}
	edge := diagram.Edges[0]
	if edge.Type != class.EdgeTypeDependency {
		t.Errorf("expected EdgeTypeDependency, got %v", edge.Type)
	}
}

// TC008: 継承エッジ
func TestClassTransformer_Extends(t *testing.T) {
	files := []*ast.SpecFile{
		{
			Component: &ast.ComponentDecl{
				Name: "Foo",
				Body: ast.ComponentBody{
					Relations: []ast.RelationDecl{
						{Kind: ast.RelationExtends, Target: "Base"},
					},
				},
			},
		},
	}

	transformer := NewClassTransformer()
	diagram, _ := transformer.Transform(files, nil)

	if diagram.Edges[0].Type != class.EdgeTypeInheritance {
		t.Errorf("expected EdgeTypeInheritance, got %v", diagram.Edges[0].Type)
	}
}

// TC009: 実装エッジ
func TestClassTransformer_Implements(t *testing.T) {
	files := []*ast.SpecFile{
		{
			Component: &ast.ComponentDecl{
				Name: "Foo",
				Body: ast.ComponentBody{
					Relations: []ast.RelationDecl{
						{Kind: ast.RelationImplements, Target: "Iface"},
					},
				},
			},
		},
	}

	transformer := NewClassTransformer()
	diagram, _ := transformer.Transform(files, nil)

	if diagram.Edges[0].Type != class.EdgeTypeImplementation {
		t.Errorf("expected EdgeTypeImplementation, got %v", diagram.Edges[0].Type)
	}
}

// TC010: コンポジションエッジ
func TestClassTransformer_Contains(t *testing.T) {
	files := []*ast.SpecFile{
		{
			Component: &ast.ComponentDecl{
				Name: "Foo",
				Body: ast.ComponentBody{
					Relations: []ast.RelationDecl{
						{Kind: ast.RelationContains, Target: "Cache"},
					},
				},
			},
		},
	}

	transformer := NewClassTransformer()
	diagram, _ := transformer.Transform(files, nil)

	if diagram.Edges[0].Type != class.EdgeTypeComposition {
		t.Errorf("expected EdgeTypeComposition, got %v", diagram.Edges[0].Type)
	}
}

// TC011: 集約エッジ
func TestClassTransformer_Aggregates(t *testing.T) {
	files := []*ast.SpecFile{
		{
			Component: &ast.ComponentDecl{
				Name: "Foo",
				Body: ast.ComponentBody{
					Relations: []ast.RelationDecl{
						{Kind: ast.RelationAggregates, Target: "Items"},
					},
				},
			},
		},
	}

	transformer := NewClassTransformer()
	diagram, _ := transformer.Transform(files, nil)

	if diagram.Edges[0].Type != class.EdgeTypeAggregation {
		t.Errorf("expected EdgeTypeAggregation, got %v", diagram.Edges[0].Type)
	}
}

// TC012: 要求インターフェース
func TestClassTransformer_RequiresInterface(t *testing.T) {
	files := []*ast.SpecFile{
		{
			Component: &ast.ComponentDecl{
				Name: "Foo",
				Body: ast.ComponentBody{
					Requires: []ast.InterfaceDecl{
						{Name: "Repository"},
					},
				},
			},
		},
	}

	transformer := NewClassTransformer()
	diagram, _ := transformer.Transform(files, nil)

	// コンポーネント + インターフェース
	if len(diagram.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(diagram.Nodes))
	}

	var ifaceNode *class.Node
	for i := range diagram.Nodes {
		if diagram.Nodes[i].Name == "Repository" {
			ifaceNode = &diagram.Nodes[i]
			break
		}
	}
	if ifaceNode == nil {
		t.Fatal("expected Repository node")
	}
	if ifaceNode.Stereotype != "interface" {
		t.Errorf("expected stereotype 'interface', got %q", ifaceNode.Stereotype)
	}
}

// TC013: 複数ファイル
func TestClassTransformer_MultipleFiles(t *testing.T) {
	files := []*ast.SpecFile{
		{
			Component: &ast.ComponentDecl{Name: "ServiceA", Body: ast.ComponentBody{}},
		},
		{
			Component: &ast.ComponentDecl{Name: "ServiceB", Body: ast.ComponentBody{}},
		},
	}

	transformer := NewClassTransformer()
	diagram, err := transformer.Transform(files, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(diagram.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(diagram.Nodes))
	}
}

// TC014: コンポーネントフィルタ
func TestClassTransformer_FilterComponents(t *testing.T) {
	files := []*ast.SpecFile{
		{Component: &ast.ComponentDecl{Name: "ServiceA", Body: ast.ComponentBody{}}},
		{Component: &ast.ComponentDecl{Name: "ServiceB", Body: ast.ComponentBody{}}},
		{Component: &ast.ComponentDecl{Name: "ServiceC", Body: ast.ComponentBody{}}},
	}

	transformer := NewClassTransformer()
	opts := &TransformOptions{FilterComponents: []string{"ServiceA", "ServiceC"}}
	diagram, _ := transformer.Transform(files, opts)

	if len(diagram.Nodes) != 2 {
		t.Errorf("expected 2 nodes (filtered), got %d", len(diagram.Nodes))
	}
}

// TC015: アノテーション伝播
func TestClassTransformer_Annotations(t *testing.T) {
	key := "file"
	files := []*ast.SpecFile{
		{
			Component: &ast.ComponentDecl{
				Name: "Foo",
				Annotations: []ast.AnnotationDecl{
					{Name: "source", Args: []ast.AnnotationArg{{Key: &key, Value: "foo.go"}}},
				},
				Body: ast.ComponentBody{},
			},
		},
	}

	transformer := NewClassTransformer()
	diagram, _ := transformer.Transform(files, nil)

	if len(diagram.Nodes[0].Annotations) != 1 {
		t.Errorf("expected 1 annotation, got %d", len(diagram.Nodes[0].Annotations))
	}
}

// TC016: エッジ装飾
func TestClassTransformer_EdgeDecorations(t *testing.T) {
	tests := []struct {
		kind       ast.RelationKind
		decoration class.Decoration
	}{
		{ast.RelationDependsOn, class.DecorationArrow},
		{ast.RelationExtends, class.DecorationTriangle},
		{ast.RelationImplements, class.DecorationTriangle},
		{ast.RelationContains, class.DecorationFilledDiamond},
		{ast.RelationAggregates, class.DecorationEmptyDiamond},
	}

	for _, tt := range tests {
		files := []*ast.SpecFile{
			{
				Component: &ast.ComponentDecl{
					Name: "Foo",
					Body: ast.ComponentBody{
						Relations: []ast.RelationDecl{{Kind: tt.kind, Target: "Bar"}},
					},
				},
			},
		}

		transformer := NewClassTransformer()
		diagram, _ := transformer.Transform(files, nil)

		if diagram.Edges[0].Decoration != tt.decoration {
			t.Errorf("for %v: expected decoration %v, got %v", tt.kind, tt.decoration, diagram.Edges[0].Decoration)
		}
	}
}

// TC017: 線スタイル
func TestClassTransformer_EdgeLineStyle(t *testing.T) {
	tests := []struct {
		kind      ast.RelationKind
		lineStyle class.LineStyle
	}{
		{ast.RelationDependsOn, class.LineStyleDashed},
		{ast.RelationExtends, class.LineStyleSolid},
		{ast.RelationImplements, class.LineStyleDashed},
		{ast.RelationContains, class.LineStyleSolid},
		{ast.RelationAggregates, class.LineStyleSolid},
	}

	for _, tt := range tests {
		files := []*ast.SpecFile{
			{
				Component: &ast.ComponentDecl{
					Name: "Foo",
					Body: ast.ComponentBody{
						Relations: []ast.RelationDecl{{Kind: tt.kind, Target: "Bar"}},
					},
				},
			},
		}

		transformer := NewClassTransformer()
		diagram, _ := transformer.Transform(files, nil)

		if diagram.Edges[0].LineStyle != tt.lineStyle {
			t.Errorf("for %v: expected line style %v, got %v", tt.kind, tt.lineStyle, diagram.Edges[0].LineStyle)
		}
	}
}
