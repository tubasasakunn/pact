package parser

import (
	"pact/internal/domain/ast"
)

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
