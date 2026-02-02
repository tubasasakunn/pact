// AST Nodes - Core node types for the Abstract Syntax Tree
// This file defines the fundamental structures that represent parsed .pact files

component ASTNodes {
    // SpecFile represents an entire .pact file
    type SpecFile {
        path: String
        imports: ImportDecl[]
        mainComponent: ComponentDecl?      // Kept for backward compatibility
        components: ComponentDecl[]        // Multiple component support
        interfaces: InterfaceDecl[]        // File-level interfaces
        types: TypeDecl[]                  // File-level type definitions
        annotations: AnnotationDecl[]      // File-level annotations
    }

    // ImportDecl represents an import statement
    type ImportDecl {
        pos: Position
        path: String
        alias: String?
    }

    // ComponentDecl represents a component declaration
    type ComponentDecl {
        pos: Position
        name: String
        annotations: AnnotationDecl[]
        body: ComponentBody
    }

    // ComponentBody represents the contents of a component
    type ComponentBody {
        types: TypeDecl[]
        relations: RelationDecl[]
        providesInterfaces: InterfaceDecl[]
        requiresInterfaces: InterfaceDecl[]
        flows: FlowDecl[]
        statesBlocks: StatesDecl[]
    }

    // AnnotationDecl represents an annotation
    type AnnotationDecl {
        pos: Position
        name: String
        args: AnnotationArg[]
    }

    // AnnotationArg represents an annotation argument
    type AnnotationArg {
        key: String?                       // nil means positional argument
        value: String
    }
}
