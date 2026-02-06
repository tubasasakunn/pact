package parser

import (
	"strconv"

	"pact/internal/domain/ast"
	"pact/internal/domain/errors"
)

// Parser は構文解析器
type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
	errors    []error
	maxErrors int // 最大エラー数（デフォルト10）
}

// NewParser は新しいParserを作成する
func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l, maxErrors: 10}
	p.nextToken()
	p.nextToken()
	return p
}

// addError はエラーを追加し、最大エラー数を超えたかを返す
func (p *Parser) addError(err error) bool {
	p.errors = append(p.errors, err)
	return len(p.errors) >= p.maxErrors
}

// synchronize はエラー回復のため、次の同期ポイントまで進む
func (p *Parser) synchronize() {
	for p.curToken.Type != TOKEN_EOF {
		// 宣言の開始トークンに到達したら停止
		switch p.curToken.Type {
		case TOKEN_COMPONENT, TOKEN_IMPORT, TOKEN_TYPE, TOKEN_ENUM,
			TOKEN_FLOW, TOKEN_STATES, TOKEN_PROVIDES, TOKEN_REQUIRES,
			TOKEN_DEPENDS, TOKEN_EXTENDS, TOKEN_IMPLEMENTS, TOKEN_CONTAINS, TOKEN_AGGREGATES:
			return
		}
		// ブロック終端に到達したら次に進む
		if p.curToken.Type == TOKEN_RBRACE {
			p.nextToken()
			return
		}
		p.nextToken()
	}
}

// getErrors は収集したエラーを返す
func (p *Parser) getErrors() error {
	if len(p.errors) == 0 {
		return nil
	}
	multiErr := &errors.MultiError{}
	for _, err := range p.errors {
		multiErr.Add(err)
	}
	return multiErr
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Parse はParseFileのエイリアス
func (p *Parser) Parse() (*ast.SpecFile, error) {
	return p.ParseFile()
}

// ParseFile はファイル全体をパースする
func (p *Parser) ParseFile() (*ast.SpecFile, error) {
	spec := &ast.SpecFile{}

	for p.curToken.Type != TOKEN_EOF {
		switch p.curToken.Type {
		case TOKEN_IMPORT:
			imp, err := p.parseImport()
			if err != nil {
				if p.addError(err) {
					return spec, p.getErrors()
				}
				p.synchronize()
				continue
			}
			spec.Imports = append(spec.Imports, *imp)

		case TOKEN_AT:
			// アノテーション付きの宣言
			annotations, err := p.parseAnnotations()
			if err != nil {
				if p.addError(err) {
					return spec, p.getErrors()
				}
				p.synchronize()
				continue
			}
			if err := p.parseDeclarationWithAnnotations(spec, annotations); err != nil {
				if p.addError(err) {
					return spec, p.getErrors()
				}
				p.synchronize()
				continue
			}

		case TOKEN_COMPONENT:
			comp, err := p.parseComponent(nil)
			if err != nil {
				if p.addError(err) {
					return spec, p.getErrors()
				}
				p.synchronize()
				continue
			}
			spec.Component = comp
			spec.Components = append(spec.Components, *comp)

		default:
			if p.addError(p.newErrorWithSuggestion("unexpected token: %v", p.curToken.Literal,
				"import", "component", "@annotation")) {
				return spec, p.getErrors()
			}
			p.synchronize()
		}
	}

	// エラーがあれば返す
	if err := p.getErrors(); err != nil {
		return spec, err
	}

	return spec, nil
}

// ParseString は文字列をパースする
func ParseString(input string) (*ast.SpecFile, error) {
	l := NewLexer(input)
	p := NewParser(l)
	return p.ParseFile()
}

func (p *Parser) parseDeclarationWithAnnotations(spec *ast.SpecFile, annotations []ast.AnnotationDecl) error {
	switch p.curToken.Type {
	case TOKEN_COMPONENT:
		comp, err := p.parseComponent(annotations)
		if err != nil {
			return err
		}
		spec.Component = comp
		spec.Components = append(spec.Components, *comp)
		return nil
	default:
		return p.newError("expected declaration after annotations")
	}
}

// =============================================================================
// Import
// =============================================================================

func (p *Parser) parseImport() (*ast.ImportDecl, error) {
	imp := &ast.ImportDecl{
		Pos: p.curPos(),
	}

	p.nextToken() // consume 'import'

	if p.curToken.Type != TOKEN_STRING {
		return nil, p.newError("expected string after 'import'")
	}
	imp.Path = p.curToken.Literal
	p.nextToken()

	// as Alias
	if p.curToken.Type == TOKEN_AS {
		p.nextToken()
		if p.curToken.Type != TOKEN_IDENT {
			return nil, p.newError("expected identifier after 'as'")
		}
		alias := p.curToken.Literal
		imp.Alias = &alias
		p.nextToken()
	}

	return imp, nil
}

// =============================================================================
// Helpers
// =============================================================================

// isReservedKeyword は現在のトークンが予約語かどうかを返す
func (p *Parser) isReservedKeyword() bool {
	switch p.curToken.Type {
	case TOKEN_COMPONENT, TOKEN_IMPORT, TOKEN_TYPE, TOKEN_ENUM,
		TOKEN_DEPENDS, TOKEN_ON, TOKEN_EXTENDS, TOKEN_IMPLEMENTS,
		TOKEN_CONTAINS, TOKEN_AGGREGATES, TOKEN_PROVIDES, TOKEN_REQUIRES,
		TOKEN_FLOW, TOKEN_STATES, TOKEN_STATE, TOKEN_PARALLEL, TOKEN_REGION,
		TOKEN_INITIAL, TOKEN_FINAL, TOKEN_ENTRY, TOKEN_EXIT,
		TOKEN_IF, TOKEN_ELSE, TOKEN_FOR, TOKEN_IN, TOKEN_WHILE,
		TOKEN_RETURN, TOKEN_THROW, TOKEN_AWAIT, TOKEN_ASYNC, TOKEN_THROWS,
		TOKEN_WHEN, TOKEN_AFTER, TOKEN_DO, TOKEN_TRUE, TOKEN_FALSE, TOKEN_NULL, TOKEN_AS:
		return true
	default:
		return false
	}
}

// expectIdentifier は識別子を期待し、予約語の場合はエラーを返す
func (p *Parser) expectIdentifier(context string) (string, error) {
	if p.curToken.Type != TOKEN_IDENT {
		if p.isReservedKeyword() {
			return "", p.newError("'%s' is a reserved keyword and cannot be used as %s", p.curToken.Literal, context)
		}
		return "", p.newError("expected %s", context)
	}
	name := p.curToken.Literal
	p.nextToken()
	return name, nil
}

// isIdentLike はトークンが識別子またはキーワード（識別子として使用可能なもの）かどうかを返す
func (p *Parser) isIdentLike() bool {
	switch p.curToken.Type {
	case TOKEN_IDENT,
		TOKEN_COMPONENT, TOKEN_IMPORT, TOKEN_TYPE, TOKEN_ENUM,
		TOKEN_DEPENDS, TOKEN_ON, TOKEN_EXTENDS, TOKEN_IMPLEMENTS,
		TOKEN_CONTAINS, TOKEN_AGGREGATES, TOKEN_PROVIDES, TOKEN_REQUIRES,
		TOKEN_FLOW, TOKEN_STATES, TOKEN_STATE, TOKEN_PARALLEL, TOKEN_REGION,
		TOKEN_INITIAL, TOKEN_FINAL, TOKEN_ENTRY, TOKEN_EXIT,
		TOKEN_IF, TOKEN_ELSE, TOKEN_FOR, TOKEN_IN, TOKEN_WHILE,
		TOKEN_RETURN, TOKEN_THROW, TOKEN_AWAIT, TOKEN_ASYNC, TOKEN_THROWS,
		TOKEN_WHEN, TOKEN_AFTER, TOKEN_DO, TOKEN_TRUE, TOKEN_FALSE, TOKEN_NULL, TOKEN_AS:
		return true
	default:
		return false
	}
}

func (p *Parser) curPos() ast.Position {
	return ast.Position{
		Line:   p.curToken.Line,
		Column: p.curToken.Column,
	}
}

func (p *Parser) newError(format string, args ...interface{}) error {
	return &errors.ParseError{
		Pos:     p.curPos(),
		Message: sprintf(format, args...),
	}
}

// newErrorWithSuggestion は期待されるトークンの候補を含むエラーを生成する
func (p *Parser) newErrorWithSuggestion(format string, got interface{}, suggestions ...string) error {
	msg := sprintf(format, got)
	if len(suggestions) > 0 {
		msg += " (expected one of: "
		for i, s := range suggestions {
			if i > 0 {
				msg += ", "
			}
			msg += "'" + s + "'"
		}
		msg += ")"
	}
	return &errors.ParseError{
		Pos:     p.curPos(),
		Message: msg,
	}
}

func sprintf(format string, args ...interface{}) string {
	if len(args) == 0 {
		return format
	}
	result := format
	for _, arg := range args {
		result = replaceFirst(result, "%v", toString(arg))
		result = replaceFirst(result, "%q", "\""+toString(arg)+"\"")
	}
	return result
}

func replaceFirst(s, old, new string) string {
	for i := 0; i <= len(s)-len(old); i++ {
		if s[i:i+len(old)] == old {
			return s[:i] + new + s[i+len(old):]
		}
	}
	return s
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	default:
		return ""
	}
}
