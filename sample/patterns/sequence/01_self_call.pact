// Sequence Pattern 01: Self calls
component Service {
    flow SelfCalls {
        a = self.stepA()
        b = self.stepB(a)
        c = self.stepC(b)
        return c
    }
}
