// Pattern 3: Multiple states (5+)
component MultipleStates {
    states TrafficLight {
        initial Red

        state Red { }
        state RedYellow { }
        state Green { }
        state Yellow { }
        state Blinking { }
        state Off { }

        Red -> RedYellow on prepare
        RedYellow -> Green on go
        Green -> Yellow on slow
        Yellow -> Red on stop
        Red -> Blinking on malfunction
        Blinking -> Off on shutdown
    }
}
