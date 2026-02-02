// Pattern 07: Component with depends on
component UserService {
    type User { id: string }
    depends on Database
}

component Database {
    type Connection { url: string }
}
