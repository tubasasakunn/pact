// Pattern 13: Multiple calls to same participant
component ReportGenerator {
    depends on Database

    flow GenerateReport {
        users = Database.queryUsers()
        orders = Database.queryOrders()
        products = Database.queryProducts()
        stats = Database.queryStatistics()
        logs = Database.queryLogs()
        return true
    }
}

component Database { }
