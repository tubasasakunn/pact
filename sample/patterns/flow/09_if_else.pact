// Pattern 9: If-else condition
component IfElseService {
    flow IfElse {
        value = self.getValue()
        if isValid {
            self.handleValid()
        } else {
            self.handleInvalid()
        }
    }
}
