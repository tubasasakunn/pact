// Pattern 20: Multiple components interacting
component Frontend {
    depends on APIGateway

    flow LoadDashboard {
        data = APIGateway.getDashboardData(userId)
        return data
    }
}

component APIGateway {
    depends on UserService
    depends on OrderService
    depends on AnalyticsService

    flow GetDashboardData {
        user = UserService.getProfile(userId)
        orders = OrderService.getRecent(userId)
        stats = AnalyticsService.getUserStats(userId)
        return combined
    }
}

component UserService { }
component OrderService { }
component AnalyticsService { }
