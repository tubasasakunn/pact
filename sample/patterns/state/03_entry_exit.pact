// State Pattern 03: Entry and exit actions
component Order {
    states OrderStatus {
        initial Pending
        final Done

        state Pending {
            entry [logEntry]
        }

        state Processing {
            entry [startTimer, notifyTeam]
            exit [stopTimer, logExit]
        }

        state Done {
            entry [sendNotification]
        }

        Pending -> Processing on start
        Processing -> Done on complete
    }
}
