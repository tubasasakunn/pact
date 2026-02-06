// Package canvas provides beautiful SVG decoration templates.
package canvas

// SVG gradient and filter definitions for pattern decorations

// GradientDef represents an SVG gradient definition
type GradientDef struct {
	ID     string
	Type   string // "linear" or "radial"
	Stops  []GradientStop
	Attrs  map[string]string // Additional attributes (x1, y1, x2, y2, etc.)
}

// GradientStop represents a gradient color stop
type GradientStop struct {
	Offset  string
	Color   string
	Opacity string
}

// FilterDef represents an SVG filter definition
type FilterDef struct {
	ID      string
	Content string // Raw SVG filter content
}

// MarkerDef represents an SVG marker definition
type MarkerDef struct {
	ID       string
	ViewBox  string
	RefX     int
	RefY     int
	Width    int
	Height   int
	Orient   string
	Content  string
}

// DecorationRegistry holds all decoration definitions
type DecorationRegistry struct {
	gradients []GradientDef
	filters   []FilterDef
	markers   []MarkerDef
	styles    []string
}

// NewDecorationRegistry creates a registry with beautiful decorations
func NewDecorationRegistry() *DecorationRegistry {
	r := &DecorationRegistry{}
	r.registerGradients()
	r.registerFilters()
	r.registerMarkers()
	r.registerStyles()
	return r
}
