// Pattern: Composition with 3 parts
// 1つのクラスが3つの部品を含むコンポジションパターン

component House {
    type HouseData {
        address: string
        sqft: int
        yearBuilt: int
    }

    contains Foundation
    contains Walls
    contains Roof

    provides HouseAPI {
        GetAddress() -> string
        GetValue() -> float
        Inspect() -> bool
    }
}

component Foundation {
    type FoundationData {
        foundationType: string
        depth: float
        material: string
    }

    provides FoundationAPI {
        CheckIntegrity() -> bool
    }
}

component Walls {
    type WallsData {
        material: string
        insulation: string
        thickness: float
    }

    provides WallsAPI {
        GetInsulationRating() -> float
    }
}

component Roof {
    type RoofData {
        style: string
        material: string
        age: int
    }

    provides RoofAPI {
        NeedsReplacement() -> bool
        GetAge() -> int
    }
}
