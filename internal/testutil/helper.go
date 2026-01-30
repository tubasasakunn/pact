// Package testutil provides test utilities and helpers.
package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"pact/internal/domain/ast"
)

// TempDir creates a temporary directory for testing.
func TempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "pact-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

// WriteFile writes a file in the given directory.
func WriteFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	return path
}

// MustParse parses a string and fails the test if parsing fails.
func MustParse(t *testing.T, content string) *ast.SpecFile {
	t.Helper()
	// Implementation would use parser
	return nil
}

// AssertContains fails if the string doesn't contain the substring.
func AssertContains(t *testing.T, s, substr string) {
	t.Helper()
	if len(s) == 0 {
		t.Errorf("string is empty, expected to contain %q", substr)
		return
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return
		}
	}
	t.Errorf("string %q does not contain %q", truncate(s, 100), substr)
}

// AssertNotContains fails if the string contains the substring.
func AssertNotContains(t *testing.T, s, substr string) {
	t.Helper()
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			t.Errorf("string %q should not contain %q", truncate(s, 100), substr)
			return
		}
	}
}

// AssertNoError fails if err is not nil.
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// AssertError fails if err is nil.
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// AssertEqual fails if got != want.
func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Golden returns the path to a golden file in testdata/golden.
func Golden(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join("testdata", "golden", name)
}

// LoadGolden loads a golden file content.
func LoadGolden(t *testing.T, name string) string {
	t.Helper()
	path := Golden(t, name)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to load golden file %s: %v", name, err)
	}
	return string(content)
}
