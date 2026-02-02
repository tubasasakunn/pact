// Pattern 4: Initial and final states
component InitialFinal {
    states ProcessLifecycle {
        initial Created
        final Terminated

        state Created { }
        state Running { }
        state Terminated { }

        Created -> Running on start
        Running -> Terminated on finish
    }
}
