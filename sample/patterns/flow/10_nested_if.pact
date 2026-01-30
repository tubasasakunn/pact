// Pattern 10: Nested if conditions
component NestedIfService {
    flow NestedIf {
        status = self.getStatus()
        if isActive {
            permissions = self.getPermissions()
            if hasAccess {
                data = self.loadData()
                if isValid {
                    self.processData(data)
                } else {
                    self.markInvalid()
                }
            } else {
                self.denyAccess()
            }
        } else {
            self.activate()
        }
    }
}
