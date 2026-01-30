package ast

// FlowDecl はフロー定義を表す
type FlowDecl struct {
	Pos         Position
	Name        string
	Annotations []AnnotationDecl
	Steps       []Step
}

// Step はフローのステップを表すインターフェース
type Step interface {
	stepNode()
	GetPos() Position
}

// AssignStep は代入ステップを表す
type AssignStep struct {
	Pos         Position
	Variable    string
	Value       Expr
	Annotations []AnnotationDecl
}

func (s *AssignStep) stepNode()        {}
func (s *AssignStep) GetPos() Position { return s.Pos }

// CallStep は呼び出しステップを表す
type CallStep struct {
	Pos         Position
	Expr        Expr
	Await       bool
	Annotations []AnnotationDecl
}

func (s *CallStep) stepNode()        {}
func (s *CallStep) GetPos() Position { return s.Pos }

// ReturnStep はreturnステップを表す
type ReturnStep struct {
	Pos         Position
	Value       Expr // nil の場合は値なし
	Annotations []AnnotationDecl
}

func (s *ReturnStep) stepNode()        {}
func (s *ReturnStep) GetPos() Position { return s.Pos }

// ThrowStep はthrowステップを表す
type ThrowStep struct {
	Pos         Position
	Error       string
	Annotations []AnnotationDecl
}

func (s *ThrowStep) stepNode()        {}
func (s *ThrowStep) GetPos() Position { return s.Pos }

// IfStep は条件分岐ステップを表す
type IfStep struct {
	Pos         Position
	Condition   Expr
	Then        []Step
	Else        []Step
	Annotations []AnnotationDecl
}

func (s *IfStep) stepNode()        {}
func (s *IfStep) GetPos() Position { return s.Pos }

// ForStep はforループステップを表す
type ForStep struct {
	Pos         Position
	Variable    string
	Iterable    Expr
	Body        []Step
	Annotations []AnnotationDecl
}

func (s *ForStep) stepNode()        {}
func (s *ForStep) GetPos() Position { return s.Pos }

// WhileStep はwhileループステップを表す
type WhileStep struct {
	Pos         Position
	Condition   Expr
	Body        []Step
	Annotations []AnnotationDecl
}

func (s *WhileStep) stepNode()        {}
func (s *WhileStep) GetPos() Position { return s.Pos }
