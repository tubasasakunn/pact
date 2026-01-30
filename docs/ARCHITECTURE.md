# Pact アーキテクチャ

## 概要

Pact は DDD（ドメイン駆動設計）とクリーンアーキテクチャに基づいて設計されています。

```
interfaces (CLI)
    │
    ▼
application (ユースケース)
    │
    ├──────────────────┐
    ▼                  ▼
domain (モデル)    infrastructure (実装)
    ▲                  │
    └──────────────────┘
```

**依存ルール**:
- `domain` は他に依存しない（純粋なビジネスロジック）
- `application` は `domain` にのみ依存
- `infrastructure` は `domain` に依存（インターフェースを実装）
- `interfaces` は `application` に依存
- `pkg` は外部向けのファサード

---

## ディレクトリ構成

```
pact/
├── cmd/
│   └── pact/
│       └── main.go
│
├── internal/
│   ├── domain/
│   │   ├── ast/
│   │   │   ├── nodes.go
│   │   │   ├── expressions.go
│   │   │   ├── statements.go
│   │   │   ├── types.go
│   │   │   ├── position.go
│   │   │   └── visitor.go
│   │   │
│   │   ├── config/
│   │   │   └── config.go
│   │   │
│   │   ├── diagram/
│   │   │   ├── common/
│   │   │   │   └── types.go
│   │   │   ├── class/
│   │   │   │   └── model.go
│   │   │   ├── sequence/
│   │   │   │   └── model.go
│   │   │   ├── state/
│   │   │   │   └── model.go
│   │   │   └── flow/
│   │   │       └── model.go
│   │   │
│   │   └── errors/
│   │       └── errors.go
│   │
│   ├── application/
│   │   ├── parser/
│   │   │   └── service.go
│   │   ├── project/
│   │   │   └── service.go
│   │   ├── transformer/
│   │   │   ├── service.go
│   │   │   ├── class.go
│   │   │   ├── sequence.go
│   │   │   ├── state.go
│   │   │   └── flow.go
│   │   └── renderer/
│   │       └── service.go
│   │
│   ├── infrastructure/
│   │   ├── config/
│   │   │   └── loader.go
│   │   ├── parser/
│   │   │   ├── lexer.go
│   │   │   ├── token.go
│   │   │   ├── parser.go
│   │   │   └── adapter.go
│   │   ├── resolver/
│   │   │   └── import.go
│   │   ├── renderer/
│   │   │   ├── svg/
│   │   │   │   ├── renderer.go
│   │   │   │   ├── class.go
│   │   │   │   ├── sequence.go
│   │   │   │   ├── state.go
│   │   │   │   └── flow.go
│   │   │   └── canvas/
│   │   │       ├── canvas.go
│   │   │       ├── shapes.go
│   │   │       └── text.go
│   │   └── filesystem/
│   │       ├── reader.go
│   │       └── writer.go
│   │
│   └── interfaces/
│       └── cli/
│           ├── root.go
│           ├── init.go
│           ├── generate.go
│           ├── validate.go
│           ├── check.go
│           └── watch.go
│
├── pkg/
│   └── pact/
│       └── api.go
│
├── testdata/
│   ├── valid/
│   │   ├── auth_service.pact
│   │   └── user_repository.pact
│   └── invalid/
│       ├── syntax_error.pact
│       └── semantic_error.pact
│
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── ARCHITECTURE.md
```

---

## cmd/ - エントリポイント

### cmd/pact/main.go

```go
package main

import (
    "os"

    "pact/internal/interfaces/cli"
)

func main() {
    cmd := cli.NewRootCmd()
    if err := cmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

**責務**: プログラムの起動のみ。ビジネスロジックを含まない。

---

## internal/domain/ - ドメイン層

ビジネスロジックの中核。外部依存なし。

### internal/domain/ast/position.go

```go
package ast

import "fmt"

// Position はソースコード内の位置を表す
type Position struct {
    File   string
    Line   int
    Column int
    Offset int
}

func (p Position) String() string {
    return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Column)
}

// NoPos は位置情報がないことを表す
var NoPos = Position{}

func (p Position) IsValid() bool {
    return p.Line > 0
}
```

**責務**: ソースコード位置の表現

---

### internal/domain/ast/nodes.go

```go
package ast

// SpecFile は .pact ファイル全体を表す
type SpecFile struct {
    Path      string
    Imports   []ImportDecl
    Component ComponentDecl
}

// ImportDecl は import 文を表す
type ImportDecl struct {
    Pos   Position
    Path  string
    Alias *string
}

// ComponentDecl は component 宣言を表す
type ComponentDecl struct {
    Pos         Position
    Name        string
    Annotations []AnnotationDecl
    Body        ComponentBody
}

// ComponentBody は component の中身を表す
type ComponentBody struct {
    Types     []TypeDecl
    Relations []RelationDecl
    Provides  []InterfaceDecl
    Requires  []InterfaceDecl
    Flows     []FlowDecl
    States    []StatesDecl
}

// AnnotationDecl はアノテーションを表す
type AnnotationDecl struct {
    Pos  Position
    Name string
    Args []AnnotationArg
}

// AnnotationArg はアノテーションの引数を表す
type AnnotationArg struct {
    Key   *string // nil なら positional
    Value string
}
```

**責務**: AST のルートノードとコンポーネント構造の定義

---

### internal/domain/ast/types.go

```go
package ast

// TypeDecl は型定義を表す
type TypeDecl struct {
    Pos         Position
    Name        string
    Kind        TypeKind
    Annotations []AnnotationDecl
    Fields      []FieldDecl // struct の場合
    Values      []string    // enum の場合
}

type TypeKind string

const (
    TypeKindStruct TypeKind = "struct"
    TypeKindEnum   TypeKind = "enum"
)

// FieldDecl はフィールド定義を表す
type FieldDecl struct {
    Pos         Position
    Name        string
    Type        TypeExpr
    Visibility  Visibility
    Annotations []AnnotationDecl
}

type Visibility string

const (
    VisibilityPublic    Visibility = "public"
    VisibilityPrivate   Visibility = "private"
    VisibilityProtected Visibility = "protected"
    VisibilityPackage   Visibility = "package"
)

// TypeExpr は型の参照を表す
type TypeExpr struct {
    Pos      Position
    Name     string
    Nullable bool
    Array    bool
}

// RelationDecl は関係定義を表す
type RelationDecl struct {
    Pos         Position
    Target      string
    Kind        RelationKind
    TargetType  TargetType
    Alias       *string
    Annotations []AnnotationDecl
}

type RelationKind string

const (
    RelationDepends    RelationKind = "depends"
    RelationExtends    RelationKind = "extends"
    RelationImplements RelationKind = "implements"
    RelationContains   RelationKind = "contains"
    RelationAggregates RelationKind = "aggregates"
)

type TargetType string

const (
    TargetComponent TargetType = "component"
    TargetDatabase  TargetType = "database"
    TargetExternal  TargetType = "external"
    TargetQueue     TargetType = "queue"
    TargetActor     TargetType = "actor"
)

// InterfaceDecl はインターフェース定義を表す
type InterfaceDecl struct {
    Pos         Position
    Name        string
    Annotations []AnnotationDecl
    Methods     []MethodDecl
}

// MethodDecl はメソッド定義を表す
type MethodDecl struct {
    Pos         Position
    Name        string
    Params      []ParamDecl
    Returns     TypeExpr
    Throws      []string
    Async       bool
    Annotations []AnnotationDecl
}

// ParamDecl はパラメータ定義を表す
type ParamDecl struct {
    Pos  Position
    Name string
    Type TypeExpr
}
```

**責務**: 型、関係、インターフェースの AST 定義

---

### internal/domain/ast/statements.go

```go
package ast

// FlowDecl はフロー定義を表す
type FlowDecl struct {
    Pos         Position
    Name        string
    Annotations []AnnotationDecl
    Steps       []StepDecl
}

// StepDecl はフロー内のステップを表すインターフェース
type StepDecl interface {
    stepNode()
    GetPos() Position
    GetAnnotations() []AnnotationDecl
}

// baseStep は共通フィールドを持つ基底構造体
type baseStep struct {
    Pos         Position
    Annotations []AnnotationDecl
}

func (b baseStep) GetPos() Position              { return b.Pos }
func (b baseStep) GetAnnotations() []AnnotationDecl { return b.Annotations }

// AssignStep は代入文を表す
type AssignStep struct {
    baseStep
    Variable    string
    Value       ExprDecl
    ThrowOnNull *string // ?? throw Error
}

func (AssignStep) stepNode() {}

// CallStep は呼び出し文を表す
type CallStep struct {
    baseStep
    Target string
    Method string
    Args   []ExprDecl
    Async  bool
}

func (CallStep) stepNode() {}

// IfStep は条件分岐を表す
type IfStep struct {
    baseStep
    Condition ExprDecl
    Then      []StepDecl
    Else      []StepDecl
}

func (IfStep) stepNode() {}

// ForStep は for ループを表す
type ForStep struct {
    baseStep
    Variable string
    Iterable ExprDecl
    Body     []StepDecl
}

func (ForStep) stepNode() {}

// WhileStep は while ループを表す
type WhileStep struct {
    baseStep
    Condition ExprDecl
    Body      []StepDecl
}

func (WhileStep) stepNode() {}

// ReturnStep は return 文を表す
type ReturnStep struct {
    baseStep
    Value ExprDecl // nil 可
}

func (ReturnStep) stepNode() {}

// ThrowStep は throw 文を表す
type ThrowStep struct {
    baseStep
    Error string
}

func (ThrowStep) stepNode() {}

// StatesDecl はステートマシン定義を表す
type StatesDecl struct {
    Pos         Position
    Name        string
    Annotations []AnnotationDecl
    Initial     string
    Finals      []string
    States      []StateDecl
    Transitions []TransitionDecl
}

// StateDecl は状態定義を表す
type StateDecl struct {
    Pos         Position
    Name        string
    Kind        StateKind
    Annotations []AnnotationDecl
    Entry       []string
    Exit        []string
    Initial     *string // compound の場合
    States      []StateDecl
    Transitions []TransitionDecl
    Regions     []RegionDecl
}

type StateKind string

const (
    StateAtomic   StateKind = "atomic"
    StateCompound StateKind = "compound"
    StateParallel StateKind = "parallel"
)

// RegionDecl はリージョン定義を表す
type RegionDecl struct {
    Pos         Position
    Name        string
    Initial     string
    Finals      []string
    States      []StateDecl
    Transitions []TransitionDecl
}

// TransitionDecl は遷移定義を表す
type TransitionDecl struct {
    Pos         Position
    From        string
    To          string
    Trigger     TriggerDecl
    Guard       ExprDecl // nil 可
    Actions     []string
    Annotations []AnnotationDecl
}

// TriggerDecl はトリガーを表すインターフェース
type TriggerDecl interface {
    triggerNode()
}

// EventTrigger はイベントトリガーを表す
type EventTrigger struct {
    Pos  Position
    Name string
}

func (EventTrigger) triggerNode() {}

// AfterTrigger は時間トリガーを表す
type AfterTrigger struct {
    Pos   Position
    Value int
    Unit  string
}

func (AfterTrigger) triggerNode() {}

// WhenTrigger は条件トリガーを表す
type WhenTrigger struct {
    Pos       Position
    Condition ExprDecl
}

func (WhenTrigger) triggerNode() {}
```

**責務**: フロー、ステートマシンの AST 定義

---

### internal/domain/ast/expressions.go

```go
package ast

// ExprDecl は式を表すインターフェース
type ExprDecl interface {
    exprNode()
    GetPos() Position
}

// LiteralExpr はリテラルを表す
type LiteralExpr struct {
    Pos   Position
    Value interface{} // string | int | float64 | bool | nil
}

func (LiteralExpr) exprNode()           {}
func (e LiteralExpr) GetPos() Position { return e.Pos }

// VariableExpr は変数参照を表す
type VariableExpr struct {
    Pos  Position
    Name string
}

func (VariableExpr) exprNode()           {}
func (e VariableExpr) GetPos() Position { return e.Pos }

// CallExpr はメソッド呼び出しを表す
type CallExpr struct {
    Pos    Position
    Target string
    Method string
    Args   []ExprDecl
}

func (CallExpr) exprNode()           {}
func (e CallExpr) GetPos() Position { return e.Pos }

// FieldExpr はフィールドアクセスを表す
type FieldExpr struct {
    Pos    Position
    Object ExprDecl
    Field  string
}

func (FieldExpr) exprNode()           {}
func (e FieldExpr) GetPos() Position { return e.Pos }

// BinaryExpr は二項演算を表す
type BinaryExpr struct {
    Pos   Position
    Op    string
    Left  ExprDecl
    Right ExprDecl
}

func (BinaryExpr) exprNode()           {}
func (e BinaryExpr) GetPos() Position { return e.Pos }

// UnaryExpr は単項演算を表す
type UnaryExpr struct {
    Pos     Position
    Op      string
    Operand ExprDecl
}

func (UnaryExpr) exprNode()           {}
func (e UnaryExpr) GetPos() Position { return e.Pos }

// TernaryExpr は三項演算を表す
type TernaryExpr struct {
    Pos       Position
    Condition ExprDecl
    Then      ExprDecl
    Else      ExprDecl
}

func (TernaryExpr) exprNode()           {}
func (e TernaryExpr) GetPos() Position { return e.Pos }
```

**責務**: 式の AST 定義

---

### internal/domain/ast/visitor.go

```go
package ast

// Visitor は AST を走査するためのインターフェース
type Visitor interface {
    VisitSpecFile(node *SpecFile) error
    VisitImport(node *ImportDecl) error
    VisitComponent(node *ComponentDecl) error
    VisitType(node *TypeDecl) error
    VisitField(node *FieldDecl) error
    VisitRelation(node *RelationDecl) error
    VisitInterface(node *InterfaceDecl) error
    VisitMethod(node *MethodDecl) error
    VisitFlow(node *FlowDecl) error
    VisitStep(node StepDecl) error
    VisitStates(node *StatesDecl) error
    VisitState(node *StateDecl) error
    VisitTransition(node *TransitionDecl) error
    VisitExpr(node ExprDecl) error
    VisitAnnotation(node *AnnotationDecl) error
}

// Walk は AST を深さ優先で走査する
func Walk(v Visitor, node interface{}) error {
    switch n := node.(type) {
    case *SpecFile:
        if err := v.VisitSpecFile(n); err != nil {
            return err
        }
        for i := range n.Imports {
            if err := Walk(v, &n.Imports[i]); err != nil {
                return err
            }
        }
        return Walk(v, &n.Component)

    case *ImportDecl:
        return v.VisitImport(n)

    case *ComponentDecl:
        if err := v.VisitComponent(n); err != nil {
            return err
        }
        for i := range n.Annotations {
            if err := Walk(v, &n.Annotations[i]); err != nil {
                return err
            }
        }
        // Body の各要素を走査
        for i := range n.Body.Types {
            if err := Walk(v, &n.Body.Types[i]); err != nil {
                return err
            }
        }
        for i := range n.Body.Relations {
            if err := Walk(v, &n.Body.Relations[i]); err != nil {
                return err
            }
        }
        for i := range n.Body.Provides {
            if err := Walk(v, &n.Body.Provides[i]); err != nil {
                return err
            }
        }
        for i := range n.Body.Requires {
            if err := Walk(v, &n.Body.Requires[i]); err != nil {
                return err
            }
        }
        for i := range n.Body.Flows {
            if err := Walk(v, &n.Body.Flows[i]); err != nil {
                return err
            }
        }
        for i := range n.Body.States {
            if err := Walk(v, &n.Body.States[i]); err != nil {
                return err
            }
        }
        return nil

    // 他のノードタイプも同様に実装...

    default:
        return nil
    }
}

// BaseVisitor は Visitor のデフォルト実装を提供する
type BaseVisitor struct{}

func (BaseVisitor) VisitSpecFile(*SpecFile) error       { return nil }
func (BaseVisitor) VisitImport(*ImportDecl) error       { return nil }
func (BaseVisitor) VisitComponent(*ComponentDecl) error { return nil }
func (BaseVisitor) VisitType(*TypeDecl) error           { return nil }
func (BaseVisitor) VisitField(*FieldDecl) error         { return nil }
func (BaseVisitor) VisitRelation(*RelationDecl) error   { return nil }
func (BaseVisitor) VisitInterface(*InterfaceDecl) error { return nil }
func (BaseVisitor) VisitMethod(*MethodDecl) error       { return nil }
func (BaseVisitor) VisitFlow(*FlowDecl) error           { return nil }
func (BaseVisitor) VisitStep(StepDecl) error            { return nil }
func (BaseVisitor) VisitStates(*StatesDecl) error       { return nil }
func (BaseVisitor) VisitState(*StateDecl) error         { return nil }
func (BaseVisitor) VisitTransition(*TransitionDecl) error { return nil }
func (BaseVisitor) VisitExpr(ExprDecl) error            { return nil }
func (BaseVisitor) VisitAnnotation(*AnnotationDecl) error { return nil }
```

**責務**: AST の走査を抽象化

---

### internal/domain/config/config.go

```go
package config

// Config はプロジェクト設定を表す
type Config struct {
    SourceRoot string   `yaml:"source_root"`
    PactRoot   string   `yaml:"pact_root"`
    OutputDir  string   `yaml:"output_dir"`
    Language   string   `yaml:"language"`
    Diagrams   []string `yaml:"diagrams"`
    Exclude    []string `yaml:"exclude"`
}

// Default はデフォルト設定を返す
func Default() *Config {
    return &Config{
        SourceRoot: "./src",
        PactRoot:   "./.pact",
        OutputDir:  "./diagrams",
        Language:   "",
        Diagrams:   []string{"class", "sequence", "state", "flow"},
        Exclude:    []string{},
    }
}

// DiagramEnabled は指定した図が有効かどうかを返す
func (c *Config) DiagramEnabled(name string) bool {
    for _, d := range c.Diagrams {
        if d == name || d == "all" {
            return true
        }
    }
    return false
}

// IsExcluded はパスが除外対象かどうかを返す
func (c *Config) IsExcluded(path string) bool {
    // glob マッチングで判定
    // 実装は infrastructure 層に委譲
    return false
}
```

**責務**: プロジェクト設定のドメインモデル

---

### internal/domain/diagram/common/types.go

```go
package common

// Annotation は図要素に付与されるアノテーション
type Annotation struct {
    Name   string
    Values map[string]string
}

// Position は図内の座標
type Position struct {
    X int
    Y int
}

// Size はサイズ
type Size struct {
    Width  int
    Height int
}

// Bounds は位置とサイズを持つ矩形
type Bounds struct {
    Position
    Size
}

// LineStyle は線のスタイル
type LineStyle string

const (
    LineSolid  LineStyle = "solid"
    LineDashed LineStyle = "dashed"
    LineDotted LineStyle = "dotted"
)
```

**責務**: 図モデル共通の型定義

---

### internal/domain/diagram/class/model.go

```go
package class

import "pact/internal/domain/diagram/common"

// Diagram はクラス図を表す
type Diagram struct {
    Nodes []Node
    Edges []Edge
}

// Node はクラス図のノード（クラス、インターフェース等）
type Node struct {
    ID           string
    Name         string
    Stereotype   *Stereotype
    Compartments Compartments
    Annotations  []common.Annotation
}

type Stereotype string

const (
    StereotypeInterface Stereotype = "interface"
    StereotypeAbstract  Stereotype = "abstract"
    StereotypeEnum      Stereotype = "enum"
    StereotypeComponent Stereotype = "component"
)

// Compartments はクラスの区画（属性、メソッド）
type Compartments struct {
    Attributes []Attribute
    Methods    []Method
}

// Attribute は属性
type Attribute struct {
    ID         string
    Name       string
    Type       string
    Visibility Visibility
}

type Visibility string

const (
    Public    Visibility = "public"
    Private   Visibility = "private"
    Protected Visibility = "protected"
    Package   Visibility = "package"
)

// Method はメソッド
type Method struct {
    ID         string
    Name       string
    Params     []Param
    ReturnType string
    Visibility Visibility
    Modifiers  []Modifier
}

// Param はパラメータ
type Param struct {
    Name string
    Type string
}

type Modifier string

const (
    ModifierAsync    Modifier = "async"
    ModifierStatic   Modifier = "static"
    ModifierAbstract Modifier = "abstract"
)

// Edge はクラス図のエッジ（関係）
type Edge struct {
    ID               string
    Source           string
    Target           string
    Type             EdgeType
    Label            *string
    SourceDecoration Decoration
    TargetDecoration Decoration
    LineStyle        common.LineStyle
}

type EdgeType string

const (
    Dependency     EdgeType = "dependency"
    Association    EdgeType = "association"
    Aggregation    EdgeType = "aggregation"
    Composition    EdgeType = "composition"
    Inheritance    EdgeType = "inheritance"
    Implementation EdgeType = "implementation"
)

type Decoration string

const (
    DecorationNone          Decoration = "none"
    DecorationArrow         Decoration = "arrow"
    DecorationTriangle      Decoration = "triangle"
    DecorationDiamond       Decoration = "diamond"
    DecorationFilledDiamond Decoration = "filled_diamond"
)
```

**責務**: クラス図のドメインモデル

---

### internal/domain/diagram/sequence/model.go

```go
package sequence

import "pact/internal/domain/diagram/common"

// Diagram はシーケンス図を表す
type Diagram struct {
    Participants []Participant
    Events       []Event
}

// Participant は参加者
type Participant struct {
    ID          string
    Name        string
    Type        ParticipantType
    Annotations []common.Annotation
}

type ParticipantType string

const (
    TypeActor     ParticipantType = "actor"
    TypeComponent ParticipantType = "component"
    TypeDatabase  ParticipantType = "database"
    TypeQueue     ParticipantType = "queue"
    TypeExternal  ParticipantType = "external"
)

// Event はシーケンス図のイベント
type Event interface {
    eventNode()
}

// MessageEvent はメッセージ
type MessageEvent struct {
    ID          string
    From        string
    To          string
    Label       string
    MessageType MessageType
    Annotations []common.Annotation
}

func (MessageEvent) eventNode() {}

type MessageType string

const (
    MessageSync    MessageType = "sync"
    MessageAsync   MessageType = "async"
    MessageReturn  MessageType = "return"
    MessageCreate  MessageType = "create"
    MessageDestroy MessageType = "destroy"
)

// ActivationEvent はアクティベーション（活性区間）
type ActivationEvent struct {
    Participant string
    Action      ActivationAction
}

func (ActivationEvent) eventNode() {}

type ActivationAction string

const (
    ActivationStart ActivationAction = "start"
    ActivationEnd   ActivationAction = "end"
)

// FragmentEvent は複合フラグメント（alt, loop 等）
type FragmentEvent struct {
    FragmentType FragmentType
    Action       FragmentAction
    Condition    *string
    Annotations  []common.Annotation
}

func (FragmentEvent) eventNode() {}

type FragmentType string

const (
    FragmentAlt      FragmentType = "alt"
    FragmentOpt      FragmentType = "opt"
    FragmentLoop     FragmentType = "loop"
    FragmentPar      FragmentType = "par"
    FragmentCritical FragmentType = "critical"
    FragmentBreak    FragmentType = "break"
)

type FragmentAction string

const (
    FragmentStart FragmentAction = "start"
    FragmentElse  FragmentAction = "else"
    FragmentEnd   FragmentAction = "end"
)

// NoteEvent はノート
type NoteEvent struct {
    Position     NotePosition
    Participants []string
    Text         string
}

func (NoteEvent) eventNode() {}

type NotePosition string

const (
    NoteLeft  NotePosition = "left"
    NoteRight NotePosition = "right"
    NoteOver  NotePosition = "over"
)
```

**責務**: シーケンス図のドメインモデル

---

### internal/domain/diagram/state/model.go

```go
package state

import "pact/internal/domain/diagram/common"

// Diagram はステートマシン図を表す
type Diagram struct {
    ID          string
    Name        string
    States      []State
    Transitions []Transition
    Annotations []common.Annotation
}

// State は状態
type State struct {
    ID          string
    Name        string
    Type        StateType
    Children    []State
    Regions     []Region
    Entry       []Action
    Exit        []Action
    Initial     *string
    Annotations []common.Annotation
}

type StateType string

const (
    TypeAtomic      StateType = "atomic"
    TypeCompound    StateType = "compound"
    TypeParallel    StateType = "parallel"
    TypeInitial     StateType = "initial"
    TypeFinal       StateType = "final"
    TypeChoice      StateType = "choice"
    TypeJunction    StateType = "junction"
    TypeHistory     StateType = "history"
    TypeDeepHistory StateType = "deep_history"
)

// Region は並行状態のリージョン
type Region struct {
    ID          string
    Name        *string
    States      []State
    Transitions []Transition
}

// Transition は遷移
type Transition struct {
    ID          string
    Source      string
    Target      string
    Trigger     Trigger
    Guard       *string
    Actions     []Action
    Annotations []common.Annotation
}

// Trigger はトリガー
type Trigger interface {
    triggerType()
}

// EventTrigger はイベントトリガー
type EventTrigger struct {
    Name string
}

func (EventTrigger) triggerType() {}

// AfterTrigger は時間トリガー
type AfterTrigger struct {
    Duration Duration
}

func (AfterTrigger) triggerType() {}

// WhenTrigger は条件トリガー
type WhenTrigger struct {
    Condition string
}

func (WhenTrigger) triggerType() {}

// Duration は時間
type Duration struct {
    Value int
    Unit  DurationUnit
}

type DurationUnit string

const (
    UnitMs DurationUnit = "ms"
    UnitS  DurationUnit = "s"
    UnitM  DurationUnit = "m"
    UnitH  DurationUnit = "h"
    UnitD  DurationUnit = "d"
)

// Action はアクション
type Action struct {
    Name string
}
```

**責務**: ステートマシン図のドメインモデル

---

### internal/domain/diagram/flow/model.go

```go
package flow

import "pact/internal/domain/diagram/common"

// Diagram はフローチャートを表す
type Diagram struct {
    ID          string
    Name        string
    Nodes       []Node
    Edges       []Edge
    Swimlanes   []Swimlane
    Annotations []common.Annotation
}

// Node はフローチャートのノード
type Node struct {
    ID          string
    Label       string
    Shape       Shape
    Swimlane    *string
    Annotations []common.Annotation
}

type Shape string

const (
    ShapeTerminal    Shape = "terminal"
    ShapeProcess     Shape = "process"
    ShapeDecision    Shape = "decision"
    ShapeIO          Shape = "io"
    ShapePredefined  Shape = "predefined"
    ShapeConnector   Shape = "connector"
    ShapeDocument    Shape = "document"
    ShapeDatabase    Shape = "database"
    ShapeManual      Shape = "manual"
    ShapePreparation Shape = "preparation"
)

// Edge はフローチャートのエッジ
type Edge struct {
    ID          string
    Source      string
    Target      string
    Label       *string
    LineStyle   common.LineStyle
    Annotations []common.Annotation
}

// Swimlane はスイムレーン
type Swimlane struct {
    ID          string
    Name        string
    Order       int
    Annotations []common.Annotation
}
```

**責務**: フローチャートのドメインモデル

---

### internal/domain/errors/errors.go

```go
package errors

import (
    "fmt"

    "pact/internal/domain/ast"
)

// ParseError はパースエラー
type ParseError struct {
    Pos     ast.Position
    Message string
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("%s: %s", e.Pos, e.Message)
}

// SemanticError は意味解析エラー
type SemanticError struct {
    Pos     ast.Position
    Message string
}

func (e *SemanticError) Error() string {
    return fmt.Sprintf("%s: %s", e.Pos, e.Message)
}

// ImportError はインポート解決エラー
type ImportError struct {
    Pos     ast.Position
    Path    string
    Message string
}

func (e *ImportError) Error() string {
    return fmt.Sprintf("%s: import %q: %s", e.Pos, e.Path, e.Message)
}

// CycleError は循環参照エラー
type CycleError struct {
    Cycle []string
}

func (e *CycleError) Error() string {
    return fmt.Sprintf("import cycle detected: %v", e.Cycle)
}

// TransformError は変換エラー
type TransformError struct {
    Source  string
    Target  string
    Message string
}

func (e *TransformError) Error() string {
    return fmt.Sprintf("transform %s -> %s: %s", e.Source, e.Target, e.Message)
}

// ConfigError は設定エラー
type ConfigError struct {
    Path    string
    Message string
}

func (e *ConfigError) Error() string {
    if e.Path != "" {
        return fmt.Sprintf("config %s: %s", e.Path, e.Message)
    }
    return fmt.Sprintf("config: %s", e.Message)
}
```

**責務**: ドメイン固有エラーの定義

---

## internal/application/ - アプリケーション層

ユースケースの実装。ドメイン層を組み合わせる。

### internal/application/parser/service.go

```go
package parser

import "pact/internal/domain/ast"

// Parser はパース処理のインターフェース
type Parser interface {
    ParseFile(path string) (*ast.SpecFile, error)
    ParseString(content string, filename string) (*ast.SpecFile, error)
}

// Service はパースのアプリケーションサービス
type Service struct {
    parser Parser
}

// NewService は新しい Service を作成する
func NewService(parser Parser) *Service {
    return &Service{parser: parser}
}

// ParseFile は単一ファイルをパースする
func (s *Service) ParseFile(path string) (*ast.SpecFile, error) {
    return s.parser.ParseFile(path)
}

// ParseString は文字列をパースする
func (s *Service) ParseString(content, filename string) (*ast.SpecFile, error) {
    return s.parser.ParseString(content, filename)
}
```

**責務**: パースのユースケース

---

### internal/application/project/service.go

```go
package project

import (
    "path/filepath"

    "pact/internal/domain/ast"
    "pact/internal/domain/config"
    "pact/internal/application/parser"
)

// ImportResolver はインポート解決のインターフェース
type ImportResolver interface {
    Resolve(files []*ast.SpecFile) ([]*ast.SpecFile, error)
}

// FileFinder はファイル検索のインターフェース
type FileFinder interface {
    FindPactFiles(dir string, exclude []string) ([]string, error)
    FindSourceFiles(dir string, exclude []string) ([]string, error)
}

// Service はプロジェクト管理のアプリケーションサービス
type Service struct {
    config       *config.Config
    parserSvc    *parser.Service
    resolver     ImportResolver
    finder       FileFinder
}

// NewService は新しい Service を作成する
func NewService(
    cfg *config.Config,
    parserSvc *parser.Service,
    resolver ImportResolver,
    finder FileFinder,
) *Service {
    return &Service{
        config:    cfg,
        parserSvc: parserSvc,
        resolver:  resolver,
        finder:    finder,
    }
}

// LoadProject はプロジェクト全体を読み込む
func (s *Service) LoadProject() ([]*ast.SpecFile, error) {
    // 1. pact_root 以下の .pact ファイルを検索
    paths, err := s.finder.FindPactFiles(s.config.PactRoot, s.config.Exclude)
    if err != nil {
        return nil, err
    }

    // 2. 各ファイルをパース
    files := make([]*ast.SpecFile, 0, len(paths))
    for _, path := range paths {
        file, err := s.parserSvc.ParseFile(path)
        if err != nil {
            return nil, err
        }
        files = append(files, file)
    }

    // 3. import を解決してソート
    return s.resolver.Resolve(files)
}

// MissingSpec は対応する .pact がないソースファイル
type MissingSpec struct {
    SourcePath   string
    ExpectedPact string
}

// CheckMissing は対応する .pact がないソースファイルを検出する
func (s *Service) CheckMissing() ([]MissingSpec, error) {
    // 1. source_root 以下のソースファイルを検索
    sources, err := s.finder.FindSourceFiles(s.config.SourceRoot, s.config.Exclude)
    if err != nil {
        return nil, err
    }

    // 2. pact_root 以下の .pact ファイルを検索
    pacts, err := s.finder.FindPactFiles(s.config.PactRoot, s.config.Exclude)
    if err != nil {
        return nil, err
    }

    // 3. .pact ファイルの相対パスをセットに
    pactSet := make(map[string]bool)
    for _, p := range pacts {
        rel, _ := filepath.Rel(s.config.PactRoot, p)
        // 拡張子を除いた名前で比較
        name := rel[:len(rel)-len(filepath.Ext(rel))]
        pactSet[name] = true
    }

    // 4. 対応する .pact がないソースを検出
    var missing []MissingSpec
    for _, src := range sources {
        rel, _ := filepath.Rel(s.config.SourceRoot, src)
        name := rel[:len(rel)-len(filepath.Ext(rel))]
        if !pactSet[name] {
            missing = append(missing, MissingSpec{
                SourcePath:   src,
                ExpectedPact: filepath.Join(s.config.PactRoot, name+".pact"),
            })
        }
    }

    return missing, nil
}

// SourceForPact は .pact ファイルに対応するソースファイルのパスを返す
func (s *Service) SourceForPact(pactPath string) string {
    rel, _ := filepath.Rel(s.config.PactRoot, pactPath)
    name := rel[:len(rel)-len(filepath.Ext(rel))]
    // 言語に応じた拡張子を付与（簡略化）
    ext := ".go"
    switch s.config.Language {
    case "typescript", "ts":
        ext = ".ts"
    case "python", "py":
        ext = ".py"
    }
    return filepath.Join(s.config.SourceRoot, name+ext)
}
```

**責務**: プロジェクト全体の管理、ファイル対応の検証

---

### internal/application/transformer/service.go

```go
package transformer

import (
    "pact/internal/domain/ast"
    "pact/internal/domain/diagram/class"
    "pact/internal/domain/diagram/sequence"
    "pact/internal/domain/diagram/state"
    "pact/internal/domain/diagram/flow"
)

// Service は変換のアプリケーションサービス
type Service struct {
    classTransformer    *ClassTransformer
    sequenceTransformer *SequenceTransformer
    stateTransformer    *StateTransformer
    flowTransformer     *FlowTransformer
}

// NewService は新しい Service を作成する
func NewService() *Service {
    return &Service{
        classTransformer:    NewClassTransformer(),
        sequenceTransformer: NewSequenceTransformer(),
        stateTransformer:    NewStateTransformer(),
        flowTransformer:     NewFlowTransformer(),
    }
}

// ClassOptions はクラス図変換のオプション
type ClassOptions struct {
    Components        []string
    IncludeTypes      bool
    IncludeInterfaces bool
}

// SequenceOptions はシーケンス図変換のオプション
type SequenceOptions struct {
    IncludeReturnMessages bool
    IncludeActivations    bool
}

// StateOptions はステートマシン図変換のオプション
type StateOptions struct {
    ExpandHierarchy bool
}

// FlowOptions はフローチャート変換のオプション
type FlowOptions struct {
    IncludeSwimlanes bool
}

// ToClassDiagram は AST からクラス図を生成する
func (s *Service) ToClassDiagram(files []*ast.SpecFile, opts ClassOptions) (*class.Diagram, error) {
    return s.classTransformer.Transform(files, opts)
}

// ToSequenceDiagram は AST からシーケンス図を生成する
func (s *Service) ToSequenceDiagram(file *ast.SpecFile, flowName string, opts SequenceOptions) (*sequence.Diagram, error) {
    return s.sequenceTransformer.Transform(file, flowName, opts)
}

// ToStateDiagram は AST からステートマシン図を生成する
func (s *Service) ToStateDiagram(file *ast.SpecFile, statesName string, opts StateOptions) (*state.Diagram, error) {
    return s.stateTransformer.Transform(file, statesName, opts)
}

// ToFlowchart は AST からフローチャートを生成する
func (s *Service) ToFlowchart(file *ast.SpecFile, flowName string, opts FlowOptions) (*flow.Diagram, error) {
    return s.flowTransformer.Transform(file, flowName, opts)
}
```

**責務**: 変換ユースケースの統合

---

### internal/application/transformer/class.go

```go
package transformer

import (
    "pact/internal/domain/ast"
    "pact/internal/domain/diagram/class"
)

// ClassTransformer は AST からクラス図への変換を行う
type ClassTransformer struct{}

// NewClassTransformer は新しい ClassTransformer を作成する
func NewClassTransformer() *ClassTransformer {
    return &ClassTransformer{}
}

// Transform は AST からクラス図を生成する
func (t *ClassTransformer) Transform(files []*ast.SpecFile, opts ClassOptions) (*class.Diagram, error) {
    diagram := &class.Diagram{
        Nodes: []class.Node{},
        Edges: []class.Edge{},
    }

    for _, file := range files {
        comp := &file.Component

        // フィルタリング
        if len(opts.Components) > 0 && !contains(opts.Components, comp.Name) {
            continue
        }

        // Component -> Node
        diagram.Nodes = append(diagram.Nodes, t.componentToNode(comp))

        // Types -> Nodes
        if opts.IncludeTypes {
            for _, typ := range comp.Body.Types {
                diagram.Nodes = append(diagram.Nodes, t.typeToNode(&typ))
            }
        }

        // Relations -> Edges
        for _, rel := range comp.Body.Relations {
            diagram.Edges = append(diagram.Edges, t.relationToEdge(comp.Name, &rel))
        }

        // Requires -> Nodes (interface)
        if opts.IncludeInterfaces {
            for _, iface := range comp.Body.Requires {
                diagram.Nodes = append(diagram.Nodes, t.interfaceToNode(&iface))
            }
        }
    }

    return diagram, nil
}

func (t *ClassTransformer) componentToNode(comp *ast.ComponentDecl) class.Node {
    node := class.Node{
        ID:   comp.Name,
        Name: comp.Name,
    }

    // Provides からメソッドを収集
    for _, iface := range comp.Body.Provides {
        for _, method := range iface.Methods {
            node.Compartments.Methods = append(node.Compartments.Methods, class.Method{
                ID:         method.Name,
                Name:       method.Name,
                ReturnType: method.Returns.Name,
                Visibility: class.Public,
            })
        }
    }

    return node
}

func (t *ClassTransformer) typeToNode(typ *ast.TypeDecl) class.Node {
    node := class.Node{
        ID:   typ.Name,
        Name: typ.Name,
    }

    if typ.Kind == ast.TypeKindEnum {
        stereo := class.StereotypeEnum
        node.Stereotype = &stereo
    }

    for _, field := range typ.Fields {
        node.Compartments.Attributes = append(node.Compartments.Attributes, class.Attribute{
            ID:         field.Name,
            Name:       field.Name,
            Type:       field.Type.Name,
            Visibility: t.convertVisibility(field.Visibility),
        })
    }

    return node
}

func (t *ClassTransformer) relationToEdge(source string, rel *ast.RelationDecl) class.Edge {
    edge := class.Edge{
        ID:     source + "-" + rel.Target,
        Source: source,
        Target: rel.Target,
    }

    switch rel.Kind {
    case ast.RelationDepends:
        edge.Type = class.Dependency
        edge.LineStyle = "dashed"
        edge.TargetDecoration = class.DecorationArrow
    case ast.RelationExtends:
        edge.Type = class.Inheritance
        edge.LineStyle = "solid"
        edge.TargetDecoration = class.DecorationTriangle
    case ast.RelationImplements:
        edge.Type = class.Implementation
        edge.LineStyle = "dashed"
        edge.TargetDecoration = class.DecorationTriangle
    case ast.RelationContains:
        edge.Type = class.Composition
        edge.LineStyle = "solid"
        edge.SourceDecoration = class.DecorationFilledDiamond
    case ast.RelationAggregates:
        edge.Type = class.Aggregation
        edge.LineStyle = "solid"
        edge.SourceDecoration = class.DecorationDiamond
    }

    return edge
}

func (t *ClassTransformer) interfaceToNode(iface *ast.InterfaceDecl) class.Node {
    stereo := class.StereotypeInterface
    node := class.Node{
        ID:         iface.Name,
        Name:       iface.Name,
        Stereotype: &stereo,
    }

    for _, method := range iface.Methods {
        node.Compartments.Methods = append(node.Compartments.Methods, class.Method{
            ID:         method.Name,
            Name:       method.Name,
            ReturnType: method.Returns.Name,
            Visibility: class.Public,
        })
    }

    return node
}

func (t *ClassTransformer) convertVisibility(v ast.Visibility) class.Visibility {
    switch v {
    case ast.VisibilityPrivate:
        return class.Private
    case ast.VisibilityProtected:
        return class.Protected
    case ast.VisibilityPackage:
        return class.Package
    default:
        return class.Public
    }
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

**責務**: AST からクラス図への変換

---

### internal/application/transformer/sequence.go

```go
package transformer

import (
    "fmt"

    "pact/internal/domain/ast"
    "pact/internal/domain/diagram/sequence"
    "pact/internal/domain/errors"
)

// SequenceTransformer は AST からシーケンス図への変換を行う
type SequenceTransformer struct{}

// NewSequenceTransformer は新しい SequenceTransformer を作成する
func NewSequenceTransformer() *SequenceTransformer {
    return &SequenceTransformer{}
}

// Transform は AST からシーケンス図を生成する
func (t *SequenceTransformer) Transform(file *ast.SpecFile, flowName string, opts SequenceOptions) (*sequence.Diagram, error) {
    comp := &file.Component

    // Flow を探す
    var flow *ast.FlowDecl
    for i := range comp.Body.Flows {
        if comp.Body.Flows[i].Name == flowName {
            flow = &comp.Body.Flows[i]
            break
        }
    }
    if flow == nil {
        return nil, &errors.TransformError{
            Source:  file.Path,
            Target:  "sequence",
            Message: fmt.Sprintf("flow %q not found", flowName),
        }
    }

    diagram := &sequence.Diagram{
        Participants: []sequence.Participant{},
        Events:       []sequence.Event{},
    }

    // 自コンポーネントを参加者として追加
    diagram.Participants = append(diagram.Participants, sequence.Participant{
        ID:   comp.Name,
        Name: comp.Name,
        Type: sequence.TypeComponent,
    })

    // 依存先を参加者として追加
    for _, rel := range comp.Body.Relations {
        if rel.Kind == ast.RelationDepends {
            diagram.Participants = append(diagram.Participants, sequence.Participant{
                ID:   rel.Target,
                Name: rel.Target,
                Type: t.convertTargetType(rel.TargetType),
            })
        }
    }

    // Steps を走査してイベントを生成
    t.processSteps(comp.Name, flow.Steps, diagram, opts)

    return diagram, nil
}

func (t *SequenceTransformer) processSteps(self string, steps []ast.StepDecl, diagram *sequence.Diagram, opts SequenceOptions) {
    for _, step := range steps {
        switch s := step.(type) {
        case *ast.AssignStep:
            if call, ok := s.Value.(*ast.CallExpr); ok {
                // メッセージを追加
                diagram.Events = append(diagram.Events, sequence.MessageEvent{
                    ID:          fmt.Sprintf("msg-%d", len(diagram.Events)),
                    From:        self,
                    To:          call.Target,
                    Label:       call.Method,
                    MessageType: sequence.MessageSync,
                })

                // return メッセージ
                if opts.IncludeReturnMessages {
                    diagram.Events = append(diagram.Events, sequence.MessageEvent{
                        ID:          fmt.Sprintf("msg-%d", len(diagram.Events)),
                        From:        call.Target,
                        To:          self,
                        Label:       s.Variable,
                        MessageType: sequence.MessageReturn,
                    })
                }
            }

        case *ast.CallStep:
            msgType := sequence.MessageSync
            if s.Async {
                msgType = sequence.MessageAsync
            }
            diagram.Events = append(diagram.Events, sequence.MessageEvent{
                ID:          fmt.Sprintf("msg-%d", len(diagram.Events)),
                From:        self,
                To:          s.Target,
                Label:       s.Method,
                MessageType: msgType,
            })

        case *ast.IfStep:
            // alt フラグメント開始
            cond := "condition"
            diagram.Events = append(diagram.Events, sequence.FragmentEvent{
                FragmentType: sequence.FragmentAlt,
                Action:       sequence.FragmentStart,
                Condition:    &cond,
            })

            // then ブロック
            t.processSteps(self, s.Then, diagram, opts)

            // else ブロック
            if len(s.Else) > 0 {
                diagram.Events = append(diagram.Events, sequence.FragmentEvent{
                    FragmentType: sequence.FragmentAlt,
                    Action:       sequence.FragmentElse,
                })
                t.processSteps(self, s.Else, diagram, opts)
            }

            // alt フラグメント終了
            diagram.Events = append(diagram.Events, sequence.FragmentEvent{
                FragmentType: sequence.FragmentAlt,
                Action:       sequence.FragmentEnd,
            })

        case *ast.ForStep:
            // loop フラグメント
            cond := fmt.Sprintf("for %s in ...", s.Variable)
            diagram.Events = append(diagram.Events, sequence.FragmentEvent{
                FragmentType: sequence.FragmentLoop,
                Action:       sequence.FragmentStart,
                Condition:    &cond,
            })
            t.processSteps(self, s.Body, diagram, opts)
            diagram.Events = append(diagram.Events, sequence.FragmentEvent{
                FragmentType: sequence.FragmentLoop,
                Action:       sequence.FragmentEnd,
            })

        case *ast.WhileStep:
            // loop フラグメント
            cond := "while condition"
            diagram.Events = append(diagram.Events, sequence.FragmentEvent{
                FragmentType: sequence.FragmentLoop,
                Action:       sequence.FragmentStart,
                Condition:    &cond,
            })
            t.processSteps(self, s.Body, diagram, opts)
            diagram.Events = append(diagram.Events, sequence.FragmentEvent{
                FragmentType: sequence.FragmentLoop,
                Action:       sequence.FragmentEnd,
            })
        }
    }
}

func (t *SequenceTransformer) convertTargetType(tt ast.TargetType) sequence.ParticipantType {
    switch tt {
    case ast.TargetDatabase:
        return sequence.TypeDatabase
    case ast.TargetExternal:
        return sequence.TypeExternal
    case ast.TargetQueue:
        return sequence.TypeQueue
    case ast.TargetActor:
        return sequence.TypeActor
    default:
        return sequence.TypeComponent
    }
}
```

**責務**: AST からシーケンス図への変換

---

### internal/application/transformer/state.go

```go
package transformer

import (
    "fmt"

    "pact/internal/domain/ast"
    "pact/internal/domain/diagram/state"
    "pact/internal/domain/errors"
)

// StateTransformer は AST からステートマシン図への変換を行う
type StateTransformer struct{}

// NewStateTransformer は新しい StateTransformer を作成する
func NewStateTransformer() *StateTransformer {
    return &StateTransformer{}
}

// Transform は AST からステートマシン図を生成する
func (t *StateTransformer) Transform(file *ast.SpecFile, statesName string, opts StateOptions) (*state.Diagram, error) {
    comp := &file.Component

    // States を探す
    var states *ast.StatesDecl
    for i := range comp.Body.States {
        if comp.Body.States[i].Name == statesName {
            states = &comp.Body.States[i]
            break
        }
    }
    if states == nil {
        return nil, &errors.TransformError{
            Source:  file.Path,
            Target:  "state",
            Message: fmt.Sprintf("states %q not found", statesName),
        }
    }

    diagram := &state.Diagram{
        ID:          states.Name,
        Name:        states.Name,
        States:      []state.State{},
        Transitions: []state.Transition{},
    }

    // 状態を収集（遷移から + 明示的定義から）
    stateMap := make(map[string]*state.State)

    // initial 状態
    stateMap[states.Initial] = &state.State{
        ID:   states.Initial,
        Name: states.Initial,
        Type: state.TypeAtomic,
    }

    // final 状態
    for _, f := range states.Finals {
        stateMap[f] = &state.State{
            ID:   f,
            Name: f,
            Type: state.TypeFinal,
        }
    }

    // 明示的に定義された状態
    for _, s := range states.States {
        stateMap[s.Name] = t.convertState(&s)
    }

    // 遷移から状態を収集
    for _, tr := range states.Transitions {
        if _, ok := stateMap[tr.From]; !ok {
            stateMap[tr.From] = &state.State{
                ID:   tr.From,
                Name: tr.From,
                Type: state.TypeAtomic,
            }
        }
        if _, ok := stateMap[tr.To]; !ok {
            stateMap[tr.To] = &state.State{
                ID:   tr.To,
                Name: tr.To,
                Type: state.TypeAtomic,
            }
        }
    }

    // initial 状態をマーク
    if s, ok := stateMap[states.Initial]; ok {
        if s.Type == state.TypeAtomic {
            // 別途 initial 疑似状態を追加
            diagram.States = append(diagram.States, state.State{
                ID:   "__initial__",
                Name: "",
                Type: state.TypeInitial,
            })
            diagram.Transitions = append(diagram.Transitions, state.Transition{
                ID:     "__initial__-" + states.Initial,
                Source: "__initial__",
                Target: states.Initial,
            })
        }
    }

    // 状態をダイアグラムに追加
    for _, s := range stateMap {
        diagram.States = append(diagram.States, *s)
    }

    // 遷移を変換
    for _, tr := range states.Transitions {
        diagram.Transitions = append(diagram.Transitions, t.convertTransition(&tr))
    }

    return diagram, nil
}

func (t *StateTransformer) convertState(s *ast.StateDecl) *state.State {
    result := &state.State{
        ID:   s.Name,
        Name: s.Name,
        Type: t.convertStateKind(s.Kind),
    }

    for _, action := range s.Entry {
        result.Entry = append(result.Entry, state.Action{Name: action})
    }
    for _, action := range s.Exit {
        result.Exit = append(result.Exit, state.Action{Name: action})
    }

    // 子状態
    for _, child := range s.States {
        result.Children = append(result.Children, *t.convertState(&child))
    }

    // リージョン
    for _, region := range s.Regions {
        result.Regions = append(result.Regions, t.convertRegion(&region))
    }

    if s.Initial != nil {
        result.Initial = s.Initial
    }

    return result
}

func (t *StateTransformer) convertRegion(r *ast.RegionDecl) state.Region {
    region := state.Region{
        ID:   r.Name,
        Name: &r.Name,
    }

    for _, s := range r.States {
        region.States = append(region.States, *t.convertState(&s))
    }

    for _, tr := range r.Transitions {
        region.Transitions = append(region.Transitions, t.convertTransition(&tr))
    }

    return region
}

func (t *StateTransformer) convertTransition(tr *ast.TransitionDecl) state.Transition {
    result := state.Transition{
        ID:     tr.From + "-" + tr.To,
        Source: tr.From,
        Target: tr.To,
    }

    if tr.Trigger != nil {
        switch trig := tr.Trigger.(type) {
        case *ast.EventTrigger:
            result.Trigger = state.EventTrigger{Name: trig.Name}
        case *ast.AfterTrigger:
            result.Trigger = state.AfterTrigger{
                Duration: state.Duration{
                    Value: trig.Value,
                    Unit:  state.DurationUnit(trig.Unit),
                },
            }
        case *ast.WhenTrigger:
            result.Trigger = state.WhenTrigger{Condition: "condition"}
        }
    }

    for _, action := range tr.Actions {
        result.Actions = append(result.Actions, state.Action{Name: action})
    }

    return result
}

func (t *StateTransformer) convertStateKind(k ast.StateKind) state.StateType {
    switch k {
    case ast.StateCompound:
        return state.TypeCompound
    case ast.StateParallel:
        return state.TypeParallel
    default:
        return state.TypeAtomic
    }
}
```

**責務**: AST からステートマシン図への変換

---

### internal/application/transformer/flow.go

```go
package transformer

import (
    "fmt"

    "pact/internal/domain/ast"
    "pact/internal/domain/diagram/flow"
    "pact/internal/domain/errors"
)

// FlowTransformer は AST からフローチャートへの変換を行う
type FlowTransformer struct {
    nodeID int
}

// NewFlowTransformer は新しい FlowTransformer を作成する
func NewFlowTransformer() *FlowTransformer {
    return &FlowTransformer{}
}

// Transform は AST からフローチャートを生成する
func (t *FlowTransformer) Transform(file *ast.SpecFile, flowName string, opts FlowOptions) (*flow.Diagram, error) {
    comp := &file.Component

    // Flow を探す
    var flowDecl *ast.FlowDecl
    for i := range comp.Body.Flows {
        if comp.Body.Flows[i].Name == flowName {
            flowDecl = &comp.Body.Flows[i]
            break
        }
    }
    if flowDecl == nil {
        return nil, &errors.TransformError{
            Source:  file.Path,
            Target:  "flow",
            Message: fmt.Sprintf("flow %q not found", flowName),
        }
    }

    t.nodeID = 0
    diagram := &flow.Diagram{
        ID:        flowDecl.Name,
        Name:      flowDecl.Name,
        Nodes:     []flow.Node{},
        Edges:     []flow.Edge{},
        Swimlanes: []flow.Swimlane{},
    }

    // Start ノード
    startID := t.nextID("start")
    diagram.Nodes = append(diagram.Nodes, flow.Node{
        ID:    startID,
        Label: "Start",
        Shape: flow.ShapeTerminal,
    })

    // Steps を処理
    lastID := t.processSteps(startID, flowDecl.Steps, diagram, opts)

    // End ノード（最後のノードが return/throw でなければ追加）
    if lastID != "" {
        endID := t.nextID("end")
        diagram.Nodes = append(diagram.Nodes, flow.Node{
            ID:    endID,
            Label: "End",
            Shape: flow.ShapeTerminal,
        })
        diagram.Edges = append(diagram.Edges, flow.Edge{
            ID:     lastID + "-" + endID,
            Source: lastID,
            Target: endID,
        })
    }

    // スイムレーンを推論
    if opts.IncludeSwimlanes {
        t.inferSwimlanes(comp, diagram)
    }

    return diagram, nil
}

func (t *FlowTransformer) processSteps(prevID string, steps []ast.StepDecl, diagram *flow.Diagram, opts FlowOptions) string {
    currentID := prevID

    for _, step := range steps {
        switch s := step.(type) {
        case *ast.AssignStep:
            nodeID := t.nextID("process")
            label := fmt.Sprintf("%s = ...", s.Variable)
            if call, ok := s.Value.(*ast.CallExpr); ok {
                label = fmt.Sprintf("%s = %s.%s()", s.Variable, call.Target, call.Method)
            }
            diagram.Nodes = append(diagram.Nodes, flow.Node{
                ID:    nodeID,
                Label: label,
                Shape: flow.ShapeProcess,
            })
            diagram.Edges = append(diagram.Edges, flow.Edge{
                ID:     currentID + "-" + nodeID,
                Source: currentID,
                Target: nodeID,
            })
            currentID = nodeID

        case *ast.CallStep:
            nodeID := t.nextID("process")
            label := fmt.Sprintf("%s.%s()", s.Target, s.Method)
            diagram.Nodes = append(diagram.Nodes, flow.Node{
                ID:    nodeID,
                Label: label,
                Shape: flow.ShapeProcess,
            })
            diagram.Edges = append(diagram.Edges, flow.Edge{
                ID:     currentID + "-" + nodeID,
                Source: currentID,
                Target: nodeID,
            })
            currentID = nodeID

        case *ast.IfStep:
            decisionID := t.nextID("decision")
            diagram.Nodes = append(diagram.Nodes, flow.Node{
                ID:    decisionID,
                Label: "condition?",
                Shape: flow.ShapeDecision,
            })
            diagram.Edges = append(diagram.Edges, flow.Edge{
                ID:     currentID + "-" + decisionID,
                Source: currentID,
                Target: decisionID,
            })

            // then ブロック
            thenEndID := t.processSteps(decisionID, s.Then, diagram, opts)

            // else ブロック
            var elseEndID string
            if len(s.Else) > 0 {
                elseEndID = t.processSteps(decisionID, s.Else, diagram, opts)
            }

            // マージポイント
            mergeID := t.nextID("merge")
            diagram.Nodes = append(diagram.Nodes, flow.Node{
                ID:    mergeID,
                Label: "",
                Shape: flow.ShapeConnector,
            })

            if thenEndID != "" {
                yesLabel := "Yes"
                diagram.Edges = append(diagram.Edges, flow.Edge{
                    ID:     thenEndID + "-" + mergeID,
                    Source: thenEndID,
                    Target: mergeID,
                    Label:  &yesLabel,
                })
            }
            if elseEndID != "" {
                noLabel := "No"
                diagram.Edges = append(diagram.Edges, flow.Edge{
                    ID:     elseEndID + "-" + mergeID,
                    Source: elseEndID,
                    Target: mergeID,
                    Label:  &noLabel,
                })
            } else {
                noLabel := "No"
                diagram.Edges = append(diagram.Edges, flow.Edge{
                    ID:     decisionID + "-" + mergeID + "-no",
                    Source: decisionID,
                    Target: mergeID,
                    Label:  &noLabel,
                })
            }

            currentID = mergeID

        case *ast.ForStep:
            decisionID := t.nextID("decision")
            label := fmt.Sprintf("for %s in ...?", s.Variable)
            diagram.Nodes = append(diagram.Nodes, flow.Node{
                ID:    decisionID,
                Label: label,
                Shape: flow.ShapeDecision,
            })
            diagram.Edges = append(diagram.Edges, flow.Edge{
                ID:     currentID + "-" + decisionID,
                Source: currentID,
                Target: decisionID,
            })

            // ループ本体
            bodyEndID := t.processSteps(decisionID, s.Body, diagram, opts)

            // ループバック
            if bodyEndID != "" {
                diagram.Edges = append(diagram.Edges, flow.Edge{
                    ID:     bodyEndID + "-" + decisionID,
                    Source: bodyEndID,
                    Target: decisionID,
                })
            }

            currentID = decisionID

        case *ast.WhileStep:
            decisionID := t.nextID("decision")
            diagram.Nodes = append(diagram.Nodes, flow.Node{
                ID:    decisionID,
                Label: "while condition?",
                Shape: flow.ShapeDecision,
            })
            diagram.Edges = append(diagram.Edges, flow.Edge{
                ID:     currentID + "-" + decisionID,
                Source: currentID,
                Target: decisionID,
            })

            bodyEndID := t.processSteps(decisionID, s.Body, diagram, opts)

            if bodyEndID != "" {
                diagram.Edges = append(diagram.Edges, flow.Edge{
                    ID:     bodyEndID + "-" + decisionID,
                    Source: bodyEndID,
                    Target: decisionID,
                })
            }

            currentID = decisionID

        case *ast.ReturnStep:
            nodeID := t.nextID("return")
            diagram.Nodes = append(diagram.Nodes, flow.Node{
                ID:    nodeID,
                Label: "return",
                Shape: flow.ShapeTerminal,
            })
            diagram.Edges = append(diagram.Edges, flow.Edge{
                ID:     currentID + "-" + nodeID,
                Source: currentID,
                Target: nodeID,
            })
            return "" // 終端

        case *ast.ThrowStep:
            nodeID := t.nextID("throw")
            diagram.Nodes = append(diagram.Nodes, flow.Node{
                ID:    nodeID,
                Label: fmt.Sprintf("throw %s", s.Error),
                Shape: flow.ShapeTerminal,
            })
            diagram.Edges = append(diagram.Edges, flow.Edge{
                ID:     currentID + "-" + nodeID,
                Source: currentID,
                Target: nodeID,
            })
            return "" // 終端
        }
    }

    return currentID
}

func (t *FlowTransformer) inferSwimlanes(comp *ast.ComponentDecl, diagram *flow.Diagram) {
    // 依存先をスイムレーンとして追加
    order := 0
    swimlaneMap := make(map[string]bool)

    // 自コンポーネント
    diagram.Swimlanes = append(diagram.Swimlanes, flow.Swimlane{
        ID:    comp.Name,
        Name:  comp.Name,
        Order: order,
    })
    swimlaneMap[comp.Name] = true
    order++

    // 依存先
    for _, rel := range comp.Body.Relations {
        if rel.Kind == ast.RelationDepends && !swimlaneMap[rel.Target] {
            diagram.Swimlanes = append(diagram.Swimlanes, flow.Swimlane{
                ID:    rel.Target,
                Name:  rel.Target,
                Order: order,
            })
            swimlaneMap[rel.Target] = true
            order++
        }
    }
}

func (t *FlowTransformer) nextID(prefix string) string {
    t.nodeID++
    return fmt.Sprintf("%s-%d", prefix, t.nodeID)
}
```

**責務**: AST からフローチャートへの変換

---

### internal/application/renderer/service.go

```go
package renderer

import (
    "io"

    "pact/internal/domain/diagram/class"
    "pact/internal/domain/diagram/flow"
    "pact/internal/domain/diagram/sequence"
    "pact/internal/domain/diagram/state"
)

// Renderer はレンダリングのインターフェース
type Renderer interface {
    RenderClass(diagram *class.Diagram, w io.Writer) error
    RenderSequence(diagram *sequence.Diagram, w io.Writer) error
    RenderState(diagram *state.Diagram, w io.Writer) error
    RenderFlow(diagram *flow.Diagram, w io.Writer) error
}

// Service はレンダリングのアプリケーションサービス
type Service struct {
    renderer Renderer
}

// NewService は新しい Service を作成する
func NewService(renderer Renderer) *Service {
    return &Service{renderer: renderer}
}

// RenderClass はクラス図をレンダリングする
func (s *Service) RenderClass(diagram *class.Diagram, w io.Writer) error {
    return s.renderer.RenderClass(diagram, w)
}

// RenderSequence はシーケンス図をレンダリングする
func (s *Service) RenderSequence(diagram *sequence.Diagram, w io.Writer) error {
    return s.renderer.RenderSequence(diagram, w)
}

// RenderState はステートマシン図をレンダリングする
func (s *Service) RenderState(diagram *state.Diagram, w io.Writer) error {
    return s.renderer.RenderState(diagram, w)
}

// RenderFlow はフローチャートをレンダリングする
func (s *Service) RenderFlow(diagram *flow.Diagram, w io.Writer) error {
    return s.renderer.RenderFlow(diagram, w)
}
```

**責務**: レンダリングのユースケース

---

## internal/infrastructure/ - インフラストラクチャ層

外部依存の実装。

### internal/infrastructure/config/loader.go

```go
package config

import (
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"

    "pact/internal/domain/config"
    "pact/internal/domain/errors"
)

const ConfigFileName = ".pactconfig"

// Loader は設定ファイルを読み込む
type Loader struct{}

// NewLoader は新しい Loader を作成する
func NewLoader() *Loader {
    return &Loader{}
}

// Load は設定ファイルを読み込む
func (l *Loader) Load(dir string) (*config.Config, error) {
    path := filepath.Join(dir, ConfigFileName)

    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return config.Default(), nil
        }
        return nil, &errors.ConfigError{
            Path:    path,
            Message: err.Error(),
        }
    }

    cfg := config.Default()
    if err := yaml.Unmarshal(data, cfg); err != nil {
        return nil, &errors.ConfigError{
            Path:    path,
            Message: err.Error(),
        }
    }

    return cfg, nil
}

// Save は設定ファイルを保存する
func (l *Loader) Save(dir string, cfg *config.Config) error {
    path := filepath.Join(dir, ConfigFileName)

    data, err := yaml.Marshal(cfg)
    if err != nil {
        return &errors.ConfigError{
            Path:    path,
            Message: err.Error(),
        }
    }

    return os.WriteFile(path, data, 0644)
}

// FindProjectRoot はプロジェクトルートを探す
func (l *Loader) FindProjectRoot(start string) (string, error) {
    dir, err := filepath.Abs(start)
    if err != nil {
        return "", err
    }

    for {
        configPath := filepath.Join(dir, ConfigFileName)
        if _, err := os.Stat(configPath); err == nil {
            return dir, nil
        }

        parent := filepath.Dir(dir)
        if parent == dir {
            return "", &errors.ConfigError{
                Message: "project root not found (no .pactconfig)",
            }
        }
        dir = parent
    }
}
```

**責務**: 設定ファイルの読み書き

---

### internal/infrastructure/parser/token.go

```go
package parser

// TokenType はトークンの種類
type TokenType int

const (
    TOKEN_EOF TokenType = iota
    TOKEN_ILLEGAL

    // リテラル
    TOKEN_IDENT
    TOKEN_STRING
    TOKEN_INT
    TOKEN_FLOAT
    TOKEN_DURATION

    // キーワード
    TOKEN_IMPORT
    TOKEN_AS
    TOKEN_COMPONENT
    TOKEN_TYPE
    TOKEN_ENUM
    TOKEN_DEPENDS
    TOKEN_ON
    TOKEN_EXTENDS
    TOKEN_IMPLEMENTS
    TOKEN_CONTAINS
    TOKEN_AGGREGATES
    TOKEN_PROVIDES
    TOKEN_REQUIRES
    TOKEN_FLOW
    TOKEN_STATES
    TOKEN_STATE
    TOKEN_PARALLEL
    TOKEN_REGION
    TOKEN_INITIAL
    TOKEN_FINAL
    TOKEN_ENTRY
    TOKEN_EXIT
    TOKEN_IF
    TOKEN_ELSE
    TOKEN_FOR
    TOKEN_IN
    TOKEN_WHILE
    TOKEN_RETURN
    TOKEN_THROW
    TOKEN_AWAIT
    TOKEN_ASYNC
    TOKEN_THROWS
    TOKEN_WHEN
    TOKEN_AFTER
    TOKEN_DO
    TOKEN_TRUE
    TOKEN_FALSE
    TOKEN_NULL

    // 演算子・記号
    TOKEN_LBRACE
    TOKEN_RBRACE
    TOKEN_LPAREN
    TOKEN_RPAREN
    TOKEN_LBRACKET
    TOKEN_RBRACKET
    TOKEN_COLON
    TOKEN_COMMA
    TOKEN_DOT
    TOKEN_ARROW
    TOKEN_AT
    TOKEN_QUESTION
    TOKEN_PLUS
    TOKEN_MINUS
    TOKEN_STAR
    TOKEN_SLASH
    TOKEN_PERCENT
    TOKEN_EQ
    TOKEN_NE
    TOKEN_LT
    TOKEN_GT
    TOKEN_LE
    TOKEN_GE
    TOKEN_AND
    TOKEN_OR
    TOKEN_NOT
    TOKEN_ASSIGN
    TOKEN_NULLISH
    TOKEN_HASH
    TOKEN_TILDE
)

// Token はトークン
type Token struct {
    Type    TokenType
    Literal string
    Line    int
    Column  int
}

var keywords = map[string]TokenType{
    "import":     TOKEN_IMPORT,
    "as":         TOKEN_AS,
    "component":  TOKEN_COMPONENT,
    "type":       TOKEN_TYPE,
    "enum":       TOKEN_ENUM,
    "depends":    TOKEN_DEPENDS,
    "on":         TOKEN_ON,
    "extends":    TOKEN_EXTENDS,
    "implements": TOKEN_IMPLEMENTS,
    "contains":   TOKEN_CONTAINS,
    "aggregates": TOKEN_AGGREGATES,
    "provides":   TOKEN_PROVIDES,
    "requires":   TOKEN_REQUIRES,
    "flow":       TOKEN_FLOW,
    "states":     TOKEN_STATES,
    "state":      TOKEN_STATE,
    "parallel":   TOKEN_PARALLEL,
    "region":     TOKEN_REGION,
    "initial":    TOKEN_INITIAL,
    "final":      TOKEN_FINAL,
    "entry":      TOKEN_ENTRY,
    "exit":       TOKEN_EXIT,
    "if":         TOKEN_IF,
    "else":       TOKEN_ELSE,
    "for":        TOKEN_FOR,
    "in":         TOKEN_IN,
    "while":      TOKEN_WHILE,
    "return":     TOKEN_RETURN,
    "throw":      TOKEN_THROW,
    "await":      TOKEN_AWAIT,
    "async":      TOKEN_ASYNC,
    "throws":     TOKEN_THROWS,
    "when":       TOKEN_WHEN,
    "after":      TOKEN_AFTER,
    "do":         TOKEN_DO,
    "true":       TOKEN_TRUE,
    "false":      TOKEN_FALSE,
    "null":       TOKEN_NULL,
}

// LookupIdent は識別子がキーワードかどうかを判定する
func LookupIdent(ident string) TokenType {
    if tok, ok := keywords[ident]; ok {
        return tok
    }
    return TOKEN_IDENT
}
```

**責務**: トークンの種類と構造を定義

---

### internal/infrastructure/parser/lexer.go

```go
package parser

// Lexer は字句解析器
type Lexer struct {
    input   string
    file    string
    pos     int
    readPos int
    ch      byte
    line    int
    column  int
}

// NewLexer は新しい Lexer を作成する
func NewLexer(input string, file string) *Lexer {
    l := &Lexer{
        input:  input,
        file:   file,
        line:   1,
        column: 0,
    }
    l.readChar()
    return l
}

// NextToken は次のトークンを返す
func (l *Lexer) NextToken() Token {
    l.skipWhitespace()
    l.skipComment()
    l.skipWhitespace()

    tok := Token{
        Line:   l.line,
        Column: l.column,
    }

    switch l.ch {
    case '{':
        tok.Type = TOKEN_LBRACE
        tok.Literal = "{"
    case '}':
        tok.Type = TOKEN_RBRACE
        tok.Literal = "}"
    case '(':
        tok.Type = TOKEN_LPAREN
        tok.Literal = "("
    case ')':
        tok.Type = TOKEN_RPAREN
        tok.Literal = ")"
    case '[':
        tok.Type = TOKEN_LBRACKET
        tok.Literal = "["
    case ']':
        tok.Type = TOKEN_RBRACKET
        tok.Literal = "]"
    case ':':
        tok.Type = TOKEN_COLON
        tok.Literal = ":"
    case ',':
        tok.Type = TOKEN_COMMA
        tok.Literal = ","
    case '.':
        tok.Type = TOKEN_DOT
        tok.Literal = "."
    case '@':
        tok.Type = TOKEN_AT
        tok.Literal = "@"
    case '+':
        tok.Type = TOKEN_PLUS
        tok.Literal = "+"
    case '*':
        tok.Type = TOKEN_STAR
        tok.Literal = "*"
    case '/':
        tok.Type = TOKEN_SLASH
        tok.Literal = "/"
    case '%':
        tok.Type = TOKEN_PERCENT
        tok.Literal = "%"
    case '#':
        tok.Type = TOKEN_HASH
        tok.Literal = "#"
    case '~':
        tok.Type = TOKEN_TILDE
        tok.Literal = "~"
    case '-':
        if l.peekChar() == '>' {
            l.readChar()
            tok.Type = TOKEN_ARROW
            tok.Literal = "->"
        } else {
            tok.Type = TOKEN_MINUS
            tok.Literal = "-"
        }
    case '?':
        if l.peekChar() == '?' {
            l.readChar()
            tok.Type = TOKEN_NULLISH
            tok.Literal = "??"
        } else {
            tok.Type = TOKEN_QUESTION
            tok.Literal = "?"
        }
    case '=':
        if l.peekChar() == '=' {
            l.readChar()
            tok.Type = TOKEN_EQ
            tok.Literal = "=="
        } else {
            tok.Type = TOKEN_ASSIGN
            tok.Literal = "="
        }
    case '!':
        if l.peekChar() == '=' {
            l.readChar()
            tok.Type = TOKEN_NE
            tok.Literal = "!="
        } else {
            tok.Type = TOKEN_NOT
            tok.Literal = "!"
        }
    case '<':
        if l.peekChar() == '=' {
            l.readChar()
            tok.Type = TOKEN_LE
            tok.Literal = "<="
        } else {
            tok.Type = TOKEN_LT
            tok.Literal = "<"
        }
    case '>':
        if l.peekChar() == '=' {
            l.readChar()
            tok.Type = TOKEN_GE
            tok.Literal = ">="
        } else {
            tok.Type = TOKEN_GT
            tok.Literal = ">"
        }
    case '&':
        if l.peekChar() == '&' {
            l.readChar()
            tok.Type = TOKEN_AND
            tok.Literal = "&&"
        } else {
            tok.Type = TOKEN_ILLEGAL
            tok.Literal = string(l.ch)
        }
    case '|':
        if l.peekChar() == '|' {
            l.readChar()
            tok.Type = TOKEN_OR
            tok.Literal = "||"
        } else {
            tok.Type = TOKEN_ILLEGAL
            tok.Literal = string(l.ch)
        }
    case '"':
        tok.Type = TOKEN_STRING
        tok.Literal = l.readString()
    case 0:
        tok.Type = TOKEN_EOF
        tok.Literal = ""
    default:
        if isLetter(l.ch) {
            tok.Literal = l.readIdentifier()
            tok.Type = LookupIdent(tok.Literal)
            return tok
        } else if isDigit(l.ch) {
            lit, tokType := l.readNumber()
            tok.Literal = lit
            tok.Type = tokType
            return tok
        } else {
            tok.Type = TOKEN_ILLEGAL
            tok.Literal = string(l.ch)
        }
    }

    l.readChar()
    return tok
}

func (l *Lexer) readChar() {
    if l.readPos >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPos]
    }
    l.pos = l.readPos
    l.readPos++

    if l.ch == '\n' {
        l.line++
        l.column = 0
    } else {
        l.column++
    }
}

func (l *Lexer) peekChar() byte {
    if l.readPos >= len(l.input) {
        return 0
    }
    return l.input[l.readPos]
}

func (l *Lexer) readString() string {
    position := l.pos + 1
    for {
        l.readChar()
        if l.ch == '"' || l.ch == 0 {
            break
        }
        if l.ch == '\\' {
            l.readChar() // エスケープ文字をスキップ
        }
    }
    return l.input[position:l.pos]
}

func (l *Lexer) readNumber() (string, TokenType) {
    position := l.pos
    tokType := TOKEN_INT

    for isDigit(l.ch) {
        l.readChar()
    }

    if l.ch == '.' && isDigit(l.peekChar()) {
        tokType = TOKEN_FLOAT
        l.readChar()
        for isDigit(l.ch) {
            l.readChar()
        }
    }

    // Duration: 24h, 30s, 500ms, etc.
    if isLetter(l.ch) {
        tokType = TOKEN_DURATION
        for isLetter(l.ch) {
            l.readChar()
        }
    }

    return l.input[position:l.pos], tokType
}

func (l *Lexer) readIdentifier() string {
    position := l.pos
    for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
        l.readChar()
    }
    return l.input[position:l.pos]
}

func (l *Lexer) skipWhitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        l.readChar()
    }
}

func (l *Lexer) skipComment() {
    if l.ch == '/' {
        if l.peekChar() == '/' {
            // 行コメント
            for l.ch != '\n' && l.ch != 0 {
                l.readChar()
            }
        } else if l.peekChar() == '*' {
            // ブロックコメント
            l.readChar()
            l.readChar()
            for {
                if l.ch == 0 {
                    break
                }
                if l.ch == '*' && l.peekChar() == '/' {
                    l.readChar()
                    l.readChar()
                    break
                }
                l.readChar()
            }
        }
    }
}

func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}
```

**責務**: ソースコードをトークン列に変換

---

### internal/infrastructure/parser/parser.go

```go
package parser

import (
    "fmt"
    "strconv"

    "pact/internal/domain/ast"
    "pact/internal/domain/errors"
)

// Parser は構文解析器
type Parser struct {
    lexer     *Lexer
    file      string
    curToken  Token
    peekToken Token
}

// NewParser は新しい Parser を作成する
func NewParser(lexer *Lexer, file string) *Parser {
    p := &Parser{
        lexer: lexer,
        file:  file,
    }
    p.nextToken()
    p.nextToken()
    return p
}

// ParseFile はファイル全体をパースする
func (p *Parser) ParseFile() (*ast.SpecFile, error) {
    spec := &ast.SpecFile{
        Path: p.file,
    }

    // import 文をパース
    for p.curTokenIs(TOKEN_IMPORT) {
        imp, err := p.parseImport()
        if err != nil {
            return nil, err
        }
        spec.Imports = append(spec.Imports, *imp)
    }

    // アノテーションをパース
    annotations, err := p.parseAnnotations()
    if err != nil {
        return nil, err
    }

    // component をパース
    if !p.curTokenIs(TOKEN_COMPONENT) {
        return nil, p.error("expected 'component'")
    }

    comp, err := p.parseComponent(annotations)
    if err != nil {
        return nil, err
    }
    spec.Component = *comp

    return spec, nil
}

func (p *Parser) parseImport() (*ast.ImportDecl, error) {
    imp := &ast.ImportDecl{
        Pos: p.curPos(),
    }

    p.nextToken() // 'import' を消費

    if !p.curTokenIs(TOKEN_STRING) {
        return nil, p.error("expected string after 'import'")
    }
    imp.Path = p.curToken.Literal

    p.nextToken()

    if p.curTokenIs(TOKEN_AS) {
        p.nextToken()
        if !p.curTokenIs(TOKEN_IDENT) {
            return nil, p.error("expected identifier after 'as'")
        }
        alias := p.curToken.Literal
        imp.Alias = &alias
        p.nextToken()
    }

    return imp, nil
}

func (p *Parser) parseComponent(annotations []ast.AnnotationDecl) (*ast.ComponentDecl, error) {
    comp := &ast.ComponentDecl{
        Pos:         p.curPos(),
        Annotations: annotations,
    }

    p.nextToken() // 'component' を消費

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected component name")
    }
    comp.Name = p.curToken.Literal

    p.nextToken()

    if !p.curTokenIs(TOKEN_LBRACE) {
        return nil, p.error("expected '{'")
    }
    p.nextToken()

    // Body をパース
    for !p.curTokenIs(TOKEN_RBRACE) && !p.curTokenIs(TOKEN_EOF) {
        annotations, err := p.parseAnnotations()
        if err != nil {
            return nil, err
        }

        switch p.curToken.Type {
        case TOKEN_TYPE:
            typ, err := p.parseType(annotations)
            if err != nil {
                return nil, err
            }
            comp.Body.Types = append(comp.Body.Types, *typ)

        case TOKEN_ENUM:
            enum, err := p.parseEnum(annotations)
            if err != nil {
                return nil, err
            }
            comp.Body.Types = append(comp.Body.Types, *enum)

        case TOKEN_DEPENDS, TOKEN_EXTENDS, TOKEN_IMPLEMENTS, TOKEN_CONTAINS, TOKEN_AGGREGATES:
            rel, err := p.parseRelation(annotations)
            if err != nil {
                return nil, err
            }
            comp.Body.Relations = append(comp.Body.Relations, *rel)

        case TOKEN_PROVIDES:
            iface, err := p.parseInterface(annotations)
            if err != nil {
                return nil, err
            }
            comp.Body.Provides = append(comp.Body.Provides, *iface)

        case TOKEN_REQUIRES:
            iface, err := p.parseInterface(annotations)
            if err != nil {
                return nil, err
            }
            comp.Body.Requires = append(comp.Body.Requires, *iface)

        case TOKEN_FLOW:
            flow, err := p.parseFlow(annotations)
            if err != nil {
                return nil, err
            }
            comp.Body.Flows = append(comp.Body.Flows, *flow)

        case TOKEN_STATES:
            states, err := p.parseStates(annotations)
            if err != nil {
                return nil, err
            }
            comp.Body.States = append(comp.Body.States, *states)

        default:
            return nil, p.error(fmt.Sprintf("unexpected token: %v", p.curToken.Literal))
        }
    }

    if !p.curTokenIs(TOKEN_RBRACE) {
        return nil, p.error("expected '}'")
    }
    p.nextToken()

    return comp, nil
}

func (p *Parser) parseAnnotations() ([]ast.AnnotationDecl, error) {
    var annotations []ast.AnnotationDecl

    for p.curTokenIs(TOKEN_AT) {
        ann, err := p.parseAnnotation()
        if err != nil {
            return nil, err
        }
        annotations = append(annotations, *ann)
    }

    return annotations, nil
}

func (p *Parser) parseAnnotation() (*ast.AnnotationDecl, error) {
    ann := &ast.AnnotationDecl{
        Pos: p.curPos(),
    }

    p.nextToken() // '@' を消費

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected annotation name")
    }
    ann.Name = p.curToken.Literal

    p.nextToken()

    // 引数があれば
    if p.curTokenIs(TOKEN_LPAREN) {
        p.nextToken()

        for !p.curTokenIs(TOKEN_RPAREN) && !p.curTokenIs(TOKEN_EOF) {
            arg := ast.AnnotationArg{}

            // 名前付き引数かどうか
            if p.curTokenIs(TOKEN_IDENT) && p.peekTokenIs(TOKEN_COLON) {
                key := p.curToken.Literal
                arg.Key = &key
                p.nextToken() // identifier
                p.nextToken() // ':'
            }

            if !p.curTokenIs(TOKEN_STRING) {
                return nil, p.error("expected string value in annotation")
            }
            arg.Value = p.curToken.Literal
            ann.Args = append(ann.Args, arg)

            p.nextToken()

            if p.curTokenIs(TOKEN_COMMA) {
                p.nextToken()
            }
        }

        if !p.curTokenIs(TOKEN_RPAREN) {
            return nil, p.error("expected ')'")
        }
        p.nextToken()
    }

    return ann, nil
}

// 以下、他のパース関数（省略可能な部分は実装を簡略化）

func (p *Parser) parseType(annotations []ast.AnnotationDecl) (*ast.TypeDecl, error) {
    typ := &ast.TypeDecl{
        Pos:         p.curPos(),
        Kind:        ast.TypeKindStruct,
        Annotations: annotations,
    }

    p.nextToken() // 'type' を消費

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected type name")
    }
    typ.Name = p.curToken.Literal
    p.nextToken()

    if !p.curTokenIs(TOKEN_LBRACE) {
        return nil, p.error("expected '{'")
    }
    p.nextToken()

    for !p.curTokenIs(TOKEN_RBRACE) && !p.curTokenIs(TOKEN_EOF) {
        field, err := p.parseField()
        if err != nil {
            return nil, err
        }
        typ.Fields = append(typ.Fields, *field)
    }

    p.nextToken() // '}' を消費

    return typ, nil
}

func (p *Parser) parseEnum(annotations []ast.AnnotationDecl) (*ast.TypeDecl, error) {
    enum := &ast.TypeDecl{
        Pos:         p.curPos(),
        Kind:        ast.TypeKindEnum,
        Annotations: annotations,
    }

    p.nextToken() // 'enum' を消費

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected enum name")
    }
    enum.Name = p.curToken.Literal
    p.nextToken()

    if !p.curTokenIs(TOKEN_LBRACE) {
        return nil, p.error("expected '{'")
    }
    p.nextToken()

    for p.curTokenIs(TOKEN_IDENT) {
        enum.Values = append(enum.Values, p.curToken.Literal)
        p.nextToken()
    }

    if !p.curTokenIs(TOKEN_RBRACE) {
        return nil, p.error("expected '}'")
    }
    p.nextToken()

    return enum, nil
}

func (p *Parser) parseField() (*ast.FieldDecl, error) {
    field := &ast.FieldDecl{
        Pos:        p.curPos(),
        Visibility: ast.VisibilityPublic,
    }

    // アノテーション
    annotations, err := p.parseAnnotations()
    if err != nil {
        return nil, err
    }
    field.Annotations = annotations

    // Visibility
    switch p.curToken.Type {
    case TOKEN_PLUS:
        field.Visibility = ast.VisibilityPublic
        p.nextToken()
    case TOKEN_MINUS:
        field.Visibility = ast.VisibilityPrivate
        p.nextToken()
    case TOKEN_HASH:
        field.Visibility = ast.VisibilityProtected
        p.nextToken()
    case TOKEN_TILDE:
        field.Visibility = ast.VisibilityPackage
        p.nextToken()
    }

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected field name")
    }
    field.Name = p.curToken.Literal
    p.nextToken()

    if !p.curTokenIs(TOKEN_COLON) {
        return nil, p.error("expected ':'")
    }
    p.nextToken()

    typeExpr, err := p.parseTypeExpr()
    if err != nil {
        return nil, err
    }
    field.Type = *typeExpr

    return field, nil
}

func (p *Parser) parseTypeExpr() (*ast.TypeExpr, error) {
    typeExpr := &ast.TypeExpr{
        Pos: p.curPos(),
    }

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected type name")
    }
    typeExpr.Name = p.curToken.Literal
    p.nextToken()

    if p.curTokenIs(TOKEN_QUESTION) {
        typeExpr.Nullable = true
        p.nextToken()
    }

    if p.curTokenIs(TOKEN_LBRACKET) {
        p.nextToken()
        if !p.curTokenIs(TOKEN_RBRACKET) {
            return nil, p.error("expected ']'")
        }
        typeExpr.Array = true
        p.nextToken()
    }

    return typeExpr, nil
}

func (p *Parser) parseRelation(annotations []ast.AnnotationDecl) (*ast.RelationDecl, error) {
    rel := &ast.RelationDecl{
        Pos:         p.curPos(),
        Annotations: annotations,
        TargetType:  ast.TargetComponent,
    }

    // Kind を判定
    switch p.curToken.Type {
    case TOKEN_DEPENDS:
        rel.Kind = ast.RelationDepends
        p.nextToken()
        if !p.curTokenIs(TOKEN_ON) {
            return nil, p.error("expected 'on' after 'depends'")
        }
        p.nextToken()
    case TOKEN_EXTENDS:
        rel.Kind = ast.RelationExtends
        p.nextToken()
    case TOKEN_IMPLEMENTS:
        rel.Kind = ast.RelationImplements
        p.nextToken()
    case TOKEN_CONTAINS:
        rel.Kind = ast.RelationContains
        p.nextToken()
    case TOKEN_AGGREGATES:
        rel.Kind = ast.RelationAggregates
        p.nextToken()
    }

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected target name")
    }
    rel.Target = p.curToken.Literal
    p.nextToken()

    // TargetType
    if p.curTokenIs(TOKEN_COLON) {
        p.nextToken()
        if !p.curTokenIs(TOKEN_IDENT) {
            return nil, p.error("expected target type")
        }
        switch p.curToken.Literal {
        case "database":
            rel.TargetType = ast.TargetDatabase
        case "external":
            rel.TargetType = ast.TargetExternal
        case "queue":
            rel.TargetType = ast.TargetQueue
        case "actor":
            rel.TargetType = ast.TargetActor
        case "component":
            rel.TargetType = ast.TargetComponent
        default:
            return nil, p.error(fmt.Sprintf("unknown target type: %s", p.curToken.Literal))
        }
        p.nextToken()
    }

    // Alias
    if p.curTokenIs(TOKEN_AS) {
        p.nextToken()
        if !p.curTokenIs(TOKEN_IDENT) {
            return nil, p.error("expected alias name")
        }
        alias := p.curToken.Literal
        rel.Alias = &alias
        p.nextToken()
    }

    return rel, nil
}

func (p *Parser) parseInterface(annotations []ast.AnnotationDecl) (*ast.InterfaceDecl, error) {
    iface := &ast.InterfaceDecl{
        Pos:         p.curPos(),
        Annotations: annotations,
    }

    p.nextToken() // 'provides' or 'requires' を消費

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected interface name")
    }
    iface.Name = p.curToken.Literal
    p.nextToken()

    if !p.curTokenIs(TOKEN_LBRACE) {
        return nil, p.error("expected '{'")
    }
    p.nextToken()

    for !p.curTokenIs(TOKEN_RBRACE) && !p.curTokenIs(TOKEN_EOF) {
        method, err := p.parseMethod()
        if err != nil {
            return nil, err
        }
        iface.Methods = append(iface.Methods, *method)
    }

    p.nextToken() // '}' を消費

    return iface, nil
}

func (p *Parser) parseMethod() (*ast.MethodDecl, error) {
    method := &ast.MethodDecl{
        Pos: p.curPos(),
    }

    annotations, err := p.parseAnnotations()
    if err != nil {
        return nil, err
    }
    method.Annotations = annotations

    if p.curTokenIs(TOKEN_ASYNC) {
        method.Async = true
        p.nextToken()
    }

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected method name")
    }
    method.Name = p.curToken.Literal
    p.nextToken()

    // パラメータ
    if !p.curTokenIs(TOKEN_LPAREN) {
        return nil, p.error("expected '('")
    }
    p.nextToken()

    for !p.curTokenIs(TOKEN_RPAREN) && !p.curTokenIs(TOKEN_EOF) {
        param, err := p.parseParam()
        if err != nil {
            return nil, err
        }
        method.Params = append(method.Params, *param)

        if p.curTokenIs(TOKEN_COMMA) {
            p.nextToken()
        }
    }

    p.nextToken() // ')' を消費

    // 戻り値
    if !p.curTokenIs(TOKEN_ARROW) {
        return nil, p.error("expected '->'")
    }
    p.nextToken()

    typeExpr, err := p.parseTypeExpr()
    if err != nil {
        return nil, err
    }
    method.Returns = *typeExpr

    // throws
    if p.curTokenIs(TOKEN_THROWS) {
        p.nextToken()
        for p.curTokenIs(TOKEN_IDENT) {
            method.Throws = append(method.Throws, p.curToken.Literal)
            p.nextToken()
            if p.curTokenIs(TOKEN_COMMA) {
                p.nextToken()
            }
        }
    }

    return method, nil
}

func (p *Parser) parseParam() (*ast.ParamDecl, error) {
    param := &ast.ParamDecl{
        Pos: p.curPos(),
    }

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected parameter name")
    }
    param.Name = p.curToken.Literal
    p.nextToken()

    if !p.curTokenIs(TOKEN_COLON) {
        return nil, p.error("expected ':'")
    }
    p.nextToken()

    typeExpr, err := p.parseTypeExpr()
    if err != nil {
        return nil, err
    }
    param.Type = *typeExpr

    return param, nil
}

func (p *Parser) parseFlow(annotations []ast.AnnotationDecl) (*ast.FlowDecl, error) {
    flow := &ast.FlowDecl{
        Pos:         p.curPos(),
        Annotations: annotations,
    }

    p.nextToken() // 'flow' を消費

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected flow name")
    }
    flow.Name = p.curToken.Literal
    p.nextToken()

    if !p.curTokenIs(TOKEN_LBRACE) {
        return nil, p.error("expected '{'")
    }
    p.nextToken()

    steps, err := p.parseSteps()
    if err != nil {
        return nil, err
    }
    flow.Steps = steps

    if !p.curTokenIs(TOKEN_RBRACE) {
        return nil, p.error("expected '}'")
    }
    p.nextToken()

    return flow, nil
}

func (p *Parser) parseSteps() ([]ast.StepDecl, error) {
    var steps []ast.StepDecl

    for !p.curTokenIs(TOKEN_RBRACE) && !p.curTokenIs(TOKEN_EOF) {
        step, err := p.parseStep()
        if err != nil {
            return nil, err
        }
        if step != nil {
            steps = append(steps, step)
        }
    }

    return steps, nil
}

func (p *Parser) parseStep() (ast.StepDecl, error) {
    annotations, err := p.parseAnnotations()
    if err != nil {
        return nil, err
    }

    switch p.curToken.Type {
    case TOKEN_IF:
        return p.parseIfStep(annotations)
    case TOKEN_FOR:
        return p.parseForStep(annotations)
    case TOKEN_WHILE:
        return p.parseWhileStep(annotations)
    case TOKEN_RETURN:
        return p.parseReturnStep(annotations)
    case TOKEN_THROW:
        return p.parseThrowStep(annotations)
    case TOKEN_AWAIT:
        return p.parseCallStep(annotations, true)
    case TOKEN_IDENT:
        // 代入 or 呼び出し
        if p.peekTokenIs(TOKEN_ASSIGN) {
            return p.parseAssignStep(annotations)
        } else if p.peekTokenIs(TOKEN_DOT) {
            return p.parseCallStep(annotations, false)
        }
    }

    return nil, p.error(fmt.Sprintf("unexpected token in flow: %v", p.curToken.Literal))
}

func (p *Parser) parseIfStep(annotations []ast.AnnotationDecl) (*ast.IfStep, error) {
    step := &ast.IfStep{
        baseStep: baseStep{
            Pos:         p.curPos(),
            Annotations: annotations,
        },
    }

    p.nextToken() // 'if' を消費

    cond, err := p.parseExpression()
    if err != nil {
        return nil, err
    }
    step.Condition = cond

    if !p.curTokenIs(TOKEN_LBRACE) {
        return nil, p.error("expected '{'")
    }
    p.nextToken()

    then, err := p.parseSteps()
    if err != nil {
        return nil, err
    }
    step.Then = then

    if !p.curTokenIs(TOKEN_RBRACE) {
        return nil, p.error("expected '}'")
    }
    p.nextToken()

    if p.curTokenIs(TOKEN_ELSE) {
        p.nextToken()
        if !p.curTokenIs(TOKEN_LBRACE) {
            return nil, p.error("expected '{'")
        }
        p.nextToken()

        els, err := p.parseSteps()
        if err != nil {
            return nil, err
        }
        step.Else = els

        if !p.curTokenIs(TOKEN_RBRACE) {
            return nil, p.error("expected '}'")
        }
        p.nextToken()
    }

    return step, nil
}

func (p *Parser) parseForStep(annotations []ast.AnnotationDecl) (*ast.ForStep, error) {
    step := &ast.ForStep{
        baseStep: baseStep{
            Pos:         p.curPos(),
            Annotations: annotations,
        },
    }

    p.nextToken() // 'for' を消費

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected variable name")
    }
    step.Variable = p.curToken.Literal
    p.nextToken()

    if !p.curTokenIs(TOKEN_IN) {
        return nil, p.error("expected 'in'")
    }
    p.nextToken()

    iter, err := p.parseExpression()
    if err != nil {
        return nil, err
    }
    step.Iterable = iter

    if !p.curTokenIs(TOKEN_LBRACE) {
        return nil, p.error("expected '{'")
    }
    p.nextToken()

    body, err := p.parseSteps()
    if err != nil {
        return nil, err
    }
    step.Body = body

    if !p.curTokenIs(TOKEN_RBRACE) {
        return nil, p.error("expected '}'")
    }
    p.nextToken()

    return step, nil
}

func (p *Parser) parseWhileStep(annotations []ast.AnnotationDecl) (*ast.WhileStep, error) {
    step := &ast.WhileStep{
        baseStep: baseStep{
            Pos:         p.curPos(),
            Annotations: annotations,
        },
    }

    p.nextToken() // 'while' を消費

    cond, err := p.parseExpression()
    if err != nil {
        return nil, err
    }
    step.Condition = cond

    if !p.curTokenIs(TOKEN_LBRACE) {
        return nil, p.error("expected '{'")
    }
    p.nextToken()

    body, err := p.parseSteps()
    if err != nil {
        return nil, err
    }
    step.Body = body

    if !p.curTokenIs(TOKEN_RBRACE) {
        return nil, p.error("expected '}'")
    }
    p.nextToken()

    return step, nil
}

func (p *Parser) parseReturnStep(annotations []ast.AnnotationDecl) (*ast.ReturnStep, error) {
    step := &ast.ReturnStep{
        baseStep: baseStep{
            Pos:         p.curPos(),
            Annotations: annotations,
        },
    }

    p.nextToken() // 'return' を消費

    if !p.curTokenIs(TOKEN_RBRACE) && !p.curTokenIs(TOKEN_EOF) {
        val, err := p.parseExpression()
        if err != nil {
            return nil, err
        }
        step.Value = val
    }

    return step, nil
}

func (p *Parser) parseThrowStep(annotations []ast.AnnotationDecl) (*ast.ThrowStep, error) {
    step := &ast.ThrowStep{
        baseStep: baseStep{
            Pos:         p.curPos(),
            Annotations: annotations,
        },
    }

    p.nextToken() // 'throw' を消費

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected error name")
    }
    step.Error = p.curToken.Literal
    p.nextToken()

    return step, nil
}

func (p *Parser) parseAssignStep(annotations []ast.AnnotationDecl) (*ast.AssignStep, error) {
    step := &ast.AssignStep{
        baseStep: baseStep{
            Pos:         p.curPos(),
            Annotations: annotations,
        },
    }

    step.Variable = p.curToken.Literal
    p.nextToken() // identifier
    p.nextToken() // '='

    val, err := p.parseExpression()
    if err != nil {
        return nil, err
    }
    step.Value = val

    // ?? throw
    if p.curTokenIs(TOKEN_NULLISH) {
        p.nextToken()
        if !p.curTokenIs(TOKEN_THROW) {
            return nil, p.error("expected 'throw' after '??'")
        }
        p.nextToken()
        if !p.curTokenIs(TOKEN_IDENT) {
            return nil, p.error("expected error name")
        }
        throwErr := p.curToken.Literal
        step.ThrowOnNull = &throwErr
        p.nextToken()
    }

    return step, nil
}

func (p *Parser) parseCallStep(annotations []ast.AnnotationDecl, async bool) (*ast.CallStep, error) {
    step := &ast.CallStep{
        baseStep: baseStep{
            Pos:         p.curPos(),
            Annotations: annotations,
        },
        Async: async,
    }

    if async {
        p.nextToken() // 'await' を消費
    }

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected target name")
    }
    step.Target = p.curToken.Literal
    p.nextToken()

    if !p.curTokenIs(TOKEN_DOT) {
        return nil, p.error("expected '.'")
    }
    p.nextToken()

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected method name")
    }
    step.Method = p.curToken.Literal
    p.nextToken()

    if !p.curTokenIs(TOKEN_LPAREN) {
        return nil, p.error("expected '('")
    }
    p.nextToken()

    for !p.curTokenIs(TOKEN_RPAREN) && !p.curTokenIs(TOKEN_EOF) {
        arg, err := p.parseExpression()
        if err != nil {
            return nil, err
        }
        step.Args = append(step.Args, arg)

        if p.curTokenIs(TOKEN_COMMA) {
            p.nextToken()
        }
    }

    p.nextToken() // ')' を消費

    return step, nil
}

func (p *Parser) parseStates(annotations []ast.AnnotationDecl) (*ast.StatesDecl, error) {
    states := &ast.StatesDecl{
        Pos:         p.curPos(),
        Annotations: annotations,
    }

    p.nextToken() // 'states' を消費

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected states name")
    }
    states.Name = p.curToken.Literal
    p.nextToken()

    if !p.curTokenIs(TOKEN_LBRACE) {
        return nil, p.error("expected '{'")
    }
    p.nextToken()

    for !p.curTokenIs(TOKEN_RBRACE) && !p.curTokenIs(TOKEN_EOF) {
        switch p.curToken.Type {
        case TOKEN_INITIAL:
            p.nextToken()
            if !p.curTokenIs(TOKEN_IDENT) {
                return nil, p.error("expected state name")
            }
            states.Initial = p.curToken.Literal
            p.nextToken()

        case TOKEN_FINAL:
            p.nextToken()
            if !p.curTokenIs(TOKEN_IDENT) {
                return nil, p.error("expected state name")
            }
            states.Finals = append(states.Finals, p.curToken.Literal)
            p.nextToken()

        case TOKEN_STATE, TOKEN_PARALLEL:
            state, err := p.parseState()
            if err != nil {
                return nil, err
            }
            states.States = append(states.States, *state)

        case TOKEN_IDENT:
            // 遷移
            tr, err := p.parseTransition()
            if err != nil {
                return nil, err
            }
            states.Transitions = append(states.Transitions, *tr)

        default:
            return nil, p.error(fmt.Sprintf("unexpected token in states: %v", p.curToken.Literal))
        }
    }

    p.nextToken() // '}' を消費

    return states, nil
}

func (p *Parser) parseState() (*ast.StateDecl, error) {
    state := &ast.StateDecl{
        Pos:  p.curPos(),
        Kind: ast.StateAtomic,
    }

    if p.curTokenIs(TOKEN_PARALLEL) {
        state.Kind = ast.StateParallel
    } else {
        state.Kind = ast.StateCompound
    }
    p.nextToken()

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected state name")
    }
    state.Name = p.curToken.Literal
    p.nextToken()

    if !p.curTokenIs(TOKEN_LBRACE) {
        return nil, p.error("expected '{'")
    }
    p.nextToken()

    for !p.curTokenIs(TOKEN_RBRACE) && !p.curTokenIs(TOKEN_EOF) {
        switch p.curToken.Type {
        case TOKEN_ENTRY:
            p.nextToken()
            actions, err := p.parseActionList()
            if err != nil {
                return nil, err
            }
            state.Entry = actions

        case TOKEN_EXIT:
            p.nextToken()
            actions, err := p.parseActionList()
            if err != nil {
                return nil, err
            }
            state.Exit = actions

        case TOKEN_INITIAL:
            p.nextToken()
            if !p.curTokenIs(TOKEN_IDENT) {
                return nil, p.error("expected state name")
            }
            initial := p.curToken.Literal
            state.Initial = &initial
            p.nextToken()

        case TOKEN_STATE, TOKEN_PARALLEL:
            child, err := p.parseState()
            if err != nil {
                return nil, err
            }
            state.States = append(state.States, *child)

        case TOKEN_REGION:
            region, err := p.parseRegion()
            if err != nil {
                return nil, err
            }
            state.Regions = append(state.Regions, *region)

        case TOKEN_IDENT:
            tr, err := p.parseTransition()
            if err != nil {
                return nil, err
            }
            state.Transitions = append(state.Transitions, *tr)

        default:
            return nil, p.error(fmt.Sprintf("unexpected token in state: %v", p.curToken.Literal))
        }
    }

    p.nextToken() // '}' を消費

    return state, nil
}

func (p *Parser) parseRegion() (*ast.RegionDecl, error) {
    region := &ast.RegionDecl{
        Pos: p.curPos(),
    }

    p.nextToken() // 'region' を消費

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected region name")
    }
    region.Name = p.curToken.Literal
    p.nextToken()

    if !p.curTokenIs(TOKEN_LBRACE) {
        return nil, p.error("expected '{'")
    }
    p.nextToken()

    for !p.curTokenIs(TOKEN_RBRACE) && !p.curTokenIs(TOKEN_EOF) {
        switch p.curToken.Type {
        case TOKEN_INITIAL:
            p.nextToken()
            if !p.curTokenIs(TOKEN_IDENT) {
                return nil, p.error("expected state name")
            }
            region.Initial = p.curToken.Literal
            p.nextToken()

        case TOKEN_FINAL:
            p.nextToken()
            if !p.curTokenIs(TOKEN_IDENT) {
                return nil, p.error("expected state name")
            }
            region.Finals = append(region.Finals, p.curToken.Literal)
            p.nextToken()

        case TOKEN_STATE:
            state, err := p.parseState()
            if err != nil {
                return nil, err
            }
            region.States = append(region.States, *state)

        case TOKEN_IDENT:
            tr, err := p.parseTransition()
            if err != nil {
                return nil, err
            }
            region.Transitions = append(region.Transitions, *tr)
        }
    }

    p.nextToken() // '}' を消費

    return region, nil
}

func (p *Parser) parseTransition() (*ast.TransitionDecl, error) {
    tr := &ast.TransitionDecl{
        Pos: p.curPos(),
    }

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected source state")
    }
    tr.From = p.curToken.Literal
    p.nextToken()

    if !p.curTokenIs(TOKEN_ARROW) {
        return nil, p.error("expected '->'")
    }
    p.nextToken()

    if !p.curTokenIs(TOKEN_IDENT) {
        return nil, p.error("expected target state")
    }
    tr.To = p.curToken.Literal
    p.nextToken()

    // Trigger
    switch p.curToken.Type {
    case TOKEN_ON:
        p.nextToken()
        if !p.curTokenIs(TOKEN_IDENT) {
            return nil, p.error("expected event name")
        }
        tr.Trigger = &ast.EventTrigger{
            Pos:  p.curPos(),
            Name: p.curToken.Literal,
        }
        p.nextToken()

    case TOKEN_AFTER:
        p.nextToken()
        if !p.curTokenIs(TOKEN_DURATION) && !p.curTokenIs(TOKEN_INT) {
            return nil, p.error("expected duration")
        }
        val, unit := p.parseDuration()
        tr.Trigger = &ast.AfterTrigger{
            Pos:   p.curPos(),
            Value: val,
            Unit:  unit,
        }

    case TOKEN_WHEN:
        p.nextToken()
        cond, err := p.parseExpression()
        if err != nil {
            return nil, err
        }
        tr.Trigger = &ast.WhenTrigger{
            Pos:       p.curPos(),
            Condition: cond,
        }
    }

    // Guard
    if p.curTokenIs(TOKEN_WHEN) {
        p.nextToken()
        guard, err := p.parseExpression()
        if err != nil {
            return nil, err
        }
        tr.Guard = guard
    }

    // Actions
    if p.curTokenIs(TOKEN_DO) {
        p.nextToken()
        actions, err := p.parseActionList()
        if err != nil {
            return nil, err
        }
        tr.Actions = actions
    }

    return tr, nil
}

func (p *Parser) parseActionList() ([]string, error) {
    if !p.curTokenIs(TOKEN_LBRACKET) {
        return nil, p.error("expected '['")
    }
    p.nextToken()

    var actions []string
    for !p.curTokenIs(TOKEN_RBRACKET) && !p.curTokenIs(TOKEN_EOF) {
        if !p.curTokenIs(TOKEN_IDENT) {
            return nil, p.error("expected action name")
        }
        actions = append(actions, p.curToken.Literal)
        p.nextToken()

        if p.curTokenIs(TOKEN_COMMA) {
            p.nextToken()
        }
    }

    p.nextToken() // ']' を消費

    return actions, nil
}

func (p *Parser) parseDuration() (int, string) {
    literal := p.curToken.Literal
    p.nextToken()

    // "3s", "500ms" などをパース
    var val int
    var unit string
    for i, ch := range literal {
        if ch >= '0' && ch <= '9' {
            continue
        }
        val, _ = strconv.Atoi(literal[:i])
        unit = literal[i:]
        break
    }

    return val, unit
}

func (p *Parser) parseExpression() (ast.ExprDecl, error) {
    // 簡略化: Pratt パーサーの完全実装は省略
    // ここでは単純な式のみ対応
    return p.parsePrimaryExpr()
}

func (p *Parser) parsePrimaryExpr() (ast.ExprDecl, error) {
    switch p.curToken.Type {
    case TOKEN_IDENT:
        name := p.curToken.Literal
        p.nextToken()

        // フィールドアクセスまたはメソッド呼び出し
        if p.curTokenIs(TOKEN_DOT) {
            p.nextToken()
            if !p.curTokenIs(TOKEN_IDENT) {
                return nil, p.error("expected identifier after '.'")
            }
            field := p.curToken.Literal
            p.nextToken()

            if p.curTokenIs(TOKEN_LPAREN) {
                // メソッド呼び出し
                p.nextToken()
                var args []ast.ExprDecl
                for !p.curTokenIs(TOKEN_RPAREN) && !p.curTokenIs(TOKEN_EOF) {
                    arg, err := p.parseExpression()
                    if err != nil {
                        return nil, err
                    }
                    args = append(args, arg)
                    if p.curTokenIs(TOKEN_COMMA) {
                        p.nextToken()
                    }
                }
                p.nextToken() // ')' を消費
                return &ast.CallExpr{
                    Pos:    p.curPos(),
                    Target: name,
                    Method: field,
                    Args:   args,
                }, nil
            }

            return &ast.FieldExpr{
                Pos:    p.curPos(),
                Object: &ast.VariableExpr{Name: name},
                Field:  field,
            }, nil
        }

        return &ast.VariableExpr{
            Pos:  p.curPos(),
            Name: name,
        }, nil

    case TOKEN_STRING:
        val := p.curToken.Literal
        p.nextToken()
        return &ast.LiteralExpr{
            Pos:   p.curPos(),
            Value: val,
        }, nil

    case TOKEN_INT:
        val, _ := strconv.Atoi(p.curToken.Literal)
        p.nextToken()
        return &ast.LiteralExpr{
            Pos:   p.curPos(),
            Value: val,
        }, nil

    case TOKEN_TRUE:
        p.nextToken()
        return &ast.LiteralExpr{
            Pos:   p.curPos(),
            Value: true,
        }, nil

    case TOKEN_FALSE:
        p.nextToken()
        return &ast.LiteralExpr{
            Pos:   p.curPos(),
            Value: false,
        }, nil

    case TOKEN_NULL:
        p.nextToken()
        return &ast.LiteralExpr{
            Pos:   p.curPos(),
            Value: nil,
        }, nil

    case TOKEN_NOT:
        p.nextToken()
        operand, err := p.parsePrimaryExpr()
        if err != nil {
            return nil, err
        }
        return &ast.UnaryExpr{
            Pos:     p.curPos(),
            Op:      "!",
            Operand: operand,
        }, nil
    }

    return nil, p.error(fmt.Sprintf("unexpected token in expression: %v", p.curToken.Literal))
}

// ヘルパー

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.lexer.NextToken()
}

func (p *Parser) curTokenIs(t TokenType) bool {
    return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
    return p.peekToken.Type == t
}

func (p *Parser) curPos() ast.Position {
    return ast.Position{
        File:   p.file,
        Line:   p.curToken.Line,
        Column: p.curToken.Column,
    }
}

func (p *Parser) error(msg string) error {
    return &errors.ParseError{
        Pos:     p.curPos(),
        Message: msg,
    }
}
```

**責務**: トークン列を AST に変換

---

### internal/infrastructure/parser/adapter.go

```go
package parser

import (
    "pact/internal/domain/ast"
    "pact/internal/infrastructure/filesystem"
)

// FileParser は application/parser.Parser インターフェースの実装
type FileParser struct {
    reader *filesystem.Reader
}

// NewFileParser は新しい FileParser を作成する
func NewFileParser(reader *filesystem.Reader) *FileParser {
    return &FileParser{reader: reader}
}

// ParseFile はファイルをパースする
func (fp *FileParser) ParseFile(path string) (*ast.SpecFile, error) {
    content, err := fp.reader.ReadFile(path)
    if err != nil {
        return nil, err
    }
    return fp.ParseString(content, path)
}

// ParseString は文字列をパースする
func (fp *FileParser) ParseString(content, filename string) (*ast.SpecFile, error) {
    lexer := NewLexer(content, filename)
    parser := NewParser(lexer, filename)
    return parser.ParseFile()
}
```

**責務**: application/parser.Parser インターフェースの実装

---

### internal/infrastructure/resolver/import.go

```go
package resolver

import (
    "path/filepath"

    "pact/internal/domain/ast"
    "pact/internal/domain/errors"
    "pact/internal/application/parser"
)

// ImportResolver はインポートを解決する
type ImportResolver struct {
    basePath string
    parser   parser.Parser
    parsed   map[string]*ast.SpecFile
    visiting map[string]bool
}

// NewImportResolver は新しい ImportResolver を作成する
func NewImportResolver(basePath string, parser parser.Parser) *ImportResolver {
    return &ImportResolver{
        basePath: basePath,
        parser:   parser,
        parsed:   make(map[string]*ast.SpecFile),
        visiting: make(map[string]bool),
    }
}

// Resolve はインポートを解決し、依存順にソートされたファイルリストを返す
func (r *ImportResolver) Resolve(files []*ast.SpecFile) ([]*ast.SpecFile, error) {
    // 既存ファイルを登録
    for _, f := range files {
        r.parsed[f.Path] = f
    }

    // 各ファイルのインポートを解決
    var result []*ast.SpecFile
    resolved := make(map[string]bool)

    for _, f := range files {
        if err := r.resolveFile(f, resolved, &result); err != nil {
            return nil, err
        }
    }

    return result, nil
}

func (r *ImportResolver) resolveFile(file *ast.SpecFile, resolved map[string]bool, result *[]*ast.SpecFile) error {
    if resolved[file.Path] {
        return nil
    }

    if r.visiting[file.Path] {
        return &errors.CycleError{Cycle: []string{file.Path}}
    }

    r.visiting[file.Path] = true

    // インポートを先に解決
    for _, imp := range file.Imports {
        importPath := r.resolvePath(file.Path, imp.Path)

        imported, ok := r.parsed[importPath]
        if !ok {
            // パースが必要
            var err error
            imported, err = r.parser.ParseFile(importPath)
            if err != nil {
                return &errors.ImportError{
                    Pos:     imp.Pos,
                    Path:    imp.Path,
                    Message: err.Error(),
                }
            }
            r.parsed[importPath] = imported
        }

        if err := r.resolveFile(imported, resolved, result); err != nil {
            return err
        }
    }

    delete(r.visiting, file.Path)
    resolved[file.Path] = true
    *result = append(*result, file)

    return nil
}

func (r *ImportResolver) resolvePath(from, importPath string) string {
    if filepath.IsAbs(importPath) {
        return importPath
    }
    dir := filepath.Dir(from)
    return filepath.Join(dir, importPath)
}
```

**責務**: インポートの解決とトポロジカルソート

---

### internal/infrastructure/filesystem/reader.go

```go
package filesystem

import (
    "os"
    "path/filepath"
)

// Reader はファイルシステムからの読み込みを行う
type Reader struct{}

// NewReader は新しい Reader を作成する
func NewReader() *Reader {
    return &Reader{}
}

// ReadFile はファイルを読み込む
func (r *Reader) ReadFile(path string) (string, error) {
    content, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }
    return string(content), nil
}

// FindPactFiles はディレクトリ内の .pact ファイルを検索する
func (r *Reader) FindPactFiles(dir string, exclude []string) ([]string, error) {
    return r.findFiles(dir, ".pact", exclude)
}

// FindSourceFiles はディレクトリ内のソースファイルを検索する
func (r *Reader) FindSourceFiles(dir string, exclude []string) ([]string, error) {
    // 一般的なソースファイル拡張子
    extensions := []string{".go", ".ts", ".js", ".py", ".java", ".rs"}
    var files []string

    for _, ext := range extensions {
        found, err := r.findFiles(dir, ext, exclude)
        if err != nil {
            return nil, err
        }
        files = append(files, found...)
    }

    return files, nil
}

func (r *Reader) findFiles(dir, ext string, exclude []string) ([]string, error) {
    var files []string

    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        if filepath.Ext(path) != ext {
            return nil
        }

        // 除外パターンをチェック
        for _, pattern := range exclude {
            matched, _ := filepath.Match(pattern, path)
            if matched {
                return nil
            }
        }

        files = append(files, path)
        return nil
    })

    return files, err
}
```

**責務**: ファイルシステムからの読み込み

---

### internal/infrastructure/filesystem/writer.go

```go
package filesystem

import (
    "os"
    "path/filepath"
)

// Writer はファイルシステムへの書き込みを行う
type Writer struct{}

// NewWriter は新しい Writer を作成する
func NewWriter() *Writer {
    return &Writer{}
}

// WriteFile はファイルに書き込む
func (w *Writer) WriteFile(path string, content []byte) error {
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    return os.WriteFile(path, content, 0644)
}

// EnsureDir はディレクトリが存在することを保証する
func (w *Writer) EnsureDir(path string) error {
    return os.MkdirAll(path, 0755)
}
```

**責務**: ファイルシステムへの書き込み

---

## internal/interfaces/cli/ - CLI

### internal/interfaces/cli/root.go

```go
package cli

import (
    "github.com/spf13/cobra"
)

// NewRootCmd はルートコマンドを作成する
func NewRootCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "pact",
        Short: "Pact - Generate diagrams from specifications",
        Long:  `Pact generates multiple diagram types from a single specification file.`,
    }

    cmd.AddCommand(
        NewInitCmd(),
        NewGenerateCmd(),
        NewValidateCmd(),
        NewCheckCmd(),
        NewWatchCmd(),
    )

    return cmd
}
```

**責務**: CLI のルート、サブコマンドの登録

---

### internal/interfaces/cli/init.go

```go
package cli

import (
    "fmt"

    "github.com/spf13/cobra"

    "pact/internal/domain/config"
    infraConfig "pact/internal/infrastructure/config"
)

// NewInitCmd は init コマンドを作成する
func NewInitCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "init",
        Short: "Initialize a new Pact project",
        RunE: func(cmd *cobra.Command, args []string) error {
            loader := infraConfig.NewLoader()
            cfg := config.Default()

            if err := loader.Save(".", cfg); err != nil {
                return err
            }

            fmt.Println("Created .pactconfig")
            return nil
        },
    }

    return cmd
}
```

**責務**: プロジェクトの初期化

---

### internal/interfaces/cli/generate.go

```go
package cli

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"

    "pact/internal/application/parser"
    "pact/internal/application/project"
    "pact/internal/application/renderer"
    "pact/internal/application/transformer"
    infraConfig "pact/internal/infrastructure/config"
    infraFilesystem "pact/internal/infrastructure/filesystem"
    infraParser "pact/internal/infrastructure/parser"
    infraResolver "pact/internal/infrastructure/resolver"
    infraRenderer "pact/internal/infrastructure/renderer/svg"
)

// NewGenerateCmd は generate コマンドを作成する
func NewGenerateCmd() *cobra.Command {
    var (
        output       string
        diagramTypes []string
    )

    cmd := &cobra.Command{
        Use:   "generate [path]",
        Short: "Generate diagrams from pact files",
        Args:  cobra.MaximumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            // 設定を読み込み
            configLoader := infraConfig.NewLoader()
            projectRoot, err := configLoader.FindProjectRoot(".")
            if err != nil {
                return err
            }

            cfg, err := configLoader.Load(projectRoot)
            if err != nil {
                return err
            }

            if output != "" {
                cfg.OutputDir = output
            }

            // 依存性を構築
            reader := infraFilesystem.NewReader()
            writer := infraFilesystem.NewWriter()
            fileParser := infraParser.NewFileParser(reader)
            parserSvc := parser.NewService(fileParser)
            resolverSvc := infraResolver.NewImportResolver(cfg.PactRoot, fileParser)
            projectSvc := project.NewService(cfg, parserSvc, resolverSvc, reader)
            transformerSvc := transformer.NewService()
            svgRenderer := infraRenderer.NewRenderer()
            rendererSvc := renderer.NewService(svgRenderer)

            // プロジェクトを読み込み
            files, err := projectSvc.LoadProject()
            if err != nil {
                return err
            }

            fmt.Printf("Loaded %d pact files\n", len(files))

            // 図を生成
            if err := writer.EnsureDir(cfg.OutputDir); err != nil {
                return err
            }

            for _, file := range files {
                compName := file.Component.Name

                // クラス図
                if shouldGenerate(diagramTypes, "class") && cfg.DiagramEnabled("class") {
                    diagram, err := transformerSvc.ToClassDiagram([]*ast.SpecFile{file}, transformer.ClassOptions{
                        IncludeTypes:      true,
                        IncludeInterfaces: true,
                    })
                    if err != nil {
                        return err
                    }

                    outPath := filepath.Join(cfg.OutputDir, compName, "class.svg")
                    if err := writer.EnsureDir(filepath.Dir(outPath)); err != nil {
                        return err
                    }

                    f, err := os.Create(outPath)
                    if err != nil {
                        return err
                    }
                    if err := rendererSvc.RenderClass(diagram, f); err != nil {
                        f.Close()
                        return err
                    }
                    f.Close()
                    fmt.Printf("Generated %s\n", outPath)
                }

                // シーケンス図（各フローに対して）
                if shouldGenerate(diagramTypes, "sequence") && cfg.DiagramEnabled("sequence") {
                    for _, flow := range file.Component.Body.Flows {
                        diagram, err := transformerSvc.ToSequenceDiagram(file, flow.Name, transformer.SequenceOptions{
                            IncludeReturnMessages: true,
                        })
                        if err != nil {
                            return err
                        }

                        outPath := filepath.Join(cfg.OutputDir, compName, fmt.Sprintf("sequence_%s.svg", flow.Name))
                        f, err := os.Create(outPath)
                        if err != nil {
                            return err
                        }
                        if err := rendererSvc.RenderSequence(diagram, f); err != nil {
                            f.Close()
                            return err
                        }
                        f.Close()
                        fmt.Printf("Generated %s\n", outPath)
                    }
                }

                // ステートマシン図
                if shouldGenerate(diagramTypes, "state") && cfg.DiagramEnabled("state") {
                    for _, states := range file.Component.Body.States {
                        diagram, err := transformerSvc.ToStateDiagram(file, states.Name, transformer.StateOptions{})
                        if err != nil {
                            return err
                        }

                        outPath := filepath.Join(cfg.OutputDir, compName, fmt.Sprintf("state_%s.svg", states.Name))
                        f, err := os.Create(outPath)
                        if err != nil {
                            return err
                        }
                        if err := rendererSvc.RenderState(diagram, f); err != nil {
                            f.Close()
                            return err
                        }
                        f.Close()
                        fmt.Printf("Generated %s\n", outPath)
                    }
                }

                // フローチャート
                if shouldGenerate(diagramTypes, "flow") && cfg.DiagramEnabled("flow") {
                    for _, flow := range file.Component.Body.Flows {
                        diagram, err := transformerSvc.ToFlowchart(file, flow.Name, transformer.FlowOptions{
                            IncludeSwimlanes: true,
                        })
                        if err != nil {
                            return err
                        }

                        outPath := filepath.Join(cfg.OutputDir, compName, fmt.Sprintf("flow_%s.svg", flow.Name))
                        f, err := os.Create(outPath)
                        if err != nil {
                            return err
                        }
                        if err := rendererSvc.RenderFlow(diagram, f); err != nil {
                            f.Close()
                            return err
                        }
                        f.Close()
                        fmt.Printf("Generated %s\n", outPath)
                    }
                }
            }

            return nil
        },
    }

    cmd.Flags().StringVarP(&output, "output", "o", "", "Output directory")
    cmd.Flags().StringSliceVarP(&diagramTypes, "type", "t", []string{"all"}, "Diagram types: class, sequence, state, flow, all")

    return cmd
}

func shouldGenerate(types []string, target string) bool {
    for _, t := range types {
        if t == "all" || t == target {
            return true
        }
    }
    return false
}
```

**責務**: 図の生成

---

### internal/interfaces/cli/validate.go

```go
package cli

import (
    "fmt"

    "github.com/spf13/cobra"

    "pact/internal/application/parser"
    infraConfig "pact/internal/infrastructure/config"
    infraFilesystem "pact/internal/infrastructure/filesystem"
    infraParser "pact/internal/infrastructure/parser"
)

// NewValidateCmd は validate コマンドを作成する
func NewValidateCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "validate [path]",
        Short: "Validate pact files",
        Args:  cobra.MaximumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            configLoader := infraConfig.NewLoader()
            projectRoot, err := configLoader.FindProjectRoot(".")
            if err != nil {
                return err
            }

            cfg, err := configLoader.Load(projectRoot)
            if err != nil {
                return err
            }

            reader := infraFilesystem.NewReader()
            fileParser := infraParser.NewFileParser(reader)
            parserSvc := parser.NewService(fileParser)

            paths, err := reader.FindPactFiles(cfg.PactRoot, cfg.Exclude)
            if err != nil {
                return err
            }

            hasError := false
            for _, path := range paths {
                _, err := parserSvc.ParseFile(path)
                if err != nil {
                    fmt.Printf("ERROR %s: %v\n", path, err)
                    hasError = true
                } else {
                    fmt.Printf("OK    %s\n", path)
                }
            }

            if hasError {
                return fmt.Errorf("validation failed")
            }

            fmt.Printf("\nAll %d files valid\n", len(paths))
            return nil
        },
    }

    return cmd
}
```

**責務**: 構文チェック

---

### internal/interfaces/cli/check.go

```go
package cli

import (
    "fmt"

    "github.com/spf13/cobra"

    "pact/internal/application/parser"
    "pact/internal/application/project"
    infraConfig "pact/internal/infrastructure/config"
    infraFilesystem "pact/internal/infrastructure/filesystem"
    infraParser "pact/internal/infrastructure/parser"
    infraResolver "pact/internal/infrastructure/resolver"
)

// NewCheckCmd は check コマンドを作成する
func NewCheckCmd() *cobra.Command {
    var missing bool

    cmd := &cobra.Command{
        Use:   "check",
        Short: "Check project consistency",
        RunE: func(cmd *cobra.Command, args []string) error {
            configLoader := infraConfig.NewLoader()
            projectRoot, err := configLoader.FindProjectRoot(".")
            if err != nil {
                return err
            }

            cfg, err := configLoader.Load(projectRoot)
            if err != nil {
                return err
            }

            reader := infraFilesystem.NewReader()
            fileParser := infraParser.NewFileParser(reader)
            parserSvc := parser.NewService(fileParser)
            resolverSvc := infraResolver.NewImportResolver(cfg.PactRoot, fileParser)
            projectSvc := project.NewService(cfg, parserSvc, resolverSvc, reader)

            if missing {
                missingSpecs, err := projectSvc.CheckMissing()
                if err != nil {
                    return err
                }

                if len(missingSpecs) == 0 {
                    fmt.Println("All source files have corresponding .pact files")
                    return nil
                }

                fmt.Printf("Found %d source files without .pact:\n", len(missingSpecs))
                for _, m := range missingSpecs {
                    fmt.Printf("  %s -> %s\n", m.SourcePath, m.ExpectedPact)
                }
                return fmt.Errorf("missing pact files")
            }

            // デフォルト: 全チェック
            fmt.Println("Running all checks...")

            // 1. パース可能か
            files, err := projectSvc.LoadProject()
            if err != nil {
                return err
            }
            fmt.Printf("✓ All %d pact files parse successfully\n", len(files))

            // 2. 対応するソースがあるか
            missingSpecs, err := projectSvc.CheckMissing()
            if err != nil {
                return err
            }
            if len(missingSpecs) > 0 {
                fmt.Printf("⚠ %d source files without .pact\n", len(missingSpecs))
            } else {
                fmt.Println("✓ All source files have .pact")
            }

            return nil
        },
    }

    cmd.Flags().BoolVar(&missing, "missing", false, "Check for source files without .pact")

    return cmd
}
```

**責務**: プロジェクトの整合性チェック

---

### internal/interfaces/cli/watch.go

```go
package cli

import (
    "fmt"
    "path/filepath"

    "github.com/fsnotify/fsnotify"
    "github.com/spf13/cobra"

    infraConfig "pact/internal/infrastructure/config"
)

// NewWatchCmd は watch コマンドを作成する
func NewWatchCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "watch",
        Short: "Watch pact files and regenerate on change",
        RunE: func(cmd *cobra.Command, args []string) error {
            configLoader := infraConfig.NewLoader()
            projectRoot, err := configLoader.FindProjectRoot(".")
            if err != nil {
                return err
            }

            cfg, err := configLoader.Load(projectRoot)
            if err != nil {
                return err
            }

            watcher, err := fsnotify.NewWatcher()
            if err != nil {
                return err
            }
            defer watcher.Close()

            // .pact ディレクトリを監視
            pactRoot := filepath.Join(projectRoot, cfg.PactRoot)
            if err := filepath.Walk(pactRoot, func(path string, info os.FileInfo, err error) error {
                if err != nil {
                    return err
                }
                if info.IsDir() {
                    return watcher.Add(path)
                }
                return nil
            }); err != nil {
                return err
            }

            fmt.Printf("Watching %s for changes...\n", pactRoot)

            for {
                select {
                case event, ok := <-watcher.Events:
                    if !ok {
                        return nil
                    }
                    if event.Op&fsnotify.Write == fsnotify.Write {
                        if filepath.Ext(event.Name) == ".pact" {
                            fmt.Printf("Changed: %s\n", event.Name)
                            // TODO: generate コマンドを実行
                        }
                    }
                case err, ok := <-watcher.Errors:
                    if !ok {
                        return nil
                    }
                    fmt.Printf("Error: %v\n", err)
                }
            }
        },
    }

    return cmd
}
```

**責務**: ファイル監視と自動再生成

---

## pkg/pact/ - 公開 API

### pkg/pact/api.go

```go
package pact

import (
    "io"

    "pact/internal/application/parser"
    "pact/internal/application/renderer"
    "pact/internal/application/transformer"
    "pact/internal/domain/ast"
    "pact/internal/domain/diagram/class"
    "pact/internal/domain/diagram/flow"
    "pact/internal/domain/diagram/sequence"
    "pact/internal/domain/diagram/state"
    infraFilesystem "pact/internal/infrastructure/filesystem"
    infraParser "pact/internal/infrastructure/parser"
    infraRenderer "pact/internal/infrastructure/renderer/svg"
)

// Client は Pact の公開 API
type Client struct {
    parserSvc      *parser.Service
    transformerSvc *transformer.Service
    rendererSvc    *renderer.Service
}

// New は新しい Client を作成する
func New() *Client {
    reader := infraFilesystem.NewReader()
    fileParser := infraParser.NewFileParser(reader)
    svgRenderer := infraRenderer.NewRenderer()

    return &Client{
        parserSvc:      parser.NewService(fileParser),
        transformerSvc: transformer.NewService(),
        rendererSvc:    renderer.NewService(svgRenderer),
    }
}

// ParseFile は .pact ファイルをパースする
func (c *Client) ParseFile(path string) (*ast.SpecFile, error) {
    return c.parserSvc.ParseFile(path)
}

// ParseString は文字列をパースする
func (c *Client) ParseString(content, filename string) (*ast.SpecFile, error) {
    return c.parserSvc.ParseString(content, filename)
}

// ToClassDiagram は AST からクラス図を生成する
func (c *Client) ToClassDiagram(files []*ast.SpecFile) (*class.Diagram, error) {
    return c.transformerSvc.ToClassDiagram(files, transformer.ClassOptions{
        IncludeTypes:      true,
        IncludeInterfaces: true,
    })
}

// ToSequenceDiagram は AST からシーケンス図を生成する
func (c *Client) ToSequenceDiagram(file *ast.SpecFile, flowName string) (*sequence.Diagram, error) {
    return c.transformerSvc.ToSequenceDiagram(file, flowName, transformer.SequenceOptions{
        IncludeReturnMessages: true,
    })
}

// ToStateDiagram は AST からステートマシン図を生成する
func (c *Client) ToStateDiagram(file *ast.SpecFile, statesName string) (*state.Diagram, error) {
    return c.transformerSvc.ToStateDiagram(file, statesName, transformer.StateOptions{})
}

// ToFlowchart は AST からフローチャートを生成する
func (c *Client) ToFlowchart(file *ast.SpecFile, flowName string) (*flow.Diagram, error) {
    return c.transformerSvc.ToFlowchart(file, flowName, transformer.FlowOptions{
        IncludeSwimlanes: true,
    })
}

// RenderClassDiagram はクラス図を SVG として出力する
func (c *Client) RenderClassDiagram(diagram *class.Diagram, w io.Writer) error {
    return c.rendererSvc.RenderClass(diagram, w)
}

// RenderSequenceDiagram はシーケンス図を SVG として出力する
func (c *Client) RenderSequenceDiagram(diagram *sequence.Diagram, w io.Writer) error {
    return c.rendererSvc.RenderSequence(diagram, w)
}

// RenderStateDiagram はステートマシン図を SVG として出力する
func (c *Client) RenderStateDiagram(diagram *state.Diagram, w io.Writer) error {
    return c.rendererSvc.RenderState(diagram, w)
}

// RenderFlowchart はフローチャートを SVG として出力する
func (c *Client) RenderFlowchart(diagram *flow.Diagram, w io.Writer) error {
    return c.rendererSvc.RenderFlow(diagram, w)
}
```

**責務**: 外部ライブラリとしての公開 API