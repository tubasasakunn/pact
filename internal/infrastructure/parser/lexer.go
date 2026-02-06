package parser

// Lexer は字句解析器
type Lexer struct {
	input       string
	pos         int
	readPos     int
	ch          byte
	line        int
	column      int
	tokenLine   int
	tokenColumn int
	err         error // レキシングエラー
}

// キーワードマップ
var keywords = map[string]TokenType{
	"component":  TOKEN_COMPONENT,
	"import":     TOKEN_IMPORT,
	"type":       TOKEN_TYPE,
	"enum":       TOKEN_ENUM,
	"depends":    TOKEN_DEPENDS,
	"on":         TOKEN_ON,
	"extends":    TOKEN_EXTENDS,
	"implements": TOKEN_IMPLEMENTS,
	"contains":   TOKEN_CONTAINS,
	"aggregates": TOKEN_AGGREGATES,
	"provides":   TOKEN_PROVIDES,
	"requires":   TOKEN_REQUIRES,
	"flow":       TOKEN_FLOW,
	"states":     TOKEN_STATES,
	"state":      TOKEN_STATE,
	"parallel":   TOKEN_PARALLEL,
	"region":     TOKEN_REGION,
	"initial":    TOKEN_INITIAL,
	"final":      TOKEN_FINAL,
	"entry":      TOKEN_ENTRY,
	"exit":       TOKEN_EXIT,
	"if":         TOKEN_IF,
	"else":       TOKEN_ELSE,
	"for":        TOKEN_FOR,
	"in":         TOKEN_IN,
	"while":      TOKEN_WHILE,
	"return":     TOKEN_RETURN,
	"throw":      TOKEN_THROW,
	"await":      TOKEN_AWAIT,
	"async":      TOKEN_ASYNC,
	"throws":     TOKEN_THROWS,
	"when":       TOKEN_WHEN,
	"after":      TOKEN_AFTER,
	"do":         TOKEN_DO,
	"true":       TOKEN_TRUE,
	"false":      TOKEN_FALSE,
	"null":       TOKEN_NULL,
	"as":         TOKEN_AS,
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
	l.skipWhitespaceAndComments()

	// トークン開始位置を記録
	l.tokenLine = l.line
	l.tokenColumn = l.column

	var tok Token
	tok.Line = l.tokenLine
	tok.Column = l.tokenColumn

	switch l.ch {
	case 0:
		tok.Type = TOKEN_EOF
		tok.Literal = ""
	case '{':
		tok = l.newToken(TOKEN_LBRACE, "{")
		l.readChar()
	case '}':
		tok = l.newToken(TOKEN_RBRACE, "}")
		l.readChar()
	case '(':
		tok = l.newToken(TOKEN_LPAREN, "(")
		l.readChar()
	case ')':
		tok = l.newToken(TOKEN_RPAREN, ")")
		l.readChar()
	case '[':
		tok = l.newToken(TOKEN_LBRACKET, "[")
		l.readChar()
	case ']':
		tok = l.newToken(TOKEN_RBRACKET, "]")
		l.readChar()
	case ':':
		tok = l.newToken(TOKEN_COLON, ":")
		l.readChar()
	case ',':
		tok = l.newToken(TOKEN_COMMA, ",")
		l.readChar()
	case '.':
		tok = l.newToken(TOKEN_DOT, ".")
		l.readChar()
	case '@':
		tok = l.newToken(TOKEN_AT, "@")
		l.readChar()
	case '?':
		if l.peekChar() == '?' {
			l.readChar()
			tok = l.newToken(TOKEN_NULLISH, "??")
			l.readChar()
		} else {
			tok = l.newToken(TOKEN_QUESTION, "?")
			l.readChar()
		}
	case '+':
		tok = l.newToken(TOKEN_PLUS, "+")
		l.readChar()
	case '-':
		if l.peekChar() == '>' {
			l.readChar()
			tok = l.newToken(TOKEN_ARROW, "->")
			l.readChar()
		} else {
			tok = l.newToken(TOKEN_MINUS, "-")
			l.readChar()
		}
	case '*':
		tok = l.newToken(TOKEN_STAR, "*")
		l.readChar()
	case '/':
		tok = l.newToken(TOKEN_SLASH, "/")
		l.readChar()
	case '%':
		tok = l.newToken(TOKEN_PERCENT, "%")
		l.readChar()
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(TOKEN_EQ, "==")
			l.readChar()
		} else {
			tok = l.newToken(TOKEN_ASSIGN, "=")
			l.readChar()
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(TOKEN_NE, "!=")
			l.readChar()
		} else {
			tok = l.newToken(TOKEN_NOT, "!")
			l.readChar()
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(TOKEN_LE, "<=")
			l.readChar()
		} else {
			tok = l.newToken(TOKEN_LT, "<")
			l.readChar()
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(TOKEN_GE, ">=")
			l.readChar()
		} else {
			tok = l.newToken(TOKEN_GT, ">")
			l.readChar()
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = l.newToken(TOKEN_AND, "&&")
			l.readChar()
		} else {
			tok = l.newToken(TOKEN_ILLEGAL, string(l.ch))
			l.readChar()
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = l.newToken(TOKEN_OR, "||")
			l.readChar()
		} else {
			tok = l.newToken(TOKEN_ILLEGAL, string(l.ch))
			l.readChar()
		}
	case '#':
		tok = l.newToken(TOKEN_HASH, "#")
		l.readChar()
	case '~':
		tok = l.newToken(TOKEN_TILDE, "~")
		l.readChar()
	case '"':
		tok.Type = TOKEN_STRING
		tok.Literal = l.readString()
		tok.Line = l.tokenLine
		tok.Column = l.tokenColumn
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tok.Literal = literal
			tok.Type = lookupIdent(literal)
			tok.Line = l.tokenLine
			tok.Column = l.tokenColumn
			return tok
		} else if isDigit(l.ch) {
			return l.readNumber()
		} else {
			tok = l.newToken(TOKEN_ILLEGAL, string(l.ch))
			l.readChar()
		}
	}

	return tok
}

func (l *Lexer) newToken(tokenType TokenType, literal string) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
		Line:    l.tokenLine,
		Column:  l.tokenColumn,
	}
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		// 空白をスキップ
		for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			l.readChar()
		}

		// コメントをスキップ
		if l.ch == '/' && l.peekChar() == '/' {
			// 行コメント
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}

		if l.ch == '/' && l.peekChar() == '*' {
			// ブロックコメント
			commentLine := l.line
			commentColumn := l.column
			l.readChar() // consume '/'
			l.readChar() // consume '*'
			closed := false
			for {
				if l.ch == 0 {
					// 閉じられていないブロックコメント
					l.err = &lexerError{
						line:    commentLine,
						column:  commentColumn,
						message: "unclosed block comment",
					}
					break
				}
				if l.ch == '*' && l.peekChar() == '/' {
					l.readChar() // consume '*'
					l.readChar() // consume '/'
					closed = true
					break
				}
				l.readChar()
			}
			if closed {
				continue
			}
			break
		}

		break
	}
}

// lexerError はレキシングエラー
type lexerError struct {
	line    int
	column  int
	message string
}

func (e *lexerError) Error() string {
	return e.message
}

// Error はLexerのエラーを返す
func (l *Lexer) Error() error {
	return l.err
}

func (l *Lexer) readIdentifier() string {
	position := l.pos
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.pos]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TOKEN_IDENT
}
