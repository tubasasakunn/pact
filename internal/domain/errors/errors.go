package errors

import (
	"fmt"
	"pact/internal/domain/ast"
)

// ParseError は構文解析エラーを表す
type ParseError struct {
	Pos     ast.Position
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Pos.String(), e.Message)
}

// SemanticError は意味解析エラーを表す
type SemanticError struct {
	Pos     ast.Position
	Message string
}

func (e *SemanticError) Error() string {
	return fmt.Sprintf("%s: %s", e.Pos.String(), e.Message)
}

// ImportError はインポートエラーを表す
type ImportError struct {
	Pos     ast.Position
	Path    string
	Message string
	Cause   error
}

func (e *ImportError) Error() string {
	msg := fmt.Sprintf("%s: import %q: %s", e.Pos.String(), e.Path, e.Message)
	if e.Cause != nil {
		msg += fmt.Sprintf(" (caused by: %v)", e.Cause)
	}
	return msg
}

func (e *ImportError) Unwrap() error {
	return e.Cause
}

// CycleError はサイクル検出エラーを表す
type CycleError struct {
	Cycle []string
}

func (e *CycleError) Error() string {
	return fmt.Sprintf("import cycle detected: %v", e.Cycle)
}

// TransformError は変換エラーを表す
type TransformError struct {
	Source  string
	Target  string
	Message string
}

func (e *TransformError) Error() string {
	return fmt.Sprintf("transform %s to %s: %s", e.Source, e.Target, e.Message)
}

// ConfigError は設定エラーを表す
type ConfigError struct {
	Path    string
	Message string
}

func (e *ConfigError) Error() string {
	if e.Path == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}
