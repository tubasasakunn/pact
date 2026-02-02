// User Service - Class Diagram Example
@version("1.0")
@author("pact")
component UserService {
    type User {
        id: string
        name: string
        email: string
        createdAt: string
    }

    type CreateUserRequest {
        name: string
        email: string
    }

    depends on Database
    depends on Logger

    provides UserAPI {
        GetUser(id: string) -> User
        CreateUser(req: CreateUserRequest) -> User
        DeleteUser(id: string) -> bool
    }
}

component Database {
    type Connection {
        host: string
        port: int
        database: string
    }

    provides DatabaseAPI {
        Query(sql: string) -> string
        Execute(sql: string) -> bool
    }
}

component Logger {
    type LogEntry {
        level: string
        message: string
        timestamp: string
    }

    provides LoggerAPI {
        Info(message: string)
        Error(message: string)
        Debug(message: string)
    }
}
