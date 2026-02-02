package integration

import (
	"testing"

	"pact/internal/application/transformer"
	"pact/internal/domain/ast"
	"pact/internal/infrastructure/parser"
)

// =============================================================================
// I001-I005: パース→変換テスト
// =============================================================================

// I001: パース→クラス図変換
func TestIntegration_ParseToClassDiagram(t *testing.T) {
	input := `
component UserService {
	type User {
		id: string
		name: string
	}

	provides API {
		GetUser(id: string) -> User
	}
}
`
	lexer := parser.NewLexer(input)
	p := parser.NewParser(lexer)

	spec, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	tr := transformer.NewClassTransformer()
	diagram, err := tr.Transform([]*ast.SpecFile{spec}, nil)
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}

	if len(diagram.Nodes) == 0 {
		t.Error("expected at least one node")
	}
}

// I002: パース→シーケンス図変換
func TestIntegration_ParseToSequenceDiagram(t *testing.T) {
	input := `
component AuthService {
	depends on Database

	flow Login {
		result = Database.validate(request)
		return result
	}
}
`
	lexer := parser.NewLexer(input)
	p := parser.NewParser(lexer)

	spec, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	tr := transformer.NewSequenceTransformer()
	diagram, err := tr.Transform([]*ast.SpecFile{spec}, &transformer.SequenceOptions{FlowName: "Login"})
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}

	if len(diagram.Participants) == 0 {
		t.Error("expected participants")
	}
}

// I003: パース→状態図変換
func TestIntegration_ParseToStateDiagram(t *testing.T) {
	input := `
component OrderService {
	states OrderState {
		initial Pending
		final Completed

		state Pending { }
		state Processing { }
		state Completed { }

		Pending -> Processing on process
		Processing -> Completed on complete
	}
}
`
	lexer := parser.NewLexer(input)
	p := parser.NewParser(lexer)

	spec, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	tr := transformer.NewStateTransformer()
	diagram, err := tr.Transform([]*ast.SpecFile{spec}, &transformer.StateOptions{StatesName: "OrderState"})
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}

	if len(diagram.States) == 0 {
		t.Error("expected states")
	}
}

// I004: パース→フローチャート変換
func TestIntegration_ParseToFlowchart(t *testing.T) {
	input := `
component PaymentService {
	flow ProcessPayment {
		valid = self.validate(input)
		if valid {
			result = self.process(payment)
		} else {
			throw ValidationError
		}
		return result
	}
}
`
	lexer := parser.NewLexer(input)
	p := parser.NewParser(lexer)

	spec, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	tr := transformer.NewFlowTransformer()
	diagram, err := tr.Transform([]*ast.SpecFile{spec}, &transformer.FlowOptions{FlowName: "ProcessPayment"})
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}

	if len(diagram.Nodes) == 0 {
		t.Error("expected nodes")
	}
}

// I005: 複数ファイルのパース
func TestIntegration_ParseMultipleFiles(t *testing.T) {
	inputs := []string{
		`component ServiceA { type Data { id: string } }`,
		`component ServiceB { type Info { name: string } }`,
		`component ServiceC { type Blob { data: byte[] } }`,
	}

	for i, input := range inputs {
		lexer := parser.NewLexer(input)
		p := parser.NewParser(lexer)

		spec, err := p.Parse()
		if err != nil {
			t.Errorf("parse error for input %d: %v", i, err)
			continue
		}

		if spec.Component == nil {
			t.Errorf("expected component for input %d", i)
		}
	}
}
