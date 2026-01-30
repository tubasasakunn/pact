// Public API - Pact Client
@version("1.0")
@package("pact")
component Client {
    depends on Parser : ParserAPI as parser
    depends on ClassTransformer : ClassTransformerAPI as classTransformer
    depends on SequenceTransformer : SequenceTransformerAPI as seqTransformer
    depends on StateTransformer : StateTransformerAPI as stateTransformer
    depends on FlowTransformer : FlowTransformerAPI as flowTransformer
    depends on ClassRenderer : ClassRendererAPI as classRenderer
    depends on SequenceRenderer : SequenceRendererAPI as seqRenderer
    depends on StateRenderer : StateRendererAPI as stateRenderer
    depends on FlowRenderer : FlowRendererAPI as flowRenderer

    provides ClientAPI {
        New() -> Client
        ParseFile(path: string) -> SpecFile throws ParseError
        ParseString(content: string) -> SpecFile throws ParseError
        ToClassDiagram(spec: SpecFile) -> ClassDiagram throws TransformError
        ToSequenceDiagram(spec: SpecFile, flowName: string) -> SequenceDiagram throws TransformError
        ToStateDiagram(spec: SpecFile, statesName: string) -> StateDiagram throws TransformError
        ToFlowchart(spec: SpecFile, flowName: string) -> FlowDiagram throws TransformError
        RenderClassDiagram(diagram: ClassDiagram, w: Writer)
        RenderSequenceDiagram(diagram: SequenceDiagram, w: Writer)
        RenderStateDiagram(diagram: StateDiagram, w: Writer)
        RenderFlowchart(diagram: FlowDiagram, w: Writer)
    }

    // Type aliases re-exported from internal packages
    type SpecFile {
        imports: ImportDecl[]
        mainComponent: ComponentDecl?
        components: ComponentDecl[]
        annotations: AnnotationDecl[]
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

    type RelationDecl {
        kind: string
        target: string
        interfaceType: string?
        alias: string?
    }

    type FlowDecl {
        name: string
        steps: Step[]
    }

    type StatesDecl {
        name: string
        initialState: string?
        finalStates: string[]
    }
}
