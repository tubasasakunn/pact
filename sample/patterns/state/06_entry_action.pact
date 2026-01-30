// Pattern 6: State with entry action
component EntryAction {
    states DoorSystem {
        initial Closed

        state Closed {
            entry [lockDoor]
        }
        state Open {
            entry [unlockDoor]
        }

        Closed -> Open on open
        Open -> Closed on close
    }
}
