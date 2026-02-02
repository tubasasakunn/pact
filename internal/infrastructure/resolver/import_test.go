package resolver

import (
	"testing"

	"pact/internal/domain/ast"
	"pact/internal/domain/errors"
)

// MockParser はテスト用のパーサー
type MockParser struct {
	files map[string]*ast.SpecFile
	err   error
}

func (m *MockParser) ParseFile(path string) (*ast.SpecFile, error) {
	if m.err != nil {
		return nil, m.err
	}
	if file, ok := m.files[path]; ok {
		return file, nil
	}
	return nil, &errors.ImportError{Path: path, Message: "file not found"}
}

// =============================================================================
// IR001-IR010: Import Resolver Tests
// =============================================================================

// IR001: インポートなし
func TestResolver_NoImports(t *testing.T) {
	file := &ast.SpecFile{
		Path:    "a.pact",
		Imports: []ast.ImportDecl{},
	}

	parser := &MockParser{files: map[string]*ast.SpecFile{}}
	resolver := NewResolver(parser)

	result, err := resolver.ResolveFile(file, ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 file, got %d", len(result))
	}
}

// IR002: 単一インポート
func TestResolver_SingleImport(t *testing.T) {
	fileA := &ast.SpecFile{
		Path: "a.pact",
		Imports: []ast.ImportDecl{
			{Path: "./b.pact"},
		},
	}
	fileB := &ast.SpecFile{
		Path:    "b.pact",
		Imports: []ast.ImportDecl{},
	}

	parser := &MockParser{
		files: map[string]*ast.SpecFile{
			"b.pact": fileB,
		},
	}
	resolver := NewResolver(parser)

	result, err := resolver.ResolveFile(fileA, ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// B が A の前に来る
	if len(result) != 2 {
		t.Fatalf("expected 2 files, got %d", len(result))
	}
	if result[0].Path != "b.pact" {
		t.Errorf("expected first file to be 'b.pact', got %q", result[0].Path)
	}
	if result[1].Path != "a.pact" {
		t.Errorf("expected second file to be 'a.pact', got %q", result[1].Path)
	}
}

// IR003: 複数インポート
func TestResolver_MultipleImports(t *testing.T) {
	fileA := &ast.SpecFile{
		Path: "a.pact",
		Imports: []ast.ImportDecl{
			{Path: "./b.pact"},
			{Path: "./c.pact"},
		},
	}
	fileB := &ast.SpecFile{Path: "b.pact", Imports: []ast.ImportDecl{}}
	fileC := &ast.SpecFile{Path: "c.pact", Imports: []ast.ImportDecl{}}

	parser := &MockParser{
		files: map[string]*ast.SpecFile{
			"b.pact": fileB,
			"c.pact": fileC,
		},
	}
	resolver := NewResolver(parser)

	result, err := resolver.ResolveFile(fileA, ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 files, got %d", len(result))
	}
	// A は最後
	if result[len(result)-1].Path != "a.pact" {
		t.Errorf("expected last file to be 'a.pact'")
	}
}

// IR004: 推移的インポート A→B→C
func TestResolver_TransitiveImports(t *testing.T) {
	fileA := &ast.SpecFile{
		Path:    "a.pact",
		Imports: []ast.ImportDecl{{Path: "./b.pact"}},
	}
	fileB := &ast.SpecFile{
		Path:    "b.pact",
		Imports: []ast.ImportDecl{{Path: "./c.pact"}},
	}
	fileC := &ast.SpecFile{
		Path:    "c.pact",
		Imports: []ast.ImportDecl{},
	}

	parser := &MockParser{
		files: map[string]*ast.SpecFile{
			"b.pact": fileB,
			"c.pact": fileC,
		},
	}
	resolver := NewResolver(parser)

	result, err := resolver.ResolveFile(fileA, ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// C→B→A 順
	if len(result) != 3 {
		t.Fatalf("expected 3 files, got %d", len(result))
	}
	if result[0].Path != "c.pact" {
		t.Errorf("expected first file to be 'c.pact', got %q", result[0].Path)
	}
	if result[1].Path != "b.pact" {
		t.Errorf("expected second file to be 'b.pact', got %q", result[1].Path)
	}
	if result[2].Path != "a.pact" {
		t.Errorf("expected third file to be 'a.pact', got %q", result[2].Path)
	}
}

// IR005: サイクル検出 A→B→A
func TestResolver_CycleDetection(t *testing.T) {
	fileA := &ast.SpecFile{
		Path:    "a.pact",
		Imports: []ast.ImportDecl{{Path: "./b.pact"}},
	}
	fileB := &ast.SpecFile{
		Path:    "b.pact",
		Imports: []ast.ImportDecl{{Path: "./a.pact"}},
	}

	parser := &MockParser{
		files: map[string]*ast.SpecFile{
			"b.pact": fileB,
			"a.pact": fileA,
		},
	}
	resolver := NewResolver(parser)

	_, err := resolver.ResolveFile(fileA, ".")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*errors.CycleError); !ok {
		t.Errorf("expected CycleError, got %T", err)
	}
}

// IR006: 自己参照
func TestResolver_SelfImport(t *testing.T) {
	fileA := &ast.SpecFile{
		Path:    "a.pact",
		Imports: []ast.ImportDecl{{Path: "./a.pact"}},
	}

	parser := &MockParser{
		files: map[string]*ast.SpecFile{
			"a.pact": fileA,
		},
	}
	resolver := NewResolver(parser)

	_, err := resolver.ResolveFile(fileA, ".")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*errors.CycleError); !ok {
		t.Errorf("expected CycleError, got %T", err)
	}
}

// IR007: ダイヤモンド依存 A→B,C B,C→D
func TestResolver_DiamondImports(t *testing.T) {
	fileA := &ast.SpecFile{
		Path: "a.pact",
		Imports: []ast.ImportDecl{
			{Path: "./b.pact"},
			{Path: "./c.pact"},
		},
	}
	fileB := &ast.SpecFile{
		Path:    "b.pact",
		Imports: []ast.ImportDecl{{Path: "./d.pact"}},
	}
	fileC := &ast.SpecFile{
		Path:    "c.pact",
		Imports: []ast.ImportDecl{{Path: "./d.pact"}},
	}
	fileD := &ast.SpecFile{
		Path:    "d.pact",
		Imports: []ast.ImportDecl{},
	}

	parser := &MockParser{
		files: map[string]*ast.SpecFile{
			"b.pact": fileB,
			"c.pact": fileC,
			"d.pact": fileD,
		},
	}
	resolver := NewResolver(parser)

	result, err := resolver.ResolveFile(fileA, ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// D は一度だけ
	dCount := 0
	for _, f := range result {
		if f.Path == "d.pact" {
			dCount++
		}
	}
	if dCount != 1 {
		t.Errorf("expected D to appear once, got %d", dCount)
	}

	// D は最初に来る
	if result[0].Path != "d.pact" {
		t.Errorf("expected first file to be 'd.pact', got %q", result[0].Path)
	}
}

// IR008: 相対パス
func TestResolver_RelativePath(t *testing.T) {
	fileA := &ast.SpecFile{
		Path:    "a.pact",
		Imports: []ast.ImportDecl{{Path: "./sub/b.pact"}},
	}
	fileB := &ast.SpecFile{
		Path:    "sub/b.pact",
		Imports: []ast.ImportDecl{},
	}

	parser := &MockParser{
		files: map[string]*ast.SpecFile{
			"sub/b.pact": fileB,
		},
	}
	resolver := NewResolver(parser)

	result, err := resolver.ResolveFile(fileA, ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 files, got %d", len(result))
	}
}

// IR009: インポート未発見
func TestResolver_ImportNotFound(t *testing.T) {
	fileA := &ast.SpecFile{
		Path:    "a.pact",
		Imports: []ast.ImportDecl{{Path: "./nonexistent.pact"}},
	}

	parser := &MockParser{files: map[string]*ast.SpecFile{}}
	resolver := NewResolver(parser)

	_, err := resolver.ResolveFile(fileA, ".")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*errors.ImportError); !ok {
		t.Errorf("expected ImportError, got %T", err)
	}
}

// IR010: インポート先エラー
func TestResolver_ImportParseError(t *testing.T) {
	fileA := &ast.SpecFile{
		Path:    "a.pact",
		Imports: []ast.ImportDecl{{Path: "./b.pact"}},
	}

	parser := &MockParser{
		files: map[string]*ast.SpecFile{},
		err:   &errors.ParseError{Message: "syntax error"},
	}
	resolver := NewResolver(parser)

	_, err := resolver.ResolveFile(fileA, ".")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	impErr, ok := err.(*errors.ImportError)
	if !ok {
		t.Fatalf("expected ImportError, got %T", err)
	}
	if impErr.Cause == nil {
		t.Error("expected Cause to be set")
	}
}
