package canvas

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
