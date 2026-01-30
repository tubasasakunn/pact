// Parser Component - Pact DSL Parser
@version("1.0")
@package("parser")
component Parser {
    depends on Lexer : LexerAPI as lexer

    type ParserState {
        curToken: Token
        peekToken: Token
        errors: ParseError[]
    }

    type ParseError {
        line: int
        column: int
        message: string
    }

    provides ParserAPI {
        NewParser(lexer: Lexer) -> Parser
        ParseFile() -> SpecFile throws ParseError
    }

    provides InternalParsing {
        parseComponent() -> ComponentDecl throws ParseError
        parseImport() -> ImportDecl throws ParseError
        parseTypeDecl() -> TypeDecl throws ParseError
        parseEnumDecl() -> TypeDecl throws ParseError
        parseField() -> FieldDecl throws ParseError
        parseRelation() -> RelationDecl throws ParseError
        parseInterface() -> InterfaceDecl throws ParseError
        parseMethod() -> MethodDecl throws ParseError
        parseFlow() -> FlowDecl throws ParseError
        parseStep() -> Step throws ParseError
        parseExpression() -> Expr throws ParseError
        parseStates() -> StatesDecl throws ParseError
        parseStateDecl() -> StateDecl throws ParseError
        parseTransition() -> TransitionDecl throws ParseError
        parseAnnotations() -> AnnotationDecl[] throws ParseError
    }
}
