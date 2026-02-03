---
paths:
  - "internal/domain/**/*.go"
---

# internal/domain/ ディレクトリルール

純粋なドメインモデル。アプリケーションの核となるデータ構造・型・バリューオブジェクト・ドメインエラーを定義する。外部依存ゼロの最も安定したレイヤー。

## 依存制約（最重要）

- **Go 標準ライブラリのみ** 依存可能
- `internal/application/`、`internal/infrastructure/`、`pkg/` への依存禁止
- 外部パッケージ（`gopkg.in/yaml.v3` 等）への依存禁止

## AST パッケージ (`ast/`)

- 全 AST ノードは `Position` フィールドを持つ（エラー報告用）
- インターフェース型（`Step`, `Expr`, `Trigger`）はマーカーメソッドで識別
- Visitor パターンで走査可能にする
- AST ノードは不変（immutable）として設計

## Diagram パッケージ (`diagram/`)

- 各ダイアグラム型は独立したサブパッケージ（`class/`, `sequence/`, `state/`, `flow/`）
- `common/` に共有型を配置
- ダイアグラムモデルは AST から独立（直接参照しない）

## Error パッケージ (`errors/`)

- ポジション情報付きエラー型、`MultiError` で複数エラー集約

## コーディング規約

- 型定義のみ（ビジネスロジックは application 層）
- メソッドはゲッター・フォーマッター等の純粋関数のみ
- ファイルごとの責務を明確に分割（1 ファイル 1 概念）
- ファイル命名: `model.go`（ダイアグラム）、`node.go`/`expr.go`/`stmt.go`（AST）

## 禁止事項

- `fmt.Println` 等の標準出力への書き込み
- `os`、`net`、`sync` パッケージの使用
- goroutine の生成、I/O 操作
- グローバル変数の変更（定数は可）
