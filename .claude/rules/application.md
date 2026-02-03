# internal/application/ ディレクトリルール

## 責務

ビジネスロジック層。AST からダイアグラムモデルへの変換、意味解析バリデーション、
パイプライン全体のオーケストレーションを担当する。

## ディレクトリ構成（リファクタリング後）

```
internal/application/
├── transformer/             # AST → Diagram 変換
│   ├── transformer.go       # 共通インターフェース・ファクトリ・ヘルパー
│   ├── class.go             # ClassTransformer
│   ├── sequence.go          # SequenceTransformer
│   ├── state.go             # StateTransformer
│   ├── flow.go              # FlowTransformer
│   └── options.go           # TransformOptions, SequenceOptions, StateOptions, FlowOptions
│
├── validator/               # 意味解析バリデーション
│   ├── validator.go         # Validator 構造体（オーケストレーター）
│   ├── component.go         # コンポーネント検証ルール
│   ├── type.go              # 型定義検証ルール（重複フィールド等）
│   ├── flow.go              # フロー検証ルール（ステップの検証等）
│   └── state.go             # 状態遷移検証ルール（遷移の整合性等）
│
└── service/                 # アプリケーションサービス
    └── pipeline.go          # DiagramService: パース → バリデーション → 変換 → レンダリング
```

## ルール

### 依存制約

- `internal/domain/` にのみ依存する
- `internal/infrastructure/` への **直接依存禁止**
- infrastructure の機能が必要な場合はインターフェースを定義し、DI で注入
- `pkg/` への依存禁止
- `cmd/` への依存禁止

### Transformer パッケージ (`transformer/`)

#### 共通インターフェース

```go
// transformer.go に定義
type Transformer[D any] interface {
    Transform(files []*ast.SpecFile, opts any) (D, error)
}
```

※ Go 1.21 では generics の制約があるため、実用的にはダイアグラム型ごとの具象型でも可。
ただし Transform メソッドのシグネチャパターンは統一する。

#### 各 Transformer の構造

```go
type ClassTransformer struct{}

func NewClassTransformer() *ClassTransformer

func (t *ClassTransformer) Transform(
    files []*ast.SpecFile,
    opts *TransformOptions,
) (*class.Diagram, error)
```

- 各 Transformer は状態を持たない（stateless）
- `New*()` コンストラクタで生成
- `Transform()` は副作用なし（純粋関数）

#### Options パターン

```go
// options.go に集約
type TransformOptions struct { ... }
type SequenceOptions struct { FlowName string }
type StateOptions struct { StatesName string }
type FlowOptions struct { FlowName string }
```

### Validator パッケージ (`validator/`)

- `Validator` がオーケストレーターとして各検証ルールを呼び出す
- 検証ルールは責務ごとにファイルを分割
- エラーは `domain/errors/` の型を使用
- 警告（Warning）とエラー（Error）を区別する

#### 検証ルールの分類

| ファイル | 検証内容 |
|---------|---------|
| `component.go` | コンポーネント名重複、関係先の存在チェック |
| `type.go` | 型名重複、フィールド名重複、enum 値重複 |
| `flow.go` | フロー名重複、ステップの型チェック、参照先の存在 |
| `state.go` | 状態名重複、遷移先の存在、initial/final 状態の整合性 |

### Service パッケージ (`service/`)

- パイプライン全体を組み立てるオーケストレーター
- infrastructure 層の機能をインターフェース経由で利用
- DI（依存性注入）でテスタビリティを確保

```go
type DiagramService struct {
    parser   Parser      // インターフェース
    renderer Renderer    // インターフェース
}

func NewDiagramService(parser Parser, renderer Renderer) *DiagramService
```

### コーディング規約

- ファイルは 500 行以内を目安
- 関数は 50 行以内を目安
- テーブル駆動テストで変換ロジックをテスト
- エラーメッセージにはコンテキスト情報（コンポーネント名、行番号等）を含める

### テスト

- 各 transformer に対応する `*_test.go` を同パッケージに配置
- テストデータは `testdata/` を使用
- ゴールデンテスト（期待出力との比較）推奨

### 禁止事項

- ファイル I/O の直接実行（infrastructure 層の責務）
- SVG/HTML 等の出力フォーマット生成（renderer の責務）
- トークン解析・字句解析（parser の責務）
- `os` パッケージの使用
- グローバル状態の変更
