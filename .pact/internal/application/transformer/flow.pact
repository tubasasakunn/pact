// FlowTransformer transforms AST to flow diagrams
component FlowTransformer {
    type FlowOptions {
        flowName: string
        includeSwimlanes: bool
    }

    provides Transformer {
        Transform(files: SpecFile[], opts: FlowOptions) -> Diagram
    }
}
