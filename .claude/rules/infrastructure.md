---
paths:
  - "internal/infrastructure/**/*.go"
---

# internal/infrastructure/ ディレクトリルール

外部 I/O・技術的関心事の実装層。パーサー、レンダラー、ファイル解決、設定ファイル読み込み、キャッシュ、エクスポート、テーマ、ファイル監視など、技術的な実装詳細を担当する。

## 依存制約

- `internal/domain/` にのみ依存する
- `internal/application/`、`pkg/`、`cmd/` への依存禁止
- infrastructure 内のサブパッケージ間の依存は最小限
  - `svg/` → `canvas/` は許可
  - `svg/` → `parser/` は禁止（レイヤー違反）

## Parser パッケージ (`parser/`)

- 再帰下降パーサー（recursive descent）
- Lexer: `NextToken()` で 1 トークンずつ返す、位置情報追跡
- Parser: Lexer からトークンを受け取り `ast.SpecFile` を生成
- エラーリカバリ: パースエラー時も可能な限り解析を継続
- `domain/errors/ParseError` でエラー報告

## Renderer パッケージ (`renderer/`)

- ダイアグラム型ごとに分離されたインターフェース（`ClassRenderer`, `SequenceRenderer` 等）
- SVG 実装は各ダイアグラム型ごとに独立ファイル
- 共通レイアウトアルゴリズムは `layout.go` に抽出
- Canvas: SVG 要素の低レベル描画プリミティブ

## その他パッケージ

- **resolver/**: import 文のファイルパス解決、循環参照検出
- **config/**: `.pactconfig` の YAML 読み込み
- **cache/**: SHA256 ベースのレンダリングキャッシュ
- **export/**: SVG → 他形式変換インターフェース
- **theme/**: ダイアグラム配色テーマ定義
- **watcher/**: ファイル変更監視インターフェース

## コーディング規約

- `io.Writer` / `io.Reader` を使い、具象的なファイルパスを受け取らない
- テスト時は `bytes.Buffer` を渡せるようにする
- マジックナンバーは定数化
- 各ファイルは 500 行以内を目標

## 禁止事項

- ビジネスロジック（変換ルール等）の実装
- AST の意味解析
- `application/` パッケージの型への依存
- 標準出力への直接書き込み（`fmt.Println` 等）
