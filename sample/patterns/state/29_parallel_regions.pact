// Pattern 29: Parallel states with multiple regions
// Demonstrates concurrent state machine regions
component ParallelRegions {
    states SmartHome {
        initial Initializing
        final ShutDown

        state Initializing {
            entry [loadConfiguration, connectSensors]
        }
        state ShutDown {
            entry [saveState, disconnectAll]
        }

        // Main active state
        state Active { }

        // Parallel region 1: Lighting control
        state LightingOff { }
        state LightingDimmed { }
        state LightingBright { }
        state LightingAuto { }

        // Parallel region 2: Climate control
        state ClimateOff { }
        state ClimateHeating { }
        state ClimateCooling { }
        state ClimateAuto { }

        // Parallel region 3: Security system
        state SecurityDisarmed { }
        state SecurityArmedHome { }
        state SecurityArmedAway { }
        state SecurityAlarm { }

        // Parallel region 4: Entertainment
        state EntertainmentOff { }
        state EntertainmentMusic { }
        state EntertainmentTV { }
        state EntertainmentMovie { }

        // System transitions
        Initializing -> Active on systemReady
        Active -> ShutDown on shutdownCommand

        // Lighting region transitions (independent)
        Active -> LightingOff on lightsOff
        LightingOff -> LightingDimmed on dimLights
        LightingDimmed -> LightingBright on brightenLights
        LightingBright -> LightingDimmed on dimLights
        LightingOff -> LightingAuto on autoLighting
        LightingAuto -> LightingOff on manualOverride
        LightingDimmed -> LightingAuto on autoLighting
        LightingBright -> LightingAuto on autoLighting

        // Climate region transitions (independent)
        Active -> ClimateOff on climateOff
        ClimateOff -> ClimateHeating on heatOn
        ClimateOff -> ClimateCooling on coolOn
        ClimateHeating -> ClimateOff on climateOff
        ClimateCooling -> ClimateOff on climateOff
        ClimateOff -> ClimateAuto on autoClimate
        ClimateHeating -> ClimateAuto on autoClimate
        ClimateCooling -> ClimateAuto on autoClimate
        ClimateAuto -> ClimateOff on manualOverride

        // Security region transitions (independent)
        Active -> SecurityDisarmed on disarmSecurity
        SecurityDisarmed -> SecurityArmedHome on armHome
        SecurityDisarmed -> SecurityArmedAway on armAway
        SecurityArmedHome -> SecurityDisarmed on disarm
        SecurityArmedAway -> SecurityDisarmed on disarm
        SecurityArmedHome -> SecurityAlarm on motionDetected
        SecurityArmedAway -> SecurityAlarm on motionDetected
        SecurityAlarm -> SecurityDisarmed on alarmDismissed

        // Entertainment region transitions (independent)
        Active -> EntertainmentOff on entertainmentOff
        EntertainmentOff -> EntertainmentMusic on playMusic
        EntertainmentOff -> EntertainmentTV on watchTV
        EntertainmentOff -> EntertainmentMovie on startMovie
        EntertainmentMusic -> EntertainmentOff on stopMedia
        EntertainmentTV -> EntertainmentOff on stopMedia
        EntertainmentMovie -> EntertainmentOff on stopMedia
        EntertainmentMusic -> EntertainmentTV on switchToTV
        EntertainmentTV -> EntertainmentMusic on switchToMusic
    }
}
