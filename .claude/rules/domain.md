# internal/domain/ ディレクトリルール

## 責務

純粋なドメインモデル。アプリケーションの核となるデータ構造・型・バリューオブジェクト・ドメインエラーを定義する。
外部依存ゼロの最も安定したレイヤー。

## ディレクトリ構成（リファクタリング後）

```
internal/domain/
├── ast/                     # 抽象構文木
│   ├── node.go              # SpecFile, ComponentDecl, ComponentBody 等のコア型
│   ├── type.go              # TypeDecl, EnumDecl, FieldDecl 等の型定義
│   ├── expr.go              # Expr インターフェースと式ノード
│   ├── stmt.go              # Step インターフェースとステートメントノード
│   ├── state.go             # StatesDecl, StateDecl, Transition 等
│   ├── position.go          # Position 型（ソースコード位置情報）
│   └── visitor.go           # Visitor インターフェースと Walk 関数
│
├── diagram/                 # ダイアグラムモデル
│   ├── common/
│   │   └── types.go         # Diagram インターフェース, Note, Annotation 等
│   ├── class/
│   │   └── model.go         # Diagram, Node, Edge, Attribute, Method
│   ├── sequence/
│   │   └── model.go         # Diagram, Participant, Event インターフェース
│   ├── state/
│   │   └── model.go         # Diagram, State, Transition, Trigger
│   └── flow/
│       └── model.go         # Diagram, Node, Edge, Swimlane
│
├── config/
│   └── config.go            # Config 構造体（YAML タグ付き値オブジェクト）
│
└── errors/
    └── errors.go            # ParseError, SemanticError, MultiError 等
```

## ルール

### 依存制約（最重要）

- **Go 標準ライブラリのみ** 依存可能
- `internal/application/` への依存禁止
- `internal/infrastructure/` への依存禁止
- `pkg/` への依存禁止
- 外部パッケージ（`gopkg.in/yaml.v3` 等）への依存禁止
  - 例外: `config.go` の YAML タグは構造体タグであり import は不要

### AST パッケージ (`ast/`)

- 全 AST ノードは `Position` フィールドを持つ（エラー報告用）
- インターフェース型（`Step`, `Expr`, `Trigger`）はマーカーメソッドで識別
  ```go
  type Step interface {
      stepNode()       // マーカーメソッド
      GetPos() Position
  }
  ```
- Visitor パターンで走査可能にする
- AST ノードは不変（immutable）として設計する

### Diagram パッケージ (`diagram/`)

- 各ダイアグラム型は独立したサブパッケージ
- `common/` に共有型を配置（`Note`, `Annotation`, `Diagram` インターフェース）
- ダイアグラムモデルは AST から独立（直接参照しない）
- 各 `Diagram` は `common.Diagram` インターフェースを実装

### Error パッケージ (`errors/`)

- ポジション情報付きエラー型
- `error` インターフェースを実装
- `MultiError` で複数エラーの集約
- エラーメッセージはユーザーフレンドリーに

### コーディング規約

- 型定義のみ（ビジネスロジックは application 層）
- メソッドはゲッター・フォーマッター等の純粋関数のみ
- goroutine を生成しない
- I/O 操作を行わない
- ファイルごとの責務を明確に分割（1 ファイル 1 概念）

### ファイル命名規則

- `model.go`: ダイアグラムモデル定義
- `node.go` / `expr.go` / `stmt.go`: AST ノード定義（概念別）
- `errors.go`: エラー型定義
- `visitor.go`: Visitor パターン実装
- `position.go`: 位置情報型

### 禁止事項

- `fmt.Println` 等の標準出力への書き込み
- `os` パッケージの使用（ファイル I/O）
- `net` パッケージの使用（ネットワーク I/O）
- グローバル変数の変更（定数は可）
- `sync` パッケージの使用（並行処理は上位レイヤー）
