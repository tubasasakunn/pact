package parser

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

	for l.ch != '"' && l.ch != 0 {
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
