package canvas

func (r *PatternRegistry) registerClassPatterns() {
	// ============================================================
	// Inheritance Tree - 2 children
	// ============================================================
	//         [Parent]
	//            │
	//      ┌─────┴─────┐
	//      │           │
	//   [Child1]   [Child2]
	r.Register(&PatternLayout{
		Type:      PatternInheritanceTree2,
		Name:      "Inheritance Tree (2 children)",
		MinWidth:  400,
		MinHeight: 280,
		Padding:   30,
		Positions: []LayoutPosition{
			{ID: "parent", X: 0.5, Y: 0.18, Width: 0.3, Height: 0.25},
			{ID: "child_0", X: 0.25, Y: 0.75, Width: 0.28, Height: 0.22},
			{ID: "child_1", X: 0.75, Y: 0.75, Width: 0.28, Height: 0.22},
		},
		Edges: []EdgePath{
			// Vertical down from parent, then horizontal split, then vertical down to children
			{FromID: "child_0", ToID: "parent", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.25, Y: 0.52}, {X: 0.5, Y: 0.52}}},
			{FromID: "child_1", ToID: "parent", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.75, Y: 0.52}, {X: 0.5, Y: 0.52}}},
		},
	})

	// ============================================================
	// Inheritance Tree - 3 children
	// ============================================================
	//              [Parent]
	//                 │
	//       ┌─────────┼─────────┐
	//       │         │         │
	//   [Child1]  [Child2]  [Child3]
	r.Register(&PatternLayout{
		Type:      PatternInheritanceTree3,
		Name:      "Inheritance Tree (3 children)",
		MinWidth:  520,
		MinHeight: 280,
		Padding:   30,
		Positions: []LayoutPosition{
			{ID: "parent", X: 0.5, Y: 0.18, Width: 0.26, Height: 0.25},
			{ID: "child_0", X: 0.17, Y: 0.75, Width: 0.24, Height: 0.22},
			{ID: "child_1", X: 0.5, Y: 0.75, Width: 0.24, Height: 0.22},
			{ID: "child_2", X: 0.83, Y: 0.75, Width: 0.24, Height: 0.22},
		},
		Edges: []EdgePath{
			{FromID: "child_0", ToID: "parent", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.17, Y: 0.52}, {X: 0.5, Y: 0.52}}},
			{FromID: "child_1", ToID: "parent", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.52}}},
			{FromID: "child_2", ToID: "parent", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.83, Y: 0.52}, {X: 0.5, Y: 0.52}}},
		},
	})

	// ============================================================
	// Inheritance Tree - 4 children
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternInheritanceTree4,
		Name:      "Inheritance Tree (4 children)",
		MinWidth:  650,
		MinHeight: 280,
		Padding:   30,
		Positions: []LayoutPosition{
			{ID: "parent", X: 0.5, Y: 0.18, Width: 0.22, Height: 0.25},
			{ID: "child_0", X: 0.125, Y: 0.75, Width: 0.2, Height: 0.22},
			{ID: "child_1", X: 0.375, Y: 0.75, Width: 0.2, Height: 0.22},
			{ID: "child_2", X: 0.625, Y: 0.75, Width: 0.2, Height: 0.22},
			{ID: "child_3", X: 0.875, Y: 0.75, Width: 0.2, Height: 0.22},
		},
		Edges: []EdgePath{
			{FromID: "child_0", ToID: "parent", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.125, Y: 0.52}, {X: 0.5, Y: 0.52}}},
			{FromID: "child_1", ToID: "parent", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.375, Y: 0.52}, {X: 0.5, Y: 0.52}}},
			{FromID: "child_2", ToID: "parent", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.625, Y: 0.52}, {X: 0.5, Y: 0.52}}},
			{FromID: "child_3", ToID: "parent", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.875, Y: 0.52}, {X: 0.5, Y: 0.52}}},
		},
	})

	// Legacy alias
	r.patterns[PatternInheritanceTree] = r.patterns[PatternInheritanceTree3]

	// ============================================================
	// Interface Implementation - 2 implementers
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternInterfaceImpl2,
		Name:      "Interface Implementation (2)",
		MinWidth:  400,
		MinHeight: 300,
		Padding:   30,
		Positions: []LayoutPosition{
			{ID: "interface", X: 0.5, Y: 0.15, Width: 0.32, Height: 0.22},
			{ID: "impl_0", X: 0.25, Y: 0.78, Width: 0.3, Height: 0.2},
			{ID: "impl_1", X: 0.75, Y: 0.78, Width: 0.3, Height: 0.2},
		},
		Edges: []EdgePath{
			{FromID: "impl_0", ToID: "interface", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.25, Y: 0.5}, {X: 0.5, Y: 0.5}}},
			{FromID: "impl_1", ToID: "interface", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.75, Y: 0.5}, {X: 0.5, Y: 0.5}}},
		},
		Decorators: []PatternDecorator{
			{Type: "label", Bounds: LayoutPosition{X: 0.5, Y: 0.05},
				Style: map[string]string{"text": "<<interface>>"}},
		},
	})

	// ============================================================
	// Interface Implementation - 3 implementers
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternInterfaceImpl3,
		Name:      "Interface Implementation (3)",
		MinWidth:  520,
		MinHeight: 300,
		Padding:   30,
		Positions: []LayoutPosition{
			{ID: "interface", X: 0.5, Y: 0.15, Width: 0.28, Height: 0.22},
			{ID: "impl_0", X: 0.17, Y: 0.78, Width: 0.26, Height: 0.2},
			{ID: "impl_1", X: 0.5, Y: 0.78, Width: 0.26, Height: 0.2},
			{ID: "impl_2", X: 0.83, Y: 0.78, Width: 0.26, Height: 0.2},
		},
		Edges: []EdgePath{
			{FromID: "impl_0", ToID: "interface", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.17, Y: 0.5}, {X: 0.5, Y: 0.5}}},
			{FromID: "impl_1", ToID: "interface", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.5}}},
			{FromID: "impl_2", ToID: "interface", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.83, Y: 0.5}, {X: 0.5, Y: 0.5}}},
		},
	})

	// ============================================================
	// Interface Implementation - 4 implementers
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternInterfaceImpl4,
		Name:      "Interface Implementation (4)",
		MinWidth:  650,
		MinHeight: 300,
		Padding:   30,
		Positions: []LayoutPosition{
			{ID: "interface", X: 0.5, Y: 0.15, Width: 0.24, Height: 0.22},
			{ID: "impl_0", X: 0.125, Y: 0.78, Width: 0.2, Height: 0.2},
			{ID: "impl_1", X: 0.375, Y: 0.78, Width: 0.2, Height: 0.2},
			{ID: "impl_2", X: 0.625, Y: 0.78, Width: 0.2, Height: 0.2},
			{ID: "impl_3", X: 0.875, Y: 0.78, Width: 0.2, Height: 0.2},
		},
		Edges: []EdgePath{
			{FromID: "impl_0", ToID: "interface", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.125, Y: 0.5}, {X: 0.5, Y: 0.5}}},
			{FromID: "impl_1", ToID: "interface", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.375, Y: 0.5}, {X: 0.5, Y: 0.5}}},
			{FromID: "impl_2", ToID: "interface", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.625, Y: 0.5}, {X: 0.5, Y: 0.5}}},
			{FromID: "impl_3", ToID: "interface", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.875, Y: 0.5}, {X: 0.5, Y: 0.5}}},
		},
	})

	r.patterns[PatternInterfaceImpl] = r.patterns[PatternInterfaceImpl3]

	// ============================================================
	// Composition - 2 parts (horizontal layout)
	// ============================================================
	//  [Owner]────◆────[Part1]
	//             │
	//             └────[Part2]
	r.Register(&PatternLayout{
		Type:      PatternComposition2,
		Name:      "Composition (2 parts)",
		MinWidth:  450,
		MinHeight: 220,
		Padding:   30,
		Positions: []LayoutPosition{
			{ID: "owner", X: 0.2, Y: 0.5, Width: 0.28, Height: 0.35},
			{ID: "part_0", X: 0.75, Y: 0.28, Width: 0.26, Height: 0.25},
			{ID: "part_1", X: 0.75, Y: 0.72, Width: 0.26, Height: 0.25},
		},
		Edges: []EdgePath{
			{FromID: "owner", ToID: "part_0", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.5}, {X: 0.5, Y: 0.28}}},
			{FromID: "owner", ToID: "part_1", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.5}, {X: 0.5, Y: 0.72}}},
		},
	})

	// ============================================================
	// Composition - 3 parts
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternComposition3,
		Name:      "Composition (3 parts)",
		MinWidth:  450,
		MinHeight: 300,
		Padding:   30,
		Positions: []LayoutPosition{
			{ID: "owner", X: 0.2, Y: 0.5, Width: 0.28, Height: 0.4},
			{ID: "part_0", X: 0.75, Y: 0.2, Width: 0.26, Height: 0.2},
			{ID: "part_1", X: 0.75, Y: 0.5, Width: 0.26, Height: 0.2},
			{ID: "part_2", X: 0.75, Y: 0.8, Width: 0.26, Height: 0.2},
		},
		Edges: []EdgePath{
			{FromID: "owner", ToID: "part_0", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.5}, {X: 0.5, Y: 0.2}}},
			{FromID: "owner", ToID: "part_1", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.5}}},
			{FromID: "owner", ToID: "part_2", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.5}, {X: 0.5, Y: 0.8}}},
		},
	})

	// ============================================================
	// Composition - 4 parts
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternComposition4,
		Name:      "Composition (4 parts)",
		MinWidth:  450,
		MinHeight: 380,
		Padding:   30,
		Positions: []LayoutPosition{
			{ID: "owner", X: 0.2, Y: 0.5, Width: 0.28, Height: 0.45},
			{ID: "part_0", X: 0.75, Y: 0.15, Width: 0.26, Height: 0.16},
			{ID: "part_1", X: 0.75, Y: 0.38, Width: 0.26, Height: 0.16},
			{ID: "part_2", X: 0.75, Y: 0.62, Width: 0.26, Height: 0.16},
			{ID: "part_3", X: 0.75, Y: 0.85, Width: 0.26, Height: 0.16},
		},
		Edges: []EdgePath{
			{FromID: "owner", ToID: "part_0", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.5}, {X: 0.5, Y: 0.15}}},
			{FromID: "owner", ToID: "part_1", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.5}, {X: 0.5, Y: 0.38}}},
			{FromID: "owner", ToID: "part_2", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.5}, {X: 0.5, Y: 0.62}}},
			{FromID: "owner", ToID: "part_3", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.5}, {X: 0.5, Y: 0.85}}},
		},
	})

	r.patterns[PatternComposition] = r.patterns[PatternComposition3]

	// ============================================================
	// Diamond Inheritance
	// ============================================================
	//       [Top]
	//         │
	//    ┌────┴────┐
	//    │         │
	// [Left]    [Right]
	//    │         │
	//    └────┬────┘
	//         │
	//      [Bottom]
	r.Register(&PatternLayout{
		Type:      PatternDiamond,
		Name:      "Diamond Inheritance",
		MinWidth:  400,
		MinHeight: 400,
		Padding:   40,
		Positions: []LayoutPosition{
			{ID: "top", X: 0.5, Y: 0.12, Width: 0.32, Height: 0.18},
			{ID: "left", X: 0.2, Y: 0.45, Width: 0.3, Height: 0.18},
			{ID: "right", X: 0.8, Y: 0.45, Width: 0.3, Height: 0.18},
			{ID: "bottom", X: 0.5, Y: 0.82, Width: 0.32, Height: 0.18},
		},
		Edges: []EdgePath{
			{FromID: "left", ToID: "top", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.2, Y: 0.28}, {X: 0.5, Y: 0.28}}},
			{FromID: "right", ToID: "top", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.8, Y: 0.28}, {X: 0.5, Y: 0.28}}},
			{FromID: "bottom", ToID: "left", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.65}, {X: 0.2, Y: 0.65}}},
			{FromID: "bottom", ToID: "right", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.65}, {X: 0.8, Y: 0.65}}},
		},
	})

	// ============================================================
	// Layered 3x2 (3 layers, 2 nodes per layer)
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternLayered3x2,
		Name:      "Layered Architecture (3x2)",
		MinWidth:  500,
		MinHeight: 400,
		Padding:   30,
		Positions: []LayoutPosition{
			// Layer 1
			{ID: "layer1_0", X: 0.3, Y: 0.12, Width: 0.28, Height: 0.16},
			{ID: "layer1_1", X: 0.7, Y: 0.12, Width: 0.28, Height: 0.16},
			// Layer 2
			{ID: "layer2_0", X: 0.3, Y: 0.45, Width: 0.28, Height: 0.16},
			{ID: "layer2_1", X: 0.7, Y: 0.45, Width: 0.28, Height: 0.16},
			// Layer 3
			{ID: "layer3_0", X: 0.3, Y: 0.78, Width: 0.28, Height: 0.16},
			{ID: "layer3_1", X: 0.7, Y: 0.78, Width: 0.28, Height: 0.16},
		},
		Decorators: []PatternDecorator{
			{Type: "divider", Bounds: LayoutPosition{X: 0.08, Y: 0.28, Width: 0.84, Height: 0},
				Style: map[string]string{"stroke": ColorSectionLine, "stroke-dasharray": "6,4"}},
			{Type: "divider", Bounds: LayoutPosition{X: 0.08, Y: 0.61, Width: 0.84, Height: 0},
				Style: map[string]string{"stroke": ColorSectionLine, "stroke-dasharray": "6,4"}},
		},
	})

	// ============================================================
	// Layered 3x3 (3 layers, 3 nodes per layer)
	// ============================================================
	r.Register(&PatternLayout{
		Type:      PatternLayered3x3,
		Name:      "Layered Architecture (3x3)",
		MinWidth:  650,
		MinHeight: 420,
		Padding:   30,
		Positions: []LayoutPosition{
			// Layer 1
			{ID: "layer1_0", X: 0.17, Y: 0.12, Width: 0.24, Height: 0.16},
			{ID: "layer1_1", X: 0.5, Y: 0.12, Width: 0.24, Height: 0.16},
			{ID: "layer1_2", X: 0.83, Y: 0.12, Width: 0.24, Height: 0.16},
			// Layer 2
			{ID: "layer2_0", X: 0.17, Y: 0.45, Width: 0.24, Height: 0.16},
			{ID: "layer2_1", X: 0.5, Y: 0.45, Width: 0.24, Height: 0.16},
			{ID: "layer2_2", X: 0.83, Y: 0.45, Width: 0.24, Height: 0.16},
			// Layer 3
			{ID: "layer3_0", X: 0.17, Y: 0.78, Width: 0.24, Height: 0.16},
			{ID: "layer3_1", X: 0.5, Y: 0.78, Width: 0.24, Height: 0.16},
			{ID: "layer3_2", X: 0.83, Y: 0.78, Width: 0.24, Height: 0.16},
		},
		Decorators: []PatternDecorator{
			{Type: "divider", Bounds: LayoutPosition{X: 0.05, Y: 0.28, Width: 0.9, Height: 0},
				Style: map[string]string{"stroke": ColorSectionLine, "stroke-dasharray": "6,4"}},
			{Type: "divider", Bounds: LayoutPosition{X: 0.05, Y: 0.61, Width: 0.9, Height: 0},
				Style: map[string]string{"stroke": ColorSectionLine, "stroke-dasharray": "6,4"}},
		},
	})

	r.patterns[PatternLayered] = r.patterns[PatternLayered3x3]
}
