---
paths:
  - "sample/**/*"
  - "docs/sample/**/*"
  - "scripts/generate-gallery.sh"
---

# Sample Files & SVG Generation

## ディレクトリ構成

```
sample/pact/              → サンプル .pact ファイル（ソース）
  class/ state/ flow/ sequence/

docs/sample/              → GitHub Pages 用（生成された SVG）
  index.html              → commit 一覧ページ
  commit/<commit-id>/     → コミットID単位でグループ化
```

## パターンテンプレート

パターン定義: `internal/infrastructure/renderer/canvas/pattern_*.go`

- **Class**: InheritanceTree2/3/4, InterfaceImpl2/3/4, Composition2/3/4, Diamond, Layered3x2/3x3
- **State**: LinearStates2/3/4, BinaryChoice, StateLoop, StarTopology
- **Flow**: IfElse, IfElseIfElse, WhileLoop, Sequential3/4
- **Sequence**: RequestResponse, Callback, Chain3/4, FanOut

## パターンプレビュー

```bash
go run ./cmd/pattern-preview    # → pattern-preview/index.html
```

## GitHub Pages

- **URL**: `https://<username>.github.io/pact/sample/`
- `sample/pact/` や `internal/` の変更が main にマージされると GitHub Actions が自動生成
- 手動: `./scripts/generate-gallery.sh`
- Pages 設定: Branch `main`, Folder `/docs`
