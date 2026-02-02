// Pattern 33: After (time-based) triggers
component SessionManager {
    states SessionLifecycle {
        initial Inactive
        final Expired

        state Inactive { }
        state Active {
            entry [startActivityTimer]
            exit [stopActivityTimer]
        }
        state Warning {
            entry [showTimeoutWarning]
            exit [hideTimeoutWarning]
        }
        state Expired { }

        // Various time units
        Inactive -> Active on login
        Active -> Warning after 25m
        Warning -> Expired after 5m
        Warning -> Active on activity
        Active -> Inactive on logout
    }
}

component RetryHandler {
    states RetryMachine {
        initial Idle
        final Failed
        final Success

        state Idle { }
        state Attempting { }
        state WaitingRetry { }
        state Success { }
        state Failed { }

        Idle -> Attempting on start
        Attempting -> Success on complete
        Attempting -> WaitingRetry on error when retriesRemaining

        // Different retry delays
        WaitingRetry -> Attempting after 1s when attempt1
        WaitingRetry -> Attempting after 5s when attempt2
        WaitingRetry -> Attempting after 30s when attempt3

        WaitingRetry -> Failed after 60s
        Attempting -> Failed on error when noRetriesRemaining
    }
}

component CacheManager {
    states CacheStates {
        initial Empty

        state Empty { }
        state Populated {
            entry [recordPopulateTime]
        }
        state Stale {
            entry [markAsStale]
        }
        state Refreshing { }

        Empty -> Populated on populate
        Populated -> Stale after 1h
        Stale -> Refreshing on access
        Refreshing -> Populated on refreshComplete
        Populated -> Empty on invalidate
        Stale -> Empty on invalidate
    }
}

component HeartbeatMonitor {
    states ConnectionHealth {
        initial Disconnected

        state Disconnected { }
        state Connected {
            entry [startHeartbeat]
            exit [stopHeartbeat]
        }
        state Suspicious {
            entry [increaseHeartbeatFrequency]
        }
        state Reconnecting { }

        Disconnected -> Connected on connect
        Connected -> Suspicious after 30s
        Suspicious -> Disconnected after 10s
        Suspicious -> Connected on heartbeatReceived
        Connected -> Connected on heartbeatReceived
        Disconnected -> Reconnecting on reconnect
        Reconnecting -> Connected on success
        Reconnecting -> Disconnected after 5s
    }
}
