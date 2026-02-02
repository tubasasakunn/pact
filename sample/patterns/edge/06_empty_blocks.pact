// Edge case: Empty blocks

component EmptyComponent { }

component WithEmptyType {
    type EmptyType { }
}

component WithEmptyStates {
    states EmptyState {
        initial Start
        state Start { }
    }
}

component WithEmptyFlow {
    provides EmptyAPI {
        EmptyMethod() -> Void
    }

    flow EmptyMethod {
    }
}
