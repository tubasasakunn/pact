// Package canvas provides pattern templates for diagram layouts.
package canvas

// PatternType represents a recognized structural pattern
type PatternType string

// Class diagram patterns
const (
	// PatternInheritanceTree: base class with multiple children vertically aligned below
	PatternInheritanceTree PatternType = "inheritance-tree"
	// PatternInterfaceImpl: interface with implementing classes in a fan shape
	PatternInterfaceImpl PatternType = "interface-impl"
	// PatternComposition: one class composed of multiple parts (1-to-many)
	PatternComposition PatternType = "composition"
	// PatternDiamond: diamond inheritance pattern (A -> B,C -> D)
	PatternDiamond PatternType = "diamond"
	// PatternLayered: horizontal layers of classes (e.g., MVC, 3-tier)
	PatternLayered PatternType = "layered"
)

// Sequence diagram patterns
const (
	// PatternRequestResponse: simple call and return between two participants
	PatternRequestResponse PatternType = "request-response"
	// PatternCallback: A calls B, B calls back to A
	PatternCallback PatternType = "callback"
	// PatternChain: sequential delegation through multiple participants
	PatternChain PatternType = "chain"
	// PatternFanOut: one participant calling multiple others
	PatternFanOut PatternType = "fan-out"
)

// State diagram patterns
const (
	// PatternLinearStates: simple sequence of states A -> B -> C
	PatternLinearStates PatternType = "linear-states"
	// PatternBinaryChoice: state with exactly two outgoing transitions
	PatternBinaryChoice PatternType = "binary-choice"
	// PatternStateLoop: state that can transition back
	PatternStateLoop PatternType = "state-loop"
	// PatternStarTopology: central state connected to multiple peripheral states
	PatternStarTopology PatternType = "star-topology"
)

// Flow diagram patterns
const (
	// PatternIfElse: decision with two branches that merge
	PatternIfElse PatternType = "if-else"
	// PatternIfElseIfElse: cascading decisions
	PatternIfElseIfElse PatternType = "if-elseif-else"
	// PatternWhileLoop: loop back to decision node
	PatternWhileLoop PatternType = "while-loop"
	// PatternSequential: simple linear flow without decisions
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
	FromID     string    // Source element ID
	ToID       string    // Target element ID
	Waypoints  []Point   // Intermediate points (relative coordinates)
	LabelPos   Point     // Label position (relative coordinates)
	CurveStyle string    // "straight", "orthogonal", "curved"
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
	// Class diagram patterns
	r.registerClassPatterns()
	// Sequence diagram patterns
	r.registerSequencePatterns()
	// State diagram patterns
	r.registerStatePatterns()
	// Flow diagram patterns
	r.registerFlowPatterns()
}

func (r *PatternRegistry) registerClassPatterns() {
	// Inheritance Tree: Parent centered on top, children spread below
	//       [Parent]
	//      /   |   \
	//  [C1]  [C2]  [C3]
	r.Register(&PatternLayout{
		Type:      PatternInheritanceTree,
		Name:      "Inheritance Tree",
		MinWidth:  600,
		MinHeight: 300,
		Padding:   40,
		Positions: []LayoutPosition{
			{ID: "parent", X: 0.5, Y: 0.15, Width: 0.25, Height: 0.3},
			{ID: "child_0", X: 0.15, Y: 0.7, Width: 0.2, Height: 0.25},
			{ID: "child_1", X: 0.4, Y: 0.7, Width: 0.2, Height: 0.25},
			{ID: "child_2", X: 0.65, Y: 0.7, Width: 0.2, Height: 0.25},
			{ID: "child_3", X: 0.9, Y: 0.7, Width: 0.2, Height: 0.25},
		},
		Edges: []EdgePath{
			{FromID: "child_0", ToID: "parent", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.15, Y: 0.5}, {X: 0.5, Y: 0.5}}},
			{FromID: "child_1", ToID: "parent", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.4, Y: 0.5}, {X: 0.5, Y: 0.5}}},
			{FromID: "child_2", ToID: "parent", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.65, Y: 0.5}, {X: 0.5, Y: 0.5}}},
			{FromID: "child_3", ToID: "parent", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.9, Y: 0.5}, {X: 0.5, Y: 0.5}}},
		},
	})

	// Interface Implementation: Interface on top, implementations in arc below
	//        <<interface>>
	//         [Interface]
	//        /     |     \
	//    [Impl1] [Impl2] [Impl3]
	r.Register(&PatternLayout{
		Type:      PatternInterfaceImpl,
		Name:      "Interface Implementation",
		MinWidth:  600,
		MinHeight: 320,
		Padding:   40,
		Positions: []LayoutPosition{
			{ID: "interface", X: 0.5, Y: 0.12, Width: 0.28, Height: 0.28},
			{ID: "impl_0", X: 0.12, Y: 0.72, Width: 0.22, Height: 0.24},
			{ID: "impl_1", X: 0.39, Y: 0.72, Width: 0.22, Height: 0.24},
			{ID: "impl_2", X: 0.66, Y: 0.72, Width: 0.22, Height: 0.24},
			{ID: "impl_3", X: 0.93, Y: 0.72, Width: 0.22, Height: 0.24},
		},
		Edges: []EdgePath{
			{FromID: "impl_0", ToID: "interface", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.12, Y: 0.5}, {X: 0.5, Y: 0.5}}},
			{FromID: "impl_1", ToID: "interface", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.39, Y: 0.5}, {X: 0.5, Y: 0.5}}},
			{FromID: "impl_2", ToID: "interface", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.66, Y: 0.5}, {X: 0.5, Y: 0.5}}},
			{FromID: "impl_3", ToID: "interface", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.93, Y: 0.5}, {X: 0.5, Y: 0.5}}},
		},
	})

	// Composition: Owner on left, parts stacked on right
	//  [Owner]────◆─────[Part1]
	//             │─────[Part2]
	//             └─────[Part3]
	r.Register(&PatternLayout{
		Type:      PatternComposition,
		Name:      "Composition",
		MinWidth:  500,
		MinHeight: 350,
		Padding:   40,
		Positions: []LayoutPosition{
			{ID: "owner", X: 0.18, Y: 0.5, Width: 0.28, Height: 0.4},
			{ID: "part_0", X: 0.75, Y: 0.18, Width: 0.22, Height: 0.2},
			{ID: "part_1", X: 0.75, Y: 0.42, Width: 0.22, Height: 0.2},
			{ID: "part_2", X: 0.75, Y: 0.66, Width: 0.22, Height: 0.2},
			{ID: "part_3", X: 0.75, Y: 0.9, Width: 0.22, Height: 0.2},
		},
		Edges: []EdgePath{
			{FromID: "owner", ToID: "part_0", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.5}, {X: 0.5, Y: 0.18}}, LabelPos: Point{X: 0.5, Y: 0.15}},
			{FromID: "owner", ToID: "part_1", CurveStyle: "straight", LabelPos: Point{X: 0.5, Y: 0.42}},
			{FromID: "owner", ToID: "part_2", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.5}, {X: 0.5, Y: 0.66}}, LabelPos: Point{X: 0.5, Y: 0.66}},
			{FromID: "owner", ToID: "part_3", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.5, Y: 0.5}, {X: 0.5, Y: 0.9}}, LabelPos: Point{X: 0.5, Y: 0.9}},
		},
	})

	// Diamond: Classic diamond inheritance
	//       [A]
	//      /   \
	//    [B]   [C]
	//      \   /
	//       [D]
	r.Register(&PatternLayout{
		Type:      PatternDiamond,
		Name:      "Diamond Inheritance",
		MinWidth:  450,
		MinHeight: 450,
		Padding:   50,
		Positions: []LayoutPosition{
			{ID: "top", X: 0.5, Y: 0.12, Width: 0.3, Height: 0.2},
			{ID: "left", X: 0.2, Y: 0.45, Width: 0.3, Height: 0.2},
			{ID: "right", X: 0.8, Y: 0.45, Width: 0.3, Height: 0.2},
			{ID: "bottom", X: 0.5, Y: 0.78, Width: 0.3, Height: 0.2},
		},
		Edges: []EdgePath{
			{FromID: "left", ToID: "top", CurveStyle: "straight"},
			{FromID: "right", ToID: "top", CurveStyle: "straight"},
			{FromID: "bottom", ToID: "left", CurveStyle: "straight"},
			{FromID: "bottom", ToID: "right", CurveStyle: "straight"},
		},
	})

	// Layered: Three-tier architecture style
	//  [Presentation Layer]
	//  ────────────────────
	//  [Business Layer]
	//  ────────────────────
	//  [Data Layer]
	r.Register(&PatternLayout{
		Type:      PatternLayered,
		Name:      "Layered Architecture",
		MinWidth:  700,
		MinHeight: 500,
		Padding:   30,
		Positions: []LayoutPosition{
			// Top layer
			{ID: "layer1_0", X: 0.25, Y: 0.1, Width: 0.2, Height: 0.18},
			{ID: "layer1_1", X: 0.5, Y: 0.1, Width: 0.2, Height: 0.18},
			{ID: "layer1_2", X: 0.75, Y: 0.1, Width: 0.2, Height: 0.18},
			// Middle layer
			{ID: "layer2_0", X: 0.25, Y: 0.42, Width: 0.2, Height: 0.18},
			{ID: "layer2_1", X: 0.5, Y: 0.42, Width: 0.2, Height: 0.18},
			{ID: "layer2_2", X: 0.75, Y: 0.42, Width: 0.2, Height: 0.18},
			// Bottom layer
			{ID: "layer3_0", X: 0.25, Y: 0.74, Width: 0.2, Height: 0.18},
			{ID: "layer3_1", X: 0.5, Y: 0.74, Width: 0.2, Height: 0.18},
			{ID: "layer3_2", X: 0.75, Y: 0.74, Width: 0.2, Height: 0.18},
		},
		Decorators: []PatternDecorator{
			{Type: "divider", Bounds: LayoutPosition{X: 0.1, Y: 0.32, Width: 0.8, Height: 0.01},
				Style: map[string]string{"stroke": ColorSectionLine, "stroke-dasharray": "5,5"}},
			{Type: "divider", Bounds: LayoutPosition{X: 0.1, Y: 0.64, Width: 0.8, Height: 0.01},
				Style: map[string]string{"stroke": ColorSectionLine, "stroke-dasharray": "5,5"}},
		},
	})
}

func (r *PatternRegistry) registerSequencePatterns() {
	// Request-Response: Simple synchronous call pattern
	//  [A]          [B]
	//   |──request──>|
	//   |<──response─|
	r.Register(&PatternLayout{
		Type:      PatternRequestResponse,
		Name:      "Request-Response",
		MinWidth:  300,
		MinHeight: 200,
		Padding:   20,
		Positions: []LayoutPosition{
			{ID: "caller", X: 0.25, Y: 0.0, Width: 0.2, Height: 0.15},
			{ID: "callee", X: 0.75, Y: 0.0, Width: 0.2, Height: 0.15},
		},
		Edges: []EdgePath{
			{FromID: "caller", ToID: "callee", CurveStyle: "straight",
				LabelPos: Point{X: 0.5, Y: 0.4}},
			{FromID: "callee", ToID: "caller", CurveStyle: "straight",
				LabelPos: Point{X: 0.5, Y: 0.7}},
		},
	})

	// Callback: Bidirectional interaction
	//  [A]          [B]
	//   |──call────>|
	//   |<──callback─|
	//   |──response─>|
	r.Register(&PatternLayout{
		Type:      PatternCallback,
		Name:      "Callback",
		MinWidth:  350,
		MinHeight: 280,
		Padding:   20,
		Positions: []LayoutPosition{
			{ID: "initiator", X: 0.25, Y: 0.0, Width: 0.2, Height: 0.12},
			{ID: "handler", X: 0.75, Y: 0.0, Width: 0.2, Height: 0.12},
		},
		Edges: []EdgePath{
			{FromID: "initiator", ToID: "handler", CurveStyle: "straight",
				LabelPos: Point{X: 0.5, Y: 0.28}},
			{FromID: "handler", ToID: "initiator", CurveStyle: "straight",
				LabelPos: Point{X: 0.5, Y: 0.52}},
			{FromID: "initiator", ToID: "handler", CurveStyle: "straight",
				LabelPos: Point{X: 0.5, Y: 0.76}},
		},
	})

	// Chain: Sequential delegation
	//  [A]     [B]     [C]     [D]
	//   |──────>|──────>|──────>|
	//   |<──────|<──────|<──────|
	r.Register(&PatternLayout{
		Type:      PatternChain,
		Name:      "Chain of Responsibility",
		MinWidth:  600,
		MinHeight: 250,
		Padding:   20,
		Positions: []LayoutPosition{
			{ID: "node_0", X: 0.125, Y: 0.0, Width: 0.15, Height: 0.12},
			{ID: "node_1", X: 0.375, Y: 0.0, Width: 0.15, Height: 0.12},
			{ID: "node_2", X: 0.625, Y: 0.0, Width: 0.15, Height: 0.12},
			{ID: "node_3", X: 0.875, Y: 0.0, Width: 0.15, Height: 0.12},
		},
	})

	// Fan-Out: One calling many
	//       [A]
	//    /   |   \
	//  [B]  [C]  [D]
	r.Register(&PatternLayout{
		Type:      PatternFanOut,
		Name:      "Fan-Out",
		MinWidth:  500,
		MinHeight: 300,
		Padding:   20,
		Positions: []LayoutPosition{
			{ID: "source", X: 0.5, Y: 0.0, Width: 0.18, Height: 0.12},
			{ID: "target_0", X: 0.2, Y: 0.0, Width: 0.15, Height: 0.12},
			{ID: "target_1", X: 0.5, Y: 0.0, Width: 0.15, Height: 0.12},
			{ID: "target_2", X: 0.8, Y: 0.0, Width: 0.15, Height: 0.12},
		},
	})
}

func (r *PatternRegistry) registerStatePatterns() {
	// Linear States: Sequential flow
	//  (●)──>[A]──>[B]──>[C]──>(◉)
	r.Register(&PatternLayout{
		Type:      PatternLinearStates,
		Name:      "Linear States",
		MinWidth:  700,
		MinHeight: 150,
		Padding:   30,
		Positions: []LayoutPosition{
			{ID: "initial", X: 0.08, Y: 0.5, Width: 0.05, Height: 0.2},
			{ID: "state_0", X: 0.25, Y: 0.5, Width: 0.15, Height: 0.35},
			{ID: "state_1", X: 0.45, Y: 0.5, Width: 0.15, Height: 0.35},
			{ID: "state_2", X: 0.65, Y: 0.5, Width: 0.15, Height: 0.35},
			{ID: "state_3", X: 0.85, Y: 0.5, Width: 0.15, Height: 0.35},
			{ID: "final", X: 0.96, Y: 0.5, Width: 0.05, Height: 0.2},
		},
		Edges: []EdgePath{
			{FromID: "initial", ToID: "state_0", CurveStyle: "straight"},
			{FromID: "state_0", ToID: "state_1", CurveStyle: "straight"},
			{FromID: "state_1", ToID: "state_2", CurveStyle: "straight"},
			{FromID: "state_2", ToID: "state_3", CurveStyle: "straight"},
			{FromID: "state_3", ToID: "final", CurveStyle: "straight"},
		},
	})

	// Binary Choice: Decision point
	//           ┌──>[TrueState]──┐
	//  [State]──┤                ├──>[Next]
	//           └──>[FalseState]─┘
	r.Register(&PatternLayout{
		Type:      PatternBinaryChoice,
		Name:      "Binary Choice",
		MinWidth:  600,
		MinHeight: 300,
		Padding:   40,
		Positions: []LayoutPosition{
			{ID: "source", X: 0.12, Y: 0.5, Width: 0.18, Height: 0.3},
			{ID: "true_branch", X: 0.5, Y: 0.22, Width: 0.18, Height: 0.25},
			{ID: "false_branch", X: 0.5, Y: 0.78, Width: 0.18, Height: 0.25},
			{ID: "target", X: 0.88, Y: 0.5, Width: 0.18, Height: 0.3},
		},
		Edges: []EdgePath{
			{FromID: "source", ToID: "true_branch", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.3, Y: 0.5}, {X: 0.3, Y: 0.22}},
				LabelPos:  Point{X: 0.22, Y: 0.35}},
			{FromID: "source", ToID: "false_branch", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.3, Y: 0.5}, {X: 0.3, Y: 0.78}},
				LabelPos:  Point{X: 0.22, Y: 0.65}},
			{FromID: "true_branch", ToID: "target", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.7, Y: 0.22}, {X: 0.7, Y: 0.5}}},
			{FromID: "false_branch", ToID: "target", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.7, Y: 0.78}, {X: 0.7, Y: 0.5}}},
		},
	})

	// State Loop: Circular pattern
	//      ┌──────────────┐
	//      │              ▼
	//  [A]───>[B]───>[C]──┘
	r.Register(&PatternLayout{
		Type:      PatternStateLoop,
		Name:      "State Loop",
		MinWidth:  500,
		MinHeight: 250,
		Padding:   40,
		Positions: []LayoutPosition{
			{ID: "start", X: 0.15, Y: 0.65, Width: 0.18, Height: 0.3},
			{ID: "middle", X: 0.5, Y: 0.65, Width: 0.18, Height: 0.3},
			{ID: "end", X: 0.85, Y: 0.65, Width: 0.18, Height: 0.3},
		},
		Edges: []EdgePath{
			{FromID: "start", ToID: "middle", CurveStyle: "straight"},
			{FromID: "middle", ToID: "end", CurveStyle: "straight"},
			{FromID: "end", ToID: "start", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.85, Y: 0.25}, {X: 0.15, Y: 0.25}},
				LabelPos:  Point{X: 0.5, Y: 0.18}},
		},
	})

	// Star Topology: Central hub with spokes
	//       [S1]
	//        │
	//  [S2]─[C]─[S3]
	//        │
	//       [S4]
	r.Register(&PatternLayout{
		Type:      PatternStarTopology,
		Name:      "Star Topology",
		MinWidth:  450,
		MinHeight: 450,
		Padding:   40,
		Positions: []LayoutPosition{
			{ID: "center", X: 0.5, Y: 0.5, Width: 0.22, Height: 0.22},
			{ID: "top", X: 0.5, Y: 0.12, Width: 0.18, Height: 0.18},
			{ID: "right", X: 0.88, Y: 0.5, Width: 0.18, Height: 0.18},
			{ID: "bottom", X: 0.5, Y: 0.88, Width: 0.18, Height: 0.18},
			{ID: "left", X: 0.12, Y: 0.5, Width: 0.18, Height: 0.18},
		},
		Edges: []EdgePath{
			{FromID: "center", ToID: "top", CurveStyle: "straight"},
			{FromID: "center", ToID: "right", CurveStyle: "straight"},
			{FromID: "center", ToID: "bottom", CurveStyle: "straight"},
			{FromID: "center", ToID: "left", CurveStyle: "straight"},
		},
	})
}

func (r *PatternRegistry) registerFlowPatterns() {
	// If-Else: Diamond decision with two branches
	//          ┌──>[Process1]──┐
	//  ◇──────┤                ├──>◯
	//          └──>[Process2]──┘
	r.Register(&PatternLayout{
		Type:      PatternIfElse,
		Name:      "If-Else",
		MinWidth:  550,
		MinHeight: 350,
		Padding:   40,
		Positions: []LayoutPosition{
			{ID: "decision", X: 0.12, Y: 0.5, Width: 0.12, Height: 0.2},
			{ID: "true_process", X: 0.45, Y: 0.22, Width: 0.2, Height: 0.18},
			{ID: "false_process", X: 0.45, Y: 0.78, Width: 0.2, Height: 0.18},
			{ID: "merge", X: 0.88, Y: 0.5, Width: 0.12, Height: 0.18},
		},
		Edges: []EdgePath{
			{FromID: "decision", ToID: "true_process", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.25, Y: 0.5}, {X: 0.25, Y: 0.22}},
				LabelPos:  Point{X: 0.18, Y: 0.35}},
			{FromID: "decision", ToID: "false_process", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.25, Y: 0.5}, {X: 0.25, Y: 0.78}},
				LabelPos:  Point{X: 0.18, Y: 0.65}},
			{FromID: "true_process", ToID: "merge", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.7, Y: 0.22}, {X: 0.7, Y: 0.5}}},
			{FromID: "false_process", ToID: "merge", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.7, Y: 0.78}, {X: 0.7, Y: 0.5}}},
		},
	})

	// If-ElseIf-Else: Cascading decisions
	//  ◇──>[P1]──────────────┐
	//  │                     │
	//  ◇──>[P2]──────────────┼──>◯
	//  │                     │
	//  └──>[P3]──────────────┘
	r.Register(&PatternLayout{
		Type:      PatternIfElseIfElse,
		Name:      "If-ElseIf-Else",
		MinWidth:  650,
		MinHeight: 450,
		Padding:   40,
		Positions: []LayoutPosition{
			{ID: "decision1", X: 0.1, Y: 0.18, Width: 0.1, Height: 0.15},
			{ID: "decision2", X: 0.1, Y: 0.5, Width: 0.1, Height: 0.15},
			{ID: "process1", X: 0.4, Y: 0.18, Width: 0.18, Height: 0.14},
			{ID: "process2", X: 0.4, Y: 0.5, Width: 0.18, Height: 0.14},
			{ID: "process3", X: 0.4, Y: 0.82, Width: 0.18, Height: 0.14},
			{ID: "merge", X: 0.85, Y: 0.5, Width: 0.1, Height: 0.15},
		},
		Edges: []EdgePath{
			{FromID: "decision1", ToID: "process1", CurveStyle: "straight",
				LabelPos: Point{X: 0.25, Y: 0.12}},
			{FromID: "decision1", ToID: "decision2", CurveStyle: "straight"},
			{FromID: "decision2", ToID: "process2", CurveStyle: "straight",
				LabelPos: Point{X: 0.25, Y: 0.44}},
			{FromID: "decision2", ToID: "process3", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.1, Y: 0.82}},
				LabelPos:  Point{X: 0.05, Y: 0.66}},
			{FromID: "process1", ToID: "merge", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.65, Y: 0.18}, {X: 0.65, Y: 0.5}}},
			{FromID: "process2", ToID: "merge", CurveStyle: "straight"},
			{FromID: "process3", ToID: "merge", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.65, Y: 0.82}, {X: 0.65, Y: 0.5}}},
		},
	})

	// While Loop: Loop back pattern
	//       ┌────────────┐
	//       │            │
	//  ──>◇────>[Process]─┘
	//       │
	//       └──>◯
	r.Register(&PatternLayout{
		Type:      PatternWhileLoop,
		Name:      "While Loop",
		MinWidth:  450,
		MinHeight: 300,
		Padding:   40,
		Positions: []LayoutPosition{
			{ID: "entry", X: 0.08, Y: 0.5, Width: 0.08, Height: 0.12},
			{ID: "condition", X: 0.28, Y: 0.5, Width: 0.12, Height: 0.18},
			{ID: "body", X: 0.6, Y: 0.5, Width: 0.2, Height: 0.18},
			{ID: "exit", X: 0.28, Y: 0.88, Width: 0.1, Height: 0.1},
		},
		Edges: []EdgePath{
			{FromID: "entry", ToID: "condition", CurveStyle: "straight"},
			{FromID: "condition", ToID: "body", CurveStyle: "straight",
				LabelPos: Point{X: 0.44, Y: 0.44}},
			{FromID: "body", ToID: "condition", CurveStyle: "orthogonal",
				Waypoints: []Point{{X: 0.6, Y: 0.2}, {X: 0.28, Y: 0.2}},
				LabelPos:  Point{X: 0.44, Y: 0.12}},
			{FromID: "condition", ToID: "exit", CurveStyle: "straight",
				LabelPos: Point{X: 0.2, Y: 0.7}},
		},
	})

	// Sequential: Simple linear flow
	//  (●)──>[P1]──>[P2]──>[P3]──>(◉)
	r.Register(&PatternLayout{
		Type:      PatternSequential,
		Name:      "Sequential",
		MinWidth:  650,
		MinHeight: 120,
		Padding:   30,
		Positions: []LayoutPosition{
			{ID: "start", X: 0.06, Y: 0.5, Width: 0.06, Height: 0.3},
			{ID: "process_0", X: 0.22, Y: 0.5, Width: 0.16, Height: 0.4},
			{ID: "process_1", X: 0.44, Y: 0.5, Width: 0.16, Height: 0.4},
			{ID: "process_2", X: 0.66, Y: 0.5, Width: 0.16, Height: 0.4},
			{ID: "process_3", X: 0.88, Y: 0.5, Width: 0.16, Height: 0.4},
			{ID: "end", X: 0.98, Y: 0.5, Width: 0.06, Height: 0.3},
		},
		Edges: []EdgePath{
			{FromID: "start", ToID: "process_0", CurveStyle: "straight"},
			{FromID: "process_0", ToID: "process_1", CurveStyle: "straight"},
			{FromID: "process_1", ToID: "process_2", CurveStyle: "straight"},
			{FromID: "process_2", ToID: "process_3", CurveStyle: "straight"},
			{FromID: "process_3", ToID: "end", CurveStyle: "straight"},
		},
	})
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
