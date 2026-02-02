// AST Expressions - Expression types for the Abstract Syntax Tree
// This file defines all expression node types used in flows and conditions

component ASTExpressions {
    // Expr is the base interface for all expressions
    // All expression types implement exprNode() and GetPos()

    // LiteralExpr represents a literal value expression
    type LiteralExpr {
        pos: Position
        value: Any                         // Can be string, number, boolean, etc.
    }

    // VariableExpr represents a variable reference
    type VariableExpr {
        pos: Position
        name: String
    }

    // FieldExpr represents a field access expression
    type FieldExpr {
        pos: Position
        object: Expr                       // The object being accessed
        field: String                      // The field name
    }

    // CallExpr represents a method call expression
    type CallExpr {
        pos: Position
        object: Expr                       // The object on which method is called
        method: String                     // The method name
        args: Expr[]                       // Call arguments
    }

    // BinaryExpr represents a binary operation expression
    type BinaryExpr {
        pos: Position
        left: Expr
        op: String                         // Operator: +, -, *, /, ==, !=, etc.
        right: Expr
    }

    // UnaryExpr represents a unary operation expression
    type UnaryExpr {
        pos: Position
        op: String                         // Operator: !, -, +, etc.
        operand: Expr
    }

    // TernaryExpr represents a ternary conditional expression
    type TernaryExpr {
        pos: Position
        condition: Expr
        thenExpr: Expr
        elseExpr: Expr
    }

    // NullishExpr represents a null coalescing expression
    type NullishExpr {
        pos: Position
        left: Expr
        right: Expr
        throwErr: String?                  // For "?? throw Error" pattern
    }
}
