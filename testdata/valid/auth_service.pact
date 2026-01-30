// Real-world example: Authentication Service
@version("1.0")
@module("auth")

type Credentials {
	username: string
	password: string
}

type Token {
	accessToken: string
	refreshToken: string
	expiresIn: int
	tokenType: string
}

type AuthResult {
	success: bool
	token: Token?
	error: string?
}

type UserSession {
	sessionId: string
	userId: string
	createdAt: timestamp
	expiresAt: timestamp
	ipAddress: string
	userAgent: string
}

interface Authenticator {
	method Authenticate(credentials: Credentials): AuthResult
	method ValidateToken(token: string): bool
	method RefreshToken(refreshToken: string): Token?
	method Logout(sessionId: string): void
}

interface TokenService {
	method GenerateToken(userId: string): Token
	method ValidateToken(token: string): bool
	method RevokeToken(token: string): void
}

@service
component AuthenticationService {
	implements Authenticator

	private userRepository: UserRepository
	private tokenService: TokenService
	private sessionStore: SessionStore
	private passwordHasher: PasswordHasher
	private eventPublisher: EventPublisher

	relation UserRepository: uses
	relation TokenService: uses
	relation SessionStore: uses

	@log(level: "info")
	@metric(name: "auth_attempts")
	public method Authenticate(credentials: Credentials): AuthResult

	@cache(ttl: 60)
	public method ValidateToken(token: string): bool

	public method RefreshToken(refreshToken: string): Token?

	@audit(action: "logout")
	public method Logout(sessionId: string): void

	private method validateCredentials(credentials: Credentials): User?
	private method createSession(user: User): UserSession

	states AuthState {
		initial -> Unauthenticated

		Unauthenticated -> Authenticating: startAuth
		Authenticating -> Authenticated: success
		Authenticating -> Failed: failure
		Authenticating -> Locked: tooManyAttempts

		Authenticated -> Refreshing: refreshToken
		Authenticated -> Unauthenticated: logout
		Authenticated -> Expired: tokenExpired

		Refreshing -> Authenticated: refreshSuccess
		Refreshing -> Unauthenticated: refreshFailed

		Expired -> Unauthenticated: forceLogout
		Expired -> Authenticated: refreshSuccess

		Failed -> Unauthenticated: reset
		Failed -> Locked: maxAttempts

		Locked -> Unauthenticated: unlock

		// Terminal states handled implicitly
	}

	flow Login {
		start: "Receive Login Request"
		validateInput: "Validate Input"
		if inputValid {
			checkRateLimit: "Check Rate Limit"
			if notLimited {
				lookupUser: "Lookup User"
				if userFound {
					verifyPassword: "Verify Password"
					if passwordValid {
						checkMFA: "Check MFA Required"
						if mfaRequired {
							sendMFACode: "Send MFA Code"
							waitMFA: "Wait for MFA"
							if mfaValid {
								createSession: "Create Session"
								generateToken: "Generate Token"
								logSuccess: "Log Successful Login"
								returnToken: "Return Token"
							} else {
								mfaFailed: "MFA Verification Failed"
								incrementAttempts: "Increment Failed Attempts"
							}
						} else {
							createSession: "Create Session"
							generateToken: "Generate Token"
							logSuccess: "Log Successful Login"
							returnToken: "Return Token"
						}
					} else {
						passwordFailed: "Password Verification Failed"
						incrementAttempts: "Increment Failed Attempts"
						checkLockout: "Check Account Lockout"
						if shouldLock {
							lockAccount: "Lock Account"
						}
					}
				} else {
					userNotFound: "User Not Found"
				}
			} else {
				rateLimited: "Return Rate Limit Error"
			}
		} else {
			invalidInput: "Return Validation Error"
		}
		end: "Complete"
	}
}

@repository
component SessionStore {
	private redis: RedisClient
	private config: SessionConfig

	method CreateSession(session: UserSession): void
	method GetSession(sessionId: string): UserSession?
	method DeleteSession(sessionId: string): void
	method GetUserSessions(userId: string): []UserSession
	method DeleteUserSessions(userId: string): void
}

@service
component JWTTokenService {
	implements TokenService

	private secretKey: string
	private issuer: string
	private accessTokenTTL: int
	private refreshTokenTTL: int

	method GenerateToken(userId: string): Token
	method ValidateToken(token: string): bool
	method RevokeToken(token: string): void
	method DecodeToken(token: string): map[string]any
}
