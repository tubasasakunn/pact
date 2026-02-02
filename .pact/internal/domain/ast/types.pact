// AST Types - Type system definitions for the Abstract Syntax Tree
// This file defines types, fields, relations, interfaces, and methods

component ASTTypes {
    // TypeDecl represents a type definition
    type TypeDecl {
        pos: Position
        name: String
        kind: TypeKind
        annotations: AnnotationDecl[]
        fields: FieldDecl[]                // For struct types
        values: String[]                   // For enum types
    }

    // TypeKind enumeration for type categories
    enum TypeKind {
        STRUCT
        ENUM
    }

    // FieldDecl represents a field definition
    type FieldDecl {
        pos: Position
        name: String
        fieldType: TypeExpr
        visibility: Visibility
        annotations: AnnotationDecl[]
    }

    // Visibility enumeration for access modifiers
    enum Visibility {
        PUBLIC
        PRIVATE
        PROTECTED
        PACKAGE
    }

    // TypeExpr represents a type reference
    type TypeExpr {
        pos: Position
        name: String
        nullable: Boolean
        array: Boolean
    }

    // RelationDecl represents a relationship definition
    type RelationDecl {
        pos: Position
        kind: RelationKind
        target: String
        targetType: String?
        alias: String?
        annotations: AnnotationDecl[]
    }

    // RelationKind enumeration for relationship types
    enum RelationKind {
        DEPENDS_ON
        EXTENDS
        IMPLEMENTS
        CONTAINS
        AGGREGATES
    }

    // InterfaceDecl represents an interface definition
    type InterfaceDecl {
        pos: Position
        name: String
        annotations: AnnotationDecl[]
        methods: MethodDecl[]
    }

    // MethodDecl represents a method definition
    type MethodDecl {
        pos: Position
        name: String
        params: ParamDecl[]
        returnType: TypeExpr?
        throwsList: String[]
        isAsync: Boolean
        annotations: AnnotationDecl[]
    }

    // ParamDecl represents a parameter definition
    type ParamDecl {
        pos: Position
        name: String
        paramType: TypeExpr
    }
}
