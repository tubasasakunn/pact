// Lexer Component - Pact DSL Tokenizer
@version("1.0")
@package("parser")
component Lexer {
    // Token types
    enum TokenType {
        EOF
        ILLEGAL
        IDENT
        STRING
        INT
        FLOAT
        DURATION
        COMPONENT
        IMPORT
        TYPE
        ENUM
        DEPENDS
        ON
        EXTENDS
        IMPLEMENTS
        CONTAINS
        AGGREGATES
        PROVIDES
        REQUIRES
        FLOW
        STATES
        STATE
        PARALLEL
        REGION
        INITIAL
        FINAL
        ENTRY
        EXIT
        IF
        ELSE
        FOR
        IN
        WHILE
        RETURN
        THROW
        AWAIT
        ASYNC
        THROWS
        WHEN
        AFTER
        DO
        TRUE
        FALSE
        NULL
        AS
        LBRACE
        RBRACE
        LPAREN
        RPAREN
        LBRACKET
        RBRACKET
        COMMA
        COLON
        DOT
        AT
        ARROW
        QUESTION
        PLUS
        MINUS
        STAR
        SLASH
        PERCENT
        EQ
        NE
        LT
        GT
        LE
        GE
        AND
        OR
        NOT
        ASSIGN
        NULLISH
    }

    type Token {
        tokenType: TokenType
        literal: string
        line: int
        column: int
    }

    type LexerState {
        input: string
        pos: int
        readPos: int
        ch: int
        line: int
        column: int
    }

    provides LexerAPI {
        NewLexer(input: string) -> Lexer
        NextToken() -> Token
    }
}
