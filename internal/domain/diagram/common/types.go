package common

// Diagram は図の共通インターフェース
type Diagram interface {
	Type() DiagramType
}

type DiagramType string

const (
	DiagramTypeClass    DiagramType = "class"
	DiagramTypeSequence DiagramType = "sequence"
	DiagramTypeState    DiagramType = "state"
	DiagramTypeFlow     DiagramType = "flow"
)

// Annotation は図の注釈
type Annotation struct {
	Name string
	Args map[string]string
}
