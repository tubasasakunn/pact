package parser

import (
	"pact/internal/domain/ast"
)

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
