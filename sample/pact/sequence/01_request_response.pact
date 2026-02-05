// Pattern: Request-Response
// シンプルなリクエスト・レスポンスパターン

component Client {
    type Request {
        id: string
        payload: string
    }

    depends on Server

    flow SendRequest {
        request = self.createRequest(data)
        response = Server.handleRequest(request)
        result = self.processResponse(response)
        return result
    }
}

component Server {
    type Response {
        id: string
        status: int
        body: string
    }

    provides ServerAPI {
        HandleRequest(request: string) -> string
    }

    flow HandleRequest {
        validated = self.validateRequest(request)
        processed = self.processRequest(validated)
        response = self.createResponse(processed)
        return response
    }
}
