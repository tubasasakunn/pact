package transformer

import (
	"pact/internal/domain/ast"
	"pact/internal/domain/diagram/state"
	"pact/internal/domain/errors"
)

// StateTransformer はASTを状態図に変換する
type StateTransformer struct{}

// NewStateTransformer は新しいStateTransformerを作成する
func NewStateTransformer() *StateTransformer {
	return &StateTransformer{}
}

// StateOptions は状態図変換オプション
type StateOptions struct {
	StatesName string
}

// Transform はASTを状態図に変換する
func (t *StateTransformer) Transform(files []*ast.SpecFile, opts *StateOptions) (*state.Diagram, error) {
	if opts == nil || opts.StatesName == "" {
		return nil, &errors.TransformError{Source: "AST", Target: "StateDiagram", Message: "states name is required"}
	}

	var targetStates *ast.StatesDecl

	for _, file := range files {
		// 単一コンポーネント
		if file.Component != nil {
			for i := range file.Component.Body.States {
				if file.Component.Body.States[i].Name == opts.StatesName {
					targetStates = &file.Component.Body.States[i]
					break
				}
			}
		}

		// 複数コンポーネント
		if targetStates == nil {
			for j := range file.Components {
				comp := &file.Components[j]
				for i := range comp.Body.States {
					if comp.Body.States[i].Name == opts.StatesName {
						targetStates = &comp.Body.States[i]
						break
					}
				}
				if targetStates != nil {
					break
				}
			}
		}
	}

	if targetStates == nil {
		return nil, &errors.TransformError{
			Source:  "AST",
			Target:  "StateDiagram",
			Message: "states not found: " + opts.StatesName,
		}
	}

	diagram := &state.Diagram{
		States:      []state.State{},
		Transitions: []state.Transition{},
	}

	// 初期疑似状態（黒丸）を追加
	initialPseudoID := "__initial__"
	if targetStates.Initial != "" {
		diagram.States = append(diagram.States, state.State{
			ID:   initialPseudoID,
			Name: "",
			Type: state.StateTypeInitial,
		})
		// 初期疑似状態から最初の状態への遷移
		diagram.Transitions = append(diagram.Transitions, state.Transition{
			From: initialPseudoID,
			To:   targetStates.Initial,
		})
	}

	// Final状態
	for _, final := range targetStates.Finals {
		diagram.States = append(diagram.States, state.State{
			ID:   final,
			Name: final,
			Type: state.StateTypeFinal,
		})
	}

	// 既に追加された状態名を記録
	stateNames := make(map[string]bool)
	for _, final := range targetStates.Finals {
		stateNames[final] = true
	}

	// 通常の状態（重複を避ける）
	for _, s := range targetStates.States {
		if !stateNames[s.Name] {
			diagram.States = append(diagram.States, t.transformState(&s))
			stateNames[s.Name] = true
		}
	}

	// 遷移から状態を収集
	for _, trans := range targetStates.Transitions {
		if !stateNames[trans.From] {
			diagram.States = append(diagram.States, state.State{
				ID:   trans.From,
				Name: trans.From,
				Type: state.StateTypeAtomic,
			})
			stateNames[trans.From] = true
		}
		if !stateNames[trans.To] {
			diagram.States = append(diagram.States, state.State{
				ID:   trans.To,
				Name: trans.To,
				Type: state.StateTypeAtomic,
			})
			stateNames[trans.To] = true
		}
	}

	// 遷移
	for _, trans := range targetStates.Transitions {
		diagram.Transitions = append(diagram.Transitions, t.transformTransition(&trans))
	}

	// 並行状態
	for _, p := range targetStates.Parallels {
		diagram.States = append(diagram.States, t.transformParallel(&p))
	}

	return diagram, nil
}

func (t *StateTransformer) transformState(s *ast.StateDecl) state.State {
	st := state.State{
		ID:    s.Name,
		Name:  s.Name,
		Type:  state.StateTypeAtomic,
		Entry: s.Entry,
		Exit:  s.Exit,
	}

	if len(s.States) > 0 {
		st.Type = state.StateTypeCompound
		for _, child := range s.States {
			st.Children = append(st.Children, t.transformState(&child))
		}
	}

	return st
}

func (t *StateTransformer) transformTransition(trans *ast.TransitionDecl) state.Transition {
	tr := state.Transition{
		From:    trans.From,
		To:      trans.To,
		Actions: trans.Actions,
	}

	if trans.Guard != nil {
		if v, ok := trans.Guard.(*ast.VariableExpr); ok {
			tr.Guard = v.Name
		}
	}

	if trans.Trigger != nil {
		switch trigger := trans.Trigger.(type) {
		case *ast.EventTrigger:
			tr.Trigger = &state.EventTrigger{Event: trigger.Event}
		case *ast.AfterTrigger:
			tr.Trigger = &state.AfterTrigger{
				Duration: state.Duration{
					Value: trigger.Duration.Value,
					Unit:  trigger.Duration.Unit,
				},
			}
		case *ast.WhenTrigger:
			if v, ok := trigger.Condition.(*ast.VariableExpr); ok {
				tr.Trigger = &state.WhenTrigger{Condition: v.Name}
			}
		}
	}

	return tr
}

func (t *StateTransformer) transformParallel(p *ast.ParallelDecl) state.State {
	st := state.State{
		ID:   p.Name,
		Name: p.Name,
		Type: state.StateTypeParallel,
	}

	for _, r := range p.Regions {
		region := state.Region{
			Name: r.Name,
		}
		for _, s := range r.States {
			region.States = append(region.States, t.transformState(&s))
		}
		for _, trans := range r.Transitions {
			region.Transitions = append(region.Transitions, t.transformTransition(&trans))
		}
		st.Regions = append(st.Regions, region)
	}

	return st
}
