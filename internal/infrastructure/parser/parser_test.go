package parser

import (
	"testing"

	"pact/internal/domain/ast"
	"pact/internal/domain/errors"
)

// isParseError はエラーがParseErrorまたはParseErrorを含むMultiErrorかどうかを確認する
func isParseError(err error) bool {
	if _, ok := err.(*errors.ParseError); ok {
		return true
	}
	if me, ok := err.(*errors.MultiError); ok {
		for _, e := range me.Errors {
			if _, ok := e.(*errors.ParseError); ok {
				return true
			}
		}
	}
	return false
}

// getFirstParseError はエラーからParseErrorを取得する
func getFirstParseError(err error) *errors.ParseError {
	if pe, ok := err.(*errors.ParseError); ok {
		return pe
	}
	if me, ok := err.(*errors.MultiError); ok {
		for _, e := range me.Errors {
			if pe, ok := e.(*errors.ParseError); ok {
				return pe
			}
		}
	}
	return nil
}

// =============================================================================
// 1.2.1 インポート文
// =============================================================================

// P001: 単純インポート
func TestParser_Import_Simple(t *testing.T) {
	input := `import "./foo.pact"`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Imports) != 1 {
		t.Fatalf("expected 1 import, got %d", len(spec.Imports))
	}
	if spec.Imports[0].Path != "./foo.pact" {
		t.Errorf("expected path './foo.pact', got %q", spec.Imports[0].Path)
	}
	if spec.Imports[0].Alias != nil {
		t.Errorf("expected no alias, got %q", *spec.Imports[0].Alias)
	}
}

// P002: エイリアス付きインポート
func TestParser_Import_WithAlias(t *testing.T) {
	input := `import "./foo.pact" as Foo`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Imports) != 1 {
		t.Fatalf("expected 1 import, got %d", len(spec.Imports))
	}
	if spec.Imports[0].Alias == nil {
		t.Fatalf("expected alias, got nil")
	}
	if *spec.Imports[0].Alias != "Foo" {
		t.Errorf("expected alias 'Foo', got %q", *spec.Imports[0].Alias)
	}
}

// P003: 複数インポート
func TestParser_Import_Multiple(t *testing.T) {
	input := `import "./a.pact"
import "./b.pact" as B`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Imports) != 2 {
		t.Fatalf("expected 2 imports, got %d", len(spec.Imports))
	}
}

// P004: パスなしエラー
func TestParser_Import_MissingPath(t *testing.T) {
	input := `import`
	_, err := ParseString(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !isParseError(err) {
		t.Errorf("expected ParseError or MultiError containing ParseError, got %T", err)
	}
}

// P005: 不正なパス
func TestParser_Import_InvalidPath(t *testing.T) {
	input := `import 123`
	_, err := ParseString(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// P006: 位置情報
func TestParser_Import_Position(t *testing.T) {
	input := `import "./foo.pact"`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.Imports[0].Pos.Line != 1 {
		t.Errorf("expected line 1, got %d", spec.Imports[0].Pos.Line)
	}
	if spec.Imports[0].Pos.Column != 1 {
		t.Errorf("expected column 1, got %d", spec.Imports[0].Pos.Column)
	}
}

// =============================================================================
// 1.2.2 コンポーネント宣言
// =============================================================================

// P010: 空コンポーネント
func TestParser_Component_Empty(t *testing.T) {
	input := `component Foo { }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.Component == nil {
		t.Fatal("expected component, got nil")
	}
	if spec.Component.Name != "Foo" {
		t.Errorf("expected name 'Foo', got %q", spec.Component.Name)
	}
}

// P011: アノテーション付きコンポーネント
func TestParser_Component_WithAnnotation(t *testing.T) {
	input := `@description("A test component")
component Foo { }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Component.Annotations) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(spec.Component.Annotations))
	}
	if spec.Component.Annotations[0].Name != "description" {
		t.Errorf("expected annotation name 'description', got %q", spec.Component.Annotations[0].Name)
	}
}

// P012: 名前なしエラー
func TestParser_Component_MissingName(t *testing.T) {
	input := `component { }`
	_, err := ParseString(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// P013: 波括弧なしエラー
func TestParser_Component_MissingBrace(t *testing.T) {
	input := `component Foo`
	_, err := ParseString(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// P014: 位置情報
func TestParser_Component_Position(t *testing.T) {
	input := `component Foo { }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.Component.Pos.Line != 1 {
		t.Errorf("expected line 1, got %d", spec.Component.Pos.Line)
	}
	if spec.Component.Pos.Column != 1 {
		t.Errorf("expected column 1, got %d", spec.Component.Pos.Column)
	}
}

// =============================================================================
// 1.2.3 型定義
// =============================================================================

// P020: 空の型
func TestParser_Type_Empty(t *testing.T) {
	input := `component Foo { type Bar { } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Component.Body.Types) != 1 {
		t.Fatalf("expected 1 type, got %d", len(spec.Component.Body.Types))
	}
	typ := spec.Component.Body.Types[0]
	if typ.Name != "Bar" {
		t.Errorf("expected name 'Bar', got %q", typ.Name)
	}
	if typ.Kind != ast.TypeKindStruct {
		t.Errorf("expected kind 'struct', got %q", typ.Kind)
	}
	if len(typ.Fields) != 0 {
		t.Errorf("expected 0 fields, got %d", len(typ.Fields))
	}
}

// P021: 単一フィールド
func TestParser_Type_SingleField(t *testing.T) {
	input := `component Foo { type Bar { name: string } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Component.Body.Types[0].Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(spec.Component.Body.Types[0].Fields))
	}
	field := spec.Component.Body.Types[0].Fields[0]
	if field.Name != "name" {
		t.Errorf("expected field name 'name', got %q", field.Name)
	}
	if field.Type.Name != "string" {
		t.Errorf("expected type 'string', got %q", field.Type.Name)
	}
}

// P022: 複数フィールド
func TestParser_Type_MultipleFields(t *testing.T) {
	input := `component Foo { type Bar { a: int b: string } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Component.Body.Types[0].Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(spec.Component.Body.Types[0].Fields))
	}
}

// P023: nullableフィールド
func TestParser_Type_NullableField(t *testing.T) {
	input := `component Foo { type Bar { name: string? } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	field := spec.Component.Body.Types[0].Fields[0]
	if !field.Type.Nullable {
		t.Errorf("expected nullable type")
	}
}

// P024: 配列フィールド
func TestParser_Type_ArrayField(t *testing.T) {
	input := `component Foo { type Bar { items: string[] } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	field := spec.Component.Body.Types[0].Fields[0]
	if !field.Type.Array {
		t.Errorf("expected array type")
	}
}

// P025: nullable配列フィールド
func TestParser_Type_NullableArrayField(t *testing.T) {
	input := `component Foo { type Bar { items: string?[] } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	field := spec.Component.Body.Types[0].Fields[0]
	if !field.Type.Nullable || !field.Type.Array {
		t.Errorf("expected nullable array type")
	}
}

// P026: public可視性
func TestParser_Type_VisibilityPublic(t *testing.T) {
	input := `component Foo { type Bar { +name: string } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	field := spec.Component.Body.Types[0].Fields[0]
	if field.Visibility != ast.VisibilityPublic {
		t.Errorf("expected public visibility, got %q", field.Visibility)
	}
}

// P027: private可視性
func TestParser_Type_VisibilityPrivate(t *testing.T) {
	input := `component Foo { type Bar { -name: string } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	field := spec.Component.Body.Types[0].Fields[0]
	if field.Visibility != ast.VisibilityPrivate {
		t.Errorf("expected private visibility, got %q", field.Visibility)
	}
}

// P028: protected可視性
func TestParser_Type_VisibilityProtected(t *testing.T) {
	input := `component Foo { type Bar { #name: string } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	field := spec.Component.Body.Types[0].Fields[0]
	if field.Visibility != ast.VisibilityProtected {
		t.Errorf("expected protected visibility, got %q", field.Visibility)
	}
}

// P029: package可視性
func TestParser_Type_VisibilityPackage(t *testing.T) {
	input := `component Foo { type Bar { ~name: string } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	field := spec.Component.Body.Types[0].Fields[0]
	if field.Visibility != ast.VisibilityPackage {
		t.Errorf("expected package visibility, got %q", field.Visibility)
	}
}

// P030: フィールドアノテーション
func TestParser_Type_FieldAnnotation(t *testing.T) {
	input := `component Foo { type Bar { @desc("x") name: string } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	field := spec.Component.Body.Types[0].Fields[0]
	if len(field.Annotations) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(field.Annotations))
	}
}

// P031: フィールド位置
func TestParser_Type_FieldPosition(t *testing.T) {
	input := `component Foo { type Bar { name: string } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	field := spec.Component.Body.Types[0].Fields[0]
	if !field.Pos.IsValid() {
		t.Errorf("expected valid position")
	}
}

// P035: 単純enum
func TestParser_Enum_Simple(t *testing.T) {
	input := `component Foo { enum Status { A B C } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	typ := spec.Component.Body.Types[0]
	if typ.Kind != ast.TypeKindEnum {
		t.Errorf("expected kind 'enum', got %q", typ.Kind)
	}
	if len(typ.Values) != 3 {
		t.Errorf("expected 3 values, got %d", len(typ.Values))
	}
}

// P036: 空enum
func TestParser_Enum_Empty(t *testing.T) {
	input := `component Foo { enum Status { } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	typ := spec.Component.Body.Types[0]
	if len(typ.Values) != 0 {
		t.Errorf("expected 0 values, got %d", len(typ.Values))
	}
}

// P037: enumアノテーション
func TestParser_Enum_WithAnnotation(t *testing.T) {
	input := `component Foo { @desc("x") enum Status { A } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	typ := spec.Component.Body.Types[0]
	if len(typ.Annotations) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(typ.Annotations))
	}
}

// =============================================================================
// 1.2.4 関係定義
// =============================================================================

// P040: depends on
func TestParser_Relation_DependsOn(t *testing.T) {
	input := `component Foo { depends on Bar }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Component.Body.Relations) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(spec.Component.Body.Relations))
	}
	rel := spec.Component.Body.Relations[0]
	if rel.Kind != ast.RelationDependsOn {
		t.Errorf("expected kind 'depends_on', got %q", rel.Kind)
	}
	if rel.Target != "Bar" {
		t.Errorf("expected target 'Bar', got %q", rel.Target)
	}
}

// P041: 型指定
func TestParser_Relation_DependsOn_WithType(t *testing.T) {
	input := `component Foo { depends on Bar: database }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rel := spec.Component.Body.Relations[0]
	if rel.TargetType == nil {
		t.Fatalf("expected target type, got nil")
	}
	if *rel.TargetType != "database" {
		t.Errorf("expected target type 'database', got %q", *rel.TargetType)
	}
}

// P042: エイリアス
func TestParser_Relation_DependsOn_WithAlias(t *testing.T) {
	input := `component Foo { depends on Bar as B }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rel := spec.Component.Body.Relations[0]
	if rel.Alias == nil {
		t.Fatalf("expected alias, got nil")
	}
	if *rel.Alias != "B" {
		t.Errorf("expected alias 'B', got %q", *rel.Alias)
	}
}

// P043: フル指定
func TestParser_Relation_DependsOn_Full(t *testing.T) {
	input := `component Foo { depends on Bar: external as B }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rel := spec.Component.Body.Relations[0]
	if rel.TargetType == nil || *rel.TargetType != "external" {
		t.Errorf("expected target type 'external'")
	}
	if rel.Alias == nil || *rel.Alias != "B" {
		t.Errorf("expected alias 'B'")
	}
}

// P044: extends
func TestParser_Relation_Extends(t *testing.T) {
	input := `component Foo { extends Base }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rel := spec.Component.Body.Relations[0]
	if rel.Kind != ast.RelationExtends {
		t.Errorf("expected kind 'extends', got %q", rel.Kind)
	}
}

// P045: implements
func TestParser_Relation_Implements(t *testing.T) {
	input := `component Foo { implements Iface }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rel := spec.Component.Body.Relations[0]
	if rel.Kind != ast.RelationImplements {
		t.Errorf("expected kind 'implements', got %q", rel.Kind)
	}
}

// P046: contains
func TestParser_Relation_Contains(t *testing.T) {
	input := `component Foo { contains Cache }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rel := spec.Component.Body.Relations[0]
	if rel.Kind != ast.RelationContains {
		t.Errorf("expected kind 'contains', got %q", rel.Kind)
	}
}

// P047: aggregates
func TestParser_Relation_Aggregates(t *testing.T) {
	input := `component Foo { aggregates Items }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rel := spec.Component.Body.Relations[0]
	if rel.Kind != ast.RelationAggregates {
		t.Errorf("expected kind 'aggregates', got %q", rel.Kind)
	}
}

// P048: アノテーション
func TestParser_Relation_WithAnnotation(t *testing.T) {
	input := `component Foo { @desc("x") depends on Bar }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rel := spec.Component.Body.Relations[0]
	if len(rel.Annotations) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(rel.Annotations))
	}
}

// P049: 全TargetType
func TestParser_Relation_AllTargetTypes(t *testing.T) {
	targetTypes := []string{"database", "external", "queue", "actor", "service"}
	for _, tt := range targetTypes {
		input := `component Foo { depends on Bar: ` + tt + ` }`
		spec, err := ParseString(input)
		if err != nil {
			t.Errorf("unexpected error for %s: %v", tt, err)
			continue
		}
		rel := spec.Component.Body.Relations[0]
		if rel.TargetType == nil || *rel.TargetType != tt {
			t.Errorf("expected target type %q", tt)
		}
	}
}

// P04A: 位置情報
func TestParser_Relation_Position(t *testing.T) {
	input := `component Foo { depends on Bar }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rel := spec.Component.Body.Relations[0]
	if !rel.Pos.IsValid() {
		t.Errorf("expected valid position")
	}
}

// =============================================================================
// 1.2.5 インターフェース定義
// =============================================================================

// P050: 空provides
func TestParser_Interface_Provides_Empty(t *testing.T) {
	input := `component Foo { provides API { } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Component.Body.Provides) != 1 {
		t.Fatalf("expected 1 provides, got %d", len(spec.Component.Body.Provides))
	}
	iface := spec.Component.Body.Provides[0]
	if iface.Name != "API" {
		t.Errorf("expected name 'API', got %q", iface.Name)
	}
	if len(iface.Methods) != 0 {
		t.Errorf("expected 0 methods, got %d", len(iface.Methods))
	}
}

// P051: 空requires
func TestParser_Interface_Requires_Empty(t *testing.T) {
	input := `component Foo { requires Query { } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Component.Body.Requires) != 1 {
		t.Fatalf("expected 1 requires, got %d", len(spec.Component.Body.Requires))
	}
}

// P052: 単一メソッド
func TestParser_Interface_SingleMethod(t *testing.T) {
	input := `component Foo { provides API { Get() -> Item } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	iface := spec.Component.Body.Provides[0]
	if len(iface.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(iface.Methods))
	}
	method := iface.Methods[0]
	if method.Name != "Get" {
		t.Errorf("expected method name 'Get', got %q", method.Name)
	}
	if method.ReturnType == nil || method.ReturnType.Name != "Item" {
		t.Errorf("expected return type 'Item'")
	}
}

// P053: パラメータ付きメソッド
func TestParser_Interface_MethodWithParams(t *testing.T) {
	input := `component Foo { provides API { Get(id: string) -> Item } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	method := spec.Component.Body.Provides[0].Methods[0]
	if len(method.Params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(method.Params))
	}
	param := method.Params[0]
	if param.Name != "id" {
		t.Errorf("expected param name 'id', got %q", param.Name)
	}
	if param.Type.Name != "string" {
		t.Errorf("expected param type 'string', got %q", param.Type.Name)
	}
}

// P054: 複数パラメータ
func TestParser_Interface_MethodMultiParams(t *testing.T) {
	input := `component Foo { provides API { Get(a: int, b: string) -> Item } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	method := spec.Component.Body.Provides[0].Methods[0]
	if len(method.Params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(method.Params))
	}
}

// P055: throws付きメソッド
func TestParser_Interface_MethodThrows(t *testing.T) {
	input := `component Foo { provides API { Get() -> Item throws NotFound } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	method := spec.Component.Body.Provides[0].Methods[0]
	if len(method.Throws) != 1 {
		t.Fatalf("expected 1 throws, got %d", len(method.Throws))
	}
	if method.Throws[0] != "NotFound" {
		t.Errorf("expected throws 'NotFound', got %q", method.Throws[0])
	}
}

// P056: 複数throws
func TestParser_Interface_MethodMultiThrows(t *testing.T) {
	input := `component Foo { provides API { Get() -> Item throws A, B } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	method := spec.Component.Body.Provides[0].Methods[0]
	if len(method.Throws) != 2 {
		t.Fatalf("expected 2 throws, got %d", len(method.Throws))
	}
}

// P057: asyncメソッド
func TestParser_Interface_MethodAsync(t *testing.T) {
	input := `component Foo { provides API { async Send() -> void } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	method := spec.Component.Body.Provides[0].Methods[0]
	if !method.Async {
		t.Errorf("expected async method")
	}
}

// P058: void戻り値
func TestParser_Interface_MethodVoid(t *testing.T) {
	input := `component Foo { provides API { Delete() -> void } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	method := spec.Component.Body.Provides[0].Methods[0]
	if method.ReturnType == nil || method.ReturnType.Name != "void" {
		t.Errorf("expected return type 'void'")
	}
}

// P059: メソッドアノテーション
func TestParser_Interface_MethodAnnotation(t *testing.T) {
	input := `component Foo { provides API { @desc("x") Get() -> Item } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	method := spec.Component.Body.Provides[0].Methods[0]
	if len(method.Annotations) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(method.Annotations))
	}
}

// P060: 複数メソッド
func TestParser_Interface_MultipleMethod(t *testing.T) {
	input := `component Foo { provides API { Get() -> Item Create(item: Item) -> Item } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	iface := spec.Component.Body.Provides[0]
	if len(iface.Methods) != 2 {
		t.Fatalf("expected 2 methods, got %d", len(iface.Methods))
	}
}

// P061: 位置情報
func TestParser_Interface_MethodPosition(t *testing.T) {
	input := `component Foo { provides API { Get() -> Item } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	method := spec.Component.Body.Provides[0].Methods[0]
	if !method.Pos.IsValid() {
		t.Errorf("expected valid position")
	}
}

// =============================================================================
// 1.2.6 フロー定義
// =============================================================================

// P070: 空フロー
func TestParser_Flow_Empty(t *testing.T) {
	input := `component Foo { flow Process { } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Component.Body.Flows) != 1 {
		t.Fatalf("expected 1 flow, got %d", len(spec.Component.Body.Flows))
	}
	flow := spec.Component.Body.Flows[0]
	if flow.Name != "Process" {
		t.Errorf("expected name 'Process', got %q", flow.Name)
	}
}

// P071: 単純代入
func TestParser_Flow_Assign_Simple(t *testing.T) {
	input := `component Foo { flow F { x = y } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Component.Body.Flows[0].Steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(spec.Component.Body.Flows[0].Steps))
	}
	step, ok := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	if !ok {
		t.Fatalf("expected AssignStep, got %T", spec.Component.Body.Flows[0].Steps[0])
	}
	if step.Variable != "x" {
		t.Errorf("expected variable 'x', got %q", step.Variable)
	}
}

// P072: メソッド呼び出し代入
func TestParser_Flow_Assign_Call(t *testing.T) {
	input := `component Foo { flow F { x = A.B() } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	_, ok := step.Value.(*ast.CallExpr)
	if !ok {
		t.Errorf("expected CallExpr, got %T", step.Value)
	}
}

// P073: 引数付き呼び出し
func TestParser_Flow_Assign_CallWithArgs(t *testing.T) {
	input := `component Foo { flow F { x = A.B(c, d) } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	call, ok := step.Value.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr, got %T", step.Value)
	}
	if len(call.Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(call.Args))
	}
}

// P074: null合体throw
func TestParser_Flow_Assign_ThrowOnNull(t *testing.T) {
	input := `component Foo { flow F { x = A.B() ?? throw E } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	nullish, ok := step.Value.(*ast.NullishExpr)
	if !ok {
		t.Fatalf("expected NullishExpr, got %T", step.Value)
	}
	if nullish.ThrowErr == nil || *nullish.ThrowErr != "E" {
		t.Errorf("expected ThrowErr 'E'")
	}
}

// P075: 単純呼び出し
func TestParser_Flow_Call(t *testing.T) {
	input := `component Foo { flow F { A.B() } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := spec.Component.Body.Flows[0].Steps[0].(*ast.CallStep)
	if !ok {
		t.Fatalf("expected CallStep, got %T", spec.Component.Body.Flows[0].Steps[0])
	}
}

// P076: await
func TestParser_Flow_Call_Await(t *testing.T) {
	input := `component Foo { flow F { await A.B() } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.CallStep)
	if !step.Await {
		t.Errorf("expected await")
	}
}

// P077: return
func TestParser_Flow_Return(t *testing.T) {
	input := `component Foo { flow F { return x } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step, ok := spec.Component.Body.Flows[0].Steps[0].(*ast.ReturnStep)
	if !ok {
		t.Fatalf("expected ReturnStep, got %T", spec.Component.Body.Flows[0].Steps[0])
	}
	if step.Value == nil {
		t.Errorf("expected return value")
	}
}

// P078: 空return
func TestParser_Flow_Return_Empty(t *testing.T) {
	input := `component Foo { flow F { return } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.ReturnStep)
	if step.Value != nil {
		t.Errorf("expected no return value")
	}
}

// P079: throw
func TestParser_Flow_Throw(t *testing.T) {
	input := `component Foo { flow F { throw E } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step, ok := spec.Component.Body.Flows[0].Steps[0].(*ast.ThrowStep)
	if !ok {
		t.Fatalf("expected ThrowStep, got %T", spec.Component.Body.Flows[0].Steps[0])
	}
	if step.Error != "E" {
		t.Errorf("expected error 'E', got %q", step.Error)
	}
}

// P080: 単純if
func TestParser_Flow_If_Simple(t *testing.T) {
	input := `component Foo { flow F { if x { } } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step, ok := spec.Component.Body.Flows[0].Steps[0].(*ast.IfStep)
	if !ok {
		t.Fatalf("expected IfStep, got %T", spec.Component.Body.Flows[0].Steps[0])
	}
	if step.Else != nil && len(step.Else) > 0 {
		t.Errorf("expected no else block")
	}
}

// P081: if-else
func TestParser_Flow_If_Else(t *testing.T) {
	input := `component Foo { flow F { if x { } else { } } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.IfStep)
	if step.Else == nil {
		t.Errorf("expected else block")
	}
}

// P082: ネストif
func TestParser_Flow_If_Nested(t *testing.T) {
	input := `component Foo { flow F { if x { if y { } } } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.IfStep)
	if len(step.Then) != 1 {
		t.Fatalf("expected 1 then step, got %d", len(step.Then))
	}
	_, ok := step.Then[0].(*ast.IfStep)
	if !ok {
		t.Errorf("expected nested IfStep")
	}
}

// P083: for
func TestParser_Flow_For(t *testing.T) {
	input := `component Foo { flow F { for x in items { } } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step, ok := spec.Component.Body.Flows[0].Steps[0].(*ast.ForStep)
	if !ok {
		t.Fatalf("expected ForStep, got %T", spec.Component.Body.Flows[0].Steps[0])
	}
	if step.Variable != "x" {
		t.Errorf("expected variable 'x', got %q", step.Variable)
	}
}

// P084: ネストfor
func TestParser_Flow_For_Nested(t *testing.T) {
	input := `component Foo { flow F { for x in a { for y in b { } } } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.ForStep)
	if len(step.Body) != 1 {
		t.Fatalf("expected 1 body step, got %d", len(step.Body))
	}
	_, ok := step.Body[0].(*ast.ForStep)
	if !ok {
		t.Errorf("expected nested ForStep")
	}
}

// P085: while
func TestParser_Flow_While(t *testing.T) {
	input := `component Foo { flow F { while cond { } } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := spec.Component.Body.Flows[0].Steps[0].(*ast.WhileStep)
	if !ok {
		t.Fatalf("expected WhileStep, got %T", spec.Component.Body.Flows[0].Steps[0])
	}
}

// P086: 複合フロー
func TestParser_Flow_Complex(t *testing.T) {
	input := `component Foo { flow F {
		x = A.Get()
		if x {
			B.Process(x)
		} else {
			throw Error
		}
		return x
	} }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Component.Body.Flows[0].Steps) != 3 {
		t.Errorf("expected 3 steps, got %d", len(spec.Component.Body.Flows[0].Steps))
	}
}

// P087: ステップアノテーション
func TestParser_Flow_StepAnnotation(t *testing.T) {
	input := `component Foo { flow F { @desc("x") a = b } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	if len(step.Annotations) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(step.Annotations))
	}
}

// P088: 位置情報
func TestParser_Flow_StepPosition(t *testing.T) {
	input := `component Foo { flow F { x = y } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0]
	if !step.GetPos().IsValid() {
		t.Errorf("expected valid position")
	}
}

// =============================================================================
// 1.2.7 式
// =============================================================================

// P100: 文字列リテラル
func TestParser_Expr_Literal_String(t *testing.T) {
	input := `component Foo { flow F { x = "hello" } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	lit, ok := step.Value.(*ast.LiteralExpr)
	if !ok {
		t.Fatalf("expected LiteralExpr, got %T", step.Value)
	}
	if lit.Value != "hello" {
		t.Errorf("expected 'hello', got %v", lit.Value)
	}
}

// P101: 整数リテラル
func TestParser_Expr_Literal_Int(t *testing.T) {
	input := `component Foo { flow F { x = 42 } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	lit, ok := step.Value.(*ast.LiteralExpr)
	if !ok {
		t.Fatalf("expected LiteralExpr, got %T", step.Value)
	}
	if lit.Value != int64(42) {
		t.Errorf("expected 42, got %v", lit.Value)
	}
}

// P102: 浮動小数点リテラル
func TestParser_Expr_Literal_Float(t *testing.T) {
	input := `component Foo { flow F { x = 3.14 } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	lit, ok := step.Value.(*ast.LiteralExpr)
	if !ok {
		t.Fatalf("expected LiteralExpr, got %T", step.Value)
	}
	if lit.Value != 3.14 {
		t.Errorf("expected 3.14, got %v", lit.Value)
	}
}

// P103: trueリテラル
func TestParser_Expr_Literal_True(t *testing.T) {
	input := `component Foo { flow F { x = true } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	lit, ok := step.Value.(*ast.LiteralExpr)
	if !ok {
		t.Fatalf("expected LiteralExpr, got %T", step.Value)
	}
	if lit.Value != true {
		t.Errorf("expected true, got %v", lit.Value)
	}
}

// P104: falseリテラル
func TestParser_Expr_Literal_False(t *testing.T) {
	input := `component Foo { flow F { x = false } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	lit := step.Value.(*ast.LiteralExpr)
	if lit.Value != false {
		t.Errorf("expected false, got %v", lit.Value)
	}
}

// P105: nullリテラル
func TestParser_Expr_Literal_Null(t *testing.T) {
	input := `component Foo { flow F { x = null } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	lit := step.Value.(*ast.LiteralExpr)
	if lit.Value != nil {
		t.Errorf("expected nil, got %v", lit.Value)
	}
}

// P106: 変数
func TestParser_Expr_Variable(t *testing.T) {
	input := `component Foo { flow F { x = foo } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	v, ok := step.Value.(*ast.VariableExpr)
	if !ok {
		t.Fatalf("expected VariableExpr, got %T", step.Value)
	}
	if v.Name != "foo" {
		t.Errorf("expected 'foo', got %q", v.Name)
	}
}

// P107: フィールドアクセス
func TestParser_Expr_Field(t *testing.T) {
	input := `component Foo { flow F { x = a.b } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	f, ok := step.Value.(*ast.FieldExpr)
	if !ok {
		t.Fatalf("expected FieldExpr, got %T", step.Value)
	}
	if f.Field != "b" {
		t.Errorf("expected field 'b', got %q", f.Field)
	}
}

// P108: チェーン
func TestParser_Expr_Field_Chain(t *testing.T) {
	input := `component Foo { flow F { x = a.b.c } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	f, ok := step.Value.(*ast.FieldExpr)
	if !ok {
		t.Fatalf("expected FieldExpr, got %T", step.Value)
	}
	if f.Field != "c" {
		t.Errorf("expected field 'c', got %q", f.Field)
	}
	f2, ok := f.Object.(*ast.FieldExpr)
	if !ok {
		t.Fatalf("expected nested FieldExpr")
	}
	if f2.Field != "b" {
		t.Errorf("expected field 'b', got %q", f2.Field)
	}
}

// P109: 呼び出し
func TestParser_Expr_Call(t *testing.T) {
	input := `component Foo { flow F { x = A.B() } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	c, ok := step.Value.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr, got %T", step.Value)
	}
	if c.Method != "B" {
		t.Errorf("expected method 'B', got %q", c.Method)
	}
}

// P110-P122: 二項演算
func TestParser_Expr_Binary_Operators(t *testing.T) {
	tests := []struct {
		op     string
		input  string
	}{
		{"+", "a + b"},
		{"-", "a - b"},
		{"*", "a * b"},
		{"/", "a / b"},
		{"%", "a % b"},
		{"==", "a == b"},
		{"!=", "a != b"},
		{"<", "a < b"},
		{">", "a > b"},
		{"<=", "a <= b"},
		{">=", "a >= b"},
		{"&&", "a && b"},
		{"||", "a || b"},
	}
	for _, tt := range tests {
		input := `component Foo { flow F { x = ` + tt.input + ` } }`
		spec, err := ParseString(input)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", tt.op, err)
			continue
		}
		step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
		b, ok := step.Value.(*ast.BinaryExpr)
		if !ok {
			t.Errorf("expected BinaryExpr for %q, got %T", tt.op, step.Value)
			continue
		}
		if b.Op != tt.op {
			t.Errorf("expected op %q, got %q", tt.op, b.Op)
		}
	}
}

// P123: 否定
func TestParser_Expr_Unary_Not(t *testing.T) {
	input := `component Foo { flow F { x = !a } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	u, ok := step.Value.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected UnaryExpr, got %T", step.Value)
	}
	if u.Op != "!" {
		t.Errorf("expected op '!', got %q", u.Op)
	}
}

// P124: 負数
func TestParser_Expr_Unary_Neg(t *testing.T) {
	input := `component Foo { flow F { x = -a } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	u, ok := step.Value.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected UnaryExpr, got %T", step.Value)
	}
	if u.Op != "-" {
		t.Errorf("expected op '-', got %q", u.Op)
	}
}

// P125: 三項演算子
func TestParser_Expr_Ternary(t *testing.T) {
	input := `component Foo { flow F { x = a ? b : c } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	_, ok := step.Value.(*ast.TernaryExpr)
	if !ok {
		t.Fatalf("expected TernaryExpr, got %T", step.Value)
	}
}

// P126: 括弧
func TestParser_Expr_Paren(t *testing.T) {
	input := `component Foo { flow F { x = (a + b) } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	b, ok := step.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", step.Value)
	}
	if b.Op != "+" {
		t.Errorf("expected op '+', got %q", b.Op)
	}
}

// P127: 優先順位（乗算優先）
func TestParser_Expr_Precedence_MulAdd(t *testing.T) {
	input := `component Foo { flow F { x = a + b * c } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	b, ok := step.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", step.Value)
	}
	// a + (b * c) なので外側は +
	if b.Op != "+" {
		t.Errorf("expected top op '+', got %q", b.Op)
	}
	// 右側が b * c
	right, ok := b.Right.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected right BinaryExpr")
	}
	if right.Op != "*" {
		t.Errorf("expected right op '*', got %q", right.Op)
	}
}

// P128: 優先順位（AND優先）
func TestParser_Expr_Precedence_AndOr(t *testing.T) {
	input := `component Foo { flow F { x = a || b && c } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	b, ok := step.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", step.Value)
	}
	// a || (b && c) なので外側は ||
	if b.Op != "||" {
		t.Errorf("expected top op '||', got %q", b.Op)
	}
}

// P129: 比較と論理
func TestParser_Expr_Precedence_Compare(t *testing.T) {
	input := `component Foo { flow F { x = a == b && c < d } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	b, ok := step.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", step.Value)
	}
	// (a == b) && (c < d) なので外側は &&
	if b.Op != "&&" {
		t.Errorf("expected top op '&&', got %q", b.Op)
	}
}

// P130: 複合式
func TestParser_Expr_Complex(t *testing.T) {
	input := `component Foo { flow F { x = a.b + c.d() * 2 } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	_, ok := step.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", step.Value)
	}
}

// P131: 位置情報
func TestParser_Expr_Position(t *testing.T) {
	input := `component Foo { flow F { x = a } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	step := spec.Component.Body.Flows[0].Steps[0].(*ast.AssignStep)
	if !step.Value.GetPos().IsValid() {
		t.Errorf("expected valid position")
	}
}

// =============================================================================
// 1.2.8 ステートマシン定義
// =============================================================================

// P140: 最小ステートマシン
func TestParser_States_Empty(t *testing.T) {
	input := `component Foo { states S { initial I } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Component.Body.States) != 1 {
		t.Fatalf("expected 1 states, got %d", len(spec.Component.Body.States))
	}
	states := spec.Component.Body.States[0]
	if states.Name != "S" {
		t.Errorf("expected name 'S', got %q", states.Name)
	}
	if states.Initial != "I" {
		t.Errorf("expected initial 'I', got %q", states.Initial)
	}
}

// P141: final付き
func TestParser_States_WithFinal(t *testing.T) {
	input := `component Foo { states S { initial I final F } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	states := spec.Component.Body.States[0]
	if len(states.Finals) != 1 || states.Finals[0] != "F" {
		t.Errorf("expected final 'F'")
	}
}

// P142: 複数final
func TestParser_States_MultipleFinal(t *testing.T) {
	input := `component Foo { states S { initial I final F1 final F2 } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	states := spec.Component.Body.States[0]
	if len(states.Finals) != 2 {
		t.Errorf("expected 2 finals, got %d", len(states.Finals))
	}
}

// P143: イベントトリガー
func TestParser_States_Transition_OnEvent(t *testing.T) {
	input := `component Foo { states S { initial I I -> A on E } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trans := spec.Component.Body.States[0].Transitions[0]
	trigger, ok := trans.Trigger.(*ast.EventTrigger)
	if !ok {
		t.Fatalf("expected EventTrigger, got %T", trans.Trigger)
	}
	if trigger.Event != "E" {
		t.Errorf("expected event 'E', got %q", trigger.Event)
	}
}

// P144: 時間トリガー
func TestParser_States_Transition_After(t *testing.T) {
	input := `component Foo { states S { initial I I -> A after 3s } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trans := spec.Component.Body.States[0].Transitions[0]
	trigger, ok := trans.Trigger.(*ast.AfterTrigger)
	if !ok {
		t.Fatalf("expected AfterTrigger, got %T", trans.Trigger)
	}
	if trigger.Duration.Value != 3 || trigger.Duration.Unit != "s" {
		t.Errorf("expected 3s, got %d%s", trigger.Duration.Value, trigger.Duration.Unit)
	}
}

// P145: 条件トリガー
func TestParser_States_Transition_When(t *testing.T) {
	input := `component Foo { states S { initial I I -> A when cond } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trans := spec.Component.Body.States[0].Transitions[0]
	_, ok := trans.Trigger.(*ast.WhenTrigger)
	if !ok {
		t.Fatalf("expected WhenTrigger, got %T", trans.Trigger)
	}
}

// P146: ガード条件
func TestParser_States_Transition_Guard(t *testing.T) {
	input := `component Foo { states S { initial I I -> A on E when g } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trans := spec.Component.Body.States[0].Transitions[0]
	if trans.Guard == nil {
		t.Errorf("expected guard")
	}
}

// P147: アクション
func TestParser_States_Transition_Actions(t *testing.T) {
	input := `component Foo { states S { initial I I -> A on E do [a, b] } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trans := spec.Component.Body.States[0].Transitions[0]
	if len(trans.Actions) != 2 {
		t.Errorf("expected 2 actions, got %d", len(trans.Actions))
	}
}

// P148: フル遷移
func TestParser_States_Transition_Full(t *testing.T) {
	input := `component Foo { states S { initial I I -> A on E when g do [a] } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trans := spec.Component.Body.States[0].Transitions[0]
	if trans.Guard == nil {
		t.Errorf("expected guard")
	}
	if len(trans.Actions) != 1 {
		t.Errorf("expected 1 action")
	}
}

// P149: 空状態定義
func TestParser_States_State_Empty(t *testing.T) {
	input := `component Foo { states S { initial I state I { } } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Component.Body.States[0].States) != 1 {
		t.Fatalf("expected 1 state, got %d", len(spec.Component.Body.States[0].States))
	}
}

// P150: entry付き
func TestParser_States_State_Entry(t *testing.T) {
	input := `component Foo { states S { initial I state I { entry [a] } } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	state := spec.Component.Body.States[0].States[0]
	if len(state.Entry) != 1 {
		t.Errorf("expected 1 entry action, got %d", len(state.Entry))
	}
}

// P151: exit付き
func TestParser_States_State_Exit(t *testing.T) {
	input := `component Foo { states S { initial I state I { exit [a] } } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	state := spec.Component.Body.States[0].States[0]
	if len(state.Exit) != 1 {
		t.Errorf("expected 1 exit action, got %d", len(state.Exit))
	}
}

// P152: entry/exit両方
func TestParser_States_State_EntryExit(t *testing.T) {
	input := `component Foo { states S { initial I state I { entry [a] exit [b] } } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	state := spec.Component.Body.States[0].States[0]
	if len(state.Entry) != 1 || len(state.Exit) != 1 {
		t.Errorf("expected entry and exit actions")
	}
}

// P153: 階層状態
func TestParser_States_Compound(t *testing.T) {
	input := `component Foo { states S { initial I state I { state J { } } } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	state := spec.Component.Body.States[0].States[0]
	if len(state.States) != 1 {
		t.Errorf("expected 1 nested state")
	}
}

// P154: 階層initial
func TestParser_States_Compound_Initial(t *testing.T) {
	input := `component Foo { states S { initial I state I { initial J state J { } } } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	state := spec.Component.Body.States[0].States[0]
	if state.Initial == nil || *state.Initial != "J" {
		t.Errorf("expected initial 'J'")
	}
}

// P155: 並行状態
func TestParser_States_Parallel(t *testing.T) {
	input := `component Foo { states S { initial I parallel P { region R1 { initial A } region R2 { initial B } } } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Component.Body.States[0].Parallels) != 1 {
		t.Fatalf("expected 1 parallel")
	}
	parallel := spec.Component.Body.States[0].Parallels[0]
	if len(parallel.Regions) != 2 {
		t.Errorf("expected 2 regions, got %d", len(parallel.Regions))
	}
}

// P156: リージョン
func TestParser_States_Region(t *testing.T) {
	input := `component Foo { states S { initial I parallel P { region R { initial A A -> B on E } } } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	region := spec.Component.Body.States[0].Parallels[0].Regions[0]
	if region.Name != "R" {
		t.Errorf("expected name 'R', got %q", region.Name)
	}
	if region.Initial != "A" {
		t.Errorf("expected initial 'A', got %q", region.Initial)
	}
	if len(region.Transitions) != 1 {
		t.Errorf("expected 1 transition")
	}
}

// P157: 複合ステートマシン
func TestParser_States_Complex(t *testing.T) {
	input := `component Foo {
		states OrderState {
			initial Pending
			final Completed
			final Cancelled

			Pending -> Processing on Approve
			Processing -> Completed on Ship when validated
			Processing -> Cancelled on Cancel do [notifyUser]

			state Processing {
				entry [startProcessing]
				exit [logExit]
			}
		}
	}`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	states := spec.Component.Body.States[0]
	if states.Initial != "Pending" {
		t.Errorf("expected initial 'Pending'")
	}
	if len(states.Finals) != 2 {
		t.Errorf("expected 2 finals")
	}
	if len(states.Transitions) != 3 {
		t.Errorf("expected 3 transitions")
	}
}

// P158: 位置情報
func TestParser_States_Position(t *testing.T) {
	input := `component Foo { states S { initial I } }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	states := spec.Component.Body.States[0]
	if !states.Pos.IsValid() {
		t.Errorf("expected valid position")
	}
}

// =============================================================================
// 1.2.9 アノテーション
// =============================================================================

// P160: フラグアノテーション
func TestParser_Annotation_Flag(t *testing.T) {
	input := `@deprecated component Foo { }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ann := spec.Component.Annotations[0]
	if ann.Name != "deprecated" {
		t.Errorf("expected name 'deprecated', got %q", ann.Name)
	}
	if len(ann.Args) != 0 {
		t.Errorf("expected 0 args, got %d", len(ann.Args))
	}
}

// P161: 単一値
func TestParser_Annotation_SingleValue(t *testing.T) {
	input := `@description("hello") component Foo { }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ann := spec.Component.Annotations[0]
	if len(ann.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(ann.Args))
	}
	if ann.Args[0].Value != "hello" {
		t.Errorf("expected value 'hello', got %q", ann.Args[0].Value)
	}
}

// P162: 複数値
func TestParser_Annotation_MultiValue(t *testing.T) {
	input := `@throws("A", "B") component Foo { }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ann := spec.Component.Annotations[0]
	if len(ann.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(ann.Args))
	}
}

// P163: 名前付き
func TestParser_Annotation_Named(t *testing.T) {
	input := `@source(file: "a.go") component Foo { }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ann := spec.Component.Annotations[0]
	if len(ann.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(ann.Args))
	}
	if ann.Args[0].Key == nil || *ann.Args[0].Key != "file" {
		t.Errorf("expected key 'file'")
	}
	if ann.Args[0].Value != "a.go" {
		t.Errorf("expected value 'a.go', got %q", ann.Args[0].Value)
	}
}

// P164: 混合
func TestParser_Annotation_Mixed(t *testing.T) {
	input := `@foo("a", b: "c") component Foo { }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ann := spec.Component.Annotations[0]
	if len(ann.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(ann.Args))
	}
	if ann.Args[0].Key != nil {
		t.Errorf("expected positional arg")
	}
	if ann.Args[1].Key == nil {
		t.Errorf("expected named arg")
	}
}

// P165: 複数アノテーション
func TestParser_Annotation_Multiple(t *testing.T) {
	input := `@deprecated @description("x") component Foo { }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Component.Annotations) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(spec.Component.Annotations))
	}
}

// P166: 位置情報
func TestParser_Annotation_Position(t *testing.T) {
	input := `@deprecated component Foo { }`
	spec, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ann := spec.Component.Annotations[0]
	if !ann.Pos.IsValid() {
		t.Errorf("expected valid position")
	}
}

// =============================================================================
// 1.2.10 エラーケース
// =============================================================================

// P200: 予期しないトークン
func TestParser_Error_UnexpectedToken(t *testing.T) {
	input := `component { }`
	_, err := ParseString(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !isParseError(err) {
		t.Fatalf("expected ParseError or MultiError containing ParseError, got %T", err)
	}
	// エラーメッセージに "identifier" が含まれることを期待
	_ = getFirstParseError(err) // 実装時にメッセージを確認
}

// P201: 括弧閉じ忘れ
func TestParser_Error_UnterminatedBrace(t *testing.T) {
	input := `component Foo {`
	_, err := ParseString(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// P202: 不正な式
func TestParser_Error_InvalidExpression(t *testing.T) {
	input := `component Foo { flow F { x = } }`
	_, err := ParseString(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// P203: エラー位置
func TestParser_Error_Position(t *testing.T) {
	input := `component { }`
	_, err := ParseString(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	pe := getFirstParseError(err)
	if pe == nil {
		t.Fatalf("expected ParseError or MultiError containing ParseError, got %T", err)
	}
	if !pe.Pos.IsValid() {
		t.Errorf("expected valid error position")
	}
}

// P204: エラー型
func TestParser_Error_Type(t *testing.T) {
	input := `component { }`
	_, err := ParseString(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !isParseError(err) {
		t.Errorf("expected ParseError or MultiError containing ParseError, got %T", err)
	}
}
