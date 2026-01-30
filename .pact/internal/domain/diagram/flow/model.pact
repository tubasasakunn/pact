// Flow diagram domain model
component FlowDiagram {
    // NodeShape represents the shape of a flow node
    enum NodeShape {
        terminal
        process
        decision
        io
        database
        connector
    }

    // Swimlane represents a swimlane in the flow diagram
    type Swimlane {
        id: string
        name: string
    }

    // Node represents a node in the flow diagram
    type Node {
        id: string
        label: string
        shape: NodeShape
        swimlane: string
    }

    // Edge represents an edge connecting nodes
    type Edge {
        fromNode: string
        toNode: string
        label: string
    }

    // Diagram represents a flow diagram
    type Diagram {
        nodes: Node[]
        edges: Edge[]
        swimlanes: Swimlane[]
    }
}
