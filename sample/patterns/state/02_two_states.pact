// Pattern 2: Two states with transition
component TwoStates {
    states BasicTransition {
        initial Off

        state Off { }
        state On { }

        Off -> On on toggle
    }
}
