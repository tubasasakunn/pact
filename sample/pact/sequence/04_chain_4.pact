// Pattern: Chain with 4 participants
// 4つのサービスが連鎖するパターン

component APIGateway {
    depends on AuthService

    flow HandleAPIRequest {
        validated = AuthService.authenticate(request)
        return validated
    }
}

component AuthService {
    depends on RateLimiter

    provides AuthAPI {
        Authenticate(request: string) -> string
    }

    flow Authenticate {
        token = self.extractToken(request)
        user = self.validateToken(token)
        allowed = RateLimiter.checkLimit(user)
        return allowed
    }
}

component RateLimiter {
    depends on BusinessLogic

    provides RateLimitAPI {
        CheckLimit(user: string) -> string
    }

    flow CheckLimit {
        current = self.getCurrentCount(user)
        if current < limit {
            self.incrementCount(user)
            result = BusinessLogic.process(user)
            return result
        } else {
            throw RateLimitExceeded
        }
    }
}

component BusinessLogic {
    provides BusinessAPI {
        Process(user: string) -> string
    }

    flow Process {
        data = self.fetchData(user)
        processed = self.transform(data)
        self.save(processed)
        return processed
    }
}
