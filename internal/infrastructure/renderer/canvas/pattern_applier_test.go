package canvas

import (
	"testing"
)

// ============================================================
// Pattern Layout Applier Tests
// ============================================================

func TestNewPatternLayoutApplier(t *testing.T) {
	r := NewPatternRegistry()
	a := NewPatternLayoutApplier(r)

	if a == nil {
		t.Fatal("NewPatternLayoutApplier returned nil")
	}
	if a.registry != r {
		t.Error("registry not set correctly")
	}
}

func TestPatternLayoutApplier_ApplyClassPattern(t *testing.T) {
	r := NewPatternRegistry()
	a := NewPatternLayoutApplier(r)

	match := ClassPatternMatch{
		Pattern: PatternInheritanceTree,
		NodeRoles: map[string]string{
			"parent":  "animal",
			"child_0": "dog",
			"child_1": "cat",
			"child_2": "bird",
		},
		Score: 0.9,
	}

	nodeWidths := map[string]int{
		"animal": 120,
		"dog":    100,
		"cat":    100,
		"bird":   100,
	}
	nodeHeights := map[string]int{
		"animal": 80,
		"dog":    60,
		"cat":    60,
		"bird":   60,
	}

	layout := a.ApplyClassPattern(match, nodeWidths, nodeHeights)

	if layout == nil {
		t.Fatal("ApplyClassPattern returned nil")
	}

	if layout.Pattern != PatternInheritanceTree {
		t.Errorf("Pattern = %s, want inheritance-tree", layout.Pattern)
	}

	if layout.Width < 400 {
		t.Errorf("Width = %d, want >= 400", layout.Width)
	}

	if len(layout.Nodes) != 4 {
		t.Errorf("expected 4 nodes, got %d", len(layout.Nodes))
	}

	// Verify node positions are within canvas bounds
	for _, node := range layout.Nodes {
		if node.X < 0 || node.X+node.Width > layout.Width {
			t.Errorf("node %s X position out of bounds: X=%d, Width=%d, CanvasWidth=%d",
				node.ID, node.X, node.Width, layout.Width)
		}
		if node.Y < 0 || node.Y+node.Height > layout.Height {
			t.Errorf("node %s Y position out of bounds: Y=%d, Height=%d, CanvasHeight=%d",
				node.ID, node.Y, node.Height, layout.Height)
		}
	}
}

func TestPatternLayoutApplier_ApplyStatePattern(t *testing.T) {
	r := NewPatternRegistry()
	a := NewPatternLayoutApplier(r)

	match := StatePatternMatch{
		Pattern: PatternLinearStates,
		StateRoles: map[string]string{
			"initial": "init",
			"state_0": "idle",
			"state_1": "running",
			"final":   "done",
		},
		Score: 0.85,
	}

	stateWidths := map[string]int{
		"init":    20,
		"idle":    100,
		"running": 100,
		"done":    24,
	}
	stateHeights := map[string]int{
		"init":    20,
		"idle":    50,
		"running": 50,
		"done":    24,
	}

	layout := a.ApplyStatePattern(match, stateWidths, stateHeights)

	if layout == nil {
		t.Fatal("ApplyStatePattern returned nil")
	}

	if layout.Pattern != PatternLinearStates {
		t.Errorf("Pattern = %s, want linear-states", layout.Pattern)
	}
}

func TestPatternLayoutApplier_ApplyFlowPattern(t *testing.T) {
	r := NewPatternRegistry()
	a := NewPatternLayoutApplier(r)

	match := FlowPatternMatch{
		Pattern: PatternIfElse,
		NodeRoles: map[string]string{
			"decision":      "dec",
			"true_process":  "yes",
			"false_process": "no",
			"merge":         "cont",
		},
		Score: 1.0,
	}

	nodeWidths := map[string]int{
		"dec":  80,
		"yes":  120,
		"no":   120,
		"cont": 100,
	}
	nodeHeights := map[string]int{
		"dec":  80,
		"yes":  60,
		"no":   60,
		"cont": 60,
	}

	layout := a.ApplyFlowPattern(match, nodeWidths, nodeHeights)

	if layout == nil {
		t.Fatal("ApplyFlowPattern returned nil")
	}

	if layout.Pattern != PatternIfElse {
		t.Errorf("Pattern = %s, want if-else", layout.Pattern)
	}

	// Should have edges connecting the nodes
	if len(layout.Edges) < 4 {
		t.Errorf("expected at least 4 edges, got %d", len(layout.Edges))
	}
}

func TestPatternLayoutApplier_ApplySequencePattern(t *testing.T) {
	r := NewPatternRegistry()
	a := NewPatternLayoutApplier(r)

	match := SequencePatternMatch{
		Pattern: PatternRequestResponse,
		ParticipantRoles: map[string]string{
			"caller": "client",
			"callee": "server",
		},
		Score: 1.0,
	}

	participantWidths := map[string]int{
		"client": 100,
		"server": 100,
	}

	layout := a.ApplySequencePattern(match, participantWidths)

	if layout == nil {
		t.Fatal("ApplySequencePattern returned nil")
	}

	if layout.Pattern != PatternRequestResponse {
		t.Errorf("Pattern = %s, want request-response", layout.Pattern)
	}
}

func TestPatternLayoutApplier_UnknownPattern(t *testing.T) {
	r := NewPatternRegistry()
	a := NewPatternLayoutApplier(r)

	match := ClassPatternMatch{
		Pattern:   PatternType("unknown"),
		NodeRoles: map[string]string{},
		Score:     0.5,
	}

	layout := a.ApplyClassPattern(match, nil, nil)

	if layout != nil {
		t.Error("expected nil for unknown pattern")
	}
}

// ============================================================
// Layout Calculation Tests
// ============================================================

func TestPatternLayoutApplier_CalculateCanvasSize(t *testing.T) {
	r := NewPatternRegistry()
	a := NewPatternLayoutApplier(r)

	layout := r.Get(PatternInheritanceTree)
	if layout == nil {
		t.Fatal("inheritance tree pattern not found")
	}

	// Test with small nodes - should use minimum size
	smallWidths := map[string]int{"parent": 50, "child_0": 40}
	smallHeights := map[string]int{"parent": 30, "child_0": 25}
	roles := map[string]string{"parent": "a", "child_0": "b"}

	w, h := a.calculateCanvasSize(layout, smallWidths, smallHeights, roles)

	if w < layout.MinWidth {
		t.Errorf("Width = %d, should be >= MinWidth %d", w, layout.MinWidth)
	}
	if h < layout.MinHeight {
		t.Errorf("Height = %d, should be >= MinHeight %d", h, layout.MinHeight)
	}
}

func TestPatternLayoutApplier_CalculateConnectionPoint(t *testing.T) {
	a := NewPatternLayoutApplier(NewPatternRegistry())

	node := NodeLayout{
		ID:     "test",
		X:      100,
		Y:      100,
		Width:  80,
		Height: 60,
	}

	tests := []struct {
		name    string
		targetX int
		targetY int
		wantX   int
		wantY   int
	}{
		{"target to right", 300, 130, 180, 130}, // Right edge
		{"target to left", 0, 130, 100, 130},    // Left edge
		{"target below", 140, 300, 140, 160},    // Bottom edge
		{"target above", 140, 0, 140, 100},      // Top edge
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := a.calculateConnectionPoint(node, tt.targetX, tt.targetY)
			if gotX != tt.wantX || gotY != tt.wantY {
				t.Errorf("calculateConnectionPoint() = (%d, %d), want (%d, %d)",
					gotX, gotY, tt.wantX, tt.wantY)
			}
		})
	}
}

// ============================================================
// NodeLayout Tests
// ============================================================

func TestNodeLayout_Fields(t *testing.T) {
	nl := NodeLayout{
		ID:     "test-node",
		X:      50,
		Y:      100,
		Width:  120,
		Height: 80,
	}

	if nl.ID != "test-node" {
		t.Error("ID not set correctly")
	}
	if nl.X != 50 || nl.Y != 100 {
		t.Error("position not set correctly")
	}
	if nl.Width != 120 || nl.Height != 80 {
		t.Error("size not set correctly")
	}
}

// ============================================================
// EdgeLayout Tests
// ============================================================

func TestEdgeLayout_Fields(t *testing.T) {
	el := EdgeLayout{
		FromID: "a",
		ToID:   "b",
		Waypoints: []AbsolutePoint{
			{X: 100, Y: 100},
			{X: 150, Y: 150},
			{X: 200, Y: 100},
		},
		LabelX: 150,
		LabelY: 125,
	}

	if el.FromID != "a" || el.ToID != "b" {
		t.Error("endpoints not set correctly")
	}
	if len(el.Waypoints) != 3 {
		t.Errorf("expected 3 waypoints, got %d", len(el.Waypoints))
	}
	if el.LabelX != 150 || el.LabelY != 125 {
		t.Error("label position not set correctly")
	}
}

// ============================================================
// AbsolutePoint Tests
// ============================================================

func TestAbsolutePoint_Fields(t *testing.T) {
	p := AbsolutePoint{X: 42, Y: 84}

	if p.X != 42 || p.Y != 84 {
		t.Errorf("AbsolutePoint = (%d, %d), want (42, 84)", p.X, p.Y)
	}
}

// ============================================================
// AppliedDecorator Tests
// ============================================================

func TestAppliedDecorator_Fields(t *testing.T) {
	d := AppliedDecorator{
		Type:   "divider",
		X:      50,
		Y:      200,
		Width:  400,
		Height: 2,
		Style: map[string]string{
			"stroke":           "#ccc",
			"stroke-dasharray": "5,5",
		},
	}

	if d.Type != "divider" {
		t.Error("Type not set correctly")
	}
	if d.Style["stroke"] != "#ccc" {
		t.Error("Style not set correctly")
	}
}

// ============================================================
// Helper Function Tests
// ============================================================

func TestAbs(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{-100, 100},
	}

	for _, tt := range tests {
		got := abs(tt.input)
		if got != tt.want {
			t.Errorf("abs(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
