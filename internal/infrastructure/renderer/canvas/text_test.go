package canvas

import "testing"

// =============================================================================
// RT001-RT002: Text Tests
// =============================================================================

// RT001: テキスト測定
func TestText_MeasureText(t *testing.T) {
	width, height := MeasureText("hello", 12)

	if width <= 0 {
		t.Error("expected positive width")
	}
	if height != 12 {
		t.Errorf("expected height 12, got %d", height)
	}

	// 長いテキストは幅が大きい
	width2, _ := MeasureText("hello world", 12)
	if width2 <= width {
		t.Error("expected longer text to have larger width")
	}
}

// RT002: テキスト折り返し
func TestText_WrapText(t *testing.T) {
	text := "This is a long text that should be wrapped"
	lines := WrapText(text, 100, 12)

	if len(lines) < 2 {
		t.Errorf("expected multiple lines, got %d", len(lines))
	}

	// 短いテキストは折り返さない
	shortText := "Hi"
	shortLines := WrapText(shortText, 100, 12)
	if len(shortLines) != 1 {
		t.Errorf("expected 1 line for short text, got %d", len(shortLines))
	}
}

// RT002b: 空テキストの折り返し
func TestText_WrapText_Empty(t *testing.T) {
	lines := WrapText("", 100, 12)
	if len(lines) != 0 {
		t.Errorf("expected 0 lines for empty text, got %d", len(lines))
	}
}

// RT002c: 幅内に収まるテキスト
func TestText_WrapText_FitsWidth(t *testing.T) {
	text := "hello"
	lines := WrapText(text, 200, 12)
	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(lines))
	}
	if lines[0] != text {
		t.Errorf("expected %q, got %q", text, lines[0])
	}
}
