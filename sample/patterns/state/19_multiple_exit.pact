// Pattern 19: State with multiple exit actions
component MultipleExit {
    states Shutdown {
        initial Running

        state Running {
            exit [saveState, closeConnections, flushBuffers, notifyClients]
        }
        state Stopped { }

        Running -> Stopped on shutdown
    }
}
