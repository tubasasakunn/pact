// Pattern 24: Real-world scenario - API call
component APIClient {
    depends on RateLimiter
    depends on Cache
    depends on HTTPClient
    depends on ResponseParser

    flow MakeAPIRequest {
        allowed = RateLimiter.checkLimit(apiKey)
        if allowed {
            cached = Cache.get(requestKey)
            if cached != null {
                return cached
            } else {
                response = HTTPClient.send(request)
                if response.status == 200 {
                    parsed = ResponseParser.parse(response.body)
                    Cache.set(requestKey, parsed)
                    return parsed
                } else {
                    throw APIError
                }
            }
        } else {
            throw RateLimitError
        }
    }
}

component RateLimiter { }
component Cache { }
component HTTPClient { }
component ResponseParser { }
