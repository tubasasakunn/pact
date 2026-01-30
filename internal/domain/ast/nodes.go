package ast

// SpecFile は .pact ファイル全体を表す
type SpecFile struct {
	Path      string
	Imports   []ImportDecl
	Component *ComponentDecl
}

// ImportDecl は import 文を表す
type ImportDecl struct {
	Pos   Position
	Path  string
	Alias *string
}

// ComponentDecl は component 宣言を表す
type ComponentDecl struct {
	Pos         Position
	Name        string
	Annotations []AnnotationDecl
	Body        ComponentBody
}

// ComponentBody は component の中身を表す
type ComponentBody struct {
	Types     []TypeDecl
	Relations []RelationDecl
	Provides  []InterfaceDecl
	Requires  []InterfaceDecl
	Flows     []FlowDecl
	States    []StatesDecl
}

// AnnotationDecl はアノテーションを表す
type AnnotationDecl struct {
	Pos  Position
	Name string
	Args []AnnotationArg
}

// AnnotationArg はアノテーションの引数を表す
type AnnotationArg struct {
	Key   *string // nil なら positional
	Value string
}
