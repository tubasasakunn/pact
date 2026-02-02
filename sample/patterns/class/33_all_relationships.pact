// Pattern 33: All relationship types combined
component BaseEntity {
    type EntityId {
        value: string
    }
}

component Repository {
    provides IRepository {
        findById(id: string) -> Entity
        save(entity: Entity) -> bool
        delete(id: string) -> bool
    }
}

component Logger {
    provides ILogger {
        log(message: string)
        error(message: string, error: Error)
    }
}

component Cache {
    type CacheEntry {
        key: string
        value: string
        ttl: int
    }
}

// Component using all relationship types
component UserService {
    // extends - inheritance
    extends BaseEntity

    // implements - interface implementation
    implements Repository

    // depends on - dependency injection
    depends on Database : IDatabase as db
    depends on Logger : ILogger as logger

    // contains - composition (owns the lifecycle)
    contains UserValidator
    contains UserMapper

    // aggregates - aggregation (shared lifecycle)
    aggregates Cache
    aggregates MetricsCollector

    type User {
        id: string
        name: string
        email: string
    }

    provides UserAPI {
        createUser(name: string, email: string) -> User
        getUser(id: string) -> User
        updateUser(user: User) -> bool
        deleteUser(id: string) -> bool
    }

    flow CreateUser {
        validated = UserValidator.validate(name, email)
        if validated {
            user = UserMapper.toEntity(name, email)
            saved = db.save(user)
            Cache.invalidate(userCacheKey)
            MetricsCollector.recordCreate()
            logger.log(userCreatedMessage)
            return user
        } else {
            logger.error(validationFailedMessage, validationError)
            throw ValidationError
        }
    }
}

// Another example with multiple inheritance-like patterns
component AdvancedService {
    extends BaseEntity
    implements Repository
    implements Logger

    depends on ExternalAPI : IAPI as api
    depends on ConfigService : IConfig as config

    contains InternalProcessor
    contains DataTransformer

    aggregates SharedCache
    aggregates ConnectionPool
}
