package canvas

import "fmt"

// --- スタイル定数 ---

// カラーパレット（モダンで視認性の高い配色）
const (
	// ノード系
	ColorNodeFill   = "#fafbfc" // ノード背景
	ColorNodeStroke = "#2d3748" // ノード枠線
	ColorNodeText   = "#1a202c" // テキスト色

	// セクション区切り
	ColorSectionLine = "#cbd5e0" // セクション区切り線

	// ヘッダー背景（スイムレーンなど）
	ColorHeaderFill = "#edf2f7"

	// ノート
	ColorNoteFill   = "#fefcbf" // 淡い黄色
	ColorNoteStroke = "#d69e2e" // 暗い黄色

	// シーケンス注釈
	ColorReturnFill   = "#c6f6d5" // return（緑系）
	ColorReturnStroke = "#276749"
	ColorThrowFill    = "#fed7d7" // throw（赤系）
	ColorThrowStroke  = "#c53030"
	ColorDefaultNote  = "#fefcbf"

	// アクティベーションバー
	ColorActivationFill = "#e2e8f0"

	// エッジ
	ColorEdge      = "#2d3748"
	ColorEdgeLabel = "#4a5568"

	// 初期状態
	ColorInitialState = "#1a202c"
)

// ストローク幅
const (
	StrokeWidthNode    = 1.5
	StrokeWidthEdge    = 1.0
	StrokeWidthSection = 1.0
)

// Template はSVG <symbol> として定義されるテンプレート
type Template struct {
	ID      string
	ViewBox string // e.g., "0 0 100 100"
	Content string // <symbol> 内部のSVGコンテンツ
}

// TemplateRegistry はテンプレートおよびSVGグローバル定義を管理する
type TemplateRegistry struct {
	templates map[string]*Template
	order     []string // 登録順序を保持
	filters   []string // <filter> 定義
	markers   []string // <marker> 定義
	styles    []string // <style> 定義
}

// NewTemplateRegistry は新しいTemplateRegistryを作成する
func NewTemplateRegistry() *TemplateRegistry {
	return &TemplateRegistry{
		templates: make(map[string]*Template),
	}
}

// Register はテンプレートを登録する
func (r *TemplateRegistry) Register(t *Template) {
	if _, exists := r.templates[t.ID]; !exists {
		r.order = append(r.order, t.ID)
	}
	r.templates[t.ID] = t
}

// Get はテンプレートを取得する
func (r *TemplateRegistry) Get(id string) (*Template, bool) {
	t, ok := r.templates[id]
	return t, ok
}

// AddFilter はSVGフィルター定義を追加する
func (r *TemplateRegistry) AddFilter(filter string) {
	r.filters = append(r.filters, filter)
}

// AddMarker はSVGマーカー定義を追加する
func (r *TemplateRegistry) AddMarker(marker string) {
	r.markers = append(r.markers, marker)
}

// AddStyle はCSS定義を追加する
func (r *TemplateRegistry) AddStyle(style string) {
	r.styles = append(r.styles, style)
}

// ApplyTo はレジストリの全定義をCanvasに追加する
func (r *TemplateRegistry) ApplyTo(c *Canvas) {
	for _, style := range r.styles {
		c.AddDef(style)
	}
	for _, filter := range r.filters {
		c.AddDef(filter)
	}
	for _, marker := range r.markers {
		c.AddDef(marker)
	}
	for _, id := range r.order {
		t := r.templates[id]
		c.AddDef(fmt.Sprintf(
			`<symbol id="%s" viewBox="%s" preserveAspectRatio="none">%s</symbol>`,
			t.ID, t.ViewBox, t.Content,
		))
	}
}

// UseTemplate はテンプレートを使用してシェイプを配置する
func (c *Canvas) UseTemplate(templateID string, x, y, width, height int, opts ...Option) {
	attrs := map[string]string{}
	applyOptions(attrs, opts)
	c.elements = append(c.elements, fmt.Sprintf(
		`<use href="#%s" x="%d" y="%d" width="%d" height="%d"%s/>`,
		templateID, x, y, width, height, attrsToString(attrs),
	))
}

// --- 組み込みテンプレートレジストリ ---

// NewBuiltinRegistry は組み込みテンプレートを持つレジストリを返す
func NewBuiltinRegistry() *TemplateRegistry {
	r := NewTemplateRegistry()

	// ドロップシャドウフィルター
	r.AddFilter(
		`<filter id="drop-shadow" x="-5%" y="-5%" width="115%" height="115%">` +
			`<feGaussianBlur in="SourceAlpha" stdDeviation="2" result="blur"/>` +
			`<feOffset in="blur" dx="1" dy="2" result="offsetBlur"/>` +
			`<feComponentTransfer in="offsetBlur" result="shadow">` +
			`<feFuncA type="linear" slope="0.12"/>` +
			`</feComponentTransfer>` +
			`<feMerge>` +
			`<feMergeNode in="shadow"/>` +
			`<feMergeNode in="SourceGraphic"/>` +
			`</feMerge>` +
			`</filter>`)

	// CSSスタイル（フォントファミリーの統一）
	r.AddStyle(
		`<style>` +
			`text { font-family: -apple-system, "Segoe UI", "Helvetica Neue", Arial, sans-serif; }` +
			`</style>`)

	// --- ノードシンボル ---

	// アクター（人型）: 固定プロポーション
	r.Register(&Template{
		ID:      "actor",
		ViewBox: "0 0 40 70",
		Content: fmt.Sprintf(
			`<circle cx="20" cy="10" r="10" fill="%s" stroke="%s" stroke-width="1.5"/>`+
				`<line x1="20" y1="20" x2="20" y2="45" stroke="%s" stroke-width="1.5"/>`+
				`<line x1="5" y1="30" x2="35" y2="30" stroke="%s" stroke-width="1.5"/>`+
				`<line x1="20" y1="45" x2="5" y2="65" stroke="%s" stroke-width="1.5"/>`+
				`<line x1="20" y1="45" x2="35" y2="65" stroke="%s" stroke-width="1.5"/>`,
			ColorNodeFill, ColorNodeStroke,
			ColorNodeStroke, ColorNodeStroke, ColorNodeStroke, ColorNodeStroke,
		),
	})

	// 初期状態（黒丸）: 固定プロポーション
	r.Register(&Template{
		ID:      "initial-state",
		ViewBox: "0 0 20 20",
		Content: fmt.Sprintf(
			`<circle cx="10" cy="10" r="10" fill="%s"/>`,
			ColorInitialState,
		),
	})

	// 終了状態（二重丸）: 固定プロポーション
	r.Register(&Template{
		ID:      "final-state",
		ViewBox: "0 0 24 24",
		Content: fmt.Sprintf(
			`<circle cx="12" cy="12" r="12" fill="none" stroke="%s" stroke-width="2"/>`+
				`<circle cx="12" cy="12" r="8" fill="%s"/>`,
			ColorNodeStroke, ColorInitialState,
		),
	})

	return r
}

// --- Canvas用追加オプション ---

// Filter はSVGフィルターを適用する
func Filter(id string) Option {
	return func(attrs map[string]string) {
		attrs["filter"] = fmt.Sprintf("url(#%s)", id)
	}
}

// FontWeight はフォントの太さを設定する
func FontWeight(weight string) Option {
	return func(attrs map[string]string) {
		attrs["font-weight"] = weight
	}
}

// FontStyle はフォントスタイルを設定する
func FontStyle(style string) Option {
	return func(attrs map[string]string) {
		attrs["font-style"] = style
	}
}

// FontSize はフォントサイズを設定する
func FontSize(size int) Option {
	return func(attrs map[string]string) {
		attrs["font-size"] = fmt.Sprintf("%d", size)
	}
}
