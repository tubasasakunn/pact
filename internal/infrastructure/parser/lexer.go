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

func (l *Lexer) readNumber() Token {
	startLine := l.tokenLine
	startColumn := l.tokenColumn
	position := l.pos
	isFloat := false

	// 整数部を読む
	for isDigit(l.ch) {
		l.readChar()
	}

	// 小数点以下があるか確認
	if l.ch == '.' && isDigit(l.peekChar()) {
		isFloat = true
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	// 指数部（scientific notation）があるか確認: e/E followed by optional +/- and digits
	if l.ch == 'e' || l.ch == 'E' {
		next := l.peekChar()
		if isDigit(next) || next == '+' || next == '-' {
			isFloat = true
			l.readChar() // consume 'e' or 'E'
			if l.ch == '+' || l.ch == '-' {
				l.readChar() // consume sign
			}
			for isDigit(l.ch) {
				l.readChar()
			}
		}
	}

	if isFloat {
		return Token{
			Type:    TOKEN_FLOAT,
			Literal: l.input[position:l.pos],
			Line:    startLine,
			Column:  startColumn,
		}
	}

	// 期間リテラルかチェック (ms, s, m, h, d)
	numLiteral := l.input[position:l.pos]
	if l.ch == 'm' && l.peekChar() == 's' {
		l.readChar() // consume 'm'
		l.readChar() // consume 's'
		return Token{
			Type:    TOKEN_DURATION,
			Literal: numLiteral + "ms",
			Line:    startLine,
			Column:  startColumn,
		}
	}
	// 有効な単位のみを期間として認識（m/s/h/d）
	if l.ch == 's' || l.ch == 'h' || l.ch == 'd' {
		unit := string(l.ch)
		l.readChar()
		return Token{
			Type:    TOKEN_DURATION,
			Literal: numLiteral + unit,
			Line:    startLine,
			Column:  startColumn,
		}
	}
	// 'm' は 'ms' でない場合のみ分として認識
	if l.ch == 'm' && l.peekChar() != 's' && !isLetter(l.peekChar()) {
		unit := string(l.ch)
		l.readChar()
		return Token{
			Type:    TOKEN_DURATION,
			Literal: numLiteral + unit,
			Line:    startLine,
			Column:  startColumn,
		}
	}

	return Token{
		Type:    TOKEN_INT,
		Literal: numLiteral,
		Line:    startLine,
		Column:  startColumn,
	}
}

func (l *Lexer) readString() string {
	var result []byte
	l.readChar() // 開始の '"' をスキップ

	for {
		if l.ch == '"' || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			case 'u':
				// Unicode escape: \uXXXX
				l.readChar()
				code := l.readHexDigits(4)
				if code >= 0 {
					result = append(result, encodeUTF8(rune(code))...)
				}
				continue
			case 'x':
				// Hex escape: \xXX
				l.readChar()
				code := l.readHexDigits(2)
				if code >= 0 {
					result = append(result, byte(code))
				}
				continue
			default:
				result = append(result, l.ch)
			}
		} else {
			result = append(result, l.ch)
		}
		l.readChar()
	}

	if l.ch == '"' {
		l.readChar() // 終了の '"' をスキップ
	}

	return string(result)
}

// readHexDigits は指定桁数の16進数を読み取る
func (l *Lexer) readHexDigits(count int) int {
	result := 0
	for i := 0; i < count; i++ {
		if !isHexDigit(l.ch) {
			return -1
		}
		result = result*16 + hexValue(l.ch)
		if i < count-1 {
			l.readChar()
		}
	}
	return result
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func hexValue(ch byte) int {
	if ch >= '0' && ch <= '9' {
		return int(ch - '0')
	}
	if ch >= 'a' && ch <= 'f' {
		return int(ch - 'a' + 10)
	}
	if ch >= 'A' && ch <= 'F' {
		return int(ch - 'A' + 10)
	}
	return 0
}

// encodeUTF8 はruneをUTF-8バイト列に変換する
func encodeUTF8(r rune) []byte {
	if r < 0x80 {
		return []byte{byte(r)}
	}
	if r < 0x800 {
		return []byte{
			byte(0xC0 | (r >> 6)),
			byte(0x80 | (r & 0x3F)),
		}
	}
	if r < 0x10000 {
		return []byte{
			byte(0xE0 | (r >> 12)),
			byte(0x80 | ((r >> 6) & 0x3F)),
			byte(0x80 | (r & 0x3F)),
		}
	}
	return []byte{
		byte(0xF0 | (r >> 18)),
		byte(0x80 | ((r >> 12) & 0x3F)),
		byte(0x80 | ((r >> 6) & 0x3F)),
		byte(0x80 | (r & 0x3F)),
	}
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
