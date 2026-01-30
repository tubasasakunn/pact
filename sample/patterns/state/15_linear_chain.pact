// Pattern 15: Linear state chain (A->B->C->D)
component LinearChain {
    states Pipeline {
        initial StageA
        final StageD

        state StageA { }
        state StageB { }
        state StageC { }
        state StageD { }

        StageA -> StageB on next
        StageB -> StageC on next
        StageC -> StageD on next
    }
}
