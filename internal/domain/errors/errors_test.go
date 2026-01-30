package errors

import (
	"strings"
	"testing"

	"pact/internal/domain/ast"
)

// =============================================================================
// DE001-DE006: Error Tests
// =============================================================================

// DE001: ParseError.Error()
func TestParseError_Error(t *testing.T) {
	err := &ParseError{
		Pos:     ast.Position{File: "test.pact", Line: 10, Column: 5},
		Message: "unexpected token",
	}
	result := err.Error()

	if !strings.Contains(result, "test.pact") {
		t.Errorf("expected error to contain file name, got %q", result)
	}
	if !strings.Contains(result, "10") {
		t.Errorf("expected error to contain line number, got %q", result)
	}
	if !strings.Contains(result, "5") {
		t.Errorf("expected error to contain column number, got %q", result)
	}
	if !strings.Contains(result, "unexpected token") {
		t.Errorf("expected error to contain message, got %q", result)
	}
}

// DE002: SemanticError.Error()
func TestSemanticError_Error(t *testing.T) {
	err := &SemanticError{
		Pos:     ast.Position{File: "test.pact", Line: 20, Column: 1},
		Message: "undefined type",
	}
	result := err.Error()

	if !strings.Contains(result, "test.pact") {
		t.Errorf("expected error to contain file name, got %q", result)
	}
	if !strings.Contains(result, "20") {
		t.Errorf("expected error to contain line number, got %q", result)
	}
	if !strings.Contains(result, "undefined type") {
		t.Errorf("expected error to contain message, got %q", result)
	}
}

// DE003: ImportError.Error()
func TestImportError_Error(t *testing.T) {
	err := &ImportError{
		Pos:     ast.Position{File: "main.pact", Line: 1, Column: 1},
		Path:    "./missing.pact",
		Message: "file not found",
	}
	result := err.Error()

	if !strings.Contains(result, "main.pact") {
		t.Errorf("expected error to contain file name, got %q", result)
	}
	if !strings.Contains(result, "./missing.pact") {
		t.Errorf("expected error to contain import path, got %q", result)
	}
	if !strings.Contains(result, "file not found") {
		t.Errorf("expected error to contain message, got %q", result)
	}
}

// DE003b: ImportError with Cause
func TestImportError_Error_WithCause(t *testing.T) {
	cause := &ParseError{Message: "syntax error"}
	err := &ImportError{
		Pos:     ast.Position{File: "main.pact", Line: 1, Column: 1},
		Path:    "./broken.pact",
		Message: "parse failed",
		Cause:   cause,
	}
	result := err.Error()

	if !strings.Contains(result, "syntax error") {
		t.Errorf("expected error to contain cause, got %q", result)
	}
}

// DE003c: ImportError.Unwrap()
func TestImportError_Unwrap(t *testing.T) {
	cause := &ParseError{Message: "syntax error"}
	err := &ImportError{
		Path:    "./broken.pact",
		Message: "parse failed",
		Cause:   cause,
	}

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("expected Unwrap to return cause")
	}
}

// DE004: CycleError.Error()
func TestCycleError_Error(t *testing.T) {
	err := &CycleError{
		Cycle: []string{"a.pact", "b.pact", "a.pact"},
	}
	result := err.Error()

	if !strings.Contains(result, "cycle") {
		t.Errorf("expected error to mention cycle, got %q", result)
	}
	if !strings.Contains(result, "a.pact") {
		t.Errorf("expected error to contain cycle files, got %q", result)
	}
	if !strings.Contains(result, "b.pact") {
		t.Errorf("expected error to contain cycle files, got %q", result)
	}
}

// DE005: TransformError.Error()
func TestTransformError_Error(t *testing.T) {
	err := &TransformError{
		Source:  "AST",
		Target:  "ClassDiagram",
		Message: "invalid component",
	}
	result := err.Error()

	if !strings.Contains(result, "AST") {
		t.Errorf("expected error to contain source, got %q", result)
	}
	if !strings.Contains(result, "ClassDiagram") {
		t.Errorf("expected error to contain target, got %q", result)
	}
	if !strings.Contains(result, "invalid component") {
		t.Errorf("expected error to contain message, got %q", result)
	}
}

// DE006: ConfigError.Error()
func TestConfigError_Error(t *testing.T) {
	err := &ConfigError{
		Path:    ".pactconfig",
		Message: "invalid yaml",
	}
	result := err.Error()

	if !strings.Contains(result, ".pactconfig") {
		t.Errorf("expected error to contain path, got %q", result)
	}
	if !strings.Contains(result, "invalid yaml") {
		t.Errorf("expected error to contain message, got %q", result)
	}
}

// DE006b: ConfigError without path
func TestConfigError_Error_NoPath(t *testing.T) {
	err := &ConfigError{
		Message: "project root not found",
	}
	result := err.Error()

	if result != "project root not found" {
		t.Errorf("expected just message, got %q", result)
	}
}
