package integration

import (
	"os"
	"path/filepath"
	"testing"

	"pact/internal/infrastructure/resolver"
)

// =============================================================================
// I006-I010: インポート解決テスト
// =============================================================================

func setupTestFiles(t *testing.T, files map[string]string) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "pact-import-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write file %s: %v", name, err)
		}
	}
	return dir
}

// I006: 単一インポート解決
func TestIntegration_ImportResolution_Single(t *testing.T) {
	files := map[string]string{
		"a.pact": `import "./b.pact"
component A { }`,
		"b.pact": `component B { }`,
	}
	dir := setupTestFiles(t, files)

	r := resolver.NewImportResolver()
	order, err := r.Resolve(filepath.Join(dir, "a.pact"))
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}

	// B should come before A
	if len(order) != 2 {
		t.Fatalf("expected 2 files, got %d", len(order))
	}
	if filepath.Base(order[0]) != "b.pact" {
		t.Error("expected b.pact first")
	}
	if filepath.Base(order[1]) != "a.pact" {
		t.Error("expected a.pact second")
	}
}

// I007: チェーンインポート解決
func TestIntegration_ImportResolution_Chain(t *testing.T) {
	files := map[string]string{
		"a.pact": `import "./b.pact"
component A { }`,
		"b.pact": `import "./c.pact"
component B { }`,
		"c.pact": `component C { }`,
	}
	dir := setupTestFiles(t, files)

	r := resolver.NewImportResolver()
	order, err := r.Resolve(filepath.Join(dir, "a.pact"))
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}

	// C -> B -> A
	if len(order) != 3 {
		t.Fatalf("expected 3 files, got %d", len(order))
	}
	if filepath.Base(order[0]) != "c.pact" {
		t.Error("expected c.pact first")
	}
	if filepath.Base(order[1]) != "b.pact" {
		t.Error("expected b.pact second")
	}
	if filepath.Base(order[2]) != "a.pact" {
		t.Error("expected a.pact third")
	}
}

// I008: ダイヤモンドインポート解決
func TestIntegration_ImportResolution_Diamond(t *testing.T) {
	files := map[string]string{
		"a.pact": `import "./b.pact"
import "./c.pact"
component A { }`,
		"b.pact": `import "./d.pact"
component B { }`,
		"c.pact": `import "./d.pact"
component C { }`,
		"d.pact": `component D { }`,
	}
	dir := setupTestFiles(t, files)

	r := resolver.NewImportResolver()
	order, err := r.Resolve(filepath.Join(dir, "a.pact"))
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}

	// D should come first, A should come last
	if len(order) != 4 {
		t.Fatalf("expected 4 files, got %d", len(order))
	}
	if filepath.Base(order[0]) != "d.pact" {
		t.Error("expected d.pact first")
	}
	if filepath.Base(order[len(order)-1]) != "a.pact" {
		t.Error("expected a.pact last")
	}
}

// I009: 循環インポートエラー
func TestIntegration_ImportResolution_Cycle(t *testing.T) {
	files := map[string]string{
		"a.pact": `import "./b.pact"
component A { }`,
		"b.pact": `import "./a.pact"
component B { }`,
	}
	dir := setupTestFiles(t, files)

	r := resolver.NewImportResolver()
	_, err := r.Resolve(filepath.Join(dir, "a.pact"))
	if err == nil {
		t.Error("expected cycle error")
	}
	if !resolver.IsCycleError(err) {
		t.Errorf("expected CycleError, got %T", err)
	}
}

// I010: インポートファイル未発見エラー
func TestIntegration_ImportResolution_NotFound(t *testing.T) {
	files := map[string]string{
		"a.pact": `import "./missing.pact"
component A { }`,
	}
	dir := setupTestFiles(t, files)

	r := resolver.NewImportResolver()
	_, err := r.Resolve(filepath.Join(dir, "a.pact"))
	if err == nil {
		t.Error("expected import error")
	}
}
