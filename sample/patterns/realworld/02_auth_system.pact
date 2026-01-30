// Real-world Pattern 2: Authentication System
// Demonstrates: Security flows, Token management, Session states

// ============ TYPE DEFINITIONS ============

component User {
	type Credentials {
		username: string
		password: string
	}

	type UserAccount {
		id: string
		username: string
		email: string
		passwordHash: string
		roles: string[]
		isActive: bool
		lastLogin: string?
		failedAttempts: int
		lockedUntil: string?
	}

	provides UserRepository {
		FindByUsername(username: string) -> UserAccount?
		FindById(userId: string) -> UserAccount
		UpdateLastLogin(userId: string, timestamp: string)
		IncrementFailedAttempts(userId: string)
		ResetFailedAttempts(userId: string)
		LockAccount(userId: string, duration: int)
	}
}

component Token {
	type AccessToken {
		token: string
		userId: string
		issuedAt: string
		expiresAt: string
		scopes: string[]
	}

	type RefreshToken {
		token: string
		userId: string
		issuedAt: string
		expiresAt: string
		deviceId: string?
	}

	type TokenPair {
		accessToken: AccessToken
		refreshToken: RefreshToken
	}

	depends on User

	provides TokenService {
		GenerateAccessToken(userId: string, scopes: string[]) -> AccessToken
		GenerateRefreshToken(userId: string, deviceId: string?) -> RefreshToken
		GenerateTokenPair(userId: string) -> TokenPair
		ValidateAccessToken(token: string) -> bool
		ValidateRefreshToken(token: string) -> bool
		RevokeToken(token: string)
		GetTokenInfo(token: string) -> AccessToken?
	}

	requires CryptoService {
		Sign(payload: string, secret: string) -> string
		Verify(token: string, secret: string) -> bool
		Hash(data: string) -> string
	}
}

component Session {
	type SessionData {
		id: string
		userId: string
		deviceInfo: string?
		ipAddress: string?
		userAgent: string?
		createdAt: string
		lastActivityAt: string
		expiresAt: string
		status: string
	}

	depends on User
	depends on Token

	provides SessionManager {
		CreateSession(userId: string, deviceInfo: string?) -> SessionData
		GetSession(sessionId: string) -> SessionData?
		GetUserSessions(userId: string) -> SessionData[]
		UpdateActivity(sessionId: string)
		InvalidateSession(sessionId: string)
		InvalidateAllSessions(userId: string)
	}

	// Session lifecycle states
	states SessionLifecycle {
		initial Active

		state Active {
			entry [recordSessionStart]
		}
		state Idle {
			entry [warnUser]
		}
		state Expired {
			entry [cleanupResources]
		}
		state Revoked {
			entry [notifyUser]
			exit [logSecurityEvent]
		}

		Active -> Idle on inactivityTimeout
		Active -> Expired on sessionTimeout
		Active -> Revoked on adminRevoked
		Active -> Revoked on userLogout
		Active -> Revoked on securityThreat
		Idle -> Active on userActivity
		Idle -> Expired on idleTimeout
		Idle -> Revoked on userLogout
	}
}

// ============ AUTHENTICATION FLOWS ============

component AuthController {
	depends on User
	depends on Token
	depends on Session
	depends on AuditLogger

	provides AuthAPI {
		Login(username: string, password: string) -> TokenPair
		Logout(sessionId: string)
		RefreshToken(refreshToken: string) -> TokenPair
		ValidateToken(accessToken: string) -> bool
		GetCurrentUser(accessToken: string) -> UserAccount
	}

	// Login flow with security measures
	flow UserLogin {
		user = User.FindByUsername(username)
		if userExists {
			if accountLocked {
				AuditLogger.LogBlockedAttempt(username)
				throw AccountLockedError
			}
			valid = self.verifyPassword(password, user.passwordHash)
			if valid {
				User.ResetFailedAttempts(user.id)
				User.UpdateLastLogin(user.id, timestamp)
				tokenPair = Token.GenerateTokenPair(user.id)
				session = Session.CreateSession(user.id, deviceInfo)
				AuditLogger.LogSuccessfulLogin(user.id, session.id)
				return tokenPair
			} else {
				User.IncrementFailedAttempts(user.id)
				attempts = self.getFailedAttempts(user.id)
				if tooManyAttempts {
					User.LockAccount(user.id, lockDuration)
					AuditLogger.LogAccountLocked(user.id)
				}
				AuditLogger.LogFailedLogin(username)
				throw InvalidCredentialsError
			}
		} else {
			AuditLogger.LogFailedLogin(username)
			throw InvalidCredentialsError
		}
	}

	// Token refresh flow
	flow RefreshTokenFlow {
		valid = Token.ValidateRefreshToken(refreshToken)
		if valid {
			tokenInfo = Token.GetTokenInfo(refreshToken)
			user = User.FindById(tokenInfo.userId)
			if userActive {
				Token.RevokeToken(refreshToken)
				newTokenPair = Token.GenerateTokenPair(user.id)
				Session.UpdateActivity(sessionId)
				AuditLogger.LogTokenRefresh(user.id)
				return newTokenPair
			} else {
				Token.RevokeToken(refreshToken)
				AuditLogger.LogRefreshDenied(tokenInfo.userId)
				throw AccountDeactivatedError
			}
		} else {
			AuditLogger.LogInvalidRefreshToken(refreshToken)
			throw InvalidTokenError
		}
	}

	// Logout flow
	flow UserLogout {
		session = Session.GetSession(sessionId)
		if sessionExists {
			Token.RevokeToken(session.accessToken)
			Token.RevokeToken(session.refreshToken)
			Session.InvalidateSession(sessionId)
			AuditLogger.LogLogout(session.userId)
			return success
		} else {
			throw SessionNotFoundError
		}
	}
}

component AuditLogger {
	type AuditEntry {
		id: string
		timestamp: string
		action: string
		userId: string?
		ipAddress: string?
		details: string?
		severity: string
	}

	provides AuditService {
		LogSuccessfulLogin(userId: string, sessionId: string)
		LogFailedLogin(username: string)
		LogBlockedAttempt(username: string)
		LogAccountLocked(userId: string)
		LogLogout(userId: string)
		LogTokenRefresh(userId: string)
		LogRefreshDenied(userId: string)
		LogInvalidRefreshToken(token: string)
		LogSecurityEvent(event: string, severity: string)
	}
}
