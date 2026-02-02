package ast

// Visitor はASTノードを走査するインターフェース
type Visitor interface {
	VisitSpecFile(node *SpecFile) error
	VisitImportDecl(node *ImportDecl) error
	VisitComponentDecl(node *ComponentDecl) error
	VisitAnnotationDecl(node *AnnotationDecl) error
	VisitTypeDecl(node *TypeDecl) error
	VisitFieldDecl(node *FieldDecl) error
	VisitRelationDecl(node *RelationDecl) error
	VisitInterfaceDecl(node *InterfaceDecl) error
	VisitMethodDecl(node *MethodDecl) error
	VisitFlowDecl(node *FlowDecl) error
	VisitStep(node Step) error
	VisitExpr(node Expr) error
	VisitStatesDecl(node *StatesDecl) error
	VisitStateDecl(node *StateDecl) error
	VisitTransitionDecl(node *TransitionDecl) error
}

// BaseVisitor はVisitorのデフォルト実装
type BaseVisitor struct{}

func (v *BaseVisitor) VisitSpecFile(node *SpecFile) error         { return nil }
func (v *BaseVisitor) VisitImportDecl(node *ImportDecl) error     { return nil }
func (v *BaseVisitor) VisitComponentDecl(node *ComponentDecl) error { return nil }
func (v *BaseVisitor) VisitAnnotationDecl(node *AnnotationDecl) error { return nil }
func (v *BaseVisitor) VisitTypeDecl(node *TypeDecl) error         { return nil }
func (v *BaseVisitor) VisitFieldDecl(node *FieldDecl) error       { return nil }
func (v *BaseVisitor) VisitRelationDecl(node *RelationDecl) error { return nil }
func (v *BaseVisitor) VisitInterfaceDecl(node *InterfaceDecl) error { return nil }
func (v *BaseVisitor) VisitMethodDecl(node *MethodDecl) error     { return nil }
func (v *BaseVisitor) VisitFlowDecl(node *FlowDecl) error         { return nil }
func (v *BaseVisitor) VisitStep(node Step) error                  { return nil }
func (v *BaseVisitor) VisitExpr(node Expr) error                  { return nil }
func (v *BaseVisitor) VisitStatesDecl(node *StatesDecl) error     { return nil }
func (v *BaseVisitor) VisitStateDecl(node *StateDecl) error       { return nil }
func (v *BaseVisitor) VisitTransitionDecl(node *TransitionDecl) error { return nil }

// Walk はASTを走査する
func Walk(v Visitor, node *SpecFile) error {
	if err := v.VisitSpecFile(node); err != nil {
		return err
	}

	for i := range node.Imports {
		if err := v.VisitImportDecl(&node.Imports[i]); err != nil {
			return err
		}
	}

	// 全コンポーネントを走査
	// 注: node.Component は node.Components の最後の要素を指すため、
	// Components のみを走査すれば全コンポーネントをカバーできる
	for i := range node.Components {
		if err := walkComponent(v, &node.Components[i]); err != nil {
			return err
		}
	}

	// Components が空で Component のみ設定されている場合（後方互換性）
	if len(node.Components) == 0 && node.Component != nil {
		if err := walkComponent(v, node.Component); err != nil {
			return err
		}
	}

	return nil
}

func walkComponent(v Visitor, node *ComponentDecl) error {
	if err := v.VisitComponentDecl(node); err != nil {
		return err
	}

	for i := range node.Annotations {
		if err := v.VisitAnnotationDecl(&node.Annotations[i]); err != nil {
			return err
		}
	}

	// Types
	for i := range node.Body.Types {
		if err := walkType(v, &node.Body.Types[i]); err != nil {
			return err
		}
	}

	// Relations
	for i := range node.Body.Relations {
		if err := v.VisitRelationDecl(&node.Body.Relations[i]); err != nil {
			return err
		}
	}

	// Provides
	for i := range node.Body.Provides {
		if err := walkInterface(v, &node.Body.Provides[i]); err != nil {
			return err
		}
	}

	// Requires
	for i := range node.Body.Requires {
		if err := walkInterface(v, &node.Body.Requires[i]); err != nil {
			return err
		}
	}

	// Flows
	for i := range node.Body.Flows {
		if err := walkFlow(v, &node.Body.Flows[i]); err != nil {
			return err
		}
	}

	// States
	for i := range node.Body.States {
		if err := walkStates(v, &node.Body.States[i]); err != nil {
			return err
		}
	}

	return nil
}

func walkType(v Visitor, node *TypeDecl) error {
	if err := v.VisitTypeDecl(node); err != nil {
		return err
	}

	for i := range node.Annotations {
		if err := v.VisitAnnotationDecl(&node.Annotations[i]); err != nil {
			return err
		}
	}

	for i := range node.Fields {
		if err := v.VisitFieldDecl(&node.Fields[i]); err != nil {
			return err
		}
	}

	return nil
}

func walkInterface(v Visitor, node *InterfaceDecl) error {
	if err := v.VisitInterfaceDecl(node); err != nil {
		return err
	}

	for i := range node.Methods {
		if err := v.VisitMethodDecl(&node.Methods[i]); err != nil {
			return err
		}
	}

	return nil
}

func walkFlow(v Visitor, node *FlowDecl) error {
	if err := v.VisitFlowDecl(node); err != nil {
		return err
	}

	for _, step := range node.Steps {
		if err := walkStep(v, step); err != nil {
			return err
		}
	}

	return nil
}

func walkStep(v Visitor, step Step) error {
	if err := v.VisitStep(step); err != nil {
		return err
	}

	switch s := step.(type) {
	case *IfStep:
		for _, inner := range s.Then {
			if err := walkStep(v, inner); err != nil {
				return err
			}
		}
		for _, inner := range s.Else {
			if err := walkStep(v, inner); err != nil {
				return err
			}
		}
	case *ForStep:
		for _, inner := range s.Body {
			if err := walkStep(v, inner); err != nil {
				return err
			}
		}
	case *WhileStep:
		for _, inner := range s.Body {
			if err := walkStep(v, inner); err != nil {
				return err
			}
		}
	}

	return nil
}

func walkStates(v Visitor, node *StatesDecl) error {
	if err := v.VisitStatesDecl(node); err != nil {
		return err
	}

	for i := range node.States {
		if err := walkState(v, &node.States[i]); err != nil {
			return err
		}
	}

	for i := range node.Transitions {
		if err := v.VisitTransitionDecl(&node.Transitions[i]); err != nil {
			return err
		}
	}

	return nil
}

func walkState(v Visitor, node *StateDecl) error {
	if err := v.VisitStateDecl(node); err != nil {
		return err
	}

	for i := range node.States {
		if err := walkState(v, &node.States[i]); err != nil {
			return err
		}
	}

	for i := range node.Transitions {
		if err := v.VisitTransitionDecl(&node.Transitions[i]); err != nil {
			return err
		}
	}

	return nil
}
