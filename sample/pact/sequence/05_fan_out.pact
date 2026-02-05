// Pattern: Fan-Out
// 1つのサービスが複数のサービスを並列呼び出しするパターン

component Orchestrator {
    type AggregatedResult {
        serviceA: string
        serviceB: string
        serviceC: string
        combined: string
    }

    depends on ServiceA
    depends on ServiceB
    depends on ServiceC

    flow FanOutRequest {
        // Parallel calls to multiple services
        resultA = ServiceA.process(request)
        resultB = ServiceB.process(request)
        resultC = ServiceC.process(request)

        // Aggregate results
        combined = self.aggregate(resultA, resultB, resultC)
        return combined
    }
}

component ServiceA {
    type ServiceAResult {
        data: string
        latency: int
    }

    provides ServiceAAPI {
        Process(request: string) -> string
    }

    flow Process {
        result = self.executeLogicA(request)
        return result
    }
}

component ServiceB {
    type ServiceBResult {
        data: string
        latency: int
    }

    provides ServiceBAPI {
        Process(request: string) -> string
    }

    flow Process {
        result = self.executeLogicB(request)
        return result
    }
}

component ServiceC {
    type ServiceCResult {
        data: string
        latency: int
    }

    provides ServiceCAPI {
        Process(request: string) -> string
    }

    flow Process {
        result = self.executeLogicC(request)
        return result
    }
}
