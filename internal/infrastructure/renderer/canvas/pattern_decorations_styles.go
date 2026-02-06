package canvas

func (r *DecorationRegistry) registerStyles() {
	// Base typography
	r.styles = append(r.styles, `
.diagram-text {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  font-size: 14px;
  fill: #1a202c;
}
.diagram-text-sm {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  font-size: 12px;
  fill: #4a5568;
}
.diagram-text-bold {
  font-weight: 600;
}
.diagram-text-italic {
  font-style: italic;
}
.diagram-text-mono {
  font-family: "SF Mono", "Consolas", "Monaco", monospace;
  font-size: 13px;
}`)

	// Node styles
	r.styles = append(r.styles, `
.node-class {
  fill: #fafbfc;
  stroke: #2d3748;
  stroke-width: 2;
}
.node-interface {
  fill: #f0fff4;
  stroke: #276749;
  stroke-width: 2;
  stroke-dasharray: 5,3;
}
.node-abstract {
  fill: #faf5ff;
  stroke: #6b46c1;
  stroke-width: 2;
}
.node-state {
  fill: #fafbfc;
  stroke: #2d3748;
  stroke-width: 2;
  rx: 10;
  ry: 10;
}
.node-process {
  fill: #fafbfc;
  stroke: #2d3748;
  stroke-width: 2;
}
.node-decision {
  fill: #fffaf0;
  stroke: #c05621;
  stroke-width: 2;
}
.node-terminal {
  fill: #2d3748;
  stroke: #1a202c;
  stroke-width: 2;
}`)

	// Edge styles
	r.styles = append(r.styles, `
.edge-solid {
  fill: none;
  stroke: #2d3748;
  stroke-width: 1.5;
}
.edge-dashed {
  fill: none;
  stroke: #2d3748;
  stroke-width: 1.5;
  stroke-dasharray: 6,4;
}
.edge-dotted {
  fill: none;
  stroke: #718096;
  stroke-width: 1.5;
  stroke-dasharray: 2,3;
}
.edge-async {
  fill: none;
  stroke: #2d3748;
  stroke-width: 1.5;
  stroke-dasharray: 8,4;
}
.edge-return {
  fill: none;
  stroke: #48bb78;
  stroke-width: 1.5;
  stroke-dasharray: 4,4;
}`)

	// Special decorations
	r.styles = append(r.styles, `
.note-box {
  fill: #fefcbf;
  stroke: #d69e2e;
  stroke-width: 1.5;
}
.fragment-box {
  fill: none;
  stroke: #4a5568;
  stroke-width: 1.5;
}
.fragment-label {
  fill: #e2e8f0;
  stroke: #4a5568;
  stroke-width: 1;
}
.swimlane-header {
  fill: #edf2f7;
  stroke: #cbd5e0;
  stroke-width: 1;
}
.swimlane-divider {
  stroke: #e2e8f0;
  stroke-width: 1;
  stroke-dasharray: 4,4;
}
.activation-bar {
  fill: #e2e8f0;
  stroke: #a0aec0;
  stroke-width: 1;
}
.lifeline {
  stroke: #a0aec0;
  stroke-width: 1;
  stroke-dasharray: 6,4;
}`)

	// Hover and interaction states
	r.styles = append(r.styles, `
.interactive:hover {
  filter: brightness(1.05);
  cursor: pointer;
}
.highlight {
  filter: url(#filter-glow-blue);
}
.selected {
  stroke-width: 3;
  stroke: #4299e1;
}`)
}
