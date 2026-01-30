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
	// 矢印の先端を描画
	arrowSize := 8
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	length := sqrt(dx*dx + dy*dy)
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
