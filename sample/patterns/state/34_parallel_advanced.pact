// Pattern 34: Advanced parallel states
component MediaPlayer {
    states PlayerStates {
        initial Off

        state Off { }

        parallel Playing {
            region AudioControl {
                initial Muted

                state Muted { }
                state Audible {
                    entry [initializeAudio]
                    exit [cleanupAudio]
                }
                state FadingIn { }
                state FadingOut { }

                Muted -> FadingIn on unmute
                FadingIn -> Audible after 500ms
                Audible -> FadingOut on mute
                FadingOut -> Muted after 500ms
                Audible -> Audible on volumeChange
            }

            region VideoControl {
                initial Hidden

                state Hidden { }
                state Visible {
                    entry [initializeVideo]
                    exit [cleanupVideo]
                }
                state Buffering { }
                state Fullscreen {
                    entry [enterFullscreen]
                    exit [exitFullscreen]
                }

                Hidden -> Buffering on show
                Buffering -> Visible on ready
                Visible -> Hidden on hide
                Visible -> Fullscreen on maximize
                Fullscreen -> Visible on minimize
                Buffering -> Hidden on error
            }

            region PlaybackControl {
                initial Paused

                state Paused { }
                state Playing {
                    entry [startPlayback]
                    exit [stopPlayback]
                }
                state Seeking { }
                state Buffering { }

                Paused -> Playing on play
                Playing -> Paused on pause
                Playing -> Seeking on seek
                Seeking -> Playing on seekComplete
                Playing -> Buffering on bufferEmpty
                Buffering -> Playing on bufferFull
            }
        }

        state Error {
            entry [logError, showErrorUI]
            exit [clearError]
        }

        Off -> Playing on start
        Playing -> Off on stop
        Playing -> Error on fatalError
        Error -> Off on dismiss
        Error -> Playing on retry
    }
}

component SmartThermostat {
    states ThermostatStates {
        initial Idle

        state Idle {
            entry [displayIdleScreen]
        }

        parallel Active {
            region TemperatureControl {
                initial Monitoring

                state Monitoring { }
                state Heating {
                    entry [activateHeater]
                    exit [deactivateHeater]
                }
                state Cooling {
                    entry [activateAC]
                    exit [deactivateAC]
                }

                Monitoring -> Heating on tooCold
                Monitoring -> Cooling on tooHot
                Heating -> Monitoring on targetReached
                Cooling -> Monitoring on targetReached
            }

            region HumidityControl {
                initial Monitoring

                state Monitoring { }
                state Humidifying {
                    entry [activateHumidifier]
                    exit [deactivateHumidifier]
                }
                state Dehumidifying {
                    entry [activateDehumidifier]
                    exit [deactivateDehumidifier]
                }

                Monitoring -> Humidifying on tooDry
                Monitoring -> Dehumidifying on tooHumid
                Humidifying -> Monitoring on optimalHumidity
                Dehumidifying -> Monitoring on optimalHumidity
            }

            region FanControl {
                initial Off

                state Off { }
                state Low {
                    entry [setFanLow]
                }
                state Medium {
                    entry [setFanMedium]
                }
                state High {
                    entry [setFanHigh]
                }
                state Auto {
                    entry [enableAutoFan]
                }

                Off -> Low on fanOn
                Low -> Medium on increase
                Medium -> High on increase
                High -> Medium on decrease
                Medium -> Low on decrease
                Low -> Off on fanOff
                Off -> Auto on autoMode
                Auto -> Off on manualMode
            }
        }

        state Sleep {
            entry [dimDisplay, setNightMode]
            exit [brightDisplay, setDayMode]
        }

        Idle -> Active on activate
        Active -> Idle on deactivate
        Active -> Sleep on nightTime
        Sleep -> Active on wakeUp
        Idle -> Sleep on nightTime
        Sleep -> Idle on wakeUp when notActive
    }
}
