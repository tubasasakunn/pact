// Pattern 17: Return value handling
component DataTransformer {
    depends on Parser
    depends on Formatter

    flow TransformData {
        parsed = Parser.parse(input)
        validated = Parser.validate(parsed)
        formatted = Formatter.format(validated)
        serialized = Formatter.serialize(formatted)
        return serialized
    }
}

component Parser { }
component Formatter { }
