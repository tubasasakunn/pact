package svg

import (
	"pact/internal/domain/diagram/state"
	"pact/internal/infrastructure/renderer/canvas"
)

// renderCompoundState は複合状態を描画する
func (r *StateRenderer) renderCompoundState(c *canvas.Canvas, s state.State, x, y int) {
	// 子状態の数に基づいてサイズを計算
	childCount := len(s.Children)
	cols := 2
	rows := (childCount + cols - 1) / cols
	if rows < 1 {
		rows = 1
	}

	width := cols*90 + 40
	height := 30 + rows*50 + 20 // ヘッダー + 子状態 + マージン

	// 複合状態の外枠
	c.RoundRect(x-width/2, y-20, width, height, 10, 10,
		canvas.Fill(canvas.ColorHeaderFill),
		canvas.Stroke(canvas.ColorNodeStroke),
		canvas.StrokeWidth(2),
		canvas.Filter("drop-shadow"),
	)

	// 状態名（上部）
	c.Text(x, y, s.Name,
		canvas.TextAnchor("middle"),
		canvas.Fill(canvas.ColorNodeText),
		canvas.FontWeight("bold"),
	)
	c.Line(x-width/2, y+10, x+width/2, y+10, canvas.Stroke(canvas.ColorSectionLine))

	// 子状態を描画
	childX := x - width/2 + 50
	childY := y + 30
	for i, child := range s.Children {
		col := i % cols
		row := i / cols
		cx := childX + col*90
		cy := childY + row*50

		// 子状態を通常の状態として描画
		c.RoundRect(cx-35, cy-15, 70, 30, 8, 8,
			canvas.Fill(canvas.ColorNodeFill),
			canvas.Stroke(canvas.ColorNodeStroke),
		)
		c.Text(cx, cy+5, child.Name,
			canvas.TextAnchor("middle"),
			canvas.Fill(canvas.ColorNodeText),
		)
	}
}

// renderParallelState は並行状態を描画する
func (r *StateRenderer) renderParallelState(c *canvas.Canvas, s state.State, x, y int) {
	regionCount := len(s.Regions)
	if regionCount == 0 {
		return
	}

	regionWidth := 100
	width := regionCount*regionWidth + 20
	height := 100

	// 並行状態の外枠
	c.RoundRect(x-width/2, y-20, width, height, 10, 10,
		canvas.Fill(canvas.ColorHeaderFill),
		canvas.Stroke(canvas.ColorNodeStroke),
		canvas.StrokeWidth(2),
		canvas.Filter("drop-shadow"),
	)

	// 状態名（上部）
	c.Text(x, y, s.Name,
		canvas.TextAnchor("middle"),
		canvas.Fill(canvas.ColorNodeText),
		canvas.FontWeight("bold"),
	)
	c.Line(x-width/2, y+10, x+width/2, y+10, canvas.Stroke(canvas.ColorSectionLine))

	// 各リージョンを描画
	for i, region := range s.Regions {
		rx := x - width/2 + 10 + i*regionWidth + regionWidth/2
		ry := y + 30

		// リージョン名
		c.Text(rx, ry, region.Name,
			canvas.TextAnchor("middle"),
			canvas.Fill(canvas.ColorEdgeLabel),
		)

		// リージョン内の状態を簡略表示
		stateY := ry + 25
		for j, child := range region.States {
			if j >= 2 { // 最大2つまで表示
				c.Text(rx, stateY, "...", canvas.TextAnchor("middle"), canvas.Fill(canvas.ColorEdgeLabel))
				break
			}
			c.RoundRect(rx-30, stateY-10, 60, 20, 5, 5,
				canvas.Fill(canvas.ColorNodeFill),
				canvas.Stroke(canvas.ColorNodeStroke),
			)
			c.Text(rx, stateY+5, child.Name,
				canvas.TextAnchor("middle"),
				canvas.Fill(canvas.ColorNodeText),
			)
			stateY += 25
		}

		// リージョン間の区切り線
		if i < regionCount-1 {
			lineX := x - width/2 + 10 + (i+1)*regionWidth
			c.Line(lineX, y+10, lineX, y-20+height, canvas.Stroke(canvas.ColorNodeStroke), canvas.Dashed())
		}
	}
}
