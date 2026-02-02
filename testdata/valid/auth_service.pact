// Real-world example: Authentication Service
@version("1.0")
component AuthenticationService {
	type Credentials {
		username: string
		password: string
	}

	type Token {
		accessToken: string
		refreshToken: string
		expiresIn: int
	}

	type AuthResult {
		success: bool
		token: Token?
		error: string?
	}

	depends on UserRepository
	depends on TokenService
	depends on SessionStore

	provides AuthAPI {
		Authenticate(credentials: Credentials) -> AuthResult
		ValidateToken(token: string) -> bool
		RefreshToken(refreshToken: string) -> Token
		Logout(sessionId: string)
	}

	states AuthState {
		initial Unauthenticated
		final Locked

		state Unauthenticated { }
		state Authenticating { }
		state Authenticated { }
		state Failed { }
		state Locked { }

		Unauthenticated -> Authenticating on startAuth
		Authenticating -> Authenticated on success
		Authenticating -> Failed on failure
		Authenticating -> Locked on tooManyAttempts
		Authenticated -> Unauthenticated on logout
		Failed -> Unauthenticated on reset
		Failed -> Locked on maxAttempts
	}

	flow Login {
		valid = self.validateInput(credentials)
		if valid {
			limited = self.checkRateLimit(credentials)
			if limited {
				throw RateLimitError
			}
			user = UserRepository.findByUsername(credentials.username)
			if user {
				passwordValid = self.verifyPassword(credentials.password, user.hash)
				if passwordValid {
					session = self.createSession(user)
					token = TokenService.generateToken(user)
					self.logSuccess(user)
					return token
				} else {
					self.incrementAttempts(user)
					throw InvalidCredentials
				}
			} else {
				throw UserNotFound
			}
		} else {
			throw ValidationError
		}
	}
}
