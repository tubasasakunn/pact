// Pattern 19: Database participant type
component UserDB { }

component UserService {
    depends on UserDB

    flow CreateUser {
        existing = UserDB.findByEmail(email)
        if existing == null {
            user = UserDB.insert(userData)
            UserDB.commit()
            return user
        } else {
            throw UserExistsError
        }
    }
}
