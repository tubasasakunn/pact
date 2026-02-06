package parser

import (
	"pact/internal/domain/ast"
)

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
