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
// Component
// =============================================================================

func (p *Parser) parseComponent(annotations []ast.AnnotationDecl) (*ast.ComponentDecl, error) {
	comp := &ast.ComponentDecl{
		Pos:         p.curPos(),
		Annotations: annotations,
	}

	p.nextToken() // consume 'component'

	name, err := p.expectIdentifier("component name")
	if err != nil {
		return nil, err
	}
	comp.Name = name

	if p.curToken.Type != TOKEN_LBRACE {
		return nil, p.newError("expected '{' after component name")
	}
	p.nextToken()

	// Parse body
	for p.curToken.Type != TOKEN_RBRACE && p.curToken.Type != TOKEN_EOF {
		if err := p.parseComponentBodyItem(&comp.Body); err != nil {
			return nil, err
		}
	}

	if p.curToken.Type != TOKEN_RBRACE {
		return nil, p.newError("expected '}' at end of component")
	}
	p.nextToken()

	return comp, nil
}

func (p *Parser) parseComponentBodyItem(body *ast.ComponentBody) error {
	// アノテーションをチェック
	var annotations []ast.AnnotationDecl
	if p.curToken.Type == TOKEN_AT {
		var err error
		annotations, err = p.parseAnnotations()
		if err != nil {
			return err
		}
	}

	switch p.curToken.Type {
	case TOKEN_TYPE:
		typ, err := p.parseTypeDecl(annotations)
		if err != nil {
			return err
		}
		body.Types = append(body.Types, *typ)

	case TOKEN_ENUM:
		typ, err := p.parseEnumDecl(annotations)
		if err != nil {
			return err
		}
		body.Types = append(body.Types, *typ)

	case TOKEN_DEPENDS:
		rel, err := p.parseDependsOn(annotations)
		if err != nil {
			return err
		}
		body.Relations = append(body.Relations, *rel)

	case TOKEN_EXTENDS:
		rel, err := p.parseSimpleRelation(ast.RelationExtends, annotations)
		if err != nil {
			return err
		}
		body.Relations = append(body.Relations, *rel)

	case TOKEN_IMPLEMENTS:
		rel, err := p.parseSimpleRelation(ast.RelationImplements, annotations)
		if err != nil {
			return err
		}
		body.Relations = append(body.Relations, *rel)

	case TOKEN_CONTAINS:
		rel, err := p.parseSimpleRelation(ast.RelationContains, annotations)
		if err != nil {
			return err
		}
		body.Relations = append(body.Relations, *rel)

	case TOKEN_AGGREGATES:
		rel, err := p.parseSimpleRelation(ast.RelationAggregates, annotations)
		if err != nil {
			return err
		}
		body.Relations = append(body.Relations, *rel)

	case TOKEN_PROVIDES:
		iface, err := p.parseInterface(annotations)
		if err != nil {
			return err
		}
		body.Provides = append(body.Provides, *iface)

	case TOKEN_REQUIRES:
		iface, err := p.parseInterface(annotations)
		if err != nil {
			return err
		}
		body.Requires = append(body.Requires, *iface)

	case TOKEN_FLOW:
		flow, err := p.parseFlow(annotations)
		if err != nil {
			return err
		}
		body.Flows = append(body.Flows, *flow)

	case TOKEN_STATES:
		states, err := p.parseStates(annotations)
		if err != nil {
			return err
		}
		body.States = append(body.States, *states)

	default:
		return p.newErrorWithSuggestion("unexpected token in component body: %v", p.curToken.Literal,
			"type", "enum", "provides", "requires", "depends", "flow", "states")
	}

	return nil
}

// =============================================================================
// Type & Enum
// =============================================================================

func (p *Parser) parseTypeDecl(annotations []ast.AnnotationDecl) (*ast.TypeDecl, error) {
	typ := &ast.TypeDecl{
		Pos:         p.curPos(),
		Kind:        ast.TypeKindStruct,
		Annotations: annotations,
	}

	p.nextToken() // consume 'type'

	name, err := p.expectIdentifier("type name")
	if err != nil {
		return nil, err
	}
	typ.Name = name

	// 型エイリアス: type UserId = string
	if p.curToken.Type == TOKEN_ASSIGN {
		p.nextToken()
		typ.Kind = ast.TypeKindAlias
		baseType, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		typ.BaseType = baseType
		return typ, nil
	}

	if p.curToken.Type != TOKEN_LBRACE {
		return nil, p.newError("expected '{' or '=' after type name")
	}
	p.nextToken()

	for p.curToken.Type != TOKEN_RBRACE && p.curToken.Type != TOKEN_EOF {
		field, err := p.parseField()
		if err != nil {
			return nil, err
		}
		typ.Fields = append(typ.Fields, *field)
	}

	if p.curToken.Type != TOKEN_RBRACE {
		return nil, p.newError("expected '}' at end of type")
	}
	p.nextToken()

	return typ, nil
}

func (p *Parser) parseEnumDecl(annotations []ast.AnnotationDecl) (*ast.TypeDecl, error) {
	typ := &ast.TypeDecl{
		Pos:         p.curPos(),
		Kind:        ast.TypeKindEnum,
		Annotations: annotations,
	}

	p.nextToken() // consume 'enum'

	name, err := p.expectIdentifier("enum name")
	if err != nil {
		return nil, err
	}
	typ.Name = name

	if p.curToken.Type != TOKEN_LBRACE {
		return nil, p.newError("expected '{' after enum name")
	}
	p.nextToken()

	for p.curToken.Type == TOKEN_IDENT {
		typ.Values = append(typ.Values, p.curToken.Literal)
		p.nextToken()
	}

	if p.curToken.Type != TOKEN_RBRACE {
		return nil, p.newError("expected '}' at end of enum")
	}
	p.nextToken()

	return typ, nil
}

func (p *Parser) parseField() (*ast.FieldDecl, error) {
	field := &ast.FieldDecl{
		Pos: p.curPos(),
	}

	// アノテーション
	if p.curToken.Type == TOKEN_AT {
		annotations, err := p.parseAnnotations()
		if err != nil {
			return nil, err
		}
		field.Annotations = annotations
	}

	// 可視性
	switch p.curToken.Type {
	case TOKEN_PLUS:
		field.Visibility = ast.VisibilityPublic
		p.nextToken()
	case TOKEN_MINUS:
		field.Visibility = ast.VisibilityPrivate
		p.nextToken()
	case TOKEN_HASH:
		field.Visibility = ast.VisibilityProtected
		p.nextToken()
	case TOKEN_TILDE:
		field.Visibility = ast.VisibilityPackage
		p.nextToken()
	}

	field.Pos = p.curPos()
	name, err := p.expectIdentifier("field name")
	if err != nil {
		return nil, err
	}
	field.Name = name

	if p.curToken.Type != TOKEN_COLON {
		return nil, p.newError("expected ':' after field name")
	}
	p.nextToken()

	typeExpr, err := p.parseTypeExpr()
	if err != nil {
		return nil, err
	}
	field.Type = *typeExpr

	return field, nil
}

func (p *Parser) parseTypeExpr() (*ast.TypeExpr, error) {
	typeExpr := &ast.TypeExpr{
		Pos: p.curPos(),
	}

	if p.curToken.Type != TOKEN_IDENT {
		return nil, p.newError("expected type name")
	}
	typeExpr.Name = p.curToken.Literal
	p.nextToken()

	// ジェネリクス型パラメータ: Type<T, U>
	if p.curToken.Type == TOKEN_LT {
		p.nextToken()
		for {
			param, err := p.parseTypeExpr()
			if err != nil {
				return nil, err
			}
			typeExpr.TypeParams = append(typeExpr.TypeParams, *param)
			if p.curToken.Type == TOKEN_COMMA {
				p.nextToken()
				continue
			}
			break
		}
		if p.curToken.Type != TOKEN_GT {
			return nil, p.newError("expected '>' after type parameters")
		}
		p.nextToken()
	}

	// 型修飾子のパース（無効なチェーンを検出）
	// 有効: Type?, Type[], Type?[], Type[]?
	// 無効: Type??, Type[][], Type?[]?, etc.

	// nullable?
	if p.curToken.Type == TOKEN_QUESTION {
		typeExpr.Nullable = true
		p.nextToken()

		// 連続した?は無効
		if p.curToken.Type == TOKEN_QUESTION {
			return nil, p.newError("invalid type modifier: multiple '?' not allowed")
		}
	}

	// array[]
	if p.curToken.Type == TOKEN_LBRACKET {
		p.nextToken()
		if p.curToken.Type != TOKEN_RBRACKET {
			return nil, p.newError("expected ']' after '['")
		}
		typeExpr.Array = true
		p.nextToken()

		// 連続した[]は無効
		if p.curToken.Type == TOKEN_LBRACKET {
			return nil, p.newError("invalid type modifier: nested arrays not allowed")
		}
	}

	// 配列の後のnullable
	if p.curToken.Type == TOKEN_QUESTION {
		if typeExpr.Nullable {
			// Type?[]? のような形式は許可しない
			return nil, p.newError("invalid type modifier: '?' already specified before array")
		}
		typeExpr.Nullable = true
		p.nextToken()

		// さらに続く修飾子は無効
		if p.curToken.Type == TOKEN_QUESTION || p.curToken.Type == TOKEN_LBRACKET {
			return nil, p.newError("invalid type modifier chain")
		}
	}

	return typeExpr, nil
}

// =============================================================================
// Relations
// =============================================================================

func (p *Parser) parseDependsOn(annotations []ast.AnnotationDecl) (*ast.RelationDecl, error) {
	rel := &ast.RelationDecl{
		Pos:         p.curPos(),
		Kind:        ast.RelationDependsOn,
		Annotations: annotations,
	}

	p.nextToken() // consume 'depends'

	if p.curToken.Type != TOKEN_ON {
		return nil, p.newError("expected 'on' after 'depends'")
	}
	p.nextToken()

	if p.curToken.Type != TOKEN_IDENT {
		return nil, p.newError("expected identifier after 'depends on'")
	}
	rel.Target = p.curToken.Literal
	p.nextToken()

	// : type
	if p.curToken.Type == TOKEN_COLON {
		p.nextToken()
		if p.curToken.Type != TOKEN_IDENT {
			return nil, p.newError("expected type after ':'")
		}
		t := p.curToken.Literal
		rel.TargetType = &t
		p.nextToken()
	}

	// as alias
	if p.curToken.Type == TOKEN_AS {
		p.nextToken()
		if p.curToken.Type != TOKEN_IDENT {
			return nil, p.newError("expected identifier after 'as'")
		}
		a := p.curToken.Literal
		rel.Alias = &a
		p.nextToken()
	}

	return rel, nil
}

func (p *Parser) parseSimpleRelation(kind ast.RelationKind, annotations []ast.AnnotationDecl) (*ast.RelationDecl, error) {
	rel := &ast.RelationDecl{
		Pos:         p.curPos(),
		Kind:        kind,
		Annotations: annotations,
	}

	p.nextToken() // consume keyword

	if p.curToken.Type != TOKEN_IDENT {
		return nil, p.newError("expected identifier")
	}
	rel.Target = p.curToken.Literal
	p.nextToken()

	return rel, nil
}

// =============================================================================
// Interface
// =============================================================================

func (p *Parser) parseInterface(annotations []ast.AnnotationDecl) (*ast.InterfaceDecl, error) {
	iface := &ast.InterfaceDecl{
		Pos:         p.curPos(),
		Annotations: annotations,
	}

	p.nextToken() // consume 'provides' or 'requires'

	name, err := p.expectIdentifier("interface name")
	if err != nil {
		return nil, err
	}
	iface.Name = name

	if p.curToken.Type != TOKEN_LBRACE {
		return nil, p.newError("expected '{' after interface name")
	}
	p.nextToken()

	for p.curToken.Type != TOKEN_RBRACE && p.curToken.Type != TOKEN_EOF {
		method, err := p.parseMethod()
		if err != nil {
			return nil, err
		}
		iface.Methods = append(iface.Methods, *method)
	}

	if p.curToken.Type != TOKEN_RBRACE {
		return nil, p.newError("expected '}' at end of interface")
	}
	p.nextToken()

	return iface, nil
}

func (p *Parser) parseMethod() (*ast.MethodDecl, error) {
	method := &ast.MethodDecl{
		Pos: p.curPos(),
	}

	// アノテーション
	if p.curToken.Type == TOKEN_AT {
		annotations, err := p.parseAnnotations()
		if err != nil {
			return nil, err
		}
		method.Annotations = annotations
	}

	// async
	if p.curToken.Type == TOKEN_ASYNC {
		method.Async = true
		p.nextToken()
	}

	method.Pos = p.curPos()
	name, err := p.expectIdentifier("method name")
	if err != nil {
		return nil, err
	}
	method.Name = name

	// パラメータ
	if p.curToken.Type != TOKEN_LPAREN {
		return nil, p.newError("expected '(' after method name")
	}
	p.nextToken()

	for p.curToken.Type != TOKEN_RPAREN && p.curToken.Type != TOKEN_EOF {
		param, err := p.parseParam()
		if err != nil {
			return nil, err
		}
		method.Params = append(method.Params, *param)

		if p.curToken.Type == TOKEN_COMMA {
			p.nextToken()
		}
	}

	if p.curToken.Type != TOKEN_RPAREN {
		return nil, p.newError("expected ')' after parameters")
	}
	p.nextToken()

	// 戻り値
	if p.curToken.Type == TOKEN_ARROW {
		p.nextToken()
		typeExpr, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		method.ReturnType = typeExpr
	}

	// throws
	if p.curToken.Type == TOKEN_THROWS {
		p.nextToken()
		for {
			if p.curToken.Type != TOKEN_IDENT {
				return nil, p.newError("expected error type after 'throws'")
			}
			method.Throws = append(method.Throws, p.curToken.Literal)
			p.nextToken()

			if p.curToken.Type != TOKEN_COMMA {
				break
			}
			p.nextToken()
		}
	}

	return method, nil
}

func (p *Parser) parseParam() (*ast.ParamDecl, error) {
	param := &ast.ParamDecl{
		Pos: p.curPos(),
	}

	name, err := p.expectIdentifier("parameter name")
	if err != nil {
		return nil, err
	}
	param.Name = name

	if p.curToken.Type != TOKEN_COLON {
		return nil, p.newError("expected ':' after parameter name")
	}
	p.nextToken()

	typeExpr, err := p.parseTypeExpr()
	if err != nil {
		return nil, err
	}
	param.Type = *typeExpr

	return param, nil
}

// =============================================================================
// Flow
// =============================================================================

func (p *Parser) parseFlow(annotations []ast.AnnotationDecl) (*ast.FlowDecl, error) {
	flow := &ast.FlowDecl{
		Pos:         p.curPos(),
		Annotations: annotations,
	}

	p.nextToken() // consume 'flow'

	name, err := p.expectIdentifier("flow name")
	if err != nil {
		return nil, err
	}
	flow.Name = name

	if p.curToken.Type != TOKEN_LBRACE {
		return nil, p.newError("expected '{' after flow name")
	}
	p.nextToken()

	steps, err := p.parseSteps()
	if err != nil {
		return nil, err
	}
	flow.Steps = steps

	if p.curToken.Type != TOKEN_RBRACE {
		return nil, p.newError("expected '}' at end of flow")
	}
	p.nextToken()

	return flow, nil
}

func (p *Parser) parseSteps() ([]ast.Step, error) {
	steps := []ast.Step{} // empty slice, not nil

	for p.curToken.Type != TOKEN_RBRACE && p.curToken.Type != TOKEN_EOF {
		step, err := p.parseStep()
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}

	return steps, nil
}

func (p *Parser) parseStep() (ast.Step, error) {
	// アノテーション
	var annotations []ast.AnnotationDecl
	if p.curToken.Type == TOKEN_AT {
		var err error
		annotations, err = p.parseAnnotations()
		if err != nil {
			return nil, err
		}
	}

	pos := p.curPos()

	switch p.curToken.Type {
	case TOKEN_RETURN:
		return p.parseReturnStep(pos, annotations)
	case TOKEN_THROW:
		return p.parseThrowStep(pos, annotations)
	case TOKEN_IF:
		return p.parseIfStep(pos, annotations)
	case TOKEN_FOR:
		return p.parseForStep(pos, annotations)
	case TOKEN_WHILE:
		return p.parseWhileStep(pos, annotations)
	case TOKEN_AWAIT:
		return p.parseAwaitStep(pos, annotations)
	case TOKEN_IDENT:
		return p.parseAssignOrCallStep(pos, annotations)
	default:
		return nil, p.newErrorWithSuggestion("unexpected token in flow: %v", p.curToken.Literal,
			"return", "throw", "if", "for", "while", "await", "identifier")
	}
}

func (p *Parser) parseReturnStep(pos ast.Position, annotations []ast.AnnotationDecl) (*ast.ReturnStep, error) {
	step := &ast.ReturnStep{
		Pos:         pos,
		Annotations: annotations,
	}

	p.nextToken() // consume 'return'

	// 値があるか
	if p.curToken.Type != TOKEN_RBRACE && p.curToken.Type != TOKEN_EOF &&
		p.curToken.Type != TOKEN_IF && p.curToken.Type != TOKEN_FOR &&
		p.curToken.Type != TOKEN_WHILE && p.curToken.Type != TOKEN_RETURN &&
		p.curToken.Type != TOKEN_THROW && p.curToken.Type != TOKEN_AT {
		expr, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		step.Value = expr
	}

	return step, nil
}

func (p *Parser) parseThrowStep(pos ast.Position, annotations []ast.AnnotationDecl) (*ast.ThrowStep, error) {
	step := &ast.ThrowStep{
		Pos:         pos,
		Annotations: annotations,
	}

	p.nextToken() // consume 'throw'

	if p.curToken.Type != TOKEN_IDENT {
		return nil, p.newError("expected error type after 'throw'")
	}
	step.Error = p.curToken.Literal
	p.nextToken()

	return step, nil
}

func (p *Parser) parseIfStep(pos ast.Position, annotations []ast.AnnotationDecl) (*ast.IfStep, error) {
	step := &ast.IfStep{
		Pos:         pos,
		Annotations: annotations,
	}

	p.nextToken() // consume 'if'

	cond, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	step.Condition = cond

	if p.curToken.Type != TOKEN_LBRACE {
		return nil, p.newError("expected '{' after if condition")
	}
	p.nextToken()

	thenSteps, err := p.parseSteps()
	if err != nil {
		return nil, err
	}
	step.Then = thenSteps

	if p.curToken.Type != TOKEN_RBRACE {
		return nil, p.newError("expected '}' after then block")
	}
	p.nextToken()

	// else
	if p.curToken.Type == TOKEN_ELSE {
		p.nextToken()
		if p.curToken.Type != TOKEN_LBRACE {
			return nil, p.newError("expected '{' after else")
		}
		p.nextToken()

		elseSteps, err := p.parseSteps()
		if err != nil {
			return nil, err
		}
		step.Else = elseSteps

		if p.curToken.Type != TOKEN_RBRACE {
			return nil, p.newError("expected '}' after else block")
		}
		p.nextToken()
	}

	return step, nil
}

func (p *Parser) parseForStep(pos ast.Position, annotations []ast.AnnotationDecl) (*ast.ForStep, error) {
	step := &ast.ForStep{
		Pos:         pos,
		Annotations: annotations,
	}

	p.nextToken() // consume 'for'

	if p.curToken.Type != TOKEN_IDENT {
		return nil, p.newError("expected variable name after 'for'")
	}
	step.Variable = p.curToken.Literal
	p.nextToken()

	if p.curToken.Type != TOKEN_IN {
		return nil, p.newError("expected 'in' after variable")
	}
	p.nextToken()

	iter, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	step.Iterable = iter

	if p.curToken.Type != TOKEN_LBRACE {
		return nil, p.newError("expected '{' after for expression")
	}
	p.nextToken()

	body, err := p.parseSteps()
	if err != nil {
		return nil, err
	}
	step.Body = body

	if p.curToken.Type != TOKEN_RBRACE {
		return nil, p.newError("expected '}' after for body")
	}
	p.nextToken()

	return step, nil
}

func (p *Parser) parseWhileStep(pos ast.Position, annotations []ast.AnnotationDecl) (*ast.WhileStep, error) {
	step := &ast.WhileStep{
		Pos:         pos,
		Annotations: annotations,
	}

	p.nextToken() // consume 'while'

	cond, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	step.Condition = cond

	if p.curToken.Type != TOKEN_LBRACE {
		return nil, p.newError("expected '{' after while condition")
	}
	p.nextToken()

	body, err := p.parseSteps()
	if err != nil {
		return nil, err
	}
	step.Body = body

	if p.curToken.Type != TOKEN_RBRACE {
		return nil, p.newError("expected '}' after while body")
	}
	p.nextToken()

	return step, nil
}

func (p *Parser) parseAwaitStep(pos ast.Position, annotations []ast.AnnotationDecl) (*ast.CallStep, error) {
	p.nextToken() // consume 'await'

	expr, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}

	return &ast.CallStep{
		Pos:         pos,
		Expr:        expr,
		Await:       true,
		Annotations: annotations,
	}, nil
}

func (p *Parser) parseAssignOrCallStep(pos ast.Position, annotations []ast.AnnotationDecl) (ast.Step, error) {
	// 最初の識別子を見る
	name := p.curToken.Literal
	p.nextToken()

	// = があれば代入
	if p.curToken.Type == TOKEN_ASSIGN {
		p.nextToken()
		expr, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}

		return &ast.AssignStep{
			Pos:         pos,
			Variable:    name,
			Value:       expr,
			Annotations: annotations,
		}, nil
	}

	// そうでなければ呼び出し (a.b() など)
	var obj ast.Expr = &ast.VariableExpr{Pos: pos, Name: name}

	for {
		if p.curToken.Type == TOKEN_DOT {
			p.nextToken()
			if p.curToken.Type != TOKEN_IDENT {
				return nil, p.newError("expected identifier after '.'")
			}
			fieldName := p.curToken.Literal
			fieldPos := p.curPos()
			p.nextToken()

			if p.curToken.Type == TOKEN_LPAREN {
				// メソッド呼び出し
				args, err := p.parseCallArgs()
				if err != nil {
					return nil, err
				}
				obj = &ast.CallExpr{
					Pos:    fieldPos,
					Object: obj,
					Method: fieldName,
					Args:   args,
				}
			} else {
				obj = &ast.FieldExpr{
					Pos:    fieldPos,
					Object: obj,
					Field:  fieldName,
				}
			}
		} else {
			break
		}
	}

	return &ast.CallStep{
		Pos:         pos,
		Expr:        obj,
		Annotations: annotations,
	}, nil
}

// =============================================================================
// Expression
// =============================================================================

func (p *Parser) parseExpression(precedence int) (ast.Expr, error) {
	left, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}

	for {
		prec := p.curPrecedence()
		if prec <= precedence {
			break
		}

		switch p.curToken.Type {
		case TOKEN_PLUS, TOKEN_MINUS, TOKEN_STAR, TOKEN_SLASH, TOKEN_PERCENT,
			TOKEN_EQ, TOKEN_NE, TOKEN_LT, TOKEN_GT, TOKEN_LE, TOKEN_GE,
			TOKEN_AND, TOKEN_OR:
			op := p.curToken.Literal
			opPrec := prec
			p.nextToken()
			right, err := p.parseExpression(opPrec)
			if err != nil {
				return nil, err
			}
			left = &ast.BinaryExpr{
				Pos:   left.GetPos(),
				Left:  left,
				Op:    op,
				Right: right,
			}

		case TOKEN_QUESTION:
			// 三項演算子
			p.nextToken()
			then, err := p.parseExpression(0)
			if err != nil {
				return nil, err
			}
			if p.curToken.Type != TOKEN_COLON {
				return nil, p.newError("expected ':' in ternary expression")
			}
			p.nextToken()
			els, err := p.parseExpression(0)
			if err != nil {
				return nil, err
			}
			left = &ast.TernaryExpr{
				Pos:       left.GetPos(),
				Condition: left,
				Then:      then,
				Else:      els,
			}

		case TOKEN_NULLISH:
			// null合体
			p.nextToken()
			if p.curToken.Type == TOKEN_THROW {
				p.nextToken()
				if p.curToken.Type != TOKEN_IDENT {
					return nil, p.newError("expected error type after 'throw'")
				}
				errName := p.curToken.Literal
				p.nextToken()
				left = &ast.NullishExpr{
					Pos:      left.GetPos(),
					Left:     left,
					ThrowErr: &errName,
				}
			} else {
				right, err := p.parseExpression(prec)
				if err != nil {
					return nil, err
				}
				left = &ast.NullishExpr{
					Pos:   left.GetPos(),
					Left:  left,
					Right: right,
				}
			}

		case TOKEN_DOT:
			p.nextToken()
			if p.curToken.Type != TOKEN_IDENT {
				return nil, p.newError("expected identifier after '.'")
			}
			fieldName := p.curToken.Literal
			fieldPos := p.curPos()
			p.nextToken()

			if p.curToken.Type == TOKEN_LPAREN {
				args, err := p.parseCallArgs()
				if err != nil {
					return nil, err
				}
				left = &ast.CallExpr{
					Pos:    fieldPos,
					Object: left,
					Method: fieldName,
					Args:   args,
				}
			} else {
				left = &ast.FieldExpr{
					Pos:    fieldPos,
					Object: left,
					Field:  fieldName,
				}
			}

		default:
			return left, nil
		}
	}

	return left, nil
}

func (p *Parser) parsePrimaryExpr() (ast.Expr, error) {
	pos := p.curPos()

	switch p.curToken.Type {
	case TOKEN_INT:
		val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
		if err != nil {
			return nil, p.newError("invalid integer literal: %s", p.curToken.Literal)
		}
		p.nextToken()
		return &ast.LiteralExpr{Pos: pos, Value: val}, nil

	case TOKEN_FLOAT:
		val, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			return nil, p.newError("invalid float literal: %s", p.curToken.Literal)
		}
		p.nextToken()
		return &ast.LiteralExpr{Pos: pos, Value: val}, nil

	case TOKEN_STRING:
		val := p.curToken.Literal
		p.nextToken()
		return &ast.LiteralExpr{Pos: pos, Value: val}, nil

	case TOKEN_TRUE:
		p.nextToken()
		return &ast.LiteralExpr{Pos: pos, Value: true}, nil

	case TOKEN_FALSE:
		p.nextToken()
		return &ast.LiteralExpr{Pos: pos, Value: false}, nil

	case TOKEN_NULL:
		p.nextToken()
		return &ast.LiteralExpr{Pos: pos, Value: nil}, nil

	case TOKEN_IDENT:
		name := p.curToken.Literal
		p.nextToken()
		return &ast.VariableExpr{Pos: pos, Name: name}, nil

	case TOKEN_NOT:
		p.nextToken()
		operand, err := p.parsePrimaryExpr()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Pos: pos, Op: "!", Operand: operand}, nil

	case TOKEN_MINUS:
		p.nextToken()
		// 負の数値リテラルの最適化: -5 は単一のLiteralExprとして扱う
		if p.curToken.Type == TOKEN_INT {
			val, err := strconv.ParseInt("-"+p.curToken.Literal, 10, 64)
			if err != nil {
				return nil, p.newError("invalid integer literal: -%s", p.curToken.Literal)
			}
			p.nextToken()
			return &ast.LiteralExpr{Pos: pos, Value: val}, nil
		}
		if p.curToken.Type == TOKEN_FLOAT {
			val, err := strconv.ParseFloat("-"+p.curToken.Literal, 64)
			if err != nil {
				return nil, p.newError("invalid float literal: -%s", p.curToken.Literal)
			}
			p.nextToken()
			return &ast.LiteralExpr{Pos: pos, Value: val}, nil
		}
		operand, err := p.parsePrimaryExpr()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Pos: pos, Op: "-", Operand: operand}, nil

	case TOKEN_LPAREN:
		p.nextToken()
		expr, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		if p.curToken.Type != TOKEN_RPAREN {
			return nil, p.newError("expected ')' after expression")
		}
		p.nextToken()
		return expr, nil

	default:
		return nil, p.newErrorWithSuggestion("unexpected token in expression: %v", p.curToken.Literal,
			"identifier", "number", "string", "true", "false", "null", "(", "!", "-")
	}
}

func (p *Parser) parseCallArgs() ([]ast.Expr, error) {
	var args []ast.Expr

	p.nextToken() // consume '('

	for p.curToken.Type != TOKEN_RPAREN && p.curToken.Type != TOKEN_EOF {
		arg, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		if p.curToken.Type == TOKEN_COMMA {
			p.nextToken()
		}
	}

	if p.curToken.Type != TOKEN_RPAREN {
		return nil, p.newError("expected ')' after arguments")
	}
	p.nextToken()

	return args, nil
}

func (p *Parser) curPrecedence() int {
	switch p.curToken.Type {
	case TOKEN_OR:
		return 1
	case TOKEN_AND:
		return 2
	case TOKEN_EQ, TOKEN_NE:
		return 3
	case TOKEN_LT, TOKEN_GT, TOKEN_LE, TOKEN_GE:
		return 4
	case TOKEN_PLUS, TOKEN_MINUS:
		return 5
	case TOKEN_STAR, TOKEN_SLASH, TOKEN_PERCENT:
		return 6
	case TOKEN_DOT:
		return 8
	case TOKEN_QUESTION:
		return 1
	case TOKEN_NULLISH:
		return 1
	default:
		return 0
	}
}

// =============================================================================
// States
// =============================================================================

func (p *Parser) parseStates(annotations []ast.AnnotationDecl) (*ast.StatesDecl, error) {
	states := &ast.StatesDecl{
		Pos:         p.curPos(),
		Annotations: annotations,
	}

	p.nextToken() // consume 'states'

	name, err := p.expectIdentifier("states name")
	if err != nil {
		return nil, err
	}
	states.Name = name

	if p.curToken.Type != TOKEN_LBRACE {
		return nil, p.newError("expected '{' after states name")
	}
	p.nextToken()

	for p.curToken.Type != TOKEN_RBRACE && p.curToken.Type != TOKEN_EOF {
		if err := p.parseStatesItem(states); err != nil {
			return nil, err
		}
	}

	if p.curToken.Type != TOKEN_RBRACE {
		return nil, p.newError("expected '}' at end of states")
	}
	p.nextToken()

	return states, nil
}

func (p *Parser) parseStatesItem(states *ast.StatesDecl) error {
	switch p.curToken.Type {
	case TOKEN_INITIAL:
		p.nextToken()
		if p.curToken.Type != TOKEN_IDENT {
			return p.newError("expected state name after 'initial'")
		}
		states.Initial = p.curToken.Literal
		p.nextToken()

	case TOKEN_FINAL:
		p.nextToken()
		if p.curToken.Type != TOKEN_IDENT {
			return p.newError("expected state name after 'final'")
		}
		states.Finals = append(states.Finals, p.curToken.Literal)
		p.nextToken()

	case TOKEN_STATE:
		state, err := p.parseStateDecl()
		if err != nil {
			return err
		}
		states.States = append(states.States, *state)

	case TOKEN_PARALLEL:
		parallel, err := p.parseParallel()
		if err != nil {
			return err
		}
		states.Parallels = append(states.Parallels, *parallel)

	case TOKEN_IDENT:
		// 遷移: From -> To ...
		trans, err := p.parseTransition()
		if err != nil {
			return err
		}
		states.Transitions = append(states.Transitions, *trans)

	default:
		return p.newErrorWithSuggestion("unexpected token in states: %v", p.curToken.Literal,
			"initial", "final", "state", "parallel", "transition (From -> To)")
	}

	return nil
}

func (p *Parser) parseStateDecl() (*ast.StateDecl, error) {
	state := &ast.StateDecl{
		Pos: p.curPos(),
	}

	p.nextToken() // consume 'state'

	name, err := p.expectIdentifier("state name")
	if err != nil {
		return nil, err
	}
	state.Name = name

	if p.curToken.Type != TOKEN_LBRACE {
		return nil, p.newError("expected '{' after state name")
	}
	p.nextToken()

	for p.curToken.Type != TOKEN_RBRACE && p.curToken.Type != TOKEN_EOF {
		switch p.curToken.Type {
		case TOKEN_ENTRY:
			p.nextToken()
			actions, err := p.parseActionList()
			if err != nil {
				return nil, err
			}
			state.Entry = actions

		case TOKEN_EXIT:
			p.nextToken()
			actions, err := p.parseActionList()
			if err != nil {
				return nil, err
			}
			state.Exit = actions

		case TOKEN_INITIAL:
			p.nextToken()
			if p.curToken.Type != TOKEN_IDENT {
				return nil, p.newError("expected state name after 'initial'")
			}
			initial := p.curToken.Literal
			state.Initial = &initial
			p.nextToken()

		case TOKEN_STATE:
			nested, err := p.parseStateDecl()
			if err != nil {
				return nil, err
			}
			state.States = append(state.States, *nested)

		case TOKEN_IDENT:
			trans, err := p.parseTransition()
			if err != nil {
				return nil, err
			}
			state.Transitions = append(state.Transitions, *trans)

		default:
			return nil, p.newErrorWithSuggestion("unexpected token in state: %v", p.curToken.Literal,
				"entry", "exit", "initial", "state", "transition (From -> To)")
		}
	}

	if p.curToken.Type != TOKEN_RBRACE {
		return nil, p.newError("expected '}' at end of state")
	}
	p.nextToken()

	return state, nil
}

func (p *Parser) parseTransition() (*ast.TransitionDecl, error) {
	trans := &ast.TransitionDecl{
		Pos: p.curPos(),
	}

	if p.curToken.Type != TOKEN_IDENT {
		return nil, p.newError("expected from state")
	}
	trans.From = p.curToken.Literal
	p.nextToken()

	if p.curToken.Type != TOKEN_ARROW {
		return nil, p.newError("expected '->' in transition")
	}
	p.nextToken()

	if p.curToken.Type != TOKEN_IDENT {
		return nil, p.newError("expected to state")
	}
	trans.To = p.curToken.Literal
	p.nextToken()

	// トリガー: on E, after 3s, when cond
	switch p.curToken.Type {
	case TOKEN_ON:
		p.nextToken()
		if p.curToken.Type != TOKEN_IDENT {
			return nil, p.newError("expected event name after 'on'")
		}
		trans.Trigger = &ast.EventTrigger{
			Pos:   p.curPos(),
			Event: p.curToken.Literal,
		}
		p.nextToken()

	case TOKEN_AFTER:
		p.nextToken()
		if p.curToken.Type != TOKEN_DURATION {
			return nil, p.newError("expected duration after 'after'")
		}
		duration, err := p.parseDuration(p.curToken.Literal)
		if err != nil {
			return nil, err
		}
		trans.Trigger = &ast.AfterTrigger{
			Pos:      p.curPos(),
			Duration: duration,
		}
		p.nextToken()

	case TOKEN_WHEN:
		p.nextToken()
		cond, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		trans.Trigger = &ast.WhenTrigger{
			Pos:       cond.GetPos(),
			Condition: cond,
		}
	}

	// ガード: when cond (on E の後)
	if p.curToken.Type == TOKEN_WHEN && trans.Trigger != nil {
		p.nextToken()
		guard, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		trans.Guard = guard
	}

	// アクション: do [a, b]
	if p.curToken.Type == TOKEN_DO {
		p.nextToken()
		actions, err := p.parseActionList()
		if err != nil {
			return nil, err
		}
		trans.Actions = actions
	}

	return trans, nil
}

func (p *Parser) parseParallel() (*ast.ParallelDecl, error) {
	parallel := &ast.ParallelDecl{
		Pos: p.curPos(),
	}

	p.nextToken() // consume 'parallel'

	name, err := p.expectIdentifier("parallel state name")
	if err != nil {
		return nil, err
	}
	parallel.Name = name

	if p.curToken.Type != TOKEN_LBRACE {
		return nil, p.newError("expected '{' after parallel name")
	}
	p.nextToken()

	for p.curToken.Type == TOKEN_REGION {
		region, err := p.parseRegion()
		if err != nil {
			return nil, err
		}
		parallel.Regions = append(parallel.Regions, *region)
	}

	if p.curToken.Type != TOKEN_RBRACE {
		return nil, p.newError("expected '}' at end of parallel")
	}
	p.nextToken()

	return parallel, nil
}

func (p *Parser) parseRegion() (*ast.RegionDecl, error) {
	region := &ast.RegionDecl{
		Pos: p.curPos(),
	}

	p.nextToken() // consume 'region'

	name, err := p.expectIdentifier("region name")
	if err != nil {
		return nil, err
	}
	region.Name = name

	if p.curToken.Type != TOKEN_LBRACE {
		return nil, p.newError("expected '{' after region name")
	}
	p.nextToken()

	for p.curToken.Type != TOKEN_RBRACE && p.curToken.Type != TOKEN_EOF {
		switch p.curToken.Type {
		case TOKEN_INITIAL:
			p.nextToken()
			if p.curToken.Type != TOKEN_IDENT {
				return nil, p.newError("expected state name after 'initial'")
			}
			region.Initial = p.curToken.Literal
			p.nextToken()

		case TOKEN_STATE:
			state, err := p.parseStateDecl()
			if err != nil {
				return nil, err
			}
			region.States = append(region.States, *state)

		case TOKEN_IDENT:
			trans, err := p.parseTransition()
			if err != nil {
				return nil, err
			}
			region.Transitions = append(region.Transitions, *trans)

		default:
			return nil, p.newErrorWithSuggestion("unexpected token in region: %v", p.curToken.Literal,
				"initial", "state", "transition (From -> To)")
		}
	}

	if p.curToken.Type != TOKEN_RBRACE {
		return nil, p.newError("expected '}' at end of region")
	}
	p.nextToken()

	return region, nil
}

func (p *Parser) parseActionList() ([]string, error) {
	var actions []string

	if p.curToken.Type != TOKEN_LBRACKET {
		return nil, p.newError("expected '[' for action list")
	}
	p.nextToken()

	for p.curToken.Type != TOKEN_RBRACKET && p.curToken.Type != TOKEN_EOF {
		if p.curToken.Type != TOKEN_IDENT {
			return nil, p.newError("expected action name")
		}
		actions = append(actions, p.curToken.Literal)
		p.nextToken()

		if p.curToken.Type == TOKEN_COMMA {
			p.nextToken()
		}
	}

	if p.curToken.Type != TOKEN_RBRACKET {
		return nil, p.newError("expected ']' at end of action list")
	}
	p.nextToken()

	return actions, nil
}

func (p *Parser) parseDuration(literal string) (ast.Duration, error) {
	// "500ms", "30s", "5m", "24h", "7d"
	var value int
	var unit string

	for i, c := range literal {
		if c >= '0' && c <= '9' {
			continue
		}
		var err error
		value, err = strconv.Atoi(literal[:i])
		if err != nil {
			return ast.Duration{}, p.newError("invalid duration value: %s", literal)
		}
		unit = literal[i:]
		break
	}

	return ast.Duration{Value: value, Unit: unit}, nil
}

// =============================================================================
// Annotations
// =============================================================================

func (p *Parser) parseAnnotations() ([]ast.AnnotationDecl, error) {
	var annotations []ast.AnnotationDecl

	for p.curToken.Type == TOKEN_AT {
		ann, err := p.parseAnnotation()
		if err != nil {
			return nil, err
		}
		annotations = append(annotations, *ann)
	}

	return annotations, nil
}

func (p *Parser) parseAnnotation() (*ast.AnnotationDecl, error) {
	ann := &ast.AnnotationDecl{
		Pos: p.curPos(),
	}

	p.nextToken() // consume '@'

	// アノテーション名はキーワードも許可する
	if !p.isIdentLike() {
		return nil, p.newError("expected annotation name after '@'")
	}
	ann.Name = p.curToken.Literal
	p.nextToken()

	// 引数があるか
	if p.curToken.Type == TOKEN_LPAREN {
		p.nextToken()

		for p.curToken.Type != TOKEN_RPAREN && p.curToken.Type != TOKEN_EOF {
			arg, err := p.parseAnnotationArg()
			if err != nil {
				return nil, err
			}
			ann.Args = append(ann.Args, *arg)

			if p.curToken.Type == TOKEN_COMMA {
				p.nextToken()
			}
		}

		if p.curToken.Type != TOKEN_RPAREN {
			return nil, p.newError("expected ')' after annotation arguments")
		}
		p.nextToken()
	}

	return ann, nil
}

func (p *Parser) parseAnnotationArg() (*ast.AnnotationArg, error) {
	arg := &ast.AnnotationArg{}

	// key: value or just value
	if p.curToken.Type == TOKEN_IDENT && p.peekToken.Type == TOKEN_COLON {
		key := p.curToken.Literal
		arg.Key = &key
		p.nextToken() // ident
		p.nextToken() // colon
	}

	if p.curToken.Type != TOKEN_STRING {
		return nil, p.newError("expected string value in annotation")
	}
	arg.Value = p.curToken.Literal
	p.nextToken()

	return arg, nil
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

