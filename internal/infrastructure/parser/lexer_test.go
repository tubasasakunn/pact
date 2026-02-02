package parser

import "testing"

// =============================================================================
// 1.1.1 基本トークン認識
// =============================================================================

// L001: 空入力でEOF
func TestLexer_EOF(t *testing.T) {
	l := NewLexer("")
	tok := l.NextToken()
	if tok.Type != TOKEN_EOF {
		t.Errorf("expected TOKEN_EOF, got %v", tok.Type)
	}
}

// L002: 単純な識別子
func TestLexer_Identifier_Simple(t *testing.T) {
	l := NewLexer("foo")
	tok := l.NextToken()
	if tok.Type != TOKEN_IDENT {
		t.Errorf("expected TOKEN_IDENT, got %v", tok.Type)
	}
	if tok.Literal != "foo" {
		t.Errorf("expected 'foo', got %q", tok.Literal)
	}
}

// L003: アンダースコア含む識別子
func TestLexer_Identifier_WithUnderscore(t *testing.T) {
	l := NewLexer("user_id")
	tok := l.NextToken()
	if tok.Type != TOKEN_IDENT {
		t.Errorf("expected TOKEN_IDENT, got %v", tok.Type)
	}
	if tok.Literal != "user_id" {
		t.Errorf("expected 'user_id', got %q", tok.Literal)
	}
}

// L004: 数字含む識別子
func TestLexer_Identifier_WithNumbers(t *testing.T) {
	l := NewLexer("token2")
	tok := l.NextToken()
	if tok.Type != TOKEN_IDENT {
		t.Errorf("expected TOKEN_IDENT, got %v", tok.Type)
	}
	if tok.Literal != "token2" {
		t.Errorf("expected 'token2', got %q", tok.Literal)
	}
}

// L005: 数字始まりは識別子でない
func TestLexer_Identifier_StartWithNumber(t *testing.T) {
	l := NewLexer("2token")
	tok1 := l.NextToken()
	if tok1.Type != TOKEN_INT {
		t.Errorf("expected TOKEN_INT, got %v", tok1.Type)
	}
	tok2 := l.NextToken()
	if tok2.Type != TOKEN_IDENT {
		t.Errorf("expected TOKEN_IDENT, got %v", tok2.Type)
	}
}

// =============================================================================
// 1.1.2 キーワード認識
// =============================================================================

func TestLexer_Keyword_Component(t *testing.T) {
	l := NewLexer("component")
	tok := l.NextToken()
	if tok.Type != TOKEN_COMPONENT {
		t.Errorf("expected TOKEN_COMPONENT, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Import(t *testing.T) {
	l := NewLexer("import")
	tok := l.NextToken()
	if tok.Type != TOKEN_IMPORT {
		t.Errorf("expected TOKEN_IMPORT, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Type(t *testing.T) {
	l := NewLexer("type")
	tok := l.NextToken()
	if tok.Type != TOKEN_TYPE {
		t.Errorf("expected TOKEN_TYPE, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Enum(t *testing.T) {
	l := NewLexer("enum")
	tok := l.NextToken()
	if tok.Type != TOKEN_ENUM {
		t.Errorf("expected TOKEN_ENUM, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Depends(t *testing.T) {
	l := NewLexer("depends")
	tok := l.NextToken()
	if tok.Type != TOKEN_DEPENDS {
		t.Errorf("expected TOKEN_DEPENDS, got %v", tok.Type)
	}
}

func TestLexer_Keyword_On(t *testing.T) {
	l := NewLexer("on")
	tok := l.NextToken()
	if tok.Type != TOKEN_ON {
		t.Errorf("expected TOKEN_ON, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Extends(t *testing.T) {
	l := NewLexer("extends")
	tok := l.NextToken()
	if tok.Type != TOKEN_EXTENDS {
		t.Errorf("expected TOKEN_EXTENDS, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Implements(t *testing.T) {
	l := NewLexer("implements")
	tok := l.NextToken()
	if tok.Type != TOKEN_IMPLEMENTS {
		t.Errorf("expected TOKEN_IMPLEMENTS, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Contains(t *testing.T) {
	l := NewLexer("contains")
	tok := l.NextToken()
	if tok.Type != TOKEN_CONTAINS {
		t.Errorf("expected TOKEN_CONTAINS, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Aggregates(t *testing.T) {
	l := NewLexer("aggregates")
	tok := l.NextToken()
	if tok.Type != TOKEN_AGGREGATES {
		t.Errorf("expected TOKEN_AGGREGATES, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Provides(t *testing.T) {
	l := NewLexer("provides")
	tok := l.NextToken()
	if tok.Type != TOKEN_PROVIDES {
		t.Errorf("expected TOKEN_PROVIDES, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Requires(t *testing.T) {
	l := NewLexer("requires")
	tok := l.NextToken()
	if tok.Type != TOKEN_REQUIRES {
		t.Errorf("expected TOKEN_REQUIRES, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Flow(t *testing.T) {
	l := NewLexer("flow")
	tok := l.NextToken()
	if tok.Type != TOKEN_FLOW {
		t.Errorf("expected TOKEN_FLOW, got %v", tok.Type)
	}
}

func TestLexer_Keyword_States(t *testing.T) {
	l := NewLexer("states")
	tok := l.NextToken()
	if tok.Type != TOKEN_STATES {
		t.Errorf("expected TOKEN_STATES, got %v", tok.Type)
	}
}

func TestLexer_Keyword_State(t *testing.T) {
	l := NewLexer("state")
	tok := l.NextToken()
	if tok.Type != TOKEN_STATE {
		t.Errorf("expected TOKEN_STATE, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Parallel(t *testing.T) {
	l := NewLexer("parallel")
	tok := l.NextToken()
	if tok.Type != TOKEN_PARALLEL {
		t.Errorf("expected TOKEN_PARALLEL, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Region(t *testing.T) {
	l := NewLexer("region")
	tok := l.NextToken()
	if tok.Type != TOKEN_REGION {
		t.Errorf("expected TOKEN_REGION, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Initial(t *testing.T) {
	l := NewLexer("initial")
	tok := l.NextToken()
	if tok.Type != TOKEN_INITIAL {
		t.Errorf("expected TOKEN_INITIAL, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Final(t *testing.T) {
	l := NewLexer("final")
	tok := l.NextToken()
	if tok.Type != TOKEN_FINAL {
		t.Errorf("expected TOKEN_FINAL, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Entry(t *testing.T) {
	l := NewLexer("entry")
	tok := l.NextToken()
	if tok.Type != TOKEN_ENTRY {
		t.Errorf("expected TOKEN_ENTRY, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Exit(t *testing.T) {
	l := NewLexer("exit")
	tok := l.NextToken()
	if tok.Type != TOKEN_EXIT {
		t.Errorf("expected TOKEN_EXIT, got %v", tok.Type)
	}
}

func TestLexer_Keyword_If(t *testing.T) {
	l := NewLexer("if")
	tok := l.NextToken()
	if tok.Type != TOKEN_IF {
		t.Errorf("expected TOKEN_IF, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Else(t *testing.T) {
	l := NewLexer("else")
	tok := l.NextToken()
	if tok.Type != TOKEN_ELSE {
		t.Errorf("expected TOKEN_ELSE, got %v", tok.Type)
	}
}

func TestLexer_Keyword_For(t *testing.T) {
	l := NewLexer("for")
	tok := l.NextToken()
	if tok.Type != TOKEN_FOR {
		t.Errorf("expected TOKEN_FOR, got %v", tok.Type)
	}
}

func TestLexer_Keyword_In(t *testing.T) {
	l := NewLexer("in")
	tok := l.NextToken()
	if tok.Type != TOKEN_IN {
		t.Errorf("expected TOKEN_IN, got %v", tok.Type)
	}
}

func TestLexer_Keyword_While(t *testing.T) {
	l := NewLexer("while")
	tok := l.NextToken()
	if tok.Type != TOKEN_WHILE {
		t.Errorf("expected TOKEN_WHILE, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Return(t *testing.T) {
	l := NewLexer("return")
	tok := l.NextToken()
	if tok.Type != TOKEN_RETURN {
		t.Errorf("expected TOKEN_RETURN, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Throw(t *testing.T) {
	l := NewLexer("throw")
	tok := l.NextToken()
	if tok.Type != TOKEN_THROW {
		t.Errorf("expected TOKEN_THROW, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Await(t *testing.T) {
	l := NewLexer("await")
	tok := l.NextToken()
	if tok.Type != TOKEN_AWAIT {
		t.Errorf("expected TOKEN_AWAIT, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Async(t *testing.T) {
	l := NewLexer("async")
	tok := l.NextToken()
	if tok.Type != TOKEN_ASYNC {
		t.Errorf("expected TOKEN_ASYNC, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Throws(t *testing.T) {
	l := NewLexer("throws")
	tok := l.NextToken()
	if tok.Type != TOKEN_THROWS {
		t.Errorf("expected TOKEN_THROWS, got %v", tok.Type)
	}
}

func TestLexer_Keyword_When(t *testing.T) {
	l := NewLexer("when")
	tok := l.NextToken()
	if tok.Type != TOKEN_WHEN {
		t.Errorf("expected TOKEN_WHEN, got %v", tok.Type)
	}
}

func TestLexer_Keyword_After(t *testing.T) {
	l := NewLexer("after")
	tok := l.NextToken()
	if tok.Type != TOKEN_AFTER {
		t.Errorf("expected TOKEN_AFTER, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Do(t *testing.T) {
	l := NewLexer("do")
	tok := l.NextToken()
	if tok.Type != TOKEN_DO {
		t.Errorf("expected TOKEN_DO, got %v", tok.Type)
	}
}

func TestLexer_Keyword_True(t *testing.T) {
	l := NewLexer("true")
	tok := l.NextToken()
	if tok.Type != TOKEN_TRUE {
		t.Errorf("expected TOKEN_TRUE, got %v", tok.Type)
	}
}

func TestLexer_Keyword_False(t *testing.T) {
	l := NewLexer("false")
	tok := l.NextToken()
	if tok.Type != TOKEN_FALSE {
		t.Errorf("expected TOKEN_FALSE, got %v", tok.Type)
	}
}

func TestLexer_Keyword_Null(t *testing.T) {
	l := NewLexer("null")
	tok := l.NextToken()
	if tok.Type != TOKEN_NULL {
		t.Errorf("expected TOKEN_NULL, got %v", tok.Type)
	}
}

func TestLexer_Keyword_As(t *testing.T) {
	l := NewLexer("as")
	tok := l.NextToken()
	if tok.Type != TOKEN_AS {
		t.Errorf("expected TOKEN_AS, got %v", tok.Type)
	}
}

// =============================================================================
// 1.1.3 リテラル
// =============================================================================

// L050: 単純な文字列
func TestLexer_String_Simple(t *testing.T) {
	l := NewLexer(`"hello"`)
	tok := l.NextToken()
	if tok.Type != TOKEN_STRING {
		t.Errorf("expected TOKEN_STRING, got %v", tok.Type)
	}
	if tok.Literal != "hello" {
		t.Errorf("expected 'hello', got %q", tok.Literal)
	}
}

// L051: 空文字列
func TestLexer_String_Empty(t *testing.T) {
	l := NewLexer(`""`)
	tok := l.NextToken()
	if tok.Type != TOKEN_STRING {
		t.Errorf("expected TOKEN_STRING, got %v", tok.Type)
	}
	if tok.Literal != "" {
		t.Errorf("expected empty string, got %q", tok.Literal)
	}
}

// L052: エスケープ（引用符）
func TestLexer_String_Escape_Quote(t *testing.T) {
	l := NewLexer(`"say \"hi\""`)
	tok := l.NextToken()
	if tok.Type != TOKEN_STRING {
		t.Errorf("expected TOKEN_STRING, got %v", tok.Type)
	}
	if tok.Literal != `say "hi"` {
		t.Errorf("expected 'say \"hi\"', got %q", tok.Literal)
	}
}

// L053: エスケープ（バックスラッシュ）
func TestLexer_String_Escape_Backslash(t *testing.T) {
	l := NewLexer(`"path\\to"`)
	tok := l.NextToken()
	if tok.Type != TOKEN_STRING {
		t.Errorf("expected TOKEN_STRING, got %v", tok.Type)
	}
	if tok.Literal != `path\to` {
		t.Errorf("expected 'path\\to', got %q", tok.Literal)
	}
}

// L054: エスケープ（改行）
func TestLexer_String_Escape_Newline(t *testing.T) {
	l := NewLexer(`"line1\nline2"`)
	tok := l.NextToken()
	if tok.Type != TOKEN_STRING {
		t.Errorf("expected TOKEN_STRING, got %v", tok.Type)
	}
	if tok.Literal != "line1\nline2" {
		t.Errorf("expected 'line1\\nline2', got %q", tok.Literal)
	}
}

// L055: エスケープ（タブ）
func TestLexer_String_Escape_Tab(t *testing.T) {
	l := NewLexer(`"col1\tcol2"`)
	tok := l.NextToken()
	if tok.Type != TOKEN_STRING {
		t.Errorf("expected TOKEN_STRING, got %v", tok.Type)
	}
	if tok.Literal != "col1\tcol2" {
		t.Errorf("expected 'col1\\tcol2', got %q", tok.Literal)
	}
}

// L056: 終端なし文字列（EOFで終了）
func TestLexer_String_Unterminated(t *testing.T) {
	l := NewLexer(`"hello`)
	tok := l.NextToken()
	if tok.Type != TOKEN_STRING {
		t.Errorf("expected TOKEN_STRING, got %v", tok.Type)
	}
	if tok.Literal != "hello" {
		t.Errorf("expected 'hello', got %q", tok.Literal)
	}
	tok2 := l.NextToken()
	if tok2.Type != TOKEN_EOF {
		t.Errorf("expected TOKEN_EOF, got %v", tok2.Type)
	}
}

// L060: ゼロ
func TestLexer_Int_Zero(t *testing.T) {
	l := NewLexer("0")
	tok := l.NextToken()
	if tok.Type != TOKEN_INT {
		t.Errorf("expected TOKEN_INT, got %v", tok.Type)
	}
	if tok.Literal != "0" {
		t.Errorf("expected '0', got %q", tok.Literal)
	}
}

// L061: 正の整数
func TestLexer_Int_Positive(t *testing.T) {
	l := NewLexer("42")
	tok := l.NextToken()
	if tok.Type != TOKEN_INT {
		t.Errorf("expected TOKEN_INT, got %v", tok.Type)
	}
	if tok.Literal != "42" {
		t.Errorf("expected '42', got %q", tok.Literal)
	}
}

// L062: 大きな整数
func TestLexer_Int_Large(t *testing.T) {
	l := NewLexer("1234567890")
	tok := l.NextToken()
	if tok.Type != TOKEN_INT {
		t.Errorf("expected TOKEN_INT, got %v", tok.Type)
	}
	if tok.Literal != "1234567890" {
		t.Errorf("expected '1234567890', got %q", tok.Literal)
	}
}

// L065: 単純な浮動小数点
func TestLexer_Float_Simple(t *testing.T) {
	l := NewLexer("3.14")
	tok := l.NextToken()
	if tok.Type != TOKEN_FLOAT {
		t.Errorf("expected TOKEN_FLOAT, got %v", tok.Type)
	}
	if tok.Literal != "3.14" {
		t.Errorf("expected '3.14', got %q", tok.Literal)
	}
}

// L066: 先頭ゼロ
func TestLexer_Float_LeadingZero(t *testing.T) {
	l := NewLexer("0.5")
	tok := l.NextToken()
	if tok.Type != TOKEN_FLOAT {
		t.Errorf("expected TOKEN_FLOAT, got %v", tok.Type)
	}
	if tok.Literal != "0.5" {
		t.Errorf("expected '0.5', got %q", tok.Literal)
	}
}

// L067: 末尾ドット（floatではない）
func TestLexer_Float_TrailingDot(t *testing.T) {
	l := NewLexer("1.")
	tok1 := l.NextToken()
	if tok1.Type != TOKEN_INT {
		t.Errorf("expected TOKEN_INT, got %v", tok1.Type)
	}
	tok2 := l.NextToken()
	if tok2.Type != TOKEN_DOT {
		t.Errorf("expected TOKEN_DOT, got %v", tok2.Type)
	}
}

// L068: 先頭ドット（floatではない）
func TestLexer_Float_LeadingDot(t *testing.T) {
	l := NewLexer(".5")
	tok1 := l.NextToken()
	if tok1.Type != TOKEN_DOT {
		t.Errorf("expected TOKEN_DOT, got %v", tok1.Type)
	}
	tok2 := l.NextToken()
	if tok2.Type != TOKEN_INT {
		t.Errorf("expected TOKEN_INT, got %v", tok2.Type)
	}
}

// L070: ミリ秒
func TestLexer_Duration_Milliseconds(t *testing.T) {
	l := NewLexer("500ms")
	tok := l.NextToken()
	if tok.Type != TOKEN_DURATION {
		t.Errorf("expected TOKEN_DURATION, got %v", tok.Type)
	}
	if tok.Literal != "500ms" {
		t.Errorf("expected '500ms', got %q", tok.Literal)
	}
}

// L071: 秒
func TestLexer_Duration_Seconds(t *testing.T) {
	l := NewLexer("30s")
	tok := l.NextToken()
	if tok.Type != TOKEN_DURATION {
		t.Errorf("expected TOKEN_DURATION, got %v", tok.Type)
	}
	if tok.Literal != "30s" {
		t.Errorf("expected '30s', got %q", tok.Literal)
	}
}

// L072: 分
func TestLexer_Duration_Minutes(t *testing.T) {
	l := NewLexer("5m")
	tok := l.NextToken()
	if tok.Type != TOKEN_DURATION {
		t.Errorf("expected TOKEN_DURATION, got %v", tok.Type)
	}
	if tok.Literal != "5m" {
		t.Errorf("expected '5m', got %q", tok.Literal)
	}
}

// L073: 時間
func TestLexer_Duration_Hours(t *testing.T) {
	l := NewLexer("24h")
	tok := l.NextToken()
	if tok.Type != TOKEN_DURATION {
		t.Errorf("expected TOKEN_DURATION, got %v", tok.Type)
	}
	if tok.Literal != "24h" {
		t.Errorf("expected '24h', got %q", tok.Literal)
	}
}

// L074: 日
func TestLexer_Duration_Days(t *testing.T) {
	l := NewLexer("7d")
	tok := l.NextToken()
	if tok.Type != TOKEN_DURATION {
		t.Errorf("expected TOKEN_DURATION, got %v", tok.Type)
	}
	if tok.Literal != "7d" {
		t.Errorf("expected '7d', got %q", tok.Literal)
	}
}

// =============================================================================
// 1.1.4 演算子・記号
// =============================================================================

func TestLexer_Symbol_LBrace(t *testing.T) {
	l := NewLexer("{")
	tok := l.NextToken()
	if tok.Type != TOKEN_LBRACE {
		t.Errorf("expected TOKEN_LBRACE, got %v", tok.Type)
	}
}

func TestLexer_Symbol_RBrace(t *testing.T) {
	l := NewLexer("}")
	tok := l.NextToken()
	if tok.Type != TOKEN_RBRACE {
		t.Errorf("expected TOKEN_RBRACE, got %v", tok.Type)
	}
}

func TestLexer_Symbol_LParen(t *testing.T) {
	l := NewLexer("(")
	tok := l.NextToken()
	if tok.Type != TOKEN_LPAREN {
		t.Errorf("expected TOKEN_LPAREN, got %v", tok.Type)
	}
}

func TestLexer_Symbol_RParen(t *testing.T) {
	l := NewLexer(")")
	tok := l.NextToken()
	if tok.Type != TOKEN_RPAREN {
		t.Errorf("expected TOKEN_RPAREN, got %v", tok.Type)
	}
}

func TestLexer_Symbol_LBracket(t *testing.T) {
	l := NewLexer("[")
	tok := l.NextToken()
	if tok.Type != TOKEN_LBRACKET {
		t.Errorf("expected TOKEN_LBRACKET, got %v", tok.Type)
	}
}

func TestLexer_Symbol_RBracket(t *testing.T) {
	l := NewLexer("]")
	tok := l.NextToken()
	if tok.Type != TOKEN_RBRACKET {
		t.Errorf("expected TOKEN_RBRACKET, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Colon(t *testing.T) {
	l := NewLexer(":")
	tok := l.NextToken()
	if tok.Type != TOKEN_COLON {
		t.Errorf("expected TOKEN_COLON, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Comma(t *testing.T) {
	l := NewLexer(",")
	tok := l.NextToken()
	if tok.Type != TOKEN_COMMA {
		t.Errorf("expected TOKEN_COMMA, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Dot(t *testing.T) {
	l := NewLexer(".")
	tok := l.NextToken()
	if tok.Type != TOKEN_DOT {
		t.Errorf("expected TOKEN_DOT, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Arrow(t *testing.T) {
	l := NewLexer("->")
	tok := l.NextToken()
	if tok.Type != TOKEN_ARROW {
		t.Errorf("expected TOKEN_ARROW, got %v", tok.Type)
	}
}

func TestLexer_Symbol_At(t *testing.T) {
	l := NewLexer("@")
	tok := l.NextToken()
	if tok.Type != TOKEN_AT {
		t.Errorf("expected TOKEN_AT, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Question(t *testing.T) {
	l := NewLexer("?")
	tok := l.NextToken()
	if tok.Type != TOKEN_QUESTION {
		t.Errorf("expected TOKEN_QUESTION, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Plus(t *testing.T) {
	l := NewLexer("+")
	tok := l.NextToken()
	if tok.Type != TOKEN_PLUS {
		t.Errorf("expected TOKEN_PLUS, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Minus(t *testing.T) {
	l := NewLexer("-")
	tok := l.NextToken()
	if tok.Type != TOKEN_MINUS {
		t.Errorf("expected TOKEN_MINUS, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Star(t *testing.T) {
	l := NewLexer("*")
	tok := l.NextToken()
	if tok.Type != TOKEN_STAR {
		t.Errorf("expected TOKEN_STAR, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Slash(t *testing.T) {
	l := NewLexer("/")
	tok := l.NextToken()
	if tok.Type != TOKEN_SLASH {
		t.Errorf("expected TOKEN_SLASH, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Percent(t *testing.T) {
	l := NewLexer("%")
	tok := l.NextToken()
	if tok.Type != TOKEN_PERCENT {
		t.Errorf("expected TOKEN_PERCENT, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Eq(t *testing.T) {
	l := NewLexer("==")
	tok := l.NextToken()
	if tok.Type != TOKEN_EQ {
		t.Errorf("expected TOKEN_EQ, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Ne(t *testing.T) {
	l := NewLexer("!=")
	tok := l.NextToken()
	if tok.Type != TOKEN_NE {
		t.Errorf("expected TOKEN_NE, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Lt(t *testing.T) {
	l := NewLexer("<")
	tok := l.NextToken()
	if tok.Type != TOKEN_LT {
		t.Errorf("expected TOKEN_LT, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Gt(t *testing.T) {
	l := NewLexer(">")
	tok := l.NextToken()
	if tok.Type != TOKEN_GT {
		t.Errorf("expected TOKEN_GT, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Le(t *testing.T) {
	l := NewLexer("<=")
	tok := l.NextToken()
	if tok.Type != TOKEN_LE {
		t.Errorf("expected TOKEN_LE, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Ge(t *testing.T) {
	l := NewLexer(">=")
	tok := l.NextToken()
	if tok.Type != TOKEN_GE {
		t.Errorf("expected TOKEN_GE, got %v", tok.Type)
	}
}

func TestLexer_Symbol_And(t *testing.T) {
	l := NewLexer("&&")
	tok := l.NextToken()
	if tok.Type != TOKEN_AND {
		t.Errorf("expected TOKEN_AND, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Or(t *testing.T) {
	l := NewLexer("||")
	tok := l.NextToken()
	if tok.Type != TOKEN_OR {
		t.Errorf("expected TOKEN_OR, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Not(t *testing.T) {
	l := NewLexer("!")
	tok := l.NextToken()
	if tok.Type != TOKEN_NOT {
		t.Errorf("expected TOKEN_NOT, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Assign(t *testing.T) {
	l := NewLexer("=")
	tok := l.NextToken()
	if tok.Type != TOKEN_ASSIGN {
		t.Errorf("expected TOKEN_ASSIGN, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Nullish(t *testing.T) {
	l := NewLexer("??")
	tok := l.NextToken()
	if tok.Type != TOKEN_NULLISH {
		t.Errorf("expected TOKEN_NULLISH, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Hash(t *testing.T) {
	l := NewLexer("#")
	tok := l.NextToken()
	if tok.Type != TOKEN_HASH {
		t.Errorf("expected TOKEN_HASH, got %v", tok.Type)
	}
}

func TestLexer_Symbol_Tilde(t *testing.T) {
	l := NewLexer("~")
	tok := l.NextToken()
	if tok.Type != TOKEN_TILDE {
		t.Errorf("expected TOKEN_TILDE, got %v", tok.Type)
	}
}

// =============================================================================
// 1.1.5 コメント
// =============================================================================

// L120: 行コメントスキップ
func TestLexer_Comment_Line(t *testing.T) {
	l := NewLexer("// comment\nfoo")
	tok := l.NextToken()
	if tok.Type != TOKEN_IDENT {
		t.Errorf("expected TOKEN_IDENT, got %v", tok.Type)
	}
	if tok.Literal != "foo" {
		t.Errorf("expected 'foo', got %q", tok.Literal)
	}
}

// L121: 行コメントで終了
func TestLexer_Comment_Line_EOF(t *testing.T) {
	l := NewLexer("// comment")
	tok := l.NextToken()
	if tok.Type != TOKEN_EOF {
		t.Errorf("expected TOKEN_EOF, got %v", tok.Type)
	}
}

// L122: ブロックコメントスキップ
func TestLexer_Comment_Block(t *testing.T) {
	l := NewLexer("/* comment */foo")
	tok := l.NextToken()
	if tok.Type != TOKEN_IDENT {
		t.Errorf("expected TOKEN_IDENT, got %v", tok.Type)
	}
	if tok.Literal != "foo" {
		t.Errorf("expected 'foo', got %q", tok.Literal)
	}
}

// L123: 複数行ブロックコメント
func TestLexer_Comment_Block_Multiline(t *testing.T) {
	l := NewLexer("/* line1\nline2 */foo")
	tok := l.NextToken()
	if tok.Type != TOKEN_IDENT {
		t.Errorf("expected TOKEN_IDENT, got %v", tok.Type)
	}
	if tok.Literal != "foo" {
		t.Errorf("expected 'foo', got %q", tok.Literal)
	}
}

// L124: ネスト非対応（最初の*/で終了）
func TestLexer_Comment_Block_NoNesting(t *testing.T) {
	l := NewLexer("/* /* inner */ outer")
	tok := l.NextToken()
	if tok.Type != TOKEN_IDENT {
		t.Errorf("expected TOKEN_IDENT, got %v", tok.Type)
	}
	if tok.Literal != "outer" {
		t.Errorf("expected 'outer', got %q", tok.Literal)
	}
}

// L125: 終端なしブロックコメント（EOFで終了）
func TestLexer_Comment_Block_Unterminated(t *testing.T) {
	l := NewLexer("/* comment")
	tok := l.NextToken()
	if tok.Type != TOKEN_EOF {
		t.Errorf("expected TOKEN_EOF, got %v", tok.Type)
	}
}

// =============================================================================
// 1.1.6 空白・位置情報
// =============================================================================

// L130: スペーススキップ
func TestLexer_Whitespace_Spaces(t *testing.T) {
	l := NewLexer("  foo  ")
	tok := l.NextToken()
	if tok.Type != TOKEN_IDENT {
		t.Errorf("expected TOKEN_IDENT, got %v", tok.Type)
	}
	if tok.Literal != "foo" {
		t.Errorf("expected 'foo', got %q", tok.Literal)
	}
}

// L131: タブスキップ
func TestLexer_Whitespace_Tabs(t *testing.T) {
	l := NewLexer("\t\tfoo\t")
	tok := l.NextToken()
	if tok.Type != TOKEN_IDENT {
		t.Errorf("expected TOKEN_IDENT, got %v", tok.Type)
	}
	if tok.Literal != "foo" {
		t.Errorf("expected 'foo', got %q", tok.Literal)
	}
}

// L132: 改行スキップ
func TestLexer_Whitespace_Newlines(t *testing.T) {
	l := NewLexer("\n\nfoo\n")
	tok := l.NextToken()
	if tok.Type != TOKEN_IDENT {
		t.Errorf("expected TOKEN_IDENT, got %v", tok.Type)
	}
	if tok.Literal != "foo" {
		t.Errorf("expected 'foo', got %q", tok.Literal)
	}
}

// L133: 混合空白スキップ
func TestLexer_Whitespace_Mixed(t *testing.T) {
	l := NewLexer("  \t\n  foo")
	tok := l.NextToken()
	if tok.Type != TOKEN_IDENT {
		t.Errorf("expected TOKEN_IDENT, got %v", tok.Type)
	}
	if tok.Literal != "foo" {
		t.Errorf("expected 'foo', got %q", tok.Literal)
	}
}

// L140: 最初のトークン位置
func TestLexer_Position_FirstToken(t *testing.T) {
	l := NewLexer("foo")
	tok := l.NextToken()
	if tok.Line != 1 {
		t.Errorf("expected line 1, got %d", tok.Line)
	}
	if tok.Column != 1 {
		t.Errorf("expected column 1, got %d", tok.Column)
	}
}

// L141: 空白後の位置
func TestLexer_Position_AfterWhitespace(t *testing.T) {
	l := NewLexer("  foo")
	tok := l.NextToken()
	if tok.Line != 1 {
		t.Errorf("expected line 1, got %d", tok.Line)
	}
	if tok.Column != 3 {
		t.Errorf("expected column 3, got %d", tok.Column)
	}
}

// L142: 改行後の位置
func TestLexer_Position_AfterNewline(t *testing.T) {
	l := NewLexer("\nfoo")
	tok := l.NextToken()
	if tok.Line != 2 {
		t.Errorf("expected line 2, got %d", tok.Line)
	}
	if tok.Column != 1 {
		t.Errorf("expected column 1, got %d", tok.Column)
	}
}

// L143: 複数トークンの位置
func TestLexer_Position_MultipleTokens(t *testing.T) {
	l := NewLexer("foo bar")
	tok1 := l.NextToken()
	if tok1.Line != 1 || tok1.Column != 1 {
		t.Errorf("expected pos {1,1}, got {%d,%d}", tok1.Line, tok1.Column)
	}
	tok2 := l.NextToken()
	if tok2.Line != 1 || tok2.Column != 5 {
		t.Errorf("expected pos {1,5}, got {%d,%d}", tok2.Line, tok2.Column)
	}
}

// L144: 複数行の位置
func TestLexer_Position_MultipleLines(t *testing.T) {
	l := NewLexer("foo\nbar\nbaz")
	tok1 := l.NextToken()
	if tok1.Line != 1 || tok1.Column != 1 {
		t.Errorf("expected pos {1,1}, got {%d,%d}", tok1.Line, tok1.Column)
	}
	tok2 := l.NextToken()
	if tok2.Line != 2 || tok2.Column != 1 {
		t.Errorf("expected pos {2,1}, got {%d,%d}", tok2.Line, tok2.Column)
	}
	tok3 := l.NextToken()
	if tok3.Line != 3 || tok3.Column != 1 {
		t.Errorf("expected pos {3,1}, got {%d,%d}", tok3.Line, tok3.Column)
	}
}

// =============================================================================
// 1.1.7 複合トークン列
// =============================================================================

// L150: 型宣言
func TestLexer_Sequence_TypeDecl(t *testing.T) {
	l := NewLexer("type Foo { }")
	expected := []TokenType{TOKEN_TYPE, TOKEN_IDENT, TOKEN_LBRACE, TOKEN_RBRACE, TOKEN_EOF}
	for _, exp := range expected {
		tok := l.NextToken()
		if tok.Type != exp {
			t.Errorf("expected %v, got %v", exp, tok.Type)
		}
	}
}

// L151: メソッド宣言
func TestLexer_Sequence_MethodDecl(t *testing.T) {
	l := NewLexer("Login(email: string) -> Token")
	expected := []TokenType{
		TOKEN_IDENT,    // Login
		TOKEN_LPAREN,   // (
		TOKEN_IDENT,    // email
		TOKEN_COLON,    // :
		TOKEN_IDENT,    // string
		TOKEN_RPAREN,   // )
		TOKEN_ARROW,    // ->
		TOKEN_IDENT,    // Token
		TOKEN_EOF,
	}
	for _, exp := range expected {
		tok := l.NextToken()
		if tok.Type != exp {
			t.Errorf("expected %v, got %v", exp, tok.Type)
		}
	}
}

// L152: アノテーション
func TestLexer_Sequence_Annotation(t *testing.T) {
	l := NewLexer(`@description("hello")`)
	expected := []TokenType{TOKEN_AT, TOKEN_IDENT, TOKEN_LPAREN, TOKEN_STRING, TOKEN_RPAREN, TOKEN_EOF}
	for _, exp := range expected {
		tok := l.NextToken()
		if tok.Type != exp {
			t.Errorf("expected %v, got %v", exp, tok.Type)
		}
	}
}

// L153: 状態遷移
func TestLexer_Sequence_Transition(t *testing.T) {
	l := NewLexer("Idle -> Active on Start")
	expected := []TokenType{
		TOKEN_IDENT, // Idle
		TOKEN_ARROW, // ->
		TOKEN_IDENT, // Active
		TOKEN_ON,    // on
		TOKEN_IDENT, // Start
		TOKEN_EOF,
	}
	for _, exp := range expected {
		tok := l.NextToken()
		if tok.Type != exp {
			t.Errorf("expected %v, got %v", exp, tok.Type)
		}
	}
}
