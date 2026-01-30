// Transformer Components
@version("1.0")
@package("transformer")
component ClassTransformer {
    depends on AST : ASTTypes as ast
    depends on ClassDiagram : DiagramBuilder as diagram

    type TransformOptions {
        filterComponents: string[]
    }

    provides ClassTransformerAPI {
        NewClassTransformer() -> ClassTransformer
        Transform(files: SpecFile[], opts: TransformOptions?) -> ClassDiagram
    }
}

component SequenceTransformer {
    depends on AST : ASTTypes as ast
    depends on SequenceDiagram : DiagramBuilder as diagram

    type SequenceOptions {
        flowName: string
        includeReturn: bool
    }

    provides SequenceTransformerAPI {
        NewSequenceTransformer() -> SequenceTransformer
        Transform(files: SpecFile[], opts: SequenceOptions) -> SequenceDiagram
    }
}

component StateTransformer {
    depends on AST : ASTTypes as ast
    depends on StateDiagram : DiagramBuilder as diagram

    type StateOptions {
        statesName: string
    }

    provides StateTransformerAPI {
        NewStateTransformer() -> StateTransformer
        Transform(files: SpecFile[], opts: StateOptions) -> StateDiagram
    }
}

component FlowTransformer {
    depends on AST : ASTTypes as ast
    depends on FlowDiagram : DiagramBuilder as diagram

    type FlowOptions {
        flowName: string
        includeSwimlanes: bool
    }

    provides FlowTransformerAPI {
        NewFlowTransformer() -> FlowTransformer
        Transform(files: SpecFile[], opts: FlowOptions) -> FlowDiagram
    }
}
