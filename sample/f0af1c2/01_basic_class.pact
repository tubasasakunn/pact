// 01: Basic Class Diagram - struct, enum, visibility modifiers
@version("1.0")
@author("pact-samples")
component BasicClassExample {
    type User {
        +id: string
        +name: string
        -password: string
        #role: string
        ~internal: int
    }

    type Address {
        street: string
        city: string
        zip: string
        country: string
    }

    enum UserRole {
        Admin
        Editor
        Viewer
        Guest
    }

    provides UserAPI {
        GetUser(id: string) -> User
        ListUsers() -> User[]
    }
}
