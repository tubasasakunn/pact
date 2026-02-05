// Pattern: Interface Implementation with 4 implementers
// 1つのインターフェースを4つのクラスが実装するパターン

component Logger {
    provides LoggerAPI {
        Log(level: string, message: string) -> bool
        SetLevel(level: string) -> bool
    }
}

component ConsoleLogger {
    type ConsoleConfig {
        colorEnabled: bool
        timestampFormat: string
    }

    implements Logger

    provides ConsoleLoggerAPI {
        EnableColors(enabled: bool) -> bool
    }
}

component FileLogger {
    type FileConfig {
        filePath: string
        maxFileSize: int
        rotationCount: int
    }

    implements Logger

    provides FileLoggerAPI {
        Rotate() -> bool
    }
}

component NetworkLogger {
    type NetworkConfig {
        serverUrl: string
        port: int
        protocol: string
    }

    implements Logger

    provides NetworkLoggerAPI {
        Connect() -> bool
        Reconnect() -> bool
    }
}

component DatabaseLogger {
    type DatabaseConfig {
        connectionString: string
        tableName: string
        batchSize: int
    }

    implements Logger

    provides DatabaseLoggerAPI {
        Flush() -> bool
        GetLogCount() -> int
    }
}
