package transformer

import (
	"testing"

	"pact/internal/domain/ast"
	"pact/internal/domain/diagram/flow"
	"pact/internal/domain/errors"
)

func createFlowTestComponent(flowSteps []ast.Step) *ast.SpecFile {
	return &ast.SpecFile{
		Component: &ast.ComponentDecl{
			Name: "TestService",
			Body: ast.ComponentBody{
				Flows: []ast.FlowDecl{
					{Name: "Process", Steps: flowSteps},
				},
			},
		},
	}
}

// TF001: 空フロー
func TestFlowTransformer_EmptyFlow(t *testing.T) {
	files := []*ast.SpecFile{createFlowTestComponent([]ast.Step{})}

	transformer := NewFlowTransformer()
	diagram, err := transformer.Transform(files, &FlowOptions{FlowName: "Process"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Start + End
	if len(diagram.Nodes) != 2 {
		t.Errorf("expected 2 nodes (start+end), got %d", len(diagram.Nodes))
	}
}

// TF002: 端子ノード
func TestFlowTransformer_StartEndNodes(t *testing.T) {
	files := []*ast.SpecFile{createFlowTestComponent([]ast.Step{})}

	transformer := NewFlowTransformer()
	diagram, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process"})

	hasStart := false
	hasEnd := false
	for _, node := range diagram.Nodes {
		if node.Label == "Start" && node.Shape == flow.NodeShapeTerminal {
			hasStart = true
		}
		if node.Label == "End" && node.Shape == flow.NodeShapeTerminal {
			hasEnd = true
		}
	}
	if !hasStart {
		t.Error("expected Start node")
	}
	if !hasEnd {
		t.Error("expected End node")
	}
}

// TF003: 代入→処理
func TestFlowTransformer_AssignStep(t *testing.T) {
	steps := []ast.Step{
		&ast.AssignStep{Variable: "x", Value: &ast.VariableExpr{Name: "y"}},
	}
	files := []*ast.SpecFile{createFlowTestComponent(steps)}

	transformer := NewFlowTransformer()
	diagram, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process"})

	hasProcess := false
	for _, node := range diagram.Nodes {
		if node.Shape == flow.NodeShapeProcess {
			hasProcess = true
			break
		}
	}
	if !hasProcess {
		t.Error("expected process node")
	}
}

// TF004: 呼び出し→処理
func TestFlowTransformer_CallStep(t *testing.T) {
	steps := []ast.Step{
		&ast.CallStep{
			Expr: &ast.CallExpr{
				Object: &ast.VariableExpr{Name: "Service"},
				Method: "DoSomething",
			},
		},
	}
	files := []*ast.SpecFile{createFlowTestComponent(steps)}

	transformer := NewFlowTransformer()
	diagram, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process"})

	hasProcess := false
	for _, node := range diagram.Nodes {
		if node.Shape == flow.NodeShapeProcess {
			hasProcess = true
			break
		}
	}
	if !hasProcess {
		t.Error("expected process node")
	}
}

// TF005: return
func TestFlowTransformer_ReturnStep(t *testing.T) {
	steps := []ast.Step{
		&ast.ReturnStep{Value: &ast.VariableExpr{Name: "x"}},
	}
	files := []*ast.SpecFile{createFlowTestComponent(steps)}

	transformer := NewFlowTransformer()
	diagram, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process"})

	returnCount := 0
	for _, node := range diagram.Nodes {
		if node.Label == "return" && node.Shape == flow.NodeShapeTerminal {
			returnCount++
		}
	}
	if returnCount != 1 {
		t.Errorf("expected 1 return node, got %d", returnCount)
	}
}

// TF006: throw
func TestFlowTransformer_ThrowStep(t *testing.T) {
	steps := []ast.Step{
		&ast.ThrowStep{Error: "NotFound"},
	}
	files := []*ast.SpecFile{createFlowTestComponent(steps)}

	transformer := NewFlowTransformer()
	diagram, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process"})

	hasThrow := false
	for _, node := range diagram.Nodes {
		if node.Shape == flow.NodeShapeTerminal && node.Label == "throw NotFound" {
			hasThrow = true
			break
		}
	}
	if !hasThrow {
		t.Error("expected throw node")
	}
}

// TF007: 判断
func TestFlowTransformer_IfStep(t *testing.T) {
	steps := []ast.Step{
		&ast.IfStep{
			Condition: &ast.VariableExpr{Name: "cond"},
			Then:      []ast.Step{},
		},
	}
	files := []*ast.SpecFile{createFlowTestComponent(steps)}

	transformer := NewFlowTransformer()
	diagram, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process"})

	hasDecision := false
	for _, node := range diagram.Nodes {
		if node.Shape == flow.NodeShapeDecision {
			hasDecision = true
			break
		}
	}
	if !hasDecision {
		t.Error("expected decision node")
	}
}

// TF008: 分岐ラベル
func TestFlowTransformer_IfStep_Labels(t *testing.T) {
	steps := []ast.Step{
		&ast.IfStep{
			Condition: &ast.VariableExpr{Name: "cond"},
			Then: []ast.Step{
				&ast.CallStep{Expr: &ast.CallExpr{Object: &ast.VariableExpr{Name: "A"}, Method: "Do"}},
			},
			Else: []ast.Step{
				&ast.CallStep{Expr: &ast.CallExpr{Object: &ast.VariableExpr{Name: "B"}, Method: "Do"}},
			},
		},
	}
	files := []*ast.SpecFile{createFlowTestComponent(steps)}

	transformer := NewFlowTransformer()
	diagram, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process"})

	hasYes := false
	hasNo := false
	for _, edge := range diagram.Edges {
		if edge.Label == "Yes" {
			hasYes = true
		}
		if edge.Label == "No" {
			hasNo = true
		}
	}
	if !hasYes {
		t.Error("expected Yes label on edge")
	}
	if !hasNo {
		t.Error("expected No label on edge")
	}
}

// TF009: if-else
func TestFlowTransformer_IfElse(t *testing.T) {
	steps := []ast.Step{
		&ast.IfStep{
			Condition: &ast.VariableExpr{Name: "cond"},
			Then: []ast.Step{
				&ast.CallStep{Expr: &ast.CallExpr{Object: &ast.VariableExpr{Name: "A"}, Method: "Do"}},
			},
			Else: []ast.Step{
				&ast.CallStep{Expr: &ast.CallExpr{Object: &ast.VariableExpr{Name: "B"}, Method: "Do"}},
			},
		},
	}
	files := []*ast.SpecFile{createFlowTestComponent(steps)}

	transformer := NewFlowTransformer()
	diagram, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process"})

	// decision, 2 process nodes, merge node + start + end = 6
	if len(diagram.Nodes) < 4 {
		t.Errorf("expected at least 4 nodes, got %d", len(diagram.Nodes))
	}
}

// TF010: forループ
func TestFlowTransformer_ForStep(t *testing.T) {
	steps := []ast.Step{
		&ast.ForStep{
			Variable: "item",
			Iterable: &ast.VariableExpr{Name: "items"},
			Body:     []ast.Step{},
		},
	}
	files := []*ast.SpecFile{createFlowTestComponent(steps)}

	transformer := NewFlowTransformer()
	diagram, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process"})

	hasDecision := false
	for _, node := range diagram.Nodes {
		if node.Shape == flow.NodeShapeDecision {
			hasDecision = true
			break
		}
	}
	if !hasDecision {
		t.Error("expected decision node for for loop")
	}
}

// TF011: whileループ
func TestFlowTransformer_WhileStep(t *testing.T) {
	steps := []ast.Step{
		&ast.WhileStep{
			Condition: &ast.VariableExpr{Name: "cond"},
			Body:      []ast.Step{},
		},
	}
	files := []*ast.SpecFile{createFlowTestComponent(steps)}

	transformer := NewFlowTransformer()
	diagram, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process"})

	hasDecision := false
	for _, node := range diagram.Nodes {
		if node.Shape == flow.NodeShapeDecision {
			hasDecision = true
			break
		}
	}
	if !hasDecision {
		t.Error("expected decision node for while loop")
	}
}

// TF013: 順次
func TestFlowTransformer_Sequential(t *testing.T) {
	steps := []ast.Step{
		&ast.CallStep{Expr: &ast.CallExpr{Object: &ast.VariableExpr{Name: "A"}, Method: "Step1"}},
		&ast.CallStep{Expr: &ast.CallExpr{Object: &ast.VariableExpr{Name: "A"}, Method: "Step2"}},
		&ast.CallStep{Expr: &ast.CallExpr{Object: &ast.VariableExpr{Name: "A"}, Method: "Step3"}},
	}
	files := []*ast.SpecFile{createFlowTestComponent(steps)}

	transformer := NewFlowTransformer()
	diagram, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process"})

	// Start + 3 process + End = 5 nodes
	if len(diagram.Nodes) != 5 {
		t.Errorf("expected 5 nodes, got %d", len(diagram.Nodes))
	}

	// 4 edges (start->1, 1->2, 2->3, 3->end)
	if len(diagram.Edges) != 4 {
		t.Errorf("expected 4 edges, got %d", len(diagram.Edges))
	}
}

// TF014: スイムレーン推論
func TestFlowTransformer_Swimlane_Infer(t *testing.T) {
	steps := []ast.Step{
		&ast.AssignStep{
			Variable: "x",
			Value: &ast.CallExpr{
				Object: &ast.VariableExpr{Name: "ServiceA"},
				Method: "Get",
			},
		},
		&ast.AssignStep{
			Variable: "y",
			Value: &ast.CallExpr{
				Object: &ast.VariableExpr{Name: "ServiceB"},
				Method: "Process",
			},
		},
	}
	files := []*ast.SpecFile{createFlowTestComponent(steps)}

	transformer := NewFlowTransformer()
	diagram, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process", IncludeSwimlanes: true})

	if len(diagram.Swimlanes) != 2 {
		t.Errorf("expected 2 swimlanes, got %d", len(diagram.Swimlanes))
	}
}

// TF015: スイムレーンオプション
func TestFlowTransformer_OptionSwimlanes(t *testing.T) {
	steps := []ast.Step{
		&ast.AssignStep{
			Variable: "x",
			Value: &ast.CallExpr{
				Object: &ast.VariableExpr{Name: "ServiceA"},
				Method: "Get",
			},
		},
	}
	files := []*ast.SpecFile{createFlowTestComponent(steps)}

	transformer := NewFlowTransformer()

	// スイムレーンなし
	diagram1, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process", IncludeSwimlanes: false})
	if len(diagram1.Swimlanes) != 0 {
		t.Errorf("expected 0 swimlanes without option, got %d", len(diagram1.Swimlanes))
	}

	// スイムレーンあり
	diagram2, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process", IncludeSwimlanes: true})
	if len(diagram2.Swimlanes) != 1 {
		t.Errorf("expected 1 swimlane with option, got %d", len(diagram2.Swimlanes))
	}
}

// TF016: 未発見
func TestFlowTransformer_FlowNotFound(t *testing.T) {
	files := []*ast.SpecFile{createFlowTestComponent([]ast.Step{})}

	transformer := NewFlowTransformer()
	_, err := transformer.Transform(files, &FlowOptions{FlowName: "NonExistent"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if _, ok := err.(*errors.TransformError); !ok {
		t.Errorf("expected TransformError, got %T", err)
	}
}

// TF017: 形状マッピング
func TestFlowTransformer_NodeShapes(t *testing.T) {
	// AssignStep -> Process
	// CallStep -> Process
	// ReturnStep -> Terminal
	// ThrowStep -> Terminal
	// IfStep -> Decision
	// ForStep -> Decision
	// WhileStep -> Decision

	steps := []ast.Step{
		&ast.AssignStep{Variable: "x", Value: &ast.VariableExpr{Name: "y"}},
	}
	files := []*ast.SpecFile{createFlowTestComponent(steps)}

	transformer := NewFlowTransformer()
	diagram, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process"})

	processCount := 0
	for _, node := range diagram.Nodes {
		if node.Shape == flow.NodeShapeProcess {
			processCount++
		}
	}
	if processCount != 1 {
		t.Errorf("expected 1 process node, got %d", processCount)
	}
}

// TF019: マージノード
func TestFlowTransformer_MergeNode(t *testing.T) {
	steps := []ast.Step{
		&ast.IfStep{
			Condition: &ast.VariableExpr{Name: "cond"},
			Then: []ast.Step{
				&ast.CallStep{Expr: &ast.CallExpr{Object: &ast.VariableExpr{Name: "A"}, Method: "Do"}},
			},
			Else: []ast.Step{
				&ast.CallStep{Expr: &ast.CallExpr{Object: &ast.VariableExpr{Name: "B"}, Method: "Do"}},
			},
		},
	}
	files := []*ast.SpecFile{createFlowTestComponent(steps)}

	transformer := NewFlowTransformer()
	diagram, _ := transformer.Transform(files, &FlowOptions{FlowName: "Process"})

	hasConnector := false
	for _, node := range diagram.Nodes {
		if node.Shape == flow.NodeShapeConnector {
			hasConnector = true
			break
		}
	}
	if !hasConnector {
		t.Error("expected connector (merge) node")
	}
}
