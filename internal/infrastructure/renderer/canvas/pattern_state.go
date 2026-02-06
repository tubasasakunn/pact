package canvas

func (r *PatternRegistry) registerStatePatterns() {
	// ============================================================
	// Linear States - 2 states
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternLinearStates2,
		Name:      "Linear States (2)",
		MinWidth:  400,
		MinHeight: 120,
		Padding:   30,
		Positions: []LayoutPosition{
			{ID: "initial", X: 0.08, Y: 0.5, Width: 0.06, Height: 0.25},
			{ID: "state_0", X: 0.35, Y: 0.5, Width: 0.22, Height: 0.4},
			{ID: "state_1", X: 0.65, Y: 0.5, Width: 0.22, Height: 0.4},
			{ID: "final", X: 0.92, Y: 0.5, Width: 0.06, Height: 0.25},
		},
		Edges: []EdgePath{
			{FromID: "initial", ToID: "state_0", CurveStyle: "orthogonal"},
			{FromID: "state_0", ToID: "state_1", CurveStyle: "orthogonal"},
			{FromID: "state_1", ToID: "final", CurveStyle: "orthogonal"},
		},
	})

	// ============================================================
	// Linear States - 3 states
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternLinearStates3,
		Name:      "Linear States (3)",
		MinWidth:  550,
		MinHeight: 120,
		Padding:   30,
		Positions: []LayoutPosition{
			{ID: "initial", X: 0.06, Y: 0.5, Width: 0.05, Height: 0.22},
			{ID: "state_0", X: 0.25, Y: 0.5, Width: 0.18, Height: 0.4},
			{ID: "state_1", X: 0.5, Y: 0.5, Width: 0.18, Height: 0.4},
			{ID: "state_2", X: 0.75, Y: 0.5, Width: 0.18, Height: 0.4},
			{ID: "final", X: 0.94, Y: 0.5, Width: 0.05, Height: 0.22},
		},
		Edges: []EdgePath{
			{FromID: "initial", ToID: "state_0", CurveStyle: "orthogonal"},
			{FromID: "state_0", ToID: "state_1", CurveStyle: "orthogonal"},
			{FromID: "state_1", ToID: "state_2", CurveStyle: "orthogonal"},
			{FromID: "state_2", ToID: "final", CurveStyle: "orthogonal"},
		},
	})

	// ============================================================
	// Linear States - 4 states
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternLinearStates4,
		Name:      "Linear States (4)",
		MinWidth:  700,
		MinHeight: 120,
		Padding:   30,
		Positions: []LayoutPosition{
			{ID: "initial", X: 0.05, Y: 0.5, Width: 0.04, Height: 0.2},
			{ID: "state_0", X: 0.2, Y: 0.5, Width: 0.15, Height: 0.4},
			{ID: "state_1", X: 0.4, Y: 0.5, Width: 0.15, Height: 0.4},
			{ID: "state_2", X: 0.6, Y: 0.5, Width: 0.15, Height: 0.4},
			{ID: "state_3", X: 0.8, Y: 0.5, Width: 0.15, Height: 0.4},
			{ID: "final", X: 0.95, Y: 0.5, Width: 0.04, Height: 0.2},
		},
		Edges: []EdgePath{
			{FromID: "initial", ToID: "state_0", CurveStyle: "orthogonal"},
			{FromID: "state_0", ToID: "state_1", CurveStyle: "orthogonal"},
			{FromID: "state_1", ToID: "state_2", CurveStyle: "orthogonal"},
			{FromID: "state_2", ToID: "state_3", CurveStyle: "orthogonal"},
			{FromID: "state_3", ToID: "final", CurveStyle: "orthogonal"},
		},
	})

	r.patterns[PatternLinearStates] = r.patterns[PatternLinearStates3]

	// ============================================================
	// Binary Choice
	// ============================================================
	//          ┌───[TrueState]───┐
	//          │                 │
	// [Source]─┤                 ├─[Target]
	//          │                 │
	//          └───[FalseState]──┘
	r.Register(&PatternLayout{
		Type:      PatternBinaryChoice,
		Name:      "Binary Choice",
		MinWidth:  550,
		MinHeight: 280,
		Padding:   35,
		Positions: []LayoutPosition{
			{ID: "source", X: 0.12, Y: 0.5, Width: 0.18, Height: 0.28},
			{ID: "true_branch", X: 0.5, Y: 0.22, Width: 0.2, Height: 0.22},
			{ID: "false_branch", X: 0.5, Y: 0.78, Width: 0.2, Height: 0.22},
			{ID: "target", X: 0.88, Y: 0.5, Width: 0.18, Height: 0.28},
		},
		Edges: []EdgePath{
			// Source to true branch: right then up
			{FromID: "source", ToID: "true_branch", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.28, Y: 0.5}, {X: 0.28, Y: 0.22}},
				LabelPos:  Point{X: 0.2, Y: 0.35}},
			// Source to false branch: right then down
			{FromID: "source", ToID: "false_branch", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.28, Y: 0.5}, {X: 0.28, Y: 0.78}},
				LabelPos:  Point{X: 0.2, Y: 0.65}},
			// True branch to target: right then down
			{FromID: "true_branch", ToID: "target", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.72, Y: 0.22}, {X: 0.72, Y: 0.5}}},
			// False branch to target: right then up
			{FromID: "false_branch", ToID: "target", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.72, Y: 0.78}, {X: 0.72, Y: 0.5}}},
		},
	})

	// ============================================================
	// State Loop
	// ============================================================
	//    ┌────────────────────────┐
	//    │                        │
	//    ▼                        │
	// [Start]───>[Middle]───>[End]┘
	r.Register(&PatternLayout{
		Type:      PatternStateLoop,
		Name:      "State Loop",
		MinWidth:  480,
		MinHeight: 220,
		Padding:   35,
		Positions: []LayoutPosition{
			{ID: "start", X: 0.17, Y: 0.65, Width: 0.2, Height: 0.32},
			{ID: "middle", X: 0.5, Y: 0.65, Width: 0.2, Height: 0.32},
			{ID: "end", X: 0.83, Y: 0.65, Width: 0.2, Height: 0.32},
		},
		Edges: []EdgePath{
			{FromID: "start", ToID: "middle", CurveStyle: "orthogonal"},
			{FromID: "middle", ToID: "end", CurveStyle: "orthogonal"},
			// Loop back: up, left, down
			{FromID: "end", ToID: "start", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.83, Y: 0.22}, {X: 0.17, Y: 0.22}},
				LabelPos:  Point{X: 0.5, Y: 0.15}},
		},
	})

	// ============================================================
	// Star Topology
	// ============================================================
	//       [Top]
	//         │
	// [Left]──[Center]──[Right]
	//         │
	//      [Bottom]
	r.Register(&PatternLayout{
		Type:      PatternStarTopology,
		Name:      "Star Topology",
		MinWidth:  400,
		MinHeight: 400,
		Padding:   35,
		Positions: []LayoutPosition{
			{ID: "center", X: 0.5, Y: 0.5, Width: 0.22, Height: 0.2},
			{ID: "top", X: 0.5, Y: 0.15, Width: 0.2, Height: 0.16},
			{ID: "right", X: 0.85, Y: 0.5, Width: 0.2, Height: 0.16},
			{ID: "bottom", X: 0.5, Y: 0.85, Width: 0.2, Height: 0.16},
			{ID: "left", X: 0.15, Y: 0.5, Width: 0.2, Height: 0.16},
		},
		Edges: []EdgePath{
			{FromID: "center", ToID: "top", CurveStyle: "orthogonal"},
			{FromID: "center", ToID: "right", CurveStyle: "orthogonal"},
			{FromID: "center", ToID: "bottom", CurveStyle: "orthogonal"},
			{FromID: "center", ToID: "left", CurveStyle: "orthogonal"},
		},
	})
}
