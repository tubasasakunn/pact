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
	v.validateTypeDuplicates(comp)

	// インターフェース名の重複チェック
	v.validateInterfaceDuplicates(comp)

	// フロー名の重複チェック
	v.validateFlowDuplicates(comp)

	// States名の重複チェック
	v.validateStatesDuplicates(comp)

	// Relation（depends_on）の重複チェック
	v.validateRelationDuplicates(comp)
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

// ValidateReferences は未定義参照を検証する
func (v *Validator) ValidateReferences(spec *ast.SpecFile) error {
	v.errors = &errors.MultiError{}

	// 定義された型、コンポーネントを収集
	definedTypes := make(map[string]bool)
	definedComponents := make(map[string]bool)

	for _, comp := range spec.Components {
		definedComponents[comp.Name] = true
		for _, typ := range comp.Body.Types {
			definedTypes[typ.Name] = true
		}
	}

	// 参照を検証
	for _, comp := range spec.Components {
		// 関係のターゲットを検証
		v.validateRelationReferences(&comp, definedComponents, definedTypes)

		// フィールド型を検証
		v.validateFieldTypeReferences(&comp, definedTypes)

		// メソッドのパラメータと戻り値型を検証
		v.validateMethodTypeReferences(&comp, definedTypes)

		// States内の状態参照を検証
		v.validateStateReferencesInComp(&comp)
	}

	return v.errors.ErrorOrNil()
}

// ValidateEmptyDeclarations は空の宣言を検出する
func (v *Validator) ValidateEmptyDeclarations(spec *ast.SpecFile) error {
	v.errors = &errors.MultiError{}

	for _, comp := range spec.Components {
		// 空のenum
		v.validateEmptyTypes(&comp)

		// 空のstates
		v.validateEmptyStates(&comp)

		// 空のprovides/requires
		v.validateEmptyInterfaces(&comp)

		// 空のflow
		v.validateEmptyFlows(&comp)
	}

	return v.errors.ErrorOrNil()
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

// hasAnnotation はアノテーションリストに指定のアノテーションがあるかを返す
func hasAnnotation(annotations []ast.AnnotationDecl, name string) bool {
	for _, ann := range annotations {
		if ann.Name == name {
			return true
		}
	}
	return false
}
