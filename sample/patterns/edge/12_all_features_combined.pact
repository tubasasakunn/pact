// Pattern 12: All features combined in one component
@version("1.0.0")
@author("test-team")
@description("Comprehensive feature test")
component AllFeaturesCombined {
    // All relationship types
    extends BaseComponent
    implements IService
    depends on Database : IDatabase as db
    depends on Logger : ILogger as logger
    contains InternalProcessor
    aggregates SharedCache

    // Struct with all field types
    type CompleteType {
        +publicString: string
        +publicInt: int
        +publicFloat: float
        +publicBool: bool
        -privateField: string
        #protectedField: int
        ~packageField: bool
        nullableString: string?
        arrayField: string[]
        complexType: NestedType
        complexNullable: NestedType?
        complexArray: NestedType[]
    }

    // Enum type
    enum Status {
        PENDING
        ACTIVE
        COMPLETED
        FAILED
        CANCELLED
    }

    // Interface with all method features
    @api
    provides CompleteAPI {
        simpleMethod() -> string
        methodWithParams(id: string, count: int) -> Result
        methodWithNullable(data: Data?) -> Result
        methodWithArray(items: Item[]) -> Item[]
        async asyncMethod(id: string) -> AsyncResult
        methodWithThrows(input: Input) -> Output throws Error1, Error2
    }

    requires ExternalDependency {
        fetch(url: string) -> Response
        async fetchAsync(url: string) -> Response throws TimeoutError
    }

    // Flow with all step types
    flow CompleteFlow {
        // Assignment
        value = self.getValue()

        // Method calls
        result = db.query(queryString)

        // If-else
        if condition {
            success = self.handleSuccess(result)
        } else {
            errorResult = self.handleError(result)
        }

        // For loop
        for item in items {
            processed = self.processItem(item)
        }

        // While loop
        while hasMore {
            next = self.getNext()
        }

        // Await
        await self.asyncOperation()
        asyncResult = self.getAsyncResult()

        // Ternary
        finalValue = isValid ? successValue : errorValue

        // Null coalescing
        safeValue = maybeNull ?? defaultValue

        // Throw
        if isError {
            throw ProcessingError
        }

        // Return
        return finalValue
    }

    // State machine with all features
    states CompleteStateMachine {
        initial Idle
        final Completed
        final Failed

        state Idle {
            entry [initializeIdle]
            exit [cleanupIdle]
        }

        state Processing {
            entry [startProcessing]
            exit [stopProcessing]

            initial Validating

            state Validating { }
            state Executing { }
            state Finalizing { }

            Validating -> Executing on valid
            Validating -> Failed on invalid
            Executing -> Finalizing on done
            Finalizing -> Completed on success
        }

        parallel Monitoring {
            region HealthCheck {
                initial Healthy
                state Healthy { }
                state Degraded { }
                Healthy -> Degraded on issue
                Degraded -> Healthy on resolved
            }

            region MetricsCollection {
                initial Collecting
                state Collecting { }
                state Paused { }
                Collecting -> Paused on pauseMetrics
                Paused -> Collecting on resumeMetrics
            }
        }

        // All trigger types
        Idle -> Processing on start
        Processing -> Idle on cancel do [logCancellation, cleanup]
        Processing -> Failed after 30s when timeoutEnabled
        Monitoring -> Idle on stop
        Failed -> Idle on reset do [clearErrors, resetState]
    }
}
