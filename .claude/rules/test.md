# test/ ディレクトリルール

## 責務

統合テスト（integration）とエンドツーエンドテスト（E2E）を格納する。
ユニットテストは各パッケージ内の `*_test.go` に配置し、ここには含めない。

## ディレクトリ構成

```
test/
├── e2e/                     # E2E テスト（CLI 全体のテスト）
│   └── cli_test.go          # pact コマンドの E2E テスト
│
└── integration/             # 統合テスト（複数パッケージの結合テスト）
    ├── parse_transform_test.go    # パース → 変換パイプライン
    ├── import_resolution_test.go  # マルチファイル import 解決
    └── full_pipeline_test.go      # 全ダイアグラム型の完全パイプライン

testdata/                    # テストデータ（test/ と同階層）
├── valid/                   # 正常系テストデータ
│   ├── order_service.pact
│   ├── complex.pact
│   ├── with_types.pact
│   ├── with_relations.pact
│   └── with_annotations.pact
│
└── invalid/                 # 異常系テストデータ
    ├── syntax/              # 構文エラー
    ├── semantic/            # 意味エラー
    └── import/              # インポートエラー
```

## ルール

### テストの分類

| 分類 | 場所 | 実行コマンド | 対象 |
|------|------|-------------|------|
| ユニットテスト | `internal/**/*_test.go` | `make test-unit` | 単一パッケージ |
| 統合テスト | `test/integration/` | `make test-integration` | 複数パッケージ |
| E2E テスト | `test/e2e/` | `make test-e2e` | CLI 全体 |
| API テスト | `pkg/**/*_test.go` | `make test-api` | 公開 API |

### テスト命名規則

```go
// テスト関数
func Test<対象>_<シナリオ>(t *testing.T) { ... }

// サブテスト
t.Run("<条件>/<期待結果>", func(t *testing.T) { ... })

// テストケース ID（統合テスト）
// I001, I002, ... : パース統合テスト
// I010, I011, ... : インポート解決テスト
// I020, I021, ... : パイプライン統合テスト
// I030, I031, ... : エラーハンドリング統合テスト
```

### テストデータ (`testdata/`)

- `.pact` ファイルのみ配置
- ファイル名で用途が分かるようにする
- テストコード内で `testdata/` のパスを使って読み込む
- テストデータは変更の影響範囲を最小化するため、できるだけ小さく保つ

### コーディング規約

- テーブル駆動テスト推奨
- `t.Helper()` でヘルパー関数をマーク
- `t.Parallel()` で並列実行可能なテストは並列化
- `-race` フラグ付きで全テスト実行
- テストヘルパーは `internal/testutil/` に集約

### 統合テスト (`test/integration/`)

- 複数の internal パッケージを結合してテスト
- パイプライン全体の動作を検証
- エラーパスも網羅的にテスト
- テストケース ID を付与して追跡可能に

### E2E テスト (`test/e2e/`)

- `os/exec` で `pact` バイナリを実行
- 標準出力・標準エラー出力・exit code を検証
- 一時ディレクトリで実行し、テスト後にクリーンアップ
- `TestMain` でバイナリのビルドを事前実行

### 禁止事項

- テストデータへのハードコードされた絶対パス
- テスト間の依存（各テストは独立して実行可能にする）
- 外部サービスへの依存（ネットワーク接続等）
- テスト結果のファイルシステムへの永続的な書き込み（一時ディレクトリを使用）
