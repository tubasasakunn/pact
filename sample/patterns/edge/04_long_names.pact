// Edge case: Very long names (50+ characters)

component ThisIsAnExtremelyLongComponentNameThatExceedsFiftyCharactersInLength {
    type ResultTypeWithExtremelyLongNameExceedingFiftyCharacters {
        fieldWithExtremelyLongNameThatExceedsFiftyCharacters: String
    }

    provides InterfaceWithExtremelyLongNameExceedingFiftyCharacters {
        thisIsAlsoAnExtremelyLongMethodNameThatExceedsFiftyCharacters(
            parameterWithVeryLongNameThatExceedsFiftyCharacters: String
        ) -> ResultTypeWithExtremelyLongNameExceedingFiftyCharacters
    }
}
