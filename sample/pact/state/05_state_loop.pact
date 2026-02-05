// Pattern: State Loop
// ループを含む状態遷移パターン

component RetryMechanism {
    type RetryData {
        maxRetries: int
        currentAttempt: int
        delayMs: int
    }

    states RetryState {
        initial Idle
        final Success
        final MaxRetriesExceeded

        state Idle {
            entry [resetCounter]
        }

        state Attempting {
            entry [executeRequest]
        }

        state Waiting {
            entry [incrementCounter, calculateDelay]
        }

        state Success {
            entry [logSuccess]
        }

        state MaxRetriesExceeded {
            entry [logFailure, alertAdmin]
        }

        Idle -> Attempting on start
        Attempting -> Success on success
        Attempting -> Waiting on failure
        Waiting -> Attempting on retryTimer
        Waiting -> MaxRetriesExceeded on maxRetriesReached
    }
}
