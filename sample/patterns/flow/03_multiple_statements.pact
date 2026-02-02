// Pattern 3: Multiple sequential statements
component MultipleStatementsService {
    flow MultipleStatements {
        first = self.getFirst()
        second = self.getSecond()
        third = self.getThird()
    }
}
