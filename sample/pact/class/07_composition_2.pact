// Pattern: Composition with 2 parts
// 1つのクラスが2つの部品を含むコンポジションパターン

component Computer {
    type ComputerData {
        id: string
        brand: string
        model: string
    }

    contains CPU
    contains Memory

    provides ComputerAPI {
        Boot() -> bool
        Shutdown() -> bool
        GetSpecs() -> ComputerData
    }
}

component CPU {
    type CPUData {
        cores: int
        clockSpeed: float
        architecture: string
    }

    provides CPUAPI {
        GetCores() -> int
        GetSpeed() -> float
    }
}

component Memory {
    type MemoryData {
        size: int
        memType: string
        speed: int
    }

    provides MemoryAPI {
        GetSize() -> int
        GetMemType() -> string
    }
}
