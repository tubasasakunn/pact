package ast

import "fmt"

// Position はソースコード内の位置を表す
type Position struct {
	File   string
	Line   int
	Column int
	Offset int
}

func (p Position) String() string {
	if p.File == "" {
		return fmt.Sprintf("%d:%d", p.Line, p.Column)
	}
	return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Column)
}

// NoPos は位置情報がないことを表す
var NoPos = Position{}

func (p Position) IsValid() bool {
	return p.Line > 0
}
