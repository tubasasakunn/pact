# Pact リファクタリング計画

## 現状の課題分析

### 構造上の問題

| ID | 重要度 | 問題 | 対象ファイル | 状態 |
|----|--------|------|-------------|------|
| S-001 | 高 | God File: CLI 全コマンドが 1 ファイル（553行） | `cmd/pact/main.go` | **完了** |
| S-002 | 高 | God File: 全 SVG レンダラーが 1 ファイル（2412行） | `infrastructure/renderer/svg/renderer.go` | **完了** |
| S-003 | 高 | 責務混在: renderer パッケージに無関係なファイル | `renderer/{cache,export,lsp,theme,watcher}.go` | **完了** |
| S-004 | 中 | Validator が 1 ファイル（869行）に全検証ロジック | `application/validator/validator.go` | **完了** |
| S-005 | 中 | インターフェース不在: 具象型への直接依存 | `pkg/pact/api.go` | **完了** |
| S-006 | 中 | DI 不在: Client 構造体が空でテスタビリティが低い | `pkg/pact/api.go` | **完了** |
| S-007 | 低 | Transformer に共通インターフェースがない | `application/transformer/` | **完了** |
| S-008 | 低 | パイプラインオーケストレーター不在 | `application/` | **完了** |

---

## リファクタリング後の構成（実装済み）

```
pact/
├── cmd/pact/                        # CLI（薄いレイヤー）
│   ├── main.go                      # エントリポイント + printUsage のみ (78行)
│   ├── cmd_generate.go              # generate コマンド
│   ├── cmd_validate.go              # validate コマンド
│   ├── cmd_check.go                 # check コマンド
│   ├── cmd_init.go                  # init コマンド
│   ├── cmd_watch.go                 # watch コマンド
│   └── helpers.go                   # expandFiles 共通ヘルパー
│
├── pkg/pact/                        # 公開 API（DI 導入済み）
│   ├── api.go                       # Client + service/renderer DI
│   └── api_test.go
│
├── internal/
│   ├── domain/                      # 純粋ドメインモデル
│   │   ├── ast/
│   │   │   ├── node.go              # (旧 nodes.go)
│   │   │   ├── type.go              # (旧 types.go)
│   │   │   ├── expr.go              # (旧 expressions.go)
│   │   │   ├── stmt.go              # (旧 statements.go)
│   │   │   ├── state.go             # (旧 states.go)
│   │   │   ├── position.go
│   │   │   └── visitor.go
│   │   ├── diagram/{common,class,sequence,state,flow}/
│   │   ├── config/
│   │   └── errors/
│   │
│   ├── application/                 # ビジネスロジック
│   │   ├── transformer/
│   │   │   ├── transformer.go       # パッケージドキュメント
│   │   │   ├── options.go           # TransformOptions, SequenceOptions, StateOptions, FlowOptions
│   │   │   ├── class.go
│   │   │   ├── sequence.go
│   │   │   ├── state.go
│   │   │   └── flow.go
│   │   │
│   │   ├── validator/
│   │   │   ├── validator.go         # Validator オーケストレーター (204行)
│   │   │   ├── component.go         # コンポーネント検証 (115行)
│   │   │   ├── type.go              # 型検証 (295行)
│   │   │   ├── flow.go              # フロー検証 (186行)
│   │   │   └── state.go             # 状態検証 (157行)
│   │   │
│   │   └── service/                 # アプリケーションサービス
│   │       ├── pipeline.go          # DiagramService (DI ベース)
│   │       └── pipeline_test.go     # モックを使ったテスト
│   │
│   ├── infrastructure/
│   │   ├── parser/                  # （変更なし）
│   │   │
│   │   ├── renderer/
│   │   │   ├── renderer.go          # ClassRenderer/SequenceRenderer/StateRenderer/FlowRenderer インターフェース
│   │   │   ├── svg/
│   │   │   │   ├── class.go         # ClassRenderer 実装
│   │   │   │   ├── sequence.go      # SequenceRenderer 実装
│   │   │   │   ├── state.go         # StateRenderer 実装
│   │   │   │   ├── flow.go          # FlowRenderer 実装
│   │   │   │   └── layout.go        # 共通レイアウト・ユーティリティ
│   │   │   └── canvas/              # （変更なし）
│   │   │
│   │   ├── resolver/                # （変更なし）
│   │   ├── config/                  # （変更なし）
│   │   │
│   │   ├── export/                  # (旧 renderer/export.go)
│   │   │   └── exporter.go
│   │   │
│   │   ├── cache/                   # (旧 renderer/cache.go)
│   │   │   └── cache.go
│   │   │
│   │   ├── theme/                   # (旧 renderer/theme.go)
│   │   │   └── theme.go
│   │   │
│   │   └── watcher/                 # (旧 renderer/watcher.go)
│   │       └── watcher.go
│   │
│   └── testutil/
│
├── test/                            # （変更なし）
├── testdata/                        # （変更なし）
├── sample/                          # （変更なし）
└── docs/                            # （変更なし）
```

**削除済み:**
- `infrastructure/renderer/lsp.go` — LSP はレンダラーの責務ではない。実装時に別パッケージに再作成

---

## フェーズ別実施結果

### Phase 1: 責務の再配置（パッケージ移動） — **完了**

#### Step 1.1: renderer パッケージから無関係なファイルを分離 — **完了**

- `renderer/cache.go` → `infrastructure/cache/cache.go`
- `renderer/export.go` → `infrastructure/export/exporter.go`
- `renderer/theme.go` → `infrastructure/theme/theme.go`
- `renderer/watcher.go` → `infrastructure/watcher/watcher.go`
- `renderer/lsp.go` → 削除

#### Step 1.2: renderer/svg/renderer.go の分割 — **完了**

2412行の God File を 5 ファイルに分割:
- `svg/class.go` — ClassRenderer（レイヤー割り当て、バリセンター法、エッジルーティング含む）
- `svg/sequence.go` — SequenceRenderer
- `svg/state.go` — StateRenderer
- `svg/flow.go` — FlowRenderer
- `svg/layout.go` — 共通ユーティリティ（maxInt, minInt, abs, sqrt）+ ノート描画関数

#### AST ファイル名の変更 — **完了**

- `nodes.go` → `node.go`、`types.go` → `type.go`、`statements.go` → `stmt.go`
- `expressions.go` → `expr.go`、`states.go` → `state.go`

---

### Phase 2: CLI の分割 — **完了**

553行の `main.go` を 7 ファイルに分割:
- `main.go` (78行) — エントリポイント + printUsage のみ
- `cmd_generate.go` — generate コマンド + ダイアグラム生成ヘルパー
- `cmd_validate.go` — validate コマンド
- `cmd_check.go` — check コマンド
- `cmd_init.go` — init コマンド
- `cmd_watch.go` — watch コマンド
- `helpers.go` — expandFiles 共通ヘルパー

---

### Phase 3: Validator の分割 — **完了**

869行の `validator.go` を 5 ファイルに分割:
- `validator.go` (204行) — オーケストレーター
- `component.go` (115行) — コンポーネント検証
- `type.go` (295行) — 型検証
- `flow.go` (186行) — フロー検証
- `state.go` (157行) — 状態検証

全ファイル 500 行以下のガイドラインを達成。

---

### Phase 4: インターフェース導入と DI — **完了**

#### Renderer インターフェース
`infrastructure/renderer/renderer.go` に 4 つのインターフェースを定義:
- `ClassRenderer`、`SequenceRenderer`、`StateRenderer`、`FlowRenderer`

#### Transformer Options 分離
`transformer/options.go` にオプション型を集約:
- `TransformOptions`、`SequenceOptions`、`StateOptions`、`FlowOptions`

#### DiagramService（パイプラインオーケストレーター）
`application/service/pipeline.go`:
- インターフェース経由で Renderer を注入
- `GenerateXxx()` — transform + render のパイプライン
- `TransformXxx()` — transform のみ

#### pkg/pact Client の DI 化
- `Client` 構造体に `service.DiagramService` を保持
- `New()` でデフォルト SVG レンダラーを注入
- 公開 API のシグネチャは変更なし（後方互換）

---

### Phase 5: テスト強化 — **完了**

- `service/pipeline_test.go` を追加（モックレンダラーを使った DI テスト）
- 全テスト通過確認:
  - ユニットテスト: `go test -race ./internal/...` — PASS
  - 統合テスト: `go test -race ./test/integration/...` — PASS
  - E2E テスト: `go test -race ./test/e2e/...` — PASS
  - API テスト: `go test -race ./pkg/...` — PASS
