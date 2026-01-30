package class

import "pact/internal/domain/diagram/common"

// Diagram はクラス図を表す
type Diagram struct {
	Nodes []Node
	Edges []Edge
}

func (d *Diagram) Type() common.DiagramType {
	return common.DiagramTypeClass
}

// Node はクラス図のノード
type Node struct {
	ID          string
	Name        string
	Stereotype  string
	Attributes  []Attribute
	Methods     []Method
	Annotations []common.Annotation
}

// Attribute はクラスの属性
type Attribute struct {
	Name       string
	Type       string
	Visibility Visibility
}

// Method はクラスのメソッド
type Method struct {
	Name       string
	Params     []Param
	ReturnType string
	Visibility Visibility
	Async      bool
}

// Param はメソッドのパラメータ
type Param struct {
	Name string
	Type string
}

// Visibility は可視性
type Visibility string

const (
	VisibilityPublic    Visibility = "public"
	VisibilityPrivate   Visibility = "private"
	VisibilityProtected Visibility = "protected"
	VisibilityPackage   Visibility = "package"
)

// Edge はクラス図のエッジ
type Edge struct {
	From       string
	To         string
	Type       EdgeType
	Label      string
	Decoration Decoration
	LineStyle  LineStyle
}

// EdgeType はエッジの種類
type EdgeType string

const (
	EdgeTypeDependency     EdgeType = "dependency"
	EdgeTypeInheritance    EdgeType = "inheritance"
	EdgeTypeImplementation EdgeType = "implementation"
	EdgeTypeComposition    EdgeType = "composition"
	EdgeTypeAggregation    EdgeType = "aggregation"
)

// Decoration はエッジの装飾
type Decoration string

const (
	DecorationNone          Decoration = "none"
	DecorationArrow         Decoration = "arrow"
	DecorationTriangle      Decoration = "triangle"
	DecorationFilledDiamond Decoration = "filled_diamond"
	DecorationEmptyDiamond  Decoration = "empty_diamond"
)

// LineStyle は線のスタイル
type LineStyle string

const (
	LineStyleSolid  LineStyle = "solid"
	LineStyleDashed LineStyle = "dashed"
)
