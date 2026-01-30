// Flow Pattern 01: Simple flow
component Service {
    flow SimpleFlow {
        result = self.doWork()
        return result
    }
}
