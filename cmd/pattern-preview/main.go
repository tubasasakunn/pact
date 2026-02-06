// Command pattern-preview generates SVG preview files for all pattern templates.
// Usage: go run ./cmd/pattern-preview
// Output: pattern-preview/ directory with SVG files
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"pact/internal/infrastructure/renderer/canvas"
)

func main() {
	outDir := "pattern-preview"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	registry := canvas.NewPatternRegistry()
	decorations := canvas.NewDecorationRegistry()

	// Get all patterns
	patterns := []canvas.PatternType{
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

	for _, pt := range patterns {
		layout := registry.Get(pt)
		if layout == nil {
			fmt.Printf("Pattern %s not found\n", pt)
			continue
		}

		svg := generatePatternPreview(layout, decorations)
		filename := filepath.Join(outDir, string(pt)+".svg")

		if err := os.WriteFile(filename, []byte(svg), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", filename, err)
			continue
		}
		fmt.Printf("Generated: %s\n", filename)
	}

	// Generate index.html for easy viewing
	generateIndex(outDir, patterns)
	fmt.Printf("\nOpen %s/index.html in a browser to view all patterns\n", outDir)
}

func generatePatternPreview(layout *canvas.PatternLayout, decorations *canvas.DecorationRegistry) string {
	c := canvas.New()
	c.SetSize(layout.MinWidth, layout.MinHeight)

	// Apply decorations
	decorations.ApplyToCanvas(c)

	// Draw background
	c.Rect(0, 0, layout.MinWidth, layout.MinHeight,
		canvas.Fill("#f8fafc"), canvas.Stroke("#e2e8f0"), canvas.StrokeWidth(1))

	// Draw title
	c.Text(layout.MinWidth/2, 25, layout.Name,
		canvas.TextAnchor("middle"), canvas.Fill("#1a202c"), canvas.FontWeight("bold"), canvas.FontSize(16))

	// Draw nodes
	for _, pos := range layout.Positions {
		x := int(pos.X * float64(layout.MinWidth))
		y := int(pos.Y * float64(layout.MinHeight))
		w := int(pos.Width * float64(layout.MinWidth))
		h := int(pos.Height * float64(layout.MinHeight))

		// Center the node
		x -= w / 2
		y -= h / 2

		// Draw node box
		c.RoundRect(x, y, w, h, 6, 6,
			canvas.Fill(canvas.ColorNodeFill),
			canvas.Stroke(canvas.ColorNodeStroke),
			canvas.StrokeWidth(2),
			canvas.Filter("drop-shadow"))

		// Draw node label
		c.Text(x+w/2, y+h/2+5, pos.ID,
			canvas.TextAnchor("middle"), canvas.Fill(canvas.ColorNodeText), canvas.FontSize(11))
	}

	// Draw edges
	for _, edge := range layout.Edges {
		drawOrthogonalEdge(c, layout, edge)
	}

	// Draw decorators
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

func drawOrthogonalEdge(c *canvas.Canvas, layout *canvas.PatternLayout, edge canvas.EdgePath) {
	// Find source and target positions
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

	// Calculate node centers
	fromX := int(fromPos.X * float64(layout.MinWidth))
	fromY := int(fromPos.Y * float64(layout.MinHeight))
	toX := int(toPos.X * float64(layout.MinWidth))
	toY := int(toPos.Y * float64(layout.MinHeight))

	// Calculate node dimensions
	fromW := int(fromPos.Width * float64(layout.MinWidth))
	fromH := int(fromPos.Height * float64(layout.MinHeight))
	toW := int(toPos.Width * float64(layout.MinWidth))
	toH := int(toPos.Height * float64(layout.MinHeight))

	// Build path through waypoints
	points := []canvas.Point{}

	// Add waypoints
	for _, wp := range edge.Waypoints {
		points = append(points, canvas.Point{
			X: wp.X * float64(layout.MinWidth),
			Y: wp.Y * float64(layout.MinHeight),
		})
	}

	// Draw the path
	if len(points) == 0 {
		// Direct connection (should be orthogonal)
		if fromY == toY {
			// Horizontal
			startX := fromX + fromW/2
			endX := toX - toW/2
			if fromX > toX {
				startX = fromX - fromW/2
				endX = toX + toW/2
			}
			c.Line(startX, fromY, endX, toY,
				canvas.Stroke(canvas.ColorEdge), canvas.StrokeWidth(2))
		} else if fromX == toX {
			// Vertical
			startY := fromY + fromH/2
			endY := toY - toH/2
			if fromY > toY {
				startY = fromY - fromH/2
				endY = toY + toH/2
			}
			c.Line(fromX, startY, toX, endY,
				canvas.Stroke(canvas.ColorEdge), canvas.StrokeWidth(2))
		} else {
			// Orthogonal bend needed
			midY := (fromY + toY) / 2
			c.Line(fromX, fromY+fromH/2, fromX, midY,
				canvas.Stroke(canvas.ColorEdge), canvas.StrokeWidth(2))
			c.Line(fromX, midY, toX, midY,
				canvas.Stroke(canvas.ColorEdge), canvas.StrokeWidth(2))
			c.Line(toX, midY, toX, toY-toH/2,
				canvas.Stroke(canvas.ColorEdge), canvas.StrokeWidth(2))
		}
	} else {
		// Draw through waypoints
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

		// Connect to target
		endY := toY - toH/2
		if prevY > toY {
			endY = toY + toH/2
		}
		c.Line(prevX, prevY, toX, endY,
			canvas.Stroke(canvas.ColorEdge), canvas.StrokeWidth(2))
	}

	// Draw arrow at end
	drawArrowhead(c, toX, toY-toH/2, "down")

	// Draw label if defined
	if edge.LabelPos.X > 0 || edge.LabelPos.Y > 0 {
		labelX := int(edge.LabelPos.X * float64(layout.MinWidth))
		labelY := int(edge.LabelPos.Y * float64(layout.MinHeight))
		c.Text(labelX, labelY, edge.FromID+"→"+edge.ToID,
			canvas.TextAnchor("middle"), canvas.Fill(canvas.ColorEdgeLabel), canvas.FontSize(9))
	}
}

func drawArrowhead(c *canvas.Canvas, x, y int, direction string) {
	size := 8
	switch direction {
	case "down":
		points := fmt.Sprintf("%d,%d %d,%d %d,%d",
			x, y, x-size/2, y-size, x+size/2, y-size)
		c.Polygon(points, canvas.Fill(canvas.ColorEdge))
	case "up":
		points := fmt.Sprintf("%d,%d %d,%d %d,%d",
			x, y, x-size/2, y+size, x+size/2, y+size)
		c.Polygon(points, canvas.Fill(canvas.ColorEdge))
	case "right":
		points := fmt.Sprintf("%d,%d %d,%d %d,%d",
			x, y, x-size, y-size/2, x-size, y+size/2)
		c.Polygon(points, canvas.Fill(canvas.ColorEdge))
	case "left":
		points := fmt.Sprintf("%d,%d %d,%d %d,%d",
			x, y, x+size, y-size/2, x+size, y+size/2)
		c.Polygon(points, canvas.Fill(canvas.ColorEdge))
	}
}

func generateIndex(outDir string, patterns []canvas.PatternType) {
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

    <h2>Class Diagram Patterns</h2>
    <div class="grid">
`
	// Class patterns
	classPatterns := []canvas.PatternType{
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
	}
	for _, p := range classPatterns {
		html += fmt.Sprintf(`        <div class="card">
            <img src="%s.svg" alt="%s">
            <div class="card-title">%s</div>
        </div>
`, p, p, p)
	}

	html += `    </div>

    <h2>State Diagram Patterns</h2>
    <div class="grid">
`
	// State patterns
	statePatterns := []canvas.PatternType{
		canvas.PatternLinearStates2,
		canvas.PatternLinearStates3,
		canvas.PatternLinearStates4,
		canvas.PatternBinaryChoice,
		canvas.PatternStateLoop,
		canvas.PatternStarTopology,
	}
	for _, p := range statePatterns {
		html += fmt.Sprintf(`        <div class="card">
            <img src="%s.svg" alt="%s">
            <div class="card-title">%s</div>
        </div>
`, p, p, p)
	}

	html += `    </div>

    <h2>Flow Diagram Patterns</h2>
    <div class="grid">
`
	// Flow patterns
	flowPatterns := []canvas.PatternType{
		canvas.PatternIfElse,
		canvas.PatternIfElseIfElse,
		canvas.PatternWhileLoop,
		canvas.PatternSequential3,
		canvas.PatternSequential4,
	}
	for _, p := range flowPatterns {
		html += fmt.Sprintf(`        <div class="card">
            <img src="%s.svg" alt="%s">
            <div class="card-title">%s</div>
        </div>
`, p, p, p)
	}

	html += `    </div>

    <h2>Sequence Diagram Patterns</h2>
    <div class="grid">
`
	// Sequence patterns
	seqPatterns := []canvas.PatternType{
		canvas.PatternRequestResponse,
		canvas.PatternCallback,
		canvas.PatternChain3,
		canvas.PatternChain4,
		canvas.PatternFanOut,
	}
	for _, p := range seqPatterns {
		html += fmt.Sprintf(`        <div class="card">
            <img src="%s.svg" alt="%s">
            <div class="card-title">%s</div>
        </div>
`, p, p, p)
	}

	html += `    </div>
</body>
</html>
`
	_ = os.WriteFile(filepath.Join(outDir, "index.html"), []byte(html), 0644)
}
