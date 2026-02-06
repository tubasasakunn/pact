package canvas

import (
	"pact/internal/domain/diagram/sequence"
)

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
