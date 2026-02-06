package svg

import (
	"fmt"
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
				inIdx, incomingCount[edge.To])
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

func (r *ClassRenderer) calculateNodeHeight(node class.Node) int {
	lineHeight := 20
	padding := 15    // 上下のパディング
	sectionGap := 10 // セクション間のギャップ

	height := padding // 上パディング

	// ステレオタイプ
	if node.Stereotype != "" {
		height += lineHeight
	}

	// 名前
	height += lineHeight

	// 属性セクション
	if len(node.Attributes) > 0 {
		height += sectionGap // 区切り線の余白
		height += len(node.Attributes) * lineHeight
	}

	// メソッドセクション
	if len(node.Methods) > 0 {
		height += sectionGap // 区切り線の余白
		height += len(node.Methods) * lineHeight
	}

	height += padding // 下パディング
	return height
}

// optimizeLayerOrder はバリセンター法でレイヤー内のノード順序を最適化する
func (r *ClassRenderer) optimizeLayerOrder(layers [][]string, incoming, outgoing map[string][]string, nodeSizes map[string]struct{ width, height int }) [][]string {
	if len(layers) <= 1 {
		return layers
	}

	// ノードのレイヤー内位置を追跡
	nodePosition := make(map[string]int) // nodeID -> position in layer

	// 初期位置を設定
	for _, layer := range layers {
		for pos, nodeID := range layer {
			nodePosition[nodeID] = pos
		}
	}

	// 反復回数をグラフの複雑さに応じて調整
	// 基本4回 + レイヤー数に応じて追加（最大20回）
	totalNodes := 0
	for _, layer := range layers {
		totalNodes += len(layer)
	}
	iterations := 4 + len(layers)/2
	if totalNodes > 20 {
		iterations += totalNodes / 10
	}
	if iterations > 20 {
		iterations = 20
	}

	for iter := 0; iter < iterations; iter++ {
		// 下向きスイープ（レイヤー1から最後まで）
		for i := 1; i < len(layers); i++ {
			layers[i] = r.reorderLayerByBarycenter(layers[i], layers[i-1], incoming, nodePosition)
			// 位置を更新
			for pos, nodeID := range layers[i] {
				nodePosition[nodeID] = pos
			}
		}

		// 上向きスイープ（最後から1つ前からレイヤー0まで）
		for i := len(layers) - 2; i >= 0; i-- {
			layers[i] = r.reorderLayerByBarycenter(layers[i], layers[i+1], outgoing, nodePosition)
			// 位置を更新
			for pos, nodeID := range layers[i] {
				nodePosition[nodeID] = pos
			}
		}
	}

	return layers
}

// reorderLayerByBarycenter はバリセンター値に基づいてレイヤー内のノードを並び替える
func (r *ClassRenderer) reorderLayerByBarycenter(layer, adjacentLayer []string, connections map[string][]string, nodePosition map[string]int) []string {
	if len(layer) <= 1 {
		return layer
	}

	// 隣接レイヤーのノード位置マップ
	adjacentPos := make(map[string]int)
	for pos, nodeID := range adjacentLayer {
		adjacentPos[nodeID] = pos
	}

	// 各ノードのバリセンター値を計算
	nodes := make([]nodeWithBarycenter, len(layer))
	for i, nodeID := range layer {
		// 隣接レイヤーへの接続を取得
		connectedNodes := connections[nodeID]
		if len(connectedNodes) == 0 {
			// 接続がない場合は現在位置を維持
			nodes[i] = nodeWithBarycenter{nodeID, float64(nodePosition[nodeID]), false}
			continue
		}

		// バリセンター（接続先の平均位置）を計算
		sum := 0.0
		count := 0
		for _, connID := range connectedNodes {
			if pos, ok := adjacentPos[connID]; ok {
				sum += float64(pos)
				count++
			}
		}

		if count > 0 {
			nodes[i] = nodeWithBarycenter{nodeID, sum / float64(count), true}
		} else {
			nodes[i] = nodeWithBarycenter{nodeID, float64(nodePosition[nodeID]), false}
		}
	}

	// バリセンター値でソート（安定ソート）
	r.stableSortByBarycenter(nodes)

	// 結果を返す
	result := make([]string, len(layer))
	for i, n := range nodes {
		result[i] = n.id
	}
	return result
}

// stableSortByBarycenter はバリセンター値で安定ソートする（マージソート: O(n log n)）
func (r *ClassRenderer) stableSortByBarycenter(nodes []nodeWithBarycenter) {
	if len(nodes) <= 1 {
		return
	}
	mergeSortBarycenter(nodes, make([]nodeWithBarycenter, len(nodes)))
}

// assignLayers はノードをレイヤーに割り当てる（Sugiyama法の簡易版）
func (r *ClassRenderer) assignLayers(nodes []class.Node, incoming, outgoing map[string][]string) [][]string {
	// ノードIDのセット
	nodeSet := make(map[string]bool)
	for _, node := range nodes {
		nodeSet[node.ID] = true
	}

	// 入次数を計算
	inDegree := make(map[string]int)
	for _, node := range nodes {
		inDegree[node.ID] = len(incoming[node.ID])
	}

	// レイヤー割り当て
	var layers [][]string
	assigned := make(map[string]bool)

	for len(assigned) < len(nodes) {
		var currentLayer []string

		// 入次数が0のノード（または未処理の依存元がないノード）を現在のレイヤーに追加
		for _, node := range nodes {
			if assigned[node.ID] {
				continue
			}

			// このノードの全ての依存元が既に割り当て済みかチェック
			allDepsAssigned := true
			for _, from := range incoming[node.ID] {
				if nodeSet[from] && !assigned[from] {
					allDepsAssigned = false
					break
				}
			}

			if allDepsAssigned {
				currentLayer = append(currentLayer, node.ID)
			}
		}

		// デッドロック防止: 何も追加できなければ残りを全部追加
		if len(currentLayer) == 0 {
			for _, node := range nodes {
				if !assigned[node.ID] {
					currentLayer = append(currentLayer, node.ID)
				}
			}
		}

		// 割り当て済みにマーク
		for _, id := range currentLayer {
			assigned[id] = true
		}

		if len(currentLayer) > 0 {
			layers = append(layers, currentLayer)
		}
	}

	return layers
}

// calculateDistributedEndpoints は複数エッジの接続点を分散配置する
func (r *ClassRenderer) calculateDistributedEndpoints(
	fromPos, toPos struct{ x, y, width, height int },
	outIdx, outTotal, inIdx, inTotal int,
) (fromX, fromY, toX, toY int) {
	fromCenterX := fromPos.x + fromPos.width/2
	fromCenterY := fromPos.y + fromPos.height/2
	toCenterX := toPos.x + toPos.width/2
	toCenterY := toPos.y + toPos.height/2

	// エッジ配置の分散幅（ノード幅の80%を使用）
	fromSpread := int(float64(fromPos.width) * 0.8)
	toSpread := int(float64(toPos.width) * 0.8)

	// 垂直方向の差が大きい場合（下向き/上向き接続）
	if abs(toCenterY-fromCenterY) > abs(toCenterX-fromCenterX) {
		// 出力点を下端/上端に分散配置（中央を基準に均等分布）
		fromOffset := 0
		if outTotal > 1 {
			fromOffset = (2*outIdx - (outTotal - 1)) * fromSpread / (2 * (outTotal - 1))
		}
		toOffset := 0
		if inTotal > 1 {
			toOffset = (2*inIdx - (inTotal - 1)) * toSpread / (2 * (inTotal - 1))
		}

		if toCenterY > fromCenterY {
			// 下向き
			fromX = fromCenterX + fromOffset
			fromY = fromPos.y + fromPos.height
			toX = toCenterX + toOffset
			toY = toPos.y
		} else {
			// 上向き
			fromX = fromCenterX + fromOffset
			fromY = fromPos.y
			toX = toCenterX + toOffset
			toY = toPos.y + toPos.height
		}
		return
	}

	// 水平方向の接続
	fromHeightSpread := int(float64(fromPos.height) * 0.6)
	toHeightSpread := int(float64(toPos.height) * 0.6)

	fromYOffset := 0
	if outTotal > 1 {
		fromYOffset = (2*outIdx - (outTotal - 1)) * fromHeightSpread / (2 * (outTotal - 1))
	}
	toYOffset := 0
	if inTotal > 1 {
		toYOffset = (2*inIdx - (inTotal - 1)) * toHeightSpread / (2 * (inTotal - 1))
	}

	if toCenterX > fromCenterX {
		// 右向き
		fromX = fromPos.x + fromPos.width
		fromY = fromCenterY + fromYOffset
		toX = toPos.x
		toY = toCenterY + toYOffset
	} else {
		// 左向き
		fromX = fromPos.x
		fromY = fromCenterY + fromYOffset
		toX = toPos.x + toPos.width
		toY = toCenterY + toYOffset
	}
	return
}

// renderEdgeImproved は改良されたエッジ描画
func (r *ClassRenderer) renderEdgeImproved(c *canvas.Canvas, edge class.Edge, x1, y1, x2, y2 int, nodePositions map[string]struct{ x, y, width, height int }) {
	opts := []canvas.Option{canvas.Stroke(canvas.ColorEdge)}
	if edge.LineStyle == class.LineStyleDashed {
		opts = append(opts, canvas.Dashed())
	}

	// 障害物リストを構築（始点・終点ノード以外）
	var obstacles []rect
	for nodeID, pos := range nodePositions {
		if nodeID == edge.From || nodeID == edge.To {
			continue
		}
		obstacles = append(obstacles, rect{pos.x, pos.y, pos.width, pos.height})
	}

	// ウェイポイントを計算
	waypoints := r.calculateWaypoints(x1, y1, x2, y2, obstacles)

	// パスを描画
	r.renderPath(c, waypoints, opts)

	// 矢印を描画
	if len(waypoints) >= 2 {
		if edge.Decoration == class.DecorationFilledDiamond || edge.Decoration == class.DecorationEmptyDiamond {
			// ダイヤモンド装飾はFrom（始点）ノード側に描画
			r.drawArrowHead(c, edge, waypoints[0].x, waypoints[0].y, waypoints[1].x, waypoints[1].y)
		} else {
			// 三角形・矢印装飾はTo（終点）ノード側に描画（最後のセグメントの方向で）
			lastIdx := len(waypoints) - 1
			r.drawArrowHead(c, edge, waypoints[lastIdx-1].x, waypoints[lastIdx-1].y, waypoints[lastIdx].x, waypoints[lastIdx].y)
		}
	}

	// ラベル描画
	if edge.Label != "" {
		midX := (x1 + x2) / 2
		midY := (y1 + y2) / 2
		c.Text(midX, midY-5, edge.Label,
			canvas.TextAnchor("middle"),
			canvas.Fill(canvas.ColorEdgeLabel),
		)
	}
}

// calculateWaypoints は障害物を回避するウェイポイントを計算（常に直交ルーティング）
func (r *ClassRenderer) calculateWaypoints(x1, y1, x2, y2 int, obstacles []rect) []point {
	start := point{x1, y1}
	end := point{x2, y2}

	// 完全に水平または垂直の場合のみ直線を許可
	if x1 == x2 || y1 == y2 {
		if !r.pathIntersectsAnyObstacle(start, end, obstacles) {
			return []point{start, end}
		}
	}

	// 直交ルーティングが必要
	margin := 20 // 障害物からのマージン

	// ルーティング戦略を選択
	// 1. L字型（1回曲がり）
	// 2. Z字型（2回曲がり）
	// 3. U字型（障害物を迂回）

	// L字型ルーティング
	// 主方向に合わせて最終セグメントの向きを決定
	// （矢印の向きがノードへの接続方向と一致するように）
	cornerA := point{x2, y1} // 水平→垂直（最終セグメントが垂直）
	cornerB := point{x1, y2} // 垂直→水平（最終セグメントが水平）

	if abs(y2-y1) >= abs(x2-x1) {
		// 垂直が主方向：最終セグメントが垂直になるcornerAを優先
		if !r.pathIntersectsAnyObstacle(start, cornerA, obstacles) &&
			!r.pathIntersectsAnyObstacle(cornerA, end, obstacles) {
			return []point{start, cornerA, end}
		}
		if !r.pathIntersectsAnyObstacle(start, cornerB, obstacles) &&
			!r.pathIntersectsAnyObstacle(cornerB, end, obstacles) {
			return []point{start, cornerB, end}
		}
	} else {
		// 水平が主方向：最終セグメントが水平になるcornerBを優先
		if !r.pathIntersectsAnyObstacle(start, cornerB, obstacles) &&
			!r.pathIntersectsAnyObstacle(cornerB, end, obstacles) {
			return []point{start, cornerB, end}
		}
		if !r.pathIntersectsAnyObstacle(start, cornerA, obstacles) &&
			!r.pathIntersectsAnyObstacle(cornerA, end, obstacles) {
			return []point{start, cornerA, end}
		}
	}

	// Z字型（中央で曲がる）
	if abs(y2-y1) > abs(x2-x1) {
		// 垂直が主方向：垂直-水平-垂直
		midY := (y1 + y2) / 2
		mid1 := point{x1, midY}
		mid2 := point{x2, midY}
		if !r.pathIntersectsAnyObstacle(start, mid1, obstacles) &&
			!r.pathIntersectsAnyObstacle(mid1, mid2, obstacles) &&
			!r.pathIntersectsAnyObstacle(mid2, end, obstacles) {
			return []point{start, mid1, mid2, end}
		}
	} else {
		// 水平が主方向：水平-垂直-水平
		midX := (x1 + x2) / 2
		mid1 := point{midX, y1}
		mid2 := point{midX, y2}
		if !r.pathIntersectsAnyObstacle(start, mid1, obstacles) &&
			!r.pathIntersectsAnyObstacle(mid1, mid2, obstacles) &&
			!r.pathIntersectsAnyObstacle(mid2, end, obstacles) {
			return []point{start, mid1, mid2, end}
		}
	}

	// U字型ルーティング（障害物の外側を通る）
	// 全障害物のバウンディングボックスを計算
	if len(obstacles) > 0 {
		minObsX, minObsY := obstacles[0].x, obstacles[0].y
		maxObsX, maxObsY := obstacles[0].x+obstacles[0].w, obstacles[0].y+obstacles[0].h

		for _, obs := range obstacles {
			if obs.x < minObsX {
				minObsX = obs.x
			}
			if obs.y < minObsY {
				minObsY = obs.y
			}
			if obs.x+obs.w > maxObsX {
				maxObsX = obs.x + obs.w
			}
			if obs.y+obs.h > maxObsY {
				maxObsY = obs.y + obs.h
			}
		}

		// 上を通るルート
		topY := minObsY - margin
		if y1 <= topY || y2 <= topY {
			mid1 := point{x1, topY}
			mid2 := point{x2, topY}
			return []point{start, mid1, mid2, end}
		}

		// 下を通るルート
		bottomY := maxObsY + margin
		mid1 := point{x1, bottomY}
		mid2 := point{x2, bottomY}
		return []point{start, mid1, mid2, end}
	}

	// フォールバック: 単純なZ字型
	midY := (y1 + y2) / 2
	return []point{start, {x1, midY}, {x2, midY}, end}
}

// pathIntersectsAnyObstacle はパスが障害物と交差するか確認
func (r *ClassRenderer) pathIntersectsAnyObstacle(p1, p2 point, obstacles []rect) bool {
	for _, obs := range obstacles {
		if r.lineIntersectsRect(p1.x, p1.y, p2.x, p2.y, obs.x, obs.y, obs.w, obs.h) {
			return true
		}
	}
	return false
}

// renderPath はウェイポイントに沿ってパスを描画
func (r *ClassRenderer) renderPath(c *canvas.Canvas, waypoints []point, opts []canvas.Option) {
	for i := 0; i < len(waypoints)-1; i++ {
		c.Line(waypoints[i].x, waypoints[i].y, waypoints[i+1].x, waypoints[i+1].y, opts...)
	}
}

// lineIntersectsRect は直線が矩形と交差するかチェック
func (r *ClassRenderer) lineIntersectsRect(x1, y1, x2, y2, rx, ry, rw, rh int) bool {
	// 簡易的な交差判定：直線のバウンディングボックスと矩形の交差
	lminX := minInt(x1, x2)
	lmaxX := maxInt(x1, x2)
	lminY := minInt(y1, y2)
	lmaxY := maxInt(y1, y2)

	// 矩形の内側にマージンを設けて判定
	margin := 10
	rectLeft := rx + margin
	rectRight := rx + rw - margin
	rectTop := ry + margin
	rectBottom := ry + rh - margin

	// バウンディングボックスが重なっているかチェック
	if lmaxX < rectLeft || lminX > rectRight || lmaxY < rectTop || lminY > rectBottom {
		return false
	}

	// 直線が矩形の内部を通過するかの簡易チェック
	// 直線の中点が矩形内にあるか
	midX := (x1 + x2) / 2
	midY := (y1 + y2) / 2
	if midX > rectLeft && midX < rectRight && midY > rectTop && midY < rectBottom {
		return true
	}

	return false
}

// renderOrthogonalEdge は直交エッジを描画
func (r *ClassRenderer) renderOrthogonalEdge(c *canvas.Canvas, edge class.Edge, x1, y1, x2, y2 int, opts []canvas.Option) {
	// 垂直が主な方向か判定
	if abs(y2-y1) > abs(x2-x1) {
		// Z字型ルーティング（垂直-水平-垂直）
		midY := (y1 + y2) / 2
		c.Line(x1, y1, x1, midY, opts...)
		c.Line(x1, midY, x2, midY, opts...)
		c.Line(x2, midY, x2, y2, opts...)
		r.drawArrowHead(c, edge, x2, midY, x2, y2)
	} else {
		// Z字型ルーティング（水平-垂直-水平）
		midX := (x1 + x2) / 2
		c.Line(x1, y1, midX, y1, opts...)
		c.Line(midX, y1, midX, y2, opts...)
		c.Line(midX, y2, x2, y2, opts...)
		r.drawArrowHead(c, edge, midX, y2, x2, y2)
	}
}

// renderEdgeWithOffset はオフセット付きでエッジを描画する
func (r *ClassRenderer) renderEdgeWithOffset(c *canvas.Canvas, edge class.Edge, x1, y1, x2, y2, offset int) {
	// オフセットを適用（垂直エッジの場合はX方向、水平の場合はY方向）
	if abs(y2-y1) > abs(x2-x1) {
		// 主に垂直
		x1 += offset
		x2 += offset
	} else {
		// 主に水平
		y1 += offset
		y2 += offset
	}

	opts := []canvas.Option{canvas.Stroke(canvas.ColorEdge)}
	if edge.LineStyle == class.LineStyleDashed {
		opts = append(opts, canvas.Dashed())
	}

	// 直交ルーティング（L字型パス）
	if abs(x2-x1) > 20 && abs(y2-y1) > 20 {
		// 中間点でL字に曲げる
		midY := (y1 + y2) / 2
		c.Line(x1, y1, x1, midY, opts...)
		c.Line(x1, midY, x2, midY, opts...)
		c.Line(x2, midY, x2, y2, opts...)

		// 矢印の先端
		r.drawArrowHead(c, edge, x2, midY, x2, y2)
	} else {
		// 直線
		c.Line(x1, y1, x2, y2, opts...)
		r.drawArrowHead(c, edge, x1, y1, x2, y2)
	}
}

// renderVerticalEdge は継承・実装エッジを垂直接続で描画する
// 子のtop → 親のbottomを常に垂直に接続し、最後のセグメントが辺に垂直になるようにする
func (r *ClassRenderer) renderVerticalEdge(c *canvas.Canvas, edge class.Edge,
	fromPos, toPos struct{ x, y, width, height int },
	outIdx, outTotal, inIdx, inTotal int) {

	fromCenterX := fromPos.x + fromPos.width/2
	toCenterX := toPos.x + toPos.width/2

	// 接続点を分散配置（中央を基準に均等分布）
	toSpread := int(float64(toPos.width) * 0.8)
	toOffset := 0
	if inTotal > 1 {
		toOffset = (2*inIdx - (inTotal - 1)) * toSpread / (2 * (inTotal - 1))
	}

	// 子のtop、親のbottomに接続
	fromX := fromCenterX
	fromY := fromPos.y
	toX := toCenterX + toOffset
	toY := toPos.y + toPos.height

	opts := []canvas.Option{canvas.Stroke(canvas.ColorEdge)}
	if edge.LineStyle == class.LineStyleDashed {
		opts = append(opts, canvas.Dashed())
	}

	if fromX == toX {
		// 真っ直ぐ垂直: 直線で接続
		c.Line(fromX, fromY, toX, toY, opts...)
		r.drawArrowHead(c, edge, fromX, fromY, toX, toY)
	} else {
		// Z字型ルーティング: 垂直→水平→垂直
		// 最後のセグメントが常に垂直になり、親のbottom辺に垂直に接続する
		midY := (fromY + toY) / 2
		c.Line(fromX, fromY, fromX, midY, opts...)
		c.Line(fromX, midY, toX, midY, opts...)
		c.Line(toX, midY, toX, toY, opts...)
		r.drawArrowHead(c, edge, toX, midY, toX, toY)
	}

	// ラベル描画
	if edge.Label != "" {
		midX := (fromX + toX) / 2
		midY := (fromY + toY) / 2
		c.Text(midX, midY-5, edge.Label,
			canvas.TextAnchor("middle"),
			canvas.Fill(canvas.ColorEdgeLabel),
		)
	}
}

// drawArrowHead はエッジの装飾（矢印先端）を描画
func (r *ClassRenderer) drawArrowHead(c *canvas.Canvas, edge class.Edge, fromX, fromY, toX, toY int) {
	switch edge.Decoration {
	case class.DecorationTriangle:
		c.Polygon(trianglePoints(toX, toY, fromX, fromY), canvas.Fill(canvas.ColorNodeFill), canvas.Stroke(canvas.ColorEdge))
	case class.DecorationFilledDiamond:
		c.Polygon(diamondPoints(fromX, fromY, toX, toY), canvas.Fill(canvas.ColorEdge))
	case class.DecorationEmptyDiamond:
		c.Polygon(diamondPoints(fromX, fromY, toX, toY), canvas.Fill(canvas.ColorNodeFill), canvas.Stroke(canvas.ColorEdge))
	default:
		c.Polygon(trianglePoints(toX, toY, fromX, fromY), canvas.Fill(canvas.ColorEdge))
	}
}

// calculateNodeWidth はテキスト内容に基づいてノード幅を計算する
func (r *ClassRenderer) calculateNodeWidth(node class.Node) int {
	minWidth := 120
	padding := 30 // 左右のパディング
	fontSize := 12

	maxTextWidth := 0

	// クラス名の幅
	nameWidth, _ := canvas.MeasureText(node.Name, fontSize)
	if nameWidth > maxTextWidth {
		maxTextWidth = nameWidth
	}

	// ステレオタイプの幅
	if node.Stereotype != "" {
		stereoWidth, _ := canvas.MeasureText("<<"+node.Stereotype+">>", fontSize)
		if stereoWidth > maxTextWidth {
			maxTextWidth = stereoWidth
		}
	}

	// 属性の幅
	for _, attr := range node.Attributes {
		text := visibilitySymbol(attr.Visibility) + attr.Name + ": " + attr.Type
		attrWidth, _ := canvas.MeasureText(text, fontSize)
		if attrWidth > maxTextWidth {
			maxTextWidth = attrWidth
		}
	}

	// メソッドの幅
	for _, method := range node.Methods {
		text := visibilitySymbol(class.Visibility(method.Visibility)) + r.formatMethod(method)
		methodWidth, _ := canvas.MeasureText(text, fontSize)
		if methodWidth > maxTextWidth {
			maxTextWidth = methodWidth
		}
	}

	width := maxTextWidth + padding
	if width < minWidth {
		width = minWidth
	}
	return width
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

// formatMethodMultiline は複数行のメソッドシグネチャを返す（E-001）
// パラメータが多い場合やシグネチャが長い場合に使用
func (r *ClassRenderer) formatMethodMultiline(method class.Method, maxWidth int) []string {
	singleLine := r.formatMethod(method)
	charWidth := 7 // 概算1文字幅

	// 1行に収まる場合はそのまま
	if len(singleLine)*charWidth <= maxWidth {
		return []string{singleLine}
	}

	// 複数行フォーマット
	var lines []string

	prefix := ""
	if method.Async {
		prefix = "async "
	}

	returnType := ""
	if method.ReturnType != "" {
		returnType = ": " + method.ReturnType
	}

	if len(method.Params) == 0 {
		lines = append(lines, prefix+method.Name+"()"+returnType)
		return lines
	}

	// 最初の行: メソッド名と開き括弧
	lines = append(lines, prefix+method.Name+"(")

	// パラメータを各行に
	for i, p := range method.Params {
		paramStr := "  "
		if p.Name != "" && p.Type != "" {
			paramStr += p.Name + ": " + p.Type
		} else if p.Type != "" {
			paramStr += p.Type
		} else if p.Name != "" {
			paramStr += p.Name
		}
		if i < len(method.Params)-1 {
			paramStr += ","
		}
		lines = append(lines, paramStr)
	}

	// 閉じ括弧と戻り値型
	lines = append(lines, ")"+returnType)

	return lines
}

func (r *ClassRenderer) renderEdge(c *canvas.Canvas, edge class.Edge, x1, y1, x2, y2 int) {
	opts := []canvas.Option{canvas.Stroke(canvas.ColorEdge)}
	if edge.LineStyle == class.LineStyleDashed {
		opts = append(opts, canvas.Dashed())
	}
	c.Line(x1, y1, x2, y2, opts...)

	// 矢印の先端を描画
	c.Polygon(trianglePoints(x2, y2, x1, y1), canvas.Fill(canvas.ColorEdge))

	// 装飾
	switch edge.Decoration {
	case class.DecorationTriangle:
		c.Polygon(trianglePoints(x2, y2, x1, y1), canvas.Fill(canvas.ColorNodeFill), canvas.Stroke(canvas.ColorEdge))
	case class.DecorationFilledDiamond:
		c.Polygon(diamondPoints(x1, y1, x2, y2), canvas.Fill(canvas.ColorEdge))
	case class.DecorationEmptyDiamond:
		c.Polygon(diamondPoints(x1, y1, x2, y2), canvas.Fill(canvas.ColorNodeFill), canvas.Stroke(canvas.ColorEdge))
	}
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

func trianglePoints(x, y, fromX, fromY int) string {
	// 矢印の方向を計算
	dx := float64(x - fromX)
	dy := float64(y - fromY)
	length := sqrt(dx*dx + dy*dy)
	if length == 0 {
		return fmt.Sprintf("%d,%d %d,%d %d,%d", x, y-5, x+10, y, x, y+5)
	}
	// 単位ベクトル
	ux := dx / length
	uy := dy / length
	// 垂直ベクトル
	px := -uy
	py := ux
	// 三角形の頂点
	size := 10.0
	ax := float64(x) - ux*size
	ay := float64(y) - uy*size
	p1x := int(ax + px*size/2)
	p1y := int(ay + py*size/2)
	p2x := int(ax - px*size/2)
	p2y := int(ay - py*size/2)
	return fmt.Sprintf("%d,%d %d,%d %d,%d", x, y, p1x, p1y, p2x, p2y)
}

func diamondPoints(x, y, toX, toY int) string {
	// ひし形を始点に描画
	size := 10
	dx := float64(toX - x)
	dy := float64(toY - y)
	length := sqrt(dx*dx + dy*dy)
	if length == 0 {
		return fmt.Sprintf("%d,%d %d,%d %d,%d %d,%d", x, y-size, x+size, y, x, y+size, x-size, y)
	}
	ux := dx / length
	uy := dy / length
	px := -uy
	py := ux
	// ひし形の4頂点
	cx := float64(x) + ux*float64(size)
	cy := float64(y) + uy*float64(size)
	return fmt.Sprintf("%d,%d %.0f,%.0f %.0f,%.0f %.0f,%.0f",
		x, y,
		cx+px*float64(size)/2, cy+py*float64(size)/2,
		float64(x)+ux*float64(size*2), float64(y)+uy*float64(size*2),
		cx-px*float64(size)/2, cy-py*float64(size)/2)
}
