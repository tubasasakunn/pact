# Pact リファクタリング計画

## 現状の課題分析

### 構造上の問題

| ID | 重要度 | 問題 | 対象ファイル |
|----|--------|------|-------------|
| S-001 | 高 | God File: CLI 全コマンドが 1 ファイル（553行） | `cmd/pact/main.go` |
| S-002 | 高 | God File: 全 SVG レンダラーが 1 ファイル（700行超） | `infrastructure/renderer/svg/renderer.go` |
| S-003 | 高 | 責務混在: renderer パッケージに無関係なファイル | `renderer/{cache,export,lsp,theme,watcher}.go` |
| S-004 | 中 | Validator が 1 ファイル（500行超）に全検証ロジック | `application/validator/validator.go` |
| S-005 | 中 | インターフェース不在: 具象型への直接依存 | `pkg/pact/api.go` |
| S-006 | 中 | DI 不在: Client 構造体が空でテスタビリティが低い | `pkg/pact/api.go` |
| S-007 | 低 | Transformer に共通インターフェースがない | `application/transformer/` |
| S-008 | 低 | パイプラインオーケストレーター不在 | `application/` |

### 現在のファイル配置の問題

```
renderer/
├── cache.go       ← レンダリングとは無関係。独立パッケージにすべき
├── export.go      ← エクスポート機能。独立パッケージにすべき
├── lsp.go         ← LSP は完全に別機能。renderer に置くべきではない
├── theme.go       ← テーマ設定。独立パッケージにすべき
├── watcher.go     ← ファイル監視。renderer に置くべきではない
├── svg/
│   └── renderer.go  ← 4 ダイアグラム分の全レンダラーが 1 ファイル
└── canvas/
```

---

## リファクタリング後の目標構成

```
pact/
├── cmd/pact/                        # CLI（薄いレイヤー）
│   ├── main.go                      # エントリポイントのみ
│   ├── cmd_generate.go              # generate コマンド
│   ├── cmd_validate.go              # validate コマンド
│   ├── cmd_check.go                 # check コマンド
│   ├── cmd_init.go                  # init コマンド
│   └── cmd_watch.go                 # watch コマンド
│
├── pkg/pact/                        # 公開 API
│   ├── pact.go                      # Client + DI
│   ├── options.go                   # 公開オプション型
│   └── pact_test.go
│
├── internal/
│   ├── domain/                      # 純粋ドメインモデル（変更なし）
│   │   ├── ast/
│   │   ├── diagram/{common,class,sequence,state,flow}/
│   │   ├── config/
│   │   └── errors/
│   │
│   ├── application/                 # ビジネスロジック
│   │   ├── transformer/
│   │   │   ├── transformer.go       # [新規] 共通インターフェース
│   │   │   ├── class.go
│   │   │   ├── sequence.go
│   │   │   ├── state.go
│   │   │   ├── flow.go
│   │   │   └── options.go           # [新規] オプション型を分離
│   │   │
│   │   ├── validator/
│   │   │   ├── validator.go         # オーケストレーター（縮小）
│   │   │   ├── component.go         # [新規] コンポーネント検証
│   │   │   ├── type.go              # [新規] 型検証
│   │   │   ├── flow.go              # [新規] フロー検証
│   │   │   └── state.go             # [新規] 状態検証
│   │   │
│   │   └── service/                 # [新規] アプリケーションサービス
│   │       └── pipeline.go          # パイプラインオーケストレーター
│   │
│   ├── infrastructure/
│   │   ├── parser/                  # （変更なし）
│   │   │   ├── token.go
│   │   │   ├── lexer.go
│   │   │   └── parser.go
│   │   │
│   │   ├── renderer/
│   │   │   ├── renderer.go          # [新規] Renderer インターフェース
│   │   │   ├── svg/
│   │   │   │   ├── class.go         # [分割] ClassRenderer
│   │   │   │   ├── sequence.go      # [分割] SequenceRenderer
│   │   │   │   ├── state.go         # [分割] StateRenderer
│   │   │   │   ├── flow.go          # [分割] FlowRenderer
│   │   │   │   └── layout.go        # [分割] 共通レイアウトアルゴリズム
│   │   │   └── canvas/              # （変更なし）
│   │   │
│   │   ├── resolver/                # （変更なし）
│   │   ├── config/                  # （変更なし）
│   │   │
│   │   ├── export/                  # [移動] renderer/export.go から
│   │   │   └── exporter.go
│   │   │
│   │   ├── cache/                   # [移動] renderer/cache.go から
│   │   │   └── cache.go
│   │   │
│   │   ├── theme/                   # [移動] renderer/theme.go から
│   │   │   └── theme.go
│   │   │
│   │   └── watcher/                 # [移動] renderer/watcher.go から
│   │       └── watcher.go
│   │
│   └── testutil/
│
├── test/                            # （変更なし）
├── testdata/                        # （変更なし）
├── sample/                          # （変更なし）
└── docs/                            # （変更なし）
```

**削除対象:**
- `infrastructure/renderer/lsp.go` → 別パッケージに移動するか、LSP 実装時まで削除

---

## フェーズ別実施計画

### Phase 1: 責務の再配置（パッケージ移動）

低リスク。ファイル移動とパッケージ名変更のみ。ロジック変更なし。

#### Step 1.1: renderer パッケージから無関係なファイルを分離

```
renderer/cache.go   → infrastructure/cache/cache.go
renderer/export.go  → infrastructure/export/exporter.go
renderer/theme.go   → infrastructure/theme/theme.go
renderer/watcher.go → infrastructure/watcher/watcher.go
renderer/lsp.go     → 削除（または infrastructure/lsp/lsp.go に移動）
```

**作業:**
1. 新パッケージディレクトリ作成
2. ファイル移動、package 宣言変更
3. import パス更新（依存元があれば）
4. `go build ./...` で確認
5. `make test` で全テスト通過確認

#### Step 1.2: renderer/svg/renderer.go の分割

```
renderer/svg/renderer.go (700行超)
  → renderer/svg/class.go      (ClassRenderer)
  → renderer/svg/sequence.go   (SequenceRenderer)
  → renderer/svg/state.go      (StateRenderer)
  → renderer/svg/flow.go       (FlowRenderer)
  → renderer/svg/layout.go     (共通レイアウトアルゴリズム)
```

**作業:**
1. `renderer.go` 内のコードを機能ごとに分割
2. 共通のレイアウトロジック（`assignLayers`, `optimizeLayerOrder` 等）を `layout.go` に抽出
3. 各テストファイルは既に分割済みのためそのまま
4. `make test` で全テスト通過確認

---

### Phase 2: CLI の分割

中リスク。main.go のリファクタリング。

#### Step 2.1: main.go をコマンドごとに分割

```
cmd/pact/main.go (553行)
  → cmd/pact/main.go          (エントリポイント + printUsage のみ)
  → cmd/pact/cmd_generate.go  (cmdGenerate, generateOptions, helper関数)
  → cmd/pact/cmd_validate.go  (cmdValidate)
  → cmd/pact/cmd_check.go     (cmdCheck)
  → cmd/pact/cmd_init.go      (cmdInit)
  → cmd/pact/cmd_watch.go     (cmdWatch)
```

**作業:**
1. 各コマンド関数を個別ファイルに移動
2. 共通ヘルパー（ファイル展開等）を `helpers.go` に抽出
3. `main.go` はコマンドディスパッチのみに
4. `make test-e2e` で確認

---

### Phase 3: Validator の分割

中リスク。大きなファイルの責務分離。

#### Step 3.1: validator.go を検証カテゴリごとに分割

```
validator/validator.go (500行超)
  → validator/validator.go     (Validator 構造体、Validate メソッド)
  → validator/component.go    (validateComponent, 名前重複チェック)
  → validator/type.go         (validateDuplicateFields, 型検証)
  → validator/flow.go         (validateFlow, ステップ検証)
  → validator/state.go        (validateStates, 遷移検証)
```

**作業:**
1. 検証メソッドをカテゴリごとに分割
2. `Validator` 構造体は `validator.go` に残す
3. private メソッドは同パッケージ内なので import 変更不要
4. `make test-unit` で確認

---

### Phase 4: インターフェース導入と DI

高リスク。設計変更を伴う。慎重に進める。

#### Step 4.1: Renderer インターフェース定義

```go
// infrastructure/renderer/renderer.go
type ClassRenderer interface {
    Render(d *class.Diagram, w io.Writer) error
}
type SequenceRenderer interface {
    Render(d *sequence.Diagram, w io.Writer) error
}
// ... 各ダイアグラム型ごと
```

#### Step 4.2: Transformer 共通インターフェース

```go
// application/transformer/transformer.go
// 注意: Go のジェネリクス制約により、完全な共通インターフェースは難しい場合がある。
// パターンの統一（メソッド名、引数の構造）を優先する。
```

#### Step 4.3: pkg/pact の Client に DI を導入

```go
type Client struct {
    parser   *parser.Parser
    // 将来的にインターフェースに変更可能な構造
}

func New(opts ...Option) *Client
```

#### Step 4.4: application/service/pipeline.go 作成

```go
type DiagramService struct { ... }
func (s *DiagramService) GenerateAll(files []string, output string) error
```

**作業:**
1. インターフェース定義
2. 既存の具象型がインターフェースを満たすことを確認
3. Client 構造体にフィールド追加
4. pipeline.go で全体オーケストレーション
5. 全テスト通過確認

---

### Phase 5: テスト強化

低リスク。コード品質向上。

#### Step 5.1: テストヘルパーの整理

- `internal/testutil/` に共通ヘルパーを集約
- AST 構築ヘルパー、比較ヘルパー等

#### Step 5.2: テストカバレッジ改善

- 新規作成ファイルのテスト追加
- `service/pipeline.go` のテスト
- エッジケースの追加

---

## 各フェーズの依存関係

```
Phase 1 (パッケージ移動)
    ↓
Phase 2 (CLI 分割)      ← Phase 1 と並行可能
    ↓
Phase 3 (Validator 分割) ← Phase 1 と並行可能
    ↓
Phase 4 (インターフェース + DI) ← Phase 1, 2, 3 完了後
    ↓
Phase 5 (テスト強化)     ← Phase 4 完了後
```

## 各フェーズの検証基準

| フェーズ | 検証コマンド | 合格基準 |
|---------|-------------|---------|
| 全フェーズ共通 | `go build ./...` | コンパイルエラーなし |
| 全フェーズ共通 | `make test` | 全テスト通過 |
| 全フェーズ共通 | `make lint` | lint エラーなし |
| Phase 1 | `go vet ./...` | vet 警告なし |
| Phase 4 | `make test-coverage` | カバレッジ低下なし |

## リスクと注意点

1. **import パスの変更**: パッケージ移動時に全 import パスを確実に更新する
2. **パッケージ名の衝突**: `renderer/renderer.go` のパッケージ名と型名の衝突に注意
3. **テストの依存**: テストが内部実装に依存している場合、リファクタリング時にテストも更新が必要
4. **Phase 4 の破壊的変更**: インターフェース導入は pkg/pact の公開 API に影響する可能性がある
5. **段階的実行**: 各 Phase を完了しテストを通してから次に進む。一度に全て変更しない

## AST ファイル名の変更（Phase 1 で実施）

現在の AST ファイル名をより明確にする:

```
ast/nodes.go       → ast/node.go       (単数形に統一)
ast/types.go       → ast/type.go       (単数形に統一)
ast/statements.go  → ast/stmt.go       (省略形に統一)
ast/expressions.go → ast/expr.go       (省略形に統一)
ast/states.go      → ast/state.go      (単数形に統一)
ast/position.go    → (変更なし)
ast/visitor.go     → (変更なし)
```

※ Go の慣習: ファイル名は単数形・小文字・短い名前が推奨。
