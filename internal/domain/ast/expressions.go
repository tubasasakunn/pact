package ast

// Expr は式を表すインターフェース
type Expr interface {
	exprNode()
	GetPos() Position
}

// LiteralExpr はリテラル式を表す
type LiteralExpr struct {
	Pos   Position
	Value interface{}
}

func (e *LiteralExpr) exprNode()        {}
func (e *LiteralExpr) GetPos() Position { return e.Pos }

// VariableExpr は変数参照を表す
type VariableExpr struct {
	Pos  Position
	Name string
}

func (e *VariableExpr) exprNode()        {}
func (e *VariableExpr) GetPos() Position { return e.Pos }

// FieldExpr はフィールドアクセスを表す
type FieldExpr struct {
	Pos    Position
	Object Expr
	Field  string
}

func (e *FieldExpr) exprNode()        {}
func (e *FieldExpr) GetPos() Position { return e.Pos }

// CallExpr はメソッド呼び出しを表す
type CallExpr struct {
	Pos      Position
	Object   Expr
	Method   string
	Args     []Expr
}

func (e *CallExpr) exprNode()        {}
func (e *CallExpr) GetPos() Position { return e.Pos }

// BinaryExpr は二項演算式を表す
type BinaryExpr struct {
	Pos   Position
	Left  Expr
	Op    string
	Right Expr
}

func (e *BinaryExpr) exprNode()        {}
func (e *BinaryExpr) GetPos() Position { return e.Pos }

// UnaryExpr は単項演算式を表す
type UnaryExpr struct {
	Pos     Position
	Op      string
	Operand Expr
}

func (e *UnaryExpr) exprNode()        {}
func (e *UnaryExpr) GetPos() Position { return e.Pos }

// TernaryExpr は三項演算式を表す
type TernaryExpr struct {
	Pos       Position
	Condition Expr
	Then      Expr
	Else      Expr
}

func (e *TernaryExpr) exprNode()        {}
func (e *TernaryExpr) GetPos() Position { return e.Pos }

// NullishExpr はnull合体演算を表す
type NullishExpr struct {
	Pos      Position
	Left     Expr
	Right    Expr
	ThrowErr *string // ?? throw Error の場合
}

func (e *NullishExpr) exprNode()        {}
func (e *NullishExpr) GetPos() Position { return e.Pos }
