// 15: State Machine with various duration triggers
component TaskScheduler {
    type Task {
        id: string
        name: string
        priority: int
        deadline: string
    }

    states TaskLifecycle {
        initial Queued
        final Completed
        final TimedOut
        final Failed

        state Queued {
            entry [addToQueue]
        }

        state Scheduled {
            entry [reserveResources]
        }

        state Running {
            entry [startExecution, startMonitor]
            exit [stopMonitor]
        }

        state Paused {
            entry [saveCheckpoint]
            exit [restoreCheckpoint]
        }

        state Retrying {
            entry [incrementRetryCount]
        }

        state Completed {
            entry [releaseResources, reportSuccess]
        }

        state TimedOut {
            entry [releaseResources, reportTimeout]
        }

        state Failed {
            entry [releaseResources, reportFailure]
        }

        Queued -> Scheduled on schedule
        Scheduled -> Running on execute
        Running -> Completed on success
        Running -> Paused on pause
        Paused -> Running on resume
        Running -> Retrying on error when retriesRemaining
        Running -> Failed on error when noRetriesLeft
        Retrying -> Running on retryReady
        Paused -> TimedOut after 1h
        Running -> TimedOut after 30m
        Queued -> TimedOut after 7d
        Scheduled -> Queued on reschedule
    }
}
