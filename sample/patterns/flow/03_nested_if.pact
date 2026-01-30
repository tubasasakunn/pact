// Flow Pattern 03: Nested if conditions
component Service {
    flow NestedConditions {
        a = self.checkA()
        if a {
            b = self.checkB()
            if b {
                c = self.checkC()
                if c {
                    return self.success()
                } else {
                    throw ErrorC
                }
            } else {
                throw ErrorB
            }
        } else {
            throw ErrorA
        }
    }
}
