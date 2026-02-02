// Parser performs syntactic analysis on token streams
// Converts tokens into an abstract syntax tree (AST)

component Parser {
    depends on Lexer
    depends on Token
    depends on ast: AST
    depends on errors: Errors

    // ParserState holds the internal state during parsing
    type ParserState {
        lexer: Lexer
        currentToken: Token
        peekToken: Token
        parseErrors: Error[]
    }

    // OperatorPrecedence defines precedence levels for operators
    enum OperatorPrecedence {
        LOWEST
        OR
        AND
        EQUALITY
        COMPARISON
        SUM
        PRODUCT
        DOT
    }

    provides ParserInterface {
        // NewParser creates a new parser with the given lexer
        NewParser(lexer: Lexer) -> Parser

        // Parse parses the input and returns the AST
        Parse() -> SpecFile throws ParseError

        // ParseFile parses an entire file
        ParseFile() -> SpecFile throws ParseError

        // ParseString parses a string and returns the AST
        ParseString(input: string) -> SpecFile throws ParseError
    }

    // Parser initialization flow
    flow Initialize {
        self.nextToken()
        self.nextToken()
    }

    // Advance to next token
    flow NextToken {
        currentToken = peekToken
        peekToken = lexer.NextToken()
    }

    // Main file parsing flow
    flow ParseFile {
        spec = self.createSpecFile()
        self.parseDeclarations(spec)
        return spec
    }

    // Parse import declaration
    flow ParseImport {
        self.nextToken()
        path = currentToken.literal
        self.nextToken()
        alias = self.parseOptionalAlias()
        return self.createImportDecl(path, alias)
    }

    // Parse component declaration
    flow ParseComponent {
        self.nextToken()
        name = currentToken.literal
        self.nextToken()
        self.expectToken(LBRACE)
        body = self.parseComponentBody()
        self.expectToken(RBRACE)
        return self.createComponentDecl(name, body)
    }

    // Parse type declaration
    flow ParseTypeDecl {
        self.nextToken()
        name = currentToken.literal
        self.nextToken()
        self.expectToken(LBRACE)
        fields = self.parseFields()
        self.expectToken(RBRACE)
        return self.createTypeDecl(name, fields)
    }

    // Parse enum declaration
    flow ParseEnumDecl {
        self.nextToken()
        name = currentToken.literal
        self.nextToken()
        self.expectToken(LBRACE)
        values = self.parseEnumValues()
        self.expectToken(RBRACE)
        return self.createEnumDecl(name, values)
    }

    // Parse field declaration
    flow ParseField {
        annotations = self.parseOptionalAnnotations()
        visibility = self.parseVisibility()
        name = currentToken.literal
        self.nextToken()
        self.expectToken(COLON)
        self.nextToken()
        typeExpr = self.parseTypeExpr()
        return self.createFieldDecl(name, typeExpr, visibility, annotations)
    }

    // Parse type expression
    flow ParseTypeExpr {
        name = currentToken.literal
        self.nextToken()
        nullable = self.parseOptionalNullable()
        isArray = self.parseOptionalArray()
        return self.createTypeExpr(name, nullable, isArray)
    }

    // Parse interface declaration
    flow ParseInterface {
        self.nextToken()
        name = currentToken.literal
        self.nextToken()
        self.expectToken(LBRACE)
        methods = self.parseMethods()
        self.expectToken(RBRACE)
        return self.createInterfaceDecl(name, methods)
    }

    // Parse method declaration
    flow ParseMethod {
        annotations = self.parseOptionalAnnotations()
        isAsync = self.parseOptionalAsync()
        name = currentToken.literal
        self.nextToken()
        params = self.parseParams()
        returnType = self.parseOptionalReturnType()
        throwsList = self.parseOptionalThrows()
        return self.createMethodDecl(name, params, returnType, throwsList, isAsync, annotations)
    }

    // Parse flow declaration
    flow ParseFlow {
        self.nextToken()
        name = currentToken.literal
        self.nextToken()
        self.expectToken(LBRACE)
        steps = self.parseSteps()
        self.expectToken(RBRACE)
        return self.createFlowDecl(name, steps)
    }

    // Parse states declaration
    flow ParseStates {
        self.nextToken()
        name = currentToken.literal
        self.nextToken()
        self.expectToken(LBRACE)
        statesBody = self.parseStatesBody()
        self.expectToken(RBRACE)
        return self.createStatesDecl(name, statesBody)
    }

    // Parse expression with precedence
    flow ParseExpression {
        left = self.parsePrimaryExpr()
        result = self.parseInfixExpressions(left, precedence)
        return result
    }

    // Parse primary expression
    flow ParsePrimaryExpr {
        expr = self.matchPrimaryExpr()
        return expr
    }

    // Parse annotations
    flow ParseAnnotations {
        annotations = self.collectAnnotations()
        return annotations
    }

    // Get current operator precedence
    flow GetCurrentPrecedence {
        precedence = self.getPrecedenceForToken(currentToken.tokenType)
        return precedence
    }
}
