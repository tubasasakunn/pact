package canvas

import (
	"pact/internal/domain/diagram/flow"
)

// FlowPatternDetector detects patterns in flow diagrams
type FlowPatternDetector struct {
	registry *PatternRegistry
}

// NewFlowPatternDetector creates a new flow pattern detector
func NewFlowPatternDetector(registry *PatternRegistry) *FlowPatternDetector {
	return &FlowPatternDetector{registry: registry}
}

// Detect analyzes a flow diagram and returns matching patterns
func (d *FlowPatternDetector) Detect(diagram *flow.Diagram) []FlowPatternMatch {
	var matches []FlowPatternMatch

	// Build node index and edge graph
	nodeByID := make(map[string]*flow.Node)
	outgoing := make(map[string][]string)
	incoming := make(map[string][]string)

	for i := range diagram.Nodes {
		n := &diagram.Nodes[i]
		nodeByID[n.ID] = n
	}
	for _, e := range diagram.Edges {
		outgoing[e.From] = append(outgoing[e.From], e.To)
		incoming[e.To] = append(incoming[e.To], e.From)
	}

	// Find terminal and decision nodes
	var startNodes, endNodes, decisionNodes, processNodes []string
	for _, n := range diagram.Nodes {
		switch n.Shape {
		case flow.NodeShapeTerminal:
			if len(outgoing[n.ID]) > 0 && len(incoming[n.ID]) == 0 {
				startNodes = append(startNodes, n.ID)
			} else if len(incoming[n.ID]) > 0 && len(outgoing[n.ID]) == 0 {
				endNodes = append(endNodes, n.ID)
			}
		case flow.NodeShapeDecision:
			decisionNodes = append(decisionNodes, n.ID)
		case flow.NodeShapeProcess:
			processNodes = append(processNodes, n.ID)
		}
	}

	// Detect if-else pattern
	if match := d.detectIfElse(decisionNodes, processNodes, nodeByID, outgoing, incoming); match != nil {
		matches = append(matches, *match)
	}

	// Detect while loop pattern
	if match := d.detectWhileLoop(decisionNodes, processNodes, outgoing, incoming); match != nil {
		matches = append(matches, *match)
	}

	// Detect sequential pattern
	if match := d.detectSequential(startNodes, endNodes, processNodes, outgoing); match != nil {
		matches = append(matches, *match)
	}

	return matches
}

func (d *FlowPatternDetector) detectIfElse(decisions, processes []string, nodeByID map[string]*flow.Node, outgoing, incoming map[string][]string) *FlowPatternMatch {
	for _, dec := range decisions {
		targets := outgoing[dec]
		if len(targets) != 2 {
			continue
		}

		trueBranch := targets[0]
		falseBranch := targets[1]

		// Check if both branches lead to process nodes
		trueNode := nodeByID[trueBranch]
		falseNode := nodeByID[falseBranch]

		if trueNode == nil || falseNode == nil {
			continue
		}
		if trueNode.Shape != flow.NodeShapeProcess || falseNode.Shape != flow.NodeShapeProcess {
			continue
		}

		// Check if both branches converge
		trueTargets := outgoing[trueBranch]
		falseTargets := outgoing[falseBranch]

		for _, tt := range trueTargets {
			for _, ft := range falseTargets {
				if tt == ft {
					return &FlowPatternMatch{
						Pattern: PatternIfElse,
						NodeRoles: map[string]string{
							"decision":      dec,
							"true_process":  trueBranch,
							"false_process": falseBranch,
							"merge":         tt,
						},
						Score: 1.0,
					}
				}
			}
		}
	}

	return nil
}

func (d *FlowPatternDetector) detectWhileLoop(decisions, processes []string, outgoing, incoming map[string][]string) *FlowPatternMatch {
	for _, dec := range decisions {
		targets := outgoing[dec]
		if len(targets) != 2 {
			continue
		}

		// Check if one branch loops back to before the decision
		for _, t := range targets {
			bodyTargets := outgoing[t]
			for _, bt := range bodyTargets {
				// Check if this target leads back to decision
				if bt == dec {
					// Find the exit path
					var exit string
					for _, t2 := range targets {
						if t2 != t {
							exit = t2
							break
						}
					}

					return &FlowPatternMatch{
						Pattern: PatternWhileLoop,
						NodeRoles: map[string]string{
							"condition": dec,
							"body":      t,
							"exit":      exit,
						},
						Score: 0.95,
					}
				}
			}
		}
	}

	return nil
}

func (d *FlowPatternDetector) detectSequential(starts, ends, processes []string, outgoing map[string][]string) *FlowPatternMatch {
	if len(starts) == 0 {
		return nil
	}

	// Follow the chain from start
	start := starts[0]
	chain := []string{}
	visited := make(map[string]bool)

	current := start
	for {
		targets := outgoing[current]
		if len(targets) != 1 {
			break
		}

		next := targets[0]
		if visited[next] {
			break
		}

		visited[next] = true
		chain = append(chain, next)
		current = next
	}

	if len(chain) < 2 {
		return nil
	}

	nodeRoles := make(map[string]string)
	nodeRoles["start"] = start
	for i, p := range chain {
		if i < 4 {
			nodeRoles[processRoleID(i)] = p
		}
	}
	if len(ends) > 0 {
		nodeRoles["end"] = ends[0]
	}

	return &FlowPatternMatch{
		Pattern:   PatternSequential,
		NodeRoles: nodeRoles,
		Score:     float64(len(chain)) / float64(len(processes)+2),
	}
}
