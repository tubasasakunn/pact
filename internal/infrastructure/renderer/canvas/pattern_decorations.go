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

func (r *DecorationRegistry) registerGradients() {
	// Elegant blue gradient for headers
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-header-blue",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#4299e1", Opacity: "1"},
			{Offset: "100%", Color: "#3182ce", Opacity: "1"},
		},
	})

	// Soft purple gradient for interfaces
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-interface-purple",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#b794f4", Opacity: "1"},
			{Offset: "100%", Color: "#9f7aea", Opacity: "1"},
		},
	})

	// Green gradient for success states
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-success-green",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#68d391", Opacity: "1"},
			{Offset: "100%", Color: "#48bb78", Opacity: "1"},
		},
	})

	// Warm orange gradient for process nodes
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-process-orange",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#fbd38d", Opacity: "1"},
			{Offset: "100%", Color: "#f6ad55", Opacity: "1"},
		},
	})

	// Subtle gray gradient for backgrounds
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-bg-subtle",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#ffffff", Opacity: "1"},
			{Offset: "100%", Color: "#f7fafc", Opacity: "1"},
		},
	})

	// Radial glow for highlights
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-glow-blue",
		Type: "radial",
		Attrs: map[string]string{
			"cx": "50%", "cy": "50%", "r": "50%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#63b3ed", Opacity: "0.4"},
			{Offset: "100%", Color: "#63b3ed", Opacity: "0"},
		},
	})

	// Cool teal gradient for database nodes
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-database-teal",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#4fd1c5", Opacity: "1"},
			{Offset: "100%", Color: "#38b2ac", Opacity: "1"},
		},
	})

	// Red gradient for error/exception states
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-error-red",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#fc8181", Opacity: "1"},
			{Offset: "100%", Color: "#f56565", Opacity: "1"},
		},
	})

	// Elegant dark gradient for terminal nodes
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-terminal-dark",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#4a5568", Opacity: "1"},
			{Offset: "100%", Color: "#2d3748", Opacity: "1"},
		},
	})

	// Pink gradient for notes and annotations
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-note-pink",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "100%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#fef3c7", Opacity: "1"},
			{Offset: "100%", Color: "#fde68a", Opacity: "1"},
		},
	})
}

func (r *DecorationRegistry) registerFilters() {
	// Soft drop shadow
	r.filters = append(r.filters, FilterDef{
		ID: "filter-shadow-soft",
		Content: `<feDropShadow dx="0" dy="2" stdDeviation="3" flood-color="#000" flood-opacity="0.15"/>`,
	})

	// Elevated shadow (more prominent)
	r.filters = append(r.filters, FilterDef{
		ID: "filter-shadow-elevated",
		Content: `<feDropShadow dx="0" dy="4" stdDeviation="6" flood-color="#000" flood-opacity="0.2"/>`,
	})

	// Inner shadow for depth
	r.filters = append(r.filters, FilterDef{
		ID: "filter-inner-shadow",
		Content: `<feOffset dx="0" dy="2" in="SourceAlpha" result="offset"/>
<feGaussianBlur in="offset" stdDeviation="2" result="blur"/>
<feComposite in="SourceGraphic" in2="blur" operator="over"/>`,
	})

	// Glow effect for highlights
	r.filters = append(r.filters, FilterDef{
		ID: "filter-glow-blue",
		Content: `<feGaussianBlur in="SourceGraphic" stdDeviation="3" result="blur"/>
<feColorMatrix in="blur" type="matrix" values="0 0 0 0 0.39  0 0 0 0 0.7  0 0 0 0 0.93  0 0 0 0.6 0"/>
<feMerge><feMergeNode/><feMergeNode in="SourceGraphic"/></feMerge>`,
	})

	// Success glow
	r.filters = append(r.filters, FilterDef{
		ID: "filter-glow-green",
		Content: `<feGaussianBlur in="SourceGraphic" stdDeviation="2" result="blur"/>
<feColorMatrix in="blur" type="matrix" values="0 0 0 0 0.28  0 0 0 0 0.73  0 0 0 0 0.47  0 0 0 0.5 0"/>
<feMerge><feMergeNode/><feMergeNode in="SourceGraphic"/></feMerge>`,
	})

	// Error glow
	r.filters = append(r.filters, FilterDef{
		ID: "filter-glow-red",
		Content: `<feGaussianBlur in="SourceGraphic" stdDeviation="2" result="blur"/>
<feColorMatrix in="blur" type="matrix" values="0 0 0 0 0.96  0 0 0 0 0.4  0 0 0 0 0.4  0 0 0 0.5 0"/>
<feMerge><feMergeNode/><feMergeNode in="SourceGraphic"/></feMerge>`,
	})

	// Emboss effect
	r.filters = append(r.filters, FilterDef{
		ID: "filter-emboss",
		Content: `<feConvolveMatrix order="3" kernelMatrix="-2 -1 0 -1 1 1 0 1 2" preserveAlpha="true"/>`,
	})

	// Blur for background elements
	r.filters = append(r.filters, FilterDef{
		ID: "filter-blur-bg",
		Content: `<feGaussianBlur in="SourceGraphic" stdDeviation="1.5"/>`,
	})
}

func (r *DecorationRegistry) registerMarkers() {
	// Elegant filled arrow
	r.markers = append(r.markers, MarkerDef{
		ID:      "marker-arrow-filled",
		ViewBox: "0 0 10 10",
		RefX:    9, RefY: 5,
		Width: 8, Height: 8,
		Orient: "auto",
		Content: `<path d="M0,0 L10,5 L0,10 z" fill="#2d3748"/>`,
	})

	// Open arrow (inheritance)
	r.markers = append(r.markers, MarkerDef{
		ID:      "marker-arrow-open",
		ViewBox: "0 0 10 10",
		RefX:    9, RefY: 5,
		Width: 10, Height: 10,
		Orient: "auto",
		Content: `<path d="M0,0 L10,5 L0,10" fill="none" stroke="#2d3748" stroke-width="1.5"/>`,
	})

	// Triangle (generalization)
	r.markers = append(r.markers, MarkerDef{
		ID:      "marker-triangle",
		ViewBox: "0 0 12 12",
		RefX:    11, RefY: 6,
		Width: 12, Height: 12,
		Orient: "auto",
		Content: `<path d="M0,0 L12,6 L0,12 z" fill="#fff" stroke="#2d3748" stroke-width="1.5"/>`,
	})

	// Filled diamond (composition)
	r.markers = append(r.markers, MarkerDef{
		ID:      "marker-diamond-filled",
		ViewBox: "0 0 14 10",
		RefX:    13, RefY: 5,
		Width: 14, Height: 10,
		Orient: "auto",
		Content: `<path d="M0,5 L7,0 L14,5 L7,10 z" fill="#2d3748"/>`,
	})

	// Open diamond (aggregation)
	r.markers = append(r.markers, MarkerDef{
		ID:      "marker-diamond-open",
		ViewBox: "0 0 14 10",
		RefX:    13, RefY: 5,
		Width: 14, Height: 10,
		Orient: "auto",
		Content: `<path d="M0,5 L7,0 L14,5 L7,10 z" fill="#fff" stroke="#2d3748" stroke-width="1.5"/>`,
	})

	// Circle marker for state diagrams
	r.markers = append(r.markers, MarkerDef{
		ID:      "marker-circle",
		ViewBox: "0 0 10 10",
		RefX:    5, RefY: 5,
		Width: 8, Height: 8,
		Orient: "auto",
		Content: `<circle cx="5" cy="5" r="4" fill="#2d3748"/>`,
	})

	// Async arrow (open with line)
	r.markers = append(r.markers, MarkerDef{
		ID:      "marker-arrow-async",
		ViewBox: "0 0 12 10",
		RefX:    11, RefY: 5,
		Width: 10, Height: 8,
		Orient: "auto",
		Content: `<path d="M0,0 L12,5 M0,10 L12,5" fill="none" stroke="#2d3748" stroke-width="1.5"/>`,
	})
}

func (r *DecorationRegistry) registerStyles() {
	// Base typography
	r.styles = append(r.styles, `
.diagram-text {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  font-size: 14px;
  fill: #1a202c;
}
.diagram-text-sm {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  font-size: 12px;
  fill: #4a5568;
}
.diagram-text-bold {
  font-weight: 600;
}
.diagram-text-italic {
  font-style: italic;
}
.diagram-text-mono {
  font-family: "SF Mono", "Consolas", "Monaco", monospace;
  font-size: 13px;
}`)

	// Node styles
	r.styles = append(r.styles, `
.node-class {
  fill: #fafbfc;
  stroke: #2d3748;
  stroke-width: 2;
}
.node-interface {
  fill: #f0fff4;
  stroke: #276749;
  stroke-width: 2;
  stroke-dasharray: 5,3;
}
.node-abstract {
  fill: #faf5ff;
  stroke: #6b46c1;
  stroke-width: 2;
}
.node-state {
  fill: #fafbfc;
  stroke: #2d3748;
  stroke-width: 2;
  rx: 10;
  ry: 10;
}
.node-process {
  fill: #fafbfc;
  stroke: #2d3748;
  stroke-width: 2;
}
.node-decision {
  fill: #fffaf0;
  stroke: #c05621;
  stroke-width: 2;
}
.node-terminal {
  fill: #2d3748;
  stroke: #1a202c;
  stroke-width: 2;
}`)

	// Edge styles
	r.styles = append(r.styles, `
.edge-solid {
  fill: none;
  stroke: #2d3748;
  stroke-width: 1.5;
}
.edge-dashed {
  fill: none;
  stroke: #2d3748;
  stroke-width: 1.5;
  stroke-dasharray: 6,4;
}
.edge-dotted {
  fill: none;
  stroke: #718096;
  stroke-width: 1.5;
  stroke-dasharray: 2,3;
}
.edge-async {
  fill: none;
  stroke: #2d3748;
  stroke-width: 1.5;
  stroke-dasharray: 8,4;
}
.edge-return {
  fill: none;
  stroke: #48bb78;
  stroke-width: 1.5;
  stroke-dasharray: 4,4;
}`)

	// Special decorations
	r.styles = append(r.styles, `
.note-box {
  fill: #fefcbf;
  stroke: #d69e2e;
  stroke-width: 1.5;
}
.fragment-box {
  fill: none;
  stroke: #4a5568;
  stroke-width: 1.5;
}
.fragment-label {
  fill: #e2e8f0;
  stroke: #4a5568;
  stroke-width: 1;
}
.swimlane-header {
  fill: #edf2f7;
  stroke: #cbd5e0;
  stroke-width: 1;
}
.swimlane-divider {
  stroke: #e2e8f0;
  stroke-width: 1;
  stroke-dasharray: 4,4;
}
.activation-bar {
  fill: #e2e8f0;
  stroke: #a0aec0;
  stroke-width: 1;
}
.lifeline {
  stroke: #a0aec0;
  stroke-width: 1;
  stroke-dasharray: 6,4;
}`)

	// Hover and interaction states
	r.styles = append(r.styles, `
.interactive:hover {
  filter: brightness(1.05);
  cursor: pointer;
}
.highlight {
  filter: url(#filter-glow-blue);
}
.selected {
  stroke-width: 3;
  stroke: #4299e1;
}`)
}

// GenerateDefs generates SVG <defs> content for all decorations
func (r *DecorationRegistry) GenerateDefs() string {
	var defs string

	// Gradients
	for _, g := range r.gradients {
		switch g.Type {
		case "linear":
			defs += `<linearGradient id="` + g.ID + `"`
			for k, v := range g.Attrs {
				defs += ` ` + k + `="` + v + `"`
			}
			defs += `>`
			for _, s := range g.Stops {
				defs += `<stop offset="` + s.Offset + `" stop-color="` + s.Color + `"`
				if s.Opacity != "" {
					defs += ` stop-opacity="` + s.Opacity + `"`
				}
				defs += `/>`
			}
			defs += `</linearGradient>`
		case "radial":
			defs += `<radialGradient id="` + g.ID + `"`
			for k, v := range g.Attrs {
				defs += ` ` + k + `="` + v + `"`
			}
			defs += `>`
			for _, s := range g.Stops {
				defs += `<stop offset="` + s.Offset + `" stop-color="` + s.Color + `"`
				if s.Opacity != "" {
					defs += ` stop-opacity="` + s.Opacity + `"`
				}
				defs += `/>`
			}
			defs += `</radialGradient>`
		}
	}

	// Filters
	for _, f := range r.filters {
		defs += `<filter id="` + f.ID + `">` + f.Content + `</filter>`
	}

	// Markers
	for _, m := range r.markers {
		defs += `<marker id="` + m.ID + `" viewBox="` + m.ViewBox + `"`
		defs += ` refX="` + itoa(m.RefX) + `" refY="` + itoa(m.RefY) + `"`
		defs += ` markerWidth="` + itoa(m.Width) + `" markerHeight="` + itoa(m.Height) + `"`
		defs += ` orient="` + m.Orient + `">` + m.Content + `</marker>`
	}

	return defs
}

// GenerateStyles generates CSS content for all styles
func (r *DecorationRegistry) GenerateStyles() string {
	var css string
	for _, s := range r.styles {
		css += s
	}
	return css
}

// ApplyToCanvas applies decorations to a Canvas
func (r *DecorationRegistry) ApplyToCanvas(c *Canvas) {
	// Add defs
	defs := r.GenerateDefs()
	c.AddDef(defs)

	// Add styles
	styles := r.GenerateStyles()
	c.AddDef(`<style>` + styles + `</style>`)
}

// NodeStyleClass returns the appropriate CSS class for a node type
func NodeStyleClass(nodeType string) string {
	switch nodeType {
	case "class":
		return "node-class"
	case "interface":
		return "node-interface"
	case "abstract":
		return "node-abstract"
	case "state":
		return "node-state"
	case "process":
		return "node-process"
	case "decision":
		return "node-decision"
	case "terminal":
		return "node-terminal"
	default:
		return "node-class"
	}
}

// EdgeStyleClass returns the appropriate CSS class for an edge type
func EdgeStyleClass(edgeType string) string {
	switch edgeType {
	case "solid", "inheritance", "composition", "aggregation":
		return "edge-solid"
	case "dashed", "dependency", "implementation":
		return "edge-dashed"
	case "dotted":
		return "edge-dotted"
	case "async":
		return "edge-async"
	case "return":
		return "edge-return"
	default:
		return "edge-solid"
	}
}

// GradientForNodeType returns the appropriate gradient ID for a node type
func GradientForNodeType(nodeType string) string {
	switch nodeType {
	case "interface":
		return "grad-interface-purple"
	case "database":
		return "grad-database-teal"
	case "process":
		return "grad-process-orange"
	case "terminal":
		return "grad-terminal-dark"
	case "success", "final":
		return "grad-success-green"
	case "error":
		return "grad-error-red"
	default:
		return "grad-bg-subtle"
	}
}

// FilterForState returns the appropriate filter ID for a node state
func FilterForState(state string) string {
	switch state {
	case "hover", "active":
		return "filter-glow-blue"
	case "success":
		return "filter-glow-green"
	case "error":
		return "filter-glow-red"
	case "elevated":
		return "filter-shadow-elevated"
	default:
		return "filter-shadow-soft"
	}
}

// MarkerForEdgeType returns the appropriate marker ID for an edge type
func MarkerForEdgeType(edgeType string) string {
	switch edgeType {
	case "inheritance", "generalization":
		return "marker-triangle"
	case "implementation", "realization":
		return "marker-triangle"
	case "composition":
		return "marker-diamond-filled"
	case "aggregation":
		return "marker-diamond-open"
	case "dependency":
		return "marker-arrow-open"
	case "async":
		return "marker-arrow-async"
	default:
		return "marker-arrow-filled"
	}
}
