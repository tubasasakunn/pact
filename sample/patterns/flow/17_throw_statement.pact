// Pattern 17: Throw statement
component ThrowStatementService {
    flow ThrowStatement {
        valid = self.validate(input)
        if valid {
            self.process(input)
        } else {
            throw ValidationError
        }
    }
}
