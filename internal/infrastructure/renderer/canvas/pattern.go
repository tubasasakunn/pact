// Package canvas provides pattern templates for diagram layouts.
package canvas

// PatternType represents a recognized structural pattern
type PatternType string

// Class diagram patterns - size variations
const (
	// Inheritance patterns by child count
	PatternInheritanceTree2 PatternType = "inheritance-tree-2"
	PatternInheritanceTree3 PatternType = "inheritance-tree-3"
	PatternInheritanceTree4 PatternType = "inheritance-tree-4"

	// Interface implementation patterns by implementer count
	PatternInterfaceImpl2 PatternType = "interface-impl-2"
	PatternInterfaceImpl3 PatternType = "interface-impl-3"
	PatternInterfaceImpl4 PatternType = "interface-impl-4"

	// Composition patterns by part count
	PatternComposition2 PatternType = "composition-2"
	PatternComposition3 PatternType = "composition-3"
	PatternComposition4 PatternType = "composition-4"

	// Special patterns
	PatternDiamond    PatternType = "diamond"
	PatternLayered3x2 PatternType = "layered-3x2"
	PatternLayered3x3 PatternType = "layered-3x3"

	// Legacy aliases
	PatternInheritanceTree PatternType = "inheritance-tree"
	PatternInterfaceImpl   PatternType = "interface-impl"
	PatternComposition     PatternType = "composition"
	PatternLayered         PatternType = "layered"
)

// Sequence diagram patterns
const (
	PatternRequestResponse PatternType = "request-response"
	PatternCallback        PatternType = "callback"
	PatternChain3          PatternType = "chain-3"
	PatternChain4          PatternType = "chain-4"
	PatternFanOut          PatternType = "fan-out"

	// Legacy alias
	PatternChain PatternType = "chain"
)

// State diagram patterns
const (
	PatternLinearStates2 PatternType = "linear-states-2"
	PatternLinearStates3 PatternType = "linear-states-3"
	PatternLinearStates4 PatternType = "linear-states-4"
	PatternBinaryChoice  PatternType = "binary-choice"
	PatternStateLoop     PatternType = "state-loop"
	PatternStarTopology  PatternType = "star-topology"

	// Legacy alias
	PatternLinearStates PatternType = "linear-states"
)

// Flow diagram patterns
const (
	PatternIfElse       PatternType = "if-else"
	PatternIfElseIfElse PatternType = "if-elseif-else"
	PatternWhileLoop    PatternType = "while-loop"
	PatternSequential3  PatternType = "sequential-3"
	PatternSequential4  PatternType = "sequential-4"

	// Legacy alias
	PatternSequential PatternType = "sequential"
)

// LayoutPosition represents a positioned element in a pattern layout
type LayoutPosition struct {
	ID     string  // Element identifier
	X      float64 // Relative X position (0.0 - 1.0)
	Y      float64 // Relative Y position (0.0 - 1.0)
	Width  float64 // Relative width (0.0 - 1.0)
	Height float64 // Relative height (0.0 - 1.0)
}

// EdgePath represents a pre-defined edge path in a pattern
type EdgePath struct {
	FromID     string  // Source element ID
	ToID       string  // Target element ID
	Waypoints  []Point // Intermediate points (relative coordinates)
	LabelPos   Point   // Label position (relative coordinates)
	CurveStyle string  // "orthogonal" only - no diagonal lines
}

// Point represents a 2D point with relative coordinates
type Point struct {
	X float64
	Y float64
}

// PatternLayout defines the complete layout for a pattern
type PatternLayout struct {
	Type       PatternType
	Name       string
	MinWidth   int
	MinHeight  int
	Padding    int
	Positions  []LayoutPosition
	Edges      []EdgePath
	Decorators []PatternDecorator
}

// PatternDecorator adds visual elements to a pattern
type PatternDecorator struct {
	Type   string            // "background", "groupbox", "divider", "label"
	Bounds LayoutPosition    // Position and size
	Style  map[string]string // Style attributes
}

// PatternRegistry holds available pattern layouts
type PatternRegistry struct {
	patterns map[PatternType]*PatternLayout
}

// NewPatternRegistry creates a new pattern registry with built-in patterns
func NewPatternRegistry() *PatternRegistry {
	r := &PatternRegistry{
		patterns: make(map[PatternType]*PatternLayout),
	}
	r.registerBuiltinPatterns()
	return r
}

// Get returns a pattern layout by type
func (r *PatternRegistry) Get(t PatternType) *PatternLayout {
	return r.patterns[t]
}

// Register adds a pattern layout to the registry
func (r *PatternRegistry) Register(layout *PatternLayout) {
	r.patterns[layout.Type] = layout
}

func (r *PatternRegistry) registerBuiltinPatterns() {
	r.registerClassPatterns()
	r.registerSequencePatterns()
	r.registerStatePatterns()
	r.registerFlowPatterns()
}

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

// PatternMatcher provides pattern detection functionality
type PatternMatcher struct {
	registry *PatternRegistry
}

// NewPatternMatcher creates a new pattern matcher
func NewPatternMatcher(registry *PatternRegistry) *PatternMatcher {
	return &PatternMatcher{registry: registry}
}

// ClassPatternMatch represents a detected class diagram pattern
type ClassPatternMatch struct {
	Pattern   PatternType
	NodeRoles map[string]string // pattern role ID -> actual node ID
	EdgeRoles map[string]int    // pattern edge index -> actual edge index
	Score     float64           // Confidence score 0.0-1.0
}

// StatePatternMatch represents a detected state diagram pattern
type StatePatternMatch struct {
	Pattern    PatternType
	StateRoles map[string]string // pattern role ID -> actual state ID
	Score      float64
}

// FlowPatternMatch represents a detected flow diagram pattern
type FlowPatternMatch struct {
	Pattern   PatternType
	NodeRoles map[string]string // pattern role ID -> actual node ID
	Score     float64
}

// SequencePatternMatch represents a detected sequence diagram pattern
type SequencePatternMatch struct {
	Pattern          PatternType
	ParticipantRoles map[string]string // pattern role ID -> actual participant ID
	Score            float64
}
