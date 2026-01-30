// EXPECTED BEHAVIOR: Two components with same name
// NOTE: Parser currently ALLOWS duplicate component names

component DuplicateName {
    provides API1 {
        Method1() -> Void
    }
}

component DuplicateName {
    provides API2 {
        Method2() -> Void
    }
}
