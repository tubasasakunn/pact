// Pattern 3: Multiple participants (5+)
component Orchestrator {
    depends on ServiceA
    depends on ServiceB
    depends on ServiceC
    depends on ServiceD
    depends on ServiceE

    flow CoordinateAll {
        a = ServiceA.process()
        b = ServiceB.process()
        c = ServiceC.process()
        d = ServiceD.process()
        e = ServiceE.process()
        return true
    }
}

component ServiceA { }
component ServiceB { }
component ServiceC { }
component ServiceD { }
component ServiceE { }
