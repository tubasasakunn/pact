package transformer

import (
	"fmt"

	"pact/internal/domain/ast"
	"pact/internal/domain/diagram/flow"
	"pact/internal/domain/errors"
)

// FlowTransformer はASTをフローチャートに変換する
type FlowTransformer struct {
	nodeCounter   int
	pendingNoEdge bool   // 次のエッジに"No"ラベルを付けるかどうか
	noEdgeFromID  string // "No"エッジの起点ノードID
}

// NewFlowTransformer は新しいFlowTransformerを作成する
func NewFlowTransformer() *FlowTransformer {
	return &FlowTransformer{}
}

// FlowOptions はフローチャート変換オプション
type FlowOptions struct {
	FlowName         string
	IncludeSwimlanes bool
}

// Transform はASTをフローチャートに変換する
func (t *FlowTransformer) Transform(files []*ast.SpecFile, opts *FlowOptions) (*flow.Diagram, error) {
	if opts == nil || opts.FlowName == "" {
		return nil, &errors.TransformError{Source: "AST", Target: "Flowchart", Message: "flow name is required"}
	}

	var targetFlow *ast.FlowDecl

	for _, file := range files {
		// 単一コンポーネント
		if file.Component != nil {
			for i := range file.Component.Body.Flows {
				if file.Component.Body.Flows[i].Name == opts.FlowName {
					targetFlow = &file.Component.Body.Flows[i]
					break
				}
			}
		}

		// 複数コンポーネント
		if targetFlow == nil {
			for j := range file.Components {
				comp := &file.Components[j]
				for i := range comp.Body.Flows {
					if comp.Body.Flows[i].Name == opts.FlowName {
						targetFlow = &comp.Body.Flows[i]
						break
					}
				}
				if targetFlow != nil {
					break
				}
			}
		}
	}

	if targetFlow == nil {
		return nil, &errors.TransformError{
			Source:  "AST",
			Target:  "Flowchart",
			Message: "flow not found: " + opts.FlowName,
		}
	}

	t.nodeCounter = 0
	t.pendingNoEdge = false
	t.noEdgeFromID = ""
	diagram := &flow.Diagram{
		Nodes:     []flow.Node{},
		Edges:     []flow.Edge{},
		Swimlanes: []flow.Swimlane{},
	}

	// 開始ノード
	startNode := t.createNode("Start", flow.NodeShapeTerminal)
	diagram.Nodes = append(diagram.Nodes, startNode)

	// ステップを変換
	lastNodeID := startNode.ID
	for _, step := range targetFlow.Steps {
		lastNodeID = t.transformStep(step, lastNodeID, diagram)
	}

	// 終了ノード
	endNode := t.createNode("End", flow.NodeShapeTerminal)
	diagram.Nodes = append(diagram.Nodes, endNode)
	if lastNodeID != "" {
		diagram.Edges = append(diagram.Edges, flow.Edge{From: lastNodeID, To: endNode.ID})
	}

	// スイムレーンを収集
	if opts.IncludeSwimlanes {
		swimlaneMap := make(map[string]bool)
		for _, node := range diagram.Nodes {
			if node.Swimlane != "" && !swimlaneMap[node.Swimlane] {
				diagram.Swimlanes = append(diagram.Swimlanes, flow.Swimlane{
					ID:   node.Swimlane,
					Name: node.Swimlane,
				})
				swimlaneMap[node.Swimlane] = true
			}
		}
	}

	return diagram, nil
}

func (t *FlowTransformer) createNode(label string, shape flow.NodeShape) flow.Node {
	t.nodeCounter++
	return flow.Node{
		ID:    fmt.Sprintf("node_%d", t.nodeCounter),
		Label: label,
		Shape: shape,
	}
}

func (t *FlowTransformer) transformStep(step ast.Step, prevNodeID string, diagram *flow.Diagram) string {
	// "No"エッジが保留中の場合、エッジにラベルを付ける
	addEdge := func(from, to string) {
		edge := flow.Edge{From: from, To: to}
		if t.pendingNoEdge && from == t.noEdgeFromID {
			edge.Label = "No"
			t.pendingNoEdge = false
			t.noEdgeFromID = ""
		}
		diagram.Edges = append(diagram.Edges, edge)
	}

	switch s := step.(type) {
	case *ast.AssignStep:
		label := s.Variable + " = " + t.formatExpr(s.Value)
		node := t.createNode(label, flow.NodeShapeProcess)
		if call, ok := s.Value.(*ast.CallExpr); ok {
			if v, ok := call.Object.(*ast.VariableExpr); ok {
				node.Swimlane = v.Name
			}
		}
		diagram.Nodes = append(diagram.Nodes, node)
		if prevNodeID != "" {
			addEdge(prevNodeID, node.ID)
		}
		return node.ID

	case *ast.CallStep:
		label := "Call"
		if call, ok := s.Expr.(*ast.CallExpr); ok {
			label = call.Method + "()"
			node := t.createNode(label, flow.NodeShapeProcess)
			if v, ok := call.Object.(*ast.VariableExpr); ok {
				node.Swimlane = v.Name
			}
			diagram.Nodes = append(diagram.Nodes, node)
			if prevNodeID != "" {
				addEdge(prevNodeID, node.ID)
			}
			return node.ID
		}
		return prevNodeID

	case *ast.ReturnStep:
		node := t.createNode("return", flow.NodeShapeTerminal)
		diagram.Nodes = append(diagram.Nodes, node)
		if prevNodeID != "" {
			addEdge(prevNodeID, node.ID)
		}
		return node.ID // returnノードからEndへ接続する

	case *ast.ThrowStep:
		node := t.createNode("throw "+s.Error, flow.NodeShapeTerminal)
		diagram.Nodes = append(diagram.Nodes, node)
		if prevNodeID != "" {
			diagram.Edges = append(diagram.Edges, flow.Edge{From: prevNodeID, To: node.ID})
		}
		return "" // throwの後は接続しない

	case *ast.IfStep:
		decisionNode := t.createNode("condition?", flow.NodeShapeDecision)
		diagram.Nodes = append(diagram.Nodes, decisionNode)
		if prevNodeID != "" {
			addEdge(prevNodeID, decisionNode.ID)
		}

		// Then分岐
		var thenEndID string
		thenEndsWithTerminal := false
		if len(s.Then) > 0 {
			currentID := decisionNode.ID
			for i, thenStep := range s.Then {
				if i == 0 {
					// 最初のエッジにYesラベル
					newID := t.transformStepWithLabel(thenStep, currentID, diagram, "Yes")
					currentID = newID
				} else {
					currentID = t.transformStep(thenStep, currentID, diagram)
				}
			}
			thenEndID = currentID
			// thenがreturn/throwで終わる場合
			if thenEndID == "" {
				thenEndsWithTerminal = true
			}
		} else {
			thenEndID = decisionNode.ID
		}

		// Else分岐
		var elseEndID string
		if len(s.Else) > 0 {
			currentID := decisionNode.ID
			for i, elseStep := range s.Else {
				if i == 0 {
					newID := t.transformStepWithLabel(elseStep, currentID, diagram, "No")
					currentID = newID
				} else {
					currentID = t.transformStep(elseStep, currentID, diagram)
				}
			}
			elseEndID = currentID
		} else {
			elseEndID = ""
		}

		// マージノード
		if thenEndID != "" || elseEndID != "" {
			mergeNode := t.createNode("", flow.NodeShapeConnector)
			diagram.Nodes = append(diagram.Nodes, mergeNode)
			if thenEndID != "" {
				diagram.Edges = append(diagram.Edges, flow.Edge{From: thenEndID, To: mergeNode.ID})
			}
			if elseEndID != "" {
				diagram.Edges = append(diagram.Edges, flow.Edge{From: elseEndID, To: mergeNode.ID})
			} else if len(s.Else) == 0 {
				// elseがない場合、decisionから直接マージへ
				diagram.Edges = append(diagram.Edges, flow.Edge{From: decisionNode.ID, To: mergeNode.ID, Label: "No"})
			}
			return mergeNode.ID
		}

		// then が終端で終わり、else がない場合、次のステップに "No" ラベルを付ける
		if thenEndsWithTerminal && len(s.Else) == 0 {
			t.pendingNoEdge = true
			t.noEdgeFromID = decisionNode.ID
		}
		return decisionNode.ID

	case *ast.ForStep:
		decisionNode := t.createNode("for "+s.Variable, flow.NodeShapeDecision)
		diagram.Nodes = append(diagram.Nodes, decisionNode)
		if prevNodeID != "" {
			diagram.Edges = append(diagram.Edges, flow.Edge{From: prevNodeID, To: decisionNode.ID})
		}

		// Body
		var bodyEndID string
		if len(s.Body) > 0 {
			currentID := decisionNode.ID
			for _, bodyStep := range s.Body {
				currentID = t.transformStep(bodyStep, currentID, diagram)
			}
			bodyEndID = currentID
			// ループバック
			if bodyEndID != "" {
				diagram.Edges = append(diagram.Edges, flow.Edge{From: bodyEndID, To: decisionNode.ID})
			}
		}

		return decisionNode.ID

	case *ast.WhileStep:
		decisionNode := t.createNode("while?", flow.NodeShapeDecision)
		diagram.Nodes = append(diagram.Nodes, decisionNode)
		if prevNodeID != "" {
			diagram.Edges = append(diagram.Edges, flow.Edge{From: prevNodeID, To: decisionNode.ID})
		}

		// Body
		var bodyEndID string
		if len(s.Body) > 0 {
			currentID := decisionNode.ID
			for _, bodyStep := range s.Body {
				currentID = t.transformStep(bodyStep, currentID, diagram)
			}
			bodyEndID = currentID
			// ループバック
			if bodyEndID != "" {
				diagram.Edges = append(diagram.Edges, flow.Edge{From: bodyEndID, To: decisionNode.ID})
			}
		}

		return decisionNode.ID
	}

	return prevNodeID
}

func (t *FlowTransformer) transformStepWithLabel(step ast.Step, prevNodeID string, diagram *flow.Diagram, label string) string {
	// 最初のノードを作成し、ラベル付きエッジで接続
	switch s := step.(type) {
	case *ast.AssignStep:
		nodeLabel := s.Variable + " = " + t.formatExpr(s.Value)
		node := t.createNode(nodeLabel, flow.NodeShapeProcess)
		diagram.Nodes = append(diagram.Nodes, node)
		if prevNodeID != "" {
			diagram.Edges = append(diagram.Edges, flow.Edge{From: prevNodeID, To: node.ID, Label: label})
		}
		return node.ID
	case *ast.CallStep:
		callLabel := "Call"
		if call, ok := s.Expr.(*ast.CallExpr); ok {
			callLabel = call.Method + "()"
		}
		node := t.createNode(callLabel, flow.NodeShapeProcess)
		diagram.Nodes = append(diagram.Nodes, node)
		if prevNodeID != "" {
			diagram.Edges = append(diagram.Edges, flow.Edge{From: prevNodeID, To: node.ID, Label: label})
		}
		return node.ID
	default:
		// その他のステップは通常の変換
		return t.transformStep(step, prevNodeID, diagram)
	}
}

// formatExpr は式を文字列に整形する
func (t *FlowTransformer) formatExpr(expr ast.Expr) string {
	if expr == nil {
		return "?"
	}

	switch e := expr.(type) {
	case *ast.LiteralExpr:
		return fmt.Sprintf("%v", e.Value)
	case *ast.VariableExpr:
		return e.Name
	case *ast.FieldExpr:
		return t.formatExpr(e.Object) + "." + e.Field
	case *ast.CallExpr:
		obj := t.formatExpr(e.Object)
		args := ""
		for i, arg := range e.Args {
			if i > 0 {
				args += ", "
			}
			args += t.formatExpr(arg)
		}
		return obj + "." + e.Method + "(" + args + ")"
	case *ast.BinaryExpr:
		return t.formatExpr(e.Left) + " " + e.Op + " " + t.formatExpr(e.Right)
	case *ast.UnaryExpr:
		return e.Op + t.formatExpr(e.Operand)
	case *ast.TernaryExpr:
		return t.formatExpr(e.Condition) + " ? " + t.formatExpr(e.Then) + " : " + t.formatExpr(e.Else)
	case *ast.NullishExpr:
		if e.ThrowErr != nil {
			return t.formatExpr(e.Left) + " ?? throw " + *e.ThrowErr
		}
		return t.formatExpr(e.Left) + " ?? " + t.formatExpr(e.Right)
	default:
		return "..."
	}
}
