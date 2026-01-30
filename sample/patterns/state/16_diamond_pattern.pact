// Pattern 16: Diamond pattern (A->B, A->C, B->D, C->D)
component DiamondPattern {
    states DecisionTree {
        initial Start
        final End

        state Start { }
        state PathLeft { }
        state PathRight { }
        state End { }

        Start -> PathLeft on goLeft
        Start -> PathRight on goRight
        PathLeft -> End on merge
        PathRight -> End on merge
    }
}
