// AST States - State machine definitions for the Abstract Syntax Tree
// This file defines state machines, states, transitions, and triggers

component ASTStates {
    // StatesDecl represents a state machine definition
    type StatesDecl {
        pos: Position
        name: String
        annotations: AnnotationDecl[]
        initialState: String
        finals: String[]
        statesBlocks: StateDecl[]
        transitions: TransitionDecl[]
        parallels: ParallelDecl[]
    }

    // StateDecl represents a state definition
    type StateDecl {
        pos: Position
        name: String
        annotations: AnnotationDecl[]
        entryActions: String[]
        exitActions: String[]
        // For hierarchical states
        initialState: String?
        statesBlocks: StateDecl[]
        transitions: TransitionDecl[]
    }

    // TransitionDecl represents a state transition definition
    type TransitionDecl {
        pos: Position
        from: String
        to: String
        trigger: Trigger
        guard: Expr?
        actions: String[]
    }

    // Trigger is the base interface for all trigger types
    // All trigger types implement triggerNode() and GetPos()

    // EventTrigger represents an event-based trigger
    type EventTrigger {
        pos: Position
        event: String
    }

    // AfterTrigger represents a time-based trigger
    type AfterTrigger {
        pos: Position
        duration: Duration
    }

    // Duration represents a time duration
    type Duration {
        value: Int
        unit: String                       // "ms", "s", "m", "h", "d"
    }

    // WhenTrigger represents a condition-based trigger
    type WhenTrigger {
        pos: Position
        condition: Expr
    }

    // ParallelDecl represents a parallel state definition
    type ParallelDecl {
        pos: Position
        name: String
        annotations: AnnotationDecl[]
        regions: RegionDecl[]
    }

    // RegionDecl represents a region within a parallel state
    type RegionDecl {
        pos: Position
        name: String
        initialState: String
        statesBlocks: StateDecl[]
        transitions: TransitionDecl[]
    }
}
