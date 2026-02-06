package canvas

import (
	"fmt"

	"pact/internal/infrastructure/renderer/geom"
)

// Diamond はひし形を描画する
func (c *Canvas) Diamond(cx, cy, width, height int, opts ...Option) {
	points := fmt.Sprintf("%d,%d %d,%d %d,%d %d,%d",
		cx, cy-height/2,
		cx+width/2, cy,
		cx, cy+height/2,
		cx-width/2, cy,
	)
	c.Polygon(points, opts...)
}

// Arrow は矢印線を描画する
func (c *Canvas) Arrow(x1, y1, x2, y2 int, opts ...Option) {
	c.Line(x1, y1, x2, y2, opts...)
	// 矢印の先端を描画
	arrowSize := 8
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	length := geom.Sqrt(dx*dx + dy*dy)
	if length == 0 {
		return
	}
	// 単位ベクトル
	ux := dx / length
	uy := dy / length
	// 垂直ベクトル
	px := -uy
	py := ux
	// 矢印の頂点
	ax := float64(x2) - ux*float64(arrowSize)
	ay := float64(y2) - uy*float64(arrowSize)
	// 矢印の両端
	p1x := int(ax + px*float64(arrowSize)/2)
	p1y := int(ay + py*float64(arrowSize)/2)
	p2x := int(ax - px*float64(arrowSize)/2)
	p2y := int(ay - py*float64(arrowSize)/2)
	c.Polygon(fmt.Sprintf("%d,%d %d,%d %d,%d", x2, y2, p1x, p1y, p2x, p2y), Fill("#000"))
}

// OrthogonalArrow は直交ルーティングの矢印を描画する（斜め線なし）
func (c *Canvas) OrthogonalArrow(x1, y1, x2, y2 int, opts ...Option) {
	// 矢印の方向に応じて直交パスを決定
	if x1 == x2 {
		// 垂直方向のみ
		c.Line(x1, y1, x2, y2, opts...)
	} else if y1 == y2 {
		// 水平方向のみ
		c.Line(x1, y1, x2, y2, opts...)
	} else {
		// L字型または逆L字型のルーティング
		// 中間点を計算（水平方向に移動してから垂直方向）
		midY := (y1 + y2) / 2
		c.Line(x1, y1, x1, midY, opts...)
		c.Line(x1, midY, x2, midY, opts...)
		c.Line(x2, midY, x2, y2, opts...)
	}

	// 矢印の先端を描画
	c.DrawArrowHead(x2, y2, x2, y1)
}

// OrthogonalArrowWithMid は中間点を指定した直交ルーティングの矢印を描画する
func (c *Canvas) OrthogonalArrowWithMid(x1, y1, midX, midY, x2, y2 int, opts ...Option) {
	c.Line(x1, y1, midX, y1, opts...)
	c.Line(midX, y1, midX, midY, opts...)
	c.Line(midX, midY, x2, midY, opts...)
	c.Line(x2, midY, x2, y2, opts...)

	// 矢印の先端を描画（最後のセグメントの方向に基づく）
	if y2 > midY {
		c.DrawArrowHead(x2, y2, x2, midY)
	} else if y2 < midY {
		c.DrawArrowHead(x2, y2, x2, midY)
	} else if x2 > midX {
		c.DrawArrowHead(x2, y2, midX, y2)
	} else {
		c.DrawArrowHead(x2, y2, midX, y2)
	}
}

// DrawArrowHead は矢印の先端を描画する
func (c *Canvas) DrawArrowHead(x2, y2, fromX, fromY int) {
	arrowSize := 8
	dx := float64(x2 - fromX)
	dy := float64(y2 - fromY)
	length := geom.Sqrt(dx*dx + dy*dy)
	if length == 0 {
		return
	}
	ux := dx / length
	uy := dy / length
	px := -uy
	py := ux
	ax := float64(x2) - ux*float64(arrowSize)
	ay := float64(y2) - uy*float64(arrowSize)
	p1x := int(ax + px*float64(arrowSize)/2)
	p1y := int(ay + py*float64(arrowSize)/2)
	p2x := int(ax - px*float64(arrowSize)/2)
	p2y := int(ay - py*float64(arrowSize)/2)
	c.Polygon(fmt.Sprintf("%d,%d %d,%d %d,%d", x2, y2, p1x, p1y, p2x, p2y), Fill("#000"))
}

// Stadium は角丸長方形（端子形状）を描画する
func (c *Canvas) Stadium(x, y, width, height int, opts ...Option) {
	radius := height / 2
	c.RoundRect(x, y, width, height, radius, radius, opts...)
}

// Cylinder は円柱（DB形状）を描画する
func (c *Canvas) Cylinder(x, y, width, height int, opts ...Option) {
	ellipseHeight := height / 6
	bodyHeight := height - ellipseHeight

	// 下部楕円
	c.Ellipse(x+width/2, y+height-ellipseHeight/2, width/2, ellipseHeight/2, opts...)
	// 本体
	c.Rect(x, y+ellipseHeight/2, width, bodyHeight, opts...)
	// 上部楕円
	c.Ellipse(x+width/2, y+ellipseHeight/2, width/2, ellipseHeight/2, opts...)
}

// Parallelogram は平行四辺形（IO形状）を描画する
func (c *Canvas) Parallelogram(x, y, width, height, skew int, opts ...Option) {
	points := fmt.Sprintf("%d,%d %d,%d %d,%d %d,%d",
		x+skew, y,
		x+width, y,
		x+width-skew, y+height,
		x, y+height,
	)
	c.Polygon(points, opts...)
}

// Note はノート形状（右上が折れた矩形）を描画する
func (c *Canvas) Note(x, y, width, height int, opts ...Option) {
	foldSize := 10
	if foldSize > width/4 {
		foldSize = width / 4
	}
	if foldSize > height/4 {
		foldSize = height / 4
	}

	// 本体のパス（右上が折れた形）
	path := fmt.Sprintf("M%d,%d L%d,%d L%d,%d L%d,%d L%d,%d L%d,%d Z",
		x, y,
		x+width-foldSize, y,
		x+width, y+foldSize,
		x+width, y+height,
		x, y+height,
		x, y,
	)
	c.Path(path, opts...)

	// 折り返し部分
	foldPath := fmt.Sprintf("M%d,%d L%d,%d L%d,%d",
		x+width-foldSize, y,
		x+width-foldSize, y+foldSize,
		x+width, y+foldSize,
	)
	c.Path(foldPath, Stroke("#000"), Fill("none"))
}
