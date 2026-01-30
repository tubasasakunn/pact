// State Pattern 01: Simple state machine
component Order {
    states OrderStatus {
        initial Pending
        final Done

        state Pending { }
        state Done { }

        Pending -> Done on complete
    }
}
