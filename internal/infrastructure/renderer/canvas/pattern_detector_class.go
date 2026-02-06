package canvas

import (
	"pact/internal/domain/diagram/class"
)

// ClassPatternDetector detects patterns in class diagrams
type ClassPatternDetector struct {
	registry *PatternRegistry
}

// NewClassPatternDetector creates a new class pattern detector
func NewClassPatternDetector(registry *PatternRegistry) *ClassPatternDetector {
	return &ClassPatternDetector{registry: registry}
}

// Detect analyzes a class diagram and returns matching patterns
func (d *ClassPatternDetector) Detect(diagram *class.Diagram) []ClassPatternMatch {
	var matches []ClassPatternMatch

	// Build adjacency information
	nodeIndex := make(map[string]int)
	for i, n := range diagram.Nodes {
		nodeIndex[n.ID] = i
	}

	// Group edges by type
	inheritanceEdges := filterEdgesByType(diagram.Edges, class.EdgeTypeInheritance)
	implementationEdges := filterEdgesByType(diagram.Edges, class.EdgeTypeImplementation)
	compositionEdges := filterEdgesByType(diagram.Edges, class.EdgeTypeComposition)
	aggregationEdges := filterEdgesByType(diagram.Edges, class.EdgeTypeAggregation)

	// Detect inheritance tree pattern
	if match := d.detectInheritanceTree(diagram.Nodes, inheritanceEdges); match != nil {
		matches = append(matches, *match)
	}

	// Detect interface implementation pattern
	if match := d.detectInterfaceImpl(diagram.Nodes, implementationEdges); match != nil {
		matches = append(matches, *match)
	}

	// Detect composition pattern
	if match := d.detectComposition(diagram.Nodes, compositionEdges, aggregationEdges); match != nil {
		matches = append(matches, *match)
	}

	// Detect diamond inheritance pattern
	if match := d.detectDiamond(diagram.Nodes, inheritanceEdges); match != nil {
		matches = append(matches, *match)
	}

	return matches
}

func (d *ClassPatternDetector) detectInheritanceTree(nodes []class.Node, edges []class.Edge) *ClassPatternMatch {
	if len(edges) < 2 {
		return nil
	}

	// Find potential parent (node with multiple incoming inheritance edges)
	incomingCount := make(map[string][]string) // parent -> children
	for _, e := range edges {
		incomingCount[e.To] = append(incomingCount[e.To], e.From)
	}

	var bestParent string
	var maxChildren int
	for parent, children := range incomingCount {
		if len(children) > maxChildren {
			maxChildren = len(children)
			bestParent = parent
		}
	}

	if maxChildren < 2 {
		return nil
	}

	nodeRoles := make(map[string]string)
	nodeRoles["parent"] = bestParent
	children := incomingCount[bestParent]
	for i, child := range children {
		if i < 4 { // Max 4 children in pattern
			nodeRoles[childRoleID(i)] = child
		}
	}

	score := float64(maxChildren) / float64(len(nodes))
	if score > 1.0 {
		score = 1.0
	}

	return &ClassPatternMatch{
		Pattern:   PatternInheritanceTree,
		NodeRoles: nodeRoles,
		Score:     score,
	}
}

func (d *ClassPatternDetector) detectInterfaceImpl(nodes []class.Node, edges []class.Edge) *ClassPatternMatch {
	if len(edges) < 2 {
		return nil
	}

	// Find interface (node with <<interface>> stereotype and multiple implementers)
	incomingCount := make(map[string][]string)
	for _, e := range edges {
		incomingCount[e.To] = append(incomingCount[e.To], e.From)
	}

	var bestInterface string
	var maxImpls int
	for iface, impls := range incomingCount {
		// Check if it's an interface
		for _, n := range nodes {
			if n.ID == iface && n.Stereotype == "interface" {
				if len(impls) > maxImpls {
					maxImpls = len(impls)
					bestInterface = iface
				}
				break
			}
		}
	}

	if maxImpls < 2 {
		return nil
	}

	nodeRoles := make(map[string]string)
	nodeRoles["interface"] = bestInterface
	impls := incomingCount[bestInterface]
	for i, impl := range impls {
		if i < 4 {
			nodeRoles[implRoleID(i)] = impl
		}
	}

	return &ClassPatternMatch{
		Pattern:   PatternInterfaceImpl,
		NodeRoles: nodeRoles,
		Score:     0.9,
	}
}

func (d *ClassPatternDetector) detectComposition(nodes []class.Node, compositionEdges, aggregationEdges []class.Edge) *ClassPatternMatch {
	allEdges := append(compositionEdges, aggregationEdges...)
	if len(allEdges) < 2 {
		return nil
	}

	// Find owner (node with multiple outgoing composition/aggregation)
	outgoingCount := make(map[string][]string)
	for _, e := range allEdges {
		outgoingCount[e.From] = append(outgoingCount[e.From], e.To)
	}

	var bestOwner string
	var maxParts int
	for owner, parts := range outgoingCount {
		if len(parts) > maxParts {
			maxParts = len(parts)
			bestOwner = owner
		}
	}

	if maxParts < 2 {
		return nil
	}

	nodeRoles := make(map[string]string)
	nodeRoles["owner"] = bestOwner
	parts := outgoingCount[bestOwner]
	for i, part := range parts {
		if i < 4 {
			nodeRoles[partRoleID(i)] = part
		}
	}

	return &ClassPatternMatch{
		Pattern:   PatternComposition,
		NodeRoles: nodeRoles,
		Score:     0.85,
	}
}

func (d *ClassPatternDetector) detectDiamond(nodes []class.Node, edges []class.Edge) *ClassPatternMatch {
	if len(nodes) < 4 || len(edges) < 4 {
		return nil
	}

	// Build adjacency for diamond detection
	children := make(map[string][]string)  // parent -> children
	parents := make(map[string][]string)   // child -> parents
	for _, e := range edges {
		children[e.To] = append(children[e.To], e.From)
		parents[e.From] = append(parents[e.From], e.To)
	}

	// Find diamond: node with 2+ parents whose parents share a common ancestor
	for bottom, bottomParents := range parents {
		if len(bottomParents) < 2 {
			continue
		}

		// Check if any two parents share a common ancestor
		for i := 0; i < len(bottomParents); i++ {
			for j := i + 1; j < len(bottomParents); j++ {
				left := bottomParents[i]
				right := bottomParents[j]

				leftParents := parents[left]
				rightParents := parents[right]

				// Find common ancestor
				for _, lp := range leftParents {
					for _, rp := range rightParents {
						if lp == rp {
							// Found diamond!
							return &ClassPatternMatch{
								Pattern: PatternDiamond,
								NodeRoles: map[string]string{
									"top":    lp,
									"left":   left,
									"right":  right,
									"bottom": bottom,
								},
								Score: 1.0,
							}
						}
					}
				}
			}
		}
	}

	return nil
}
