// Pattern: Interface Implementation with 3 implementers
// 1つのインターフェースを3つのクラスが実装するパターン

component Storage {
    provides StorageAPI {
        Save(key: string, value: string) -> bool
        Load(key: string) -> string
        Delete(key: string) -> bool
    }
}

component FileStorage {
    type FileStorageConfig {
        basePath: string
        extension: string
    }

    implements Storage

    provides FileStorageAPI {
        SetBasePath(path: string) -> bool
    }
}

component MemoryStorage {
    type MemoryStorageConfig {
        maxSize: int
        evictionPolicy: string
    }

    implements Storage

    provides MemoryStorageAPI {
        Clear() -> bool
        GetSize() -> int
    }
}

component DatabaseStorage {
    type DatabaseStorageConfig {
        connectionString: string
        tableName: string
    }

    implements Storage

    provides DatabaseStorageAPI {
        Connect() -> bool
        Disconnect() -> bool
    }
}
