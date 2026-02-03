// Package renderer defines interfaces for diagram rendering.
package renderer

import (
	"io"

	"pact/internal/domain/diagram/class"
	"pact/internal/domain/diagram/flow"
	"pact/internal/domain/diagram/sequence"
	"pact/internal/domain/diagram/state"
)

// ClassRenderer renders class diagrams.
type ClassRenderer interface {
	Render(d *class.Diagram, w io.Writer) error
}

// SequenceRenderer renders sequence diagrams.
type SequenceRenderer interface {
	Render(d *sequence.Diagram, w io.Writer) error
}

// StateRenderer renders state diagrams.
type StateRenderer interface {
	Render(d *state.Diagram, w io.Writer) error
}

// FlowRenderer renders flow diagrams.
type FlowRenderer interface {
	Render(d *flow.Diagram, w io.Writer) error
}
