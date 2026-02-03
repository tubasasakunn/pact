package ast

// TypeDecl は型定義を表す
type TypeDecl struct {
	Pos         Position
	Name        string
	Kind        TypeKind
	Annotations []AnnotationDecl
	Fields      []FieldDecl // struct の場合
	Values      []string    // enum の場合
	BaseType    *TypeExpr   // alias の場合
}

type TypeKind string

const (
	TypeKindStruct TypeKind = "struct"
	TypeKindEnum   TypeKind = "enum"
	TypeKindAlias  TypeKind = "alias"
)

// FieldDecl はフィールド定義を表す
type FieldDecl struct {
	Pos         Position
	Name        string
	Type        TypeExpr
	Visibility  Visibility
	Annotations []AnnotationDecl
}

type Visibility string

const (
	VisibilityPublic    Visibility = "public"
	VisibilityPrivate   Visibility = "private"
	VisibilityProtected Visibility = "protected"
	VisibilityPackage   Visibility = "package"
)

// TypeExpr は型の参照を表す
type TypeExpr struct {
	Pos        Position
	Name       string
	Nullable   bool
	Array      bool
	TypeParams []TypeExpr // ジェネリクス型パラメータ（例: List<T>）
}

// RelationDecl は関係定義を表す
type RelationDecl struct {
	Pos         Position
	Kind        RelationKind
	Target      string
	TargetType  *string
	Alias       *string
	Annotations []AnnotationDecl
}

type RelationKind string

const (
	RelationDependsOn   RelationKind = "depends_on"
	RelationExtends     RelationKind = "extends"
	RelationImplements  RelationKind = "implements"
	RelationContains    RelationKind = "contains"
	RelationAggregates  RelationKind = "aggregates"
)

// InterfaceDecl はインターフェース定義を表す
type InterfaceDecl struct {
	Pos         Position
	Name        string
	Annotations []AnnotationDecl
	Methods     []MethodDecl
}

// MethodDecl はメソッド定義を表す
type MethodDecl struct {
	Pos         Position
	Name        string
	Params      []ParamDecl
	ReturnType  *TypeExpr
	Throws      []string
	Async       bool
	Annotations []AnnotationDecl
}

// ParamDecl はパラメータ定義を表す
type ParamDecl struct {
	Pos  Position
	Name string
	Type TypeExpr
}
