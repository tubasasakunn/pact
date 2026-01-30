// Common types shared across diagram domain models
component DiagramCommon {
    // DiagramType represents the type of diagram
    enum DiagramType {
        ClassDiagram
        SequenceDiagram
        StateDiagram
        FlowDiagram
    }

    // Annotation represents a diagram annotation
    type Annotation {
        name: string
        arguments: AnnotationArg[]
    }

    // AnnotationArg represents a key-value pair for annotation arguments
    type AnnotationArg {
        key: string
        value: string
    }
}
