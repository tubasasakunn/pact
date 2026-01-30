package parser

// Lexer は字句解析器
type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
	line    int
	column  int
}

// NewLexer は新しいLexerを作成する
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// NextToken は次のトークンを返す
func (l *Lexer) NextToken() Token {
	// TODO: 実装
	return Token{Type: TOKEN_EOF}
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
	l.column++
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}
