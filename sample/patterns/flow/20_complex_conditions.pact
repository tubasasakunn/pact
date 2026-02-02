// Pattern 20: Complex condition expressions
component ComplexConditionsService {
    flow ComplexConditions {
        user = self.getUser(userId)
        permissions = self.getPermissions(user)
        resource = self.getResource(resourceId)
        if userActive {
            if hasPermission {
                if resourceAvailable {
                    if quotaRemaining {
                        result = self.executeOperation(user, resource)
                        self.logAccess(user, resource)
                        return result
                    } else {
                        throw QuotaExceededError
                    }
                } else {
                    throw ResourceUnavailableError
                }
            } else {
                throw PermissionDeniedError
            }
        } else {
            throw UserInactiveError
        }
    }
}
