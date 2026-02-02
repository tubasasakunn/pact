// Pattern 19: Try-catch like pattern (using if-else for error handling)
component TryCatchPatternService {
    depends on ExternalService

    flow TryCatchPattern {
        result = ExternalService.call(request)
        if success {
            processed = self.processResult(result)
            self.saveResult(processed)
            return processed
        } else {
            self.logError(result)
            self.notifyFailure()
            throw ServiceError
        }
    }
}
