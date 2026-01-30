// Pattern 1: Single participant (self calls only)
component SoloService {
    flow ProcessInternally {
        data = self.loadData()
        validated = self.validate(data)
        transformed = self.transform(validated)
        result = self.save(transformed)
        return result
    }
}
