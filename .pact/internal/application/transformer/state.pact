// StateTransformer transforms AST to state diagrams
component StateTransformer {
    type StateOptions {
        statesName: string
    }

    provides Transformer {
        Transform(files: SpecFile[], opts: StateOptions) -> Diagram
    }
}
