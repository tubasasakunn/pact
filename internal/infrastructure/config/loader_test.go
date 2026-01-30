package config

import (
	"os"
	"path/filepath"
	"testing"

	domainConfig "pact/internal/domain/config"
	"pact/internal/domain/errors"
)

// =============================================================================
// CL001-CL006: Load/Save
// =============================================================================

// CL001: ファイル不在時デフォルト
func TestLoader_Load_NotFound(t *testing.T) {
	loader := NewLoader()
	cfg, err := loader.Load("/nonexistent/path/.pactconfig")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// デフォルト値が返されることを確認
	if cfg.SourceRoot != "./src" {
		t.Errorf("expected default SourceRoot './src', got %q", cfg.SourceRoot)
	}
}

// CL002: 有効ファイル
func TestLoader_Load_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".pactconfig")

	content := `source_root: ./custom/src
output_dir: ./custom/out
diagrams:
  - class
  - sequence
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	loader := NewLoader()
	cfg, err := loader.Load(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.SourceRoot != "./custom/src" {
		t.Errorf("expected SourceRoot './custom/src', got %q", cfg.SourceRoot)
	}
	if cfg.OutputDir != "./custom/out" {
		t.Errorf("expected OutputDir './custom/out', got %q", cfg.OutputDir)
	}
	if len(cfg.Diagrams) != 2 {
		t.Errorf("expected 2 diagrams, got %d", len(cfg.Diagrams))
	}
}

// CL003: 不正ファイル
func TestLoader_Load_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".pactconfig")

	content := `invalid yaml: [[[`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	loader := NewLoader()
	_, err := loader.Load(configPath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*errors.ConfigError); !ok {
		t.Errorf("expected ConfigError, got %T", err)
	}
}

// CL004: 部分設定
func TestLoader_Load_Partial(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".pactconfig")

	// SourceRootのみ指定
	content := `source_root: ./my/src`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	loader := NewLoader()
	cfg, err := loader.Load(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 指定した値
	if cfg.SourceRoot != "./my/src" {
		t.Errorf("expected SourceRoot './my/src', got %q", cfg.SourceRoot)
	}
	// デフォルト値が維持される
	if cfg.OutputDir != "./diagrams" {
		t.Errorf("expected default OutputDir './diagrams', got %q", cfg.OutputDir)
	}
}

// CL005: 保存成功
func TestLoader_Save_Success(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".pactconfig")

	cfg := &domainConfig.Config{
		SourceRoot: "./src",
		OutputDir:  "./out",
		Diagrams:   []string{"class"},
	}

	loader := NewLoader()
	if err := loader.Save(configPath, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ファイルが作成されたことを確認
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("expected file to be created")
	}
}

// CL006: ラウンドトリップ
func TestLoader_Save_ReadBack(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".pactconfig")

	original := &domainConfig.Config{
		SourceRoot: "./custom/src",
		PactRoot:   "./custom/pact",
		OutputDir:  "./custom/out",
		Language:   "typescript",
		Diagrams:   []string{"class", "sequence"},
		Exclude:    []string{"vendor", "node_modules"},
	}

	loader := NewLoader()
	if err := loader.Save(configPath, original); err != nil {
		t.Fatalf("save error: %v", err)
	}

	loaded, err := loader.Load(configPath)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}

	if loaded.SourceRoot != original.SourceRoot {
		t.Errorf("SourceRoot mismatch: %q vs %q", loaded.SourceRoot, original.SourceRoot)
	}
	if loaded.PactRoot != original.PactRoot {
		t.Errorf("PactRoot mismatch: %q vs %q", loaded.PactRoot, original.PactRoot)
	}
	if loaded.OutputDir != original.OutputDir {
		t.Errorf("OutputDir mismatch: %q vs %q", loaded.OutputDir, original.OutputDir)
	}
	if loaded.Language != original.Language {
		t.Errorf("Language mismatch: %q vs %q", loaded.Language, original.Language)
	}
	if len(loaded.Diagrams) != len(original.Diagrams) {
		t.Errorf("Diagrams count mismatch: %d vs %d", len(loaded.Diagrams), len(original.Diagrams))
	}
	if len(loaded.Exclude) != len(original.Exclude) {
		t.Errorf("Exclude count mismatch: %d vs %d", len(loaded.Exclude), len(original.Exclude))
	}
}

// =============================================================================
// CL007-CL009: FindProjectRoot
// =============================================================================

// CL007: ルート発見
func TestLoader_FindProjectRoot_Found(t *testing.T) {
	tmpDir := t.TempDir()

	// プロジェクトルートに.pactconfigを作成
	configPath := filepath.Join(tmpDir, ".pactconfig")
	if err := os.WriteFile(configPath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	loader := NewLoader()
	root, err := loader.FindProjectRoot(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root != tmpDir {
		t.Errorf("expected root %q, got %q", tmpDir, root)
	}
}

// CL008: ルート未発見
func TestLoader_FindProjectRoot_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	// .pactconfigを作成しない

	loader := NewLoader()
	_, err := loader.FindProjectRoot(tmpDir)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*errors.ConfigError); !ok {
		t.Errorf("expected ConfigError, got %T", err)
	}
}

// CL009: 親ディレクトリ検索
func TestLoader_FindProjectRoot_Nested(t *testing.T) {
	tmpDir := t.TempDir()

	// ルートに.pactconfigを作成
	configPath := filepath.Join(tmpDir, ".pactconfig")
	if err := os.WriteFile(configPath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	// サブディレクトリを作成
	subDir := filepath.Join(tmpDir, "sub", "nested")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	loader := NewLoader()
	root, err := loader.FindProjectRoot(subDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root != tmpDir {
		t.Errorf("expected root %q, got %q", tmpDir, root)
	}
}
