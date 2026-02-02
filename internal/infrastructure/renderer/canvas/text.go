package canvas

// MeasureText はテキストの概算サイズを返す
func MeasureText(text string, fontSize int) (width, height int) {
	// 概算: 1文字あたり fontSize * 0.6 の幅
	charWidth := float64(fontSize) * 0.6
	width = int(float64(len(text)) * charWidth)
	height = fontSize
	return
}

// WrapText はテキストを指定幅で折り返す
func WrapText(text string, maxWidth, fontSize int) []string {
	charWidth := float64(fontSize) * 0.6
	charsPerLine := int(float64(maxWidth) / charWidth)

	if charsPerLine <= 0 {
		charsPerLine = 1
	}

	var lines []string
	for len(text) > 0 {
		if len(text) <= charsPerLine {
			lines = append(lines, text)
			break
		}

		// 単語境界で折り返しを試みる
		breakPoint := charsPerLine
		for i := charsPerLine; i > 0; i-- {
			if text[i] == ' ' {
				breakPoint = i
				break
			}
		}

		lines = append(lines, text[:breakPoint])
		text = text[breakPoint:]
		// 先頭の空白を削除
		for len(text) > 0 && text[0] == ' ' {
			text = text[1:]
		}
	}

	return lines
}
