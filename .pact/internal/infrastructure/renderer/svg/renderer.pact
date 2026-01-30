// SVG Renderers - Diagram-to-SVG rendering components
// Provides renderers for class, sequence, state, and flow diagrams

@layer("infrastructure")
@package("svg")
component ClassRenderer {
    depends on Canvas

    type ClassRenderer {
        // Stateless renderer for class diagrams
    }

    provides ClassRendererFactory {
        // Creates a new ClassRenderer instance
        NewClassRenderer() -> ClassRenderer
    }

    provides ClassRendererAPI {
        // Renders a class diagram to SVG output
        // Positions nodes in a 3-column grid layout
        // Draws edges with appropriate decorations (triangles, diamonds)
        Render(diagram: ClassDiagram, w: Writer) -> error
    }

    provides ClassRendererInternal {
        // Renders a single class node with stereotype, name, attributes, and methods
        renderNode(c: Canvas, node: ClassNode, x: int, y: int)

        // Renders an edge between classes with line style and decoration
        renderEdge(c: Canvas, edge: ClassEdge, x1: int, y1: int, x2: int, y2: int)
    }

    provides VisibilityHelper {
        // Converts visibility enum to UML symbol (+ public, - private, # protected, ~ package)
        visibilitySymbol(v: Visibility) -> string

        // Generates SVG points for a triangle arrowhead
        trianglePoints(x: int, y: int) -> string

        // Generates SVG points for a diamond decoration
        diamondPoints(x: int, y: int) -> string
    }
}

component SequenceRenderer {
    depends on Canvas

    type SequenceRenderer {
        // Stateless renderer for sequence diagrams
    }

    provides SequenceRendererFactory {
        // Creates a new SequenceRenderer instance
        NewSequenceRenderer() -> SequenceRenderer
    }

    provides SequenceRendererAPI {
        // Renders a sequence diagram to SVG output
        // Positions participants horizontally with vertical lifelines
        // Draws messages with appropriate arrow styles
        Render(diagram: SequenceDiagram, w: Writer) -> error
    }

    provides SequenceRendererInternal {
        // Recursively renders sequence events (messages and fragments)
        // Handles sync, async, and return message types
        // Supports fragment boxes for alt, loop, opt constructs
        renderEvents(c: Canvas, events: Event[], participantX: map, y: int)

        // Renders a participant based on type (actor, database, default box)
        // Draws lifeline extending downward
        renderParticipant(c: Canvas, p: Participant, x: int, y: int)
    }
}

component StateRenderer {
    depends on Canvas

    type StateRenderer {
        // Stateless renderer for state diagrams
    }

    provides StateRendererFactory {
        // Creates a new StateRenderer instance
        NewStateRenderer() -> StateRenderer
    }

    provides StateRendererAPI {
        // Renders a state diagram to SVG output
        // Positions states in a 3-column grid layout
        // Draws transitions with arrows and labels
        Render(diagram: StateDiagram, w: Writer) -> error
    }

    provides StateRendererInternal {
        // Renders a state based on type (initial=filled circle, final=bullseye, normal=rounded rect)
        renderState(c: Canvas, s: State, x: int, y: int)

        // Renders a transition arrow with trigger label
        // Supports EventTrigger and WhenTrigger labels
        renderTransition(c: Canvas, t: Transition, x1: int, y1: int, x2: int, y2: int)
    }
}

component FlowRenderer {
    depends on Canvas

    type FlowRenderer {
        // Stateless renderer for flowcharts
    }

    provides FlowRendererFactory {
        // Creates a new FlowRenderer instance
        NewFlowRenderer() -> FlowRenderer
    }

    provides FlowRendererAPI {
        // Renders a flowchart to SVG output
        // Positions nodes vertically in a single column
        // Dynamically calculates height based on node count
        Render(diagram: FlowDiagram, w: Writer) -> error
    }

    provides FlowRendererInternal {
        // Renders a flow node based on shape type
        // Shapes: terminal (stadium), process (rect), decision (diamond),
        //         database (cylinder), io (parallelogram)
        renderFlowNode(c: Canvas, node: FlowNode, x: int, y: int)

        // Renders a flow edge with arrow and optional label
        // Labels appear at edge midpoint (e.g., Yes/No for decisions)
        renderFlowEdge(c: Canvas, edge: FlowEdge, x1: int, y1: int, x2: int, y2: int)
    }
}
