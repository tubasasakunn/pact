// Pattern: Inheritance Tree with 4 children
// 親クラスを4つの子クラスが継承するパターン

component Shape {
    type ShapeData {
        id: string
        color: string
        filled: bool
    }

    provides ShapeAPI {
        Draw() -> bool
        GetArea() -> float
        GetPerimeter() -> float
    }
}

component Circle {
    extends Shape

    type CircleData {
        radius: float
        centerX: float
        centerY: float
    }

    provides CircleAPI {
        GetRadius() -> float
        GetDiameter() -> float
    }
}

component Rectangle {
    extends Shape

    type RectangleData {
        width: float
        height: float
    }

    provides RectangleAPI {
        GetWidth() -> float
        GetHeight() -> float
    }
}

component Triangle {
    extends Shape

    type TriangleData {
        sideA: float
        sideB: float
        sideC: float
    }

    provides TriangleAPI {
        GetSides() -> string
        IsEquilateral() -> bool
    }
}

component Polygon {
    extends Shape

    type PolygonData {
        numSides: int
        sideLength: float
    }

    provides PolygonAPI {
        GetNumSides() -> int
        IsRegular() -> bool
    }
}
