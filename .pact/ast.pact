// AST Domain - Abstract Syntax Tree Types
@version("1.0")
@package("ast")
component AST {
    // Root node
    type SpecFile {
        imports: ImportDecl[]
        mainComponent: ComponentDecl?
        components: ComponentDecl[]
        annotations: AnnotationDecl[]
    }

    type ImportDecl {
        path: string
        alias: string?
    }

    type ComponentDecl {
        name: string
        body: ComponentBody
        annotations: AnnotationDecl[]
    }

    type ComponentBody {
        types: TypeDecl[]
        relations: RelationDecl[]
        interfaces: InterfaceDecl[]
        flows: FlowDecl[]
        stateBlocks: StatesDecl[]
    }

    // Type declarations
    type TypeDecl {
        name: string
        kind: TypeKind
        fields: FieldDecl[]
        values: string[]
        annotations: AnnotationDecl[]
    }

    enum TypeKind {
        STRUCT
        ENUM
    }

    type FieldDecl {
        name: string
        fieldType: TypeExpr
        visibility: Visibility
        annotations: AnnotationDecl[]
    }

    type TypeExpr {
        name: string
        isNullable: bool
        isArray: bool
    }

    enum Visibility {
        PUBLIC
        PRIVATE
        PROTECTED
        PACKAGE
    }

    // Relations
    type RelationDecl {
        kind: RelationKind
        target: string
        interfaceType: string?
        alias: string?
        annotations: AnnotationDecl[]
    }

    enum RelationKind {
        DEPENDS_ON
        EXTENDS
        IMPLEMENTS
        CONTAINS
        AGGREGATES
    }

    // Interfaces
    type InterfaceDecl {
        name: string
        kind: InterfaceKind
        methods: MethodDecl[]
        annotations: AnnotationDecl[]
    }

    enum InterfaceKind {
        PROVIDES
        REQUIRES
    }

    type MethodDecl {
        name: string
        params: ParamDecl[]
        returnType: TypeExpr?
        isAsync: bool
        throwsList: string[]
        annotations: AnnotationDecl[]
    }

    type ParamDecl {
        name: string
        paramType: TypeExpr
    }

    // Flows
    type FlowDecl {
        name: string
        steps: Step[]
        annotations: AnnotationDecl[]
    }

    type AssignStep {
        variable: string
        value: Expr
    }

    type CallStep {
        expr: CallExpr
    }

    type ReturnStep {
        value: Expr?
    }

    type ThrowStep {
        errorType: string
    }

    type IfStep {
        condition: Expr
        thenSteps: Step[]
        elseSteps: Step[]
    }

    type ForStep {
        variable: string
        iterable: Expr
        body: Step[]
    }

    type WhileStep {
        condition: Expr
        body: Step[]
    }

    // Expressions
    type LiteralExpr {
        value: string
        kind: LiteralKind
    }

    enum LiteralKind {
        STRING
        INT
        FLOAT
        BOOL
        NULL
    }

    type VariableExpr {
        name: string
    }

    type FieldExpr {
        object: Expr
        field: string
    }

    type CallExpr {
        object: Expr
        methodName: string
        args: Expr[]
    }

    type BinaryExpr {
        left: Expr
        operator: string
        right: Expr
    }

    type UnaryExpr {
        operator: string
        operand: Expr
    }

    type TernaryExpr {
        condition: Expr
        consequent: Expr
        alternate: Expr
    }

    type NullishExpr {
        left: Expr
        right: Expr
        throwOnNull: bool
        errorType: string?
    }

    // State machines
    type StatesDecl {
        name: string
        initialState: string?
        finalStates: string[]
        stateDecls: StateDecl[]
        transitions: TransitionDecl[]
        parallels: ParallelDecl[]
        annotations: AnnotationDecl[]
    }

    type StateDecl {
        name: string
        entryActions: string[]
        exitActions: string[]
        nested: StateDecl[]
        annotations: AnnotationDecl[]
    }

    type TransitionDecl {
        fromState: string
        toState: string
        trigger: Trigger?
        guard: Expr?
        actions: string[]
    }

    type EventTrigger {
        eventName: string
    }

    type AfterTrigger {
        duration: Duration
    }

    type WhenTrigger {
        condition: Expr
    }

    type Duration {
        value: int
        unit: DurationUnit
    }

    enum DurationUnit {
        MS
        S
        M
        H
        D
    }

    type ParallelDecl {
        name: string
        regions: RegionDecl[]
    }

    type RegionDecl {
        name: string
        stateDecls: StateDecl[]
        transitions: TransitionDecl[]
    }

    // Annotations
    type AnnotationDecl {
        name: string
        args: AnnotationArg[]
    }

    type AnnotationArg {
        key: string?
        value: string
    }
}
