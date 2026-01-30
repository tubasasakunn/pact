// Pattern 12: Self-transition (state to itself)
component SelfTransition {
    states Counter {
        initial Counting

        state Counting { }
        state Done { }

        Counting -> Counting on increment
        Counting -> Done on finish
    }
}
