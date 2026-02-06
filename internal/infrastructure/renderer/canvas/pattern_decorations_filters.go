package canvas

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
