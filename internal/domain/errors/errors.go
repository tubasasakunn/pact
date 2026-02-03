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

// MultiError は複数のエラーを表す
type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return "no errors"
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	msg := fmt.Sprintf("%d errors:\n", len(e.Errors))
	for i, err := range e.Errors {
		msg += fmt.Sprintf("  %d. %s\n", i+1, err.Error())
	}
	return msg
}

// Add はエラーを追加する
func (e *MultiError) Add(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

// HasErrors はエラーがあるかどうかを返す
func (e *MultiError) HasErrors() bool {
	return len(e.Errors) > 0
}

// ErrorOrNil はエラーがあればMultiErrorを、なければnilを返す
func (e *MultiError) ErrorOrNil() error {
	if !e.HasErrors() {
		return nil
	}
	return e
}

// ValidationError は検証エラーを表す（重複、未定義参照など）
type ValidationError struct {
	Pos     ast.Position
	Type    string // "duplicate", "undefined", "invalid" など
	Name    string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s '%s': %s", e.Pos.String(), e.Type, e.Name, e.Message)
}

// Warning は警告を表す（エラーではないが注意が必要な事項）
type Warning struct {
	Pos     ast.Position
	Code    string // "unused-import", "unused-type", "deprecated", etc.
	Message string
}

func (w *Warning) Error() string {
	return fmt.Sprintf("%s: warning[%s]: %s", w.Pos.String(), w.Code, w.Message)
}

// WarningList は警告のリスト
type WarningList struct {
	Warnings []*Warning
}

// Add は警告を追加する
func (wl *WarningList) Add(w *Warning) {
	if w != nil {
		wl.Warnings = append(wl.Warnings, w)
	}
}

// HasWarnings は警告があるかどうかを返す
func (wl *WarningList) HasWarnings() bool {
	return len(wl.Warnings) > 0
}

// String は警告リストの文字列表現を返す
func (wl *WarningList) String() string {
	if len(wl.Warnings) == 0 {
		return ""
	}
	msg := fmt.Sprintf("%d warning(s):\n", len(wl.Warnings))
	for i, w := range wl.Warnings {
		msg += fmt.Sprintf("  %d. %s\n", i+1, w.Error())
	}
	return msg
}
