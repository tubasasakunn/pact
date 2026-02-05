// Pattern: Chain with 3 participants
// 3つのサービスが連鎖するパターン

component Gateway {
    type GatewayRequest {
        clientId: string
        action: string
    }

    depends on Middleware

    flow ProcessRequest {
        validated = self.validateClient(request)
        enriched = Middleware.enrich(validated)
        result = self.formatResponse(enriched)
        return result
    }
}

component Middleware {
    type MiddlewareData {
        metadata: string
        context: string
    }

    depends on Backend

    provides MiddlewareAPI {
        Enrich(data: string) -> string
    }

    flow Enrich {
        context = self.buildContext(data)
        processed = Backend.execute(data, context)
        enriched = self.addMetadata(processed)
        return enriched
    }
}

component Backend {
    type BackendResult {
        data: string
        status: string
    }

    provides BackendAPI {
        Execute(data: string, context: string) -> string
    }

    flow Execute {
        validated = self.validate(data)
        result = self.process(validated, context)
        self.audit(result)
        return result
    }
}
