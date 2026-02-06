package parser

import (
	"strconv"

	"pact/internal/domain/ast"
)

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
