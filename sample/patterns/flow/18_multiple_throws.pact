// Pattern 18: Multiple throw statements (different errors)
component MultipleThrowsService {
    flow MultipleThrows {
        authenticated = self.authenticate(user)
        if authenticated {
            authorized = self.authorize(user, resource)
            if authorized {
                data = self.getData(resource)
                if data {
                    return data
                } else {
                    throw NotFoundError
                }
            } else {
                throw AuthorizationError
            }
        } else {
            throw AuthenticationError
        }
    }
}
