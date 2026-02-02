// Pattern 04: Component with nullable fields
component UserService {
    type User {
        id: string
        nickname: string?
        avatar: string?
    }
}
