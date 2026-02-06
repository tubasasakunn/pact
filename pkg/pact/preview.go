package pact

import (
	"fmt"
	"os"
	"path/filepath"

	"pact/internal/infrastructure/renderer/canvas"
)

// PatternPreviewConfig holds configuration for pattern preview generation.
type PatternPreviewConfig struct {
	OutputDir string
}

// GeneratePatternPreviews generates SVG preview files for all pattern templates.
func GeneratePatternPreviews(cfg PatternPreviewConfig) error {
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	registry := canvas.NewPatternRegistry()
	decorations := canvas.NewDecorationRegistry()

	patterns := allPatternTypes()

	for _, pt := range patterns {
		layout := registry.Get(pt)
		if layout == nil {
			continue
		}

		svg := renderPatternPreview(layout, decorations)
		filename := filepath.Join(cfg.OutputDir, string(pt)+".svg")

		if err := os.WriteFile(filename, []byte(svg), 0644); err != nil {
			return fmt.Errorf("write %s: %w", filename, err)
		}
	}

	generatePreviewIndex(cfg.OutputDir, patterns)
	return nil
}

func allPatternTypes() []canvas.PatternType {
	return []canvas.PatternType{
		// Class patterns
		canvas.PatternInheritanceTree2,
		canvas.PatternInheritanceTree3,
		canvas.PatternInheritanceTree4,
		canvas.PatternInterfaceImpl2,
		canvas.PatternInterfaceImpl3,
		canvas.PatternInterfaceImpl4,
		canvas.PatternComposition2,
		canvas.PatternComposition3,
		canvas.PatternComposition4,
		canvas.PatternDiamond,
		canvas.PatternLayered3x2,
		canvas.PatternLayered3x3,
		// State patterns
		canvas.PatternLinearStates2,
		canvas.PatternLinearStates3,
		canvas.PatternLinearStates4,
		canvas.PatternBinaryChoice,
		canvas.PatternStateLoop,
		canvas.PatternStarTopology,
		// Flow patterns
		canvas.PatternIfElse,
		canvas.PatternIfElseIfElse,
		canvas.PatternWhileLoop,
		canvas.PatternSequential3,
		canvas.PatternSequential4,
		// Sequence patterns
		canvas.PatternRequestResponse,
		canvas.PatternCallback,
		canvas.PatternChain3,
		canvas.PatternChain4,
		canvas.PatternFanOut,
	}
}

func renderPatternPreview(layout *canvas.PatternLayout, decorations *canvas.DecorationRegistry) string {
	c := canvas.New()
	c.SetSize(layout.MinWidth, layout.MinHeight)

	decorations.ApplyToCanvas(c)

	c.Rect(0, 0, layout.MinWidth, layout.MinHeight,
		canvas.Fill("#f8fafc"), canvas.Stroke("#e2e8f0"), canvas.StrokeWidth(1))

	c.Text(layout.MinWidth/2, 25, layout.Name,
		canvas.TextAnchor("middle"), canvas.Fill("#1a202c"), canvas.FontWeight("bold"), canvas.FontSize(16))

	for _, pos := range layout.Positions {
		x := int(pos.X * float64(layout.MinWidth))
		y := int(pos.Y * float64(layout.MinHeight))
		w := int(pos.Width * float64(layout.MinWidth))
		h := int(pos.Height * float64(layout.MinHeight))

		x -= w / 2
		y -= h / 2

		c.RoundRect(x, y, w, h, 6, 6,
			canvas.Fill(canvas.ColorNodeFill),
			canvas.Stroke(canvas.ColorNodeStroke),
			canvas.StrokeWidth(2),
			canvas.Filter("drop-shadow"))

		c.Text(x+w/2, y+h/2+5, pos.ID,
			canvas.TextAnchor("middle"), canvas.Fill(canvas.ColorNodeText), canvas.FontSize(11))
	}

	for _, edge := range layout.Edges {
		drawPreviewEdge(c, layout, edge)
	}

	for _, dec := range layout.Decorators {
		x := int(dec.Bounds.X * float64(layout.MinWidth))
		y := int(dec.Bounds.Y * float64(layout.MinHeight))
		w := int(dec.Bounds.Width * float64(layout.MinWidth))
		h := int(dec.Bounds.Height * float64(layout.MinHeight))

		switch dec.Type {
		case "divider":
			stroke := dec.Style["stroke"]
			if stroke == "" {
				stroke = canvas.ColorSectionLine
			}
			dasharray := dec.Style["stroke-dasharray"]
			if dasharray != "" {
				c.Line(x, y, x+w, y+h, canvas.Stroke(stroke), canvas.StrokeDasharray(dasharray))
			} else {
				c.Line(x, y, x+w, y+h, canvas.Stroke(stroke))
			}
		case "label":
			c.Text(x, y, dec.Style["text"],
				canvas.Fill(canvas.ColorEdgeLabel), canvas.FontSize(10))
		}
	}

	return c.String()
}

func drawPreviewEdge(c *canvas.Canvas, layout *canvas.PatternLayout, edge canvas.EdgePath) {
	var fromPos, toPos *canvas.LayoutPosition
	for i := range layout.Positions {
		if layout.Positions[i].ID == edge.FromID {
			fromPos = &layout.Positions[i]
		}
		if layout.Positions[i].ID == edge.ToID {
			toPos = &layout.Positions[i]
		}
	}
	if fromPos == nil || toPos == nil {
		return
	}

	fromX := int(fromPos.X * float64(layout.MinWidth))
	fromY := int(fromPos.Y * float64(layout.MinHeight))
	toX := int(toPos.X * float64(layout.MinWidth))
	toY := int(toPos.Y * float64(layout.MinHeight))

	fromW := int(fromPos.Width * float64(layout.MinWidth))
	fromH := int(fromPos.Height * float64(layout.MinHeight))
	toW := int(toPos.Width * float64(layout.MinWidth))
	toH := int(toPos.Height * float64(layout.MinHeight))

	points := []canvas.Point{}
	for _, wp := range edge.Waypoints {
		points = append(points, canvas.Point{
			X: wp.X * float64(layout.MinWidth),
			Y: wp.Y * float64(layout.MinHeight),
		})
	}

	if len(points) == 0 {
		if fromY == toY {
			startX := fromX + fromW/2
			endX := toX - toW/2
			if fromX > toX {
				startX = fromX - fromW/2
				endX = toX + toW/2
			}
			c.Line(startX, fromY, endX, toY,
				canvas.Stroke(canvas.ColorEdge), canvas.StrokeWidth(2))
		} else if fromX == toX {
			startY := fromY + fromH/2
			endY := toY - toH/2
			if fromY > toY {
				startY = fromY - fromH/2
				endY = toY + toH/2
			}
			c.Line(fromX, startY, toX, endY,
				canvas.Stroke(canvas.ColorEdge), canvas.StrokeWidth(2))
		} else {
			midY := (fromY + toY) / 2
			c.Line(fromX, fromY+fromH/2, fromX, midY,
				canvas.Stroke(canvas.ColorEdge), canvas.StrokeWidth(2))
			c.Line(fromX, midY, toX, midY,
				canvas.Stroke(canvas.ColorEdge), canvas.StrokeWidth(2))
			c.Line(toX, midY, toX, toY-toH/2,
				canvas.Stroke(canvas.ColorEdge), canvas.StrokeWidth(2))
		}
	} else {
		prevX := fromX
		prevY := fromY + fromH/2
		if len(points) > 0 && points[0].Y < float64(fromY) {
			prevY = fromY - fromH/2
		}

		for _, p := range points {
			c.Line(prevX, prevY, int(p.X), int(p.Y),
				canvas.Stroke(canvas.ColorEdge), canvas.StrokeWidth(2))
			prevX = int(p.X)
			prevY = int(p.Y)
		}

		endY := toY - toH/2
		if prevY > toY {
			endY = toY + toH/2
		}
		c.Line(prevX, prevY, toX, endY,
			canvas.Stroke(canvas.ColorEdge), canvas.StrokeWidth(2))
	}

	// Draw arrow at end
	size := 8
	ay := toY - toH/2
	arrowPoints := fmt.Sprintf("%d,%d %d,%d %d,%d",
		toX, ay, toX-size/2, ay-size, toX+size/2, ay-size)
	c.Polygon(arrowPoints, canvas.Fill(canvas.ColorEdge))

	if edge.LabelPos.X > 0 || edge.LabelPos.Y > 0 {
		labelX := int(edge.LabelPos.X * float64(layout.MinWidth))
		labelY := int(edge.LabelPos.Y * float64(layout.MinHeight))
		c.Text(labelX, labelY, edge.FromID+"→"+edge.ToID,
			canvas.TextAnchor("middle"), canvas.Fill(canvas.ColorEdgeLabel), canvas.FontSize(9))
	}
}

func generatePreviewIndex(outDir string, patterns []canvas.PatternType) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Pattern Templates Preview</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            margin: 20px;
            background: #f0f4f8;
        }
        h1 { color: #1a202c; }
        h2 { color: #2d3748; margin-top: 40px; border-bottom: 2px solid #e2e8f0; padding-bottom: 10px; }
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
            gap: 20px;
            margin-top: 20px;
        }
        .card {
            background: white;
            border-radius: 12px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .card img {
            width: 100%;
            height: auto;
            display: block;
        }
        .card-title {
            padding: 12px 16px;
            font-weight: 600;
            color: #2d3748;
            background: #f7fafc;
            border-top: 1px solid #e2e8f0;
        }
    </style>
</head>
<body>
    <h1>Pattern Templates Preview</h1>
`

	sections := []struct {
		title    string
		patterns []canvas.PatternType
	}{
		{"Class Diagram Patterns", []canvas.PatternType{
			canvas.PatternInheritanceTree2, canvas.PatternInheritanceTree3, canvas.PatternInheritanceTree4,
			canvas.PatternInterfaceImpl2, canvas.PatternInterfaceImpl3, canvas.PatternInterfaceImpl4,
			canvas.PatternComposition2, canvas.PatternComposition3, canvas.PatternComposition4,
			canvas.PatternDiamond, canvas.PatternLayered3x2, canvas.PatternLayered3x3,
		}},
		{"State Diagram Patterns", []canvas.PatternType{
			canvas.PatternLinearStates2, canvas.PatternLinearStates3, canvas.PatternLinearStates4,
			canvas.PatternBinaryChoice, canvas.PatternStateLoop, canvas.PatternStarTopology,
		}},
		{"Flow Diagram Patterns", []canvas.PatternType{
			canvas.PatternIfElse, canvas.PatternIfElseIfElse, canvas.PatternWhileLoop,
			canvas.PatternSequential3, canvas.PatternSequential4,
		}},
		{"Sequence Diagram Patterns", []canvas.PatternType{
			canvas.PatternRequestResponse, canvas.PatternCallback,
			canvas.PatternChain3, canvas.PatternChain4, canvas.PatternFanOut,
		}},
	}

	for _, sec := range sections {
		html += fmt.Sprintf("\n    <h2>%s</h2>\n    <div class=\"grid\">\n", sec.title)
		for _, p := range sec.patterns {
			html += fmt.Sprintf(`        <div class="card">
            <img src="%s.svg" alt="%s">
            <div class="card-title">%s</div>
        </div>
`, p, p, p)
		}
		html += "    </div>\n"
	}

	html += `</body>
</html>
`
	_ = os.WriteFile(filepath.Join(outDir, "index.html"), []byte(html), 0644)
}
