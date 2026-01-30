// Pattern 6: Method call on self
component MethodOnSelfService {
    flow MethodOnSelf {
        self.initialize()
        data = self.loadData()
        self.processData(data)
        self.cleanup()
    }
}
