// 17: Inheritance and Implementation Chains
component Serializable {
    provides SerializableAPI {
        Serialize() -> string
        Deserialize(data: string) -> bool
    }
}

component Cacheable {
    provides CacheableAPI {
        GetCacheKey() -> string
        GetTTL() -> int
    }
}

component Auditable {
    provides AuditableAPI {
        GetAuditLog() -> string[]
    }
}

component BaseEntity {
    type EntityBase {
        +id: string
        +createdAt: string
        +updatedAt: string
        #version: int
    }

    implements Serializable

    provides EntityAPI {
        GetId() -> string
        GetVersion() -> int
    }
}

component VersionedEntity {
    type VersionInfo {
        version: int
        modifiedBy: string
        modifiedAt: string
    }

    extends BaseEntity
    implements Auditable
}

component CachedEntity {
    extends VersionedEntity
    implements Cacheable
}

component UserEntity {
    type User {
        +id: string
        +name: string
        +email: string
        +role: string
    }

    extends CachedEntity

    provides UserEntityAPI {
        GetName() -> string
        GetEmail() -> string
    }
}

component AdminEntity {
    type AdminPrivilege {
        scope: string
        level: int
    }

    extends UserEntity

    provides AdminAPI {
        GetPrivileges() -> AdminPrivilege[]
        HasPermission(scope: string) -> bool
    }
}
