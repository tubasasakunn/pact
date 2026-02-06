package canvas

func (r *DecorationRegistry) registerGradients() {
	// Elegant blue gradient for headers
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-header-blue",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#4299e1", Opacity: "1"},
			{Offset: "100%", Color: "#3182ce", Opacity: "1"},
		},
	})

	// Soft purple gradient for interfaces
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-interface-purple",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#b794f4", Opacity: "1"},
			{Offset: "100%", Color: "#9f7aea", Opacity: "1"},
		},
	})

	// Green gradient for success states
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-success-green",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#68d391", Opacity: "1"},
			{Offset: "100%", Color: "#48bb78", Opacity: "1"},
		},
	})

	// Warm orange gradient for process nodes
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-process-orange",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#fbd38d", Opacity: "1"},
			{Offset: "100%", Color: "#f6ad55", Opacity: "1"},
		},
	})

	// Subtle gray gradient for backgrounds
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-bg-subtle",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#ffffff", Opacity: "1"},
			{Offset: "100%", Color: "#f7fafc", Opacity: "1"},
		},
	})

	// Radial glow for highlights
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-glow-blue",
		Type: "radial",
		Attrs: map[string]string{
			"cx": "50%", "cy": "50%", "r": "50%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#63b3ed", Opacity: "0.4"},
			{Offset: "100%", Color: "#63b3ed", Opacity: "0"},
		},
	})

	// Cool teal gradient for database nodes
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-database-teal",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#4fd1c5", Opacity: "1"},
			{Offset: "100%", Color: "#38b2ac", Opacity: "1"},
		},
	})

	// Red gradient for error/exception states
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-error-red",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#fc8181", Opacity: "1"},
			{Offset: "100%", Color: "#f56565", Opacity: "1"},
		},
	})

	// Elegant dark gradient for terminal nodes
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-terminal-dark",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "0%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#4a5568", Opacity: "1"},
			{Offset: "100%", Color: "#2d3748", Opacity: "1"},
		},
	})

	// Pink gradient for notes and annotations
	r.gradients = append(r.gradients, GradientDef{
		ID:   "grad-note-pink",
		Type: "linear",
		Attrs: map[string]string{
			"x1": "0%", "y1": "0%", "x2": "100%", "y2": "100%",
		},
		Stops: []GradientStop{
			{Offset: "0%", Color: "#fef3c7", Opacity: "1"},
			{Offset: "100%", Color: "#fde68a", Opacity: "1"},
		},
	})
}
