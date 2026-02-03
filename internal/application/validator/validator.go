package validator

import (
	"pact/internal/domain/ast"
	"pact/internal/domain/errors"
)

// Validator はASTの検証を行う
type Validator struct {
	errors   *errors.MultiError
	warnings *errors.WarningList
}

// validDurationUnits は有効な期間単位
var validDurationUnits = map[string]bool{
	"ms": true, "s": true, "m": true, "h": true, "d": true,
}

// NewValidator は新しいValidatorを作成する
func NewValidator() *Validator {
	return &Validator{
		errors:   &errors.MultiError{},
		warnings: &errors.WarningList{},
	}
}

// GetWarnings は収集した警告を返す
func (v *Validator) GetWarnings() *errors.WarningList {
	return v.warnings
}

// Validate はSpecFileを検証する
func (v *Validator) Validate(spec *ast.SpecFile) error {
	v.errors = &errors.MultiError{}

	// 各コンポーネントを検証
	for i := range spec.Components {
		v.validateComponent(&spec.Components[i])
	}

	return v.errors.ErrorOrNil()
}

// validateComponent はコンポーネントを検証する
func (v *Validator) validateComponent(comp *ast.ComponentDecl) {
	// 型名の重複チェック
	typeNames := make(map[string]ast.Position)
	for _, typ := range comp.Body.Types {
		if pos, exists := typeNames[typ.Name]; exists {
			v.errors.Add(&errors.ValidationError{
				Pos:     typ.Pos,
				Type:    "duplicate",
				Name:    typ.Name,
				Message: "type already defined at " + pos.String(),
			})
		} else {
			typeNames[typ.Name] = typ.Pos
		}

		// フィールド名の重複チェック
		v.validateDuplicateFields(&typ)
	}

	// インターフェース名の重複チェック
	ifaceNames := make(map[string]ast.Position)
	for _, iface := range comp.Body.Provides {
		if pos, exists := ifaceNames[iface.Name]; exists {
			v.errors.Add(&errors.ValidationError{
				Pos:     iface.Pos,
				Type:    "duplicate",
				Name:    iface.Name,
				Message: "interface already defined at " + pos.String(),
			})
		} else {
			ifaceNames[iface.Name] = iface.Pos
		}

		// メソッド名の重複チェック
		v.validateDuplicateMethods(&iface)
	}

	for _, iface := range comp.Body.Requires {
		if pos, exists := ifaceNames[iface.Name]; exists {
			v.errors.Add(&errors.ValidationError{
				Pos:     iface.Pos,
				Type:    "duplicate",
				Name:    iface.Name,
				Message: "interface already defined at " + pos.String(),
			})
		} else {
			ifaceNames[iface.Name] = iface.Pos
		}

		// メソッド名の重複チェック
		v.validateDuplicateMethods(&iface)
	}

	// フロー名の重複チェック
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

	// States名の重複チェック
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

	// Relation（depends_on）の重複チェック
	relationTargets := make(map[string]ast.Position)
	for _, rel := range comp.Body.Relations {
		key := string(rel.Kind) + ":" + rel.Target
		if pos, exists := relationTargets[key]; exists {
			v.errors.Add(&errors.ValidationError{
				Pos:     rel.Pos,
				Type:    "duplicate",
				Name:    rel.Target,
				Message: "relation already defined at " + pos.String(),
			})
		} else {
			relationTargets[key] = rel.Pos
		}
	}
}

// validateDuplicateFields はフィールド名の重複をチェックする
func (v *Validator) validateDuplicateFields(typ *ast.TypeDecl) {
	fieldNames := make(map[string]ast.Position)
	for _, field := range typ.Fields {
		if pos, exists := fieldNames[field.Name]; exists {
			v.errors.Add(&errors.ValidationError{
				Pos:     field.Pos,
				Type:    "duplicate",
				Name:    field.Name,
				Message: "field already defined at " + pos.String(),
			})
		} else {
			fieldNames[field.Name] = field.Pos
		}
	}
}

// validateDuplicateMethods はメソッド名の重複をチェックする
func (v *Validator) validateDuplicateMethods(iface *ast.InterfaceDecl) {
	methodNames := make(map[string]ast.Position)
	for _, method := range iface.Methods {
		if pos, exists := methodNames[method.Name]; exists {
			v.errors.Add(&errors.ValidationError{
				Pos:     method.Pos,
				Type:    "duplicate",
				Name:    method.Name,
				Message: "method already defined at " + pos.String(),
			})
		} else {
			methodNames[method.Name] = method.Pos
		}

		// パラメータ名の重複チェック
		v.validateDuplicateParams(&method)
	}
}

// validateDuplicateParams はパラメータ名の重複をチェックする
func (v *Validator) validateDuplicateParams(method *ast.MethodDecl) {
	paramNames := make(map[string]ast.Position)
	for _, param := range method.Params {
		if pos, exists := paramNames[param.Name]; exists {
			v.errors.Add(&errors.ValidationError{
				Pos:     param.Pos,
				Type:    "duplicate",
				Name:    param.Name,
				Message: "parameter already defined at " + pos.String(),
			})
		} else {
			paramNames[param.Name] = param.Pos
		}
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

// ValidateReferences は未定義参照を検証する
func (v *Validator) ValidateReferences(spec *ast.SpecFile) error {
	v.errors = &errors.MultiError{}

	// 定義された型、コンポーネント、状態を収集
	definedTypes := make(map[string]bool)
	definedComponents := make(map[string]bool)

	// 組み込み型
	builtinTypes := map[string]bool{
		"string": true, "String": true,
		"int": true, "Int": true, "integer": true, "Integer": true,
		"float": true, "Float": true, "double": true, "Double": true,
		"bool": true, "Bool": true, "boolean": true, "Boolean": true,
		"void": true, "Void": true,
		"any": true, "Any": true,
		"object": true, "Object": true,
		"array": true, "Array": true,
		"map": true, "Map": true,
		"date": true, "Date": true,
		"datetime": true, "DateTime": true,
		"time": true, "Time": true,
		"uuid": true, "UUID": true,
		"bytes": true, "Bytes": true,
	}

	// 定義された型を収集
	for _, comp := range spec.Components {
		definedComponents[comp.Name] = true
		for _, typ := range comp.Body.Types {
			definedTypes[typ.Name] = true
		}
	}

	// 参照を検証
	for _, comp := range spec.Components {
		// 関係のターゲットを検証
		for _, rel := range comp.Body.Relations {
			if !definedComponents[rel.Target] && !definedTypes[rel.Target] {
				v.errors.Add(&errors.ValidationError{
					Pos:     rel.Pos,
					Type:    "undefined",
					Name:    rel.Target,
					Message: "referenced component/type is not defined",
				})
			}
		}

		// フィールド型を検証
		for _, typ := range comp.Body.Types {
			for _, field := range typ.Fields {
				typeName := field.Type.Name
				if !definedTypes[typeName] && !builtinTypes[typeName] {
					v.errors.Add(&errors.ValidationError{
						Pos:     field.Pos,
						Type:    "undefined",
						Name:    typeName,
						Message: "type is not defined",
					})
				}
			}
		}

		// メソッドのパラメータと戻り値型を検証
		for _, iface := range comp.Body.Provides {
			v.validateMethodTypes(&iface, definedTypes, builtinTypes)
		}
		for _, iface := range comp.Body.Requires {
			v.validateMethodTypes(&iface, definedTypes, builtinTypes)
		}

		// States内の状態参照を検証
		for _, states := range comp.Body.States {
			v.validateStateReferences(&states)
		}
	}

	return v.errors.ErrorOrNil()
}

// validateMethodTypes はメソッドの型参照を検証する
func (v *Validator) validateMethodTypes(iface *ast.InterfaceDecl, definedTypes, builtinTypes map[string]bool) {
	for _, method := range iface.Methods {
		// パラメータ型
		for _, param := range method.Params {
			typeName := param.Type.Name
			if !definedTypes[typeName] && !builtinTypes[typeName] {
				v.errors.Add(&errors.ValidationError{
					Pos:     param.Pos,
					Type:    "undefined",
					Name:    typeName,
					Message: "parameter type is not defined",
				})
			}
		}
		// 戻り値型
		if method.ReturnType != nil {
			typeName := method.ReturnType.Name
			if !definedTypes[typeName] && !builtinTypes[typeName] {
				v.errors.Add(&errors.ValidationError{
					Pos:     method.Pos,
					Type:    "undefined",
					Name:    typeName,
					Message: "return type is not defined",
				})
			}
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

// MaxNestingDepth は式のネスト深度の上限
const MaxNestingDepth = 50

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

// ValidateEmptyDeclarations は空の宣言を検出する
func (v *Validator) ValidateEmptyDeclarations(spec *ast.SpecFile) error {
	v.errors = &errors.MultiError{}

	for _, comp := range spec.Components {
		// 空のenum
		for _, typ := range comp.Body.Types {
			if typ.Kind == "enum" && len(typ.Values) == 0 {
				v.errors.Add(&errors.ValidationError{
					Pos:     typ.Pos,
					Type:    "warning",
					Name:    typ.Name,
					Message: "empty enum declaration",
				})
			}
		}

		// 空のstates
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

		// 空のprovides/requires
		for _, iface := range comp.Body.Provides {
			if len(iface.Methods) == 0 {
				v.errors.Add(&errors.ValidationError{
					Pos:     iface.Pos,
					Type:    "warning",
					Name:    iface.Name,
					Message: "empty interface declaration in provides",
				})
			}
		}
		for _, iface := range comp.Body.Requires {
			if len(iface.Methods) == 0 {
				v.errors.Add(&errors.ValidationError{
					Pos:     iface.Pos,
					Type:    "warning",
					Name:    iface.Name,
					Message: "empty interface declaration in requires",
				})
			}
		}

		// 空のflow
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

	return v.errors.ErrorOrNil()
}

// ValidateAll は全ての検証を実行する
func (v *Validator) ValidateAll(spec *ast.SpecFile) error {
	multiErr := &errors.MultiError{}

	if err := v.Validate(spec); err != nil {
		if me, ok := err.(*errors.MultiError); ok {
			for _, e := range me.Errors {
				multiErr.Add(e)
			}
		} else {
			multiErr.Add(err)
		}
	}

	if err := v.ValidateReferences(spec); err != nil {
		if me, ok := err.(*errors.MultiError); ok {
			for _, e := range me.Errors {
				multiErr.Add(e)
			}
		} else {
			multiErr.Add(err)
		}
	}

	if err := v.ValidateExpressionDepth(spec); err != nil {
		if me, ok := err.(*errors.MultiError); ok {
			for _, e := range me.Errors {
				multiErr.Add(e)
			}
		} else {
			multiErr.Add(err)
		}
	}

	if err := v.ValidateDeadCode(spec); err != nil {
		if me, ok := err.(*errors.MultiError); ok {
			for _, e := range me.Errors {
				multiErr.Add(e)
			}
		} else {
			multiErr.Add(err)
		}
	}

	if err := v.ValidateEmptyDeclarations(spec); err != nil {
		if me, ok := err.(*errors.MultiError); ok {
			for _, e := range me.Errors {
				multiErr.Add(e)
			}
		} else {
			multiErr.Add(err)
		}
	}

	// 追加のバリデーション
	if err := v.ValidateDurationUnits(spec); err != nil {
		if me, ok := err.(*errors.MultiError); ok {
			for _, e := range me.Errors {
				multiErr.Add(e)
			}
		} else {
			multiErr.Add(err)
		}
	}

	// 警告を収集
	v.CollectWarnings(spec)

	return multiErr.ErrorOrNil()
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

// CollectWarnings は警告を収集する（L-005, L-007）
func (v *Validator) CollectWarnings(spec *ast.SpecFile) {
	v.warnings = &errors.WarningList{}

	// 未使用のimportを検出
	v.checkUnusedImports(spec)

	// 未使用の型を検出
	v.checkUnusedTypes(spec)

	// @deprecated アノテーションの使用チェック
	v.checkDeprecatedUsage(spec)
}

// checkUnusedImports は未使用のimportを検出する
func (v *Validator) checkUnusedImports(spec *ast.SpecFile) {
	if len(spec.Imports) == 0 {
		return
	}

	// importされたパスを収集
	for _, imp := range spec.Imports {
		// エイリアスがある場合はエイリアス名で参照されるべき
		// 簡易チェック: importが存在するが、対応するコンポーネントがない場合
		alias := imp.Path
		if imp.Alias != nil {
			alias = *imp.Alias
		}
		used := false

		// コンポーネント内の関係で使用されているか
		for _, comp := range spec.Components {
			for _, rel := range comp.Body.Relations {
				if rel.Target == alias {
					used = true
					break
				}
			}
			if used {
				break
			}
		}

		if !used {
			v.warnings.Add(&errors.Warning{
				Pos:     imp.Pos,
				Code:    "unused-import",
				Message: "import '" + imp.Path + "' is not referenced",
			})
		}
	}
}

// checkUnusedTypes は未使用の型を検出する
func (v *Validator) checkUnusedTypes(spec *ast.SpecFile) {
	// 定義された型と使用された型を追跡
	for _, comp := range spec.Components {
		definedTypes := make(map[string]ast.Position)
		usedTypes := make(map[string]bool)

		// 定義を収集
		for _, typ := range comp.Body.Types {
			definedTypes[typ.Name] = typ.Pos
		}

		// フィールド型の使用を収集
		for _, typ := range comp.Body.Types {
			for _, field := range typ.Fields {
				usedTypes[field.Type.Name] = true
				for _, tp := range field.Type.TypeParams {
					usedTypes[tp.Name] = true
				}
			}
			// エイリアスの基底型
			if typ.BaseType != nil {
				usedTypes[typ.BaseType.Name] = true
			}
		}

		// メソッドのパラメータ・戻り値型の使用を収集
		for _, iface := range comp.Body.Provides {
			for _, method := range iface.Methods {
				for _, param := range method.Params {
					usedTypes[param.Type.Name] = true
				}
				if method.ReturnType != nil {
					usedTypes[method.ReturnType.Name] = true
				}
			}
		}
		for _, iface := range comp.Body.Requires {
			for _, method := range iface.Methods {
				for _, param := range method.Params {
					usedTypes[param.Type.Name] = true
				}
				if method.ReturnType != nil {
					usedTypes[method.ReturnType.Name] = true
				}
			}
		}

		// 関係ターゲットの使用を収集
		for _, rel := range comp.Body.Relations {
			usedTypes[rel.Target] = true
		}

		// 未使用の型を警告
		for name, pos := range definedTypes {
			if !usedTypes[name] {
				v.warnings.Add(&errors.Warning{
					Pos:     pos,
					Code:    "unused-type",
					Message: "type '" + name + "' is defined but not referenced",
				})
			}
		}
	}
}

// checkDeprecatedUsage は @deprecated アノテーションの付いた要素の使用を検出する（L-005）
func (v *Validator) checkDeprecatedUsage(spec *ast.SpecFile) {
	// 非推奨の型・メソッドを収集
	deprecatedTypes := make(map[string]bool)
	deprecatedMethods := make(map[string]bool) // "InterfaceName.MethodName"

	for _, comp := range spec.Components {
		for _, typ := range comp.Body.Types {
			if hasAnnotation(typ.Annotations, "deprecated") {
				deprecatedTypes[typ.Name] = true
			}
		}

		for _, iface := range comp.Body.Provides {
			for _, method := range iface.Methods {
				if hasAnnotation(method.Annotations, "deprecated") {
					deprecatedMethods[iface.Name+"."+method.Name] = true
				}
			}
		}
	}

	// 非推奨の型が参照されている場合は警告
	for _, comp := range spec.Components {
		for _, typ := range comp.Body.Types {
			for _, field := range typ.Fields {
				if deprecatedTypes[field.Type.Name] {
					v.warnings.Add(&errors.Warning{
						Pos:     field.Pos,
						Code:    "deprecated",
						Message: "type '" + field.Type.Name + "' is deprecated",
					})
				}
			}
		}
	}
}

// hasAnnotation はアノテーションリストに指定のアノテーションがあるかを返す
func hasAnnotation(annotations []ast.AnnotationDecl, name string) bool {
	for _, ann := range annotations {
		if ann.Name == name {
			return true
		}
	}
	return false
}
