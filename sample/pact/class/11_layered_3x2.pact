// Pattern: Layered Architecture 3x2
// 3層2列のレイヤードアーキテクチャパターン

// Presentation Layer
component WebController {
    type RequestData {
        path: string
        method: string
    }

    depends on UserService

    provides WebAPI {
        HandleRequest() -> string
    }
}

component APIController {
    type APIRequestData {
        endpoint: string
        payload: string
    }

    depends on OrderService

    provides RESTAPI {
        HandleAPI() -> string
    }
}

// Business Layer
component UserService {
    type UserData {
        id: string
        name: string
    }

    depends on UserRepository

    provides UserServiceAPI {
        GetUser(id: string) -> string
        CreateUser(data: string) -> bool
    }
}

component OrderService {
    type OrderData {
        id: string
        total: float
    }

    depends on OrderRepository

    provides OrderServiceAPI {
        GetOrder(id: string) -> string
        CreateOrder(data: string) -> bool
    }
}

// Data Layer
component UserRepository {
    type UserEntity {
        id: string
        data: string
    }

    provides UserRepoAPI {
        Find(id: string) -> string
        Save(entity: string) -> bool
    }
}

component OrderRepository {
    type OrderEntity {
        id: string
        data: string
    }

    provides OrderRepoAPI {
        Find(id: string) -> string
        Save(entity: string) -> bool
    }
}
