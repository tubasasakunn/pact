package flow

import "pact/internal/domain/diagram/common"

// Diagram はフローチャートを表す
type Diagram struct {
	Nodes     []Node
	Edges     []Edge
	Swimlanes []Swimlane
}

func (d *Diagram) Type() common.DiagramType {
	return common.DiagramTypeFlow
}

// Node はフローチャートのノード
type Node struct {
	ID       string
	Label    string
	Shape    NodeShape
	Swimlane string
}

// NodeShape はノードの形状
type NodeShape string

const (
	NodeShapeTerminal  NodeShape = "terminal"  // 端子（開始/終了）
	NodeShapeProcess   NodeShape = "process"   // 処理
	NodeShapeDecision  NodeShape = "decision"  // 判断
	NodeShapeIO        NodeShape = "io"        // 入出力
	NodeShapeDatabase  NodeShape = "database"  // データベース
	NodeShapeConnector NodeShape = "connector" // 接続子
)

// Edge はフローチャートのエッジ
type Edge struct {
	From  string
	To    string
	Label string
}

// Swimlane はスイムレーン
type Swimlane struct {
	ID   string
	Name string
}
