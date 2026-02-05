// Pattern: Star Topology
// 中央の状態から放射状に遷移するパターン

component MediaPlayer {
    type PlayerData {
        currentTrack: string
        volume: int
        position: float
    }

    states PlayerState {
        initial Idle

        state Idle {
            entry [showReady]
        }

        state Playing {
            entry [startPlayback]
            exit [updatePosition]
        }

        state Paused {
            entry [pausePlayback]
        }

        state Buffering {
            entry [showLoading]
        }

        state Error {
            entry [showError]
        }

        Idle -> Playing on play
        Idle -> Buffering on load
        Playing -> Paused on pause
        Playing -> Idle on stop
        Playing -> Error on playbackError
        Paused -> Playing on resume
        Paused -> Idle on stop
        Buffering -> Playing on bufferComplete
        Buffering -> Error on bufferError
        Error -> Idle on reset
    }
}
