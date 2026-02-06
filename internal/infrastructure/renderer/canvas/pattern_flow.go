package canvas

func (r *PatternRegistry) registerFlowPatterns() {
	// ============================================================
	// If-Else
	// ============================================================
	//         ┌───[TrueProc]───┐
	//         │                │
	// [Dec]───┤                ├───[Merge]
	//         │                │
	//         └───[FalseProc]──┘
	r.Register(&PatternLayout{
		Type:      PatternIfElse,
		Name:      "If-Else",
		MinWidth:  520,
		MinHeight: 300,
		Padding:   35,
		Positions: []LayoutPosition{
			{ID: "decision", X: 0.12, Y: 0.5, Width: 0.14, Height: 0.22},
			{ID: "true_process", X: 0.45, Y: 0.22, Width: 0.22, Height: 0.18},
			{ID: "false_process", X: 0.45, Y: 0.78, Width: 0.22, Height: 0.18},
			{ID: "merge", X: 0.85, Y: 0.5, Width: 0.14, Height: 0.2},
		},
		Edges: []EdgePath{
			{FromID: "decision", ToID: "true_process", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.25, Y: 0.5}, {X: 0.25, Y: 0.22}},
				LabelPos:  Point{X: 0.18, Y: 0.35}},
			{FromID: "decision", ToID: "false_process", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.25, Y: 0.5}, {X: 0.25, Y: 0.78}},
				LabelPos:  Point{X: 0.18, Y: 0.65}},
			{FromID: "true_process", ToID: "merge", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.68, Y: 0.22}, {X: 0.68, Y: 0.5}}},
			{FromID: "false_process", ToID: "merge", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.68, Y: 0.78}, {X: 0.68, Y: 0.5}}},
		},
		Decorators: []PatternDecorator{
			{Type: "label", Bounds: LayoutPosition{X: 0.18, Y: 0.32},
				Style: map[string]string{"text": "Yes"}},
			{Type: "label", Bounds: LayoutPosition{X: 0.18, Y: 0.68},
				Style: map[string]string{"text": "No"}},
		},
	})

	// ============================================================
	// If-ElseIf-Else
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternIfElseIfElse,
		Name:      "If-ElseIf-Else",
		MinWidth:  600,
		MinHeight: 420,
		Padding:   35,
		Positions: []LayoutPosition{
			{ID: "decision1", X: 0.1, Y: 0.18, Width: 0.12, Height: 0.14},
			{ID: "decision2", X: 0.1, Y: 0.5, Width: 0.12, Height: 0.14},
			{ID: "process1", X: 0.42, Y: 0.18, Width: 0.2, Height: 0.12},
			{ID: "process2", X: 0.42, Y: 0.5, Width: 0.2, Height: 0.12},
			{ID: "process3", X: 0.42, Y: 0.82, Width: 0.2, Height: 0.12},
			{ID: "merge", X: 0.85, Y: 0.5, Width: 0.12, Height: 0.14},
		},
		Edges: []EdgePath{
			// Decision1 -> Process1 (horizontal)
			{FromID: "decision1", ToID: "process1", CurveStyle: "orthogonal",
				LabelPos: Point{X: 0.26, Y: 0.12}},
			// Decision1 -> Decision2 (vertical down)
			{FromID: "decision1", ToID: "decision2", CurveStyle: "orthogonal"},
			// Decision2 -> Process2 (horizontal)
			{FromID: "decision2", ToID: "process2", CurveStyle: "orthogonal",
				LabelPos: Point{X: 0.26, Y: 0.44}},
			// Decision2 -> Process3 (down then right)
			{FromID: "decision2", ToID: "process3", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.1, Y: 0.82}},
				LabelPos:  Point{X: 0.05, Y: 0.66}},
			// All processes -> Merge
			{FromID: "process1", ToID: "merge", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.68, Y: 0.18}, {X: 0.68, Y: 0.5}}},
			{FromID: "process2", ToID: "merge", CurveStyle: "orthogonal"},
			{FromID: "process3", ToID: "merge", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.68, Y: 0.82}, {X: 0.68, Y: 0.5}}},
		},
	})

	// ============================================================
	// While Loop
	// ============================================================
	//    ┌─────────────────┐
	//    │                 │
	//    ▼                 │
	// [Cond]───>[Body]─────┘
	//    │
	//    ▼
	// [Exit]
	r.Register(&PatternLayout{
		Type:      PatternWhileLoop,
		Name:      "While Loop",
		MinWidth:  420,
		MinHeight: 280,
		Padding:   35,
		Positions: []LayoutPosition{
			{ID: "condition", X: 0.25, Y: 0.35, Width: 0.18, Height: 0.2},
			{ID: "body", X: 0.7, Y: 0.35, Width: 0.22, Height: 0.18},
			{ID: "exit", X: 0.25, Y: 0.82, Width: 0.16, Height: 0.14},
		},
		Edges: []EdgePath{
			// Condition -> Body (horizontal)
			{FromID: "condition", ToID: "body", CurveStyle: "orthogonal",
				LabelPos: Point{X: 0.47, Y: 0.28}},
			// Body -> Condition (loop back: up, left, down)
			{FromID: "body", ToID: "condition", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.7, Y: 0.12}, {X: 0.25, Y: 0.12}},
				LabelPos:  Point{X: 0.47, Y: 0.06}},
			// Condition -> Exit (vertical down)
			{FromID: "condition", ToID: "exit", CurveStyle: "orthogonal",
				LabelPos: Point{X: 0.18, Y: 0.58}},
		},
		Decorators: []PatternDecorator{
			{Type: "label", Bounds: LayoutPosition{X: 0.47, Y: 0.26},
				Style: map[string]string{"text": "Yes"}},
			{Type: "label", Bounds: LayoutPosition{X: 0.16, Y: 0.56},
				Style: map[string]string{"text": "No"}},
		},
	})

	// ============================================================
	// Sequential - 3 steps
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternSequential3,
		Name:      "Sequential (3 steps)",
		MinWidth:  550,
		MinHeight: 100,
		Padding:   25,
		Positions: []LayoutPosition{
			{ID: "start", X: 0.06, Y: 0.5, Width: 0.06, Height: 0.35},
			{ID: "process_0", X: 0.28, Y: 0.5, Width: 0.18, Height: 0.45},
			{ID: "process_1", X: 0.52, Y: 0.5, Width: 0.18, Height: 0.45},
			{ID: "process_2", X: 0.76, Y: 0.5, Width: 0.18, Height: 0.45},
			{ID: "end", X: 0.94, Y: 0.5, Width: 0.06, Height: 0.35},
		},
		Edges: []EdgePath{
			{FromID: "start", ToID: "process_0", CurveStyle: "orthogonal"},
			{FromID: "process_0", ToID: "process_1", CurveStyle: "orthogonal"},
			{FromID: "process_1", ToID: "process_2", CurveStyle: "orthogonal"},
			{FromID: "process_2", ToID: "end", CurveStyle: "orthogonal"},
		},
	})

	// ============================================================
	// Sequential - 4 steps
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternSequential4,
		Name:      "Sequential (4 steps)",
		MinWidth:  700,
		MinHeight: 100,
		Padding:   25,
		Positions: []LayoutPosition{
			{ID: "start", X: 0.05, Y: 0.5, Width: 0.05, Height: 0.35},
			{ID: "process_0", X: 0.22, Y: 0.5, Width: 0.15, Height: 0.45},
			{ID: "process_1", X: 0.42, Y: 0.5, Width: 0.15, Height: 0.45},
			{ID: "process_2", X: 0.62, Y: 0.5, Width: 0.15, Height: 0.45},
			{ID: "process_3", X: 0.82, Y: 0.5, Width: 0.15, Height: 0.45},
			{ID: "end", X: 0.95, Y: 0.5, Width: 0.05, Height: 0.35},
		},
		Edges: []EdgePath{
			{FromID: "start", ToID: "process_0", CurveStyle: "orthogonal"},
			{FromID: "process_0", ToID: "process_1", CurveStyle: "orthogonal"},
			{FromID: "process_1", ToID: "process_2", CurveStyle: "orthogonal"},
			{FromID: "process_2", ToID: "process_3", CurveStyle: "orthogonal"},
			{FromID: "process_3", ToID: "end", CurveStyle: "orthogonal"},
		},
	})

	r.patterns[PatternSequential] = r.patterns[PatternSequential3]
}
