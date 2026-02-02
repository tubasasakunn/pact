// Pattern 20: Complex event names
component ComplexEvents {
    states NetworkConnection {
        initial Disconnected

        state Disconnected { }
        state Connecting { }
        state Connected { }
        state Reconnecting { }

        Disconnected -> Connecting on userRequestedConnect
        Connecting -> Connected on connectionEstablished
        Connected -> Disconnected on connectionLostUnexpectedly
        Connected -> Reconnecting on temporaryNetworkFailure
        Reconnecting -> Connected on reconnectionSuccessful
        Reconnecting -> Disconnected on maxRetriesExceeded
    }
}
