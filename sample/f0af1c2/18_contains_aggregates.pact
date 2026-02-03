// 18: Composition (contains) and Aggregation (aggregates)
component Car {
    type CarSpec {
        make: string
        model: string
        year: int
        vin: string
    }

    contains Engine
    contains Transmission
    contains Chassis
    aggregates GPS
    aggregates EntertainmentSystem
    aggregates InsurancePolicy

    provides CarAPI {
        Start() -> bool
        Stop() -> bool
        GetStatus() -> string
    }
}

component Engine {
    type EngineSpec {
        displacement: float
        horsepower: int
        fuelType: string
    }

    contains FuelInjector
    contains Turbocharger

    provides EngineAPI {
        Ignite() -> bool
        GetRPM() -> int
    }
}

component FuelInjector {
    type InjectorConfig {
        nozzleSize: float
        pressure: float
    }
}

component Turbocharger {
    type TurboSpec {
        maxBoost: float
        wastegateSize: float
    }
}

component Transmission {
    type GearRatio {
        gear: int
        ratio: float
    }

    enum TransmissionType {
        Manual
        Automatic
        CVT
        DCT
    }
}

component Chassis {
    type ChassisSpec {
        material: string
        weight: float
    }
}

component GPS {
    type Location {
        lat: float
        lng: float
    }

    provides NavigationAPI {
        GetCurrentLocation() -> Location
        Navigate(destination: string) -> string[]
    }
}

component EntertainmentSystem {
    provides MediaAPI {
        PlayMusic(track: string) -> bool
        SetVolume(level: int) -> bool
    }
}

component InsurancePolicy {
    type Policy {
        provider: string
        policyNumber: string
        expiresAt: string
    }
}
