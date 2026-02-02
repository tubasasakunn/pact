// Pattern 5: Call to another component
component UserController {
    depends on UserRepository

    flow GetUser {
        user = UserRepository.findById(userId)
        return user
    }
}

component UserRepository {
    provides UserDataAccess {
        FindById(id: string)
    }
}
