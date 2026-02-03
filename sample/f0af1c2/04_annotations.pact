// 04: Various Annotations
@version("2.0.0")
@author("backend-team")
@description("Annotated service example")
component AnnotatedService {
    @entity
    @table("users")
    type User {
        @primaryKey
        +id: string
        @index
        +email: string
        +name: string
        @deprecated("Use fullName instead")
        +displayName: string
    }

    @serializable
    type Config {
        dbHost: string
        dbPort: int
        cacheEnabled: bool
    }

    enum Priority {
        Low
        Medium
        High
        Critical
    }

    depends on Database
    depends on CacheService

    @authenticated
    @rateLimit(max: "100", window: "60s")
    provides RestAPI {
        @get("/api/users")
        ListUsers() -> User[]

        @get("/api/users/{id}")
        GetUser(id: string) -> User

        @post("/api/users")
        @validate
        CreateUser(name: string, email: string) -> User throws ValidationError

        @put("/api/users/{id}")
        @transaction
        UpdateUser(id: string, data: string) -> User

        @delete("/api/users/{id}")
        @adminOnly
        DeleteUser(id: string) -> bool
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
    }
}
