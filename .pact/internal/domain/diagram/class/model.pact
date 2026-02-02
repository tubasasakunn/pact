// Class diagram domain model
component ClassDiagram {
    // Visibility represents the visibility of class members
    enum Visibility {
        public
        private
        protected
        package
    }

    // EdgeType represents the type of edge/relationship between classes
    enum EdgeType {
        dependency
        inheritance
        implementation
        composition
        aggregation
    }

    // Decoration represents edge decoration style
    enum Decoration {
        none
        arrow
        triangle
        filled_diamond
        empty_diamond
    }

    // LineStyle represents the line style for edges
    enum LineStyle {
        solid
        dashed
    }

    // Param represents a method parameter
    type Param {
        name: string
        paramType: string
    }

    // Attribute represents a class attribute
    type Attribute {
        name: string
        attrType: string
        visibility: Visibility
    }

    // Method represents a class method
    type Method {
        name: string
        params: Param[]
        returnType: string
        visibility: Visibility
        isAsync: bool
    }

    // Node represents a class node in the diagram
    type Node {
        id: string
        name: string
        stereotype: string
        attributes: Attribute[]
        methods: Method[]
    }

    // Edge represents a relationship between classes
    type Edge {
        fromNode: string
        toNode: string
        edgeType: EdgeType
        label: string
        decoration: Decoration
        lineStyle: LineStyle
    }

    // Diagram represents a class diagram
    type Diagram {
        nodes: Node[]
        edges: Edge[]
    }
}
