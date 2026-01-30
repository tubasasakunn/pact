package canvas

import (
	"bytes"
	"strings"
	"testing"
)

// =============================================================================
// RC001-RC017: Canvas Tests
// =============================================================================

// RC001: 空キャンバス
func TestCanvas_Empty(t *testing.T) {
	c := New()
	svg := c.String()

	if !strings.Contains(svg, "<svg") {
		t.Error("expected svg element")
	}
	if !strings.Contains(svg, "</svg>") {
		t.Error("expected closing svg tag")
	}
}

// RC002: サイズ設定
func TestCanvas_SetSize(t *testing.T) {
	c := New()
	c.SetSize(800, 600)
	svg := c.String()

	if !strings.Contains(svg, `viewBox="0 0 800 600"`) {
		t.Error("expected viewBox with specified size")
	}
}

// RC003: 矩形
func TestCanvas_Rect(t *testing.T) {
	c := New()
	c.Rect(10, 20, 100, 50)
	svg := c.String()

	if !strings.Contains(svg, "<rect") {
		t.Error("expected rect element")
	}
	if !strings.Contains(svg, `x="10"`) {
		t.Error("expected x attribute")
	}
	if !strings.Contains(svg, `y="20"`) {
		t.Error("expected y attribute")
	}
	if !strings.Contains(svg, `width="100"`) {
		t.Error("expected width attribute")
	}
	if !strings.Contains(svg, `height="50"`) {
		t.Error("expected height attribute")
	}
}

// RC004: 角丸矩形
func TestCanvas_RoundRect(t *testing.T) {
	c := New()
	c.RoundRect(10, 20, 100, 50, 5, 5)
	svg := c.String()

	if !strings.Contains(svg, `rx="5"`) {
		t.Error("expected rx attribute")
	}
	if !strings.Contains(svg, `ry="5"`) {
		t.Error("expected ry attribute")
	}
}

// RC005: 円
func TestCanvas_Circle(t *testing.T) {
	c := New()
	c.Circle(50, 50, 25)
	svg := c.String()

	if !strings.Contains(svg, "<circle") {
		t.Error("expected circle element")
	}
	if !strings.Contains(svg, `cx="50"`) {
		t.Error("expected cx attribute")
	}
	if !strings.Contains(svg, `cy="50"`) {
		t.Error("expected cy attribute")
	}
	if !strings.Contains(svg, `r="25"`) {
		t.Error("expected r attribute")
	}
}

// RC006: 楕円
func TestCanvas_Ellipse(t *testing.T) {
	c := New()
	c.Ellipse(50, 50, 30, 20)
	svg := c.String()

	if !strings.Contains(svg, "<ellipse") {
		t.Error("expected ellipse element")
	}
	if !strings.Contains(svg, `rx="30"`) {
		t.Error("expected rx attribute")
	}
	if !strings.Contains(svg, `ry="20"`) {
		t.Error("expected ry attribute")
	}
}

// RC007: 線
func TestCanvas_Line(t *testing.T) {
	c := New()
	c.Line(0, 0, 100, 100)
	svg := c.String()

	if !strings.Contains(svg, "<line") {
		t.Error("expected line element")
	}
	if !strings.Contains(svg, `x1="0"`) {
		t.Error("expected x1 attribute")
	}
	if !strings.Contains(svg, `y1="0"`) {
		t.Error("expected y1 attribute")
	}
	if !strings.Contains(svg, `x2="100"`) {
		t.Error("expected x2 attribute")
	}
	if !strings.Contains(svg, `y2="100"`) {
		t.Error("expected y2 attribute")
	}
}

// RC008: パス
func TestCanvas_Path(t *testing.T) {
	c := New()
	c.Path("M0,0 L100,100")
	svg := c.String()

	if !strings.Contains(svg, "<path") {
		t.Error("expected path element")
	}
	if !strings.Contains(svg, `d="M0,0 L100,100"`) {
		t.Error("expected d attribute")
	}
}

// RC009: 多角形
func TestCanvas_Polygon(t *testing.T) {
	c := New()
	c.Polygon("0,0 50,25 0,50")
	svg := c.String()

	if !strings.Contains(svg, "<polygon") {
		t.Error("expected polygon element")
	}
	if !strings.Contains(svg, `points="0,0 50,25 0,50"`) {
		t.Error("expected points attribute")
	}
}

// RC010: テキスト
func TestCanvas_Text(t *testing.T) {
	c := New()
	c.Text(10, 20, "hello")
	svg := c.String()

	if !strings.Contains(svg, "<text") {
		t.Error("expected text element")
	}
	if !strings.Contains(svg, ">hello</text>") {
		t.Error("expected text content")
	}
}

// RC011: 塗りオプション
func TestCanvas_Option_Fill(t *testing.T) {
	c := New()
	c.Rect(0, 0, 100, 100, Fill("#ff0000"))
	svg := c.String()

	if !strings.Contains(svg, `fill="#ff0000"`) {
		t.Error("expected fill attribute")
	}
}

// RC012: 線色オプション
func TestCanvas_Option_Stroke(t *testing.T) {
	c := New()
	c.Rect(0, 0, 100, 100, Stroke("#000000"))
	svg := c.String()

	if !strings.Contains(svg, `stroke="#000000"`) {
		t.Error("expected stroke attribute")
	}
}

// RC013: 線幅オプション
func TestCanvas_Option_StrokeWidth(t *testing.T) {
	c := New()
	c.Rect(0, 0, 100, 100, StrokeWidth(2))
	svg := c.String()

	if !strings.Contains(svg, `stroke-width="2"`) {
		t.Error("expected stroke-width attribute")
	}
}

// RC014: クラスオプション
func TestCanvas_Option_Class(t *testing.T) {
	c := New()
	c.Rect(0, 0, 100, 100, Class("node"))
	svg := c.String()

	if !strings.Contains(svg, `class="node"`) {
		t.Error("expected class attribute")
	}
}

// RC015: 複数オプション
func TestCanvas_Option_Multiple(t *testing.T) {
	c := New()
	c.Rect(0, 0, 100, 100, Fill("#fff"), Stroke("#000"), StrokeWidth(2))
	svg := c.String()

	if !strings.Contains(svg, `fill="#fff"`) {
		t.Error("expected fill attribute")
	}
	if !strings.Contains(svg, `stroke="#000"`) {
		t.Error("expected stroke attribute")
	}
	if !strings.Contains(svg, `stroke-width="2"`) {
		t.Error("expected stroke-width attribute")
	}
}

// RC016: 定義
func TestCanvas_Defs(t *testing.T) {
	c := New()
	c.AddDef(`<marker id="arrow"/>`)
	svg := c.String()

	if !strings.Contains(svg, "<defs>") {
		t.Error("expected defs element")
	}
	if !strings.Contains(svg, `<marker id="arrow"/>`) {
		t.Error("expected marker in defs")
	}
}

// RC017: 出力
func TestCanvas_WriteTo(t *testing.T) {
	c := New()
	c.Rect(0, 0, 100, 100)

	var buf bytes.Buffer
	n, err := c.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == 0 {
		t.Error("expected non-zero bytes written")
	}

	svg := buf.String()
	if !strings.Contains(svg, "<svg") {
		t.Error("expected valid SVG output")
	}
}
