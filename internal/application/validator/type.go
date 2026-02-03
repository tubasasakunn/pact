package validator

import (
	"pact/internal/domain/ast"
	"pact/internal/domain/errors"
)

// builtinTypes は組み込み型
var builtinTypes = map[string]bool{
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

// validateTypeDuplicates は型名の重複をチェックする
func (v *Validator) validateTypeDuplicates(comp *ast.ComponentDecl) {
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
}

// validateInterfaceDuplicates はインターフェース名の重複をチェックする
func (v *Validator) validateInterfaceDuplicates(comp *ast.ComponentDecl) {
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

// validateFieldTypeReferences はフィールド型の参照を検証する
func (v *Validator) validateFieldTypeReferences(comp *ast.ComponentDecl, definedTypes map[string]bool) {
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
}

// validateMethodTypeReferences はメソッドの型参照を検証する
func (v *Validator) validateMethodTypeReferences(comp *ast.ComponentDecl, definedTypes map[string]bool) {
	for _, iface := range comp.Body.Provides {
		v.validateMethodTypes(&iface, definedTypes)
	}
	for _, iface := range comp.Body.Requires {
		v.validateMethodTypes(&iface, definedTypes)
	}
}

// validateMethodTypes はメソッドの型参照を検証する
func (v *Validator) validateMethodTypes(iface *ast.InterfaceDecl, definedTypes map[string]bool) {
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

// validateEmptyTypes は空の型宣言を検出する
func (v *Validator) validateEmptyTypes(comp *ast.ComponentDecl) {
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
}

// validateEmptyInterfaces は空のインターフェース宣言を検出する
func (v *Validator) validateEmptyInterfaces(comp *ast.ComponentDecl) {
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
