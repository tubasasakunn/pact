package svg

import (
	"pact/internal/domain/diagram/class"
	"pact/internal/infrastructure/renderer/canvas"
)

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
