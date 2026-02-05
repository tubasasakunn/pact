package canvas

import (
	"testing"
)

// ============================================================
// Pattern Registry Tests
// ============================================================

func TestNewPatternRegistry(t *testing.T) {
	r := NewPatternRegistry()
	if r == nil {
		t.Fatal("NewPatternRegistry returned nil")
	}
	if r.patterns == nil {
		t.Error("patterns map should be initialized")
	}
}

func TestPatternRegistry_Get(t *testing.T) {
	r := NewPatternRegistry()

	tests := []struct {
		name        string
		patternType PatternType
		wantNil     bool
	}{
		// Size variations
		{"inheritance tree 2 exists", PatternInheritanceTree2, false},
		{"inheritance tree 3 exists", PatternInheritanceTree3, false},
		{"inheritance tree 4 exists", PatternInheritanceTree4, false},
		{"interface impl 2 exists", PatternInterfaceImpl2, false},
		{"interface impl 3 exists", PatternInterfaceImpl3, false},
		{"interface impl 4 exists", PatternInterfaceImpl4, false},
		{"composition 2 exists", PatternComposition2, false},
		{"composition 3 exists", PatternComposition3, false},
		{"composition 4 exists", PatternComposition4, false},
		{"diamond exists", PatternDiamond, false},
		{"layered 3x2 exists", PatternLayered3x2, false},
		{"layered 3x3 exists", PatternLayered3x3, false},
		// Legacy aliases
		{"inheritance tree exists", PatternInheritanceTree, false},
		{"interface impl exists", PatternInterfaceImpl, false},
		{"composition exists", PatternComposition, false},
		{"layered exists", PatternLayered, false},
		// Sequence patterns
		{"request-response exists", PatternRequestResponse, false},
		{"callback exists", PatternCallback, false},
		{"chain 3 exists", PatternChain3, false},
		{"chain 4 exists", PatternChain4, false},
		{"chain exists", PatternChain, false},
		{"fan-out exists", PatternFanOut, false},
		// State patterns
		{"linear states 2 exists", PatternLinearStates2, false},
		{"linear states 3 exists", PatternLinearStates3, false},
		{"linear states 4 exists", PatternLinearStates4, false},
		{"linear states exists", PatternLinearStates, false},
		{"binary choice exists", PatternBinaryChoice, false},
		{"state loop exists", PatternStateLoop, false},
		{"star topology exists", PatternStarTopology, false},
		// Flow patterns
		{"if-else exists", PatternIfElse, false},
		{"if-elseif-else exists", PatternIfElseIfElse, false},
		{"while loop exists", PatternWhileLoop, false},
		{"sequential 3 exists", PatternSequential3, false},
		{"sequential 4 exists", PatternSequential4, false},
		{"sequential exists", PatternSequential, false},
		{"unknown pattern returns nil", PatternType("unknown"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.Get(tt.patternType)
			if (got == nil) != tt.wantNil {
				t.Errorf("Get(%q) returned nil=%v, want nil=%v", tt.patternType, got == nil, tt.wantNil)
			}
		})
	}
}

func TestPatternRegistry_Register(t *testing.T) {
	r := NewPatternRegistry()

	customPattern := &PatternLayout{
		Type:      PatternType("custom-test"),
		Name:      "Custom Test",
		MinWidth:  100,
		MinHeight: 100,
	}

	r.Register(customPattern)

	got := r.Get(PatternType("custom-test"))
	if got == nil {
		t.Fatal("registered pattern should be retrievable")
	}
	if got.Name != "Custom Test" {
		t.Errorf("Name = %q, want %q", got.Name, "Custom Test")
	}
}

// ============================================================
// Pattern Layout Tests
// ============================================================

func TestPatternLayout_InheritanceTree(t *testing.T) {
	r := NewPatternRegistry()
	layout := r.Get(PatternInheritanceTree)

	if layout == nil {
		t.Fatal("inheritance tree pattern not found")
	}

	if layout.MinWidth < 400 {
		t.Errorf("MinWidth = %d, want >= 400", layout.MinWidth)
	}

	if len(layout.Positions) < 2 {
		t.Errorf("should have at least 2 positions (parent + children)")
	}

	// Check for parent position
	hasParent := false
	for _, pos := range layout.Positions {
		if pos.ID == "parent" {
			hasParent = true
			if pos.Y > 0.5 {
				t.Error("parent should be in top half")
			}
		}
	}
	if !hasParent {
		t.Error("should have 'parent' position")
	}
}

func TestPatternLayout_IfElse(t *testing.T) {
	r := NewPatternRegistry()
	layout := r.Get(PatternIfElse)

	if layout == nil {
		t.Fatal("if-else pattern not found")
	}

	// Check for decision, true_process, false_process, merge
	requiredIDs := []string{"decision", "true_process", "false_process", "merge"}
	posIDs := make(map[string]bool)
	for _, pos := range layout.Positions {
		posIDs[pos.ID] = true
	}

	for _, id := range requiredIDs {
		if !posIDs[id] {
			t.Errorf("missing position: %s", id)
		}
	}

	// Check edges connect correctly
	if len(layout.Edges) < 4 {
		t.Errorf("should have at least 4 edges, got %d", len(layout.Edges))
	}
}

func TestPatternLayout_LinearStates(t *testing.T) {
	r := NewPatternRegistry()
	layout := r.Get(PatternLinearStates)

	if layout == nil {
		t.Fatal("linear states pattern not found")
	}

	// All positions should have similar Y (horizontal layout)
	var firstY float64
	for i, pos := range layout.Positions {
		if i == 0 {
			firstY = pos.Y
		} else {
			// Allow some variance but should be roughly horizontal
			if pos.Y < firstY-0.2 || pos.Y > firstY+0.2 {
				t.Logf("Note: position %s has Y=%f, first Y=%f", pos.ID, pos.Y, firstY)
			}
		}
	}
}

// ============================================================
// Pattern Matcher Tests
// ============================================================

func TestNewPatternMatcher(t *testing.T) {
	r := NewPatternRegistry()
	m := NewPatternMatcher(r)

	if m == nil {
		t.Fatal("NewPatternMatcher returned nil")
	}
	if m.registry != r {
		t.Error("registry not set correctly")
	}
}

// ============================================================
// Pattern Type Constants Tests
// ============================================================

func TestPatternTypeConstants(t *testing.T) {
	// Verify all pattern types are unique
	patterns := []PatternType{
		// Class patterns (size variations)
		PatternInheritanceTree2,
		PatternInheritanceTree3,
		PatternInheritanceTree4,
		PatternInterfaceImpl2,
		PatternInterfaceImpl3,
		PatternInterfaceImpl4,
		PatternComposition2,
		PatternComposition3,
		PatternComposition4,
		PatternDiamond,
		PatternLayered3x2,
		PatternLayered3x3,
		// Sequence patterns
		PatternRequestResponse,
		PatternCallback,
		PatternChain3,
		PatternChain4,
		PatternFanOut,
		// State patterns
		PatternLinearStates2,
		PatternLinearStates3,
		PatternLinearStates4,
		PatternBinaryChoice,
		PatternStateLoop,
		PatternStarTopology,
		// Flow patterns
		PatternIfElse,
		PatternIfElseIfElse,
		PatternWhileLoop,
		PatternSequential3,
		PatternSequential4,
	}

	seen := make(map[PatternType]bool)
	for _, p := range patterns {
		if seen[p] {
			t.Errorf("duplicate pattern type: %s", p)
		}
		seen[p] = true
	}
}

// ============================================================
// LayoutPosition Tests
// ============================================================

func TestLayoutPosition_RelativeCoordinates(t *testing.T) {
	r := NewPatternRegistry()

	for patternType, layout := range r.patterns {
		for _, pos := range layout.Positions {
			// X and Y should be 0.0 - 1.0 (relative)
			if pos.X < 0.0 || pos.X > 1.0 {
				t.Errorf("%s: position %s has X=%f outside [0,1]", patternType, pos.ID, pos.X)
			}
			if pos.Y < 0.0 || pos.Y > 1.0 {
				t.Errorf("%s: position %s has Y=%f outside [0,1]", patternType, pos.ID, pos.Y)
			}
			// Width and Height should be positive and reasonable
			if pos.Width <= 0 || pos.Width > 1.0 {
				t.Errorf("%s: position %s has invalid Width=%f", patternType, pos.ID, pos.Width)
			}
			if pos.Height <= 0 || pos.Height > 1.0 {
				t.Errorf("%s: position %s has invalid Height=%f", patternType, pos.ID, pos.Height)
			}
		}
	}
}

// ============================================================
// EdgePath Tests
// ============================================================

func TestEdgePath_CurveStyles(t *testing.T) {
	r := NewPatternRegistry()

	validStyles := map[string]bool{
		"straight":   true,
		"orthogonal": true,
		"curved":     true,
		"":           true, // empty is allowed
	}

	for patternType, layout := range r.patterns {
		for i, edge := range layout.Edges {
			if !validStyles[edge.CurveStyle] {
				t.Errorf("%s: edge %d has invalid CurveStyle=%q", patternType, i, edge.CurveStyle)
			}
		}
	}
}

// ============================================================
// GetBestMatch Tests
// ============================================================

func TestGetBestMatch_ClassPattern(t *testing.T) {
	matches := []ClassPatternMatch{
		{Pattern: PatternInheritanceTree, Score: 0.5},
		{Pattern: PatternDiamond, Score: 0.9},
		{Pattern: PatternComposition, Score: 0.3},
	}

	best, found := GetBestMatch(matches)
	if !found {
		t.Fatal("expected to find best match")
	}
	if best.Pattern != PatternDiamond {
		t.Errorf("expected diamond pattern (score 0.9), got %s", best.Pattern)
	}
}

func TestGetBestMatch_Empty(t *testing.T) {
	var matches []ClassPatternMatch
	_, found := GetBestMatch(matches)
	if found {
		t.Error("expected not to find match in empty slice")
	}
}

func TestGetBestMatch_StatePattern(t *testing.T) {
	matches := []StatePatternMatch{
		{Pattern: PatternLinearStates, Score: 0.7},
		{Pattern: PatternStateLoop, Score: 0.95},
	}

	best, found := GetBestMatch(matches)
	if !found {
		t.Fatal("expected to find best match")
	}
	if best.Pattern != PatternStateLoop {
		t.Errorf("expected state loop pattern, got %s", best.Pattern)
	}
}

// ============================================================
// Score Accessor Tests
// ============================================================

func TestClassPatternMatch_GetScore(t *testing.T) {
	m := ClassPatternMatch{Score: 0.75}
	if m.GetScore() != 0.75 {
		t.Errorf("GetScore() = %f, want 0.75", m.GetScore())
	}
}

func TestStatePatternMatch_GetScore(t *testing.T) {
	m := StatePatternMatch{Score: 0.8}
	if m.GetScore() != 0.8 {
		t.Errorf("GetScore() = %f, want 0.8", m.GetScore())
	}
}

func TestFlowPatternMatch_GetScore(t *testing.T) {
	m := FlowPatternMatch{Score: 0.65}
	if m.GetScore() != 0.65 {
		t.Errorf("GetScore() = %f, want 0.65", m.GetScore())
	}
}

func TestSequencePatternMatch_GetScore(t *testing.T) {
	m := SequencePatternMatch{Score: 0.9}
	if m.GetScore() != 0.9 {
		t.Errorf("GetScore() = %f, want 0.9", m.GetScore())
	}
}

// ============================================================
// Helper Function Tests
// ============================================================

func TestHelperRoleIDs(t *testing.T) {
	tests := []struct {
		fn   func(int) string
		i    int
		want string
	}{
		{childRoleID, 0, "child_0"},
		{childRoleID, 3, "child_3"},
		{implRoleID, 0, "impl_0"},
		{implRoleID, 2, "impl_2"},
		{partRoleID, 1, "part_1"},
		{stateRoleID, 0, "state_0"},
		{nodeRoleID, 3, "node_3"},
		{processRoleID, 2, "process_2"},
	}

	for _, tt := range tests {
		got := tt.fn(tt.i)
		if got != tt.want {
			t.Errorf("got %q, want %q", got, tt.want)
		}
	}
}

func TestItoa(t *testing.T) {
	tests := []struct {
		i    int
		want string
	}{
		{0, "0"},
		{1, "1"},
		{5, "5"},
		{9, "9"},
	}

	for _, tt := range tests {
		got := itoa(tt.i)
		if got != tt.want {
			t.Errorf("itoa(%d) = %q, want %q", tt.i, got, tt.want)
		}
	}
}
