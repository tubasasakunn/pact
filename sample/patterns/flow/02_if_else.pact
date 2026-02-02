// Flow Pattern 02: If-else flow
component Service {
    flow ConditionalFlow {
        valid = self.validate(input)
        if valid {
            result = self.process(input)
            return result
        } else {
            throw ValidationError
        }
    }
}
