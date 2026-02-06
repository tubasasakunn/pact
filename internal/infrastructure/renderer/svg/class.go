package svg

import (
	"io"

	"pact/internal/domain/diagram/class"
	"pact/internal/infrastructure/renderer/canvas"
)

// ClassRenderer はクラス図をSVGにレンダリングする
type ClassRenderer struct{}

// NewClassRenderer は新しいClassRendererを作成する
func NewClassRenderer() *ClassRenderer {
	return &ClassRenderer{}
}

// Render はクラス図をSVGにレンダリングする
func (r *ClassRenderer) Render(diagram *class.Diagram, w io.Writer) error {
	c := canvas.New()

	// テンプレートレジストリを適用（シャドウ、フォント、カラーテーマ）
	registry := canvas.NewBuiltinRegistry()
	registry.ApplyTo(c)

	// 各ノードのサイズを事前計算
	nodeSizes := make(map[string]struct{ width, height int })
	for _, node := range diagram.Nodes {
		w := r.calculateNodeWidth(node)
		h := r.calculateNodeHeight(node)
		nodeSizes[node.ID] = struct{ width, height int }{w, h}
	}

	// エッジ接続情報を構築（レイアウト用）
	// 継承・実装エッジはレイアウト上の方向を反転させる
	// （親クラス・インターフェースが上、子クラス・実装が下に配置されるように）
	outgoing := make(map[string][]string) // from -> []to
	incoming := make(map[string][]string) // to -> []from
	for _, edge := range diagram.Edges {
		if edge.Type == class.EdgeTypeInheritance || edge.Type == class.EdgeTypeImplementation {
			// 継承・実装: 親/インターフェースを上に配置するため方向を反転
			outgoing[edge.To] = append(outgoing[edge.To], edge.From)
			incoming[edge.From] = append(incoming[edge.From], edge.To)
		} else {
			outgoing[edge.From] = append(outgoing[edge.From], edge.To)
			incoming[edge.To] = append(incoming[edge.To], edge.From)
		}
	}

	// レイヤー割り当て（トポロジカルソート風）
	layers := r.assignLayers(diagram.Nodes, incoming, outgoing)

	// バリセンター法でレイヤー内のノード順序を最適化（交差最小化）
	layers = r.optimizeLayerOrder(layers, incoming, outgoing, nodeSizes)

	// 各レイヤーの幅と高さを計算
	layerWidths := make([]int, len(layers))
	layerHeights := make([]int, len(layers))
	for i, layer := range layers {
		totalWidth := 0
		maxHeight := 0
		for j, nodeID := range layer {
			size := nodeSizes[nodeID]
			totalWidth += size.width
			if j < len(layer)-1 {
				totalWidth += 40 // ノード間マージン
			}
			if size.height > maxHeight {
				maxHeight = size.height
			}
		}
		layerWidths[i] = totalWidth
		layerHeights[i] = maxHeight
	}

	// 最大幅を計算
	maxLayerWidth := 0
	for _, w := range layerWidths {
		if w > maxLayerWidth {
			maxLayerWidth = w
		}
	}
	canvasWidth := maxLayerWidth + 100
	if canvasWidth < 800 {
		canvasWidth = 800
	}

	// ノード位置を計算（レイヤーベース）
	nodePositions := make(map[string]struct{ x, y, width, height int })
	y := 50
	layerMargin := 60 // レイヤー間のマージン

	for i, layer := range layers {
		// レイヤーを中央揃え
		layerWidth := layerWidths[i]
		startX := (canvasWidth - layerWidth) / 2
		if startX < 50 {
			startX = 50
		}

		x := startX
		for _, nodeID := range layer {
			size := nodeSizes[nodeID]
			nodePositions[nodeID] = struct{ x, y, width, height int }{x, y, size.width, size.height}
			x += size.width + 40
		}

		y += layerHeights[i] + layerMargin
	}

	// ノードを描画
	nodeMap := make(map[string]class.Node)
	for _, node := range diagram.Nodes {
		nodeMap[node.ID] = node
	}
	for _, node := range diagram.Nodes {
		pos := nodePositions[node.ID]
		r.renderNode(c, node, pos.x, pos.y, pos.width)
	}

	// ノートの位置を考慮してキャンバスサイズを計算
	totalHeight := y + 50
	totalWidth := canvasWidth

	// ノートが存在する場合、その位置を考慮
	if len(diagram.Notes) > 0 {
		noteWidth := 100
		noteHeight := 40
		simplePositions := make(map[string]struct{ x, y int })
		for id, pos := range nodePositions {
			simplePositions[id] = struct{ x, y int }{pos.x, pos.y}
		}
		for _, note := range diagram.Notes {
			noteX, noteY := calculateNotePosition(note, simplePositions, noteWidth, noteHeight)
			// 右端と下端を更新
			if noteX+noteWidth+50 > totalWidth {
				totalWidth = noteX + noteWidth + 50
			}
			if noteY+noteHeight+50 > totalHeight {
				totalHeight = noteY + noteHeight + 50
			}
		}
	}

	if totalHeight < 600 {
		totalHeight = 600
	}
	if totalWidth < 800 {
		totalWidth = 800
	}
	c.SetSize(totalWidth, totalHeight)

	// エッジをレンダリング（改良版：接続点の分散配置）
	// 各ノードからの出力エッジ数と入力エッジ数をカウント
	outgoingCount := make(map[string]int)
	incomingCount := make(map[string]int)
	outgoingIndex := make(map[string]int)
	incomingIndex := make(map[string]int)

	for _, edge := range diagram.Edges {
		outgoingCount[edge.From]++
		incomingCount[edge.To]++
	}

	for _, edge := range diagram.Edges {
		fromPos, fromOk := nodePositions[edge.From]
		toPos, toOk := nodePositions[edge.To]
		if !fromOk || !toOk {
			continue
		}

		// 接続点のインデックスを取得
		outIdx := outgoingIndex[edge.From]
		inIdx := incomingIndex[edge.To]
		outgoingIndex[edge.From]++
		incomingIndex[edge.To]++

		if edge.Type == class.EdgeTypeInheritance || edge.Type == class.EdgeTypeImplementation {
			// 継承・実装エッジは常に垂直接続（子のtop → 親のbottom）
			r.renderVerticalEdge(c, edge, fromPos, toPos,
				outIdx, outgoingCount[edge.From],
				inIdx, incomingCount[edge.To],
				nodePositions)
		} else {
			// その他のエッジは従来の分散接続点計算
			fromX, fromY, toX, toY := r.calculateDistributedEndpoints(
				fromPos, toPos,
				outIdx, outgoingCount[edge.From],
				inIdx, incomingCount[edge.To],
			)
			r.renderEdgeImproved(c, edge, fromX, fromY, toX, toY, nodePositions)
		}
	}

	// ノートをレンダリング
	if len(diagram.Notes) > 0 {
		simplePositions := make(map[string]struct{ x, y int })
		for id, pos := range nodePositions {
			simplePositions[id] = struct{ x, y int }{pos.x, pos.y}
		}
		renderNotes(c, diagram.Notes, simplePositions)
	}

	_, err := c.WriteTo(w)
	return err
}

func (r *ClassRenderer) renderNode(c *canvas.Canvas, node class.Node, x, y, width int) {
	height := r.calculateNodeHeight(node)
	lineHeight := 20
	padding := 15
	sectionGap := 10

	// ノード本体（テーマカラー＋ドロップシャドウ）
	c.Rect(x, y, width, height,
		canvas.Fill(canvas.ColorNodeFill),
		canvas.Stroke(canvas.ColorNodeStroke),
		canvas.StrokeWidth(2),
		canvas.Filter("drop-shadow"),
	)

	centerX := x + width/2
	textY := y + padding + 12 // ベースライン調整

	// ステレオタイプ
	if node.Stereotype != "" {
		c.Text(centerX, textY, "<<"+node.Stereotype+">>",
			canvas.TextAnchor("middle"),
			canvas.Fill(canvas.ColorEdgeLabel),
			canvas.FontStyle("italic"),
		)
		textY += lineHeight
	}

	// 名前
	c.Text(centerX, textY, node.Name,
		canvas.TextAnchor("middle"),
		canvas.Fill(canvas.ColorNodeText),
		canvas.FontWeight("bold"),
	)
	textY += lineHeight

	// 属性セクション
	if len(node.Attributes) > 0 {
		// 区切り線
		lineY := textY - lineHeight/2 + sectionGap/2
		c.Line(x, lineY, x+width, lineY, canvas.Stroke(canvas.ColorSectionLine))
		textY += sectionGap

		for _, attr := range node.Attributes {
			vis := visibilitySymbol(attr.Visibility)
			c.Text(x+10, textY, vis+attr.Name+": "+attr.Type,
				canvas.Fill(canvas.ColorNodeText),
			)
			textY += lineHeight
		}
	}

	// メソッドセクション
	if len(node.Methods) > 0 {
		// 区切り線
		lineY := textY - lineHeight/2 + sectionGap/2
		c.Line(x, lineY, x+width, lineY, canvas.Stroke(canvas.ColorSectionLine))
		textY += sectionGap

		for _, method := range node.Methods {
			vis := visibilitySymbol(class.Visibility(method.Visibility))
			methodStr := r.formatMethod(method)
			c.Text(x+10, textY, vis+methodStr,
				canvas.Fill(canvas.ColorNodeText),
			)
			textY += lineHeight
		}
	}
}

// formatMethod はメソッドシグネチャを整形する
func (r *ClassRenderer) formatMethod(method class.Method) string {
	// パラメータリストを構築
	params := ""
	for i, p := range method.Params {
		if i > 0 {
			params += ", "
		}
		if p.Name != "" && p.Type != "" {
			params += p.Name + ": " + p.Type
		} else if p.Type != "" {
			params += p.Type
		} else if p.Name != "" {
			params += p.Name
		}
	}

	// asyncプレフィックス
	prefix := ""
	if method.Async {
		prefix = "async "
	}

	// 戻り型
	returnType := ""
	if method.ReturnType != "" {
		returnType = ": " + method.ReturnType
	}

	// throws句
	throws := ""
	if len(method.Throws) > 0 {
		throws = " throws "
		for i, t := range method.Throws {
			if i > 0 {
				throws += ", "
			}
			throws += t
		}
	}

	return prefix + method.Name + "(" + params + ")" + returnType + throws
}

func visibilitySymbol(v class.Visibility) string {
	switch v {
	case class.VisibilityPublic:
		return "+ "
	case class.VisibilityPrivate:
		return "- "
	case class.VisibilityProtected:
		return "# "
	case class.VisibilityPackage:
		return "~ "
	default:
		return "+ "
	}
}
