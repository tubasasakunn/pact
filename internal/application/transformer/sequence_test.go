package transformer

import (
	"testing"

	"pact/internal/domain/ast"
	"pact/internal/domain/diagram/sequence"
	"pact/internal/domain/errors"
)

// =============================================================================
// TS001-TS019: SequenceTransformer Tests
// =============================================================================

func createTestComponent(flowSteps []ast.Step, relations []ast.RelationDecl) *ast.SpecFile {
	return &ast.SpecFile{
		Component: &ast.ComponentDecl{
			Name: "TestService",
			Body: ast.ComponentBody{
				Relations: relations,
				Flows: []ast.FlowDecl{
					{Name: "Process", Steps: flowSteps},
				},
			},
		},
	}
}

// TS001: 空フロー
func TestSeqTransformer_EmptyFlow(t *testing.T) {
	files := []*ast.SpecFile{createTestComponent([]ast.Step{}, nil)}

	transformer := NewSequenceTransformer()
	diagram, err := transformer.Transform(files, &SequenceOptions{FlowName: "Process"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(diagram.Participants) != 1 {
		t.Errorf("expected 1 participant, got %d", len(diagram.Participants))
	}
	if len(diagram.Events) != 0 {
		t.Errorf("expected 0 events, got %d", len(diagram.Events))
	}
}

// TS002: 単一呼び出し
func TestSeqTransformer_SingleCall(t *testing.T) {
	steps := []ast.Step{
		&ast.CallStep{
			Expr: &ast.CallExpr{
				Object: &ast.VariableExpr{Name: "DB"},
				Method: "Query",
			},
		},
	}
	files := []*ast.SpecFile{createTestComponent(steps, []ast.RelationDecl{
		{Kind: ast.RelationDependsOn, Target: "DB"},
	})}

	transformer := NewSequenceTransformer()
	diagram, err := transformer.Transform(files, &SequenceOptions{FlowName: "Process"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(diagram.Participants) != 2 {
		t.Errorf("expected 2 participants, got %d", len(diagram.Participants))
	}
	if len(diagram.Events) != 1 {
		t.Errorf("expected 1 event, got %d", len(diagram.Events))
	}
}

// TS003: 戻りメッセージ
func TestSeqTransformer_CallWithReturn(t *testing.T) {
	steps := []ast.Step{
		&ast.AssignStep{
			Variable: "result",
			Value: &ast.CallExpr{
				Object: &ast.VariableExpr{Name: "DB"},
				Method: "Query",
			},
		},
	}
	files := []*ast.SpecFile{createTestComponent(steps, []ast.RelationDecl{
		{Kind: ast.RelationDependsOn, Target: "DB"},
	})}

	transformer := NewSequenceTransformer()
	diagram, err := transformer.Transform(files, &SequenceOptions{FlowName: "Process", IncludeReturn: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 呼び出し + 戻り
	if len(diagram.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(diagram.Events))
	}

	returnEvent := diagram.Events[1].(*sequence.MessageEvent)
	if returnEvent.MessageType != sequence.MessageTypeReturn {
		t.Errorf("expected return message type")
	}
}

// TS004: 非同期呼び出し
func TestSeqTransformer_AsyncCall(t *testing.T) {
	steps := []ast.Step{
		&ast.CallStep{
			Await: true,
			Expr: &ast.CallExpr{
				Object: &ast.VariableExpr{Name: "Queue"},
				Method: "Send",
			},
		},
	}
	files := []*ast.SpecFile{createTestComponent(steps, []ast.RelationDecl{
		{Kind: ast.RelationDependsOn, Target: "Queue"},
	})}

	transformer := NewSequenceTransformer()
	diagram, _ := transformer.Transform(files, &SequenceOptions{FlowName: "Process"})

	msg := diagram.Events[0].(*sequence.MessageEvent)
	if msg.MessageType != sequence.MessageTypeAsync {
		t.Errorf("expected async message type, got %v", msg.MessageType)
	}
}

// TS005: 同期呼び出し
func TestSeqTransformer_SyncCall(t *testing.T) {
	steps := []ast.Step{
		&ast.CallStep{
			Await: false,
			Expr: &ast.CallExpr{
				Object: &ast.VariableExpr{Name: "Service"},
				Method: "Call",
			},
		},
	}
	files := []*ast.SpecFile{createTestComponent(steps, []ast.RelationDecl{
		{Kind: ast.RelationDependsOn, Target: "Service"},
	})}

	transformer := NewSequenceTransformer()
	diagram, _ := transformer.Transform(files, &SequenceOptions{FlowName: "Process"})

	msg := diagram.Events[0].(*sequence.MessageEvent)
	if msg.MessageType != sequence.MessageTypeSync {
		t.Errorf("expected sync message type, got %v", msg.MessageType)
	}
}

// TS006: 複数参加者
func TestSeqTransformer_MultipleParticipants(t *testing.T) {
	relations := []ast.RelationDecl{
		{Kind: ast.RelationDependsOn, Target: "ServiceA"},
		{Kind: ast.RelationDependsOn, Target: "ServiceB"},
		{Kind: ast.RelationDependsOn, Target: "Database"},
	}
	files := []*ast.SpecFile{createTestComponent([]ast.Step{}, relations)}

	transformer := NewSequenceTransformer()
	diagram, _ := transformer.Transform(files, &SequenceOptions{FlowName: "Process"})

	// 自身 + 3つの依存先
	if len(diagram.Participants) != 4 {
		t.Errorf("expected 4 participants, got %d", len(diagram.Participants))
	}
}

// TS007-TS010: 参加者タイプ
func TestSeqTransformer_ParticipantTypes(t *testing.T) {
	tests := []struct {
		targetType string
		expected   sequence.ParticipantType
	}{
		{"database", sequence.ParticipantTypeDatabase},
		{"external", sequence.ParticipantTypeExternal},
		{"queue", sequence.ParticipantTypeQueue},
		{"actor", sequence.ParticipantTypeActor},
	}

	for _, tt := range tests {
		targetType := tt.targetType
		relations := []ast.RelationDecl{
			{Kind: ast.RelationDependsOn, Target: "Target", TargetType: &targetType},
		}
		files := []*ast.SpecFile{createTestComponent([]ast.Step{}, relations)}

		transformer := NewSequenceTransformer()
		diagram, _ := transformer.Transform(files, &SequenceOptions{FlowName: "Process"})

		var found *sequence.Participant
		for i := range diagram.Participants {
			if diagram.Participants[i].Name == "Target" {
				found = &diagram.Participants[i]
				break
			}
		}
		if found == nil {
			t.Errorf("expected Target participant")
			continue
		}
		if found.Type != tt.expected {
			t.Errorf("for %s: expected %v, got %v", tt.targetType, tt.expected, found.Type)
		}
	}
}

// TS011: 条件分岐
func TestSeqTransformer_IfStep(t *testing.T) {
	steps := []ast.Step{
		&ast.IfStep{
			Condition: &ast.VariableExpr{Name: "cond"},
			Then:      []ast.Step{},
		},
	}
	files := []*ast.SpecFile{createTestComponent(steps, nil)}

	transformer := NewSequenceTransformer()
	diagram, _ := transformer.Transform(files, &SequenceOptions{FlowName: "Process"})

	if len(diagram.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(diagram.Events))
	}

	fragment, ok := diagram.Events[0].(*sequence.FragmentEvent)
	if !ok {
		t.Fatalf("expected FragmentEvent, got %T", diagram.Events[0])
	}
	if fragment.Type != sequence.FragmentTypeAlt {
		t.Errorf("expected alt fragment, got %v", fragment.Type)
	}
}

// TS012: if-else
func TestSeqTransformer_IfElse(t *testing.T) {
	steps := []ast.Step{
		&ast.IfStep{
			Condition: &ast.VariableExpr{Name: "cond"},
			Then:      []ast.Step{},
			Else:      []ast.Step{},
		},
	}
	files := []*ast.SpecFile{createTestComponent(steps, nil)}

	transformer := NewSequenceTransformer()
	diagram, _ := transformer.Transform(files, &SequenceOptions{FlowName: "Process"})

	fragment := diagram.Events[0].(*sequence.FragmentEvent)
	if fragment.AltLabel != "else" {
		t.Errorf("expected else label")
	}
}

// TS013: forループ
func TestSeqTransformer_ForStep(t *testing.T) {
	steps := []ast.Step{
		&ast.ForStep{
			Variable: "item",
			Iterable: &ast.VariableExpr{Name: "items"},
			Body:     []ast.Step{},
		},
	}
	files := []*ast.SpecFile{createTestComponent(steps, nil)}

	transformer := NewSequenceTransformer()
	diagram, _ := transformer.Transform(files, &SequenceOptions{FlowName: "Process"})

	fragment := diagram.Events[0].(*sequence.FragmentEvent)
	if fragment.Type != sequence.FragmentTypeLoop {
		t.Errorf("expected loop fragment, got %v", fragment.Type)
	}
}

// TS014: whileループ
func TestSeqTransformer_WhileStep(t *testing.T) {
	steps := []ast.Step{
		&ast.WhileStep{
			Condition: &ast.VariableExpr{Name: "cond"},
			Body:      []ast.Step{},
		},
	}
	files := []*ast.SpecFile{createTestComponent(steps, nil)}

	transformer := NewSequenceTransformer()
	diagram, _ := transformer.Transform(files, &SequenceOptions{FlowName: "Process"})

	fragment := diagram.Events[0].(*sequence.FragmentEvent)
	if fragment.Type != sequence.FragmentTypeLoop {
		t.Errorf("expected loop fragment, got %v", fragment.Type)
	}
}

// TS017: フロー未発見
func TestSeqTransformer_FlowNotFound(t *testing.T) {
	files := []*ast.SpecFile{createTestComponent([]ast.Step{}, nil)}

	transformer := NewSequenceTransformer()
	_, err := transformer.Transform(files, &SequenceOptions{FlowName: "NonExistent"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if _, ok := err.(*errors.TransformError); !ok {
		t.Errorf("expected TransformError, got %T", err)
	}
}

// TS016: returnなしオプション
func TestSeqTransformer_OptionReturnMessages(t *testing.T) {
	steps := []ast.Step{
		&ast.AssignStep{
			Variable: "result",
			Value: &ast.CallExpr{
				Object: &ast.VariableExpr{Name: "DB"},
				Method: "Query",
			},
		},
	}
	files := []*ast.SpecFile{createTestComponent(steps, []ast.RelationDecl{
		{Kind: ast.RelationDependsOn, Target: "DB"},
	})}

	transformer := NewSequenceTransformer()
	diagram, _ := transformer.Transform(files, &SequenceOptions{FlowName: "Process", IncludeReturn: false})

	// returnなしなので1イベントのみ
	if len(diagram.Events) != 1 {
		t.Errorf("expected 1 event, got %d", len(diagram.Events))
	}
}
