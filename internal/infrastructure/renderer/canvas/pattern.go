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
