// Pattern: Inheritance Tree with 3 children
// 親クラスを3つの子クラスが継承するパターン

component Vehicle {
    type VehicleData {
        id: string
        manufacturer: string
        year: int
    }

    provides VehicleAPI {
        Start() -> bool
        Stop() -> bool
        GetInfo() -> VehicleData
    }
}

component Car {
    extends Vehicle

    type CarData {
        numDoors: int
        fuelType: string
    }

    provides CarAPI {
        Drive() -> bool
        Park() -> bool
    }
}

component Motorcycle {
    extends Vehicle

    type MotorcycleData {
        engineCC: int
        hasABS: bool
    }

    provides MotorcycleAPI {
        Wheelie() -> bool
        LeanAngle() -> int
    }
}

component Truck {
    extends Vehicle

    type TruckData {
        cargoCapacity: float
        numAxles: int
    }

    provides TruckAPI {
        LoadCargo() -> bool
        Haul() -> bool
    }
}
