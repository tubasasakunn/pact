package svg

import (
	"io"

	"pact/internal/domain/diagram/class"
	"pact/internal/domain/diagram/flow"
	"pact/internal/domain/diagram/sequence"
	"pact/internal/domain/diagram/state"
	"pact/internal/infrastructure/renderer/canvas"
)

// ClassRenderer はクラス図をSVGにレンダリングする
type ClassRenderer struct{}

// NewClassRenderer は新しいClassRendererを作成する
func NewClassRenderer() *ClassRenderer {
	return &ClassRenderer{}
}

// Render はクラス図をSVGにレンダリングする
func (r *ClassRenderer) Render(diagram *class.Diagram, w io.Writer) error {
	c := canvas.New()
	c.SetSize(800, 600)

	y := 50
	nodePositions := make(map[string]struct{ x, y int })

	// ノードをレンダリング
	for i, node := range diagram.Nodes {
		x := 50 + (i%3)*250
		if i > 0 && i%3 == 0 {
			y += 200
		}

		nodePositions[node.ID] = struct{ x, y int }{x, y}

		// ノードの描画
		r.renderNode(c, node, x, y)
	}

	// エッジをレンダリング
	for _, edge := range diagram.Edges {
		fromPos := nodePositions[edge.From]
		toPos := nodePositions[edge.To]
		r.renderEdge(c, edge, fromPos.x+100, fromPos.y+50, toPos.x, toPos.y+50)
	}

	_, err := c.WriteTo(w)
	return err
}

func (r *ClassRenderer) renderNode(c *canvas.Canvas, node class.Node, x, y int) {
	height := 60
	if len(node.Attributes) > 0 {
		height += len(node.Attributes) * 20
	}
	if len(node.Methods) > 0 {
		height += len(node.Methods) * 20
	}

	// ノード本体
	c.Rect(x, y, 200, height, canvas.Fill("#fff"), canvas.Stroke("#000"), canvas.StrokeWidth(1))

	// ステレオタイプ
	textY := y + 20
	if node.Stereotype != "" {
		c.Text(x+100, textY, "<<"+node.Stereotype+">>")
		textY += 20
	}

	// 名前
	c.Text(x+100, textY, node.Name)
	textY += 20

	// 属性
	if len(node.Attributes) > 0 {
		c.Line(x, textY-5, x+200, textY-5, canvas.Stroke("#000"))
		for _, attr := range node.Attributes {
			vis := visibilitySymbol(attr.Visibility)
			c.Text(x+10, textY+5, vis+attr.Name+": "+attr.Type)
			textY += 20
		}
	}

	// メソッド
	if len(node.Methods) > 0 {
		c.Line(x, textY-5, x+200, textY-5, canvas.Stroke("#000"))
		for _, method := range node.Methods {
			vis := visibilitySymbol(class.Visibility(method.Visibility))
			c.Text(x+10, textY+5, vis+method.Name+"()")
			textY += 20
		}
	}
}

func (r *ClassRenderer) renderEdge(c *canvas.Canvas, edge class.Edge, x1, y1, x2, y2 int) {
	opts := []canvas.Option{canvas.Stroke("#000")}
	if edge.LineStyle == class.LineStyleDashed {
		// 点線は簡易的に実線で表現（本来はstroke-dasharray）
		opts = append(opts, canvas.Class("dashed"))
	}
	c.Line(x1, y1, x2, y2, opts...)

	// 装飾
	switch edge.Decoration {
	case class.DecorationTriangle:
		c.Polygon(trianglePoints(x2, y2), canvas.Fill("#fff"), canvas.Stroke("#000"))
	case class.DecorationFilledDiamond:
		c.Polygon(diamondPoints(x1, y1), canvas.Fill("#000"))
	case class.DecorationEmptyDiamond:
		c.Polygon(diamondPoints(x1, y1), canvas.Fill("#fff"), canvas.Stroke("#000"))
	}
}

func visibilitySymbol(v class.Visibility) string {
	switch v {
	case class.VisibilityPublic:
		return "+ "
	case class.VisibilityPrivate:
		return "- "
	case class.VisibilityProtected:
		return "# "
	case class.VisibilityPackage:
		return "~ "
	default:
		return ""
	}
}

func trianglePoints(x, y int) string {
	return "" // 簡易実装
}

func diamondPoints(x, y int) string {
	return "" // 簡易実装
}

// SequenceRenderer はシーケンス図をSVGにレンダリングする
type SequenceRenderer struct{}

// NewSequenceRenderer は新しいSequenceRendererを作成する
func NewSequenceRenderer() *SequenceRenderer {
	return &SequenceRenderer{}
}

// Render はシーケンス図をSVGにレンダリングする
func (r *SequenceRenderer) Render(diagram *sequence.Diagram, w io.Writer) error {
	c := canvas.New()
	c.SetSize(800, 600)

	// 参加者をレンダリング
	for i, p := range diagram.Participants {
		x := 100 + i*150
		r.renderParticipant(c, p, x, 50)
	}

	_, err := c.WriteTo(w)
	return err
}

func (r *SequenceRenderer) renderParticipant(c *canvas.Canvas, p sequence.Participant, x, y int) {
	switch p.Type {
	case sequence.ParticipantTypeActor:
		// 人型を描画
		c.Circle(x, y, 10, canvas.Stroke("#000"), canvas.Fill("#fff"))
		c.Line(x, y+10, x, y+30, canvas.Stroke("#000"))
		c.Line(x-10, y+15, x+10, y+15, canvas.Stroke("#000"))
		c.Line(x, y+30, x-10, y+45, canvas.Stroke("#000"))
		c.Line(x, y+30, x+10, y+45, canvas.Stroke("#000"))
	case sequence.ParticipantTypeDatabase:
		c.Cylinder(x-30, y, 60, 50)
	default:
		c.Rect(x-40, y, 80, 40, canvas.Fill("#fff"), canvas.Stroke("#000"))
		c.Text(x, y+25, p.Name)
	}

	// ライフライン
	c.Line(x, y+50, x, 500, canvas.Stroke("#000"), canvas.Class("dashed"))
}

// StateRenderer は状態図をSVGにレンダリングする
type StateRenderer struct{}

// NewStateRenderer は新しいStateRendererを作成する
func NewStateRenderer() *StateRenderer {
	return &StateRenderer{}
}

// Render は状態図をSVGにレンダリングする
func (r *StateRenderer) Render(diagram *state.Diagram, w io.Writer) error {
	c := canvas.New()
	c.SetSize(800, 600)

	y := 50
	for i, s := range diagram.States {
		x := 100 + (i%3)*200
		if i > 0 && i%3 == 0 {
			y += 150
		}
		r.renderState(c, s, x, y)
	}

	_, err := c.WriteTo(w)
	return err
}

func (r *StateRenderer) renderState(c *canvas.Canvas, s state.State, x, y int) {
	switch s.Type {
	case state.StateTypeInitial:
		c.Circle(x, y, 10, canvas.Fill("#000"))
	case state.StateTypeFinal:
		c.Circle(x, y, 12, canvas.Stroke("#000"), canvas.StrokeWidth(2))
		c.Circle(x, y, 8, canvas.Fill("#000"))
	default:
		c.RoundRect(x-40, y-20, 80, 40, 10, 10, canvas.Fill("#fff"), canvas.Stroke("#000"))
		c.Text(x, y+5, s.Name)
	}
}

// FlowRenderer はフローチャートをSVGにレンダリングする
type FlowRenderer struct{}

// NewFlowRenderer は新しいFlowRendererを作成する
func NewFlowRenderer() *FlowRenderer {
	return &FlowRenderer{}
}

// Render はフローチャートをSVGにレンダリングする
func (r *FlowRenderer) Render(diagram *flow.Diagram, w io.Writer) error {
	c := canvas.New()
	c.SetSize(800, 600)

	y := 50
	for i, node := range diagram.Nodes {
		x := 400
		r.renderFlowNode(c, node, x, y+i*80)
	}

	_, err := c.WriteTo(w)
	return err
}

func (r *FlowRenderer) renderFlowNode(c *canvas.Canvas, node flow.Node, x, y int) {
	switch node.Shape {
	case flow.NodeShapeTerminal:
		c.Stadium(x-50, y, 100, 40, canvas.Fill("#fff"), canvas.Stroke("#000"))
	case flow.NodeShapeProcess:
		c.Rect(x-50, y, 100, 40, canvas.Fill("#fff"), canvas.Stroke("#000"))
	case flow.NodeShapeDecision:
		c.Diamond(x, y+20, 80, 40, canvas.Fill("#fff"), canvas.Stroke("#000"))
	case flow.NodeShapeDatabase:
		c.Cylinder(x-30, y, 60, 50)
	case flow.NodeShapeIO:
		c.Parallelogram(x-50, y, 100, 40, 15, canvas.Fill("#fff"), canvas.Stroke("#000"))
	default:
		c.Rect(x-50, y, 100, 40, canvas.Fill("#fff"), canvas.Stroke("#000"))
	}

	if node.Label != "" {
		c.Text(x, y+25, node.Label)
	}
}
