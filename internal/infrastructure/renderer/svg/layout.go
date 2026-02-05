package svg

import (
	"strings"

	"pact/internal/domain/diagram/common"
	"pact/internal/infrastructure/renderer/canvas"
)

// maxInt は2つの整数の大きい方を返す
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// minInt は2つの整数の小さい方を返す
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// abs は整数の絶対値を返す
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// sqrt はニュートン法で平方根を計算する
func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
	}
	return z
}

// rect は矩形を表す
type rect struct {
	x, y, w, h int
}

// point は点を表す
type point struct {
	x, y int
}

// nodeWithBarycenter はバリセンター値を持つノード
type nodeWithBarycenter struct {
	id         string
	barycenter float64
	hasConn    bool
}

// mergeSortBarycenter はマージソートの実装
func mergeSortBarycenter(nodes, buf []nodeWithBarycenter) {
	n := len(nodes)
	if n <= 1 {
		return
	}
	mid := n / 2
	mergeSortBarycenter(nodes[:mid], buf[:mid])
	mergeSortBarycenter(nodes[mid:], buf[mid:])

	// マージ
	copy(buf, nodes)
	i, j, k := 0, mid, 0
	for i < mid && j < n {
		if buf[i].barycenter <= buf[j].barycenter {
			nodes[k] = buf[i]
			i++
		} else {
			nodes[k] = buf[j]
			j++
		}
		k++
	}
	for i < mid {
		nodes[k] = buf[i]
		i++
		k++
	}
	for j < n {
		nodes[k] = buf[j]
		j++
		k++
	}
}

// calculateNotePosition はノートの位置を計算する
func calculateNotePosition(note common.Note, elementPositions map[string]struct{ x, y int }, noteWidth, noteHeight int) (int, int) {
	if note.AttachTo != "" {
		if pos, ok := elementPositions[note.AttachTo]; ok {
			switch note.Position {
			case common.NotePositionLeft:
				return pos.x - noteWidth - 30, pos.y
			case common.NotePositionRight:
				return pos.x + 120, pos.y
			case common.NotePositionTop:
				return pos.x, pos.y - noteHeight - 20
			case common.NotePositionBottom:
				return pos.x, pos.y + 60
			default:
				return pos.x + 120, pos.y
			}
		}
	}
	return 600, 50
}

// noteRect はノートの矩形を表す
type noteRect struct {
	x, y, w, h int
}

// rectsCollide は2つの矩形が衝突するかチェックする
func rectsCollide(a, b noteRect) bool {
	return a.x < b.x+b.w && a.x+a.w > b.x && a.y < b.y+b.h && a.y+a.h > b.y
}

// findNonCollidingPosition は衝突しない位置を探す
func findNonCollidingPosition(x, y, w, h int, occupiedRects []noteRect, elementPositions map[string]struct{ x, y int }) (int, int) {
	testRect := noteRect{x, y, w, h}

	// ノード位置もチェック
	nodeRects := make([]noteRect, 0, len(elementPositions))
	for _, pos := range elementPositions {
		nodeRects = append(nodeRects, noteRect{pos.x, pos.y, 120, 60}) // 概算のノードサイズ
	}

	// 衝突がなければそのまま返す
	hasCollision := false
	for _, rect := range occupiedRects {
		if rectsCollide(testRect, rect) {
			hasCollision = true
			break
		}
	}
	if !hasCollision {
		for _, rect := range nodeRects {
			if rectsCollide(testRect, rect) {
				hasCollision = true
				break
			}
		}
	}

	if !hasCollision {
		return x, y
	}

	// 衝突がある場合、下方向にずらして再試行
	offsets := []struct{ dx, dy int }{
		{0, 50},    // 下
		{0, -50},   // 上
		{100, 0},   // 右
		{-100, 0},  // 左
		{100, 50},  // 右下
		{-100, 50}, // 左下
		{0, 100},   // 更に下
	}

	for _, off := range offsets {
		testRect.x = x + off.dx
		testRect.y = y + off.dy

		collision := false
		for _, rect := range occupiedRects {
			if rectsCollide(testRect, rect) {
				collision = true
				break
			}
		}
		if collision {
			continue
		}
		for _, rect := range nodeRects {
			if rectsCollide(testRect, rect) {
				collision = true
				break
			}
		}

		if !collision {
			return testRect.x, testRect.y
		}
	}

	// 全てのオフセットで衝突する場合は元の位置を返す
	return x, y
}

// renderNotes はノートを描画する共通関数（衝突検出付き）
func renderNotes(c *canvas.Canvas, notes []common.Note, elementPositions map[string]struct{ x, y int }) {
	noteWidth := 100
	noteHeight := 40
	occupiedRects := make([]noteRect, 0, len(notes))

	for _, note := range notes {
		x, y := calculateNotePosition(note, elementPositions, noteWidth, noteHeight)

		// 衝突検出：他のノートやノードと重ならない位置を探す
		x, y = findNonCollidingPosition(x, y, noteWidth, noteHeight, occupiedRects, elementPositions)

		// この位置を占有済みとして記録
		occupiedRects = append(occupiedRects, noteRect{x, y, noteWidth, noteHeight})

		// 関連付け要素がある場合は接続線を描画
		if note.AttachTo != "" {
			if pos, ok := elementPositions[note.AttachTo]; ok {
				c.Line(pos.x+60, pos.y+20, x, y+noteHeight/2,
					canvas.Stroke(canvas.ColorNoteStroke), canvas.Dashed(),
				)
			}
		}

		// ノートを描画（テーマカラー）
		c.Note(x, y, noteWidth, noteHeight,
			canvas.Fill(canvas.ColorNoteFill),
			canvas.Stroke(canvas.ColorNoteStroke),
		)

		// テキストを描画（複数行対応）
		lines := strings.Split(note.Text, "\n")
		textY := y + 15
		for _, line := range lines {
			c.Text(x+5, textY, line, canvas.Fill(canvas.ColorNodeText))
			textY += 15
		}
	}
}
