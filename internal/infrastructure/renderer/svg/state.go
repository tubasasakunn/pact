package svg

import (
	"fmt"
	"io"

	"pact/internal/domain/diagram/state"
	"pact/internal/infrastructure/renderer/canvas"
)

// StateRenderer は状態図をSVGにレンダリングする
type StateRenderer struct{}

// stateRect は状態のバウンディングボックスを表す
type stateRect struct {
	x, y, w, h int
}

// NewStateRenderer は新しいStateRendererを作成する
func NewStateRenderer() *StateRenderer {
	return &StateRenderer{}
}

// Render は状態図をSVGにレンダリングする
func (r *StateRenderer) Render(diagram *state.Diagram, w io.Writer) error {
	c := canvas.New()

	// テンプレートレジストリを適用
	registry := canvas.NewBuiltinRegistry()
	registry.ApplyTo(c)

	// 状態をタイプ別に分類
	var initialState *state.State
	var finalStates []*state.State
	var normalStates []*state.State

	for i := range diagram.States {
		s := &diagram.States[i]
		switch s.Type {
		case state.StateTypeInitial:
			initialState = s
		case state.StateTypeFinal:
			finalStates = append(finalStates, s)
		default:
			normalStates = append(normalStates, s)
		}
	}

	// 各状態の幅を事前計算
	stateWidths := make(map[string]int)
	for _, s := range normalStates {
		stateWidths[s.ID] = r.calculateStateWidth(*s)
	}

	// 列ごとの最大幅を計算
	cols := 3
	rows := (len(normalStates) + cols - 1) / cols
	colMaxWidths := make([]int, cols)
	for i, s := range normalStates {
		col := i % cols
		w := stateWidths[s.ID]
		if w > colMaxWidths[col] {
			colMaxWidths[col] = w
		}
	}

	// 列の開始X座標を計算（動的間隔）
	margin := 50 // 左マージン
	colGap := 40 // 列間のギャップ
	colCenters := make([]int, cols)
	x := margin
	for col := 0; col < cols; col++ {
		colWidth := colMaxWidths[col]
		if colWidth < 100 {
			colWidth = 100
		}
		colCenters[col] = x + colWidth/2
		x += colWidth + colGap
	}

	// キャンバスサイズを計算
	totalWidth := x + margin
	height := 100 + rows*120 + 100
	if len(finalStates) > 0 {
		height += 80
	}
	if totalWidth < 800 {
		totalWidth = 800
	}
	if height < 600 {
		height = 600
	}
	c.SetSize(totalWidth, height)

	// 状態の位置とサイズを記録
	statePositions := make(map[string]struct{ x, y int })
	stateSizes := make(map[string]struct{ w, h int })

	// 初期状態を描画（テンプレート使用）
	if initialState != nil {
		ix := colCenters[0]
		iy := 50
		statePositions[initialState.ID] = struct{ x, y int }{ix, iy}
		stateSizes[initialState.ID] = struct{ w, h int }{20, 20}
		c.UseTemplate("initial-state", ix-10, iy-10, 20, 20)
	}

	// 通常状態を描画
	startY := 120
	for i, s := range normalStates {
		col := i % cols
		row := i / cols
		sx := colCenters[col]
		sy := startY + row*120
		statePositions[s.ID] = struct{ x, y int }{sx, sy}
		sw := stateWidths[s.ID]
		sh := r.calculateStateHeight(*s)
		stateSizes[s.ID] = struct{ w, h int }{sw, sh}
		r.renderState(c, *s, sx, sy)
	}

	// 終了状態を描画（テンプレート使用）
	finalY := startY + rows*120 + 40
	for i, s := range finalStates {
		fx := colCenters[0] + i*150
		statePositions[s.ID] = struct{ x, y int }{fx, finalY}
		stateSizes[s.ID] = struct{ w, h int }{24, 24}
		c.UseTemplate("final-state", fx-12, finalY-12, 24, 24)
	}

	// 遷移を描画（直交ルーティング）
	labelOffset := make(map[string]int)
	// ノードのバウンディングボックスリストを作成
	var nodeBounds []stateRect
	for id, pos := range statePositions {
		size := stateSizes[id]
		nodeBounds = append(nodeBounds, stateRect{
			x: pos.x - size.w/2,
			y: pos.y - size.h/2,
			w: size.w,
			h: size.h,
		})
	}

	for _, t := range diagram.Transitions {
		fromPos, fromOk := statePositions[t.From]
		toPos, toOk := statePositions[t.To]
		if fromOk && toOk {
			fromSize := stateSizes[t.From]
			toSize := stateSizes[t.To]
			key := fmt.Sprintf("%d,%d-%d,%d", fromPos.x, fromPos.y, toPos.x, toPos.y)
			offset := labelOffset[key]
			labelOffset[key] = offset + 15
			r.renderOrthogonalTransition(c, t, fromPos.x, fromPos.y, fromSize.w, fromSize.h,
				toPos.x, toPos.y, toSize.w, toSize.h, offset, nodeBounds)
		}
	}

	// ノートをレンダリング
	if len(diagram.Notes) > 0 {
		renderNotes(c, diagram.Notes, statePositions)
	}

	_, err := c.WriteTo(w)
	return err
}

// calculateStateHeight は状態ボックスの高さを計算する
func (r *StateRenderer) calculateStateHeight(s state.State) int {
	hasActions := len(s.Entry) > 0 || len(s.Exit) > 0
	if hasActions {
		return 60 + len(s.Entry)*15 + len(s.Exit)*15
	}
	return 40
}

// calculateStateWidth は状態ボックスの幅を計算する
func (r *StateRenderer) calculateStateWidth(s state.State) int {
	minWidth := 80
	padding := 20
	fontSize := 12

	maxWidth := 0

	// 状態名の幅
	nameWidth, _ := canvas.MeasureText(s.Name, fontSize)
	if nameWidth > maxWidth {
		maxWidth = nameWidth
	}

	// Entry/Exitアクションの幅
	for _, entry := range s.Entry {
		actionWidth, _ := canvas.MeasureText("entry/ "+entry, fontSize)
		if actionWidth > maxWidth {
			maxWidth = actionWidth
		}
	}
	for _, exit := range s.Exit {
		actionWidth, _ := canvas.MeasureText("exit/ "+exit, fontSize)
		if actionWidth > maxWidth {
			maxWidth = actionWidth
		}
	}

	width := maxWidth + padding
	if width < minWidth {
		width = minWidth
	}
	return width
}

func (r *StateRenderer) renderState(c *canvas.Canvas, s state.State, x, y int) {
	// 複合状態の場合
	if s.Type == state.StateTypeCompound && len(s.Children) > 0 {
		r.renderCompoundState(c, s, x, y)
		return
	}

	// 並行状態の場合
	if s.Type == state.StateTypeParallel && len(s.Regions) > 0 {
		r.renderParallelState(c, s, x, y)
		return
	}

	// 通常の状態
	width := r.calculateStateWidth(s)
	hasActions := len(s.Entry) > 0 || len(s.Exit) > 0
	height := 40
	if hasActions {
		height = 60 + len(s.Entry)*15 + len(s.Exit)*15
	}

	c.RoundRect(x-width/2, y-20, width, height, 10, 10,
		canvas.Fill(canvas.ColorNodeFill),
		canvas.Stroke(canvas.ColorNodeStroke),
		canvas.StrokeWidth(2),
		canvas.Filter("drop-shadow"),
	)
	c.Text(x, y+5, s.Name,
		canvas.TextAnchor("middle"),
		canvas.Fill(canvas.ColorNodeText),
		canvas.FontWeight("bold"),
	)

	// Entry/Exitアクションを描画
	if hasActions {
		c.Line(x-width/2, y+15, x+width/2, y+15, canvas.Stroke(canvas.ColorSectionLine))
		actionY := y + 30
		for _, entry := range s.Entry {
			c.Text(x-width/2+5, actionY, "entry/ "+entry, canvas.Fill(canvas.ColorNodeText))
			actionY += 15
		}
		for _, exit := range s.Exit {
			c.Text(x-width/2+5, actionY, "exit/ "+exit, canvas.Fill(canvas.ColorNodeText))
			actionY += 15
		}
	}
}
