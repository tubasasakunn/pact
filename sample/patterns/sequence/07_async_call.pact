// Pattern 7: Async call (await simulation)
component AsyncHandler {
    depends on ExternalAPI

    flow FetchAsyncData {
        request = self.prepareRequest()
        response = ExternalAPI.asyncFetch(request)
        processed = self.processResponse(response)
        return processed
    }
}

component ExternalAPI {
    provides AsyncAPI {
        AsyncFetch(request: string)
    }
}
