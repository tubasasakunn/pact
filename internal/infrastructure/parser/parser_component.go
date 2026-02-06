package parser

import (
	"pact/internal/domain/ast"
)

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
