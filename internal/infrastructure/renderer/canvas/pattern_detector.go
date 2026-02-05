// Package canvas provides pattern detection for diagram layouts.
package canvas

import (
	"pact/internal/domain/diagram/class"
	"pact/internal/domain/diagram/flow"
	"pact/internal/domain/diagram/sequence"
	"pact/internal/domain/diagram/state"
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

func filterEdgesByType(edges []class.Edge, edgeType class.EdgeType) []class.Edge {
	var result []class.Edge
	for _, e := range edges {
		if e.Type == edgeType {
			result = append(result, e)
		}
	}
	return result
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

// Helper functions for role IDs
func childRoleID(i int) string  { return "child_" + itoa(i) }
func implRoleID(i int) string   { return "impl_" + itoa(i) }
func partRoleID(i int) string   { return "part_" + itoa(i) }
func stateRoleID(i int) string  { return "state_" + itoa(i) }
func nodeRoleID(i int) string   { return "node_" + itoa(i) }
func processRoleID(i int) string { return "process_" + itoa(i) }

func itoa(i int) string {
	return string(rune('0' + i))
}

// StatePatternDetector detects patterns in state diagrams
type StatePatternDetector struct {
	registry *PatternRegistry
}

// NewStatePatternDetector creates a new state pattern detector
func NewStatePatternDetector(registry *PatternRegistry) *StatePatternDetector {
	return &StatePatternDetector{registry: registry}
}

// Detect analyzes a state diagram and returns matching patterns
func (d *StatePatternDetector) Detect(diagram *state.Diagram) []StatePatternMatch {
	var matches []StatePatternMatch

	// Build transition graph
	outgoing := make(map[string][]string)
	incoming := make(map[string][]string)
	for _, t := range diagram.Transitions {
		outgoing[t.From] = append(outgoing[t.From], t.To)
		incoming[t.To] = append(incoming[t.To], t.From)
	}

	// Find initial and final states
	var initialState, finalState string
	atomicStates := make([]string, 0)
	for _, s := range diagram.States {
		switch s.Type {
		case state.StateTypeInitial:
			initialState = s.ID
		case state.StateTypeFinal:
			finalState = s.ID
		case state.StateTypeAtomic:
			atomicStates = append(atomicStates, s.ID)
		}
	}

	// Detect linear pattern
	if match := d.detectLinear(atomicStates, initialState, finalState, outgoing); match != nil {
		matches = append(matches, *match)
	}

	// Detect binary choice
	if match := d.detectBinaryChoice(atomicStates, outgoing, incoming); match != nil {
		matches = append(matches, *match)
	}

	// Detect loop pattern
	if match := d.detectLoop(atomicStates, outgoing); match != nil {
		matches = append(matches, *match)
	}

	// Detect star topology
	if match := d.detectStar(atomicStates, outgoing, incoming); match != nil {
		matches = append(matches, *match)
	}

	return matches
}

func (d *StatePatternDetector) detectLinear(states []string, initial, final string, outgoing map[string][]string) *StatePatternMatch {
	if len(states) < 2 {
		return nil
	}

	// Check if each state has exactly one outgoing transition
	linearChain := make([]string, 0)

	// Start from initial state's target
	if initial != "" {
		targets := outgoing[initial]
		if len(targets) == 1 {
			current := targets[0]
			visited := make(map[string]bool)

			for current != "" && current != final && !visited[current] {
				visited[current] = true
				linearChain = append(linearChain, current)

				targets := outgoing[current]
				if len(targets) != 1 {
					break
				}
				current = targets[0]
			}
		}
	}

	if len(linearChain) < 2 {
		return nil
	}

	stateRoles := make(map[string]string)
	if initial != "" {
		stateRoles["initial"] = initial
	}
	for i, s := range linearChain {
		if i < 4 {
			stateRoles[stateRoleID(i)] = s
		}
	}
	if final != "" {
		stateRoles["final"] = final
	}

	return &StatePatternMatch{
		Pattern:    PatternLinearStates,
		StateRoles: stateRoles,
		Score:      float64(len(linearChain)) / float64(len(states)),
	}
}

func (d *StatePatternDetector) detectBinaryChoice(states []string, outgoing, incoming map[string][]string) *StatePatternMatch {
	// Find a state with exactly 2 outgoing transitions to different targets
	// that both eventually converge
	for _, source := range states {
		targets := outgoing[source]
		if len(targets) != 2 {
			continue
		}

		true_branch := targets[0]
		false_branch := targets[1]

		// Check if both branches have outgoing edges
		trueTargets := outgoing[true_branch]
		falseTargets := outgoing[false_branch]

		if len(trueTargets) == 0 || len(falseTargets) == 0 {
			continue
		}

		// Check for convergence
		for _, tt := range trueTargets {
			for _, ft := range falseTargets {
				if tt == ft {
					return &StatePatternMatch{
						Pattern: PatternBinaryChoice,
						StateRoles: map[string]string{
							"source":       source,
							"true_branch":  true_branch,
							"false_branch": false_branch,
							"target":       tt,
						},
						Score: 1.0,
					}
				}
			}
		}
	}

	return nil
}

func (d *StatePatternDetector) detectLoop(states []string, outgoing map[string][]string) *StatePatternMatch {
	// Find a cycle in the state graph
	for _, start := range states {
		visited := make(map[string]bool)
		path := []string{start}

		current := start
		for {
			targets := outgoing[current]
			if len(targets) == 0 {
				break
			}

			// Look for a back edge to start
			for _, t := range targets {
				if t == start && len(path) >= 2 {
					stateRoles := make(map[string]string)
					stateRoles["start"] = path[0]
					if len(path) > 1 {
						stateRoles["middle"] = path[len(path)/2]
					}
					stateRoles["end"] = path[len(path)-1]

					return &StatePatternMatch{
						Pattern:    PatternStateLoop,
						StateRoles: stateRoles,
						Score:      0.9,
					}
				}
			}

			// Continue traversal
			nextFound := false
			for _, t := range targets {
				if !visited[t] && t != start {
					visited[t] = true
					path = append(path, t)
					current = t
					nextFound = true
					break
				}
			}

			if !nextFound {
				break
			}
		}
	}

	return nil
}

func (d *StatePatternDetector) detectStar(states []string, outgoing, incoming map[string][]string) *StatePatternMatch {
	// Find a central hub state connected to multiple peripheral states
	for _, center := range states {
		out := outgoing[center]
		in := incoming[center]

		// Hub should have multiple connections
		totalConnections := len(out) + len(in)
		if totalConnections < 4 {
			continue
		}

		stateRoles := make(map[string]string)
		stateRoles["center"] = center

		// Add peripheral states
		peripherals := make(map[string]bool)
		for _, s := range out {
			peripherals[s] = true
		}
		for _, s := range in {
			peripherals[s] = true
		}

		positions := []string{"top", "right", "bottom", "left"}
		i := 0
		for p := range peripherals {
			if i < 4 {
				stateRoles[positions[i]] = p
				i++
			}
		}

		if i >= 3 {
			return &StatePatternMatch{
				Pattern:    PatternStarTopology,
				StateRoles: stateRoles,
				Score:      float64(i) / 4.0,
			}
		}
	}

	return nil
}

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

// SequencePatternDetector detects patterns in sequence diagrams
type SequencePatternDetector struct {
	registry *PatternRegistry
}

// NewSequencePatternDetector creates a new sequence pattern detector
func NewSequencePatternDetector(registry *PatternRegistry) *SequencePatternDetector {
	return &SequencePatternDetector{registry: registry}
}

// Detect analyzes a sequence diagram and returns matching patterns
func (d *SequencePatternDetector) Detect(diagram *sequence.Diagram) []SequencePatternMatch {
	var matches []SequencePatternMatch

	// Extract messages
	messages := extractMessages(diagram.Events)
	if len(messages) < 2 {
		return matches
	}

	// Detect request-response pattern
	if match := d.detectRequestResponse(messages, diagram.Participants); match != nil {
		matches = append(matches, *match)
	}

	// Detect callback pattern
	if match := d.detectCallback(messages, diagram.Participants); match != nil {
		matches = append(matches, *match)
	}

	// Detect chain pattern
	if match := d.detectChain(messages, diagram.Participants); match != nil {
		matches = append(matches, *match)
	}

	return matches
}

func extractMessages(events []sequence.Event) []*sequence.MessageEvent {
	var messages []*sequence.MessageEvent
	for _, e := range events {
		if msg, ok := e.(*sequence.MessageEvent); ok {
			messages = append(messages, msg)
		}
	}
	return messages
}

func (d *SequencePatternDetector) detectRequestResponse(messages []*sequence.MessageEvent, participants []sequence.Participant) *SequencePatternMatch {
	if len(messages) < 2 {
		return nil
	}

	// Look for A->B followed by B->A pattern
	for i := 0; i < len(messages)-1; i++ {
		req := messages[i]
		resp := messages[i+1]

		if req.From == resp.To && req.To == resp.From {
			// Check if response is a return message
			if resp.MessageType == sequence.MessageTypeReturn {
				return &SequencePatternMatch{
					Pattern: PatternRequestResponse,
					ParticipantRoles: map[string]string{
						"caller": req.From,
						"callee": req.To,
					},
					Score: 1.0,
				}
			}
		}
	}

	return nil
}

func (d *SequencePatternDetector) detectCallback(messages []*sequence.MessageEvent, participants []sequence.Participant) *SequencePatternMatch {
	if len(messages) < 3 {
		return nil
	}

	// Look for A->B, B->A, A->B pattern
	for i := 0; i < len(messages)-2; i++ {
		m1 := messages[i]
		m2 := messages[i+1]
		m3 := messages[i+2]

		if m1.From == m2.To && m1.To == m2.From &&
			m2.From == m3.To && m2.To == m3.From {
			return &SequencePatternMatch{
				Pattern: PatternCallback,
				ParticipantRoles: map[string]string{
					"initiator": m1.From,
					"handler":   m1.To,
				},
				Score: 0.95,
			}
		}
	}

	return nil
}

func (d *SequencePatternDetector) detectChain(messages []*sequence.MessageEvent, participants []sequence.Participant) *SequencePatternMatch {
	if len(messages) < 3 || len(participants) < 3 {
		return nil
	}

	// Look for A->B->C->... chain
	chain := []string{messages[0].From}

	for _, msg := range messages {
		if msg.MessageType == sequence.MessageTypeReturn {
			continue
		}

		last := chain[len(chain)-1]
		if msg.From == last && msg.To != last {
			// Check if target is not already in chain (no loop back)
			isNew := true
			for _, c := range chain {
				if c == msg.To {
					isNew = false
					break
				}
			}
			if isNew {
				chain = append(chain, msg.To)
			}
		}
	}

	if len(chain) < 3 {
		return nil
	}

	roles := make(map[string]string)
	for i, p := range chain {
		if i < 4 {
			roles[nodeRoleID(i)] = p
		}
	}

	return &SequencePatternMatch{
		Pattern:          PatternChain,
		ParticipantRoles: roles,
		Score:            float64(len(chain)) / float64(len(participants)),
	}
}
