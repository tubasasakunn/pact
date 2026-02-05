package parser

// TokenType はトークンの種類を表す
type TokenType int

const (
	TOKEN_EOF TokenType = iota
	TOKEN_ILLEGAL

	// 識別子・リテラル
	TOKEN_IDENT
	TOKEN_STRING
	TOKEN_INT
	TOKEN_FLOAT
	TOKEN_DURATION

	// キーワード
	TOKEN_COMPONENT
	TOKEN_IMPORT
	TOKEN_TYPE
	TOKEN_ENUM
	TOKEN_DEPENDS
	TOKEN_ON
	TOKEN_EXTENDS
	TOKEN_IMPLEMENTS
	TOKEN_CONTAINS
	TOKEN_AGGREGATES
	TOKEN_PROVIDES
	TOKEN_REQUIRES
	TOKEN_FLOW
	TOKEN_STATES
	TOKEN_STATE
	TOKEN_PARALLEL
	TOKEN_REGION
	TOKEN_INITIAL
	TOKEN_FINAL
	TOKEN_ENTRY
	TOKEN_EXIT
	TOKEN_IF
	TOKEN_ELSE
	TOKEN_FOR
	TOKEN_IN
	TOKEN_WHILE
	TOKEN_RETURN
	TOKEN_THROW
	TOKEN_AWAIT
	TOKEN_ASYNC
	TOKEN_THROWS
	TOKEN_WHEN
	TOKEN_AFTER
	TOKEN_DO
	TOKEN_TRUE
	TOKEN_FALSE
	TOKEN_NULL
	TOKEN_AS

	// 記号・演算子
	TOKEN_LBRACE   // {
	TOKEN_RBRACE   // }
	TOKEN_LPAREN   // (
	TOKEN_RPAREN   // )
	TOKEN_LBRACKET // [
	TOKEN_RBRACKET // ]
	TOKEN_COLON    // :
	TOKEN_COMMA    // ,
	TOKEN_DOT      // .
	TOKEN_ARROW    // ->
	TOKEN_AT       // @
	TOKEN_QUESTION // ?
	TOKEN_PLUS     // +
	TOKEN_MINUS    // -
	TOKEN_STAR     // *
	TOKEN_SLASH    // /
	TOKEN_PERCENT  // %
	TOKEN_EQ       // ==
	TOKEN_NE       // !=
	TOKEN_LT       // <
	TOKEN_GT       // >
	TOKEN_LE       // <=
	TOKEN_GE       // >=
	TOKEN_AND      // &&
	TOKEN_OR       // ||
	TOKEN_NOT      // !
	TOKEN_ASSIGN   // =
	TOKEN_NULLISH  // ??
	TOKEN_HASH     // #
	TOKEN_TILDE    // ~
)

// Token は字句解析の結果トークンを表す
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Position はソースコード内の位置
type Position struct {
	Line   int
	Column int
}

// DurationValue は期間リテラルの値
type DurationValue struct {
	Value int
	Unit  string
}
