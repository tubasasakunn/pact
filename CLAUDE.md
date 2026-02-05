# CLAUDE.md - Pact Project Rules

## Project Overview

Pact は `.pact` DSL ファイルから UML ダイアグラム（クラス図・シーケンス図・状態遷移図・フローチャート）を生成する Go 製 CLI ツール兼ライブラリ。

## Build & Test Commands

```bash
# Build
make build                    # go build -o bin/pact ./cmd/pact

# Test (全体)
make test                     # lint + unit + integration + e2e

# Test (個別)
make test-unit                # go test -v -race ./internal/...
make test-integration         # go test -v -race ./test/integration/...
make test-e2e                 # go test -v -race ./test/e2e/...
make test-api                 # go test -v -race ./pkg/...

# 特定パッケージ
make test-pkg PKG=internal/infrastructure/parser

# 特定テスト
make test-run RUN=TestParseComponent

# Lint
make lint                     # golangci-lint run ./...

# Format
make fmt                      # go fmt ./...

# Coverage
make test-coverage            # coverage.html を生成
```

## Architecture

DDD レイヤードアーキテクチャ + Clean Architecture の原則に従う。

```
cmd/pact/           → CLI エントリポイント（薄いレイヤー）
pkg/pact/           → 公開 API（ファサード）
internal/
  domain/           → 純粋なドメインモデル（外部依存ゼロ）
  application/      → ビジネスロジック（transformer, validator, service）
  infrastructure/   → 外部I/O（parser, renderer, resolver, config）
```

### 依存ルール（最重要）

```
cmd/ → pkg/ → application/ → domain/
                    ↓
              infrastructure/  → domain/
```

- **domain/** は他のどのレイヤーにも依存してはならない（Go 標準ライブラリのみ）
- **application/** は domain/ にのみ依存する。infrastructure/ への依存はインターフェース経由
- **infrastructure/** は domain/ にのみ依存する
- **pkg/** は application/ と infrastructure/ を組み合わせて公開 API を提供する
- **cmd/** は pkg/ にのみ依存する

## Go Coding Conventions

### 全般

- Go 1.21 以上
- `go fmt` / `goimports` でフォーマット統一
- `golangci-lint` で静的解析通過必須
- コメントは日本語可（ただし公開パッケージの godoc は英語）
- エラーは `fmt.Errorf("context: %w", err)` でラップ
- `panic` は使用禁止（テスト以外）

### 命名規則

- パッケージ名: 小文字、単一単語（`transformer`、`parser`、`svg`）
- インターフェース: `-er` suffix（`Parser`、`Renderer`、`Resolver`）
- コンストラクタ: `New` prefix（`NewParser()`、`NewClassTransformer()`）
- テストファイル: `*_test.go`
- テスト関数: `Test<対象>_<条件>` 形式

### インターフェース設計

- **"Accept interfaces, return structs"** の原則に従う
- インターフェースは消費側で定義する（実装側ではなく）
- 不必要に大きなインターフェースを定義しない（Interface Segregation）

### エラーハンドリング

- `domain/errors/` に定義されたエラー型を使用
- ポジション情報付きエラー（`ParseError`、`SemanticError`）でユーザーフレンドリーなメッセージ
- `MultiError` で複数エラーの集約
- sentinel error よりカスタムエラー型を優先

### テスト

- テーブル駆動テスト推奨
- `testdata/` ディレクトリに `.pact` テストファイルを配置
- テストヘルパーは `internal/testutil/` に集約
- `-race` フラグ付きでテスト実行

## File Size Guidelines

- 1 ファイル 500 行以内を目安とする
- 500 行を超える場合は責務ごとに分割を検討
- テストファイルは例外（テーブル駆動テストで長くなることは許容）

## Data Flow

```
.pact file → Lexer → Parser → AST → Validator → Transformer → Diagram Model → Renderer → SVG
```

## Key Design Patterns

- **Visitor Pattern**: AST 走査（`ast/visitor.go`）
- **Options Pattern**: Transformer へのオプション渡し
- **Factory Pattern**: `New*()` コンストラクタ
- **Strategy Pattern**: 各ダイアグラム型ごとの Transformer/Renderer

## Dependencies

- 外部依存は最小限に保つ（現在は `gopkg.in/yaml.v3` のみ）
- 新しい外部依存追加時はレビュー必須
- 標準ライブラリで実現可能なら標準ライブラリを使う

## Sample Files & SVG Generation

### ディレクトリ構成

```
sample/
  pact/                 → サンプル .pact ファイル（ソース）
    class/              → クラス図パターン
    state/              → 状態遷移図パターン
    flow/               → フローチャートパターン
    sequence/           → シーケンス図パターン
  commit/               → 生成された SVG ファイル（コミット対象）
    class/
    state/
    flow/
    sequence/
```

### サンプル生成コマンド

コミット時に以下のコマンドで SVG を再生成する：

```bash
# 全サンプルの SVG 生成
for f in sample/pact/class/*.pact; do ./bin/pact generate -o sample/commit/class/ -t class "$f"; done
for f in sample/pact/state/*.pact; do ./bin/pact generate -o sample/commit/state/ -t state "$f"; done
for f in sample/pact/flow/*.pact; do ./bin/pact generate -o sample/commit/flow/ -t flow "$f"; done
for f in sample/pact/sequence/*.pact; do ./bin/pact generate -o sample/commit/sequence/ -t sequence "$f"; done
```

### パターンテンプレート

SVG 生成には以下のパターンテンプレートが使用される（`internal/infrastructure/renderer/canvas/pattern.go`）：

**Class Patterns:**
- InheritanceTree2/3/4 - 継承ツリー（2〜4子クラス）
- InterfaceImpl2/3/4 - インターフェース実装（2〜4実装クラス）
- Composition2/3/4 - コンポジション（2〜4パーツ）
- Diamond - ダイヤモンド依存
- Layered3x2/3x3 - レイヤードアーキテクチャ

**State Patterns:**
- LinearStates2/3/4 - 直線状態遷移（2〜4状態）
- BinaryChoice - 二項選択
- StateLoop - ループ状態
- StarTopology - 星型トポロジー

**Flow Patterns:**
- IfElse, IfElseIfElse - 条件分岐
- WhileLoop - ループ
- Sequential3/4 - 順次処理（3〜4ステップ）

**Sequence Patterns:**
- RequestResponse - リクエスト/レスポンス
- Callback - コールバック
- Chain3/4 - チェーン（3〜4参加者）
- FanOut - ファンアウト

### パターンプレビューツール

パターンテンプレートのプレビューを生成：

```bash
go run ./cmd/pattern-preview
# → pattern-preview/index.html をブラウザで開く
```
