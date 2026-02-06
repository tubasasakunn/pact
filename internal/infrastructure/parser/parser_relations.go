package parser

import (
	"pact/internal/domain/ast"
)

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
