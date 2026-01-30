// Pattern 29: Multiple throw statements with different error types
component ErrorHandlingService {
    depends on InputValidator
    depends on ResourceManager
    depends on NetworkClient
    depends on DataProcessor

    flow ErrorHandling {
        // Validate input format
        formatValid = InputValidator.checkFormat(input)
        if formatInvalid {
            throw InvalidFormatError
        }

        // Validate input size
        sizeValid = InputValidator.checkSize(input)
        if sizeExceeded {
            throw SizeLimitExceededError
        }

        // Check resource availability
        resourceAvailable = ResourceManager.checkAvailability()
        if resourceUnavailable {
            throw ResourceUnavailableError
        }

        // Acquire resource with timeout
        resource = ResourceManager.acquire(timeout)
        if acquisitionTimeout {
            throw ResourceTimeoutError
        }

        // Network operation
        connection = NetworkClient.connect(endpoint)
        if connectionFailed {
            ResourceManager.release(resource)
            throw NetworkConnectionError
        }

        // Send request
        response = NetworkClient.send(request)
        if sendFailed {
            NetworkClient.disconnect(connection)
            ResourceManager.release(resource)
            throw NetworkSendError
        }

        // Process response
        processed = DataProcessor.process(response)
        if processingFailed {
            NetworkClient.disconnect(connection)
            ResourceManager.release(resource)
            throw DataProcessingError
        }

        // Validate result
        resultValid = DataProcessor.validate(processed)
        if resultInvalid {
            NetworkClient.disconnect(connection)
            ResourceManager.release(resource)
            throw ResultValidationError
        }

        // Cleanup and return
        NetworkClient.disconnect(connection)
        ResourceManager.release(resource)
        return processed
    }
}
