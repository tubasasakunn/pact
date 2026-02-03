package validator

import (
	"pact/internal/domain/ast"
	"pact/internal/domain/errors"
)

// validDurationUnits は有効な期間単位
var validDurationUnits = map[string]bool{
	"ms": true, "s": true, "m": true, "h": true, "d": true,
}

// validateStatesDuplicates はStates名の重複をチェックする
func (v *Validator) validateStatesDuplicates(comp *ast.ComponentDecl) {
	statesNames := make(map[string]ast.Position)
	for _, states := range comp.Body.States {
		if pos, exists := statesNames[states.Name]; exists {
			v.errors.Add(&errors.ValidationError{
				Pos:     states.Pos,
				Type:    "duplicate",
				Name:    states.Name,
				Message: "states already defined at " + pos.String(),
			})
		} else {
			statesNames[states.Name] = states.Pos
		}

		// State名の重複チェック
		v.validateDuplicateStates(&states)
	}
}

// validateDuplicateStates はState名の重複をチェックする
func (v *Validator) validateDuplicateStates(states *ast.StatesDecl) {
	stateNames := make(map[string]ast.Position)
	for _, state := range states.States {
		if pos, exists := stateNames[state.Name]; exists {
			v.errors.Add(&errors.ValidationError{
				Pos:     state.Pos,
				Type:    "duplicate",
				Name:    state.Name,
				Message: "state already defined at " + pos.String(),
			})
		} else {
			stateNames[state.Name] = state.Pos
		}
	}
}

// validateStateReferences はStates内の状態参照を検証する
func (v *Validator) validateStateReferences(states *ast.StatesDecl) {
	// 定義された状態名を収集
	definedStates := make(map[string]bool)
	for _, state := range states.States {
		definedStates[state.Name] = true
	}

	// initial状態を検証
	if states.Initial != "" && !definedStates[states.Initial] {
		v.errors.Add(&errors.ValidationError{
			Pos:     states.Pos,
			Type:    "undefined",
			Name:    states.Initial,
			Message: "initial state is not defined",
		})
	}

	// final状態を検証
	for _, final := range states.Finals {
		if !definedStates[final] {
			v.errors.Add(&errors.ValidationError{
				Pos:     states.Pos,
				Type:    "undefined",
				Name:    final,
				Message: "final state is not defined",
			})
		}
	}

	// 遷移の参照を検証
	for _, trans := range states.Transitions {
		if !definedStates[trans.From] {
			v.errors.Add(&errors.ValidationError{
				Pos:     trans.Pos,
				Type:    "undefined",
				Name:    trans.From,
				Message: "transition source state is not defined",
			})
		}
		if !definedStates[trans.To] {
			v.errors.Add(&errors.ValidationError{
				Pos:     trans.Pos,
				Type:    "undefined",
				Name:    trans.To,
				Message: "transition target state is not defined",
			})
		}
	}
}

// validateStateReferencesInComp はコンポーネント内のStates参照を検証する
func (v *Validator) validateStateReferencesInComp(comp *ast.ComponentDecl) {
	for _, states := range comp.Body.States {
		v.validateStateReferences(&states)
	}
}

// ValidateDurationUnits は状態遷移のduration単位を検証する（M-003）
func (v *Validator) ValidateDurationUnits(spec *ast.SpecFile) error {
	v.errors = &errors.MultiError{}

	for _, comp := range spec.Components {
		for _, states := range comp.Body.States {
			v.validateTransitionDurations(states.Transitions)
			for _, state := range states.States {
				v.validateTransitionDurations(state.Transitions)
			}
			for _, parallel := range states.Parallels {
				for _, region := range parallel.Regions {
					v.validateTransitionDurations(region.Transitions)
				}
			}
		}
	}

	return v.errors.ErrorOrNil()
}

// validateTransitionDurations は遷移のduration単位を検証する
func (v *Validator) validateTransitionDurations(transitions []ast.TransitionDecl) {
	for _, trans := range transitions {
		if trigger, ok := trans.Trigger.(*ast.AfterTrigger); ok {
			if !validDurationUnits[trigger.Duration.Unit] {
				v.errors.Add(&errors.ValidationError{
					Pos:     trans.Pos,
					Type:    "invalid",
					Name:    trigger.Duration.Unit,
					Message: "invalid duration unit (valid units: ms, s, m, h, d)",
				})
			}
		}
	}
}

// validateEmptyStates は空のstates宣言を検出する
func (v *Validator) validateEmptyStates(comp *ast.ComponentDecl) {
	for _, states := range comp.Body.States {
		if len(states.States) == 0 {
			v.errors.Add(&errors.ValidationError{
				Pos:     states.Pos,
				Type:    "warning",
				Name:    states.Name,
				Message: "empty states declaration",
			})
		}
	}
}
