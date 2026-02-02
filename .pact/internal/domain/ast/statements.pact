// AST Statements - Flow and step definitions for the Abstract Syntax Tree
// This file defines flow declarations and all step types

component ASTStatements {
    // FlowDecl represents a flow definition
    type FlowDecl {
        pos: Position
        name: String
        annotations: AnnotationDecl[]
        steps: Step[]
    }

    // Step is the base interface for all flow steps
    // All step types implement stepNode() and GetPos()

    // AssignStep represents an assignment step
    type AssignStep {
        pos: Position
        variable: String
        value: Expr
        annotations: AnnotationDecl[]
    }

    // CallStep represents a call step
    type CallStep {
        pos: Position
        expr: Expr
        isAwaited: Boolean
        annotations: AnnotationDecl[]
    }

    // ReturnStep represents a return step
    type ReturnStep {
        pos: Position
        value: Expr?                       // nil means no return value
        annotations: AnnotationDecl[]
    }

    // ThrowStep represents a throw step
    type ThrowStep {
        pos: Position
        error: String
        annotations: AnnotationDecl[]
    }

    // IfStep represents a conditional branching step
    type IfStep {
        pos: Position
        condition: Expr
        thenSteps: Step[]
        elseSteps: Step[]
        annotations: AnnotationDecl[]
    }

    // ForStep represents a for loop step
    type ForStep {
        pos: Position
        variable: String
        iterable: Expr
        body: Step[]
        annotations: AnnotationDecl[]
    }

    // WhileStep represents a while loop step
    type WhileStep {
        pos: Position
        condition: Expr
        body: Step[]
        annotations: AnnotationDecl[]
    }
}
