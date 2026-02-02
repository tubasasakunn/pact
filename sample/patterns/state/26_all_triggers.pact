// Pattern 26: State machine using on, when, after triggers
// Demonstrates all three trigger types: event (on), condition (when), time (after)
component AllTriggers {
    states TaskScheduler {
        initial Idle
        final Completed

        state Idle {
            entry [initializeScheduler]
        }
        state Waiting { }
        state Ready { }
        state Processing {
            entry [startProcessing]
            exit [logCompletion]
        }
        state TimedOut { }
        state Completed { }

        // Event trigger (on)
        Idle -> Waiting on taskReceived
        Waiting -> Ready on resourcesAvailable
        Ready -> Processing on startSignal
        Processing -> Completed on taskFinished

        // Condition trigger (when)
        Idle -> Ready when queueNotEmpty
        Waiting -> Processing when priorityTask

        // Time trigger (after)
        Waiting -> TimedOut after 30s
        Processing -> TimedOut after 5m
        TimedOut -> Idle after 10s
    }
}
