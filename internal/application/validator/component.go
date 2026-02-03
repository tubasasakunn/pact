package validator

import (
	"pact/internal/domain/ast"
	"pact/internal/domain/errors"
)

// validateRelationDuplicates はRelation（depends_on）の重複をチェックする
func (v *Validator) validateRelationDuplicates(comp *ast.ComponentDecl) {
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

// validateRelationReferences は関係のターゲットを検証する
func (v *Validator) validateRelationReferences(comp *ast.ComponentDecl, definedComponents, definedTypes map[string]bool) {
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
