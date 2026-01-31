package svg

import (
	"fmt"
	"io"
	"strings"

	"pact/internal/domain/diagram/class"
	"pact/internal/domain/diagram/common"
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

	// ノード数に応じてサイズを計算
	cols := 3
	rows := (len(diagram.Nodes) + cols - 1) / cols
	width := 50 + cols*250 + 50
	height := 50 + rows*200 + 50
	if width < 800 {
		width = 800
	}
	if height < 600 {
		height = 600
	}
	c.SetSize(width, height)

	y := 50
	nodePositions := make(map[string]struct{ x, y, width, height int })

	// ノードをレンダリング
	for i, node := range diagram.Nodes {
		x := 50 + (i%cols)*250
		if i > 0 && i%cols == 0 {
			y += 200
		}

		nodeHeight := r.calculateNodeHeight(node)
		nodePositions[node.ID] = struct{ x, y, width, height int }{x, y, 200, nodeHeight}

		// ノードの描画
		r.renderNode(c, node, x, y)
	}

	// エッジをレンダリング
	for _, edge := range diagram.Edges {
		fromPos, fromOk := nodePositions[edge.From]
		toPos, toOk := nodePositions[edge.To]
		if !fromOk || !toOk {
			continue
		}
		// 接続点を計算（右端から左端へ）
		fromX := fromPos.x + fromPos.width
		fromY := fromPos.y + fromPos.height/2
		toX := toPos.x
		toY := toPos.y + toPos.height/2
		r.renderEdge(c, edge, fromX, fromY, toX, toY)
	}

	// ノートをレンダリング
	if len(diagram.Notes) > 0 {
		simplePositions := make(map[string]struct{ x, y int })
		for id, pos := range nodePositions {
			simplePositions[id] = struct{ x, y int }{pos.x, pos.y}
		}
		renderNotes(c, diagram.Notes, simplePositions)
	}

	_, err := c.WriteTo(w)
	return err
}

func (r *ClassRenderer) calculateNodeHeight(node class.Node) int {
	height := 60
	if len(node.Attributes) > 0 {
		height += len(node.Attributes) * 20
	}
	if len(node.Methods) > 0 {
		height += len(node.Methods) * 20
	}
	return height
}

func (r *ClassRenderer) renderNode(c *canvas.Canvas, node class.Node, x, y int) {
	height := r.calculateNodeHeight(node)

	// ノード本体
	c.Rect(x, y, 200, height, canvas.Fill("#fff"), canvas.Stroke("#000"), canvas.StrokeWidth(1))

	// ステレオタイプ
	textY := y + 20
	if node.Stereotype != "" {
		c.Text(x+100, textY, "<<"+node.Stereotype+">>", canvas.TextAnchor("middle"))
		textY += 20
	}

	// 名前
	c.Text(x+100, textY, node.Name, canvas.TextAnchor("middle"))
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
			methodStr := r.formatMethod(method)
			c.Text(x+10, textY+5, vis+methodStr)
			textY += 20
		}
	}
}

// formatMethod はメソッドシグネチャを整形する
func (r *ClassRenderer) formatMethod(method class.Method) string {
	// パラメータリストを構築
	params := ""
	for i, p := range method.Params {
		if i > 0 {
			params += ", "
		}
		if p.Name != "" && p.Type != "" {
			params += p.Name + ": " + p.Type
		} else if p.Type != "" {
			params += p.Type
		} else if p.Name != "" {
			params += p.Name
		}
	}

	// asyncプレフィックス
	prefix := ""
	if method.Async {
		prefix = "async "
	}

	// 戻り型
	returnType := ""
	if method.ReturnType != "" {
		returnType = ": " + method.ReturnType
	}

	return prefix + method.Name + "(" + params + ")" + returnType
}

func (r *ClassRenderer) renderEdge(c *canvas.Canvas, edge class.Edge, x1, y1, x2, y2 int) {
	opts := []canvas.Option{canvas.Stroke("#000")}
	if edge.LineStyle == class.LineStyleDashed {
		opts = append(opts, canvas.Dashed())
	}
	c.Line(x1, y1, x2, y2, opts...)

	// 矢印の先端を描画
	c.Polygon(trianglePoints(x2, y2, x1, y1), canvas.Fill("#000"))

	// 装飾
	switch edge.Decoration {
	case class.DecorationTriangle:
		c.Polygon(trianglePoints(x2, y2, x1, y1), canvas.Fill("#fff"), canvas.Stroke("#000"))
	case class.DecorationFilledDiamond:
		c.Polygon(diamondPoints(x1, y1, x2, y2), canvas.Fill("#000"))
	case class.DecorationEmptyDiamond:
		c.Polygon(diamondPoints(x1, y1, x2, y2), canvas.Fill("#fff"), canvas.Stroke("#000"))
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
		return "+ "
	}
}

func trianglePoints(x, y, fromX, fromY int) string {
	// 矢印の方向を計算
	dx := float64(x - fromX)
	dy := float64(y - fromY)
	length := sqrt(dx*dx + dy*dy)
	if length == 0 {
		return fmt.Sprintf("%d,%d %d,%d %d,%d", x, y-5, x+10, y, x, y+5)
	}
	// 単位ベクトル
	ux := dx / length
	uy := dy / length
	// 垂直ベクトル
	px := -uy
	py := ux
	// 三角形の頂点
	size := 10.0
	ax := float64(x) - ux*size
	ay := float64(y) - uy*size
	p1x := int(ax + px*size/2)
	p1y := int(ay + py*size/2)
	p2x := int(ax - px*size/2)
	p2y := int(ay - py*size/2)
	return fmt.Sprintf("%d,%d %d,%d %d,%d", x, y, p1x, p1y, p2x, p2y)
}

func diamondPoints(x, y, toX, toY int) string {
	// ひし形を始点に描画
	size := 10
	dx := float64(toX - x)
	dy := float64(toY - y)
	length := sqrt(dx*dx + dy*dy)
	if length == 0 {
		return fmt.Sprintf("%d,%d %d,%d %d,%d %d,%d", x, y-size, x+size, y, x, y+size, x-size, y)
	}
	ux := dx / length
	uy := dy / length
	px := -uy
	py := ux
	// ひし形の4頂点
	cx := float64(x) + ux*float64(size)
	cy := float64(y) + uy*float64(size)
	return fmt.Sprintf("%d,%d %.0f,%.0f %.0f,%.0f %.0f,%.0f",
		x, y,
		cx+px*float64(size)/2, cy+py*float64(size)/2,
		float64(x)+ux*float64(size*2), float64(y)+uy*float64(size*2),
		cx-px*float64(size)/2, cy-py*float64(size)/2)
}

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
	}
	return z
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

	// サイズを計算
	width := 100 + len(diagram.Participants)*150 + 50
	if width < 800 {
		width = 800
	}
	c.SetSize(width, 600)

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

	// ノートをレンダリング
	if len(diagram.Notes) > 0 {
		simplePositions := make(map[string]struct{ x, y int })
		for id, x := range participantX {
			simplePositions[id] = struct{ x, y int }{x, 50}
		}
		renderNotes(c, diagram.Notes, simplePositions)
	}

	_, err := c.WriteTo(w)
	return err
}

func (r *SequenceRenderer) renderEvents(c *canvas.Canvas, events []sequence.Event, participantX map[string]int, y *int) {
	// アクティベーション状態を追跡
	activations := make(map[string]int) // participant -> activation start Y

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
				c.Line(fromX, *y, toX, *y, canvas.Stroke("#000"), canvas.Dashed())
				r.drawOpenArrow(c, fromX, toX, *y)
			case sequence.MessageTypeReturn:
				c.Line(fromX, *y, toX, *y, canvas.Stroke("#000"), canvas.Dashed())
				r.drawOpenArrow(c, fromX, toX, *y)
				// returnでアクティベーション終了
				if startY, ok := activations[e.From]; ok {
					r.drawActivationBar(c, fromX, startY, *y)
					delete(activations, e.From)
				}
			default: // sync
				c.Arrow(fromX, *y, toX, *y, canvas.Stroke("#000"))
				// syncでターゲットをアクティベート
				if _, ok := activations[e.To]; !ok {
					activations[e.To] = *y
				}
			}

			// ラベル
			midX := (fromX + toX) / 2
			c.Text(midX, *y-5, e.Label, canvas.TextAnchor("middle"))

			*y += 40

		case *sequence.FragmentEvent:
			// フラグメント（alt, loop, opt）の枠を描画
			startY := *y
			r.renderEvents(c, e.Events, participantX, y)
			// 枠を描画
			c.Rect(50, startY-10, 700, *y-startY+20, canvas.Stroke("#000"), canvas.Fill("none"))
			c.Text(60, startY, "["+string(e.Type)+"] "+e.Label)

		case *sequence.ActivationEvent:
			x, ok := participantX[e.Participant]
			if !ok {
				continue
			}
			if e.Active {
				// アクティベーション開始
				activations[e.Participant] = *y
			} else {
				// アクティベーション終了
				if startY, ok := activations[e.Participant]; ok {
					r.drawActivationBar(c, x, startY, *y)
					delete(activations, e.Participant)
				}
			}
		}
	}

	// 残っているアクティベーションを閉じる
	for participant, startY := range activations {
		if x, ok := participantX[participant]; ok {
			r.drawActivationBar(c, x, startY, *y)
		}
	}
}

// drawActivationBar はアクティベーションバーを描画する
func (r *SequenceRenderer) drawActivationBar(c *canvas.Canvas, x, startY, endY int) {
	barWidth := 10
	c.Rect(x-barWidth/2, startY, barWidth, endY-startY, canvas.Fill("#fff"), canvas.Stroke("#000"))
}

func (r *SequenceRenderer) drawOpenArrow(c *canvas.Canvas, fromX, toX, y int) {
	if toX > fromX {
		c.Line(toX-8, y-5, toX, y, canvas.Stroke("#000"))
		c.Line(toX-8, y+5, toX, y, canvas.Stroke("#000"))
	} else {
		c.Line(toX+8, y-5, toX, y, canvas.Stroke("#000"))
		c.Line(toX+8, y+5, toX, y, canvas.Stroke("#000"))
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
		c.Text(x, y+60, p.Name, canvas.TextAnchor("middle"))
	case sequence.ParticipantTypeDatabase:
		c.Cylinder(x-30, y, 60, 50)
		c.Text(x, y+60, p.Name, canvas.TextAnchor("middle"))
	default:
		c.Rect(x-40, y, 80, 40, canvas.Fill("#fff"), canvas.Stroke("#000"))
		c.Text(x, y+25, p.Name, canvas.TextAnchor("middle"))
	}

	// ライフライン（破線）
	c.Line(x, y+50, x, 500, canvas.Stroke("#000"), canvas.Dashed())
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

	// 状態をタイプ別に分類
	var initialState *state.State
	var finalStates []*state.State
	var normalStates []*state.State

	for i := range diagram.States {
		s := &diagram.States[i]
		switch s.Type {
		case state.StateTypeInitial:
			initialState = s
		case state.StateTypeFinal:
			finalStates = append(finalStates, s)
		default:
			normalStates = append(normalStates, s)
		}
	}

	// レイアウト計算
	cols := 3
	rows := (len(normalStates) + cols - 1) / cols
	width := 100 + cols*200 + 100
	height := 100 + rows*120 + 100
	if len(finalStates) > 0 {
		height += 80
	}
	if width < 800 {
		width = 800
	}
	if height < 600 {
		height = 600
	}
	c.SetSize(width, height)

	// 状態の位置を記録
	statePositions := make(map[string]struct{ x, y int })

	// 初期状態を描画
	if initialState != nil {
		x := 100
		y := 50
		statePositions[initialState.ID] = struct{ x, y int }{x, y}
		c.Circle(x, y, 10, canvas.Fill("#000"))
	}

	// 通常状態を描画
	startY := 120
	for i, s := range normalStates {
		col := i % cols
		row := i / cols
		x := 100 + col*200
		y := startY + row*120
		statePositions[s.ID] = struct{ x, y int }{x, y}
		r.renderState(c, *s, x, y)
	}

	// 終了状態を描画
	finalY := startY + rows*120 + 40
	for i, s := range finalStates {
		x := 100 + i*150
		statePositions[s.ID] = struct{ x, y int }{x, finalY}
		c.Circle(x, finalY, 12, canvas.Stroke("#000"), canvas.StrokeWidth(2))
		c.Circle(x, finalY, 8, canvas.Fill("#000"))
	}

	// 遷移を描画（ラベル位置をずらして重複を避ける）
	labelOffset := make(map[string]int)
	for _, t := range diagram.Transitions {
		fromPos, fromOk := statePositions[t.From]
		toPos, toOk := statePositions[t.To]
		if fromOk && toOk {
			// ラベルオフセットを計算（同じ遷移元のラベルをずらす）
			key := fmt.Sprintf("%d,%d-%d,%d", fromPos.x, fromPos.y, toPos.x, toPos.y)
			offset := labelOffset[key]
			labelOffset[key] = offset + 15
			r.renderTransition(c, t, fromPos.x, fromPos.y, toPos.x, toPos.y, offset)
		}
	}

	// ノートをレンダリング
	if len(diagram.Notes) > 0 {
		renderNotes(c, diagram.Notes, statePositions)
	}

	_, err := c.WriteTo(w)
	return err
}

func (r *StateRenderer) renderState(c *canvas.Canvas, s state.State, x, y int) {
	// 複合状態の場合
	if s.Type == state.StateTypeCompound && len(s.Children) > 0 {
		r.renderCompoundState(c, s, x, y)
		return
	}

	// 並行状態の場合
	if s.Type == state.StateTypeParallel && len(s.Regions) > 0 {
		r.renderParallelState(c, s, x, y)
		return
	}

	// 通常の状態
	hasActions := len(s.Entry) > 0 || len(s.Exit) > 0
	height := 40
	if hasActions {
		height = 60 + len(s.Entry)*15 + len(s.Exit)*15
	}

	c.RoundRect(x-40, y-20, 80, height, 10, 10, canvas.Fill("#fff"), canvas.Stroke("#000"))
	c.Text(x, y+5, s.Name, canvas.TextAnchor("middle"))

	// Entry/Exitアクションを描画
	if hasActions {
		c.Line(x-40, y+15, x+40, y+15, canvas.Stroke("#000"))
		actionY := y + 30
		for _, entry := range s.Entry {
			c.Text(x-35, actionY, "entry/ "+entry)
			actionY += 15
		}
		for _, exit := range s.Exit {
			c.Text(x-35, actionY, "exit/ "+exit)
			actionY += 15
		}
	}
}

// renderCompoundState は複合状態を描画する
func (r *StateRenderer) renderCompoundState(c *canvas.Canvas, s state.State, x, y int) {
	// 子状態の数に基づいてサイズを計算
	childCount := len(s.Children)
	cols := 2
	rows := (childCount + cols - 1) / cols
	if rows < 1 {
		rows = 1
	}

	width := cols*90 + 40
	height := 30 + rows*50 + 20 // ヘッダー + 子状態 + マージン

	// 複合状態の外枠
	c.RoundRect(x-width/2, y-20, width, height, 10, 10, canvas.Fill("#f8f8f8"), canvas.Stroke("#000"))

	// 状態名（上部）
	c.Text(x, y, s.Name, canvas.TextAnchor("middle"))
	c.Line(x-width/2, y+10, x+width/2, y+10, canvas.Stroke("#000"))

	// 子状態を描画
	childX := x - width/2 + 50
	childY := y + 30
	for i, child := range s.Children {
		col := i % cols
		row := i / cols
		cx := childX + col*90
		cy := childY + row*50

		// 子状態を通常の状態として描画
		c.RoundRect(cx-35, cy-15, 70, 30, 8, 8, canvas.Fill("#fff"), canvas.Stroke("#000"))
		c.Text(cx, cy+5, child.Name, canvas.TextAnchor("middle"))
	}
}

// renderParallelState は並行状態を描画する
func (r *StateRenderer) renderParallelState(c *canvas.Canvas, s state.State, x, y int) {
	regionCount := len(s.Regions)
	if regionCount == 0 {
		return
	}

	regionWidth := 100
	width := regionCount*regionWidth + 20
	height := 100

	// 並行状態の外枠
	c.RoundRect(x-width/2, y-20, width, height, 10, 10, canvas.Fill("#f8f8f8"), canvas.Stroke("#000"))

	// 状態名（上部）
	c.Text(x, y, s.Name, canvas.TextAnchor("middle"))
	c.Line(x-width/2, y+10, x+width/2, y+10, canvas.Stroke("#000"))

	// 各リージョンを描画
	for i, region := range s.Regions {
		rx := x - width/2 + 10 + i*regionWidth + regionWidth/2
		ry := y + 30

		// リージョン名
		c.Text(rx, ry, region.Name, canvas.TextAnchor("middle"))

		// リージョン内の状態を簡略表示
		stateY := ry + 25
		for j, child := range region.States {
			if j >= 2 { // 最大2つまで表示
				c.Text(rx, stateY, "...", canvas.TextAnchor("middle"))
				break
			}
			c.RoundRect(rx-30, stateY-10, 60, 20, 5, 5, canvas.Fill("#fff"), canvas.Stroke("#000"))
			c.Text(rx, stateY+5, child.Name, canvas.TextAnchor("middle"))
			stateY += 25
		}

		// リージョン間の区切り線
		if i < regionCount-1 {
			lineX := x - width/2 + 10 + (i+1)*regionWidth
			c.Line(lineX, y+10, lineX, y-20+height, canvas.Stroke("#000"), canvas.Dashed())
		}
	}
}

func (r *StateRenderer) renderTransition(c *canvas.Canvas, t state.Transition, x1, y1, x2, y2 int, labelOffset int) {
	// 矢印を描画
	c.Arrow(x1, y1, x2, y2, canvas.Stroke("#000"))

	// ラベルを構築: トリガー [ガード] / アクション
	midX := (x1 + x2) / 2
	midY := (y1+y2)/2 - 10 - labelOffset

	// トリガー部分
	triggerLabel := ""
	if t.Trigger != nil {
		switch trig := t.Trigger.(type) {
		case *state.EventTrigger:
			triggerLabel = trig.Event
		case *state.WhenTrigger:
			triggerLabel = "when(" + trig.Condition + ")"
		case *state.AfterTrigger:
			triggerLabel = fmt.Sprintf("after %d%s", trig.Duration.Value, trig.Duration.Unit)
		}
	}

	// ガード条件
	guardLabel := ""
	if t.Guard != "" {
		guardLabel = "[" + t.Guard + "]"
	}

	// アクション
	actionLabel := ""
	if len(t.Actions) > 0 {
		actionLabel = "/ "
		for i, action := range t.Actions {
			if i > 0 {
				actionLabel += ", "
			}
			actionLabel += action
		}
	}

	// ラベルを結合
	label := triggerLabel
	if guardLabel != "" {
		if label != "" {
			label += " "
		}
		label += guardLabel
	}
	if actionLabel != "" {
		label += " " + actionLabel
	}

	if label != "" {
		c.Text(midX, midY, label, canvas.TextAnchor("middle"))
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
		c.Text(x, y+25, node.Label, canvas.TextAnchor("middle"))
	}
}

func (r *FlowRenderer) renderFlowEdge(c *canvas.Canvas, edge flow.Edge, x1, y1, x2, y2 int) {
	// 矢印を描画
	c.Arrow(x1, y1, x2, y2, canvas.Stroke("#000"))

	// ラベルがある場合は表示
	if edge.Label != "" {
		midX := (x1 + x2) / 2
		midY := (y1 + y2) / 2
		c.Text(midX+10, midY, edge.Label)
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

// renderNotes はノートを描画する共通関数
func renderNotes(c *canvas.Canvas, notes []common.Note, elementPositions map[string]struct{ x, y int }) {
	noteWidth := 100
	noteHeight := 40

	for _, note := range notes {
		var x, y int

		// 関連付け要素がある場合はその近くに配置
		if note.AttachTo != "" {
			if pos, ok := elementPositions[note.AttachTo]; ok {
				switch note.Position {
				case common.NotePositionLeft:
					x = pos.x - noteWidth - 30
					y = pos.y
				case common.NotePositionRight:
					x = pos.x + 120
					y = pos.y
				case common.NotePositionTop:
					x = pos.x
					y = pos.y - noteHeight - 20
				case common.NotePositionBottom:
					x = pos.x
					y = pos.y + 60
				default:
					x = pos.x + 120
					y = pos.y
				}

				// 接続線を描画
				c.Line(pos.x+60, pos.y+20, x, y+noteHeight/2, canvas.Stroke("#000"), canvas.Dashed())
			} else {
				// 要素が見つからない場合はデフォルト位置
				x = 600
				y = 50
			}
		} else {
			// 関連付けがない場合は右上に配置
			x = 600
			y = 50
		}

		// ノートを描画
		c.Note(x, y, noteWidth, noteHeight, canvas.Fill("#ffffcc"), canvas.Stroke("#000"))

		// テキストを描画（複数行対応）
		lines := strings.Split(note.Text, "\n")
		textY := y + 15
		for _, line := range lines {
			c.Text(x+5, textY, line)
			textY += 15
		}
	}
}
