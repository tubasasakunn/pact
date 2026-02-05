package canvas

import (
	"testing"

	"pact/internal/domain/diagram/class"
	"pact/internal/domain/diagram/flow"
	"pact/internal/domain/diagram/sequence"
	"pact/internal/domain/diagram/state"
)

// ============================================================
// Class Pattern Detector Tests
// ============================================================

func TestNewClassPatternDetector(t *testing.T) {
	r := NewPatternRegistry()
	d := NewClassPatternDetector(r)

	if d == nil {
		t.Fatal("NewClassPatternDetector returned nil")
	}
}

func TestClassPatternDetector_DetectInheritanceTree(t *testing.T) {
	r := NewPatternRegistry()
	d := NewClassPatternDetector(r)

	diagram := &class.Diagram{
		Nodes: []class.Node{
			{ID: "animal", Name: "Animal"},
			{ID: "dog", Name: "Dog"},
			{ID: "cat", Name: "Cat"},
			{ID: "bird", Name: "Bird"},
		},
		Edges: []class.Edge{
			{From: "dog", To: "animal", Type: class.EdgeTypeInheritance},
			{From: "cat", To: "animal", Type: class.EdgeTypeInheritance},
			{From: "bird", To: "animal", Type: class.EdgeTypeInheritance},
		},
	}

	matches := d.Detect(diagram)

	found := false
	for _, m := range matches {
		if m.Pattern == PatternInheritanceTree {
			found = true
			if m.NodeRoles["parent"] != "animal" {
				t.Errorf("expected parent=animal, got %s", m.NodeRoles["parent"])
			}
			if m.Score <= 0 {
				t.Error("score should be positive")
			}
		}
	}

	if !found {
		t.Error("expected to detect inheritance tree pattern")
	}
}

func TestClassPatternDetector_DetectInterfaceImpl(t *testing.T) {
	r := NewPatternRegistry()
	d := NewClassPatternDetector(r)

	diagram := &class.Diagram{
		Nodes: []class.Node{
			{ID: "reader", Name: "Reader", Stereotype: "interface"},
			{ID: "file_reader", Name: "FileReader"},
			{ID: "string_reader", Name: "StringReader"},
			{ID: "buf_reader", Name: "BufReader"},
		},
		Edges: []class.Edge{
			{From: "file_reader", To: "reader", Type: class.EdgeTypeImplementation},
			{From: "string_reader", To: "reader", Type: class.EdgeTypeImplementation},
			{From: "buf_reader", To: "reader", Type: class.EdgeTypeImplementation},
		},
	}

	matches := d.Detect(diagram)

	found := false
	for _, m := range matches {
		if m.Pattern == PatternInterfaceImpl {
			found = true
			if m.NodeRoles["interface"] != "reader" {
				t.Errorf("expected interface=reader, got %s", m.NodeRoles["interface"])
			}
		}
	}

	if !found {
		t.Error("expected to detect interface implementation pattern")
	}
}

func TestClassPatternDetector_DetectComposition(t *testing.T) {
	r := NewPatternRegistry()
	d := NewClassPatternDetector(r)

	diagram := &class.Diagram{
		Nodes: []class.Node{
			{ID: "car", Name: "Car"},
			{ID: "engine", Name: "Engine"},
			{ID: "wheel", Name: "Wheel"},
			{ID: "door", Name: "Door"},
		},
		Edges: []class.Edge{
			{From: "car", To: "engine", Type: class.EdgeTypeComposition},
			{From: "car", To: "wheel", Type: class.EdgeTypeComposition},
			{From: "car", To: "door", Type: class.EdgeTypeComposition},
		},
	}

	matches := d.Detect(diagram)

	found := false
	for _, m := range matches {
		if m.Pattern == PatternComposition {
			found = true
			if m.NodeRoles["owner"] != "car" {
				t.Errorf("expected owner=car, got %s", m.NodeRoles["owner"])
			}
		}
	}

	if !found {
		t.Error("expected to detect composition pattern")
	}
}

func TestClassPatternDetector_DetectDiamond(t *testing.T) {
	r := NewPatternRegistry()
	d := NewClassPatternDetector(r)

	// Diamond: A <- B, A <- C, B <- D, C <- D
	diagram := &class.Diagram{
		Nodes: []class.Node{
			{ID: "a", Name: "A"},
			{ID: "b", Name: "B"},
			{ID: "c", Name: "C"},
			{ID: "d", Name: "D"},
		},
		Edges: []class.Edge{
			{From: "b", To: "a", Type: class.EdgeTypeInheritance},
			{From: "c", To: "a", Type: class.EdgeTypeInheritance},
			{From: "d", To: "b", Type: class.EdgeTypeInheritance},
			{From: "d", To: "c", Type: class.EdgeTypeInheritance},
		},
	}

	matches := d.Detect(diagram)

	found := false
	for _, m := range matches {
		if m.Pattern == PatternDiamond {
			found = true
			if m.NodeRoles["top"] != "a" {
				t.Errorf("expected top=a, got %s", m.NodeRoles["top"])
			}
			if m.NodeRoles["bottom"] != "d" {
				t.Errorf("expected bottom=d, got %s", m.NodeRoles["bottom"])
			}
			if m.Score != 1.0 {
				t.Errorf("diamond should have score 1.0, got %f", m.Score)
			}
		}
	}

	if !found {
		t.Error("expected to detect diamond pattern")
	}
}

func TestClassPatternDetector_NoPattern(t *testing.T) {
	r := NewPatternRegistry()
	d := NewClassPatternDetector(r)

	// Simple diagram with no recognizable pattern
	diagram := &class.Diagram{
		Nodes: []class.Node{
			{ID: "a", Name: "A"},
			{ID: "b", Name: "B"},
		},
		Edges: []class.Edge{
			{From: "a", To: "b", Type: class.EdgeTypeDependency},
		},
	}

	matches := d.Detect(diagram)

	// Should return empty or minimal matches
	if len(matches) > 0 {
		for _, m := range matches {
			t.Logf("Detected: %s with score %f", m.Pattern, m.Score)
		}
	}
}

// ============================================================
// State Pattern Detector Tests
// ============================================================

func TestNewStatePatternDetector(t *testing.T) {
	r := NewPatternRegistry()
	d := NewStatePatternDetector(r)

	if d == nil {
		t.Fatal("NewStatePatternDetector returned nil")
	}
}

func TestStatePatternDetector_DetectLinear(t *testing.T) {
	r := NewPatternRegistry()
	d := NewStatePatternDetector(r)

	diagram := &state.Diagram{
		States: []state.State{
			{ID: "init", Name: "", Type: state.StateTypeInitial},
			{ID: "idle", Name: "Idle", Type: state.StateTypeAtomic},
			{ID: "running", Name: "Running", Type: state.StateTypeAtomic},
			{ID: "done", Name: "Done", Type: state.StateTypeAtomic},
			{ID: "final", Name: "", Type: state.StateTypeFinal},
		},
		Transitions: []state.Transition{
			{From: "init", To: "idle"},
			{From: "idle", To: "running"},
			{From: "running", To: "done"},
			{From: "done", To: "final"},
		},
	}

	matches := d.Detect(diagram)

	found := false
	for _, m := range matches {
		if m.Pattern == PatternLinearStates {
			found = true
			if m.StateRoles["initial"] != "init" {
				t.Errorf("expected initial=init, got %s", m.StateRoles["initial"])
			}
		}
	}

	if !found {
		t.Error("expected to detect linear states pattern")
	}
}

func TestStatePatternDetector_DetectBinaryChoice(t *testing.T) {
	r := NewPatternRegistry()
	d := NewStatePatternDetector(r)

	diagram := &state.Diagram{
		States: []state.State{
			{ID: "check", Name: "Check", Type: state.StateTypeAtomic},
			{ID: "valid", Name: "Valid", Type: state.StateTypeAtomic},
			{ID: "invalid", Name: "Invalid", Type: state.StateTypeAtomic},
			{ID: "next", Name: "Next", Type: state.StateTypeAtomic},
		},
		Transitions: []state.Transition{
			{From: "check", To: "valid"},
			{From: "check", To: "invalid"},
			{From: "valid", To: "next"},
			{From: "invalid", To: "next"},
		},
	}

	matches := d.Detect(diagram)

	found := false
	for _, m := range matches {
		if m.Pattern == PatternBinaryChoice {
			found = true
			if m.StateRoles["source"] != "check" {
				t.Errorf("expected source=check, got %s", m.StateRoles["source"])
			}
			if m.StateRoles["target"] != "next" {
				t.Errorf("expected target=next, got %s", m.StateRoles["target"])
			}
		}
	}

	if !found {
		t.Error("expected to detect binary choice pattern")
	}
}

func TestStatePatternDetector_DetectLoop(t *testing.T) {
	r := NewPatternRegistry()
	d := NewStatePatternDetector(r)

	diagram := &state.Diagram{
		States: []state.State{
			{ID: "a", Name: "A", Type: state.StateTypeAtomic},
			{ID: "b", Name: "B", Type: state.StateTypeAtomic},
			{ID: "c", Name: "C", Type: state.StateTypeAtomic},
		},
		Transitions: []state.Transition{
			{From: "a", To: "b"},
			{From: "b", To: "c"},
			{From: "c", To: "a"}, // Loop back
		},
	}

	matches := d.Detect(diagram)

	found := false
	for _, m := range matches {
		if m.Pattern == PatternStateLoop {
			found = true
		}
	}

	if !found {
		t.Error("expected to detect state loop pattern")
	}
}

// ============================================================
// Flow Pattern Detector Tests
// ============================================================

func TestNewFlowPatternDetector(t *testing.T) {
	r := NewPatternRegistry()
	d := NewFlowPatternDetector(r)

	if d == nil {
		t.Fatal("NewFlowPatternDetector returned nil")
	}
}

func TestFlowPatternDetector_DetectIfElse(t *testing.T) {
	r := NewPatternRegistry()
	d := NewFlowPatternDetector(r)

	diagram := &flow.Diagram{
		Nodes: []flow.Node{
			{ID: "dec", Label: "Is Valid?", Shape: flow.NodeShapeDecision},
			{ID: "yes", Label: "Process A", Shape: flow.NodeShapeProcess},
			{ID: "no", Label: "Process B", Shape: flow.NodeShapeProcess},
			{ID: "merge", Label: "Continue", Shape: flow.NodeShapeProcess},
		},
		Edges: []flow.Edge{
			{From: "dec", To: "yes", Label: "yes"},
			{From: "dec", To: "no", Label: "no"},
			{From: "yes", To: "merge"},
			{From: "no", To: "merge"},
		},
	}

	matches := d.Detect(diagram)

	found := false
	for _, m := range matches {
		if m.Pattern == PatternIfElse {
			found = true
			if m.NodeRoles["decision"] != "dec" {
				t.Errorf("expected decision=dec, got %s", m.NodeRoles["decision"])
			}
			if m.NodeRoles["merge"] != "merge" {
				t.Errorf("expected merge=merge, got %s", m.NodeRoles["merge"])
			}
		}
	}

	if !found {
		t.Error("expected to detect if-else pattern")
	}
}

func TestFlowPatternDetector_DetectWhileLoop(t *testing.T) {
	r := NewPatternRegistry()
	d := NewFlowPatternDetector(r)

	diagram := &flow.Diagram{
		Nodes: []flow.Node{
			{ID: "cond", Label: "Continue?", Shape: flow.NodeShapeDecision},
			{ID: "body", Label: "Process", Shape: flow.NodeShapeProcess},
			{ID: "exit", Label: "End", Shape: flow.NodeShapeTerminal},
		},
		Edges: []flow.Edge{
			{From: "cond", To: "body", Label: "yes"},
			{From: "cond", To: "exit", Label: "no"},
			{From: "body", To: "cond"}, // Loop back
		},
	}

	matches := d.Detect(diagram)

	found := false
	for _, m := range matches {
		if m.Pattern == PatternWhileLoop {
			found = true
			if m.NodeRoles["condition"] != "cond" {
				t.Errorf("expected condition=cond, got %s", m.NodeRoles["condition"])
			}
			if m.NodeRoles["body"] != "body" {
				t.Errorf("expected body=body, got %s", m.NodeRoles["body"])
			}
		}
	}

	if !found {
		t.Error("expected to detect while loop pattern")
	}
}

func TestFlowPatternDetector_DetectSequential(t *testing.T) {
	r := NewPatternRegistry()
	d := NewFlowPatternDetector(r)

	diagram := &flow.Diagram{
		Nodes: []flow.Node{
			{ID: "start", Label: "Start", Shape: flow.NodeShapeTerminal},
			{ID: "p1", Label: "Step 1", Shape: flow.NodeShapeProcess},
			{ID: "p2", Label: "Step 2", Shape: flow.NodeShapeProcess},
			{ID: "p3", Label: "Step 3", Shape: flow.NodeShapeProcess},
			{ID: "end", Label: "End", Shape: flow.NodeShapeTerminal},
		},
		Edges: []flow.Edge{
			{From: "start", To: "p1"},
			{From: "p1", To: "p2"},
			{From: "p2", To: "p3"},
			{From: "p3", To: "end"},
		},
	}

	matches := d.Detect(diagram)

	found := false
	for _, m := range matches {
		if m.Pattern == PatternSequential {
			found = true
			if m.NodeRoles["start"] != "start" {
				t.Errorf("expected start=start, got %s", m.NodeRoles["start"])
			}
		}
	}

	if !found {
		t.Error("expected to detect sequential pattern")
	}
}

// ============================================================
// Sequence Pattern Detector Tests
// ============================================================

func TestNewSequencePatternDetector(t *testing.T) {
	r := NewPatternRegistry()
	d := NewSequencePatternDetector(r)

	if d == nil {
		t.Fatal("NewSequencePatternDetector returned nil")
	}
}

func TestSequencePatternDetector_DetectRequestResponse(t *testing.T) {
	r := NewPatternRegistry()
	d := NewSequencePatternDetector(r)

	diagram := &sequence.Diagram{
		Participants: []sequence.Participant{
			{ID: "client", Name: "Client"},
			{ID: "server", Name: "Server"},
		},
		Events: []sequence.Event{
			&sequence.MessageEvent{From: "client", To: "server", Label: "request", MessageType: sequence.MessageTypeSync},
			&sequence.MessageEvent{From: "server", To: "client", Label: "response", MessageType: sequence.MessageTypeReturn},
		},
	}

	matches := d.Detect(diagram)

	found := false
	for _, m := range matches {
		if m.Pattern == PatternRequestResponse {
			found = true
			if m.ParticipantRoles["caller"] != "client" {
				t.Errorf("expected caller=client, got %s", m.ParticipantRoles["caller"])
			}
			if m.ParticipantRoles["callee"] != "server" {
				t.Errorf("expected callee=server, got %s", m.ParticipantRoles["callee"])
			}
		}
	}

	if !found {
		t.Error("expected to detect request-response pattern")
	}
}

func TestSequencePatternDetector_DetectChain(t *testing.T) {
	r := NewPatternRegistry()
	d := NewSequencePatternDetector(r)

	diagram := &sequence.Diagram{
		Participants: []sequence.Participant{
			{ID: "ui", Name: "UI"},
			{ID: "controller", Name: "Controller"},
			{ID: "service", Name: "Service"},
			{ID: "db", Name: "Database"},
		},
		Events: []sequence.Event{
			&sequence.MessageEvent{From: "ui", To: "controller", Label: "click", MessageType: sequence.MessageTypeSync},
			&sequence.MessageEvent{From: "controller", To: "service", Label: "process", MessageType: sequence.MessageTypeSync},
			&sequence.MessageEvent{From: "service", To: "db", Label: "query", MessageType: sequence.MessageTypeSync},
		},
	}

	matches := d.Detect(diagram)

	found := false
	for _, m := range matches {
		if m.Pattern == PatternChain {
			found = true
			if m.ParticipantRoles["node_0"] != "ui" {
				t.Errorf("expected node_0=ui, got %s", m.ParticipantRoles["node_0"])
			}
		}
	}

	if !found {
		t.Error("expected to detect chain pattern")
	}
}

// ============================================================
// Edge Cases
// ============================================================

func TestClassPatternDetector_EmptyDiagram(t *testing.T) {
	r := NewPatternRegistry()
	d := NewClassPatternDetector(r)

	diagram := &class.Diagram{}
	matches := d.Detect(diagram)

	if len(matches) != 0 {
		t.Errorf("expected no matches for empty diagram, got %d", len(matches))
	}
}

func TestStatePatternDetector_EmptyDiagram(t *testing.T) {
	r := NewPatternRegistry()
	d := NewStatePatternDetector(r)

	diagram := &state.Diagram{}
	matches := d.Detect(diagram)

	if len(matches) != 0 {
		t.Errorf("expected no matches for empty diagram, got %d", len(matches))
	}
}

func TestFlowPatternDetector_EmptyDiagram(t *testing.T) {
	r := NewPatternRegistry()
	d := NewFlowPatternDetector(r)

	diagram := &flow.Diagram{}
	matches := d.Detect(diagram)

	if len(matches) != 0 {
		t.Errorf("expected no matches for empty diagram, got %d", len(matches))
	}
}

func TestSequencePatternDetector_EmptyDiagram(t *testing.T) {
	r := NewPatternRegistry()
	d := NewSequencePatternDetector(r)

	diagram := &sequence.Diagram{}
	matches := d.Detect(diagram)

	if len(matches) != 0 {
		t.Errorf("expected no matches for empty diagram, got %d", len(matches))
	}
}
