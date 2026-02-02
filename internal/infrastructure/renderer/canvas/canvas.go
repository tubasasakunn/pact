package canvas

import (
	"bytes"
	"fmt"
	"html"
	"io"
)

// Canvas はSVGを生成するためのキャンバス
type Canvas struct {
	width    int
	height   int
	elements []string
	defs     []string
}

// New は新しいCanvasを作成する
func New() *Canvas {
	return &Canvas{
		width:    800,
		height:   600,
		elements: []string{},
		defs:     []string{},
	}
}

// SetSize はキャンバスサイズを設定する
func (c *Canvas) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// Option は描画オプション
type Option func(attrs map[string]string)

// Fill は塗りつぶし色を設定する
func Fill(color string) Option {
	return func(attrs map[string]string) {
		attrs["fill"] = color
	}
}

// Stroke は線色を設定する
func Stroke(color string) Option {
	return func(attrs map[string]string) {
		attrs["stroke"] = color
	}
}

// StrokeWidth は線幅を設定する
func StrokeWidth(width int) Option {
	return func(attrs map[string]string) {
		attrs["stroke-width"] = fmt.Sprintf("%d", width)
	}
}

// Class はクラス属性を設定する
func Class(className string) Option {
	return func(attrs map[string]string) {
		attrs["class"] = className
	}
}

// StrokeDasharray は破線パターンを設定する
func StrokeDasharray(pattern string) Option {
	return func(attrs map[string]string) {
		attrs["stroke-dasharray"] = pattern
	}
}

// Dashed は破線を設定する（デフォルトパターン: 5,5）
func Dashed() Option {
	return StrokeDasharray("5,5")
}

// TextAnchor はテキストのアンカーを設定する
func TextAnchor(anchor string) Option {
	return func(attrs map[string]string) {
		attrs["text-anchor"] = anchor
	}
}

func applyOptions(attrs map[string]string, opts []Option) {
	for _, opt := range opts {
		opt(attrs)
	}
}

func attrsToString(attrs map[string]string) string {
	var result string
	for k, v := range attrs {
		result += fmt.Sprintf(` %s="%s"`, k, v)
	}
	return result
}

// Rect は矩形を描画する
func (c *Canvas) Rect(x, y, width, height int, opts ...Option) {
	attrs := map[string]string{}
	applyOptions(attrs, opts)
	c.elements = append(c.elements, fmt.Sprintf(
		`<rect x="%d" y="%d" width="%d" height="%d"%s/>`,
		x, y, width, height, attrsToString(attrs),
	))
}

// RoundRect は角丸矩形を描画する
func (c *Canvas) RoundRect(x, y, width, height, rx, ry int, opts ...Option) {
	attrs := map[string]string{}
	applyOptions(attrs, opts)
	c.elements = append(c.elements, fmt.Sprintf(
		`<rect x="%d" y="%d" width="%d" height="%d" rx="%d" ry="%d"%s/>`,
		x, y, width, height, rx, ry, attrsToString(attrs),
	))
}

// Circle は円を描画する
func (c *Canvas) Circle(cx, cy, r int, opts ...Option) {
	attrs := map[string]string{}
	applyOptions(attrs, opts)
	c.elements = append(c.elements, fmt.Sprintf(
		`<circle cx="%d" cy="%d" r="%d"%s/>`,
		cx, cy, r, attrsToString(attrs),
	))
}

// Ellipse は楕円を描画する
func (c *Canvas) Ellipse(cx, cy, rx, ry int, opts ...Option) {
	attrs := map[string]string{}
	applyOptions(attrs, opts)
	c.elements = append(c.elements, fmt.Sprintf(
		`<ellipse cx="%d" cy="%d" rx="%d" ry="%d"%s/>`,
		cx, cy, rx, ry, attrsToString(attrs),
	))
}

// Line は線を描画する
func (c *Canvas) Line(x1, y1, x2, y2 int, opts ...Option) {
	attrs := map[string]string{}
	applyOptions(attrs, opts)
	c.elements = append(c.elements, fmt.Sprintf(
		`<line x1="%d" y1="%d" x2="%d" y2="%d"%s/>`,
		x1, y1, x2, y2, attrsToString(attrs),
	))
}

// Path はパスを描画する
func (c *Canvas) Path(d string, opts ...Option) {
	attrs := map[string]string{}
	applyOptions(attrs, opts)
	c.elements = append(c.elements, fmt.Sprintf(
		`<path d="%s"%s/>`,
		d, attrsToString(attrs),
	))
}

// Polygon は多角形を描画する
func (c *Canvas) Polygon(points string, opts ...Option) {
	attrs := map[string]string{}
	applyOptions(attrs, opts)
	c.elements = append(c.elements, fmt.Sprintf(
		`<polygon points="%s"%s/>`,
		points, attrsToString(attrs),
	))
}

// Text はテキストを描画する
func (c *Canvas) Text(x, y int, text string, opts ...Option) {
	attrs := map[string]string{}
	applyOptions(attrs, opts)
	c.elements = append(c.elements, fmt.Sprintf(
		`<text x="%d" y="%d"%s>%s</text>`,
		x, y, attrsToString(attrs), html.EscapeString(text),
	))
}

// TextWrapped は最大幅で折り返したテキストを描画する
// maxWidth: 最大幅（ピクセル）、lineHeight: 行の高さ
// 戻り値: 描画に使用した行数
func (c *Canvas) TextWrapped(x, y int, text string, maxWidth, lineHeight int, opts ...Option) int {
	fontSize := 12 // デフォルトフォントサイズ
	lines := WrapText(text, maxWidth, fontSize)

	attrs := map[string]string{}
	applyOptions(attrs, opts)

	for i, line := range lines {
		c.elements = append(c.elements, fmt.Sprintf(
			`<text x="%d" y="%d"%s>%s</text>`,
			x, y+i*lineHeight, attrsToString(attrs), html.EscapeString(line),
		))
	}

	return len(lines)
}

// AddDef は定義を追加する
func (c *Canvas) AddDef(def string) {
	c.defs = append(c.defs, def)
}

// WriteTo はSVGを書き出す
func (c *Canvas) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`,
		c.width, c.height, c.width, c.height,
	))
	buf.WriteString("\n")

	if len(c.defs) > 0 {
		buf.WriteString("<defs>\n")
		for _, def := range c.defs {
			buf.WriteString(def)
			buf.WriteString("\n")
		}
		buf.WriteString("</defs>\n")
	}

	for _, elem := range c.elements {
		buf.WriteString(elem)
		buf.WriteString("\n")
	}

	buf.WriteString("</svg>")

	n, err := w.Write(buf.Bytes())
	return int64(n), err
}

// String はSVGを文字列として返す
func (c *Canvas) String() string {
	var buf bytes.Buffer
	c.WriteTo(&buf)
	return buf.String()
}
