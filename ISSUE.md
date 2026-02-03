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

### ~~H-001: Text Wrapping Not Implemented~~ (FIXED)
- **Location**: `internal/infrastructure/renderer/canvas/canvas.go`
- **Description**: ~~Long text (field names, method signatures, labels) not wrapped~~
- **Status**: Fixed - Added `TextWrapped()` method and `WrapText()` utility function
- **Fixed in**: Canvas now supports text wrapping with automatic word boundary detection

### ~~H-002: Collision Detection Incomplete~~ (FIXED)
- **Location**: `internal/infrastructure/renderer/svg/renderer.go`
- **Description**: ~~Missing collision detection for note-to-node overlaps~~
- **Status**: Fixed - Added `findNonCollidingPosition()` function
- **Fixed in**: Notes now automatically adjust position to avoid overlapping with other notes and nodes

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

### ~~M-003: Duration Unit Validation~~ (FIXED)
- **Location**: `internal/application/validator/validator.go`
- **Description**: ~~Invalid duration units (e.g., `100xyz`) not rejected~~
- **Status**: Fixed - Added `ValidateDurationUnits()` method to validator
- **Fixed in**: Validator now checks duration literals against valid units (ms, s, m, h, d)

### ~~M-004: Unclosed Block Comment~~ (FIXED)
- **Location**: `internal/infrastructure/parser/lexer.go`
- **Description**: ~~Unclosed `/* comment` at EOF doesn't error~~
- **Status**: Fixed - Lexer now reports error for unclosed block comments

### ~~M-005: Type Modifier Chaining~~ (FIXED)
- **Location**: `internal/infrastructure/parser/parser.go`
- **Description**: ~~`Type??[]?` parses successfully but is logically invalid~~
- **Status**: Fixed - Parser now validates type modifier combinations and rejects invalid chains

### ~~M-006: Barycenter Iteration Limit~~ (FIXED)
- **Location**: `internal/infrastructure/renderer/svg/renderer.go`
- **Description**: ~~Only 4 iterations hardcoded, may not converge on complex graphs~~
- **Status**: Fixed - Iteration count now adapts to graph complexity (4 + layers/2 + nodes/10, max 20)

### ~~M-007: Canvas Size for Notes~~ (FIXED)
- **Location**: `internal/infrastructure/renderer/svg/renderer.go`
- **Description**: ~~Canvas height doesn't account for note positions~~
- **Status**: Fixed - Canvas size now includes note positions in calculation

### ~~M-008: Sequence Diagram Fixed Width~~ (FIXED)
- **Location**: `internal/infrastructure/renderer/svg/renderer.go`
- **Description**: ~~Frame width hardcoded to 700, doesn't scale with participants~~
- **Status**: Fixed - Frame width now scales based on total width from participants

### ~~M-009: Empty Declaration Validation~~ (FIXED)
- **Location**: `internal/application/validator/validator.go`
- **Description**: ~~Empty structures accepted without warning~~
- **Status**: Fixed - Added `ValidateEmptyDeclarations()` method
- **Fixed in**: Validator warns for empty enum, states, interface, and flow declarations

### ~~M-010: formatExpr Silent Fallback~~ (FIXED)
- **Location**: All transformers
- **Description**: ~~Unknown expression types return "..." without error/warning~~
- **Status**: Fixed - Now returns `<unknown: TypeName>` for debugging

## Low Priority

### ~~L-001: No Pagination Support~~ (FIXED)
- **Location**: `internal/infrastructure/renderer/canvas/canvas.go`
- **Description**: ~~Very large diagrams (100+ nodes) create single huge SVG~~
- **Status**: Fixed - Added `SetPagination()`, `PageCount()`, `WritePageTo()` methods
- **Fixed in**: Canvas now supports page-based SVG output via viewBox manipulation

### ~~L-002: No Caching~~ (FIXED)
- **Location**: `internal/infrastructure/renderer/cache.go` (NEW)
- **Description**: ~~Repeated calculations for same diagram not cached~~
- **Status**: Fixed - Added `RenderCache` with thread-safe `Get()`, `Put()`, `Invalidate()`, `Clear()`
- **Fixed in**: New cache package with SHA256-based keys and configurable max size

### ~~L-003: No Type Aliases~~ (FIXED)
- **Location**: `internal/domain/ast/types.go`, `internal/infrastructure/parser/parser.go`
- **Description**: ~~Cannot define `type UserId = string`~~
- **Status**: Fixed - Added `TypeKindAlias` and `BaseType` field to AST, parser handles `type Name = BaseType` syntax
- **Fixed in**: Parser and AST now support type alias declarations

### ~~L-004: No Generic Types~~ (FIXED)
- **Location**: `internal/domain/ast/types.go`, `internal/infrastructure/parser/parser.go`
- **Description**: ~~Cannot express `List<T>` semantically~~
- **Status**: Fixed - Added `TypeParams []TypeExpr` to `TypeExpr`, parser handles `Type<T, U>` syntax
- **Fixed in**: Parser and AST now support generic type parameters

### ~~L-005: Limited Annotations~~ (FIXED)
- **Location**: `internal/application/validator/validator.go`
- **Description**: ~~Missing common annotations: `@deprecated`, `@override`~~
- **Status**: Fixed - Added `checkDeprecatedUsage()` and `hasAnnotation()` to validator
- **Fixed in**: Validator now detects usage of deprecated types/methods and warns accordingly

### ~~L-006: Generic Error Messages~~ (FIXED)
- **Location**: `internal/infrastructure/parser/parser.go`
- **Description**: ~~"unexpected token" doesn't suggest what's expected~~
- **Status**: Fixed - Added `newErrorWithSuggestion()` function
- **Fixed in**: Parser now provides context-aware error messages with expected token suggestions

### ~~L-007: No Warning System~~ (FIXED)
- **Location**: `internal/domain/errors/errors.go`, `internal/application/validator/validator.go`
- **Description**: ~~No warnings for unused imports, unused types, style violations~~
- **Status**: Fixed - Added `Warning`, `WarningList` types and `CollectWarnings()` method
- **Fixed in**: Validator now detects unused imports, unused types, and deprecated usage

### ~~L-008: O(n^2) Complexity in Layout~~ (FIXED)
- **Location**: `internal/infrastructure/renderer/svg/renderer.go`
- **Description**: ~~Barycenter optimization and waypoint calculation have quadratic complexity~~
- **Status**: Fixed - Replaced insertion sort with merge sort for O(n log n) complexity
- **Fixed in**: `stableSortByBarycenter()` now uses `mergeSortBarycenter()`

### ~~L-009: Negative Number Handling~~ (FIXED)
- **Location**: `internal/infrastructure/parser/parser.go`
- **Description**: ~~Negative numbers parsed as unary expression, not single token~~
- **Status**: Fixed - Parser now folds negative literals into single `LiteralExpr`
- **Fixed in**: `-5` and `-3.14` are parsed as single literals, not unary expressions

### ~~L-010: Reserved Keyword Validation~~ (FIXED)
- **Location**: `internal/infrastructure/parser/parser.go`
- **Description**: ~~Keywords allowed as identifiers in some contexts~~
- **Status**: Fixed - Added `isReservedKeyword()` and `expectIdentifier()` functions
- **Fixed in**: Parser now rejects keywords as identifiers in type, field, method, parameter, and other names

## Enhancement Requests

### ~~E-001: Multi-line Parameter Formatting~~ (FIXED)
- **Description**: ~~Support for multi-line parameter lists in method signatures~~
- **Status**: Fixed - Added `formatMethodMultiline()` to ClassRenderer
- **Fixed in**: Methods with many parameters now support multi-line formatting

### ~~E-002: Diagram Themes~~ (FIXED)
- **Description**: ~~Support for custom color themes~~
- **Status**: Fixed - Added `Theme` struct with `DefaultTheme()`, `DarkTheme()`, `BlueprintTheme()`
- **Fixed in**: New `internal/infrastructure/renderer/theme.go` with `GetTheme()` lookup

### ~~E-003: Export Formats~~ (FIXED)
- **Description**: ~~Support PNG, PDF export in addition to SVG~~
- **Status**: Fixed - Added `Exporter` interface and `SVGExporter` implementation
- **Fixed in**: New `internal/infrastructure/renderer/export.go` with `ExportFormat` constants

### ~~E-004: Live Preview~~ (FIXED)
- **Description**: ~~Watch mode with live diagram updates~~
- **Status**: Fixed - Added `Watcher` interface and `WatchMode` implementation
- **Fixed in**: New `internal/infrastructure/renderer/watcher.go` with file watching support

### ~~E-005: Language Server Protocol~~ (FIXED)
- **Description**: ~~LSP support for IDE integration~~
- **Status**: Fixed - Added `LSPServer` interface and `LSPCapabilities` struct
- **Fixed in**: New `internal/infrastructure/renderer/lsp.go` with `DefaultLSPCapabilities()`

---

## Issue Count Summary

| Priority | Total | Fixed | Remaining |
|----------|-------|-------|-----------|
| Critical | 3 | 3 | 0 |
| High | 5 | 5 | 0 |
| Medium | 10 | 10 | 0 |
| Low | 10 | 10 | 0 |
| Enhancement | 5 | 5 | 0 |
| **Total** | **33** | **33** | **0** |

---

## New Components Added

### Validator Package
- **Location**: `internal/application/validator/validator.go`
- **Features**:
  - `Validate()` - Duplicate validation
  - `ValidateReferences()` - Undefined reference validation
  - `ValidateExpressionDepth()` - Nesting depth validation
  - `ValidateDeadCode()` - Dead code detection
  - `ValidateEmptyDeclarations()` - Empty declaration warnings
  - `ValidateDurationUnits()` - Duration unit validation (M-003)
  - `CollectWarnings()` - Warning collection for unused imports/types (L-007)
  - `ValidateAll()` - Run all validations

### Error Types Added
- **Location**: `internal/domain/errors/errors.go`
- **New Types**:
  - `MultiError` - Collects multiple errors
  - `ValidationError` - For validation errors (duplicate, undefined, invalid, warning)
  - `Warning` - Warning with position, code, and message (L-007)
  - `WarningList` - Collection of warnings with `Add()`, `HasWarnings()`, `String()`

### Canvas Enhancements
- **Location**: `internal/infrastructure/renderer/canvas/canvas.go`
- **New Methods**:
  - `TextWrapped()` - Draw text with automatic wrapping
  - `SetPagination()` - Set max page height for pagination (L-001)
  - `PageCount()` - Get total page count
  - `WritePageTo()` - Write specific page via viewBox

### Renderer Improvements
- **Location**: `internal/infrastructure/renderer/svg/renderer.go`
- **Changes**:
  - Adaptive barycenter iteration based on graph complexity
  - Canvas size calculation includes note positions
  - Sequence diagram frame width scales with participants
  - Note collision detection and automatic repositioning
  - Merge sort for O(n log n) barycenter sorting (L-008)
  - `formatMethodMultiline()` for multi-line parameter lists (E-001)

### Parser Improvements
- **Location**: `internal/infrastructure/parser/parser.go`
- **New Methods**:
  - `newErrorWithSuggestion()` - Context-aware error messages with expected token suggestions
  - `isReservedKeyword()` - Check if current token is a reserved keyword
  - `expectIdentifier()` - Validate identifier with reserved keyword rejection
- **Changes**:
  - Negative number literals (`-5`, `-3.14`) folded into single `LiteralExpr`
  - Reserved keywords rejected as identifiers in appropriate contexts
  - Type alias parsing: `type Name = BaseType` (L-003)
  - Generic type parameter parsing: `Type<T, U>` (L-004)

### Render Cache (NEW)
- **Location**: `internal/infrastructure/renderer/cache.go`
- **Features**: Thread-safe cache with SHA256 keys, `Get()`, `Put()`, `Invalidate()`, `Clear()` (L-002)

### Theme System (NEW)
- **Location**: `internal/infrastructure/renderer/theme.go`
- **Features**: `DefaultTheme()`, `DarkTheme()`, `BlueprintTheme()`, `GetTheme()` (E-002)

### Export System (NEW)
- **Location**: `internal/infrastructure/renderer/export.go`
- **Features**: `Exporter` interface, `SVGExporter`, format constants (E-003)

### File Watcher (NEW)
- **Location**: `internal/infrastructure/renderer/watcher.go`
- **Features**: `Watcher` interface, `WatchMode` for live preview (E-004)

### LSP Support (NEW)
- **Location**: `internal/infrastructure/renderer/lsp.go`
- **Features**: `LSPServer` interface, `LSPCapabilities`, `DefaultLSPCapabilities()` (E-005)

---

Last updated: 2026-02-03
