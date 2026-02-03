// 02: All Relationship Types
component BaseEntity {
    type Entity {
        id: string
        createdAt: string
        updatedAt: string
    }
}

component Auditable {
    provides AuditAPI {
        GetAuditLog() -> string[]
    }
}

component Cache {
    type CacheEntry {
        key: string
        value: string
        ttl: int
    }

    provides CacheAPI {
        Get(key: string) -> string
        Set(key: string, value: string) -> bool
    }
}

component Database {
    type Connection {
        host: string
        port: int
    }

    provides DatabaseAPI {
        Query(sql: string) -> string
        Execute(sql: string) -> bool
    }
}

component Logger {
    provides LogAPI {
        Info(msg: string)
        Error(msg: string)
    }
}

component OrderProcessor {
    type InternalWorker {
        id: string
        status: string
    }

    extends BaseEntity
    implements Auditable
    depends on Database
    depends on Logger
    contains InternalWorker
    aggregates Cache
}
