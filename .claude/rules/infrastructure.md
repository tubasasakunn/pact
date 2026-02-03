# internal/infrastructure/ ディレクトリルール

## 責務

外部 I/O・技術的関心事の実装層。パーサー、レンダラー、ファイル解決、設定ファイル読み込み、
キャッシュ、エクスポート、テーマ、ファイル監視など、技術的な実装詳細を担当する。

## ディレクトリ構成（リファクタリング後）

```
internal/infrastructure/
├── parser/                  # 字句解析・構文解析
│   ├── token.go             # TokenType 定義、Token 構造体
│   ├── lexer.go             # Lexer（字句解析器）
│   └── parser.go            # Parser（構文解析器）→ AST 生成
│
├── renderer/                # ダイアグラムレンダリング
│   ├── renderer.go          # Renderer インターフェース定義
│   ├── svg/                 # SVG 実装
│   │   ├── class.go         # ClassRenderer
│   │   ├── sequence.go      # SequenceRenderer
│   │   ├── state.go         # StateRenderer
│   │   ├── flow.go          # FlowRenderer
│   │   └── layout.go        # 共通レイアウトアルゴリズム（Sugiyama 等）
│   └── canvas/              # 2D 描画プリミティブ
│       ├── canvas.go        # Canvas 構造体（SVG 要素バッファ）
│       ├── shapes.go        # 図形描画（矩形、円、菱形等）
│       └── text.go          # テキスト描画・メトリクス
│
├── resolver/                # インポート解決
│   └── resolver.go          # ImportResolver（依存順序解決、循環検出）
│
├── config/                  # 設定ファイル読み込み
│   └── loader.go            # .pactconfig の YAML 読み込み
│
├── export/                  # エクスポート機能
│   └── exporter.go          # SVG → PNG/PDF 変換（Exporter インターフェース）
│
├── cache/                   # レンダリングキャッシュ
│   └── cache.go             # RenderCache（SHA256 ベース）
│
├── theme/                   # テーマ管理
│   └── theme.go             # Theme 構造体、DefaultTheme/DarkTheme/BlueprintTheme
│
└── watcher/                 # ファイル監視
    └── watcher.go           # Watcher インターフェース、WatchEvent
```

## ルール

### 依存制約

- `internal/domain/` にのみ依存する
- `internal/application/` への依存禁止（逆依存にならないように）
- `pkg/` への依存禁止
- `cmd/` への依存禁止
- infrastructure 内のサブパッケージ間の依存は最小限に
  - `svg/` → `canvas/` は許可（描画プリミティブの利用）
  - `svg/` → `parser/` は禁止（レイヤー違反）

### Parser パッケージ (`parser/`)

#### Lexer

- 入力文字列 → トークンストリーム変換
- `NextToken()` で 1 トークンずつ返す
- 位置情報（行・列）を追跡
- 文字列リテラル、数値リテラル、識別子、キーワード、演算子をサポート

#### Parser

- 再帰下降パーサー（recursive descent）
- Lexer からトークンを受け取り `ast.SpecFile` を生成
- エラーリカバリ: パースエラー時も可能な限り解析を継続
- `domain/errors/ParseError` でエラー報告

#### Token

- `TokenType` は `iota` ベースの定数
- 88 種のトークン型（キーワード・演算子・リテラル・記号）
- キーワードマップで識別子とキーワードを区別

### Renderer パッケージ (`renderer/`)

#### インターフェース (`renderer.go`)

```go
type Renderer interface {
    RenderClass(d *class.Diagram, w io.Writer) error
    RenderSequence(d *sequence.Diagram, w io.Writer) error
    RenderState(d *state.Diagram, w io.Writer) error
    RenderFlow(d *flow.Diagram, w io.Writer) error
}
```

※ インターフェースが大きすぎる場合はダイアグラム型ごとに分割：

```go
type ClassRenderer interface {
    Render(d *class.Diagram, w io.Writer) error
}
```

#### SVG 実装 (`svg/`)

- 各ダイアグラム型ごとに独立したファイル（`class.go`, `sequence.go` 等）
- 共通レイアウトアルゴリズムは `layout.go` に抽出
- レイアウトアルゴリズム:
  - レイヤー割り当て（トポロジカルソート）
  - バリセンター法（交差最小化）
  - ノード位置計算
  - エッジルーティング

#### Canvas (`canvas/`)

- SVG 要素の低レベル描画プリミティブ
- テキストメトリクス計算
- 図形描画（矩形、角丸矩形、円、菱形、矢印等）

### Resolver パッケージ (`resolver/`)

- `import` 文のファイルパス解決
- 依存グラフの構築
- 循環参照の検出（`CycleError`）
- 深さ優先探索で解決順序を決定

### その他パッケージ

#### Cache (`cache/`)
- `RenderCache`: SHA256 ベースのキャッシュキー
- `sync.RWMutex` でスレッドセーフ
- LRU に近い eviction（将来改善予定）

#### Export (`export/`)
- SVG → 他形式の変換インターフェース
- 現在は SVG のみサポート（PNG/PDF は将来拡張）

#### Theme (`theme/`)
- ダイアグラムの配色テーマ定義
- `GetTheme(name)` でテーマ取得
- ビルトイン: default, dark, blueprint

#### Watcher (`watcher/`)
- ファイル変更監視のインターフェース定義
- `WatchEvent` でイベント通知
- 実装は将来追加（現在はインターフェースのみ）

### コーディング規約

- `io.Writer` / `io.Reader` を使い、具象的なファイルパスを受け取らない
- テスト時は `bytes.Buffer` を渡せるようにする
- マジックナンバーは定数化（SVG のマージン値、フォントサイズ等）
- レンダラーの各ファイルは 500 行以内を目標
  - 現在の `renderer.go`（700行超）は分割対象

### テスト

- ユニットテストは同パッケージの `*_test.go` に配置
- テストデータは `testdata/` ディレクトリを使用
- SVG 出力のスナップショットテスト推奨
- Lexer/Parser は境界値テストを充実させる

### 禁止事項

- ビジネスロジック（変換ルール等）の実装
- AST の意味解析
- `application/` パッケージの型への依存
- 標準出力への直接書き込み（`fmt.Println` 等）
