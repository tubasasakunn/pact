# Pact DSL - Remaining Issues

This document tracks known issues and areas for improvement in the Pact DSL project.

## Critical Priority

### ~~C-001: Parser Error Recovery~~ (FIXED)
- **Location**: `internal/infrastructure/parser/parser.go`
- **Description**: ~~Parser stops on first error without recovering or collecting multiple errors~~
- **Status**: Fixed - Parser now collects multiple errors and continues parsing with error recovery
- **Fixed in**: `internal/infrastructure/parser/parser.go` - Added `addError()`, `synchronize()`, and `getErrors()` methods

### ~~C-002: Duplicate Validation Missing~~ (FIXED)
- **Location**: `internal/application/validator/validator.go` (NEW)
- **Description**: ~~No validation for duplicate field names, method names, parameter names, etc.~~
- **Status**: Fixed - Added comprehensive duplicate validation
- **Fixed in**: New validator package with `Validate()` method checking duplicates for:
  - Field names in types
  - Method names in interfaces
  - Parameter names in methods
  - Type names, flow names, states names
  - Relation declarations

### ~~C-003: Circular Dependency Detection~~ (ALREADY IMPLEMENTED)
- **Location**: `internal/infrastructure/resolver/resolver.go`
- **Description**: ~~No cycle detection for import dependencies~~
- **Status**: Already implemented - `CycleError` is returned when circular imports are detected
- **Note**: The resolver already had cycle detection via `inProgress` map tracking

## High Priority

### H-001: Text Wrapping Not Implemented
- **Location**: `internal/infrastructure/renderer/svg/renderer.go`
- **Description**: Long text (field names, method signatures, labels) not wrapped
- **Impact**: Text overflows node boundaries or creates excessively wide nodes

### H-002: Collision Detection Incomplete
- **Location**: `internal/infrastructure/renderer/svg/renderer.go`
- **Description**: Missing collision detection for:
  - Note-to-node overlaps
  - Edge label-to-edge label overlaps
  - Participant name overlaps in sequence diagrams
- **Impact**: Visual overlapping in complex diagrams

### ~~H-003: Undefined Reference Validation~~ (FIXED)
- **Location**: `internal/application/validator/validator.go` (NEW)
- **Description**: ~~References to undefined types/components/states accepted without validation~~
- **Status**: Fixed - Added `ValidateReferences()` method
- **Fixed in**: Validator now checks:
  - Relation targets reference defined components/types
  - Field types are defined or builtin
  - Method parameter and return types are defined
  - State transitions reference defined states

### ~~H-004: Expression Nesting Limit~~ (FIXED)
- **Location**: `internal/application/validator/validator.go` (NEW)
- **Description**: ~~No limit on expression nesting depth~~
- **Status**: Fixed - Added `ValidateExpressionDepth()` method with `MaxNestingDepth = 50`
- **Fixed in**: Validator now checks nesting depth for expressions and control flow statements

### ~~H-005: Dead Code Detection~~ (FIXED)
- **Location**: `internal/application/validator/validator.go` (NEW)
- **Description**: ~~Steps after unconditional return/throw not flagged as dead code~~
- **Status**: Fixed - Added `ValidateDeadCode()` method
- **Fixed in**: Validator detects unreachable code after return/throw statements

## Medium Priority

### ~~M-001: String Escape Sequences Limited~~ (FIXED)
- **Location**: `internal/infrastructure/parser/lexer.go`
- **Description**: ~~Only supports `\n, \t, \r, ", \\`. Missing `\u` Unicode and `\x` hex escapes~~
- **Status**: Fixed - Added `\uXXXX` Unicode escapes and `\xXX` hex escapes

### ~~M-002: Scientific Notation Not Supported~~ (FIXED)
- **Location**: `internal/infrastructure/parser/lexer.go`
- **Description**: ~~Cannot parse `1.5e-10` or `2E+5`~~
- **Status**: Fixed - Lexer now supports scientific notation (e.g., `1.5e-10`, `2E+5`)

### M-003: Duration Unit Validation
- **Location**: `internal/infrastructure/parser/lexer.go`
- **Description**: Invalid duration units (e.g., `100xyz`) not rejected
- **Impact**: Silent acceptance of invalid durations
- **Note**: Partially improved - only valid units (ms, s, m, h, d) are now recognized

### ~~M-004: Unclosed Block Comment~~ (FIXED)
- **Location**: `internal/infrastructure/parser/lexer.go`
- **Description**: ~~Unclosed `/* comment` at EOF doesn't error~~
- **Status**: Fixed - Lexer now reports error for unclosed block comments

### ~~M-005: Type Modifier Chaining~~ (FIXED)
- **Location**: `internal/infrastructure/parser/parser.go`
- **Description**: ~~`Type??[]?` parses successfully but is logically invalid~~
- **Status**: Fixed - Parser now validates type modifier combinations and rejects invalid chains

### M-006: Barycenter Iteration Limit
- **Location**: `internal/infrastructure/renderer/svg/renderer.go`
- **Description**: Only 4 iterations hardcoded, may not converge on complex graphs
- **Impact**: Suboptimal edge crossing reduction

### M-007: Canvas Size for Notes
- **Location**: `internal/infrastructure/renderer/svg/renderer.go`
- **Description**: Canvas height doesn't account for note positions
- **Impact**: Notes may be cut off at diagram edges

### M-008: Sequence Diagram Fixed Width
- **Location**: `internal/infrastructure/renderer/svg/renderer.go`
- **Description**: Frame width hardcoded to 700, doesn't scale with participants
- **Impact**: Crowded or sparse layouts depending on participant count

### M-009: Empty Declaration Validation
- **Location**: `internal/infrastructure/parser/parser.go`
- **Description**: Empty structures accepted without warning:
  - `enum Empty {}`
  - `states Empty {}`
  - `provides EmptyAPI {}`
- **Impact**: Useless declarations pollute diagrams

### M-010: formatExpr Silent Fallback
- **Location**: All transformers
- **Description**: Unknown expression types return "..." without error/warning
- **Impact**: Silent data loss in complex expressions

## Low Priority

### L-001: No Pagination Support
- **Location**: `internal/infrastructure/renderer/svg/renderer.go`
- **Description**: Very large diagrams (100+ nodes) create single huge SVG
- **Impact**: Performance and usability issues for large projects

### L-002: No Caching
- **Location**: Multiple locations
- **Description**: Repeated calculations for same diagram not cached
- **Impact**: Slower repeated renders

### L-003: No Type Aliases
- **Location**: AST/Parser
- **Description**: Cannot define `type UserId = string`
- **Impact**: Limited type expressiveness

### L-004: No Generic Types
- **Location**: AST/Parser
- **Description**: Cannot express `List<T>` semantically
- **Impact**: Limited type system

### L-005: Limited Annotations
- **Location**: AST/Parser
- **Description**: Missing common annotations:
  - `@deprecated`
  - `@override`
- **Impact**: Limited metadata expressiveness

### L-006: Generic Error Messages
- **Location**: `internal/infrastructure/parser/parser.go`
- **Description**: "unexpected token" doesn't suggest what's expected
- **Impact**: Poor developer experience

### L-007: No Warning System
- **Location**: Multiple locations
- **Description**: No warnings for:
  - Unused imports
  - Unused types
  - Style violations
- **Impact**: No early feedback for code quality

### L-008: O(n^2) Complexity in Layout
- **Location**: `internal/infrastructure/renderer/svg/renderer.go`
- **Description**: Barycenter optimization and waypoint calculation have quadratic complexity
- **Impact**: Slow rendering for large diagrams

### L-009: Negative Number Handling
- **Location**: `internal/infrastructure/parser/lexer.go`
- **Description**: Negative numbers parsed as unary expression, not single token
- **Impact**: Inconsistent AST for negative literals

### L-010: Reserved Keyword Validation
- **Location**: `internal/infrastructure/parser/parser.go`
- **Description**: Keywords allowed as identifiers in some contexts
- **Impact**: Potential parsing ambiguity

## Enhancement Requests

### E-001: Multi-line Parameter Formatting
- **Description**: Support for multi-line parameter lists in method signatures
- **Benefit**: Better readability for complex methods

### E-002: Diagram Themes
- **Description**: Support for custom color themes
- **Benefit**: Better visual customization

### E-003: Export Formats
- **Description**: Support PNG, PDF export in addition to SVG
- **Benefit**: More output options

### E-004: Live Preview
- **Description**: Watch mode with live diagram updates
- **Benefit**: Better development workflow

### E-005: Language Server Protocol
- **Description**: LSP support for IDE integration
- **Benefit**: Better editing experience

---

## Issue Count Summary

| Priority | Total | Fixed | Remaining |
|----------|-------|-------|-----------|
| Critical | 3 | 3 | 0 |
| High | 5 | 3 | 2 |
| Medium | 10 | 5 | 5 |
| Low | 10 | 0 | 10 |
| Enhancement | 5 | 0 | 5 |
| **Total** | **33** | **11** | **22** |

---

## New Components Added

### Validator Package
- **Location**: `internal/application/validator/validator.go`
- **Features**:
  - `Validate()` - Duplicate validation
  - `ValidateReferences()` - Undefined reference validation
  - `ValidateExpressionDepth()` - Nesting depth validation
  - `ValidateDeadCode()` - Dead code detection
  - `ValidateAll()` - Run all validations

### Error Types Added
- **Location**: `internal/domain/errors/errors.go`
- **New Types**:
  - `MultiError` - Collects multiple errors
  - `ValidationError` - For validation errors (duplicate, undefined, invalid)

---

Last updated: 2026-02-02
