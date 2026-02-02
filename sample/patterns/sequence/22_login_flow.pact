// Pattern 22: Real-world scenario - Login flow
component User { }

component LoginController {
    depends on AuthService
    depends on SessionManager
    depends on AuditLogger

    flow UserLogin {
        credentials = self.extractCredentials(request)
        valid = AuthService.validateCredentials(credentials)
        if valid {
            user = AuthService.getUser(credentials.username)
            token = AuthService.generateToken(user)
            SessionManager.createSession(user, token)
            AuditLogger.logLogin(user)
            return token
        } else {
            AuditLogger.logFailedLogin(credentials.username)
            throw AuthenticationError
        }
    }
}

component AuthService { }
component SessionManager { }
component AuditLogger { }
