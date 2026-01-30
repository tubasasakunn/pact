package canvas

import "fmt"

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
	// 矢印の先端を描画（簡易版）
	arrowSize := 8
	dx := x2 - x1
	dy := y2 - y1
	// 簡易的な矢印の描画
	if dx != 0 || dy != 0 {
		c.Path(fmt.Sprintf("M%d,%d L%d,%d L%d,%d Z",
			x2, y2,
			x2-arrowSize, y2-arrowSize/2,
			x2-arrowSize, y2+arrowSize/2,
		), opts...)
	}
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
