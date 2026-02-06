# CLAUDE.md - Pact Project Rules

## Project Overview

`.pact` DSL → UML ダイアグラム（クラス図・シーケンス図・状態遷移図・フローチャート）を生成する Go CLI + ライブラリ。

## Build & Test

```bash
make build                    # go build -o bin/pact ./cmd/pact
make test                     # lint + unit + integration + e2e
make test-unit                # go test -v -race ./internal/...
make test-integration         # go test -v -race ./test/integration/...
make test-e2e                 # go test -v -race ./test/e2e/...
make test-api                 # go test -v -race ./pkg/...
make test-pkg PKG=internal/infrastructure/parser  # 特定パッケージ
make test-run RUN=TestParseComponent              # 特定テスト
make lint                     # golangci-lint run ./...
make fmt                      # go fmt ./...
```

## Architecture

DDD レイヤード + Clean Architecture。

```
cmd/pact/           → CLI エントリポイント（薄いレイヤー）
pkg/pact/           → 公開 API（ファサード）
internal/
  domain/           → 純粋なドメインモデル（外部依存ゼロ）
  application/      → ビジネスロジック（transformer, validator, service）
  infrastructure/   → 外部I/O（parser, renderer, resolver, config）
```

### 依存ルール（MUST: 違反禁止）

```
cmd/ → pkg/ → application/ → domain/
                    ↓
              infrastructure/  → domain/
```

- **cmd/** は `pkg/` のみ import 可。internal/ への直接 import 禁止
- **pkg/** は `application/` + `infrastructure/` を組み合わせて公開 API 提供
- **application/** は `domain/` のみ。infrastructure/ はインターフェース経由
- **infrastructure/** は `domain/` のみ
- **domain/** は Go 標準ライブラリのみ。外部パッケージ依存ゼロ

レイヤー詳細は `.claude/rules/{domain,application,infrastructure,pkg,cmd}.md` 参照。

### Data Flow

```
.pact file → Lexer → Parser → AST → Validator → Transformer → Diagram Model → Renderer → SVG
```

## Key File Locations

| 責務 | ファイル |
|---|---|
| Lexer | `parser/lexer.go`, `lexer_read.go`, `token.go` |
| Parser コア | `parser/parser.go`（オーケストレーション・ヘルパー） |
| Parser 文法要素 | `parser/parser_{component,types,relations,flow,expressions,states,annotations}.go` |
| AST 定義 | `domain/ast/{node,stmt,expr,type,state,visitor}.go` |
| Transformer | `application/transformer/{class,sequence,state,flow}.go` |
| Validator | `application/validator/{validator,component,type,flow,state}.go` |
| Class図 SVG | `svg/class.go`(描画), `class_layout.go`(配置), `class_edge.go`(エッジ) |
| State図 SVG | `svg/state.go`, `state_compound.go`, `state_transition.go` |
| SVG 共通ユーティリティ | `svg/layout.go`（abs, maxInt, minInt, sqrt, 衝突検出, ノート描画） |
| Canvas 描画 | `canvas/canvas.go`(プリミティブ), `canvas/template.go`(テンプレート) |
| パターン定義 | `canvas/pattern.go`(型+レジストリ), `pattern_{class,state,flow,sequence}.go` |
| パターン検出 | `canvas/pattern_detector_{class,state,flow,sequence}.go` |
| デコレーション | `canvas/pattern_decorations.go`(型), `_gradients/_filters/_styles/_render.go` |
| 公開 API | `pkg/pact/api.go`, `preview.go` |
| エラー型 | `domain/errors/errors.go` |

## Coding Conventions

- Go 1.21+, `golangci-lint` 通過必須
- エラーは `fmt.Errorf("context: %w", err)` でラップ。`panic` 禁止（テスト以外）
- コメントは日本語可（公開パッケージの godoc は英語）
- パッケージ: 小文字単一単語。インターフェース: `-er` suffix。コンストラクタ: `New` prefix
- テスト: `Test<対象>_<条件>` 形式、テーブル駆動推奨、`testdata/` 配置、`-race` 付き
- "Accept interfaces, return structs"。インターフェースは消費側で定義
- カスタムエラー型（`domain/errors/`）を sentinel error より優先

### ファイル分割ルール

- **1ファイル 500行以内**（テストファイルは例外）
- 分割命名: `<base>_<concern>.go`（例: `parser_flow.go`, `class_edge.go`, `state_transition.go`）
- Go 同一パッケージ内の分割は `package` 宣言のみで OK。import 変更不要

## Design Patterns

- **Visitor**: AST 走査（`ast/visitor.go`）
- **Options**: Transformer オプション（`transformer/options.go`）
- **Factory**: `New*()` コンストラクタ
- **Strategy**: ダイアグラム型ごとの Transformer/Renderer

## Dependencies

外部依存は `gopkg.in/yaml.v3` のみ。追加時はレビュー必須。標準ライブラリ優先。

## Gotchas

- `cmd/pattern-preview/` は `.gitignore` 対象 → git 操作時は `git add -f` が必要
- `svg/layout.go` に共通ユーティリティ（`abs`, `maxInt`, `minInt`, `sqrt` 等）あり → 他ファイルで再定義しない
- `canvas/pattern.go` にパターン型定義（`PatternRegistry`, `ClassPatternMatch` 等）集約
- Parser エラーリカバリ: `synchronize()` で `}` まで読み飛ばして解析継続
- サンプル/SVG/GitHub Pages 関連は `.claude/rules/samples.md` 参照
