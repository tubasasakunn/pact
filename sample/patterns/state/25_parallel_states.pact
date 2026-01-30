// Pattern 25: Parallel states simulation (flat representation)
// Note: Pact may not support true parallel regions, simulating with combined states
component ParallelStates {
    states WashingMachine {
        initial Off
        final Complete

        state Off { }
        state OnIdle { }
        state WashingFilling { }
        state WashingAgitating { }
        state WashingDraining { }
        state RinseFilling { }
        state RinseAgitating { }
        state RinseDraining { }
        state SpinningLow { }
        state SpinningHigh { }
        state Complete { }

        Off -> OnIdle on powerOn
        OnIdle -> WashingFilling on startWash
        WashingFilling -> WashingAgitating on waterFull
        WashingAgitating -> WashingDraining on washComplete
        WashingDraining -> RinseFilling on drainComplete
        RinseFilling -> RinseAgitating on waterFull
        RinseAgitating -> RinseDraining on rinseComplete
        RinseDraining -> SpinningLow on drainComplete
        SpinningLow -> SpinningHigh on spinUp
        SpinningHigh -> Complete on cycleComplete
        Complete -> Off on powerOff
        OnIdle -> Off on powerOff
    }
}
