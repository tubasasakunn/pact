// 19: Complex Expressions in Flows - ternary, null coalescing, binary ops
component Calculator {
    type Result {
        value: float
        error: string?
    }

    depends on MathEngine
    depends on Logger

    flow ComputeDiscount {
        basePrice = self.getBasePrice(productId)
        memberDiscount = self.getMemberDiscount(userId)

        // Ternary expression
        discount = memberDiscount > 20 ? 20 : memberDiscount

        // Binary operations
        discountedPrice = basePrice - (basePrice * discount / 100)
        tax = discountedPrice * 0.08
        total = discountedPrice + tax

        if total < 0 {
            Logger.Error("Negative total calculated")
            throw CalculationError
        }

        return total
    }

    flow ProcessData {
        raw = MathEngine.fetchData(sourceId)

        // Null coalescing
        value = raw ?? 0
        config = self.getConfig(key) ?? self.defaultConfig()

        // Complex conditional
        if value > 100 {
            result = MathEngine.processLarge(value)
        } else {
            if value > 0 {
                result = MathEngine.processSmall(value)
            } else {
                result = MathEngine.processZero()
            }
        }

        Logger.Info("Processed data")
        return result
    }
}

component MathEngine {
    provides MathAPI {
        FetchData(sourceId: string) -> float
        ProcessLarge(value: float) -> float
        ProcessSmall(value: float) -> float
        ProcessZero() -> float
    }
}

component Logger {
    provides LogAPI {
        Info(msg: string)
        Error(msg: string)
    }
}
