// Pattern 22: Variable assignment from method call
component VariableFromMethodService {
    depends on ConfigService
    depends on DatabaseService

    flow VariableFromMethod {
        config = ConfigService.load(appName)
        connectionString = ConfigService.get(dbKey)
        connection = DatabaseService.connect(connectionString)
        query = self.buildQuery(params)
        result = DatabaseService.execute(query)
        formatted = self.formatResult(result)
        return formatted
    }
}
