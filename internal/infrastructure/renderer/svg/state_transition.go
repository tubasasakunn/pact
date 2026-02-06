package svg

import (
	"fmt"

	"pact/internal/domain/diagram/state"
	"pact/internal/infrastructure/renderer/canvas"
	"pact/internal/infrastructure/renderer/geom"
)

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
		c.Line(startX, startY, startX+30, startY, canvas.Stroke(canvas.ColorEdge))
		c.Line(startX+30, startY, startX+30, startY-30, canvas.Stroke(canvas.ColorEdge))
		c.Line(startX+30, startY-30, endX, startY-30, canvas.Stroke(canvas.ColorEdge))
		c.Line(endX, startY-30, endX, endY, canvas.Stroke(canvas.ColorEdge))
		c.DrawArrowHead(endX, endY, endX, startY-30)
		midX = startX + 15
		midY = startY - 15
	} else if geom.Abs(dx) > geom.Abs(dy) {
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
		c.Text(labelX, labelY, label,
			canvas.TextAnchor("middle"),
			canvas.Fill(canvas.ColorEdgeLabel),
		)
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
		c.Line(waypoints[i].x, waypoints[i].y, waypoints[i+1].x, waypoints[i+1].y, canvas.Stroke(canvas.ColorEdge))
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
		{0, 0},     // 元の位置
		{0, -20},   // 上にずらす
		{0, 20},    // 下にずらす
		{30, 0},    // 右にずらす
		{-30, 0},   // 左にずらす
		{30, -15},  // 右上
		{-30, -15}, // 左上
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
