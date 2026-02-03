---
paths:
  - "cmd/**/*.go"
---

# cmd/ ディレクトリルール

CLI エントリポイント。ユーザー入力の受け取り・引数パース・出力表示のみを担当する薄いレイヤー。

## 依存制約

- `pkg/pact/` のみに依存する
- `internal/` パッケージを直接 import してはならない
- ビジネスロジックを含めてはならない

## コーディング規約

- 各コマンドは `cmd_<command>.go` に配置
- `main.go` はコマンドディスパッチと `printUsage()` のみ
- 引数パースはコマンドファイル内の `parse*Options()` 関数で行う
- エラー出力は `fmt.Fprintf(os.Stderr, ...)` で統一
- `.pact` ファイルのグロブ展開ロジックは `expandFiles()` で共通化

## コマンドの構造パターン

```go
type <command>Options struct { /* コマンド固有のオプション */ }
func parse<Command>Options(args []string) (*<command>Options, error)
func cmd<Command>(args []string) error
```

## 禁止事項

- AST の直接操作
- パーサー・レンダラーの直接呼び出し
- ビジネスロジック（変換・バリデーション）の実装
