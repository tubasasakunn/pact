// Pattern: Interface Implementation with 2 implementers
// 1つのインターフェースを2つのクラスが実装するパターン

component Serializer {
    provides SerializerAPI {
        Serialize(data: string) -> string
        Deserialize(data: string) -> string
    }
}

component JSONSerializer {
    type JSONConfig {
        prettyPrint: bool
        indent: int
    }

    implements Serializer

    provides JSONSerializerAPI {
        SetConfig(config: string) -> bool
    }
}

component XMLSerializer {
    type XMLConfig {
        includeDeclaration: bool
        encoding: string
    }

    implements Serializer

    provides XMLSerializerAPI {
        SetConfig(config: string) -> bool
    }
}
