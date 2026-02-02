// Pattern 11: Call chain (A->B->C)
component ServiceA {
    depends on ServiceB

    flow ChainedCall {
        result = ServiceB.processAndForward(data)
        return result
    }
}

component ServiceB {
    depends on ServiceC

    flow ProcessAndForward {
        processed = self.process(data)
        result = ServiceC.finalize(processed)
        return result
    }
}

component ServiceC {
    flow Finalize {
        completed = self.complete(data)
        return completed
    }
}
