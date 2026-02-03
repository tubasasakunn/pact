---
paths:
  - "internal/application/**/*.go"
---

# internal/application/ ディレクトリルール

ビジネスロジック層。AST からダイアグラムモデルへの変換、意味解析バリデーション、パイプライン全体のオーケストレーションを担当する。

## 依存制約

- `internal/domain/` にのみ依存する
- `internal/infrastructure/` への **直接依存禁止**（インターフェース経由で DI）
- `pkg/`、`cmd/` への依存禁止

## Transformer パッケージ (`transformer/`)

- 各 Transformer は状態を持たない（stateless）
- `New*()` コンストラクタで生成、`Transform()` は副作用なし（純粋関数）
- Transform メソッドのシグネチャパターンは統一する
- Options 型は `options.go` に集約

```go
type ClassTransformer struct{}
func NewClassTransformer() *ClassTransformer
func (t *ClassTransformer) Transform(files []*ast.SpecFile, opts *TransformOptions) (*class.Diagram, error)
```

## Validator パッケージ (`validator/`)

- `Validator` がオーケストレーターとして各検証ルールを呼び出す
- 検証ルールは責務ごとにファイル分割: `component.go`, `type.go`, `flow.go`, `state.go`
- エラーは `domain/errors/` の型を使用
- 警告（Warning）とエラー（Error）を区別

## Service パッケージ (`service/`)

- パイプライン全体を組み立てるオーケストレーター
- infrastructure 層の機能をインターフェース経由で利用、DI でテスタビリティ確保

## コーディング規約

- ファイルは 500 行以内、関数は 50 行以内を目安
- テーブル駆動テストで変換ロジックをテスト
- エラーメッセージにはコンテキスト情報（コンポーネント名、行番号等）を含める

## 禁止事項

- ファイル I/O の直接実行（infrastructure 層の責務）
- SVG/HTML 等の出力フォーマット生成（renderer の責務）
- トークン解析・字句解析（parser の責務）
- `os` パッケージの使用、グローバル状態の変更
