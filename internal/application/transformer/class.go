package transformer

import (
	"pact/internal/domain/ast"
	"pact/internal/domain/diagram/class"
	"pact/internal/domain/diagram/common"
)

// ClassTransformer はASTをクラス図に変換する
type ClassTransformer struct{}

// NewClassTransformer は新しいClassTransformerを作成する
func NewClassTransformer() *ClassTransformer {
	return &ClassTransformer{}
}

// TransformOptions は変換オプション
type TransformOptions struct {
	FilterComponents []string
}

// Transform はASTをクラス図に変換する
func (t *ClassTransformer) Transform(files []*ast.SpecFile, opts *TransformOptions) (*class.Diagram, error) {
	diagram := &class.Diagram{
		Nodes: []class.Node{},
		Edges: []class.Edge{},
	}

	for _, file := range files {
		if file.Component == nil {
			continue
		}

		if opts != nil && len(opts.FilterComponents) > 0 {
			if !contains(opts.FilterComponents, file.Component.Name) {
				continue
			}
		}

		// コンポーネントをノードに変換
		node := t.transformComponent(file.Component)
		diagram.Nodes = append(diagram.Nodes, node)

		// 型をノードに変換
		for _, typ := range file.Component.Body.Types {
			typeNode := t.transformType(&typ)
			diagram.Nodes = append(diagram.Nodes, typeNode)
		}

		// 関係をエッジに変換
		for _, rel := range file.Component.Body.Relations {
			edge := t.transformRelation(file.Component.Name, &rel)
			diagram.Edges = append(diagram.Edges, edge)
		}

		// インターフェースをノードに変換
		for _, iface := range file.Component.Body.Requires {
			ifaceNode := t.transformInterface(&iface)
			diagram.Nodes = append(diagram.Nodes, ifaceNode)
		}
	}

	return diagram, nil
}

func (t *ClassTransformer) transformComponent(comp *ast.ComponentDecl) class.Node {
	node := class.Node{
		ID:         comp.Name,
		Name:       comp.Name,
		Stereotype: "component",
		Methods:    []class.Method{},
	}

	// providesからメソッドを収集
	for _, iface := range comp.Body.Provides {
		for _, method := range iface.Methods {
			node.Methods = append(node.Methods, t.transformMethod(&method))
		}
	}

	// アノテーションを変換
	for _, ann := range comp.Annotations {
		node.Annotations = append(node.Annotations, transformAnnotation(&ann))
	}

	return node
}

func (t *ClassTransformer) transformType(typ *ast.TypeDecl) class.Node {
	node := class.Node{
		ID:         typ.Name,
		Name:       typ.Name,
		Attributes: []class.Attribute{},
		Methods:    []class.Method{},
	}

	if typ.Kind == ast.TypeKindEnum {
		node.Stereotype = "enum"
	}

	for _, field := range typ.Fields {
		node.Attributes = append(node.Attributes, class.Attribute{
			Name:       field.Name,
			Type:       field.Type.Name,
			Visibility: convertVisibility(field.Visibility),
		})
	}

	for _, ann := range typ.Annotations {
		node.Annotations = append(node.Annotations, transformAnnotation(&ann))
	}

	return node
}

func (t *ClassTransformer) transformMethod(method *ast.MethodDecl) class.Method {
	m := class.Method{
		Name:       method.Name,
		Visibility: class.VisibilityPublic,
		Async:      method.Async,
	}

	if method.ReturnType != nil {
		m.ReturnType = method.ReturnType.Name
	}

	for _, param := range method.Params {
		m.Params = append(m.Params, class.Param{
			Name: param.Name,
			Type: param.Type.Name,
		})
	}

	return m
}

func (t *ClassTransformer) transformRelation(fromName string, rel *ast.RelationDecl) class.Edge {
	edge := class.Edge{
		From: fromName,
		To:   rel.Target,
	}

	switch rel.Kind {
	case ast.RelationDependsOn:
		edge.Type = class.EdgeTypeDependency
		edge.Decoration = class.DecorationArrow
		edge.LineStyle = class.LineStyleDashed
	case ast.RelationExtends:
		edge.Type = class.EdgeTypeInheritance
		edge.Decoration = class.DecorationTriangle
		edge.LineStyle = class.LineStyleSolid
	case ast.RelationImplements:
		edge.Type = class.EdgeTypeImplementation
		edge.Decoration = class.DecorationTriangle
		edge.LineStyle = class.LineStyleDashed
	case ast.RelationContains:
		edge.Type = class.EdgeTypeComposition
		edge.Decoration = class.DecorationFilledDiamond
		edge.LineStyle = class.LineStyleSolid
	case ast.RelationAggregates:
		edge.Type = class.EdgeTypeAggregation
		edge.Decoration = class.DecorationEmptyDiamond
		edge.LineStyle = class.LineStyleSolid
	}

	return edge
}

func (t *ClassTransformer) transformInterface(iface *ast.InterfaceDecl) class.Node {
	node := class.Node{
		ID:         iface.Name,
		Name:       iface.Name,
		Stereotype: "interface",
		Methods:    []class.Method{},
	}

	for _, method := range iface.Methods {
		node.Methods = append(node.Methods, t.transformMethod(&method))
	}

	return node
}

func convertVisibility(v ast.Visibility) class.Visibility {
	switch v {
	case ast.VisibilityPublic:
		return class.VisibilityPublic
	case ast.VisibilityPrivate:
		return class.VisibilityPrivate
	case ast.VisibilityProtected:
		return class.VisibilityProtected
	case ast.VisibilityPackage:
		return class.VisibilityPackage
	default:
		return class.VisibilityPublic
	}
}

func transformAnnotation(ann *ast.AnnotationDecl) common.Annotation {
	result := common.Annotation{
		Name: ann.Name,
		Args: make(map[string]string),
	}
	for _, arg := range ann.Args {
		if arg.Key != nil {
			result.Args[*arg.Key] = arg.Value
		}
	}
	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
