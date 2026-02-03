package svg

import (
	"io"

	"pact/internal/domain/diagram/flow"
	"pact/internal/infrastructure/renderer/canvas"
)

// FlowRenderer はフローチャートをSVGにレンダリングする
type FlowRenderer struct{}

// NewFlowRenderer は新しいFlowRendererを作成する
func NewFlowRenderer() *FlowRenderer {
	return &FlowRenderer{}
}

// Render はフローチャートをSVGにレンダリングする
func (r *FlowRenderer) Render(diagram *flow.Diagram, w io.Writer) error {
	c := canvas.New()

	// ノードの位置を計算
	nodePositions := make(map[string]struct{ x, y int })
	nodeInfo := make(map[string]flow.Node)

	// ノードをIDでマップ
	for _, node := range diagram.Nodes {
		nodeInfo[node.ID] = node
	}

	// エッジから接続情報を構築
	outEdges := make(map[string][]flow.Edge)
	for _, edge := range diagram.Edges {
		outEdges[edge.From] = append(outEdges[edge.From], edge)
	}

	// スイムレーンがあるかチェック
	hasSwimlanes := len(diagram.Swimlanes) > 0

	if hasSwimlanes {
		r.renderWithSwimlanes(c, diagram, nodePositions, nodeInfo)
	} else {
		r.renderWithoutSwimlanes(c, diagram, nodePositions, nodeInfo)
	}

	// エッジを描画
	for _, edge := range diagram.Edges {
		fromPos, fromOk := nodePositions[edge.From]
		toPos, toOk := nodePositions[edge.To]
		if !fromOk || !toOk {
			continue
		}

		fromNode := nodeInfo[edge.From]
		fromHeight := 40

		if fromNode.Shape == flow.NodeShapeDecision {
			// 分岐ノードの場合
			if edge.Label == "Yes" {
				// 下向き
				r.renderFlowEdge(c, edge, fromPos.x, fromPos.y+fromHeight, toPos.x, toPos.y)
			} else if edge.Label == "No" {
				// 右向きに曲げる
				r.renderBranchEdge(c, edge, fromPos.x+40, fromPos.y+20, toPos.x, toPos.y)
			} else {
				r.renderFlowEdge(c, edge, fromPos.x, fromPos.y+fromHeight, toPos.x, toPos.y)
			}
		} else {
			r.renderFlowEdge(c, edge, fromPos.x, fromPos.y+fromHeight, toPos.x, toPos.y)
		}
	}

	// ノートをレンダリング
	if len(diagram.Notes) > 0 {
		renderNotes(c, diagram.Notes, nodePositions)
	}

	_, err := c.WriteTo(w)
	return err
}

// renderWithSwimlanes はスイムレーン付きでレンダリングする
func (r *FlowRenderer) renderWithSwimlanes(c *canvas.Canvas, diagram *flow.Diagram, nodePositions map[string]struct{ x, y int }, nodeInfo map[string]flow.Node) {
	swimlaneWidth := 200
	headerHeight := 40
	yStep := 80
	startY := headerHeight + 30

	// スイムレーンIDからインデックスへのマップ
	swimlaneIndex := make(map[string]int)
	for i, sl := range diagram.Swimlanes {
		swimlaneIndex[sl.ID] = i
	}

	// 各スイムレーンのノードをY座標でトラッキング
	swimlaneY := make(map[string]int)
	for _, sl := range diagram.Swimlanes {
		swimlaneY[sl.ID] = startY
	}

	// ノードを順番に配置
	maxY := startY
	for _, node := range diagram.Nodes {
		slIdx, ok := swimlaneIndex[node.Swimlane]
		if !ok {
			slIdx = 0 // デフォルトは最初のスイムレーン
		}

		x := 50 + slIdx*swimlaneWidth + swimlaneWidth/2
		y := swimlaneY[node.Swimlane]

		// 分岐先のNoラベルがある場合はY位置を調整
		for _, edge := range diagram.Edges {
			if edge.To == node.ID && edge.Label == "No" {
				// Noの場合は少し右にオフセット
				x += 60
				break
			}
		}

		nodePositions[node.ID] = struct{ x, y int }{x, y}
		r.renderFlowNode(c, node, x, y)

		swimlaneY[node.Swimlane] = y + yStep
		if y+yStep > maxY {
			maxY = y + yStep
		}
	}

	// 高さを計算
	height := maxY + 50
	if height < 600 {
		height = 600
	}
	width := 50 + len(diagram.Swimlanes)*swimlaneWidth + 50
	if width < 800 {
		width = 800
	}
	c.SetSize(width, height)

	// スイムレーンの枠とヘッダーを描画
	for i, sl := range diagram.Swimlanes {
		x := 50 + i*swimlaneWidth

		// ヘッダー背景
		c.Rect(x, 0, swimlaneWidth, headerHeight, canvas.Fill("#e0e0e0"), canvas.Stroke("#000"))
		// ヘッダーテキスト
		c.Text(x+swimlaneWidth/2, headerHeight/2+5, sl.Name, canvas.TextAnchor("middle"))

		// スイムレーンの縦線
		c.Line(x, 0, x, height, canvas.Stroke("#000"))
	}
	// 最後の縦線
	c.Line(50+len(diagram.Swimlanes)*swimlaneWidth, 0, 50+len(diagram.Swimlanes)*swimlaneWidth, height, canvas.Stroke("#000"))
	// ヘッダーの下線
	c.Line(50, headerHeight, 50+len(diagram.Swimlanes)*swimlaneWidth, headerHeight, canvas.Stroke("#000"))
}

// renderWithoutSwimlanes はスイムレーンなしでレンダリングする
func (r *FlowRenderer) renderWithoutSwimlanes(c *canvas.Canvas, diagram *flow.Diagram, nodePositions map[string]struct{ x, y int }, nodeInfo map[string]flow.Node) {
	mainX := 400
	branchX := 550
	y := 50
	yStep := 80

	for _, node := range diagram.Nodes {
		x := mainX
		for _, edge := range diagram.Edges {
			if edge.To == node.ID && edge.Label == "No" {
				x = branchX
				break
			}
		}
		nodePositions[node.ID] = struct{ x, y int }{x, y}
		r.renderFlowNode(c, node, x, y)
		y += yStep
	}

	height := y + 50
	if height < 600 {
		height = 600
	}
	c.SetSize(800, height)
}

// calculateFlowNodeWidth はフローノードの幅を計算する
func (r *FlowRenderer) calculateFlowNodeWidth(node flow.Node) int {
	minWidth := 100
	padding := 30
	fontSize := 12

	if node.Label == "" {
		return minWidth
	}

	textWidth, _ := canvas.MeasureText(node.Label, fontSize)
	width := textWidth + padding
	if width < minWidth {
		width = minWidth
	}
	return width
}

func (r *FlowRenderer) renderFlowNode(c *canvas.Canvas, node flow.Node, x, y int) {
	width := r.calculateFlowNodeWidth(node)
	r.renderFlowNodeWithWidth(c, node, x, y, width)
}

func (r *FlowRenderer) renderFlowNodeWithWidth(c *canvas.Canvas, node flow.Node, x, y, width int) {
	switch node.Shape {
	case flow.NodeShapeTerminal:
		c.Stadium(x-width/2, y, width, 40, canvas.Fill("#fff"), canvas.Stroke("#000"))
	case flow.NodeShapeProcess:
		c.Rect(x-width/2, y, width, 40, canvas.Fill("#fff"), canvas.Stroke("#000"))
	case flow.NodeShapeDecision:
		// 判断ノードは幅を少し大きめに
		diamondWidth := width
		if diamondWidth < 80 {
			diamondWidth = 80
		}
		c.Diamond(x, y+20, diamondWidth, 40, canvas.Fill("#fff"), canvas.Stroke("#000"))
	case flow.NodeShapeDatabase:
		c.Cylinder(x-width/2, y, width, 50)
	case flow.NodeShapeIO:
		c.Parallelogram(x-width/2, y, width, 40, 15, canvas.Fill("#fff"), canvas.Stroke("#000"))
	default:
		c.Rect(x-width/2, y, width, 40, canvas.Fill("#fff"), canvas.Stroke("#000"))
	}

	if node.Label != "" {
		c.Text(x, y+25, node.Label, canvas.TextAnchor("middle"))
	}
}

func (r *FlowRenderer) renderFlowEdge(c *canvas.Canvas, edge flow.Edge, x1, y1, x2, y2 int) {
	// 直交ルーティングで矢印を描画
	if x1 == x2 {
		// 垂直方向のみ
		c.Line(x1, y1, x2, y2, canvas.Stroke("#000"))
		c.DrawArrowHead(x2, y2, x1, y1)
	} else if y1 == y2 {
		// 水平方向のみ
		c.Line(x1, y1, x2, y2, canvas.Stroke("#000"))
		c.DrawArrowHead(x2, y2, x1, y1)
	} else {
		// L字型ルーティング
		midY := (y1 + y2) / 2
		c.Line(x1, y1, x1, midY, canvas.Stroke("#000"))
		c.Line(x1, midY, x2, midY, canvas.Stroke("#000"))
		c.Line(x2, midY, x2, y2, canvas.Stroke("#000"))
		c.DrawArrowHead(x2, y2, x2, midY)
	}

	// ラベルがある場合は表示
	if edge.Label != "" {
		var labelX, labelY int
		if x1 == x2 {
			labelX = x1 + 10
			labelY = (y1 + y2) / 2
		} else {
			labelX = (x1 + x2) / 2
			labelY = (y1 + y2) / 2
		}
		c.Text(labelX, labelY, edge.Label)
	}
}

func (r *FlowRenderer) renderBranchEdge(c *canvas.Canvas, edge flow.Edge, x1, y1, x2, y2 int) {
	// L字型のパスを描画（右に出てから下に曲がる）
	midX := x2
	c.Line(x1, y1, midX, y1, canvas.Stroke("#000"))
	c.Arrow(midX, y1, x2, y2, canvas.Stroke("#000"))

	// ラベル
	if edge.Label != "" {
		c.Text(x1+20, y1-5, edge.Label)
	}
}
