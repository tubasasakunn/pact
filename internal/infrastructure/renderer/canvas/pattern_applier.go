// Package canvas provides pattern layout application for diagrams.
package canvas

import (
	"pact/internal/infrastructure/renderer/geom"
)

// NodeLayout represents calculated absolute positions for a node
type NodeLayout struct {
	ID     string
	X      int
	Y      int
	Width  int
	Height int
}

// EdgeLayout represents calculated waypoints for an edge
type EdgeLayout struct {
	FromID    string
	ToID      string
	Waypoints []AbsolutePoint
	LabelX    int
	LabelY    int
}

// AbsolutePoint represents an absolute pixel coordinate
type AbsolutePoint struct {
	X int
	Y int
}

// AppliedLayout contains the complete layout result for a pattern
type AppliedLayout struct {
	Pattern     PatternType
	Width       int
	Height      int
	Nodes       []NodeLayout
	Edges       []EdgeLayout
	Decorators  []AppliedDecorator
}

// AppliedDecorator contains calculated decorator positions
type AppliedDecorator struct {
	Type   string
	X      int
	Y      int
	Width  int
	Height int
	Style  map[string]string
}

// PatternLayoutApplier applies pattern layouts to diagrams
type PatternLayoutApplier struct {
	registry *PatternRegistry
}

// NewPatternLayoutApplier creates a new pattern layout applier
func NewPatternLayoutApplier(registry *PatternRegistry) *PatternLayoutApplier {
	return &PatternLayoutApplier{registry: registry}
}

// ApplyClassPattern applies a class pattern layout
func (a *PatternLayoutApplier) ApplyClassPattern(match ClassPatternMatch, nodeWidths, nodeHeights map[string]int) *AppliedLayout {
	layout := a.registry.Get(match.Pattern)
	if layout == nil {
		return nil
	}

	// Calculate canvas size based on content
	canvasWidth, canvasHeight := a.calculateCanvasSize(layout, nodeWidths, nodeHeights, match.NodeRoles)

	// Apply node positions
	nodes := a.applyNodePositions(layout, match.NodeRoles, nodeWidths, nodeHeights, canvasWidth, canvasHeight)

	// Apply edge layouts
	edges := a.applyEdgeLayouts(layout, match.NodeRoles, nodes, canvasWidth, canvasHeight)

	// Apply decorators
	decorators := a.applyDecorators(layout, canvasWidth, canvasHeight)

	return &AppliedLayout{
		Pattern:    match.Pattern,
		Width:      canvasWidth,
		Height:     canvasHeight,
		Nodes:      nodes,
		Edges:      edges,
		Decorators: decorators,
	}
}

// ApplyStatePattern applies a state pattern layout
func (a *PatternLayoutApplier) ApplyStatePattern(match StatePatternMatch, stateWidths, stateHeights map[string]int) *AppliedLayout {
	layout := a.registry.Get(match.Pattern)
	if layout == nil {
		return nil
	}

	canvasWidth, canvasHeight := a.calculateCanvasSize(layout, stateWidths, stateHeights, match.StateRoles)
	nodes := a.applyNodePositions(layout, match.StateRoles, stateWidths, stateHeights, canvasWidth, canvasHeight)
	edges := a.applyEdgeLayouts(layout, match.StateRoles, nodes, canvasWidth, canvasHeight)
	decorators := a.applyDecorators(layout, canvasWidth, canvasHeight)

	return &AppliedLayout{
		Pattern:    match.Pattern,
		Width:      canvasWidth,
		Height:     canvasHeight,
		Nodes:      nodes,
		Edges:      edges,
		Decorators: decorators,
	}
}

// ApplyFlowPattern applies a flow pattern layout
func (a *PatternLayoutApplier) ApplyFlowPattern(match FlowPatternMatch, nodeWidths, nodeHeights map[string]int) *AppliedLayout {
	layout := a.registry.Get(match.Pattern)
	if layout == nil {
		return nil
	}

	canvasWidth, canvasHeight := a.calculateCanvasSize(layout, nodeWidths, nodeHeights, match.NodeRoles)
	nodes := a.applyNodePositions(layout, match.NodeRoles, nodeWidths, nodeHeights, canvasWidth, canvasHeight)
	edges := a.applyEdgeLayouts(layout, match.NodeRoles, nodes, canvasWidth, canvasHeight)
	decorators := a.applyDecorators(layout, canvasWidth, canvasHeight)

	return &AppliedLayout{
		Pattern:    match.Pattern,
		Width:      canvasWidth,
		Height:     canvasHeight,
		Nodes:      nodes,
		Edges:      edges,
		Decorators: decorators,
	}
}

// ApplySequencePattern applies a sequence pattern layout
func (a *PatternLayoutApplier) ApplySequencePattern(match SequencePatternMatch, participantWidths map[string]int) *AppliedLayout {
	layout := a.registry.Get(match.Pattern)
	if layout == nil {
		return nil
	}

	// Sequence diagrams have fixed height per participant
	heights := make(map[string]int)
	for id := range participantWidths {
		heights[id] = 60 // Standard participant box height
	}

	canvasWidth, canvasHeight := a.calculateCanvasSize(layout, participantWidths, heights, match.ParticipantRoles)
	nodes := a.applyNodePositions(layout, match.ParticipantRoles, participantWidths, heights, canvasWidth, canvasHeight)
	edges := a.applyEdgeLayouts(layout, match.ParticipantRoles, nodes, canvasWidth, canvasHeight)

	return &AppliedLayout{
		Pattern: match.Pattern,
		Width:   canvasWidth,
		Height:  canvasHeight,
		Nodes:   nodes,
		Edges:   edges,
	}
}

func (a *PatternLayoutApplier) calculateCanvasSize(layout *PatternLayout, widths, heights map[string]int, roles map[string]string) (int, int) {
	// Start with minimum size
	width := layout.MinWidth
	height := layout.MinHeight

	// Calculate required size based on actual node dimensions
	for roleID, actualID := range roles {
		nodeW, hasW := widths[actualID]
		nodeH, hasH := heights[actualID]

		if !hasW || !hasH {
			continue
		}

		// Find the position definition for this role
		for _, pos := range layout.Positions {
			if pos.ID == roleID {
				// Calculate required canvas size to fit this node
				requiredWidth := int(float64(nodeW) / pos.Width)
				requiredHeight := int(float64(nodeH) / pos.Height)

				if requiredWidth > width {
					width = requiredWidth
				}
				if requiredHeight > height {
					height = requiredHeight
				}
				break
			}
		}
	}

	// Add padding
	width += layout.Padding * 2
	height += layout.Padding * 2

	return width, height
}

func (a *PatternLayoutApplier) applyNodePositions(layout *PatternLayout, roles map[string]string, widths, heights map[string]int, canvasWidth, canvasHeight int) []NodeLayout {
	var nodes []NodeLayout

	for roleID, actualID := range roles {
		// Find position definition
		for _, pos := range layout.Positions {
			if pos.ID == roleID {
				nodeW := widths[actualID]
				nodeH := heights[actualID]

				// Calculate center position
				centerX := int(pos.X * float64(canvasWidth))
				centerY := int(pos.Y * float64(canvasHeight))

				// Calculate top-left corner
				x := centerX - nodeW/2
				y := centerY - nodeH/2

				nodes = append(nodes, NodeLayout{
					ID:     actualID,
					X:      x,
					Y:      y,
					Width:  nodeW,
					Height: nodeH,
				})
				break
			}
		}
	}

	return nodes
}

func (a *PatternLayoutApplier) applyEdgeLayouts(layout *PatternLayout, roles map[string]string, nodes []NodeLayout, canvasWidth, canvasHeight int) []EdgeLayout {
	var edges []EdgeLayout

	// Build node lookup
	nodeByID := make(map[string]NodeLayout)
	for _, n := range nodes {
		nodeByID[n.ID] = n
	}

	// Reverse role lookup
	roleToActual := roles
	actualToRole := make(map[string]string)
	for role, actual := range roles {
		actualToRole[actual] = role
	}

	for _, edgeDef := range layout.Edges {
		fromActual, hasFrom := roleToActual[edgeDef.FromID]
		toActual, hasTo := roleToActual[edgeDef.ToID]

		if !hasFrom || !hasTo {
			continue
		}

		fromNode, hasFromNode := nodeByID[fromActual]
		toNode, hasToNode := nodeByID[toActual]

		if !hasFromNode || !hasToNode {
			continue
		}

		// Calculate waypoints
		waypoints := a.calculateWaypoints(edgeDef, fromNode, toNode, canvasWidth, canvasHeight)

		// Calculate label position
		labelX := int(edgeDef.LabelPos.X * float64(canvasWidth))
		labelY := int(edgeDef.LabelPos.Y * float64(canvasHeight))

		edges = append(edges, EdgeLayout{
			FromID:    fromActual,
			ToID:      toActual,
			Waypoints: waypoints,
			LabelX:    labelX,
			LabelY:    labelY,
		})
	}

	return edges
}

func (a *PatternLayoutApplier) calculateWaypoints(edgeDef EdgePath, from, to NodeLayout, canvasWidth, canvasHeight int) []AbsolutePoint {
	var points []AbsolutePoint

	// Start point (center of from node edge)
	fromCenterX := from.X + from.Width/2
	fromCenterY := from.Y + from.Height/2
	toCenterX := to.X + to.Width/2
	toCenterY := to.Y + to.Height/2

	// Determine exit point from source node
	startX, startY := a.calculateConnectionPoint(from, toCenterX, toCenterY)
	points = append(points, AbsolutePoint{X: startX, Y: startY})

	// Add intermediate waypoints if defined
	switch edgeDef.CurveStyle {
	case "orthogonal":
		for _, wp := range edgeDef.Waypoints {
			points = append(points, AbsolutePoint{
				X: int(wp.X * float64(canvasWidth)),
				Y: int(wp.Y * float64(canvasHeight)),
			})
		}
	case "curved":
		// For curved, we add control points
		midX := (fromCenterX + toCenterX) / 2
		midY := (fromCenterY + toCenterY) / 2
		points = append(points, AbsolutePoint{X: midX, Y: midY})
	}
	// "straight" case: no intermediate points

	// End point (center of to node edge)
	endX, endY := a.calculateConnectionPoint(to, fromCenterX, fromCenterY)
	points = append(points, AbsolutePoint{X: endX, Y: endY})

	return points
}

func (a *PatternLayoutApplier) calculateConnectionPoint(node NodeLayout, targetX, targetY int) (int, int) {
	centerX := node.X + node.Width/2
	centerY := node.Y + node.Height/2

	dx := targetX - centerX
	dy := targetY - centerY

	// Determine which edge to connect to
	if geom.Abs(dx)*node.Height > geom.Abs(dy)*node.Width {
		// Connect to left or right edge
		if dx > 0 {
			return node.X + node.Width, centerY // Right edge
		}
		return node.X, centerY // Left edge
	}
	// Connect to top or bottom edge
	if dy > 0 {
		return centerX, node.Y + node.Height // Bottom edge
	}
	return centerX, node.Y // Top edge
}

func (a *PatternLayoutApplier) applyDecorators(layout *PatternLayout, canvasWidth, canvasHeight int) []AppliedDecorator {
	var decorators []AppliedDecorator

	for _, dec := range layout.Decorators {
		decorators = append(decorators, AppliedDecorator{
			Type:   dec.Type,
			X:      int(dec.Bounds.X * float64(canvasWidth)),
			Y:      int(dec.Bounds.Y * float64(canvasHeight)),
			Width:  int(dec.Bounds.Width * float64(canvasWidth)),
			Height: int(dec.Bounds.Height * float64(canvasHeight)),
			Style:  dec.Style,
		})
	}

	return decorators
}

// GetBestMatch returns the highest scoring match from a list
func GetBestMatch[T interface{ GetScore() float64 }](matches []T) (T, bool) {
	var best T
	var found bool
	var bestScore float64

	for _, m := range matches {
		score := m.GetScore()
		if score > bestScore {
			bestScore = score
			best = m
			found = true
		}
	}

	return best, found
}

// GetScore implements score accessor for ClassPatternMatch
func (m ClassPatternMatch) GetScore() float64 { return m.Score }

// GetScore implements score accessor for StatePatternMatch
func (m StatePatternMatch) GetScore() float64 { return m.Score }

// GetScore implements score accessor for FlowPatternMatch
func (m FlowPatternMatch) GetScore() float64 { return m.Score }

// GetScore implements score accessor for SequencePatternMatch
func (m SequencePatternMatch) GetScore() float64 { return m.Score }
