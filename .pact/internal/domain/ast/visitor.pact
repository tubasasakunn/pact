// AST Visitor - Visitor pattern for AST traversal
// This file defines the visitor interface and walk functions

component ASTVisitor {
    // Visitor type for traversing AST nodes
    type VisitorCallbacks {
        visitSpecFile: func
        visitImportDecl: func
        visitComponentDecl: func
        visitAnnotationDecl: func
        visitTypeDecl: func
        visitFieldDecl: func
        visitRelationDecl: func
        visitInterfaceDecl: func
        visitMethodDecl: func
        visitFlowDecl: func
        visitStep: func
        visitExpr: func
        visitStatesDecl: func
        visitStateDecl: func
        visitTransitionDecl: func
    }

    // BaseVisitor provides default no-op implementations
    type BaseVisitor {
        callbacks: VisitorCallbacks
    }

    provides VisitorAPI {
        // Walk traverses the entire AST starting from a SpecFile
        Walk(visitor: VisitorCallbacks, node: SpecFile) -> Error?

        // WalkComponent traverses a component declaration
        WalkComponent(visitor: VisitorCallbacks, node: ComponentDecl) -> Error?

        // WalkFlow traverses a flow declaration
        WalkFlow(visitor: VisitorCallbacks, node: FlowDecl) -> Error?

        // WalkStates traverses a states declaration
        WalkStates(visitor: VisitorCallbacks, node: StatesDecl) -> Error?
    }

    // Walk flow - visits nodes in order:
    // 1. SpecFile itself
    // 2. All imports
    // 3. Component (if present)
    //    - Annotations
    //    - Types (with their annotations and fields)
    //    - Relations
    //    - Provides interfaces (with methods)
    //    - Requires interfaces (with methods)
    //    - Flows (with steps, recursively for nested steps)
    //    - States (with states and transitions)
    flow WalkAST {
        err = visitor.visitSpecFile(node)
        if err != null {
            return err
        }
        self.walkImports(visitor, node)
        self.walkComponents(visitor, node)
        return null
    }
}
