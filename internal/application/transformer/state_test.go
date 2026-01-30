package transformer

import (
	"testing"

	"pact/internal/domain/ast"
	"pact/internal/domain/diagram/state"
	"pact/internal/domain/errors"
)

func createStateTestComponent(states ast.StatesDecl) *ast.SpecFile {
	return &ast.SpecFile{
		Component: &ast.ComponentDecl{
			Name: "TestService",
			Body: ast.ComponentBody{
				States: []ast.StatesDecl{states},
			},
		},
	}
}

// TST001: 最小ステートマシン
func TestStateTransformer_MinimalStates(t *testing.T) {
	statesDecl := ast.StatesDecl{
		Name:    "OrderState",
		Initial: "Pending",
	}
	files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

	transformer := NewStateTransformer()
	diagram, err := transformer.Transform(files, &StateOptions{StatesName: "OrderState"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Initial状態のみ
	if len(diagram.States) < 1 {
		t.Errorf("expected at least 1 state, got %d", len(diagram.States))
	}
}

// TST002: 初期状態
func TestStateTransformer_InitialState(t *testing.T) {
	statesDecl := ast.StatesDecl{
		Name:    "OrderState",
		Initial: "Pending",
	}
	files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

	transformer := NewStateTransformer()
	diagram, _ := transformer.Transform(files, &StateOptions{StatesName: "OrderState"})

	var initialState *state.State
	for i := range diagram.States {
		if diagram.States[i].Name == "Pending" {
			initialState = &diagram.States[i]
			break
		}
	}
	if initialState == nil {
		t.Fatal("expected initial state")
	}
	if initialState.Type != state.StateTypeInitial {
		t.Errorf("expected initial type, got %v", initialState.Type)
	}
}

// TST003: 最終状態
func TestStateTransformer_FinalState(t *testing.T) {
	statesDecl := ast.StatesDecl{
		Name:    "OrderState",
		Initial: "Pending",
		Finals:  []string{"Completed"},
	}
	files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

	transformer := NewStateTransformer()
	diagram, _ := transformer.Transform(files, &StateOptions{StatesName: "OrderState"})

	var finalState *state.State
	for i := range diagram.States {
		if diagram.States[i].Name == "Completed" {
			finalState = &diagram.States[i]
			break
		}
	}
	if finalState == nil {
		t.Fatal("expected final state")
	}
	if finalState.Type != state.StateTypeFinal {
		t.Errorf("expected final type, got %v", finalState.Type)
	}
}

// TST004: 原子状態
func TestStateTransformer_AtomicStates(t *testing.T) {
	statesDecl := ast.StatesDecl{
		Name:    "OrderState",
		Initial: "Pending",
		Transitions: []ast.TransitionDecl{
			{From: "Pending", To: "Processing", Trigger: &ast.EventTrigger{Event: "Start"}},
		},
	}
	files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

	transformer := NewStateTransformer()
	diagram, _ := transformer.Transform(files, &StateOptions{StatesName: "OrderState"})

	var processingState *state.State
	for i := range diagram.States {
		if diagram.States[i].Name == "Processing" {
			processingState = &diagram.States[i]
			break
		}
	}
	if processingState == nil {
		t.Fatal("expected Processing state")
	}
	if processingState.Type != state.StateTypeAtomic {
		t.Errorf("expected atomic type, got %v", processingState.Type)
	}
}

// TST005: イベントトリガー
func TestStateTransformer_Transition_Event(t *testing.T) {
	statesDecl := ast.StatesDecl{
		Name:    "OrderState",
		Initial: "Pending",
		Transitions: []ast.TransitionDecl{
			{From: "Pending", To: "Processing", Trigger: &ast.EventTrigger{Event: "Start"}},
		},
	}
	files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

	transformer := NewStateTransformer()
	diagram, _ := transformer.Transform(files, &StateOptions{StatesName: "OrderState"})

	if len(diagram.Transitions) != 1 {
		t.Fatalf("expected 1 transition, got %d", len(diagram.Transitions))
	}

	trigger, ok := diagram.Transitions[0].Trigger.(*state.EventTrigger)
	if !ok {
		t.Fatalf("expected EventTrigger, got %T", diagram.Transitions[0].Trigger)
	}
	if trigger.Event != "Start" {
		t.Errorf("expected event 'Start', got %q", trigger.Event)
	}
}

// TST006: 時間トリガー
func TestStateTransformer_Transition_After(t *testing.T) {
	statesDecl := ast.StatesDecl{
		Name:    "OrderState",
		Initial: "Pending",
		Transitions: []ast.TransitionDecl{
			{From: "Pending", To: "Timeout", Trigger: &ast.AfterTrigger{Duration: ast.Duration{Value: 30, Unit: "s"}}},
		},
	}
	files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

	transformer := NewStateTransformer()
	diagram, _ := transformer.Transform(files, &StateOptions{StatesName: "OrderState"})

	trigger, ok := diagram.Transitions[0].Trigger.(*state.AfterTrigger)
	if !ok {
		t.Fatalf("expected AfterTrigger, got %T", diagram.Transitions[0].Trigger)
	}
	if trigger.Duration.Value != 30 || trigger.Duration.Unit != "s" {
		t.Errorf("expected 30s, got %d%s", trigger.Duration.Value, trigger.Duration.Unit)
	}
}

// TST007: 条件トリガー
func TestStateTransformer_Transition_When(t *testing.T) {
	statesDecl := ast.StatesDecl{
		Name:    "OrderState",
		Initial: "Pending",
		Transitions: []ast.TransitionDecl{
			{From: "Pending", To: "Active", Trigger: &ast.WhenTrigger{Condition: &ast.VariableExpr{Name: "isReady"}}},
		},
	}
	files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

	transformer := NewStateTransformer()
	diagram, _ := transformer.Transform(files, &StateOptions{StatesName: "OrderState"})

	trigger, ok := diagram.Transitions[0].Trigger.(*state.WhenTrigger)
	if !ok {
		t.Fatalf("expected WhenTrigger, got %T", diagram.Transitions[0].Trigger)
	}
	if trigger.Condition != "isReady" {
		t.Errorf("expected condition 'isReady', got %q", trigger.Condition)
	}
}

// TST008: ガード
func TestStateTransformer_Transition_Guard(t *testing.T) {
	statesDecl := ast.StatesDecl{
		Name:    "OrderState",
		Initial: "Pending",
		Transitions: []ast.TransitionDecl{
			{
				From:    "Pending",
				To:      "Active",
				Trigger: &ast.EventTrigger{Event: "Start"},
				Guard:   &ast.VariableExpr{Name: "isValid"},
			},
		},
	}
	files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

	transformer := NewStateTransformer()
	diagram, _ := transformer.Transform(files, &StateOptions{StatesName: "OrderState"})

	if diagram.Transitions[0].Guard != "isValid" {
		t.Errorf("expected guard 'isValid', got %q", diagram.Transitions[0].Guard)
	}
}

// TST009: アクション
func TestStateTransformer_Transition_Actions(t *testing.T) {
	statesDecl := ast.StatesDecl{
		Name:    "OrderState",
		Initial: "Pending",
		Transitions: []ast.TransitionDecl{
			{
				From:    "Pending",
				To:      "Active",
				Trigger: &ast.EventTrigger{Event: "Start"},
				Actions: []string{"notifyUser", "logEvent"},
			},
		},
	}
	files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

	transformer := NewStateTransformer()
	diagram, _ := transformer.Transform(files, &StateOptions{StatesName: "OrderState"})

	if len(diagram.Transitions[0].Actions) != 2 {
		t.Errorf("expected 2 actions, got %d", len(diagram.Transitions[0].Actions))
	}
}

// TST010: entry
func TestStateTransformer_State_Entry(t *testing.T) {
	statesDecl := ast.StatesDecl{
		Name:    "OrderState",
		Initial: "Pending",
		States: []ast.StateDecl{
			{Name: "Processing", Entry: []string{"startProcessing"}},
		},
	}
	files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

	transformer := NewStateTransformer()
	diagram, _ := transformer.Transform(files, &StateOptions{StatesName: "OrderState"})

	var processingState *state.State
	for i := range diagram.States {
		if diagram.States[i].Name == "Processing" {
			processingState = &diagram.States[i]
			break
		}
	}
	if processingState == nil {
		t.Fatal("expected Processing state")
	}
	if len(processingState.Entry) != 1 {
		t.Errorf("expected 1 entry action, got %d", len(processingState.Entry))
	}
}

// TST011: exit
func TestStateTransformer_State_Exit(t *testing.T) {
	statesDecl := ast.StatesDecl{
		Name:    "OrderState",
		Initial: "Pending",
		States: []ast.StateDecl{
			{Name: "Processing", Exit: []string{"cleanup"}},
		},
	}
	files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

	transformer := NewStateTransformer()
	diagram, _ := transformer.Transform(files, &StateOptions{StatesName: "OrderState"})

	var processingState *state.State
	for i := range diagram.States {
		if diagram.States[i].Name == "Processing" {
			processingState = &diagram.States[i]
			break
		}
	}
	if len(processingState.Exit) != 1 {
		t.Errorf("expected 1 exit action, got %d", len(processingState.Exit))
	}
}

// TST012: 階層状態
func TestStateTransformer_CompoundState(t *testing.T) {
	statesDecl := ast.StatesDecl{
		Name:    "OrderState",
		Initial: "Pending",
		States: []ast.StateDecl{
			{
				Name: "Active",
				States: []ast.StateDecl{
					{Name: "SubState1"},
					{Name: "SubState2"},
				},
			},
		},
	}
	files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

	transformer := NewStateTransformer()
	diagram, _ := transformer.Transform(files, &StateOptions{StatesName: "OrderState"})

	var activeState *state.State
	for i := range diagram.States {
		if diagram.States[i].Name == "Active" {
			activeState = &diagram.States[i]
			break
		}
	}
	if activeState == nil {
		t.Fatal("expected Active state")
	}
	if activeState.Type != state.StateTypeCompound {
		t.Errorf("expected compound type, got %v", activeState.Type)
	}
	if len(activeState.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(activeState.Children))
	}
}

// TST014: 並行状態
func TestStateTransformer_ParallelState(t *testing.T) {
	statesDecl := ast.StatesDecl{
		Name:    "OrderState",
		Initial: "Pending",
		Parallels: []ast.ParallelDecl{
			{
				Name: "ParallelState",
				Regions: []ast.RegionDecl{
					{Name: "Region1", Initial: "A"},
					{Name: "Region2", Initial: "B"},
				},
			},
		},
	}
	files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

	transformer := NewStateTransformer()
	diagram, _ := transformer.Transform(files, &StateOptions{StatesName: "OrderState"})

	var parallelState *state.State
	for i := range diagram.States {
		if diagram.States[i].Name == "ParallelState" {
			parallelState = &diagram.States[i]
			break
		}
	}
	if parallelState == nil {
		t.Fatal("expected ParallelState")
	}
	if parallelState.Type != state.StateTypeParallel {
		t.Errorf("expected parallel type, got %v", parallelState.Type)
	}
	if len(parallelState.Regions) != 2 {
		t.Errorf("expected 2 regions, got %d", len(parallelState.Regions))
	}
}

// TST016: 未発見
func TestStateTransformer_StatesNotFound(t *testing.T) {
	statesDecl := ast.StatesDecl{Name: "OrderState", Initial: "Pending"}
	files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

	transformer := NewStateTransformer()
	_, err := transformer.Transform(files, &StateOptions{StatesName: "NonExistent"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if _, ok := err.(*errors.TransformError); !ok {
		t.Errorf("expected TransformError, got %T", err)
	}
}

// TST018: 時間単位
func TestStateTransformer_DurationUnits(t *testing.T) {
	units := []string{"ms", "s", "m", "h", "d"}

	for _, unit := range units {
		statesDecl := ast.StatesDecl{
			Name:    "OrderState",
			Initial: "Pending",
			Transitions: []ast.TransitionDecl{
				{From: "Pending", To: "Timeout", Trigger: &ast.AfterTrigger{Duration: ast.Duration{Value: 5, Unit: unit}}},
			},
		}
		files := []*ast.SpecFile{createStateTestComponent(statesDecl)}

		transformer := NewStateTransformer()
		diagram, _ := transformer.Transform(files, &StateOptions{StatesName: "OrderState"})

		trigger := diagram.Transitions[0].Trigger.(*state.AfterTrigger)
		if trigger.Duration.Unit != unit {
			t.Errorf("expected unit %q, got %q", unit, trigger.Duration.Unit)
		}
	}
}
