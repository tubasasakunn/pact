package ast

// StatesDecl はステートマシン定義を表す
type StatesDecl struct {
	Pos         Position
	Name        string
	Annotations []AnnotationDecl
	Initial     string
	Finals      []string
	States      []StateDecl
	Transitions []TransitionDecl
	Parallels   []ParallelDecl
}

// StateDecl は状態定義を表す
type StateDecl struct {
	Pos         Position
	Name        string
	Annotations []AnnotationDecl
	Entry       []string
	Exit        []string
	// 階層状態の場合
	Initial     *string
	States      []StateDecl
	Transitions []TransitionDecl
}

// TransitionDecl は状態遷移定義を表す
type TransitionDecl struct {
	Pos     Position
	From    string
	To      string
	Trigger Trigger
	Guard   Expr
	Actions []string
}

// Trigger はトリガーを表すインターフェース
type Trigger interface {
	triggerNode()
	GetPos() Position
}

// EventTrigger はイベントトリガーを表す
type EventTrigger struct {
	Pos   Position
	Event string
}

func (t *EventTrigger) triggerNode()     {}
func (t *EventTrigger) GetPos() Position { return t.Pos }

// AfterTrigger は時間トリガーを表す
type AfterTrigger struct {
	Pos      Position
	Duration Duration
}

func (t *AfterTrigger) triggerNode()     {}
func (t *AfterTrigger) GetPos() Position { return t.Pos }

// Duration は期間を表す
type Duration struct {
	Value int
	Unit  string // "ms", "s", "m", "h", "d"
}

// WhenTrigger は条件トリガーを表す
type WhenTrigger struct {
	Pos       Position
	Condition Expr
}

func (t *WhenTrigger) triggerNode()     {}
func (t *WhenTrigger) GetPos() Position { return t.Pos }

// ParallelDecl は並行状態を表す
type ParallelDecl struct {
	Pos         Position
	Name        string
	Annotations []AnnotationDecl
	Regions     []RegionDecl
}

// RegionDecl はリージョンを表す
type RegionDecl struct {
	Pos         Position
	Name        string
	Initial     string
	States      []StateDecl
	Transitions []TransitionDecl
}
