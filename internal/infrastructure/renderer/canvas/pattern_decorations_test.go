package canvas

import (
	"strings"
	"testing"
)

// ============================================================
// Decoration Registry Tests
// ============================================================

func TestNewDecorationRegistry(t *testing.T) {
	r := NewDecorationRegistry()

	if r == nil {
		t.Fatal("NewDecorationRegistry returned nil")
	}

	if len(r.gradients) == 0 {
		t.Error("gradients should be registered")
	}
	if len(r.filters) == 0 {
		t.Error("filters should be registered")
	}
	if len(r.markers) == 0 {
		t.Error("markers should be registered")
	}
	if len(r.styles) == 0 {
		t.Error("styles should be registered")
	}
}

// ============================================================
// Gradient Tests
// ============================================================

func TestDecorationRegistry_Gradients(t *testing.T) {
	r := NewDecorationRegistry()

	expectedGradients := []string{
		"grad-header-blue",
		"grad-interface-purple",
		"grad-success-green",
		"grad-process-orange",
		"grad-bg-subtle",
		"grad-glow-blue",
		"grad-database-teal",
		"grad-error-red",
		"grad-terminal-dark",
		"grad-note-pink",
	}

	gradientIDs := make(map[string]bool)
	for _, g := range r.gradients {
		gradientIDs[g.ID] = true
	}

	for _, id := range expectedGradients {
		if !gradientIDs[id] {
			t.Errorf("missing gradient: %s", id)
		}
	}
}

func TestGradientDef_LinearGradient(t *testing.T) {
	r := NewDecorationRegistry()

	var linearGradient *GradientDef
	for i := range r.gradients {
		if r.gradients[i].Type == "linear" {
			linearGradient = &r.gradients[i]
			break
		}
	}

	if linearGradient == nil {
		t.Fatal("no linear gradient found")
	}

	if len(linearGradient.Stops) < 2 {
		t.Error("linear gradient should have at least 2 stops")
	}

	if linearGradient.Attrs == nil {
		t.Error("linear gradient should have attributes")
	}
}

func TestGradientDef_RadialGradient(t *testing.T) {
	r := NewDecorationRegistry()

	var radialGradient *GradientDef
	for i := range r.gradients {
		if r.gradients[i].Type == "radial" {
			radialGradient = &r.gradients[i]
			break
		}
	}

	if radialGradient == nil {
		t.Fatal("no radial gradient found")
	}

	if len(radialGradient.Stops) < 2 {
		t.Error("radial gradient should have at least 2 stops")
	}
}

// ============================================================
// Filter Tests
// ============================================================

func TestDecorationRegistry_Filters(t *testing.T) {
	r := NewDecorationRegistry()

	expectedFilters := []string{
		"filter-shadow-soft",
		"filter-shadow-elevated",
		"filter-inner-shadow",
		"filter-glow-blue",
		"filter-glow-green",
		"filter-glow-red",
		"filter-emboss",
		"filter-blur-bg",
	}

	filterIDs := make(map[string]bool)
	for _, f := range r.filters {
		filterIDs[f.ID] = true
	}

	for _, id := range expectedFilters {
		if !filterIDs[id] {
			t.Errorf("missing filter: %s", id)
		}
	}
}

func TestFilterDef_HasContent(t *testing.T) {
	r := NewDecorationRegistry()

	for _, f := range r.filters {
		if f.Content == "" {
			t.Errorf("filter %s has empty content", f.ID)
		}
	}
}

// ============================================================
// Marker Tests
// ============================================================

func TestDecorationRegistry_Markers(t *testing.T) {
	r := NewDecorationRegistry()

	expectedMarkers := []string{
		"marker-arrow-filled",
		"marker-arrow-open",
		"marker-triangle",
		"marker-diamond-filled",
		"marker-diamond-open",
		"marker-circle",
		"marker-arrow-async",
	}

	markerIDs := make(map[string]bool)
	for _, m := range r.markers {
		markerIDs[m.ID] = true
	}

	for _, id := range expectedMarkers {
		if !markerIDs[id] {
			t.Errorf("missing marker: %s", id)
		}
	}
}

func TestMarkerDef_ValidDimensions(t *testing.T) {
	r := NewDecorationRegistry()

	for _, m := range r.markers {
		if m.Width <= 0 || m.Height <= 0 {
			t.Errorf("marker %s has invalid dimensions: %dx%d", m.ID, m.Width, m.Height)
		}
		if m.ViewBox == "" {
			t.Errorf("marker %s has empty ViewBox", m.ID)
		}
		if m.Content == "" {
			t.Errorf("marker %s has empty Content", m.ID)
		}
	}
}

// ============================================================
// Style Tests
// ============================================================

func TestDecorationRegistry_Styles(t *testing.T) {
	r := NewDecorationRegistry()

	if len(r.styles) == 0 {
		t.Fatal("no styles registered")
	}

	// Check for expected CSS classes
	allStyles := strings.Join(r.styles, "")

	expectedClasses := []string{
		".diagram-text",
		".node-class",
		".node-interface",
		".node-state",
		".node-process",
		".node-decision",
		".node-terminal",
		".edge-solid",
		".edge-dashed",
		".note-box",
		".fragment-box",
		".swimlane-header",
		".activation-bar",
		".lifeline",
	}

	for _, class := range expectedClasses {
		if !strings.Contains(allStyles, class) {
			t.Errorf("missing CSS class: %s", class)
		}
	}
}

// ============================================================
// GenerateDefs Tests
// ============================================================

func TestDecorationRegistry_GenerateDefs(t *testing.T) {
	r := NewDecorationRegistry()
	defs := r.GenerateDefs()

	if defs == "" {
		t.Fatal("GenerateDefs returned empty string")
	}

	// Should contain gradients
	if !strings.Contains(defs, "<linearGradient") {
		t.Error("defs should contain linearGradient")
	}
	if !strings.Contains(defs, "<radialGradient") {
		t.Error("defs should contain radialGradient")
	}

	// Should contain filters
	if !strings.Contains(defs, "<filter") {
		t.Error("defs should contain filter")
	}

	// Should contain markers
	if !strings.Contains(defs, "<marker") {
		t.Error("defs should contain marker")
	}

	// Verify gradient structure
	if !strings.Contains(defs, "<stop") {
		t.Error("gradients should contain stop elements")
	}
}

// ============================================================
// GenerateStyles Tests
// ============================================================

func TestDecorationRegistry_GenerateStyles(t *testing.T) {
	r := NewDecorationRegistry()
	css := r.GenerateStyles()

	if css == "" {
		t.Fatal("GenerateStyles returned empty string")
	}

	// Should contain font-family definitions
	if !strings.Contains(css, "font-family") {
		t.Error("styles should contain font-family")
	}

	// Should contain fill and stroke
	if !strings.Contains(css, "fill") {
		t.Error("styles should contain fill")
	}
	if !strings.Contains(css, "stroke") {
		t.Error("styles should contain stroke")
	}
}

// ============================================================
// ApplyToCanvas Tests
// ============================================================

func TestDecorationRegistry_ApplyToCanvas(t *testing.T) {
	r := NewDecorationRegistry()
	c := New()

	r.ApplyToCanvas(c)

	// Canvas should now have defs
	svg := c.String()

	if !strings.Contains(svg, "<defs>") {
		t.Error("canvas should contain defs")
	}
	if !strings.Contains(svg, "<style>") {
		t.Error("canvas should contain style")
	}
}

// ============================================================
// Helper Function Tests
// ============================================================

func TestNodeStyleClass(t *testing.T) {
	tests := []struct {
		nodeType string
		want     string
	}{
		{"class", "node-class"},
		{"interface", "node-interface"},
		{"abstract", "node-abstract"},
		{"state", "node-state"},
		{"process", "node-process"},
		{"decision", "node-decision"},
		{"terminal", "node-terminal"},
		{"unknown", "node-class"},
		{"", "node-class"},
	}

	for _, tt := range tests {
		t.Run(tt.nodeType, func(t *testing.T) {
			got := NodeStyleClass(tt.nodeType)
			if got != tt.want {
				t.Errorf("NodeStyleClass(%q) = %q, want %q", tt.nodeType, got, tt.want)
			}
		})
	}
}

func TestEdgeStyleClass(t *testing.T) {
	tests := []struct {
		edgeType string
		want     string
	}{
		{"solid", "edge-solid"},
		{"inheritance", "edge-solid"},
		{"composition", "edge-solid"},
		{"aggregation", "edge-solid"},
		{"dashed", "edge-dashed"},
		{"dependency", "edge-dashed"},
		{"implementation", "edge-dashed"},
		{"dotted", "edge-dotted"},
		{"async", "edge-async"},
		{"return", "edge-return"},
		{"unknown", "edge-solid"},
		{"", "edge-solid"},
	}

	for _, tt := range tests {
		t.Run(tt.edgeType, func(t *testing.T) {
			got := EdgeStyleClass(tt.edgeType)
			if got != tt.want {
				t.Errorf("EdgeStyleClass(%q) = %q, want %q", tt.edgeType, got, tt.want)
			}
		})
	}
}

func TestGradientForNodeType(t *testing.T) {
	tests := []struct {
		nodeType string
		want     string
	}{
		{"interface", "grad-interface-purple"},
		{"database", "grad-database-teal"},
		{"process", "grad-process-orange"},
		{"terminal", "grad-terminal-dark"},
		{"success", "grad-success-green"},
		{"final", "grad-success-green"},
		{"error", "grad-error-red"},
		{"class", "grad-bg-subtle"},
		{"unknown", "grad-bg-subtle"},
	}

	for _, tt := range tests {
		t.Run(tt.nodeType, func(t *testing.T) {
			got := GradientForNodeType(tt.nodeType)
			if got != tt.want {
				t.Errorf("GradientForNodeType(%q) = %q, want %q", tt.nodeType, got, tt.want)
			}
		})
	}
}

func TestFilterForState(t *testing.T) {
	tests := []struct {
		state string
		want  string
	}{
		{"hover", "filter-glow-blue"},
		{"active", "filter-glow-blue"},
		{"success", "filter-glow-green"},
		{"error", "filter-glow-red"},
		{"elevated", "filter-shadow-elevated"},
		{"normal", "filter-shadow-soft"},
		{"", "filter-shadow-soft"},
	}

	for _, tt := range tests {
		t.Run(tt.state, func(t *testing.T) {
			got := FilterForState(tt.state)
			if got != tt.want {
				t.Errorf("FilterForState(%q) = %q, want %q", tt.state, got, tt.want)
			}
		})
	}
}

func TestMarkerForEdgeType(t *testing.T) {
	tests := []struct {
		edgeType string
		want     string
	}{
		{"inheritance", "marker-triangle"},
		{"generalization", "marker-triangle"},
		{"implementation", "marker-triangle"},
		{"realization", "marker-triangle"},
		{"composition", "marker-diamond-filled"},
		{"aggregation", "marker-diamond-open"},
		{"dependency", "marker-arrow-open"},
		{"async", "marker-arrow-async"},
		{"association", "marker-arrow-filled"},
		{"", "marker-arrow-filled"},
	}

	for _, tt := range tests {
		t.Run(tt.edgeType, func(t *testing.T) {
			got := MarkerForEdgeType(tt.edgeType)
			if got != tt.want {
				t.Errorf("MarkerForEdgeType(%q) = %q, want %q", tt.edgeType, got, tt.want)
			}
		})
	}
}

// ============================================================
// GradientStop Tests
// ============================================================

func TestGradientStop_Fields(t *testing.T) {
	stop := GradientStop{
		Offset:  "50%",
		Color:   "#ff0000",
		Opacity: "0.8",
	}

	if stop.Offset != "50%" {
		t.Error("Offset not set correctly")
	}
	if stop.Color != "#ff0000" {
		t.Error("Color not set correctly")
	}
	if stop.Opacity != "0.8" {
		t.Error("Opacity not set correctly")
	}
}

// ============================================================
// Integration Test
// ============================================================

func TestDecorationRegistry_FullIntegration(t *testing.T) {
	// Create registry and canvas
	r := NewDecorationRegistry()
	c := New()
	c.SetSize(800, 600)

	// Apply decorations
	r.ApplyToCanvas(c)

	// Add some elements that use the decorations
	c.Rect(100, 100, 200, 100, Fill("url(#grad-header-blue)"), Class("node-class"))
	c.Line(100, 250, 300, 250, Stroke("#000"), Class("edge-solid"))
	c.Circle(400, 300, 50, Fill("url(#grad-success-green)"))

	svg := c.String()

	// Verify the output is valid SVG
	if !strings.HasPrefix(svg, "<svg") {
		t.Error("output should start with <svg")
	}
	if !strings.HasSuffix(svg, "</svg>") {
		t.Error("output should end with </svg>")
	}

	// Verify decorations are present
	if !strings.Contains(svg, "grad-header-blue") {
		t.Error("SVG should reference gradient")
	}
	if !strings.Contains(svg, "node-class") {
		t.Error("SVG should reference CSS class")
	}
}
