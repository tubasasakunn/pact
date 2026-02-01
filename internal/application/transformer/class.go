package transformer

import (
	"fmt"

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
		Notes: []common.Note{},
	}

	noteCounter := 0

	for _, file := range files {
		// ファイル内の全コンポーネントを収集
		components := t.getComponents(file)

		for _, comp := range components {
			if opts != nil && len(opts.FilterComponents) > 0 {
				if !contains(opts.FilterComponents, comp.Name) {
					continue
				}
			}

			// コンポーネントをノードに変換
			node := t.transformComponent(comp)
			diagram.Nodes = append(diagram.Nodes, node)

			// コンポーネントのアノテーションからNotesを抽出
			for _, ann := range comp.Annotations {
				if note := t.extractNote(&ann, comp.Name, &noteCounter); note != nil {
					diagram.Notes = append(diagram.Notes, *note)
				}
			}

			// 型をノードに変換
			for _, typ := range comp.Body.Types {
				typeNode := t.transformType(&typ)
				diagram.Nodes = append(diagram.Nodes, typeNode)

				// 型のアノテーションからNotesを抽出
				for _, ann := range typ.Annotations {
					if note := t.extractNote(&ann, typ.Name, &noteCounter); note != nil {
						diagram.Notes = append(diagram.Notes, *note)
					}
				}
			}

			// 関係をエッジに変換
			for _, rel := range comp.Body.Relations {
				edge := t.transformRelation(comp.Name, &rel)
				diagram.Edges = append(diagram.Edges, edge)
			}

			// インターフェースをノードに変換
			for _, iface := range comp.Body.Requires {
				ifaceNode := t.transformInterface(&iface)
				diagram.Nodes = append(diagram.Nodes, ifaceNode)
			}
		}
	}

	return diagram, nil
}

// getComponents はファイルから全コンポーネントを取得する
func (t *ClassTransformer) getComponents(file *ast.SpecFile) []*ast.ComponentDecl {
	// Componentsがある場合はそれを使用
	if len(file.Components) > 0 {
		result := make([]*ast.ComponentDecl, len(file.Components))
		for i := range file.Components {
			result[i] = &file.Components[i]
		}
		return result
	}
	// 単一Componentの場合
	if file.Component != nil {
		return []*ast.ComponentDecl{file.Component}
	}
	return nil
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
		Throws:     method.Throws,
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

	// Aliasがある場合はラベルとして使用
	if rel.Alias != nil {
		edge.Label = *rel.Alias
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

// extractNote は@noteまたは@descriptionアノテーションからNoteを抽出する
func (t *ClassTransformer) extractNote(ann *ast.AnnotationDecl, attachTo string, counter *int) *common.Note {
	// @note または @description アノテーションを処理
	if ann.Name != "note" && ann.Name != "description" {
		return nil
	}

	// 最初の引数をテキストとして使用
	text := ""
	position := common.NotePositionRight // デフォルト位置

	for _, arg := range ann.Args {
		if arg.Key == nil {
			// 位置指定なしの引数はテキスト
			text = arg.Value
		} else if *arg.Key == "position" {
			switch arg.Value {
			case "left":
				position = common.NotePositionLeft
			case "right":
				position = common.NotePositionRight
			case "top":
				position = common.NotePositionTop
			case "bottom":
				position = common.NotePositionBottom
			}
		} else if *arg.Key == "text" {
			text = arg.Value
		}
	}

	if text == "" {
		return nil
	}

	*counter++
	return &common.Note{
		ID:       fmt.Sprintf("note_%d", *counter),
		Text:     text,
		Position: position,
		AttachTo: attachTo,
	}
}
