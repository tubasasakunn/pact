// ClassTransformer transforms AST to class diagrams
component ClassTransformer {
    type TransformOptions {
        filterComponents: string[]
    }

    provides Transformer {
        Transform(files: SpecFile[], opts: TransformOptions) -> Diagram
    }
}
