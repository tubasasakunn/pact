package canvas

func (r *PatternRegistry) registerSequencePatterns() {
	// ============================================================
	// Request-Response
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternRequestResponse,
		Name:      "Request-Response",
		MinWidth:  350,
		MinHeight: 200,
		Padding:   20,
		Positions: []LayoutPosition{
			{ID: "caller", X: 0.25, Y: 0.15, Width: 0.25, Height: 0.18},
			{ID: "callee", X: 0.75, Y: 0.15, Width: 0.25, Height: 0.18},
		},
		Edges: []EdgePath{
			{FromID: "caller", ToID: "callee", CurveStyle: "orthogonal",
				LabelPos: Point{X: 0.5, Y: 0.45}},
			{FromID: "callee", ToID: "caller", CurveStyle: "orthogonal",
				LabelPos: Point{X: 0.5, Y: 0.75}},
		},
	})

	// ============================================================
	// Callback
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternCallback,
		Name:      "Callback Pattern",
		MinWidth:  380,
		MinHeight: 280,
		Padding:   20,
		Positions: []LayoutPosition{
			{ID: "initiator", X: 0.25, Y: 0.12, Width: 0.25, Height: 0.15},
			{ID: "handler", X: 0.75, Y: 0.12, Width: 0.25, Height: 0.15},
		},
		Edges: []EdgePath{
			{FromID: "initiator", ToID: "handler", CurveStyle: "orthogonal",
				LabelPos: Point{X: 0.5, Y: 0.32}},
			{FromID: "handler", ToID: "initiator", CurveStyle: "orthogonal",
				LabelPos: Point{X: 0.5, Y: 0.55}},
			{FromID: "initiator", ToID: "handler", CurveStyle: "orthogonal",
				LabelPos: Point{X: 0.5, Y: 0.78}},
		},
	})

	// ============================================================
	// Chain - 3 participants
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternChain3,
		Name:      "Chain (3 participants)",
		MinWidth:  500,
		MinHeight: 220,
		Padding:   20,
		Positions: []LayoutPosition{
			{ID: "node_0", X: 0.17, Y: 0.15, Width: 0.2, Height: 0.15},
			{ID: "node_1", X: 0.5, Y: 0.15, Width: 0.2, Height: 0.15},
			{ID: "node_2", X: 0.83, Y: 0.15, Width: 0.2, Height: 0.15},
		},
	})

	// ============================================================
	// Chain - 4 participants
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternChain4,
		Name:      "Chain (4 participants)",
		MinWidth:  650,
		MinHeight: 220,
		Padding:   20,
		Positions: []LayoutPosition{
			{ID: "node_0", X: 0.125, Y: 0.15, Width: 0.18, Height: 0.15},
			{ID: "node_1", X: 0.375, Y: 0.15, Width: 0.18, Height: 0.15},
			{ID: "node_2", X: 0.625, Y: 0.15, Width: 0.18, Height: 0.15},
			{ID: "node_3", X: 0.875, Y: 0.15, Width: 0.18, Height: 0.15},
		},
	})

	r.patterns[PatternChain] = r.patterns[PatternChain3]

	// ============================================================
	// Fan-Out
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternFanOut,
		Name:      "Fan-Out Pattern",
		MinWidth:  550,
		MinHeight: 250,
		Padding:   20,
		Positions: []LayoutPosition{
			{ID: "source", X: 0.15, Y: 0.15, Width: 0.2, Height: 0.15},
			{ID: "target_0", X: 0.5, Y: 0.15, Width: 0.18, Height: 0.15},
			{ID: "target_1", X: 0.72, Y: 0.15, Width: 0.18, Height: 0.15},
			{ID: "target_2", X: 0.94, Y: 0.15, Width: 0.18, Height: 0.15},
		},
	})
}
