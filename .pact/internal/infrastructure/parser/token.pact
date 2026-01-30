// Token types and structures for lexical analysis
// This component defines the token system used by the lexer and parser

component Token {
    // TokenType represents the kind of token
    enum TokenType {
        EOF
        ILLEGAL

        // Identifiers and literals
        IDENT
        STRING
        INT
        FLOAT
        DURATION

        // Keywords
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

        // Symbols and operators
        LBRACE
        RBRACE
        LPAREN
        RPAREN
        LBRACKET
        RBRACKET
        COLON
        COMMA
        DOT
        ARROW
        AT
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
        HASH
        TILDE
    }

    // Token represents a lexical token with its metadata
    type Token {
        tokenType: TokenType
        literal: string
        lineNumber: int
        columnNumber: int
    }

    // Position represents a location in source code
    type Position {
        lineNumber: int
        columnNumber: int
    }

    // DurationValue represents a duration literal value
    type DurationValue {
        amount: int
        unit: string
    }
}
