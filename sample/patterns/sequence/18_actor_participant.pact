// Pattern 18: Actor participant type
component User { }

component WebApp {
    depends on AuthService

    flow UserLogin {
        credentials = self.getCredentials()
        token = AuthService.authenticate(credentials)
        session = AuthService.createSession(token)
        return session
    }
}

component AuthService { }
