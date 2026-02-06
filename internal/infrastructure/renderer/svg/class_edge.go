package svg

import (
	"fmt"

	"pact/internal/domain/diagram/class"
	"pact/internal/infrastructure/renderer/canvas"
)

// calculateDistributedEndpoints は複数エッジの接続点を分散配置する
func (r *ClassRenderer) calculateDistributedEndpoints(
	fromPos, toPos struct{ x, y, width, height int },
	outIdx, outTotal, inIdx, inTotal int,
) (fromX, fromY, toX, toY int) {
	fromCenterX := fromPos.x + fromPos.width/2
	fromCenterY := fromPos.y + fromPos.height/2
	toCenterX := toPos.x + toPos.width/2
	toCenterY := toPos.y + toPos.height/2

	// エッジ配置の分散幅（ノード幅の50%を使用、端に寄りすぎないように）
	fromSpread := int(float64(fromPos.width) * 0.5)
	toSpread := int(float64(toPos.width) * 0.5)

	// 垂直方向の差が大きい場合（下向き/上向き接続）
	if abs(toCenterY-fromCenterY) > abs(toCenterX-fromCenterX) {
		// 出力点を下端/上端に分散配置（中央を基準に均等分布）
		fromOffset := 0
		if outTotal > 1 {
			fromOffset = (2*outIdx - (outTotal - 1)) * fromSpread / (2 * (outTotal - 1))
		}
		toOffset := 0
		if inTotal > 1 {
			toOffset = (2*inIdx - (inTotal - 1)) * toSpread / (2 * (inTotal - 1))
		}

		if toCenterY > fromCenterY {
			// 下向き
			fromX = fromCenterX + fromOffset
			fromY = fromPos.y + fromPos.height
			toX = toCenterX + toOffset
			toY = toPos.y
		} else {
			// 上向き
			fromX = fromCenterX + fromOffset
			fromY = fromPos.y
			toX = toCenterX + toOffset
			toY = toPos.y + toPos.height
		}
		return
	}

	// 水平方向の接続
	fromHeightSpread := int(float64(fromPos.height) * 0.4)
	toHeightSpread := int(float64(toPos.height) * 0.4)

	fromYOffset := 0
	if outTotal > 1 {
		fromYOffset = (2*outIdx - (outTotal - 1)) * fromHeightSpread / (2 * (outTotal - 1))
	}
	toYOffset := 0
	if inTotal > 1 {
		toYOffset = (2*inIdx - (inTotal - 1)) * toHeightSpread / (2 * (inTotal - 1))
	}

	if toCenterX > fromCenterX {
		// 右向き
		fromX = fromPos.x + fromPos.width
		fromY = fromCenterY + fromYOffset
		toX = toPos.x
		toY = toCenterY + toYOffset
	} else {
		// 左向き
		fromX = fromPos.x
		fromY = fromCenterY + fromYOffset
		toX = toPos.x + toPos.width
		toY = toCenterY + toYOffset
	}
	return
}

// renderEdgeImproved は改良されたエッジ描画
func (r *ClassRenderer) renderEdgeImproved(c *canvas.Canvas, edge class.Edge, x1, y1, x2, y2 int, nodePositions map[string]struct{ x, y, width, height int }) {
	opts := []canvas.Option{canvas.Stroke(canvas.ColorEdge)}
	if edge.LineStyle == class.LineStyleDashed {
		opts = append(opts, canvas.Dashed())
	}

	// 障害物リストを構築（始点・終点ノード以外）
	var obstacles []rect
	for nodeID, pos := range nodePositions {
		if nodeID == edge.From || nodeID == edge.To {
			continue
		}
		obstacles = append(obstacles, rect{pos.x, pos.y, pos.width, pos.height})
	}

	// ウェイポイントを計算
	waypoints := r.calculateWaypoints(x1, y1, x2, y2, obstacles)

	// パスを描画
	r.renderPath(c, waypoints, opts)

	// 矢印を描画
	if len(waypoints) >= 2 {
		if edge.Decoration == class.DecorationFilledDiamond || edge.Decoration == class.DecorationEmptyDiamond {
			// ダイヤモンド装飾はFrom（始点）ノード側に描画
			r.drawArrowHead(c, edge, waypoints[0].x, waypoints[0].y, waypoints[1].x, waypoints[1].y)
		} else {
			// 三角形・矢印装飾はTo（終点）ノード側に描画（最後のセグメントの方向で）
			lastIdx := len(waypoints) - 1
			r.drawArrowHead(c, edge, waypoints[lastIdx-1].x, waypoints[lastIdx-1].y, waypoints[lastIdx].x, waypoints[lastIdx].y)
		}
	}

	// ラベル描画
	if edge.Label != "" {
		midX := (x1 + x2) / 2
		midY := (y1 + y2) / 2
		c.Text(midX, midY-5, edge.Label,
			canvas.TextAnchor("middle"),
			canvas.Fill(canvas.ColorEdgeLabel),
		)
	}
}

// calculateWaypoints は障害物を回避するウェイポイントを計算（常に直交ルーティング）
func (r *ClassRenderer) calculateWaypoints(x1, y1, x2, y2 int, obstacles []rect) []point {
	start := point{x1, y1}
	end := point{x2, y2}

	// 完全に水平または垂直の場合のみ直線を許可
	if x1 == x2 || y1 == y2 {
		if !r.pathIntersectsAnyObstacle(start, end, obstacles) {
			return []point{start, end}
		}
	}

	// 直交ルーティングが必要
	margin := 20 // 障害物からのマージン

	// ルーティング戦略を選択
	// 1. L字型（1回曲がり）
	// 2. Z字型（2回曲がり）
	// 3. U字型（障害物を迂回）

	// L字型はスキップし、Z字型ルーティングで最初と最後のセグメントが
	// 接続辺に垂直になるようにする（矢印が辺に並行にならない）

	// Z字型（両方向を試す：V-H-V と H-V-H）
	// 主方向を先に試し、ダメならもう一方の向きも試す
	tryVHV := func() []point {
		// 垂直-水平-垂直
		zCandidates := []int{
			(y1 + y2) / 2,
			y1 + (y2-y1)/4,
			y1 + (y2-y1)*3/4,
			y1 + (y2-y1)/3,
			y1 + (y2-y1)*2/3,
		}
		for _, midY := range zCandidates {
			mid1 := point{x1, midY}
			mid2 := point{x2, midY}
			if !r.pathIntersectsAnyObstacle(start, mid1, obstacles) &&
				!r.pathIntersectsAnyObstacle(mid1, mid2, obstacles) &&
				!r.pathIntersectsAnyObstacle(mid2, end, obstacles) {
				return []point{start, mid1, mid2, end}
			}
		}
		return nil
	}

	tryHVH := func() []point {
		// 水平-垂直-水平
		zCandidates := []int{
			(x1 + x2) / 2,
			x1 + (x2-x1)/4,
			x1 + (x2-x1)*3/4,
			x1 + (x2-x1)/3,
			x1 + (x2-x1)*2/3,
		}
		for _, midX := range zCandidates {
			mid1 := point{midX, y1}
			mid2 := point{midX, y2}
			if !r.pathIntersectsAnyObstacle(start, mid1, obstacles) &&
				!r.pathIntersectsAnyObstacle(mid1, mid2, obstacles) &&
				!r.pathIntersectsAnyObstacle(mid2, end, obstacles) {
				return []point{start, mid1, mid2, end}
			}
		}
		return nil
	}

	// 主方向を先に試し、ダメなら逆方向も試す
	if abs(y2-y1) > abs(x2-x1) {
		if result := tryVHV(); result != nil {
			return result
		}
		if result := tryHVH(); result != nil {
			return result
		}
	} else {
		if result := tryHVH(); result != nil {
			return result
		}
		if result := tryVHV(); result != nil {
			return result
		}
	}

	// U字型ルーティング（障害物の外側を通る）
	// 全障害物のバウンディングボックスを計算
	if len(obstacles) > 0 {
		minObsX, minObsY := obstacles[0].x, obstacles[0].y
		maxObsX, maxObsY := obstacles[0].x+obstacles[0].w, obstacles[0].y+obstacles[0].h

		for _, obs := range obstacles {
			if obs.x < minObsX {
				minObsX = obs.x
			}
			if obs.y < minObsY {
				minObsY = obs.y
			}
			if obs.x+obs.w > maxObsX {
				maxObsX = obs.x + obs.w
			}
			if obs.y+obs.h > maxObsY {
				maxObsY = obs.y + obs.h
			}
		}

		// 上を通るルート
		topY := minObsY - margin
		mid1 := point{x1, topY}
		mid2 := point{x2, topY}
		if !r.pathIntersectsAnyObstacle(start, mid1, obstacles) &&
			!r.pathIntersectsAnyObstacle(mid1, mid2, obstacles) &&
			!r.pathIntersectsAnyObstacle(mid2, end, obstacles) {
			return []point{start, mid1, mid2, end}
		}

		// 下を通るルート
		bottomY := maxObsY + margin
		mid1 = point{x1, bottomY}
		mid2 = point{x2, bottomY}
		if !r.pathIntersectsAnyObstacle(start, mid1, obstacles) &&
			!r.pathIntersectsAnyObstacle(mid1, mid2, obstacles) &&
			!r.pathIntersectsAnyObstacle(mid2, end, obstacles) {
			return []point{start, mid1, mid2, end}
		}

		// 左を通るルート
		leftX := minObsX - margin
		mid1 = point{leftX, y1}
		mid2 = point{leftX, y2}
		if !r.pathIntersectsAnyObstacle(start, mid1, obstacles) &&
			!r.pathIntersectsAnyObstacle(mid1, mid2, obstacles) &&
			!r.pathIntersectsAnyObstacle(mid2, end, obstacles) {
			return []point{start, mid1, mid2, end}
		}

		// 右を通るルート
		rightX := maxObsX + margin
		mid1 = point{rightX, y1}
		mid2 = point{rightX, y2}
		if !r.pathIntersectsAnyObstacle(start, mid1, obstacles) &&
			!r.pathIntersectsAnyObstacle(mid1, mid2, obstacles) &&
			!r.pathIntersectsAnyObstacle(mid2, end, obstacles) {
			return []point{start, mid1, mid2, end}
		}

		// 5セグメントルーティング: 上/下に出て障害物を回避して戻る
		for _, bypassY := range []int{topY, bottomY} {
			for _, bypassX := range []int{leftX, rightX} {
				wp := []point{
					start,
					{x1, bypassY},
					{bypassX, bypassY},
					{bypassX, y2},
					end,
				}
				allClear := true
				for i := 0; i < len(wp)-1; i++ {
					if r.pathIntersectsAnyObstacle(wp[i], wp[i+1], obstacles) {
						allClear = false
						break
					}
				}
				if allClear {
					return wp
				}
			}
		}
	}

	// フォールバック: 単純なZ字型
	midY := (y1 + y2) / 2
	return []point{start, {x1, midY}, {x2, midY}, end}
}

// pathIntersectsAnyObstacle はパスが障害物と交差するか確認
func (r *ClassRenderer) pathIntersectsAnyObstacle(p1, p2 point, obstacles []rect) bool {
	for _, obs := range obstacles {
		if r.lineIntersectsRect(p1.x, p1.y, p2.x, p2.y, obs.x, obs.y, obs.w, obs.h) {
			return true
		}
	}
	return false
}

// renderPath はウェイポイントに沿ってパスを描画
func (r *ClassRenderer) renderPath(c *canvas.Canvas, waypoints []point, opts []canvas.Option) {
	for i := 0; i < len(waypoints)-1; i++ {
		c.Line(waypoints[i].x, waypoints[i].y, waypoints[i+1].x, waypoints[i+1].y, opts...)
	}
}

// lineIntersectsRect は直線（直交セグメント）が矩形と交差するかチェック
func (r *ClassRenderer) lineIntersectsRect(x1, y1, x2, y2, rx, ry, rw, rh int) bool {
	// 矩形の外側にマージンを設けて判定（線が矩形に近づきすぎないように）
	margin := 5
	rectLeft := rx - margin
	rectRight := rx + rw + margin
	rectTop := ry - margin
	rectBottom := ry + rh + margin

	lminX := minInt(x1, x2)
	lmaxX := maxInt(x1, x2)
	lminY := minInt(y1, y2)
	lmaxY := maxInt(y1, y2)

	// バウンディングボックスが重なっていなければ交差しない
	if lmaxX < rectLeft || lminX > rectRight || lmaxY < rectTop || lminY > rectBottom {
		return false
	}

	// 水平線セグメントの場合
	if y1 == y2 {
		return y1 >= rectTop && y1 <= rectBottom &&
			lmaxX >= rectLeft && lminX <= rectRight
	}

	// 垂直線セグメントの場合
	if x1 == x2 {
		return x1 >= rectLeft && x1 <= rectRight &&
			lmaxY >= rectTop && lminY <= rectBottom
	}

	// 斜め線の場合（直交ルーティングでは通常発生しないが安全のため）
	// 線分の任意の点が矩形内にあるかチェック
	// 端点チェック
	if pointInRect(x1, y1, rectLeft, rectTop, rectRight, rectBottom) ||
		pointInRect(x2, y2, rectLeft, rectTop, rectRight, rectBottom) {
		return true
	}

	// 矩形の4辺と線分の交差チェック
	if segmentsIntersect(x1, y1, x2, y2, rectLeft, rectTop, rectRight, rectTop) ||
		segmentsIntersect(x1, y1, x2, y2, rectRight, rectTop, rectRight, rectBottom) ||
		segmentsIntersect(x1, y1, x2, y2, rectLeft, rectBottom, rectRight, rectBottom) ||
		segmentsIntersect(x1, y1, x2, y2, rectLeft, rectTop, rectLeft, rectBottom) {
		return true
	}

	return false
}

// renderVerticalEdge は継承・実装エッジを垂直接続で描画する
// 子のtop → 親のbottomを常に垂直に接続し、最後のセグメントが辺に垂直になるようにする
func (r *ClassRenderer) renderVerticalEdge(c *canvas.Canvas, edge class.Edge,
	fromPos, toPos struct{ x, y, width, height int },
	outIdx, outTotal, inIdx, inTotal int,
	nodePositions map[string]struct{ x, y, width, height int }) {

	fromCenterX := fromPos.x + fromPos.width/2
	toCenterX := toPos.x + toPos.width/2

	// 接続点を分散配置（中央を基準に均等分布、端に寄りすぎないように）
	toSpread := int(float64(toPos.width) * 0.5)
	toOffset := 0
	if inTotal > 1 {
		toOffset = (2*inIdx - (inTotal - 1)) * toSpread / (2 * (inTotal - 1))
	}

	// 子のtop、親のbottomに接続
	fromX := fromCenterX
	fromY := fromPos.y
	toX := toCenterX + toOffset
	toY := toPos.y + toPos.height

	opts := []canvas.Option{canvas.Stroke(canvas.ColorEdge)}
	if edge.LineStyle == class.LineStyleDashed {
		opts = append(opts, canvas.Dashed())
	}

	// 障害物リストを構築（始点・終点ノード以外）
	var obstacles []rect
	for nodeID, pos := range nodePositions {
		if nodeID == edge.From || nodeID == edge.To {
			continue
		}
		obstacles = append(obstacles, rect{pos.x, pos.y, pos.width, pos.height})
	}

	if fromX == toX {
		// 真っ直ぐ垂直: 障害物チェック付き
		if !r.pathIntersectsAnyObstacle(point{fromX, fromY}, point{toX, toY}, obstacles) {
			c.Line(fromX, fromY, toX, toY, opts...)
			r.drawArrowHead(c, edge, fromX, fromY, toX, toY)
		} else {
			// 障害物を回避するためウェイポイント計算
			waypoints := r.calculateWaypoints(fromX, fromY, toX, toY, obstacles)
			r.renderPath(c, waypoints, opts)
			if len(waypoints) >= 2 {
				lastIdx := len(waypoints) - 1
				r.drawArrowHead(c, edge, waypoints[lastIdx-1].x, waypoints[lastIdx-1].y, waypoints[lastIdx].x, waypoints[lastIdx].y)
			}
		}
	} else {
		// Z字型ルーティング: 障害物を回避する中間Y座標を探す
		midY := r.findSafeVerticalMidY(fromX, fromY, toX, toY, obstacles)
		c.Line(fromX, fromY, fromX, midY, opts...)
		c.Line(fromX, midY, toX, midY, opts...)
		c.Line(toX, midY, toX, toY, opts...)
		r.drawArrowHead(c, edge, toX, midY, toX, toY)
	}

	// ラベル描画
	if edge.Label != "" {
		midX := (fromX + toX) / 2
		midY := (fromY + toY) / 2
		c.Text(midX, midY-5, edge.Label,
			canvas.TextAnchor("middle"),
			canvas.Fill(canvas.ColorEdgeLabel),
		)
	}
}

// findSafeVerticalMidY はZ字型ルーティングで障害物を回避する中間Y座標を見つける
func (r *ClassRenderer) findSafeVerticalMidY(fromX, fromY, toX, toY int, obstacles []rect) int {
	// デフォルトの中間Y（2点の中間）
	midY := (fromY + toY) / 2

	// 3セグメント全て（垂直→水平→垂直）が障害物を回避できるかチェック
	if !r.zShapeIntersectsObstacles(fromX, fromY, toX, toY, midY, obstacles) {
		return midY
	}

	// 中間Yをずらして障害物を回避する候補を試す
	// fromY（上）とtoY（下）の間で複数のY座標を試す
	minY := minInt(fromY, toY)
	maxY := maxInt(fromY, toY)
	margin := 15

	// 上端寄り・下端寄りを含む複数候補を試す
	candidates := []int{
		minY + margin,           // 上端近く
		maxY - margin,           // 下端近く
		minY + (maxY-minY)/4,   // 1/4
		minY + (maxY-minY)*3/4, // 3/4
		minY + (maxY-minY)/3,   // 1/3
		minY + (maxY-minY)*2/3, // 2/3
	}

	for _, candidateY := range candidates {
		if candidateY <= minY || candidateY >= maxY {
			continue
		}
		if !r.zShapeIntersectsObstacles(fromX, fromY, toX, toY, candidateY, obstacles) {
			return candidateY
		}
	}

	// 外側を通るルート
	if len(obstacles) > 0 {
		minObsY := obstacles[0].y
		maxObsY := obstacles[0].y + obstacles[0].h
		for _, obs := range obstacles {
			if obs.y < minObsY {
				minObsY = obs.y
			}
			if obs.y+obs.h > maxObsY {
				maxObsY = obs.y + obs.h
			}
		}
		// 障害物の上を通る
		topY := minObsY - 20
		if !r.zShapeIntersectsObstacles(fromX, fromY, toX, toY, topY, obstacles) {
			return topY
		}
		// 障害物の下を通る
		bottomY := maxObsY + 20
		if !r.zShapeIntersectsObstacles(fromX, fromY, toX, toY, bottomY, obstacles) {
			return bottomY
		}
	}

	// フォールバック: デフォルトの中間Y
	return midY
}

// zShapeIntersectsObstacles はZ字型パス（垂直→水平→垂直）が障害物と交差するかチェック
func (r *ClassRenderer) zShapeIntersectsObstacles(fromX, fromY, toX, toY, midY int, obstacles []rect) bool {
	seg1Start := point{fromX, fromY}
	seg1End := point{fromX, midY}
	seg2Start := point{fromX, midY}
	seg2End := point{toX, midY}
	seg3Start := point{toX, midY}
	seg3End := point{toX, toY}

	return r.pathIntersectsAnyObstacle(seg1Start, seg1End, obstacles) ||
		r.pathIntersectsAnyObstacle(seg2Start, seg2End, obstacles) ||
		r.pathIntersectsAnyObstacle(seg3Start, seg3End, obstacles)
}

// drawArrowHead はエッジの装飾（矢印先端）を描画
func (r *ClassRenderer) drawArrowHead(c *canvas.Canvas, edge class.Edge, fromX, fromY, toX, toY int) {
	switch edge.Decoration {
	case class.DecorationTriangle:
		c.Polygon(trianglePoints(toX, toY, fromX, fromY), canvas.Fill(canvas.ColorNodeFill), canvas.Stroke(canvas.ColorEdge))
	case class.DecorationFilledDiamond:
		c.Polygon(diamondPoints(fromX, fromY, toX, toY), canvas.Fill(canvas.ColorEdge))
	case class.DecorationEmptyDiamond:
		c.Polygon(diamondPoints(fromX, fromY, toX, toY), canvas.Fill(canvas.ColorNodeFill), canvas.Stroke(canvas.ColorEdge))
	default:
		c.Polygon(trianglePoints(toX, toY, fromX, fromY), canvas.Fill(canvas.ColorEdge))
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
