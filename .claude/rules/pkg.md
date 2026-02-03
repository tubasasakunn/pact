---
paths:
  - "pkg/**/*.go"
---

# pkg/ ディレクトリルール

外部ライブラリとして利用される公開 API（ファサードパターン）。内部実装の詳細を隠蔽し、安定した API を提供する。

## 依存制約

- `internal/application/` と `internal/infrastructure/` を組み合わせる
- `internal/domain/` の型を公開型エイリアスとして re-export 可能
- 外部パッケージへの依存は最小限

## API 設計

- `Client` 構造体を通じて全機能を提供
- メソッドシグネチャは安定させる（破壊的変更を避ける）
- Options パターンで拡張性を確保
- 戻り値の型は domain 層の型を type alias で公開

```go
type (
    SpecFile      = ast.SpecFile
    ComponentDecl = ast.ComponentDecl
    // ... 公開が必要な型のみ
)
```

## コーディング規約

- godoc コメントは英語で記述
- 全公開メソッドに godoc コメント必須
- メソッドは `error` を返す（panic しない）
- `io.Writer` を受け取り、出力先をユーザーが制御可能にする
- コンストラクタは `New()` で統一

## 禁止事項

- ビジネスロジックの実装（application 層に委譲）
- 内部型の直接公開（type alias を使う）
- グローバル変数・グローバル状態
- `init()` 関数の使用
