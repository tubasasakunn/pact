package parser

import (
	"strconv"

	"pact/internal/domain/ast"
)

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
