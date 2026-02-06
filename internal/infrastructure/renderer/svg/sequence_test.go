package svg

import (
	"bytes"
	"strings"
	"testing"

	"pact/internal/domain/diagram/sequence"
)

// =============================================================================
// RSQ001-RSQ015: SequenceRenderer Tests
// =============================================================================

// RSQ001: 空図
func TestSeqRenderer_EmptyDiagram(t *testing.T) {
	diagram := &sequence.Diagram{}

	renderer := NewSequenceRenderer()
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

// RSQ002: 単一参加者
func TestSeqRenderer_SingleParticipant(t *testing.T) {
	diagram := &sequence.Diagram{
		Participants: []sequence.Participant{
			{ID: "User", Name: "User", Type: sequence.ParticipantTypeDefault},
		},
	}

	renderer := NewSequenceRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(diagram, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	// 矩形とライフライン
	if !strings.Contains(svg, "<rect") {
		t.Error("expected rect for participant")
	}
	lineCount := strings.Count(svg, "<line")
	if lineCount < 1 {
		t.Error("expected line for lifeline")
	}
}

// RSQ003: 複数参加者
func TestSeqRenderer_MultipleParticipants(t *testing.T) {
	diagram := &sequence.Diagram{
		Participants: []sequence.Participant{
			{ID: "A", Name: "A", Type: sequence.ParticipantTypeDefault},
			{ID: "B", Name: "B", Type: sequence.ParticipantTypeDefault},
			{ID: "C", Name: "C", Type: sequence.ParticipantTypeDefault},
		},
	}

	renderer := NewSequenceRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(diagram, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	rectCount := strings.Count(svg, "<rect")
	if rectCount < 3 {
		t.Errorf("expected at least 3 rect elements, got %d", rectCount)
	}
}

// RSQ004: アクター
func TestSeqRenderer_ParticipantType_Actor(t *testing.T) {
	diagram := &sequence.Diagram{
		Participants: []sequence.Participant{
			{ID: "User", Name: "User", Type: sequence.ParticipantTypeActor},
		},
	}

	renderer := NewSequenceRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(diagram, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	// 人型は円と線で構成
	if !strings.Contains(svg, "<circle") {
		t.Error("expected circle for actor head")
	}
}

// RSQ005: データベース
func TestSeqRenderer_ParticipantType_Database(t *testing.T) {
	diagram := &sequence.Diagram{
		Participants: []sequence.Participant{
			{ID: "DB", Name: "DB", Type: sequence.ParticipantTypeDatabase},
		},
	}

	renderer := NewSequenceRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(diagram, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svg := buf.String()
	// 円柱は楕円で構成
	if !strings.Contains(svg, "<ellipse") {
		t.Error("expected ellipse for database cylinder")
	}
}
