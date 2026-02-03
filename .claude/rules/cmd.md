# cmd/ ディレクトリルール

## 責務

CLI エントリポイント。ユーザー入力の受け取り・引数パース・出力表示のみを担当する薄いレイヤー。

## ディレクトリ構成（リファクタリング後）

```
cmd/pact/
├── main.go              # エントリポイント（コマンドディスパッチのみ）
├── cmd_generate.go      # generate コマンド
├── cmd_validate.go      # validate コマンド
├── cmd_check.go         # check コマンド
├── cmd_init.go          # init コマンド
└── cmd_watch.go         # watch コマンド
```

## ルール

### 依存制約

- `pkg/pact/` のみに依存する
- `internal/` パッケージを直接 import してはならない
- ビジネスロジックを含めてはならない

### コーディング規約

- 各コマンドは独立したファイルに配置（`cmd_<command>.go`）
- `main.go` はコマンドディスパッチと `printUsage()` のみ
- 引数パースはコマンドファイル内の `parse*Options()` 関数で行う
- エラー出力は `fmt.Fprintf(os.Stderr, ...)` で統一
- エラー時の exit code は `os.Exit(1)` で統一
- 成功メッセージは `fmt.Println()` / `fmt.Printf()` で統一

### コマンドの構造パターン

```go
// cmd_<command>.go

type <command>Options struct {
    // コマンド固有のオプション
}

func parse<Command>Options(args []string) (*<command>Options, error) {
    // 引数パース
}

func cmd<Command>(args []string) error {
    opts, err := parse<Command>Options(args)
    if err != nil {
        return err
    }
    // pkg/pact の API を呼び出す
    return nil
}
```

### ファイルパターン展開

- `.pact` ファイルのグロブ展開ロジックは共通化する（`expandFiles()` ヘルパー）
- ディレクトリ指定時は `*.pact` を自動展開

### 禁止事項

- AST の直接操作
- パーサーの直接呼び出し
- レンダラーの直接呼び出し
- ビジネスロジック（変換・バリデーション）の実装
