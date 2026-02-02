// Pattern 11: Transition with action
component TransitionAction {
    states ATM {
        initial Idle

        state Idle { }
        state CardInserted { }
        state PinEntered { }
        state Dispensing { }

        Idle -> CardInserted on insertCard
        CardInserted -> PinEntered on enterPin
        PinEntered -> Dispensing on withdraw
        Dispensing -> Idle on complete
    }
}
