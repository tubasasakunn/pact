package config

import "testing"

// =============================================================================
// DC001-DC005: デフォルト設定
// =============================================================================

// DC001: デフォルト設定
func TestConfig_Default(t *testing.T) {
	cfg := Default()
	if cfg == nil {
		t.Fatal("expected config, got nil")
	}
}

// DC002: デフォルトSourceRoot
func TestConfig_Default_SourceRoot(t *testing.T) {
	cfg := Default()
	if cfg.SourceRoot != "./src" {
		t.Errorf("expected './src', got %q", cfg.SourceRoot)
	}
}

// DC003: デフォルトPactRoot
func TestConfig_Default_PactRoot(t *testing.T) {
	cfg := Default()
	if cfg.PactRoot != "./.pact" {
		t.Errorf("expected './.pact', got %q", cfg.PactRoot)
	}
}

// DC004: デフォルト出力先
func TestConfig_Default_OutputDir(t *testing.T) {
	cfg := Default()
	if cfg.OutputDir != "./diagrams" {
		t.Errorf("expected './diagrams', got %q", cfg.OutputDir)
	}
}

// DC005: デフォルト図種類
func TestConfig_Default_Diagrams(t *testing.T) {
	cfg := Default()
	expected := []string{"class", "sequence", "state", "flow"}
	if len(cfg.Diagrams) != len(expected) {
		t.Fatalf("expected %d diagrams, got %d", len(expected), len(cfg.Diagrams))
	}
	for i, exp := range expected {
		if cfg.Diagrams[i] != exp {
			t.Errorf("expected diagrams[%d]=%q, got %q", i, exp, cfg.Diagrams[i])
		}
	}
}

// =============================================================================
// DC006-DC008: DiagramEnabled
// =============================================================================

// DC006: 有効な図
func TestConfig_DiagramEnabled_Exists(t *testing.T) {
	cfg := Default()
	if !cfg.DiagramEnabled("class") {
		t.Error("expected class diagram to be enabled")
	}
}

// DC007: 無効な図
func TestConfig_DiagramEnabled_NotExists(t *testing.T) {
	cfg := Default()
	if cfg.DiagramEnabled("foo") {
		t.Error("expected foo diagram to be disabled")
	}
}

// DC008: all指定
func TestConfig_DiagramEnabled_All(t *testing.T) {
	cfg := &Config{
		Diagrams: []string{"all"},
	}
	if !cfg.DiagramEnabled("class") {
		t.Error("expected class to be enabled with 'all'")
	}
	if !cfg.DiagramEnabled("sequence") {
		t.Error("expected sequence to be enabled with 'all'")
	}
	if !cfg.DiagramEnabled("foo") {
		t.Error("expected any diagram to be enabled with 'all'")
	}
}

// =============================================================================
// DC009: IsExcluded
// =============================================================================

// DC009: 除外判定
func TestConfig_IsExcluded(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		path     string
		expected bool
	}{
		{
			name:     "exact match",
			patterns: []string{"vendor"},
			path:     "vendor",
			expected: true,
		},
		{
			name:     "pattern match",
			patterns: []string{"*.tmp"},
			path:     "file.tmp",
			expected: true,
		},
		{
			name:     "no match",
			patterns: []string{"vendor"},
			path:     "src",
			expected: false,
		},
		{
			name:     "basename match",
			patterns: []string{"node_modules"},
			path:     "project/node_modules",
			expected: true,
		},
		{
			name:     "empty patterns",
			patterns: []string{},
			path:     "anything",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Exclude: tt.patterns}
			result := cfg.IsExcluded(tt.path)
			if result != tt.expected {
				t.Errorf("IsExcluded(%q) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}
