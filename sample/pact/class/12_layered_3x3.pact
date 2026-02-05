// Pattern: Layered Architecture 3x3
// 3層3列のレイヤードアーキテクチャパターン

// Presentation Layer
component WebUI {
    depends on ProductService

    provides WebUIAPI {
        RenderPage() -> string
    }
}

component MobileApp {
    depends on CartService

    provides MobileAppAPI {
        ShowScreen() -> string
    }
}

component AdminPanel {
    depends on ReportService

    provides AdminAPI {
        ShowDashboard() -> string
    }
}

// Business Layer
component ProductService {
    depends on ProductRepo

    provides ProductServiceAPI {
        GetProduct(id: string) -> string
        ListProducts() -> string
    }
}

component CartService {
    depends on CartRepo

    provides CartServiceAPI {
        AddItem(item: string) -> bool
        GetCart() -> string
    }
}

component ReportService {
    depends on ReportRepo

    provides ReportServiceAPI {
        GenerateReport() -> string
        GetStats() -> string
    }
}

// Data Layer
component ProductRepo {
    provides ProductRepoAPI {
        FindProduct(id: string) -> string
        SaveProduct(data: string) -> bool
    }
}

component CartRepo {
    provides CartRepoAPI {
        GetCart(userId: string) -> string
        UpdateCart(data: string) -> bool
    }
}

component ReportRepo {
    provides ReportRepoAPI {
        QueryData(query: string) -> string
        AggregateData() -> string
    }
}
