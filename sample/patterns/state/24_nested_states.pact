// Pattern 24: Nested states simulation (flat representation)
// Note: Pact may not support true nested states, simulating with naming
component NestedStates {
    states PhoneCall {
        initial Idle
        final Terminated

        state Idle { }
        state Ringing { }
        state Connected { }
        state ConnectedActive { }
        state ConnectedOnHold { }
        state ConnectedMuted { }
        state Terminated { }

        Idle -> Ringing on incomingCall
        Ringing -> Connected on answer
        Ringing -> Terminated on reject
        Connected -> ConnectedActive on startTalking
        ConnectedActive -> ConnectedOnHold on holdPressed
        ConnectedOnHold -> ConnectedActive on resumePressed
        ConnectedActive -> ConnectedMuted on mutePressed
        ConnectedMuted -> ConnectedActive on unmutePressed
        ConnectedActive -> Terminated on hangup
        ConnectedOnHold -> Terminated on hangup
        ConnectedMuted -> Terminated on hangup
    }
}
