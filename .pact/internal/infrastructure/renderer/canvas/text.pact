// Text - Text measurement and wrapping utilities
// Provides functions for calculating text dimensions and line wrapping

@layer("infrastructure")
@package("canvas")
component Text {
    // TextDimensions represents width and height measurements
    type TextDimensions {
        width: int
        height: int
    }

    provides TextMeasurement {
        // Measures the approximate dimensions of text at a given font size
        // Uses heuristic: character width = fontSize * 0.6
        // Returns TextDimensions with width and height in pixels
        MeasureText(text: string, fontSize: int) -> TextDimensions
    }

    provides TextWrapping {
        // Wraps text to fit within a maximum width
        // Attempts to break at word boundaries (spaces)
        // Returns array of lines
        WrapText(text: string, maxWidth: int, fontSize: int) -> string[]
    }
}
