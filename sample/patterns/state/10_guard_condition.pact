// Pattern 10: Transition with guard condition
component GuardCondition {
    states Elevator {
        initial Idle

        state Idle { }
        state MovingUp { }
        state MovingDown { }
        state DoorOpen { }

        Idle -> MovingUp on requestUp
        Idle -> MovingDown on requestDown
        MovingUp -> DoorOpen on arrived
        MovingDown -> DoorOpen on arrived
        DoorOpen -> Idle on doorClosed
    }
}
