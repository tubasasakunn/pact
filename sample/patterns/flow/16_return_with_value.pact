// Pattern 16: Return with value
component ReturnWithValueService {
    flow ReturnWithValue {
        data = self.loadData()
        processed = self.process(data)
        return processed
    }
}
