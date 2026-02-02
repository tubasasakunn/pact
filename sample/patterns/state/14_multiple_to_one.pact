// Pattern 14: Multiple transitions to one state
component MultipleToOne {
    states Collector {
        initial Idle
        final Complete

        state Idle { }
        state SourceA { }
        state SourceB { }
        state SourceC { }
        state Complete { }

        Idle -> SourceA on startA
        Idle -> SourceB on startB
        Idle -> SourceC on startC
        SourceA -> Complete on finish
        SourceB -> Complete on finish
        SourceC -> Complete on finish
    }
}
