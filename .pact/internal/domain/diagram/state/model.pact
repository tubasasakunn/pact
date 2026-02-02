// State diagram domain model
component StateDiagram {
    // StateType represents the type of state
    enum StateType {
        InitialState
        FinalState
        AtomicState
        CompoundState
        ParallelState
    }

    // Duration represents a time duration
    type Duration {
        value: int
        unit: string
    }

    // EventTrigger represents an event-based trigger
    type EventTrigger {
        event: string
    }

    // AfterTrigger represents a time-based trigger
    type AfterTrigger {
        duration: Duration
    }

    // WhenTrigger represents a condition-based trigger
    type WhenTrigger {
        condition: string
    }

    // Transition represents a state transition
    type Transition {
        fromState: string
        toState: string
        guard: string
        actions: string[]
    }

    // Region represents a region in parallel states
    type Region {
        name: string
    }

    // State represents a state in the diagram
    type State {
        id: string
        name: string
        stateType: StateType
        entryActions: string[]
        exitActions: string[]
    }

    // Diagram represents a state diagram
    type Diagram {
        statesBlocks: State[]
        transitions: Transition[]
    }
}
