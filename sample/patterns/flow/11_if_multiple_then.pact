// Pattern 11: If with multiple statements in then block
component IfMultipleThenService {
    flow IfMultipleThen {
        data = self.getData()
        if isReady {
            result = self.processFirst(data)
            transformed = self.transform(result)
            validated = self.validate(transformed)
            self.save(validated)
            self.notifyComplete()
        }
    }
}
