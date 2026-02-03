package svg

import (
	"fmt"
	"io"

	"pact/internal/domain/diagram/state"
	"pact/internal/infrastructure/renderer/canvas"
)

// StateRenderer は状態図をSVGにレンダリングする
type StateRenderer struct{}

// stateRect は状態のバウンディングボックスを表す
type stateRect struct {
	x, y, w, h int
}

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

	// 各状態の幅を事前計算
	stateWidths := make(map[string]int)
	for _, s := range normalStates {
		stateWidths[s.ID] = r.calculateStateWidth(*s)
	}

	// 列ごとの最大幅を計算
	cols := 3
	rows := (len(normalStates) + cols - 1) / cols
	colMaxWidths := make([]int, cols)
	for i, s := range normalStates {
		col := i % cols
		w := stateWidths[s.ID]
		if w > colMaxWidths[col] {
			colMaxWidths[col] = w
		}
	}

	// 列の開始X座標を計算（動的間隔）
	margin := 50   // 左マージン
	colGap := 40   // 列間のギャップ
	colCenters := make([]int, cols)
	x := margin
	for col := 0; col < cols; col++ {
		colWidth := colMaxWidths[col]
		if colWidth < 100 {
			colWidth = 100
		}
		colCenters[col] = x + colWidth/2
		x += colWidth + colGap
	}

	// キャンバスサイズを計算
	totalWidth := x + margin
	height := 100 + rows*120 + 100
	if len(finalStates) > 0 {
		height += 80
	}
	if totalWidth < 800 {
		totalWidth = 800
	}
	if height < 600 {
		height = 600
	}
	c.SetSize(totalWidth, height)

	// 状態の位置とサイズを記録
	statePositions := make(map[string]struct{ x, y int })
	stateSizes := make(map[string]struct{ w, h int })

	// 初期状態を描画
	if initialState != nil {
		ix := colCenters[0]
		iy := 50
		statePositions[initialState.ID] = struct{ x, y int }{ix, iy}
		stateSizes[initialState.ID] = struct{ w, h int }{20, 20}
		c.Circle(ix, iy, 10, canvas.Fill("#000"))
	}

	// 通常状態を描画
	startY := 120
	for i, s := range normalStates {
		col := i % cols
		row := i / cols
		sx := colCenters[col]
		sy := startY + row*120
		statePositions[s.ID] = struct{ x, y int }{sx, sy}
		sw := stateWidths[s.ID]
		sh := r.calculateStateHeight(*s)
		stateSizes[s.ID] = struct{ w, h int }{sw, sh}
		r.renderState(c, *s, sx, sy)
	}

	// 終了状態を描画
	finalY := startY + rows*120 + 40
	for i, s := range finalStates {
		fx := colCenters[0] + i*150
		statePositions[s.ID] = struct{ x, y int }{fx, finalY}
		stateSizes[s.ID] = struct{ w, h int }{24, 24}
		c.Circle(fx, finalY, 12, canvas.Stroke("#000"), canvas.StrokeWidth(2))
		c.Circle(fx, finalY, 8, canvas.Fill("#000"))
	}

	// 遷移を描画（直交ルーティング）
	labelOffset := make(map[string]int)
	// ノードのバウンディングボックスリストを作成
	var nodeBounds []stateRect
	for id, pos := range statePositions {
		size := stateSizes[id]
		nodeBounds = append(nodeBounds, stateRect{
			x: pos.x - size.w/2,
			y: pos.y - size.h/2,
			w: size.w,
			h: size.h,
		})
	}

	for _, t := range diagram.Transitions {
		fromPos, fromOk := statePositions[t.From]
		toPos, toOk := statePositions[t.To]
		if fromOk && toOk {
			fromSize := stateSizes[t.From]
			toSize := stateSizes[t.To]
			key := fmt.Sprintf("%d,%d-%d,%d", fromPos.x, fromPos.y, toPos.x, toPos.y)
			offset := labelOffset[key]
			labelOffset[key] = offset + 15
			r.renderOrthogonalTransition(c, t, fromPos.x, fromPos.y, fromSize.w, fromSize.h,
				toPos.x, toPos.y, toSize.w, toSize.h, offset, nodeBounds)
		}
	}

	// ノートをレンダリング
	if len(diagram.Notes) > 0 {
		renderNotes(c, diagram.Notes, statePositions)
	}

	_, err := c.WriteTo(w)
	return err
}

// calculateStateHeight は状態ボックスの高さを計算する
func (r *StateRenderer) calculateStateHeight(s state.State) int {
	hasActions := len(s.Entry) > 0 || len(s.Exit) > 0
	if hasActions {
		return 60 + len(s.Entry)*15 + len(s.Exit)*15
	}
	return 40
}

// calculateStateWidth は状態ボックスの幅を計算する
func (r *StateRenderer) calculateStateWidth(s state.State) int {
	minWidth := 80
	padding := 20
	fontSize := 12

	maxWidth := 0

	// 状態名の幅
	nameWidth, _ := canvas.MeasureText(s.Name, fontSize)
	if nameWidth > maxWidth {
		maxWidth = nameWidth
	}

	// Entry/Exitアクションの幅
	for _, entry := range s.Entry {
		actionWidth, _ := canvas.MeasureText("entry/ "+entry, fontSize)
		if actionWidth > maxWidth {
			maxWidth = actionWidth
		}
	}
	for _, exit := range s.Exit {
		actionWidth, _ := canvas.MeasureText("exit/ "+exit, fontSize)
		if actionWidth > maxWidth {
			maxWidth = actionWidth
		}
	}

	width := maxWidth + padding
	if width < minWidth {
		width = minWidth
	}
	return width
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
	width := r.calculateStateWidth(s)
	hasActions := len(s.Entry) > 0 || len(s.Exit) > 0
	height := 40
	if hasActions {
		height = 60 + len(s.Entry)*15 + len(s.Exit)*15
	}

	c.RoundRect(x-width/2, y-20, width, height, 10, 10, canvas.Fill("#fff"), canvas.Stroke("#000"))
	c.Text(x, y+5, s.Name, canvas.TextAnchor("middle"))

	// Entry/Exitアクションを描画
	if hasActions {
		c.Line(x-width/2, y+15, x+width/2, y+15, canvas.Stroke("#000"))
		actionY := y + 30
		for _, entry := range s.Entry {
			c.Text(x-width/2+5, actionY, "entry/ "+entry)
			actionY += 15
		}
		for _, exit := range s.Exit {
			c.Text(x-width/2+5, actionY, "exit/ "+exit)
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

// renderOrthogonalTransition は直交ルーティングで遷移を描画する
func (r *StateRenderer) renderOrthogonalTransition(c *canvas.Canvas, t state.Transition,
	x1, y1, w1, h1, x2, y2, w2, h2 int, labelOffset int, nodeBounds []stateRect) {

	// 出発点と到着点を計算（ボックスの端）
	var startX, startY, endX, endY int
	var midX, midY int

	// 方向を判定
	dx := x2 - x1
	dy := y2 - y1

	// 始点・終点ノードを除いた障害物リスト
	obstacles := r.filterObstacles(nodeBounds, x1, y1, w1, h1, x2, y2, w2, h2)

	if dx == 0 && dy == 0 {
		// 自己遷移
		startX = x1 + w1/2
		startY = y1
		endX = x1
		endY = y1 - h1/2 - 10
		c.Line(startX, startY, startX+30, startY, canvas.Stroke("#000"))
		c.Line(startX+30, startY, startX+30, startY-30, canvas.Stroke("#000"))
		c.Line(startX+30, startY-30, endX, startY-30, canvas.Stroke("#000"))
		c.Line(endX, startY-30, endX, endY, canvas.Stroke("#000"))
		c.DrawArrowHead(endX, endY, endX, startY-30)
		midX = startX + 15
		midY = startY - 15
	} else if abs(dx) > abs(dy) {
		// 主に水平方向
		if dx > 0 {
			startX = x1 + w1/2
			startY = y1
			endX = x2 - w2/2
			endY = y2
		} else {
			startX = x1 - w1/2
			startY = y1
			endX = x2 + w2/2
			endY = y2
		}

		// ウェイポイントを計算（障害物回避付き）
		waypoints := r.calculateStateWaypoints(startX, startY, endX, endY, obstacles)
		midX, midY = r.drawWaypointsAndGetMid(c, waypoints)
	} else {
		// 主に垂直方向
		if dy > 0 {
			startX = x1
			startY = y1 + h1/2
			endX = x2
			endY = y2 - h2/2
		} else {
			startX = x1
			startY = y1 - h1/2
			endX = x2
			endY = y2 + h2/2
		}

		// ウェイポイントを計算（障害物回避付き）
		waypoints := r.calculateStateWaypoints(startX, startY, endX, endY, obstacles)
		midX, midY = r.drawWaypointsAndGetMid(c, waypoints)
	}

	// ラベルを構築
	label := r.buildTransitionLabel(t)
	if label != "" {
		labelX, labelY := r.findSafeLabelPosition(midX, midY-5-labelOffset, label, nodeBounds)
		c.Text(labelX, labelY, label, canvas.TextAnchor("middle"))
	}
}

// filterObstacles は始点・終点ノードを除いた障害物リストを返す
func (r *StateRenderer) filterObstacles(nodeBounds []stateRect, x1, y1, w1, h1, x2, y2, w2, h2 int) []stateRect {
	var obstacles []stateRect
	// 始点ノードの矩形
	srcRect := stateRect{x: x1 - w1/2, y: y1 - h1/2, w: w1, h: h1}
	// 終点ノードの矩形
	dstRect := stateRect{x: x2 - w2/2, y: y2 - h2/2, w: w2, h: h2}

	for _, node := range nodeBounds {
		// 始点・終点と重なるノードは除外
		if r.rectsOverlap(node, srcRect) || r.rectsOverlap(node, dstRect) {
			continue
		}
		obstacles = append(obstacles, node)
	}
	return obstacles
}

// rectsOverlap は2つの矩形が重なるかチェック
func (r *StateRenderer) rectsOverlap(a, b stateRect) bool {
	return a.x < b.x+b.w && a.x+a.w > b.x && a.y < b.y+b.h && a.y+a.h > b.y
}

// calculateStateWaypoints は障害物を回避するウェイポイントを計算
func (r *StateRenderer) calculateStateWaypoints(x1, y1, x2, y2 int, obstacles []stateRect) []struct{ x, y int } {
	start := struct{ x, y int }{x1, y1}
	end := struct{ x, y int }{x2, y2}

	// 直線でOKかチェック（水平または垂直のみ）
	if x1 == x2 || y1 == y2 {
		if !r.pathIntersectsObstacles(x1, y1, x2, y2, obstacles) {
			return []struct{ x, y int }{start, end}
		}
	}

	// L字型を試す（縦→横）
	corner1 := struct{ x, y int }{x1, y2}
	if !r.pathIntersectsObstacles(x1, y1, x1, y2, obstacles) &&
		!r.pathIntersectsObstacles(x1, y2, x2, y2, obstacles) {
		return []struct{ x, y int }{start, corner1, end}
	}

	// L字型を試す（横→縦）
	corner2 := struct{ x, y int }{x2, y1}
	if !r.pathIntersectsObstacles(x1, y1, x2, y1, obstacles) &&
		!r.pathIntersectsObstacles(x2, y1, x2, y2, obstacles) {
		return []struct{ x, y int }{start, corner2, end}
	}

	// Z字型を試す
	midY := (y1 + y2) / 2
	mid1 := struct{ x, y int }{x1, midY}
	mid2 := struct{ x, y int }{x2, midY}
	if !r.pathIntersectsObstacles(x1, y1, x1, midY, obstacles) &&
		!r.pathIntersectsObstacles(x1, midY, x2, midY, obstacles) &&
		!r.pathIntersectsObstacles(x2, midY, x2, y2, obstacles) {
		return []struct{ x, y int }{start, mid1, mid2, end}
	}

	// 迂回ルート（上または下を回る）
	margin := 30
	// 全障害物の最小Y/最大Yを求めて迂回
	minY, maxY := y1, y1
	for _, obs := range obstacles {
		if obs.y < minY {
			minY = obs.y
		}
		if obs.y+obs.h > maxY {
			maxY = obs.y + obs.h
		}
	}

	// 上を回る
	topY := minY - margin
	if topY > 0 {
		wp1 := struct{ x, y int }{x1, topY}
		wp2 := struct{ x, y int }{x2, topY}
		return []struct{ x, y int }{start, wp1, wp2, end}
	}

	// 下を回る
	bottomY := maxY + margin
	wp1 := struct{ x, y int }{x1, bottomY}
	wp2 := struct{ x, y int }{x2, bottomY}
	return []struct{ x, y int }{start, wp1, wp2, end}
}

// pathIntersectsObstacles はパスが障害物と交差するかチェック
func (r *StateRenderer) pathIntersectsObstacles(x1, y1, x2, y2 int, obstacles []stateRect) bool {
	// 簡易判定：パスのバウンディングボックスと各障害物の重なりチェック
	pminX, pmaxX := x1, x2
	if x1 > x2 {
		pminX, pmaxX = x2, x1
	}
	pminY, pmaxY := y1, y2
	if y1 > y2 {
		pminY, pmaxY = y2, y1
	}

	// パスに厚みを持たせる
	padding := 5
	pminX -= padding
	pmaxX += padding
	pminY -= padding
	pmaxY += padding

	for _, obs := range obstacles {
		// 障害物のバウンディングボックスとの交差チェック
		if pmaxX > obs.x && pminX < obs.x+obs.w && pmaxY > obs.y && pminY < obs.y+obs.h {
			return true
		}
	}
	return false
}

// drawWaypointsAndGetMid はウェイポイントを描画し、中間点を返す
func (r *StateRenderer) drawWaypointsAndGetMid(c *canvas.Canvas, waypoints []struct{ x, y int }) (int, int) {
	if len(waypoints) < 2 {
		return 0, 0
	}

	// パスを描画
	for i := 0; i < len(waypoints)-1; i++ {
		c.Line(waypoints[i].x, waypoints[i].y, waypoints[i+1].x, waypoints[i+1].y, canvas.Stroke("#000"))
	}

	// 矢印を描画
	last := waypoints[len(waypoints)-1]
	prev := waypoints[len(waypoints)-2]
	c.DrawArrowHead(last.x, last.y, prev.x, prev.y)

	// 中間点を計算（パスの中央セグメント上）
	midIdx := len(waypoints) / 2
	if midIdx > 0 {
		midX := (waypoints[midIdx-1].x + waypoints[midIdx].x) / 2
		midY := (waypoints[midIdx-1].y + waypoints[midIdx].y) / 2
		return midX, midY
	}
	return (waypoints[0].x + waypoints[1].x) / 2, (waypoints[0].y + waypoints[1].y) / 2
}

// findSafeLabelPosition はラベルがノードと重ならない位置を探す
func (r *StateRenderer) findSafeLabelPosition(x, y int, label string, nodeBounds []stateRect) (int, int) {
	// ラベルのおおよそのサイズを推定
	labelWidth, _ := canvas.MeasureText(label, 12)
	labelHeight := 15

	// 候補位置のリスト（元の位置、上、下、右にオフセット）
	offsets := []struct{ dx, dy int }{
		{0, 0},      // 元の位置
		{0, -20},    // 上にずらす
		{0, 20},     // 下にずらす
		{30, 0},     // 右にずらす
		{-30, 0},    // 左にずらす
		{30, -15},   // 右上
		{-30, -15},  // 左上
	}

	for _, off := range offsets {
		testX := x + off.dx
		testY := y + off.dy
		if !r.labelOverlapsNodes(testX, testY, labelWidth, labelHeight, nodeBounds) {
			return testX, testY
		}
	}

	// どの位置も重なる場合は元の位置を返す（上にさらにオフセット）
	return x, y - 30
}

// labelOverlapsNodes はラベルがノードと重なるかチェック
func (r *StateRenderer) labelOverlapsNodes(x, y, labelW, labelH int, nodeBounds []stateRect) bool {
	// ラベルの矩形（中央揃えなのでx - labelW/2から開始）
	lx := x - labelW/2
	ly := y - labelH
	lw := labelW
	lh := labelH

	for _, node := range nodeBounds {
		// 矩形の重なりチェック
		if lx < node.x+node.w && lx+lw > node.x &&
			ly < node.y+node.h && ly+lh > node.y {
			return true
		}
	}
	return false
}

// buildTransitionLabel は遷移ラベルを構築する
func (r *StateRenderer) buildTransitionLabel(t state.Transition) string {
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

	guardLabel := ""
	if t.Guard != "" {
		guardLabel = "[" + t.Guard + "]"
	}

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
	return label
}
