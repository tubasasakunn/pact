// Pattern 11: Special characters in strings
component StringPatterns {
    type StringData {
        normalString: string
        withEscapes: string
    }

    flow TestStrings {
        // Normal string
        simple = self.process(normalString)

        // String with special meaning
        withQuote = self.process(stringWithQuote)
        withNewline = self.process(stringWithNewline)
        withTab = self.process(stringWithTab)

        return simple
    }
}
