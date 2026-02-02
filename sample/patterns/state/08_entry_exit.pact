// Pattern 8: State with both entry and exit actions
component EntryExit {
    states SessionManager {
        initial Inactive

        state Inactive { }
        state Active {
            entry [startSession, logSessionStart]
            exit [endSession, logSessionEnd]
        }

        Inactive -> Active on login
        Active -> Inactive on logout
    }
}
