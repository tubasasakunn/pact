package pact

import (
	"bytes"
	"strings"
	"testing"
)

// =============================================================================
// A001-A015: 公開APIテスト
// =============================================================================

// A001: クライアント初期化
func TestAPI_New(t *testing.T) {
	client := New()
	if client == nil {
		t.Error("expected non-nil client")
	}
}

// A002: ファイルパース
func TestAPI_ParseFile(t *testing.T) {
	t.Skip("requires test file setup")
}

// A003: 存在しないファイルパース
func TestAPI_ParseFile_NotFound(t *testing.T) {
	client := New()
	_, err := client.ParseFile("/nonexistent/path.pact")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

// A004: 文字列パース
func TestAPI_ParseString(t *testing.T) {
	client := New()
	spec, err := client.ParseString(`component User { type Data { id: string } }`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if spec == nil {
		t.Error("expected non-nil spec")
	}
}

// A005: 不正な文字列パース
func TestAPI_ParseString_Invalid(t *testing.T) {
	client := New()
	_, err := client.ParseString(`component { invalid }`)
	if err == nil {
		t.Error("expected parse error")
	}
}

// A006: クラス図変換
func TestAPI_ToClassDiagram(t *testing.T) {
	client := New()
	spec, err := client.ParseString(`
component User {
	type UserData {
		id: string
		name: string
	}
}
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	diagram, err := client.ToClassDiagram(spec)
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}
	if diagram == nil {
		t.Error("expected non-nil diagram")
	}
}

// A007: シーケンス図変換
func TestAPI_ToSequenceDiagram(t *testing.T) {
	client := New()
	spec, err := client.ParseString(`
component Auth {
	depends on Client

	flow Login {
		result = Client.request(data)
		return result
	}
}
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	diagram, err := client.ToSequenceDiagram(spec, "Login")
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}
	if diagram == nil {
		t.Error("expected non-nil diagram")
	}
}

// A008: 存在しないフローのシーケンス図変換
func TestAPI_ToSequenceDiagram_NotFound(t *testing.T) {
	client := New()
	spec, err := client.ParseString(`component Empty { }`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	_, err = client.ToSequenceDiagram(spec, "NonExistent")
	if err == nil {
		t.Error("expected error for non-existent flow")
	}
}

// A009: 状態図変換
func TestAPI_ToStateDiagram(t *testing.T) {
	client := New()
	spec, err := client.ParseString(`
component Order {
	states OrderState {
		initial Pending
		final Complete

		state Pending { }
		state Complete { }

		Pending -> Complete on finish
	}
}
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	diagram, err := client.ToStateDiagram(spec, "OrderState")
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}
	if diagram == nil {
		t.Error("expected non-nil diagram")
	}
}

// A010: 存在しないステートの状態図変換
func TestAPI_ToStateDiagram_NotFound(t *testing.T) {
	client := New()
	spec, err := client.ParseString(`component Empty { }`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	_, err = client.ToStateDiagram(spec, "NonExistent")
	if err == nil {
		t.Error("expected error for non-existent states")
	}
}

// A011: フローチャート変換
func TestAPI_ToFlowchart(t *testing.T) {
	client := New()
	spec, err := client.ParseString(`
component Process {
	flow Main {
		result = self.doWork()
		return result
	}
}
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	diagram, err := client.ToFlowchart(spec, "Main")
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}
	if diagram == nil {
		t.Error("expected non-nil diagram")
	}
}

// A012: クラス図レンダー
func TestAPI_RenderClassDiagram(t *testing.T) {
	client := New()
	spec, err := client.ParseString(`component User { type Data { id: string } }`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	diagram, err := client.ToClassDiagram(spec)
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}

	var buf bytes.Buffer
	if err := client.RenderClassDiagram(diagram, &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	if !strings.Contains(buf.String(), "<svg") {
		t.Error("expected valid SVG output")
	}
}

// A013: シーケンス図レンダー
func TestAPI_RenderSequenceDiagram(t *testing.T) {
	client := New()
	spec, err := client.ParseString(`
component Auth {
	depends on Client

	flow Login {
		result = Client.request(data)
	}
}
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	diagram, err := client.ToSequenceDiagram(spec, "Login")
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}

	var buf bytes.Buffer
	if err := client.RenderSequenceDiagram(diagram, &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	if !strings.Contains(buf.String(), "<svg") {
		t.Error("expected valid SVG output")
	}
}

// A014: 状態図レンダー
func TestAPI_RenderStateDiagram(t *testing.T) {
	client := New()
	spec, err := client.ParseString(`
component Order {
	states State {
		initial Active
		final Done

		state Active { }
		state Done { }

		Active -> Done on finish
	}
}
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	diagram, err := client.ToStateDiagram(spec, "State")
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}

	var buf bytes.Buffer
	if err := client.RenderStateDiagram(diagram, &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	if !strings.Contains(buf.String(), "<svg") {
		t.Error("expected valid SVG output")
	}
}

// A015: フローチャートレンダー
func TestAPI_RenderFlowchart(t *testing.T) {
	client := New()
	spec, err := client.ParseString(`
component Process {
	flow Main {
		result = self.doWork()
		return result
	}
}
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	diagram, err := client.ToFlowchart(spec, "Main")
	if err != nil {
		t.Fatalf("transform error: %v", err)
	}

	var buf bytes.Buffer
	if err := client.RenderFlowchart(diagram, &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	if !strings.Contains(buf.String(), "<svg") {
		t.Error("expected valid SVG output")
	}
}
