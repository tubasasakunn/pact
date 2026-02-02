// Lexer performs lexical analysis on source code
// Converts input string into a stream of tokens

component Lexer {
    depends on Token

    // LexerState holds the internal state during lexical analysis
    type LexerState {
        input: string
        pos: int
        readPos: int
        currentChar: int
        lineNumber: int
        columnNumber: int
        tokenLineNumber: int
        tokenColumnNumber: int
    }

    // KeywordEntry maps a keyword string to its token type
    type KeywordEntry {
        keyword: string
        tokenType: TokenType
    }

    provides LexerInterface {
        // NewLexer creates a new lexer for the given input string
        NewLexer(input: string) -> Lexer

        // NextToken returns the next token from the input
        NextToken() -> Token

        // ReadChar advances to the next character
        ReadChar()

        // SkipWhitespace skips whitespace characters
        SkipWhitespace()

        // ReadIdentifier reads an identifier or keyword
        ReadIdentifier() -> string

        // ReadNumber reads a numeric literal
        ReadNumber() -> Token

        // ReadString reads a string literal
        ReadString() -> string

        // LookupIdent checks if identifier is a keyword
        LookupIdent(ident: string) -> TokenType
    }

    // Lexer initialization flow
    flow Initialize {
        lineNumber = 1
        columnNumber = 0
        self.ReadChar()
    }

    // Main tokenization flow
    flow Tokenize {
        self.SkipWhitespace()
        self.recordTokenPosition()
        token = self.matchToken()
        return token
    }
}
