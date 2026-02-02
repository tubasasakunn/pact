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

// Note は図上のノート（コメント）
type Note struct {
	ID       string
	Text     string
	Position NotePosition
	AttachTo string // 関連付ける要素のID（オプション）
}

// NotePosition はノートの配置位置
type NotePosition string

const (
	NotePositionLeft   NotePosition = "left"
	NotePositionRight  NotePosition = "right"
	NotePositionTop    NotePosition = "top"
	NotePositionBottom NotePosition = "bottom"
	NotePositionOver   NotePosition = "over" // シーケンス図用
)
