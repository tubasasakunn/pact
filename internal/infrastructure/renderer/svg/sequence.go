package svg

import (
	"io"

	"pact/internal/domain/diagram/sequence"
	"pact/internal/infrastructure/renderer/canvas"
)

// SequenceRenderer はシーケンス図をSVGにレンダリングする
type SequenceRenderer struct{}

// NewSequenceRenderer は新しいSequenceRendererを作成する
func NewSequenceRenderer() *SequenceRenderer {
	return &SequenceRenderer{}
}

// Render はシーケンス図をSVGにレンダリングする
func (r *SequenceRenderer) Render(diagram *sequence.Diagram, w io.Writer) error {
	c := canvas.New()

	// テンプレートレジストリを適用
	registry := canvas.NewBuiltinRegistry()
	registry.ApplyTo(c)

	// 各参加者のボックス幅を計算
	participantWidths := make(map[string]int)
	minWidth := 80
	padding := 20
	fontSize := 12

	for _, p := range diagram.Participants {
		textWidth, _ := canvas.MeasureText(p.Name, fontSize)
		width := textWidth + padding
		if width < minWidth {
			width = minWidth
		}
		participantWidths[p.ID] = width
	}

	// 参加者の位置を計算（動的間隔）
	participantX := make(map[string]int)
	margin := 30 // 参加者間のマージン
	x := 50

	for _, p := range diagram.Participants {
		w := participantWidths[p.ID]
		participantX[p.ID] = x + w/2 // 中心位置
		x += w + margin
	}

	// キャンバス幅を計算
	totalWidth := x + 50
	if totalWidth < 800 {
		totalWidth = 800
	}
	c.SetSize(totalWidth, 600)

	// 参加者をレンダリング
	for _, p := range diagram.Participants {
		px := participantX[p.ID]
		pw := participantWidths[p.ID]
		r.renderParticipantWithWidth(c, p, px, 50, pw)
	}

	// メッセージをレンダリング
	messageY := 120
	frameWidth := totalWidth - 100 // 左右マージンを引いた幅
	if frameWidth < 700 {
		frameWidth = 700
	}
	r.renderEvents(c, diagram.Events, participantX, &messageY, frameWidth)

	// ノートをレンダリング
	if len(diagram.Notes) > 0 {
		simplePositions := make(map[string]struct{ x, y int })
		for id, x := range participantX {
			simplePositions[id] = struct{ x, y int }{x, 50}
		}
		renderNotes(c, diagram.Notes, simplePositions)
	}

	_, err := c.WriteTo(w)
	return err
}

func (r *SequenceRenderer) renderEvents(c *canvas.Canvas, events []sequence.Event, participantX map[string]int, y *int, frameWidth int) {
	// アクティベーション状態を追跡
	activations := make(map[string]int) // participant -> activation start Y

	for _, event := range events {
		switch e := event.(type) {
		case *sequence.MessageEvent:
			fromX, fromOk := participantX[e.From]
			toX, toOk := participantX[e.To]
			if !fromOk || !toOk {
				continue
			}

			// メッセージの矢印を描画
			switch e.MessageType {
			case sequence.MessageTypeAsync:
				c.Line(fromX, *y, toX, *y, canvas.Stroke(canvas.ColorEdge), canvas.Dashed())
				r.drawOpenArrow(c, fromX, toX, *y)
			case sequence.MessageTypeReturn:
				c.Line(fromX, *y, toX, *y, canvas.Stroke(canvas.ColorEdge), canvas.Dashed())
				r.drawOpenArrow(c, fromX, toX, *y)
				// returnでアクティベーション終了
				if startY, ok := activations[e.From]; ok {
					r.drawActivationBar(c, fromX, startY, *y)
					delete(activations, e.From)
				}
			default: // sync
				c.Arrow(fromX, *y, toX, *y, canvas.Stroke(canvas.ColorEdge))
				// syncでターゲットをアクティベート
				if _, ok := activations[e.To]; !ok {
					activations[e.To] = *y
				}
			}

			// ラベル
			midX := (fromX + toX) / 2
			c.Text(midX, *y-5, e.Label,
				canvas.TextAnchor("middle"),
				canvas.Fill(canvas.ColorNodeText),
			)

			*y += 40

		case *sequence.FragmentEvent:
			// フラグメント（alt, loop, opt）の枠を描画
			startY := *y

			// メイン(then)部分のイベントをレンダリング
			r.renderEvents(c, e.Events, participantX, y, frameWidth)

			// alt フラグメントで AltEvents がある場合
			altSeparatorY := 0
			if e.Type == sequence.FragmentTypeAlt && len(e.AltEvents) > 0 {
				altSeparatorY = *y
				// 区切り線用に少しスペースを空ける
				*y += 20
				// else ラベル
				if e.AltLabel != "" {
					c.Text(60, *y-5, "["+e.AltLabel+"]",
						canvas.Fill(canvas.ColorEdgeLabel),
					)
				}
				// else 部分のイベントをレンダリング
				r.renderEvents(c, e.AltEvents, participantX, y, frameWidth)
			}

			// 枠を描画（参加者数に応じた幅）
			c.Rect(50, startY-10, frameWidth, *y-startY+20,
				canvas.Stroke(canvas.ColorNodeStroke), canvas.Fill("none"),
			)
			c.Text(60, startY, "["+string(e.Type)+"] "+e.Label,
				canvas.Fill(canvas.ColorEdgeLabel),
				canvas.FontWeight("bold"),
			)

			// alt の区切り線を描画
			if altSeparatorY > 0 {
				c.Line(50, altSeparatorY+5, 50+frameWidth, altSeparatorY+5,
					canvas.Stroke(canvas.ColorNodeStroke), canvas.Dashed(),
				)
			}

		case *sequence.ActivationEvent:
			x, ok := participantX[e.Participant]
			if !ok {
				continue
			}
			if e.Active {
				// アクティベーション開始
				activations[e.Participant] = *y
			} else {
				// アクティベーション終了
				if startY, ok := activations[e.Participant]; ok {
					r.drawActivationBar(c, x, startY, *y)
					delete(activations, e.Participant)
				}
			}

		case *sequence.NoteEvent:
			// 注釈（return/throw）の描画
			x, ok := participantX[e.Participant]
			if !ok {
				x = 100 // デフォルト位置
			}

			noteWidth, _ := canvas.MeasureText(e.Text, 12)
			noteWidth += 20 // パディング
			noteHeight := 25

			// 注釈の種類によって色を変える
			fillColor := canvas.ColorDefaultNote
			strokeColor := canvas.ColorNodeStroke
			if e.NoteType == sequence.NoteTypeThrow {
				fillColor = canvas.ColorThrowFill
				strokeColor = canvas.ColorThrowStroke
			} else if e.NoteType == sequence.NoteTypeReturn {
				fillColor = canvas.ColorReturnFill
				strokeColor = canvas.ColorReturnStroke
			}

			// 注釈ボックスを描画
			noteX := x + 20
			noteY := *y - 10
			c.Rect(noteX, noteY, noteWidth, noteHeight,
				canvas.Fill(fillColor), canvas.Stroke(strokeColor),
			)
			c.Text(noteX+10, *y+5, e.Text, canvas.Fill(canvas.ColorNodeText))

			*y += 30
		}
	}

	// 残っているアクティベーションを閉じる
	for participant, startY := range activations {
		if x, ok := participantX[participant]; ok {
			r.drawActivationBar(c, x, startY, *y)
		}
	}
}

// drawActivationBar はアクティベーションバーを描画する
func (r *SequenceRenderer) drawActivationBar(c *canvas.Canvas, x, startY, endY int) {
	barWidth := 10
	c.Rect(x-barWidth/2, startY, barWidth, endY-startY,
		canvas.Fill(canvas.ColorActivationFill),
		canvas.Stroke(canvas.ColorNodeStroke),
	)
}

func (r *SequenceRenderer) drawOpenArrow(c *canvas.Canvas, fromX, toX, y int) {
	if toX > fromX {
		c.Line(toX-8, y-5, toX, y, canvas.Stroke(canvas.ColorEdge))
		c.Line(toX-8, y+5, toX, y, canvas.Stroke(canvas.ColorEdge))
	} else {
		c.Line(toX+8, y-5, toX, y, canvas.Stroke(canvas.ColorEdge))
		c.Line(toX+8, y+5, toX, y, canvas.Stroke(canvas.ColorEdge))
	}
}

func (r *SequenceRenderer) renderParticipant(c *canvas.Canvas, p sequence.Participant, x, y int) {
	r.renderParticipantWithWidth(c, p, x, y, 80)
}

func (r *SequenceRenderer) renderParticipantWithWidth(c *canvas.Canvas, p sequence.Participant, x, y, width int) {
	switch p.Type {
	case sequence.ParticipantTypeActor:
		// アクターテンプレートを使用（固定プロポーション）
		c.UseTemplate("actor", x-20, y, 40, 55)
		c.Text(x, y+60, p.Name,
			canvas.TextAnchor("middle"),
			canvas.Fill(canvas.ColorNodeText),
			canvas.FontWeight("bold"),
		)
	case sequence.ParticipantTypeDatabase:
		c.Cylinder(x-width/2, y, width, 50,
			canvas.Fill(canvas.ColorNodeFill),
			canvas.Stroke(canvas.ColorNodeStroke),
		)
		c.Text(x, y+60, p.Name,
			canvas.TextAnchor("middle"),
			canvas.Fill(canvas.ColorNodeText),
			canvas.FontWeight("bold"),
		)
	default:
		c.Rect(x-width/2, y, width, 40,
			canvas.Fill(canvas.ColorNodeFill),
			canvas.Stroke(canvas.ColorNodeStroke),
			canvas.StrokeWidth(2),
			canvas.Filter("drop-shadow"),
		)
		c.Text(x, y+25, p.Name,
			canvas.TextAnchor("middle"),
			canvas.Fill(canvas.ColorNodeText),
			canvas.FontWeight("bold"),
		)
	}

	// ライフライン（破線）
	c.Line(x, y+50, x, 500, canvas.Stroke(canvas.ColorNodeStroke), canvas.Dashed())
}
