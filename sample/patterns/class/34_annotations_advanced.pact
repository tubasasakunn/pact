// Pattern 34: Advanced annotation patterns
@version("2.0.0")
@author("platform-team")
@deprecated("Use NewService instead")
@since("1.5.0")
component AnnotatedService {

    @description("User data transfer object")
    @serializable
    type UserDTO {
        @required
        @minLength("1")
        @maxLength("100")
        id: string

        @email
        @unique
        email: string

        @nullable
        @default("anonymous")
        displayName: string?

        @range("0-150")
        age: int?

        @pattern("[A-Z]{2}")
        countryCode: string
    }

    @entity
    @table("users")
    type UserEntity {
        @primaryKey
        @autoGenerate
        id: string

        @column("user_email")
        @index
        email: string

        @column("created_at")
        @timestamp
        createdAt: string

        @oneToMany("Order")
        orders: Order[]
    }

    @enumType
    @description("User status values")
    enum UserStatus {
        ACTIVE
        INACTIVE
        SUSPENDED
        DELETED
    }

    @api
    @authenticated
    @rateLimit("100/1m")
    provides UserAPI {
        @get("/users/{id}")
        @cache("ttl:300")
        getUser(id: string) -> UserDTO

        @post("/users")
        @validate
        @transaction
        createUser(user: UserDTO) -> UserDTO

        @put("/users/{id}")
        @idempotent
        updateUser(id: string, user: UserDTO) -> UserDTO

        @delete("/users/{id}")
        @adminOnly
        deleteUser(id: string) -> bool
    }

    @internal
    @retry("maxAttempts:3,delay:1000")
    flow ProcessUser {
        @log("info")
        user = self.fetchUser(userId)

        @measure("validation_time")
        validated = self.validateUser(user)

        @transactional
        if validated {
            @audit("user_update")
            result = self.saveUser(user)
            return result
        } else {
            @alert("warning")
            throw ValidationError
        }
    }
}

@service
@singleton
@injectable
component SingletonService {
    @inject
    depends on Database : IDatabase as db

    @lazy
    depends on HeavyService : IHeavy as heavy

    @config("app.timeout")
    type Config {
        timeout: int
        retries: int
    }
}
