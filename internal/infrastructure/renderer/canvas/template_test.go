package canvas

import (
	"strings"
	"testing"
)

// =============================================================================
// RT010-RT030: Template System Tests
// =============================================================================

// RT010: TemplateRegistry の基本操作
func TestTemplateRegistry_RegisterAndGet(t *testing.T) {
	r := NewTemplateRegistry()

	tmpl := &Template{
		ID:      "test-shape",
		ViewBox: "0 0 100 100",
		Content: `<rect x="0" y="0" width="100" height="100"/>`,
	}
	r.Register(tmpl)

	got, ok := r.Get("test-shape")
	if !ok {
		t.Fatal("expected template to be registered")
	}
	if got.ID != "test-shape" {
		t.Errorf("expected ID 'test-shape', got '%s'", got.ID)
	}
}

// RT011: 未登録テンプレートの取得
func TestTemplateRegistry_GetNotFound(t *testing.T) {
	r := NewTemplateRegistry()

	_, ok := r.Get("nonexistent")
	if ok {
		t.Error("expected template not found")
	}
}

// RT012: 同一IDの上書き
func TestTemplateRegistry_OverwriteSameID(t *testing.T) {
	r := NewTemplateRegistry()

	r.Register(&Template{ID: "shape", ViewBox: "0 0 100 100", Content: "original"})
	r.Register(&Template{ID: "shape", ViewBox: "0 0 100 100", Content: "updated"})

	got, _ := r.Get("shape")
	if got.Content != "updated" {
		t.Errorf("expected content 'updated', got '%s'", got.Content)
	}
}

// RT013: ApplyTo でCanvasにdefs追加
func TestTemplateRegistry_ApplyTo(t *testing.T) {
	r := NewTemplateRegistry()
	r.AddFilter(`<filter id="f1"><feGaussianBlur/></filter>`)
	r.AddMarker(`<marker id="m1"><path d="M0,0"/></marker>`)
	r.AddStyle(`<style>.test{fill:red}</style>`)
	r.Register(&Template{
		ID:      "my-shape",
		ViewBox: "0 0 50 50",
		Content: `<circle cx="25" cy="25" r="25"/>`,
	})

	c := New()
	r.ApplyTo(c)

	svg := c.String()
	if !strings.Contains(svg, "<defs>") {
		t.Error("expected <defs> in SVG output")
	}
	if !strings.Contains(svg, `<filter id="f1"`) {
		t.Error("expected filter definition in SVG")
	}
	if !strings.Contains(svg, `<marker id="m1"`) {
		t.Error("expected marker definition in SVG")
	}
	if !strings.Contains(svg, `<style>`) {
		t.Error("expected style definition in SVG")
	}
	if !strings.Contains(svg, `<symbol id="my-shape"`) {
		t.Error("expected symbol definition in SVG")
	}
	if !strings.Contains(svg, `preserveAspectRatio="none"`) {
		t.Error("expected preserveAspectRatio on symbol")
	}
}

// RT014: UseTemplate で <use> 要素を生成
func TestCanvas_UseTemplate(t *testing.T) {
	c := New()
	c.UseTemplate("test-shape", 10, 20, 100, 50)

	svg := c.String()
	if !strings.Contains(svg, `<use href="#test-shape"`) {
		t.Error("expected <use> element with href")
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

// RT015: UseTemplate にオプション適用
func TestCanvas_UseTemplate_WithOptions(t *testing.T) {
	c := New()
	c.UseTemplate("test-shape", 0, 0, 100, 100, Class("highlight"))

	svg := c.String()
	if !strings.Contains(svg, `class="highlight"`) {
		t.Error("expected class attribute on <use>")
	}
}

// RT016: Filter オプション
func TestOption_Filter(t *testing.T) {
	c := New()
	c.Rect(0, 0, 100, 100, Filter("drop-shadow"))

	svg := c.String()
	if !strings.Contains(svg, `filter="url(#drop-shadow)"`) {
		t.Error("expected filter attribute")
	}
}

// RT017: FontWeight オプション
func TestOption_FontWeight(t *testing.T) {
	c := New()
	c.Text(0, 0, "bold text", FontWeight("bold"))

	svg := c.String()
	if !strings.Contains(svg, `font-weight="bold"`) {
		t.Error("expected font-weight attribute")
	}
}

// RT018: FontStyle オプション
func TestOption_FontStyle(t *testing.T) {
	c := New()
	c.Text(0, 0, "italic text", FontStyle("italic"))

	svg := c.String()
	if !strings.Contains(svg, `font-style="italic"`) {
		t.Error("expected font-style attribute")
	}
}

// RT019: FontSize オプション
func TestOption_FontSize(t *testing.T) {
	c := New()
	c.Text(0, 0, "sized text", FontSize(14))

	svg := c.String()
	if !strings.Contains(svg, `font-size="14"`) {
		t.Error("expected font-size attribute")
	}
}

// RT020: 組み込みレジストリの基本テスト
func TestNewBuiltinRegistry(t *testing.T) {
	r := NewBuiltinRegistry()

	// 組み込みテンプレートの存在確認
	templates := []string{"actor", "initial-state", "final-state"}
	for _, id := range templates {
		if _, ok := r.Get(id); !ok {
			t.Errorf("expected built-in template '%s' to exist", id)
		}
	}
}

// RT021: 組み込みレジストリ適用でフィルターとスタイルが含まれる
func TestBuiltinRegistry_ApplyToCanvas(t *testing.T) {
	r := NewBuiltinRegistry()
	c := New()
	r.ApplyTo(c)

	svg := c.String()
	if !strings.Contains(svg, `<filter id="drop-shadow"`) {
		t.Error("expected drop-shadow filter in built-in registry")
	}
	if !strings.Contains(svg, `<style>`) {
		t.Error("expected style definition in built-in registry")
	}
	if !strings.Contains(svg, `<symbol id="actor"`) {
		t.Error("expected actor symbol in built-in registry")
	}
}

// RT022: カラー定数が空でないこと
func TestColorConstants(t *testing.T) {
	colors := []struct {
		name  string
		value string
	}{
		{"ColorNodeFill", ColorNodeFill},
		{"ColorNodeStroke", ColorNodeStroke},
		{"ColorNodeText", ColorNodeText},
		{"ColorSectionLine", ColorSectionLine},
		{"ColorNoteFill", ColorNoteFill},
		{"ColorNoteStroke", ColorNoteStroke},
		{"ColorEdge", ColorEdge},
	}

	for _, c := range colors {
		if c.value == "" {
			t.Errorf("color constant %s should not be empty", c.name)
		}
		if c.value[0] != '#' {
			t.Errorf("color constant %s should start with '#', got '%s'", c.name, c.value)
		}
	}
}

// RT023: ApplyTo の定義順序（style → filter → marker → symbol）
func TestTemplateRegistry_ApplyToOrder(t *testing.T) {
	r := NewTemplateRegistry()
	r.AddStyle(`<style>.a{}</style>`)
	r.AddFilter(`<filter id="f"/>`)
	r.AddMarker(`<marker id="m"/>`)
	r.Register(&Template{ID: "s", ViewBox: "0 0 10 10", Content: ""})

	c := New()
	r.ApplyTo(c)

	svg := c.String()
	styleIdx := strings.Index(svg, "<style>")
	filterIdx := strings.Index(svg, "<filter")
	markerIdx := strings.Index(svg, "<marker")
	symbolIdx := strings.Index(svg, "<symbol")

	if styleIdx > filterIdx {
		t.Error("expected style before filter")
	}
	if filterIdx > markerIdx {
		t.Error("expected filter before marker")
	}
	if markerIdx > symbolIdx {
		t.Error("expected marker before symbol")
	}
}
