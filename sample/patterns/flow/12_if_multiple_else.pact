// Pattern 12: If with multiple statements in else block
component IfMultipleElseService {
    flow IfMultipleElse {
        status = self.checkStatus()
        if isSuccess {
            self.recordSuccess()
        } else {
            error = self.getError()
            self.logError(error)
            self.notifyAdmin(error)
            self.scheduleRetry()
            self.updateStatus()
        }
    }
}
