# Pact DSL Changelog

This document tracks all issues found and fixes made to the Pact DSL project.

---

## Latest Changes (2026-02-02)

### Issue #14: Silent Failures in Parser
- **Problem**: Parser ignored errors from `strconv.ParseInt`, `strconv.ParseFloat`, `strconv.Atoi`
- **Impact**: Invalid numeric literals silently converted to zero, masking overflow errors
- **Fix**: Added proper error handling with descriptive error messages
- **Commit**: `3b804a2`
- **Files Changed**: `internal/infrastructure/parser/parser.go`

### Issue #13: Incomplete Swimlane Assignment
- **Problem**: Start/End/Decision nodes in flow diagrams had no swimlane
- **Impact**: Nodes not associated with any participant in swimlane view
- **Fix**: Assign component name as default swimlane to all unassigned nodes
- **Commit**: `834cb03`
- **Files Changed**: `internal/application/transformer/flow.go`, `flow_test.go`

### Issue #12: Type Modifiers Not Processed
- **Problem**: `Nullable` and `Array` type modifiers ignored in class transformer
- **Impact**: Types displayed as `Type` instead of `Type?` or `Type[]`
- **Fix**: Added `formatTypeExpr` helper to format types with modifiers
- **Commit**: `a68a260`
- **Files Changed**: `internal/application/transformer/class.go`

### Issue #11: Participants Not in Dependencies
- **Problem**: Call targets not in dependencies weren't added as participants
- **Impact**: Messages to undeclared targets failed silently
- **Fix**: Auto-discover participants from call targets in all step types
- **Commit**: `bf866bf`
- **Files Changed**: `internal/application/transformer/sequence.go`

### Issue #10: AST Validation Missing
- **Problem**: No validation for required Initial state in state diagrams
- **Impact**: Invalid state machines accepted without error
- **Fix**: Added validation + implicit state creation for Initial state
- **Commit**: `7bbe454`
- **Files Changed**: `internal/application/transformer/state.go`

### Issue #9: State Annotations Ignored
- **Problem**: Annotations on state declarations not transformed
- **Impact**: Metadata lost in state diagrams
- **Fix**: Added `transformAnnotation` to StateTransformer
- **Commit**: `3efa7e8`
- **Files Changed**: `internal/application/transformer/state.go`

### Issue #8: Guard Expressions Limited
- **Problem**: Guard conditions only supported simple variable names
- **Impact**: Complex guards like `a && b` displayed as "..."
- **Fix**: Added `formatExpr` to StateTransformer for complex expressions
- **Commit**: `e6bf9ba`
- **Files Changed**: `internal/application/transformer/state.go`

### Issue #7: RelationDecl.Alias Unused
- **Problem**: Relationship aliases not displayed on edges
- **Impact**: `depends on X as alias` didn't show "alias" label
- **Fix**: Set `edge.Label` from `rel.Alias`
- **Commit**: `fd4a88e`
- **Files Changed**: `internal/application/transformer/class.go`

### Issue #6: Method.Throws Lost
- **Problem**: Method throw declarations not included in diagram
- **Impact**: Exception information missing from class diagrams
- **Fix**: Added `Throws []string` to class.Method, updated renderer
- **Commit**: `9f08172`
- **Files Changed**: `internal/domain/diagram/class/model.go`, `internal/application/transformer/class.go`, `internal/infrastructure/renderer/svg/renderer.go`

### Issue #5: Notes Unused
- **Problem**: `@note` and `@description` annotations not converted to Notes
- **Impact**: Documentation metadata not visible in diagrams
- **Fix**: Added `extractNote` method to ClassTransformer
- **Commit**: `63d4ac8`
- **Files Changed**: `internal/application/transformer/class.go`

### Issue #4: TernaryExpr/NullishExpr Not Formatted
- **Problem**: Complex expressions returned "..." in `formatExpr`
- **Impact**: Ternary and nullish expressions displayed incorrectly
- **Fix**: Added TernaryExpr and NullishExpr cases to formatExpr
- **Commit**: `23cff4e`
- **Files Changed**: `internal/application/transformer/flow.go`, `sequence.go`

### Issue #3: ReturnStep/ThrowStep Not in Sequence
- **Problem**: Return and throw steps not converted to sequence events
- **Impact**: Missing return/throw visualization in sequence diagrams
- **Fix**: Added NoteEvent type with color-coded boxes
- **Commit**: `bad1ccb`
- **Files Changed**: `internal/domain/diagram/sequence/model.go`, `internal/application/transformer/sequence.go`, `internal/infrastructure/renderer/svg/renderer.go`

### Issue #2: Visitor Not Walking Components[]
- **Problem**: AST Visitor only walked `Component`, not `Components[]`
- **Impact**: Multi-component files partially processed
- **Fix**: Updated `ast.Walk` to iterate over `Components[]` array
- **Commit**: `07e802a`
- **Files Changed**: `internal/domain/ast/visitor.go`, all transformers

### Issue #1: FragmentEvent.Events Empty
- **Problem**: Alt/loop/opt fragment events not populated with content
- **Impact**: Empty fragment boxes in sequence diagrams
- **Fix**: Changed `transformSteps` to return events, added `AltEvents` field
- **Commit**: `feecd04`
- **Files Changed**: `internal/domain/diagram/sequence/model.go`, `internal/application/transformer/sequence.go`

---

## Rendering Improvements

### Diagram Rendering Issues Fix
- **Problem**: Multiple visual issues in rendered diagrams
- **Commit**: `505b61d`
- **Changes**:
  - Fixed node/edge overlaps
  - Fixed label collisions
  - Improved text formatting

### Orthogonal Routing
- **Problem**: Edges crossed nodes and overlapped
- **Commit**: `25470f5`
- **Fix**: Implemented orthogonal routing for all diagram types

### Auto-sizing
- **Problem**: Node sizes didn't fit content
- **Commit**: `8100365`
- **Fix**: Implemented auto-sizing based on text content

### Barycenter Optimization
- **Problem**: Too many edge crossings
- **Commit**: `719e3b9`
- **Fix**: Implemented barycenter method for edge crossing minimization

### Advanced Orthogonal Routing
- **Problem**: Edges overlapped obstacles
- **Commit**: `4ccb11f`
- **Fix**: Added obstacle avoidance to edge routing

### Distributed Edge Endpoints
- **Problem**: Multiple edges shared same endpoint
- **Commit**: `7efcc2f`
- **Fix**: Distributed edge endpoints along node boundaries

### Layer-based Layout
- **Problem**: Nodes positioned randomly
- **Commit**: `e229d23`
- **Fix**: Implemented Sugiyama-style layer-based layout

### Text-based Auto-sizing
- **Problem**: Node sizes didn't account for text length
- **Commit**: `caa73c8`
- **Fix**: Calculate node size from actual text measurements

### Phase 1 & 2 Enhancements
- **Commit**: `82e849b`
- **Changes**:
  - Phase 1: Basic rendering improvements
  - Phase 2: Layout algorithm enhancements

---

## Core Fixes

### State and Flow Transformers
- **Problem**: Transformers not generating proper diagrams
- **Commit**: `8282a8b`
- **Fix**: Fixed transformation logic for state and flow diagrams

### SVG Renderers
- **Problem**: Renderers producing invalid/empty SVG
- **Commit**: `9c52014`
- **Fix**: Fixed all SVG renderers for proper output

### Transitions, Messages, and Edges
- **Problem**: These elements not drawn
- **Commit**: `f05be24`
- **Fix**: Added rendering logic for transitions, messages, edges

### SVG XML Escape
- **Problem**: Special characters not escaped in SVG
- **Commit**: `76b8f29`
- **Fix**: Proper XML escaping for text content

### Transformer Syntax
- **Problem**: Transformers didn't match parser output
- **Commit**: `b734ef6`
- **Fix**: Updated transformers to match parser AST structure

---

## Initial Implementation

### CLI Implementation
- **Commit**: `69c518c`
- Implemented command-line interface for diagram generation

### Import Resolver
- **Commit**: `ba48fc8`
- Implemented dependency resolution for imports

### Parser
- **Commit**: `d7d8d12`
- Implemented full parser for Pact DSL

### Lexer
- **Commit**: `e820d20`
- Implemented complete lexer with all token types

### TDD Test Suite
- **Commit**: `14cf283`
- Added comprehensive test suite based on TEST_SPEC.md

---

## Summary Statistics

| Category | Count |
|----------|-------|
| Critical Issues Fixed | 3 |
| High Priority Issues Fixed | 6 |
| Medium Priority Issues Fixed | 5 |
| Rendering Improvements | 10 |
| Core Fixes | 5 |
| Initial Implementation | 5 |
| **Total Commits** | **50+** |

---

## Known Remaining Issues

See [ISSUE.md](/ISSUE.md) for a comprehensive list of known issues and enhancement requests.

---

Last updated: 2026-02-02
