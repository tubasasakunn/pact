package canvas

import (
	"pact/internal/domain/diagram/state"
)

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
