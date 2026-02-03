package validator

import (
	"pact/internal/domain/ast"
	"pact/internal/domain/errors"
)

// MaxNestingDepth は式のネスト深度の上限
const MaxNestingDepth = 50

// validateFlowDuplicates はフロー名の重複をチェックする
func (v *Validator) validateFlowDuplicates(comp *ast.ComponentDecl) {
	flowNames := make(map[string]ast.Position)
	for _, flow := range comp.Body.Flows {
		if pos, exists := flowNames[flow.Name]; exists {
			v.errors.Add(&errors.ValidationError{
				Pos:     flow.Pos,
				Type:    "duplicate",
				Name:    flow.Name,
				Message: "flow already defined at " + pos.String(),
			})
		} else {
			flowNames[flow.Name] = flow.Pos
		}
	}
}

// ValidateExpressionDepth は式のネスト深度を検証する
func (v *Validator) ValidateExpressionDepth(spec *ast.SpecFile) error {
	v.errors = &errors.MultiError{}

	for _, comp := range spec.Components {
		for _, flow := range comp.Body.Flows {
			for _, step := range flow.Steps {
				v.checkStepDepth(step, 0)
			}
		}
	}

	return v.errors.ErrorOrNil()
}

// checkStepDepth はステップ内の式の深度をチェックする
func (v *Validator) checkStepDepth(step ast.Step, depth int) {
	if depth > MaxNestingDepth {
		v.errors.Add(&errors.ValidationError{
			Pos:     step.GetPos(),
			Type:    "invalid",
			Name:    "nesting",
			Message: "expression nesting depth exceeds maximum allowed",
		})
		return
	}

	switch s := step.(type) {
	case *ast.AssignStep:
		v.checkExprDepth(s.Value, s.Pos, 0)
	case *ast.CallStep:
		v.checkExprDepth(s.Expr, s.Pos, 0)
	case *ast.ReturnStep:
		if s.Value != nil {
			v.checkExprDepth(s.Value, s.Pos, 0)
		}
	case *ast.IfStep:
		v.checkExprDepth(s.Condition, s.Pos, 0)
		for _, thenStep := range s.Then {
			v.checkStepDepth(thenStep, depth+1)
		}
		for _, elseStep := range s.Else {
			v.checkStepDepth(elseStep, depth+1)
		}
	case *ast.ForStep:
		v.checkExprDepth(s.Iterable, s.Pos, 0)
		for _, bodyStep := range s.Body {
			v.checkStepDepth(bodyStep, depth+1)
		}
	case *ast.WhileStep:
		v.checkExprDepth(s.Condition, s.Pos, 0)
		for _, bodyStep := range s.Body {
			v.checkStepDepth(bodyStep, depth+1)
		}
	}
}

// checkExprDepth は式の深度をチェックする
func (v *Validator) checkExprDepth(expr ast.Expr, pos ast.Position, depth int) {
	if expr == nil {
		return
	}
	if depth > MaxNestingDepth {
		v.errors.Add(&errors.ValidationError{
			Pos:     pos,
			Type:    "invalid",
			Name:    "nesting",
			Message: "expression nesting depth exceeds maximum allowed",
		})
		return
	}

	switch e := expr.(type) {
	case *ast.BinaryExpr:
		v.checkExprDepth(e.Left, pos, depth+1)
		v.checkExprDepth(e.Right, pos, depth+1)
	case *ast.UnaryExpr:
		v.checkExprDepth(e.Operand, pos, depth+1)
	case *ast.TernaryExpr:
		v.checkExprDepth(e.Condition, pos, depth+1)
		v.checkExprDepth(e.Then, pos, depth+1)
		v.checkExprDepth(e.Else, pos, depth+1)
	case *ast.CallExpr:
		v.checkExprDepth(e.Object, pos, depth+1)
		for _, arg := range e.Args {
			v.checkExprDepth(arg, pos, depth+1)
		}
	case *ast.FieldExpr:
		v.checkExprDepth(e.Object, pos, depth+1)
	case *ast.NullishExpr:
		v.checkExprDepth(e.Left, pos, depth+1)
		if e.Right != nil {
			v.checkExprDepth(e.Right, pos, depth+1)
		}
	}
}

// ValidateDeadCode はデッドコードを検出する
func (v *Validator) ValidateDeadCode(spec *ast.SpecFile) error {
	v.errors = &errors.MultiError{}

	for _, comp := range spec.Components {
		for _, flow := range comp.Body.Flows {
			v.checkDeadCodeInSteps(flow.Steps)
		}
	}

	return v.errors.ErrorOrNil()
}

// checkDeadCodeInSteps はステップリスト内のデッドコードをチェックする
func (v *Validator) checkDeadCodeInSteps(steps []ast.Step) {
	for i, step := range steps {
		// 終端ステップの後にコードがあればデッドコード
		if v.isTerminalStep(step) && i < len(steps)-1 {
			nextStep := steps[i+1]
			v.errors.Add(&errors.ValidationError{
				Pos:     nextStep.GetPos(),
				Type:    "deadcode",
				Name:    "unreachable",
				Message: "code after return/throw is unreachable",
			})
		}

		// 制御構文内のデッドコードもチェック
		switch s := step.(type) {
		case *ast.IfStep:
			v.checkDeadCodeInSteps(s.Then)
			v.checkDeadCodeInSteps(s.Else)
		case *ast.ForStep:
			v.checkDeadCodeInSteps(s.Body)
		case *ast.WhileStep:
			v.checkDeadCodeInSteps(s.Body)
		}
	}
}

// isTerminalStep はステップが終端ステップ（return/throw）かを返す
func (v *Validator) isTerminalStep(step ast.Step) bool {
	switch step.(type) {
	case *ast.ReturnStep, *ast.ThrowStep:
		return true
	}
	return false
}

// validateEmptyFlows は空のフロー宣言を検出する
func (v *Validator) validateEmptyFlows(comp *ast.ComponentDecl) {
	for _, flow := range comp.Body.Flows {
		if len(flow.Steps) == 0 {
			v.errors.Add(&errors.ValidationError{
				Pos:     flow.Pos,
				Type:    "warning",
				Name:    flow.Name,
				Message: "empty flow declaration",
			})
		}
	}
}
