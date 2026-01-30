// Pattern 7: State with exit action
component ExitAction {
    states AlarmSystem {
        initial Armed

        state Armed {
            exit [disableAlarm]
        }
        state Disarmed {
            exit [enableAlarm]
        }

        Armed -> Disarmed on disarm
        Disarmed -> Armed on arm
    }
}
