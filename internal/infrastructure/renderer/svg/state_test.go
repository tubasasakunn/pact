package svg

import (
	"bytes"
	"strings"
	"testing"

	"pact/internal/domain/diagram/state"
)

// =============================================================================
// RST001-RST014: StateRenderer Tests
// =============================================================================

// RST001: 空図
func TestStateRenderer_EmptyDiagram(t *testing.T) {
	diagram := &state.Diagram{}

	renderer := NewStateRenderer()
	var buf bytes.Buffer
	err := renderer.Render(diagram, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	if !strings.Contains(svg, "<svg") {
		t.Error("expected valid SVG output")
	}
}

// RST002: 初期状態
func TestStateRenderer_InitialState(t *testing.T) {
	diagram := &state.Diagram{
		States: []state.State{
			{ID: "initial", Name: "Start", Type: state.StateTypeInitial},
		},
	}

	renderer := NewStateRenderer()
	var buf bytes.Buffer
	renderer.Render(diagram, &buf)

	svg := buf.String()
	// 初期状態のシンボルテンプレートが使用されている
	if !strings.Contains(svg, `initial-state`) {
		t.Error("expected initial-state template for initial state")
	}
}

// RST003: 最終状態
func TestStateRenderer_FinalState(t *testing.T) {
	diagram := &state.Diagram{
		States: []state.State{
			{ID: "final", Name: "End", Type: state.StateTypeFinal},
		},
	}

	renderer := NewStateRenderer()
	var buf bytes.Buffer
	renderer.Render(diagram, &buf)

	svg := buf.String()
	// 二重丸（2つの円）
	circleCount := strings.Count(svg, "<circle")
	if circleCount < 2 {
		t.Errorf("expected 2 circles for final state, got %d", circleCount)
	}
}

// RST004: 原子状態
func TestStateRenderer_AtomicState(t *testing.T) {
	diagram := &state.Diagram{
		States: []state.State{
			{ID: "pending", Name: "Pending", Type: state.StateTypeAtomic},
		},
	}

	renderer := NewStateRenderer()
	var buf bytes.Buffer
	renderer.Render(diagram, &buf)

	svg := buf.String()
	// 角丸矩形
	if !strings.Contains(svg, "<rect") {
		t.Error("expected rect for atomic state")
	}
}

// RST005: 状態名
func TestStateRenderer_StateWithName(t *testing.T) {
	diagram := &state.Diagram{
		States: []state.State{
			{ID: "processing", Name: "Processing", Type: state.StateTypeAtomic},
		},
	}

	renderer := NewStateRenderer()
	var buf bytes.Buffer
	renderer.Render(diagram, &buf)

	svg := buf.String()
	if !strings.Contains(svg, "Processing") {
		t.Error("expected state name in output")
	}
}
