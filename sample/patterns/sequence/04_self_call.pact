// Pattern 4: Self-call (self.method())
component Calculator {
    flow CalculateComplex {
        sum = self.add(a, b)
        product = self.multiply(sum, c)
        result = self.format(product)
        return result
    }
}
