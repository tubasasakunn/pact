// Pattern: Linear States with 3 states
// 3つの状態が直線的に遷移するパターン

component TrafficLight {
    type LightData {
        id: string
        intersection: string
    }

    states LightState {
        initial Red
        final Green

        state Red {
            entry [stopTraffic]
        }

        state Yellow {
            entry [prepareToChange]
        }

        state Green {
            entry [allowTraffic]
        }

        Red -> Green on timer
        Green -> Yellow on timer
        Yellow -> Red on timer
    }
}
