// Pattern 13: Multiple transitions from one state
component MultipleFromOne {
    states Router {
        initial Receiving

        state Receiving { }
        state ProcessingA { }
        state ProcessingB { }
        state ProcessingC { }
        state Error { }

        Receiving -> ProcessingA on routeA
        Receiving -> ProcessingB on routeB
        Receiving -> ProcessingC on routeC
        Receiving -> Error on invalidRoute
    }
}
