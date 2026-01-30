// Pattern 17: Cycle in states (A->B->C->A)
component CyclePattern {
    states RotatingBuffer {
        initial BufferA

        state BufferA { }
        state BufferB { }
        state BufferC { }

        BufferA -> BufferB on rotate
        BufferB -> BufferC on rotate
        BufferC -> BufferA on rotate
    }
}
