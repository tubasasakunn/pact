// Package canvas provides pattern detection for diagram layouts.
package canvas

import (
	"pact/internal/domain/diagram/class"
	"pact/internal/domain/diagram/sequence"
)

func filterEdgesByType(edges []class.Edge, edgeType class.EdgeType) []class.Edge {
	var result []class.Edge
	for _, e := range edges {
		if e.Type == edgeType {
			result = append(result, e)
		}
	}
	return result
}

// Helper functions for role IDs
func childRoleID(i int) string   { return "child_" + itoa(i) }
func implRoleID(i int) string    { return "impl_" + itoa(i) }
func partRoleID(i int) string    { return "part_" + itoa(i) }
func stateRoleID(i int) string   { return "state_" + itoa(i) }
func nodeRoleID(i int) string    { return "node_" + itoa(i) }
func processRoleID(i int) string { return "process_" + itoa(i) }

func itoa(i int) string {
	return string(rune('0' + i))
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
