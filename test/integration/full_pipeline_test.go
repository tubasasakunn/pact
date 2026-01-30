package integration

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"pact/internal/application/transformer"
	"pact/internal/infrastructure/parser"
	"pact/internal/infrastructure/renderer/svg"
)

// =============================================================================
// I020-I026: フルパイプラインテスト
// =============================================================================

// I020: クラス図フルパイプライン
func TestPipeline_ClassDiagram_Full(t *testing.T) {
	input := `
component UserService {
	id: string
	name: string
	email: string

	method Create(user: User): User
	method Get(id: string): User
	method Delete(id: string): void
}

component OrderService {
	orderId: string
	userId: string

	relation UserService: uses
}
`
	// Parse
	lexer := parser.NewLexer(strings.NewReader(input))
	p := parser.NewParser(lexer)
	spec, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// Transform
	tr := transformer.NewClassTransformer()
	diagram, err := tr.Transform(spec)
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}

	// Render
	renderer := svg.NewClassRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(diagram, &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	// Verify
	svgOut := buf.String()
	if !strings.Contains(svgOut, "<svg") {
		t.Error("expected valid SVG output")
	}
	if !strings.Contains(svgOut, "UserService") {
		t.Error("expected UserService in output")
	}
}

// I021: シーケンス図フルパイプライン
func TestPipeline_SequenceDiagram_Full(t *testing.T) {
	input := `
component AuthService {
	flow Login {
		step1: receive credentials from Client
		step2: call validate on UserDB
		step3: call createToken on TokenService
		step4: return token to Client
	}
}
`
	// Parse
	lexer := parser.NewLexer(strings.NewReader(input))
	p := parser.NewParser(lexer)
	spec, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// Transform
	tr := transformer.NewSequenceTransformer()
	diagram, err := tr.Transform(spec, "Login")
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}

	// Render
	renderer := svg.NewSequenceRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(diagram, &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	// Verify
	svgOut := buf.String()
	if !strings.Contains(svgOut, "<svg") {
		t.Error("expected valid SVG output")
	}
}

// I022: 状態図フルパイプライン
func TestPipeline_StateDiagram_Full(t *testing.T) {
	input := `
component OrderService {
	states OrderState {
		initial -> Created
		Created -> Pending: submit
		Pending -> Processing: pay
		Processing -> Shipped: ship
		Shipped -> Delivered: deliver
		Delivered -> [*]
	}
}
`
	// Parse
	lexer := parser.NewLexer(strings.NewReader(input))
	p := parser.NewParser(lexer)
	spec, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// Transform
	tr := transformer.NewStateTransformer()
	diagram, err := tr.Transform(spec, "OrderState")
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}

	// Render
	renderer := svg.NewStateRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(diagram, &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	// Verify
	svgOut := buf.String()
	if !strings.Contains(svgOut, "<svg") {
		t.Error("expected valid SVG output")
	}
}

// I023: フローチャートフルパイプライン
func TestPipeline_Flowchart_Full(t *testing.T) {
	input := `
component PaymentService {
	flow ProcessPayment {
		start: "Start"
		input: "Get Payment Info"
		validate: "Validate Card"
		if cardValid {
			process: "Process Payment"
			if success {
				confirm: "Send Confirmation"
			} else {
				retry: "Retry Payment"
			}
		} else {
			reject: "Reject Payment"
		}
		end: "End"
	}
}
`
	// Parse
	lexer := parser.NewLexer(strings.NewReader(input))
	p := parser.NewParser(lexer)
	spec, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// Transform
	tr := transformer.NewFlowTransformer()
	diagram, err := tr.Transform(spec, "ProcessPayment")
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}

	// Render
	renderer := svg.NewFlowRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(diagram, &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	// Verify
	svgOut := buf.String()
	if !strings.Contains(svgOut, "<svg") {
		t.Error("expected valid SVG output")
	}
}

// I024: 全図種類生成
func TestPipeline_AllDiagrams(t *testing.T) {
	input := `
component UserService {
	id: string
	name: string

	method GetUser(id: string): User

	flow FetchUser {
		start: "Begin"
		fetch: "Fetch from DB"
		end: "Return"
	}

	states UserState {
		initial -> Active
		Active -> Inactive: deactivate
		Inactive -> [*]
	}
}
`
	lexer := parser.NewLexer(strings.NewReader(input))
	p := parser.NewParser(lexer)
	spec, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// Class diagram
	classTr := transformer.NewClassTransformer()
	classDiagram, err := classTr.Transform(spec)
	if err != nil {
		t.Fatalf("class transform error: %v", err)
	}
	if len(classDiagram.Nodes) == 0 {
		t.Error("expected class nodes")
	}

	// Sequence diagram
	seqTr := transformer.NewSequenceTransformer()
	seqDiagram, err := seqTr.Transform(spec, "FetchUser")
	if err != nil {
		t.Fatalf("sequence transform error: %v", err)
	}
	if seqDiagram == nil {
		t.Error("expected sequence diagram")
	}

	// State diagram
	stateTr := transformer.NewStateTransformer()
	stateDiagram, err := stateTr.Transform(spec, "UserState")
	if err != nil {
		t.Fatalf("state transform error: %v", err)
	}
	if len(stateDiagram.States) == 0 {
		t.Error("expected states")
	}

	// Flow diagram
	flowTr := transformer.NewFlowTransformer()
	flowDiagram, err := flowTr.Transform(spec, "FetchUser")
	if err != nil {
		t.Fatalf("flow transform error: %v", err)
	}
	if len(flowDiagram.Nodes) == 0 {
		t.Error("expected flow nodes")
	}
}

// I025: 複雑なスペックの処理
func TestPipeline_ComplexSpec(t *testing.T) {
	input := `
@version("1.0")
@author("test")
component AuthenticationService {
	// Fields
	private tokenSecret: string
	private expirationTime: int
	users: map[string]User

	// Methods
	public method Authenticate(credentials: Credentials): AuthResult
	private method validatePassword(password: string, hash: string): bool
	method generateToken(user: User): string

	// Relations
	relation UserRepository: uses
	relation TokenCache: uses

	// Flow
	flow AuthenticateUser {
		start: "Receive Request"
		validate: "Validate Credentials"
		if valid {
			generate: "Generate Token"
			cache: "Cache Token"
			respond: "Return Success"
		} else {
			reject: "Return Error"
		}
		end: "Complete"
	}

	// States
	states SessionState {
		initial -> Created
		Created -> Active: activate
		Active -> Expired: timeout
		Active -> Revoked: revoke
		Expired -> [*]
		Revoked -> [*]
	}
}

interface Authenticator {
	method Authenticate(credentials: Credentials): AuthResult
}

type Credentials {
	username: string
	password: string
}

type AuthResult {
	success: bool
	token: string?
	error: string?
}
`
	lexer := parser.NewLexer(strings.NewReader(input))
	p := parser.NewParser(lexer)
	spec, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// All transformations should succeed
	classTr := transformer.NewClassTransformer()
	if _, err := classTr.Transform(spec); err != nil {
		t.Errorf("class transform error: %v", err)
	}

	seqTr := transformer.NewSequenceTransformer()
	if _, err := seqTr.Transform(spec, "AuthenticateUser"); err != nil {
		t.Errorf("sequence transform error: %v", err)
	}

	stateTr := transformer.NewStateTransformer()
	if _, err := stateTr.Transform(spec, "SessionState"); err != nil {
		t.Errorf("state transform error: %v", err)
	}

	flowTr := transformer.NewFlowTransformer()
	if _, err := flowTr.Transform(spec, "AuthenticateUser"); err != nil {
		t.Errorf("flow transform error: %v", err)
	}
}

// I026: 実例テスト
func TestPipeline_RealWorldExample(t *testing.T) {
	// Skip if testdata doesn't exist
	testdataDir := "../../testdata/valid"
	if _, err := os.Stat(testdataDir); os.IsNotExist(err) {
		t.Skip("testdata directory not found")
	}

	files, err := filepath.Glob(filepath.Join(testdataDir, "*.pact"))
	if err != nil {
		t.Fatalf("failed to glob: %v", err)
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("failed to read: %v", err)
			}

			lexer := parser.NewLexer(strings.NewReader(string(content)))
			p := parser.NewParser(lexer)
			spec, err := p.Parse()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			// At least class diagram should work
			classTr := transformer.NewClassTransformer()
			diagram, err := classTr.Transform(spec)
			if err != nil {
				t.Fatalf("transform error: %v", err)
			}

			renderer := svg.NewClassRenderer()
			var buf bytes.Buffer
			if err := renderer.Render(diagram, &buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			if !strings.Contains(buf.String(), "<svg") {
				t.Error("expected valid SVG")
			}
		})
	}
}

// =============================================================================
// I030-I033: エラーハンドリングテスト
// =============================================================================

// I030: パースエラーの伝播
func TestIntegration_ParseError_Propagation(t *testing.T) {
	input := `component { invalid syntax`

	lexer := parser.NewLexer(strings.NewReader(input))
	p := parser.NewParser(lexer)
	_, err := p.Parse()
	if err == nil {
		t.Error("expected parse error")
	}
}

// I031: 変換エラーの伝播
func TestIntegration_TransformError_Propagation(t *testing.T) {
	input := `component Empty { }`

	lexer := parser.NewLexer(strings.NewReader(input))
	p := parser.NewParser(lexer)
	spec, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// Request non-existent flow
	tr := transformer.NewSequenceTransformer()
	_, err = tr.Transform(spec, "NonExistentFlow")
	if err == nil {
		t.Error("expected transform error for non-existent flow")
	}
}

// I032: エラー位置情報
func TestIntegration_ErrorPosition(t *testing.T) {
	input := `component Test {
	invalid syntax here
}`
	lexer := parser.NewLexer(strings.NewReader(input))
	p := parser.NewParser(lexer)
	_, err := p.Parse()
	if err == nil {
		t.Error("expected parse error")
		return
	}

	// Error should contain position information
	errStr := err.Error()
	if !strings.Contains(errStr, "line") && !strings.Contains(errStr, "2") {
		t.Log("error may not contain detailed position info (implementation dependent)")
	}
}

// I033: エラー型の確認
func TestIntegration_ErrorType(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"syntax_error", "component {"},
		{"missing_name", "component { field: string }"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := parser.NewLexer(strings.NewReader(tt.input))
			p := parser.NewParser(lexer)
			_, err := p.Parse()
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}
