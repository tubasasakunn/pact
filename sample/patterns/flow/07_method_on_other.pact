// Pattern 7: Method call on other component
component MethodOnOtherService {
    depends on DatabaseService
    depends on CacheService

    flow MethodOnOther {
        cached = CacheService.get(key)
        result = DatabaseService.query(params)
        CacheService.set(key, result)
    }
}
