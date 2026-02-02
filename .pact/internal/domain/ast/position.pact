// AST Position - Source code position tracking
// This file defines position information for AST nodes

component ASTPosition {
    // Position represents a location in source code
    type Position {
        file: String
        line: Int
        column: Int
        offset: Int
    }

    // Position provides a String() method for formatting:
    // - If file is empty: "line:column"
    // - If file is present: "file:line:column"

    // NoPos represents the absence of position information
    // A Position is valid when line > 0
}
