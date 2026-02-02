// Pattern 23: Flow with 10+ statements
component ManyStatementsService {
    depends on LogService
    depends on MetricsService
    depends on CacheService

    flow ManyStatements {
        LogService.info(startMessage)
        MetricsService.startTimer(operationId)
        cached = CacheService.get(cacheKey)
        config = self.loadConfig()
        validated = self.validateInput(input)
        prepared = self.prepareData(validated)
        enriched = self.enrichData(prepared)
        transformed = self.transformData(enriched)
        processed = self.processData(transformed)
        result = self.finalizeData(processed)
        CacheService.set(cacheKey, result)
        MetricsService.stopTimer(operationId)
        LogService.info(endMessage)
        return result
    }
}
