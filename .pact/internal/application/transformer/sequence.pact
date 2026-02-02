// SequenceTransformer transforms AST to sequence diagrams
component SequenceTransformer {
    type SequenceOptions {
        flowName: string
        includeReturn: bool
    }

    provides Transformer {
        Transform(files: SpecFile[], opts: SequenceOptions) -> Diagram
    }
}
