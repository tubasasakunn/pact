package svg

import (
	"fmt"

	"pact/internal/domain/diagram/class"
	"pact/internal/infrastructure/renderer/canvas"
	"pact/internal/infrastructure/renderer/geom"
)

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

// renderPath はウェイポイントに沿ってパスを描画
func (r *ClassRenderer) renderPath(c *canvas.Canvas, waypoints []point, opts []canvas.Option) {
	for i := 0; i < len(waypoints)-1; i++ {
		c.Line(waypoints[i].x, waypoints[i].y, waypoints[i+1].x, waypoints[i+1].y, opts...)
	}
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
	length := geom.Sqrt(dx*dx + dy*dy)
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
	length := geom.Sqrt(dx*dx + dy*dy)
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
