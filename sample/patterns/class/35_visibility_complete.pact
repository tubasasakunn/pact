// Pattern 35: Complete visibility modifiers
component VisibilityDemo {
    type PublicType {
        +publicField: string
        +anotherPublic: int
    }

    type PrivateType {
        -privateField: string
        -secretKey: string
        -internalId: int
    }

    type ProtectedType {
        #protectedField: string
        #inheritableValue: int
    }

    type PackageType {
        ~packageField: string
        ~internalOnly: bool
    }

    type MixedVisibility {
        +id: string
        +name: string
        -password: string
        -apiKey: string
        #baseUrl: string
        ~internalData: string
    }

    type FullyAnnotated {
        +publicId: string
        +publicName: string
        +publicItems: Item[]
        -privateSecret: string
        -privateConfig: Config
        -privateList: string[]
        #protectedBase: Base
        #protectedArray: Base[]
        ~packageInternal: Internal
        ~packageData: Data[]
    }

    provides PublicAPI {
        getPublicData() -> Data
        setPublicData(data: Data) -> bool
    }

    provides ProtectedAPI {
        getProtectedData() -> Data
        setProtectedData(data: Data) -> bool
    }

    provides PrivateAPI {
        getPrivateData() -> Data
        setPrivateData(data: Data) -> bool
    }

    provides PackageAPI {
        getPackageData() -> Data
        setPackageData(data: Data) -> bool
    }
}
