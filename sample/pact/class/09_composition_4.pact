// Pattern: Composition with 4 parts
// 1つのクラスが4つの部品を含むコンポジションパターン

component Smartphone {
    type SmartphoneData {
        brand: string
        model: string
        osVersion: string
    }

    contains Display
    contains Battery
    contains Camera
    contains Processor

    provides SmartphoneAPI {
        PowerOn() -> bool
        PowerOff() -> bool
        GetBatteryLevel() -> int
    }
}

component Display {
    type DisplayData {
        resolution: string
        size: float
        displayType: string
    }

    provides DisplayAPI {
        SetBrightness(level: int) -> bool
        GetResolution() -> string
    }
}

component Battery {
    type BatteryData {
        capacity: int
        batteryType: string
        healthPercent: int
    }

    provides BatteryAPI {
        GetCapacity() -> int
        GetHealth() -> int
    }
}

component Camera {
    type CameraData {
        megapixels: int
        hasOIS: bool
        aperture: float
    }

    provides CameraAPI {
        TakePhoto() -> string
        RecordVideo() -> bool
    }
}

component Processor {
    type ProcessorData {
        model: string
        cores: int
        speed: float
    }

    provides ProcessorAPI {
        GetModel() -> string
        GetPerformance() -> float
    }
}
