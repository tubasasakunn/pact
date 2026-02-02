// Pattern 09: Component with provides interface
component UserService {
    type User { id: string }

    provides UserAPI {
        GetUser(id: string) -> User
        CreateUser(name: string) -> User
        DeleteUser(id: string) -> bool
    }
}
