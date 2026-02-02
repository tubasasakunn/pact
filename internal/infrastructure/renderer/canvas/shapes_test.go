package canvas

import (
	"strings"
	"testing"
)

// =============================================================================
// RS001-RS005: Shapes Tests
// =============================================================================

// RS001: ひし形
func TestShapes_Diamond(t *testing.T) {
	c := New()
	c.Diamond(50, 50, 60, 40)
	svg := c.String()

	if !strings.Contains(svg, "<polygon") {
		t.Error("expected polygon element for diamond")
	}
	if !strings.Contains(svg, "points=") {
		t.Error("expected points attribute")
	}
}

// RS002: 矢印
func TestShapes_Arrow(t *testing.T) {
	c := New()
	c.Arrow(0, 0, 100, 0)
	svg := c.String()

	if !strings.Contains(svg, "<line") {
		t.Error("expected line element")
	}
	if !strings.Contains(svg, "<polygon") {
		t.Error("expected polygon element for arrowhead")
	}
}

// RS003: 端子形状
func TestShapes_Stadium(t *testing.T) {
	c := New()
	c.Stadium(10, 10, 100, 40)
	svg := c.String()

	if !strings.Contains(svg, "<rect") {
		t.Error("expected rect element")
	}
	if !strings.Contains(svg, "rx=") {
		t.Error("expected rx attribute for rounded corners")
	}
	if !strings.Contains(svg, "ry=") {
		t.Error("expected ry attribute for rounded corners")
	}
}

// RS004: 円柱（DB形状）
func TestShapes_Cylinder(t *testing.T) {
	c := New()
	c.Cylinder(10, 10, 60, 80)
	svg := c.String()

	// 円柱は楕円2つと矩形1つで構成
	ellipseCount := strings.Count(svg, "<ellipse")
	if ellipseCount != 2 {
		t.Errorf("expected 2 ellipse elements, got %d", ellipseCount)
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("expected rect element for cylinder body")
	}
}

// RS005: 平行四辺形（IO形状）
func TestShapes_Parallelogram(t *testing.T) {
	c := New()
	c.Parallelogram(10, 10, 100, 50, 15)
	svg := c.String()

	if !strings.Contains(svg, "<polygon") {
		t.Error("expected polygon element for parallelogram")
	}
	if !strings.Contains(svg, "points=") {
		t.Error("expected points attribute")
	}
}
