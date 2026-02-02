# Pact Architecture

## Overview

Pact is a DSL-to-diagram generator with a clean layered architecture following Domain-Driven Design principles.

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI (cmd/pact)                       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Public API (pkg/pact)                   │
└─────────────────────────────────────────────────────────────┘
                              │
         ┌────────────────────┼────────────────────┐
         ▼                    ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   Application   │  │     Domain      │  │ Infrastructure  │
│  (transformer)  │  │   (ast, model)  │  │ (parser, render)│
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

## Directory Structure

```
pact/
├── cmd/pact/           # CLI entry point
├── pkg/pact/           # Public API
├── internal/
│   ├── application/
│   │   └── transformer/  # AST to diagram transformation
│   ├── domain/
│   │   ├── ast/          # Abstract Syntax Tree types
│   │   ├── config/       # Configuration types
│   │   ├── diagram/      # Diagram model types
│   │   │   ├── class/
│   │   │   ├── sequence/
│   │   │   ├── flow/
│   │   │   ├── state/
│   │   │   └── common/
│   │   └── errors/       # Error types
│   └── infrastructure/
│       ├── config/       # Config file loading
│       ├── parser/       # Lexer and Parser
│       ├── renderer/     # SVG rendering
│       │   ├── svg/
│       │   └── canvas/
│       └── resolver/     # Import resolution
├── test/
│   ├── e2e/              # End-to-end tests
│   └── integration/      # Integration tests
└── docs/                 # Documentation
```

## Layer Responsibilities

### 1. Domain Layer (`internal/domain/`)

Pure domain models with no external dependencies.

#### AST (`domain/ast/`)

```go
// Core AST types
type SpecFile struct { ... }
type ComponentDecl struct { ... }
type TypeDecl struct { ... }
type FlowDecl struct { ... }
type StatesDecl struct { ... }

// Expressions
type Expr interface { ... }
type LiteralExpr struct { ... }
type CallExpr struct { ... }
type BinaryExpr struct { ... }

// Steps
type Step interface { ... }
type IfStep struct { ... }
type ForStep struct { ... }
```

#### Diagram Models (`domain/diagram/`)

Each diagram type has its own model:

**Class Diagram**
```go
type Diagram struct {
    Nodes []Node
    Edges []Edge
    Notes []common.Note
}

type Node struct {
    ID, Name     string
    Stereotype   string
    Attributes   []Attribute
    Methods      []Method
    Annotations  []common.Annotation
}
```

**Sequence Diagram**
```go
type Diagram struct {
    Participants []Participant
    Events       []Event
}

type Event interface { ... }
type MessageEvent struct { ... }
type FragmentEvent struct { ... }  // alt, loop, opt
type NoteEvent struct { ... }      // note, return, throw
```

**Flow Diagram**
```go
type Diagram struct {
    Nodes     []Node
    Edges     []Edge
    Swimlanes []Swimlane
}

type Node struct {
    ID, Label string
    Shape     NodeShape  // terminal, process, decision, etc.
    Swimlane  string
}
```

**State Diagram**
```go
type Diagram struct {
    States      []State
    Transitions []Transition
}

type State struct {
    ID, Name string
    Type     StateType  // initial, final, atomic, compound, parallel
    Entry    []string
    Exit     []string
    Children []State
    Regions  []Region
}
```

### 2. Application Layer (`internal/application/`)

Transforms AST to diagram models.

#### Transformers

```go
// Class transformer
type ClassTransformer struct{}
func (t *ClassTransformer) Transform(files []*ast.SpecFile, opts *TransformOptions) (*class.Diagram, error)

// Sequence transformer
type SequenceTransformer struct{}
func (t *SequenceTransformer) Transform(files []*ast.SpecFile, opts *SequenceOptions) (*sequence.Diagram, error)

// Flow transformer
type FlowTransformer struct{}
func (t *FlowTransformer) Transform(files []*ast.SpecFile, opts *FlowOptions) (*flow.Diagram, error)

// State transformer
type StateTransformer struct{}
func (t *StateTransformer) Transform(files []*ast.SpecFile, opts *StateOptions) (*state.Diagram, error)
```

### 3. Infrastructure Layer (`internal/infrastructure/`)

External concerns: parsing, rendering, file I/O.

#### Parser (`infrastructure/parser/`)

**Lexer**
```go
type Lexer struct { ... }
func NewLexer(input string) *Lexer
func (l *Lexer) NextToken() Token

type Token struct {
    Type    TokenType
    Literal string
    Pos     Position
}
```

**Parser**
```go
type Parser struct { ... }
func NewParser(input string) *Parser
func (p *Parser) Parse() (*ast.SpecFile, error)
```

#### Renderer (`infrastructure/renderer/`)

**SVG Renderer**
```go
type Renderer struct{}
func (r *Renderer) RenderClass(d *class.Diagram) (string, error)
func (r *Renderer) RenderSequence(d *sequence.Diagram) (string, error)
func (r *Renderer) RenderFlow(d *flow.Diagram) (string, error)
func (r *Renderer) RenderState(d *state.Diagram) (string, error)
```

**Layout Algorithms**
- Layer-based layout (Sugiyama-style)
- Barycenter method for edge crossing reduction
- Orthogonal edge routing with obstacle avoidance

#### Resolver (`infrastructure/resolver/`)

```go
type Resolver struct { ... }
func (r *Resolver) ResolveImports(file *ast.SpecFile) ([]*ast.SpecFile, error)
```

### 4. Public API (`pkg/pact/`)

Clean public interface for library users.

```go
type Pact struct { ... }

func New() *Pact

func (p *Pact) ParseFile(path string) (*ast.SpecFile, error)
func (p *Pact) ParseString(input string) (*ast.SpecFile, error)

func (p *Pact) GenerateClassDiagram(files []*ast.SpecFile, opts ClassOptions) (string, error)
func (p *Pact) GenerateSequenceDiagram(files []*ast.SpecFile, opts SequenceOptions) (string, error)
func (p *Pact) GenerateFlowDiagram(files []*ast.SpecFile, opts FlowOptions) (string, error)
func (p *Pact) GenerateStateDiagram(files []*ast.SpecFile, opts StateOptions) (string, error)
```

### 5. CLI (`cmd/pact/`)

Command-line interface.

```bash
# Generate all diagrams
pact generate input.pact -o output/

# Generate specific diagram type
pact generate input.pact --class -o class.svg
pact generate input.pact --sequence --flow=ProcessOrder -o sequence.svg
pact generate input.pact --state --states=OrderState -o state.svg
pact generate input.pact --flow --flow=ProcessOrder -o flow.svg
```

## Data Flow

```
.pact file
    │
    ▼
┌─────────┐
│  Lexer  │  → Token stream
└─────────┘
    │
    ▼
┌─────────┐
│ Parser  │  → AST (SpecFile)
└─────────┘
    │
    ▼
┌──────────┐
│ Resolver │  → Resolved AST (with imports)
└──────────┘
    │
    ├─────────────────┬─────────────────┬─────────────────┐
    ▼                 ▼                 ▼                 ▼
┌───────────┐   ┌───────────┐   ┌───────────┐   ┌───────────┐
│  Class    │   │ Sequence  │   │   Flow    │   │   State   │
│Transformer│   │Transformer│   │Transformer│   │Transformer│
└───────────┘   └───────────┘   └───────────┘   └───────────┘
    │                 │                 │                 │
    ▼                 ▼                 ▼                 ▼
┌───────────────────────────────────────────────────────────┐
│                      SVG Renderer                          │
└───────────────────────────────────────────────────────────┘
    │
    ▼
  .svg file
```

## Key Design Decisions

### 1. Separation of AST and Diagram Models

The AST represents the source code structure, while diagram models represent visual elements. This separation allows:
- Different diagram types from the same AST
- Easier testing of each layer
- Clean transformation logic

### 2. Interface-based Step and Expression Types

Using interfaces for `Step` and `Expr` allows:
- Type-safe pattern matching
- Easy addition of new step/expression types
- Clean visitor pattern support

### 3. Options Pattern for Transformers

Each transformer accepts an options struct:
```go
type SequenceOptions struct {
    FlowName      string
    IncludeReturn bool
}
```

This allows:
- Selective diagram generation
- Customization without API changes
- Clear parameter semantics

### 4. Layer-based Layout Algorithm

The renderer uses a Sugiyama-style layout:
1. Assign nodes to layers
2. Order nodes within layers (barycenter)
3. Calculate positions
4. Route edges with waypoints

### 5. Visitor Pattern for AST Traversal

```go
type Visitor interface {
    VisitComponent(comp *ComponentDecl) error
    VisitType(typ *TypeDecl) error
    VisitFlow(flow *FlowDecl) error
    // ...
}

func Walk(v Visitor, node interface{}) error
```

## Error Handling

Errors are defined in `domain/errors/`:

```go
type ParseError struct {
    Pos     Position
    Message string
}

type TransformError struct {
    Source  string
    Target  string
    Message string
}

type RenderError struct {
    DiagramType string
    Message     string
}
```

## Configuration

Configuration in `pact.yaml`:

```yaml
output:
  dir: "./diagrams"
  format: "svg"

diagrams:
  class:
    filterComponents: ["OrderService", "UserService"]
  sequence:
    includeReturn: true
  flow:
    includeSwimlanes: true
```

---

Last updated: 2026-02-02
