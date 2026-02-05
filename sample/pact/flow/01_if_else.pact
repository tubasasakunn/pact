// Pattern: If-Else Flow
// 条件分岐を含むフローパターン

component AuthService {
    type AuthResult {
        success: bool
        token: string
        message: string
    }

    flow Authenticate {
        credentials = self.getCredentials(request)
        user = self.findUser(credentials.username)

        if user {
            valid = self.validatePassword(user, credentials.password)
            if valid {
                token = self.generateToken(user)
                result = self.createAuthResult(token)
                return result
            } else {
                throw InvalidPasswordError
            }
        } else {
            throw UserNotFoundError
        }
    }
}
