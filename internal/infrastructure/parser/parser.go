package parser

import (
	"pact/internal/domain/ast"
	"pact/internal/domain/errors"
)

// Parser は構文解析器
type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
	errors    []error
}

// NewParser は新しいParserを作成する
func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	// 2つのトークンを読み込んで curToken と peekToken を設定
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseFile はファイル全体をパースする
func (p *Parser) ParseFile() (*ast.SpecFile, error) {
	// TODO: 実装
	return nil, &errors.ParseError{Message: "not implemented"}
}

// Parse はParseFileのエイリアス
func (p *Parser) Parse() (*ast.SpecFile, error) {
	return p.ParseFile()
}

// ParseString は文字列をパースする
func ParseString(input string) (*ast.SpecFile, error) {
	l := NewLexer(input)
	p := NewParser(l)
	return p.ParseFile()
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, &errors.ParseError{
		Pos: ast.Position{
			Line:   p.curToken.Line,
			Column: p.curToken.Column,
		},
		Message: msg,
	})
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}
