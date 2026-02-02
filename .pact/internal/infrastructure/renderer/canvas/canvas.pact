// Canvas - SVG drawing canvas for diagram rendering
// Provides a fluent API for building SVG elements with configurable options

@layer("infrastructure")
@package("canvas")
component Canvas {
    type Canvas {
        width: int
        height: int
        elements: string[]
        defs: string[]
    }

    type Option {
        // Function type that modifies attribute map
        // func(attrs map[string]string)
    }

    provides CanvasFactory {
        // Creates a new Canvas with default size (800x600)
        New() -> Canvas
    }

    provides CanvasAPI {
        // Sets the canvas dimensions
        SetSize(width: int, height: int)

        // Draws a rectangle at (x, y) with given dimensions
        Rect(x: int, y: int, width: int, height: int, opts: Option[])

        // Draws a rounded rectangle with corner radii rx, ry
        RoundRect(x: int, y: int, width: int, height: int, rx: int, ry: int, opts: Option[])

        // Draws a circle at center (cx, cy) with radius r
        Circle(cx: int, cy: int, r: int, opts: Option[])

        // Draws an ellipse at center (cx, cy) with radii rx, ry
        Ellipse(cx: int, cy: int, rx: int, ry: int, opts: Option[])

        // Draws a line from (x1, y1) to (x2, y2)
        Line(x1: int, y1: int, x2: int, y2: int, opts: Option[])

        // Draws an SVG path with the given path data
        Path(d: string, opts: Option[])

        // Draws a polygon with the given points string
        Polygon(points: string, opts: Option[])

        // Draws text at position (x, y)
        Text(x: int, y: int, text: string, opts: Option[])

        // Adds a definition to the SVG defs section
        AddDef(def: string)

        // Writes the complete SVG to the given writer
        WriteTo(w: Writer) -> int64

        // Returns the SVG as a string
        String() -> string
    }

    provides OptionFactory {
        // Creates an option that sets fill color
        Fill(color: string) -> Option

        // Creates an option that sets stroke color
        Stroke(color: string) -> Option

        // Creates an option that sets stroke width
        StrokeWidth(width: int) -> Option

        // Creates an option that sets CSS class
        Class(className: string) -> Option
    }
}
