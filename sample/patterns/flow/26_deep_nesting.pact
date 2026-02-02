// Pattern 26: 5-level nested if statements with complex logic
component DeepNestingService {
    depends on SecurityService
    depends on ValidationService
    depends on ConfigService

    flow DeepNesting {
        // Level 1: Security check
        securityContext = SecurityService.getContext(request)
        if securityEnabled {
            // Level 2: Authentication
            authResult = SecurityService.authenticate(credentials)
            if authenticated {
                // Level 3: Authorization
                permissions = SecurityService.getPermissions(user)
                if authorized {
                    // Level 4: Validation
                    validationResult = ValidationService.validate(data)
                    if validationPassed {
                        // Level 5: Configuration
                        config = ConfigService.getConfig(operation)
                        if configValid {
                            result = self.executeOperation(data, config)
                            auditLog = self.logSuccess(result)
                            return result
                        } else {
                            self.logConfigError(config)
                            throw ConfigurationError
                        }
                    } else {
                        errors = ValidationService.getErrors(validationResult)
                        self.logValidationErrors(errors)
                        throw ValidationError
                    }
                } else {
                    self.logAuthorizationFailure(user, resource)
                    throw AuthorizationError
                }
            } else {
                self.logAuthenticationFailure(credentials)
                throw AuthenticationError
            }
        } else {
            self.logSecurityBypass(request)
            throw SecurityDisabledError
        }
    }
}
