// 14: Microservices Architecture
@version("1.0")
@description("Microservices with gateway pattern")
component APIGateway {
    type Route {
        path: string
        method: string
        service: string
    }

    depends on AuthService
    depends on UserService
    depends on ProductService
    depends on OrderService
    depends on NotificationService

    provides PublicAPI {
        HandleRequest(path: string, method: string, body: string) -> string throws AuthError
    }

    flow RouteRequest {
        token = self.extractToken(request)
        auth = AuthService.validate(token)
        if auth.valid {
            route = self.matchRoute(request)
            if route.service == "users" {
                result = UserService.handle(route, request)
                return result
            } else {
                if route.service == "products" {
                    result = ProductService.handle(route, request)
                    return result
                } else {
                    result = OrderService.handle(route, request)
                    return result
                }
            }
        } else {
            throw AuthError
        }
    }
}

component AuthService {
    type AuthToken {
        userId: string
        roles: string[]
        expiresAt: string
    }

    depends on UserService

    provides AuthAPI {
        Validate(token: string) -> AuthToken
        Login(email: string, password: string) -> AuthToken throws InvalidCredentials
        Logout(token: string) -> bool
    }

    states SessionLifecycle {
        initial Unauthenticated
        final Expired

        state Unauthenticated {
            entry [clearSession]
        }

        state Authenticated {
            entry [setSessionCookie]
        }

        state Refreshing {
            entry [validateRefreshToken]
        }

        state Expired {
            entry [notifyClient]
        }

        Unauthenticated -> Authenticated on loginSuccess
        Authenticated -> Refreshing on tokenExpiring
        Refreshing -> Authenticated on refreshSuccess
        Refreshing -> Expired on refreshFailed
        Authenticated -> Unauthenticated on logout
        Authenticated -> Expired after 24h
    }
}

component UserService {
    type User {
        +id: string
        +email: string
        +name: string
        -passwordHash: string
        +roles: string[]
    }

    depends on Database

    provides UserAPI {
        GetUser(id: string) -> User
        CreateUser(data: string) -> User throws ValidationError
        UpdateUser(id: string, data: string) -> User
        DeleteUser(id: string) -> bool
    }
}

component ProductService {
    type Product {
        +id: string
        +name: string
        +price: float
        +stock: int
    }

    depends on Database
    depends on CacheService

    provides ProductAPI {
        ListProducts(category: string) -> Product[]
        GetProduct(id: string) -> Product
    }
}

component OrderService {
    type Order {
        +id: string
        +userId: string
        +items: string[]
        +total: float
    }

    depends on Database
    depends on ProductService
    depends on NotificationService

    provides OrderAPI {
        CreateOrder(userId: string, items: string[]) -> Order
        GetOrder(id: string) -> Order
    }
}

component NotificationService {
    provides NotifyAPI {
        SendEmail(to: string, subject: string, body: string) -> bool
        SendPush(userId: string, message: string) -> bool
    }
}

component Database {
    provides DBAPI {
        Query(sql: string) -> string
    }
}

component CacheService {
    provides CacheAPI {
        Get(key: string) -> string
        Set(key: string, value: string, ttl: int) -> bool
    }
}
