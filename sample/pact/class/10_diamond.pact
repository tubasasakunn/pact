// Pattern: Diamond Dependency
// ダイヤモンド継承パターン（上位が2つの中間を経由して下位に依存）

component Application {
    type AppData {
        name: string
        version: string
    }

    depends on ServiceA
    depends on ServiceB

    provides AppAPI {
        Run() -> bool
        GetVersion() -> string
    }
}

component ServiceA {
    type ServiceAData {
        configA: string
    }

    depends on CoreLibrary

    provides ServiceAAPI {
        ProcessA() -> string
    }
}

component ServiceB {
    type ServiceBData {
        configB: string
    }

    depends on CoreLibrary

    provides ServiceBAPI {
        ProcessB() -> string
    }
}

component CoreLibrary {
    type CoreData {
        version: string
        initialized: bool
    }

    provides CoreAPI {
        Initialize() -> bool
        GetVersion() -> string
        Cleanup() -> bool
    }
}
