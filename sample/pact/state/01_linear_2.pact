// Pattern: Linear States with 2 states
// 2つの状態が直線的に遷移するパターン

component LightSwitch {
    type SwitchData {
        id: string
        location: string
    }

    states SwitchState {
        initial Off
        final On

        state Off {
            entry [turnOffLight]
        }

        state On {
            entry [turnOnLight]
        }

        Off -> On on toggle
        On -> Off on toggle
    }
}
