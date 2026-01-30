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

	// 参加者の位置を記録
	participantX := make(map[string]int)

	// 参加者をレンダリング
	for i, p := range diagram.Participants {
		x := 100 + i*150
		participantX[p.ID] = x
		r.renderParticipant(c, p, x, 50)
	}

	// メッセージをレンダリング
	messageY := 120
	r.renderEvents(c, diagram.Events, participantX, &messageY)

	_, err := c.WriteTo(w)
	return err
}

func (r *SequenceRenderer) renderEvents(c *canvas.Canvas, events []sequence.Event, participantX map[string]int, y *int) {
	for _, event := range events {
		switch e := event.(type) {
		case *sequence.MessageEvent:
			fromX, fromOk := participantX[e.From]
			toX, toOk := participantX[e.To]
			if !fromOk || !toOk {
				continue
			}

			// メッセージの矢印を描画
			switch e.MessageType {
			case sequence.MessageTypeAsync:
				c.Line(fromX, *y, toX, *y, canvas.Stroke("#000"), canvas.Class("dashed"))
				// 矢印の先端
				if toX > fromX {
					c.Line(toX-8, *y-5, toX, *y, canvas.Stroke("#000"))
					c.Line(toX-8, *y+5, toX, *y, canvas.Stroke("#000"))
				} else {
					c.Line(toX+8, *y-5, toX, *y, canvas.Stroke("#000"))
					c.Line(toX+8, *y+5, toX, *y, canvas.Stroke("#000"))
				}
			case sequence.MessageTypeReturn:
				c.Line(fromX, *y, toX, *y, canvas.Stroke("#000"), canvas.Class("dashed"))
				// 開いた矢印
				if toX > fromX {
					c.Line(toX-8, *y-5, toX, *y, canvas.Stroke("#000"))
					c.Line(toX-8, *y+5, toX, *y, canvas.Stroke("#000"))
				} else {
					c.Line(toX+8, *y-5, toX, *y, canvas.Stroke("#000"))
					c.Line(toX+8, *y+5, toX, *y, canvas.Stroke("#000"))
				}
			default: // sync
				c.Arrow(fromX, *y, toX, *y, canvas.Stroke("#000"))
			}

			// ラベル
			midX := (fromX + toX) / 2
			c.Text(midX, *y-5, e.Label)

			*y += 40

		case *sequence.FragmentEvent:
			// フラグメント（alt, loop, opt）の枠を描画
			startY := *y
			r.renderEvents(c, e.Events, participantX, y)
			// 枠を描画
			c.Rect(50, startY-10, 700, *y-startY+20, canvas.Stroke("#000"), canvas.Fill("none"))
			c.Text(60, startY, "["+string(e.Type)+"] "+e.Label)
		}
	}
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

	// 状態の位置を記録
	statePositions := make(map[string]struct{ x, y int })

	y := 50
	for i, s := range diagram.States {
		x := 100 + (i%3)*200
		if i > 0 && i%3 == 0 {
			y += 150
		}
		statePositions[s.ID] = struct{ x, y int }{x, y}
		r.renderState(c, s, x, y)
	}

	// 遷移を描画
	for _, t := range diagram.Transitions {
		fromPos, fromOk := statePositions[t.From]
		toPos, toOk := statePositions[t.To]
		if fromOk && toOk {
			r.renderTransition(c, t, fromPos.x, fromPos.y, toPos.x, toPos.y)
		}
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

func (r *StateRenderer) renderTransition(c *canvas.Canvas, t state.Transition, x1, y1, x2, y2 int) {
	// 矢印を描画
	c.Arrow(x1, y1, x2, y2, canvas.Stroke("#000"))

	// ラベルを描画（イベント名）
	midX := (x1 + x2) / 2
	midY := (y1 + y2) / 2
	label := ""
	if t.Trigger != nil {
		switch trig := t.Trigger.(type) {
		case *state.EventTrigger:
			label = trig.Event
		case *state.WhenTrigger:
			label = "[" + trig.Condition + "]"
		}
	}
	if label != "" {
		c.Text(midX, midY-5, label)
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

	// 必要な高さを計算
	height := 100 + len(diagram.Nodes)*80
	if height < 600 {
		height = 600
	}
	c.SetSize(800, height)

	// ノードの位置を記録
	nodePositions := make(map[string]struct{ x, y int })

	y := 50
	for i, node := range diagram.Nodes {
		x := 400
		nodeY := y + i*80
		nodePositions[node.ID] = struct{ x, y int }{x, nodeY}
		r.renderFlowNode(c, node, x, nodeY)
	}

	// エッジを描画
	for _, edge := range diagram.Edges {
		fromPos, fromOk := nodePositions[edge.From]
		toPos, toOk := nodePositions[edge.To]
		if fromOk && toOk {
			r.renderFlowEdge(c, edge, fromPos.x, fromPos.y+40, toPos.x, toPos.y)
		}
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

func (r *FlowRenderer) renderFlowEdge(c *canvas.Canvas, edge flow.Edge, x1, y1, x2, y2 int) {
	// 矢印を描画
	c.Arrow(x1, y1, x2, y2, canvas.Stroke("#000"))

	// ラベルがある場合は表示（条件分岐のYes/Noなど）
	if edge.Label != "" {
		midX := (x1 + x2) / 2
		midY := (y1 + y2) / 2
		c.Text(midX+10, midY, edge.Label)
	}
}
