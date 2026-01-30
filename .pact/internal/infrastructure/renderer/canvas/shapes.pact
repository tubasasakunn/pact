// Shapes - Extended shape drawing methods for Canvas
// Provides higher-level shape primitives for diagram rendering

@layer("infrastructure")
@package("canvas")
component Shapes {
    depends on Canvas

    provides ShapeAPI {
        // Draws a diamond (rhombus) centered at (cx, cy)
        // Used for decision nodes in flowcharts
        Diamond(cx: int, cy: int, width: int, height: int, opts: Option[])

        // Draws a line with an arrowhead at the end point
        // Used for directed edges and transitions
        Arrow(x1: int, y1: int, x2: int, y2: int, opts: Option[])

        // Draws a stadium shape (rectangle with fully rounded ends)
        // Used for terminal nodes in flowcharts
        Stadium(x: int, y: int, width: int, height: int, opts: Option[])

        // Draws a cylinder shape (3D database representation)
        // Used for database nodes in diagrams
        Cylinder(x: int, y: int, width: int, height: int, opts: Option[])

        // Draws a parallelogram with specified skew
        // Used for I/O nodes in flowcharts
        Parallelogram(x: int, y: int, width: int, height: int, skew: int, opts: Option[])
    }

    provides MathUtil {
        // Calculates square root using Newton-Raphson method
        // Used internally for arrow direction calculations
        sqrt(x: float) -> float
    }
}
