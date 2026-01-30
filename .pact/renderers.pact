// Renderer Components
@version("1.0")
@package("renderer")
component Canvas {
    type CanvasState {
        width: int
        height: int
        elements: string[]
        defs: string[]
    }

    type StyleOption {
        fill: string?
        stroke: string?
        strokeWidth: int?
        className: string?
    }

    provides CanvasAPI {
        New() -> Canvas
        SetSize(width: int, height: int)
        Rect(x: int, y: int, w: int, h: int, opts: StyleOption?)
        RoundRect(x: int, y: int, w: int, h: int, rx: int, ry: int, opts: StyleOption?)
        Circle(x: int, y: int, r: int, opts: StyleOption?)
        Ellipse(x: int, y: int, rx: int, ry: int, opts: StyleOption?)
        Line(x1: int, y1: int, x2: int, y2: int, opts: StyleOption?)
        Path(d: string, opts: StyleOption?)
        Polygon(points: string, opts: StyleOption?)
        Text(x: int, y: int, text: string, opts: StyleOption?)
        Diamond(x: int, y: int, w: int, h: int, opts: StyleOption?)
        Arrow(x1: int, y1: int, x2: int, y2: int, opts: StyleOption?)
        Stadium(x: int, y: int, w: int, h: int, opts: StyleOption?)
        Cylinder(x: int, y: int, w: int, h: int, opts: StyleOption?)
        Parallelogram(x: int, y: int, w: int, h: int, opts: StyleOption?)
        AddDef(def: string)
        WriteTo(w: Writer) -> int
        String() -> string
    }

    provides TextUtils {
        MeasureText(text: string, fontSize: int) -> TextMetrics
        WrapText(text: string, maxWidth: int, fontSize: int) -> string[]
    }

    type TextMetrics {
        width: int
        height: int
    }
}

component ClassRenderer {
    depends on Canvas : CanvasAPI as canvas

    provides ClassRendererAPI {
        NewClassRenderer() -> ClassRenderer
        Render(diagram: ClassDiagram, w: Writer)
    }
}

component SequenceRenderer {
    depends on Canvas : CanvasAPI as canvas

    provides SequenceRendererAPI {
        NewSequenceRenderer() -> SequenceRenderer
        Render(diagram: SequenceDiagram, w: Writer)
    }
}

component StateRenderer {
    depends on Canvas : CanvasAPI as canvas

    provides StateRendererAPI {
        NewStateRenderer() -> StateRenderer
        Render(diagram: StateDiagram, w: Writer)
    }
}

component FlowRenderer {
    depends on Canvas : CanvasAPI as canvas

    provides FlowRendererAPI {
        NewFlowRenderer() -> FlowRenderer
        Render(diagram: FlowDiagram, w: Writer)
    }
}
