// Pattern 9: Transition with event trigger
component EventTrigger {
    states MediaPlayer {
        initial Stopped

        state Stopped { }
        state Playing { }
        state Paused { }

        Stopped -> Playing on play
        Playing -> Paused on pause
        Paused -> Playing on resume
        Playing -> Stopped on stop
        Paused -> Stopped on stop
    }
}
