// Pattern: Inheritance Tree with 2 children
// 親クラスを2つの子クラスが継承するパターン

component Animal {
    type AnimalData {
        name: string
        age: int
    }

    provides AnimalAPI {
        GetInfo() -> AnimalData
        MakeSound() -> string
    }
}

component Dog {
    extends Animal

    type DogData {
        breed: string
        isGoodBoy: bool
    }

    provides DogAPI {
        Bark() -> string
        Fetch() -> bool
    }
}

component Cat {
    extends Animal

    type CatData {
        color: string
        isIndoor: bool
    }

    provides CatAPI {
        Meow() -> string
        Purr() -> bool
    }
}
