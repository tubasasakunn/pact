// Diagram Domain Models
@version("1.0")
@package("diagram")
component DiagramModels {
    // Class Diagram
    type ClassDiagram {
        nodes: ClassNode[]
        edges: ClassEdge[]
    }

    type ClassNode {
        id: string
        name: string
        stereotype: string?
        attributes: Attribute[]
        methods: Method[]
        annotations: Annotation[]
    }

    type Attribute {
        name: string
        attrType: string
        visibility: Visibility
    }

    type Method {
        name: string
        params: Param[]
        returnType: string?
        visibility: Visibility
        isAsync: bool
    }

    type Param {
        name: string
        paramType: string
    }

    enum Visibility {
        PUBLIC
        PRIVATE
        PROTECTED
        PACKAGE
    }

    type ClassEdge {
        fromId: string
        toId: string
        edgeType: EdgeType
        label: string?
        decoration: Decoration
        lineStyle: LineStyle
    }

    enum EdgeType {
        DEPENDENCY
        INHERITANCE
        IMPLEMENTATION
        COMPOSITION
        AGGREGATION
    }

    enum Decoration {
        NONE
        ARROW
        TRIANGLE
        FILLED_DIAMOND
        EMPTY_DIAMOND
    }

    enum LineStyle {
        SOLID
        DASHED
    }

    // Sequence Diagram
    type SequenceDiagram {
        participants: Participant[]
        events: SequenceEvent[]
    }

    type Participant {
        id: string
        name: string
        participantType: ParticipantType
    }

    enum ParticipantType {
        DEFAULT
        ACTOR
        DATABASE
        QUEUE
        EXTERNAL
    }

    type MessageEvent {
        fromId: string
        toId: string
        label: string
        messageType: MessageType
    }

    enum MessageType {
        SYNC
        ASYNC
        RETURN
    }

    type FragmentEvent {
        fragmentType: FragmentType
        label: string
        events: SequenceEvent[]
        altLabel: string?
    }

    enum FragmentType {
        ALT
        LOOP
        OPT
    }

    // State Diagram
    type StateDiagram {
        stateNodes: StateNode[]
        stateTransitions: StateTransition[]
    }

    type StateNode {
        id: string
        name: string
        nodeType: StateType
        entryActions: string[]
        exitActions: string[]
        children: StateNode[]
        regions: Region[]
    }

    enum StateType {
        INITIAL
        FINAL
        ATOMIC
        COMPOUND
        PARALLEL
    }

    type Region {
        name: string
        stateNodes: StateNode[]
        transitions: StateTransition[]
    }

    type StateTransition {
        fromId: string
        toId: string
        trigger: Trigger?
        guard: string?
        actions: string[]
    }

    type EventTrigger {
        eventName: string
    }

    type AfterTrigger {
        duration: Duration
    }

    type WhenTrigger {
        condition: string
    }

    type Duration {
        value: int
        unit: string
    }

    // Flow Diagram
    type FlowDiagram {
        flowNodes: FlowNode[]
        flowEdges: FlowEdge[]
        swimlanes: Swimlane[]
    }

    type FlowNode {
        id: string
        label: string
        shape: NodeShape
        swimlane: string?
    }

    enum NodeShape {
        TERMINAL
        PROCESS
        DECISION
        IO
        DATABASE
        CONNECTOR
    }

    type FlowEdge {
        fromId: string
        toId: string
        label: string?
    }

    type Swimlane {
        id: string
        name: string
    }

    type Annotation {
        name: string
        value: string?
    }
}
