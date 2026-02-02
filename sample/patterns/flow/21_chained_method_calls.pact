// Pattern 21: Chained method calls
component ChainedMethodService {
    depends on DataService
    depends on TransformService
    depends on ValidationService

    flow ChainedMethods {
        raw = DataService.fetch(source)
        transformed = TransformService.apply(raw)
        validated = ValidationService.check(transformed)
        result = self.finalize(validated)
        return result
    }
}
