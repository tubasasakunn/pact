package ast

import (
	"errors"
	"testing"
)

// CountingVisitor は各ノードの訪問回数をカウントする
type CountingVisitor struct {
	BaseVisitor
	SpecFileCount     int
	ImportCount       int
	ComponentCount    int
	AnnotationCount   int
	TypeCount         int
	FieldCount        int
	RelationCount     int
	InterfaceCount    int
	MethodCount       int
	FlowCount         int
	StepCount         int
	ExprCount         int
	StatesCount       int
	StateCount        int
	TransitionCount   int
}

func (v *CountingVisitor) VisitSpecFile(node *SpecFile) error {
	v.SpecFileCount++
	return nil
}

func (v *CountingVisitor) VisitImportDecl(node *ImportDecl) error {
	v.ImportCount++
	return nil
}

func (v *CountingVisitor) VisitComponentDecl(node *ComponentDecl) error {
	v.ComponentCount++
	return nil
}

func (v *CountingVisitor) VisitAnnotationDecl(node *AnnotationDecl) error {
	v.AnnotationCount++
	return nil
}

func (v *CountingVisitor) VisitTypeDecl(node *TypeDecl) error {
	v.TypeCount++
	return nil
}

func (v *CountingVisitor) VisitFieldDecl(node *FieldDecl) error {
	v.FieldCount++
	return nil
}

func (v *CountingVisitor) VisitRelationDecl(node *RelationDecl) error {
	v.RelationCount++
	return nil
}

func (v *CountingVisitor) VisitInterfaceDecl(node *InterfaceDecl) error {
	v.InterfaceCount++
	return nil
}

func (v *CountingVisitor) VisitMethodDecl(node *MethodDecl) error {
	v.MethodCount++
	return nil
}

func (v *CountingVisitor) VisitFlowDecl(node *FlowDecl) error {
	v.FlowCount++
	return nil
}

func (v *CountingVisitor) VisitStep(node Step) error {
	v.StepCount++
	return nil
}

func (v *CountingVisitor) VisitExpr(node Expr) error {
	v.ExprCount++
	return nil
}

func (v *CountingVisitor) VisitStatesDecl(node *StatesDecl) error {
	v.StatesCount++
	return nil
}

func (v *CountingVisitor) VisitStateDecl(node *StateDecl) error {
	v.StateCount++
	return nil
}

func (v *CountingVisitor) VisitTransitionDecl(node *TransitionDecl) error {
	v.TransitionCount++
	return nil
}

// =============================================================================
// AV001-AV004: Visitor Tests
// =============================================================================

// AV001: SpecFile訪問
func TestVisitor_Walk_SpecFile(t *testing.T) {
	spec := &SpecFile{
		Path: "test.pact",
	}

	v := &CountingVisitor{}
	if err := Walk(v, spec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v.SpecFileCount != 1 {
		t.Errorf("expected SpecFileCount 1, got %d", v.SpecFileCount)
	}
}

// AV002: 全ノード走査
func TestVisitor_Walk_AllNodes(t *testing.T) {
	spec := &SpecFile{
		Path: "test.pact",
		Imports: []ImportDecl{
			{Path: "./a.pact"},
			{Path: "./b.pact"},
		},
		Component: &ComponentDecl{
			Name: "TestComponent",
			Annotations: []AnnotationDecl{
				{Name: "description"},
			},
			Body: ComponentBody{
				Types: []TypeDecl{
					{
						Name: "User",
						Kind: TypeKindStruct,
						Fields: []FieldDecl{
							{Name: "id"},
							{Name: "name"},
						},
					},
				},
				Relations: []RelationDecl{
					{Kind: RelationDependsOn, Target: "Foo"},
				},
				Provides: []InterfaceDecl{
					{
						Name: "API",
						Methods: []MethodDecl{
							{Name: "Get"},
							{Name: "Set"},
						},
					},
				},
				Flows: []FlowDecl{
					{
						Name: "Process",
						Steps: []Step{
							&CallStep{},
							&IfStep{
								Then: []Step{&CallStep{}},
								Else: []Step{&ReturnStep{}},
							},
						},
					},
				},
				States: []StatesDecl{
					{
						Name:    "OrderState",
						Initial: "Pending",
						States: []StateDecl{
							{Name: "Pending"},
							{Name: "Processing"},
						},
						Transitions: []TransitionDecl{
							{From: "Pending", To: "Processing"},
						},
					},
				},
			},
		},
	}

	v := &CountingVisitor{}
	if err := Walk(v, spec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v.SpecFileCount != 1 {
		t.Errorf("expected SpecFileCount 1, got %d", v.SpecFileCount)
	}
	if v.ImportCount != 2 {
		t.Errorf("expected ImportCount 2, got %d", v.ImportCount)
	}
	if v.ComponentCount != 1 {
		t.Errorf("expected ComponentCount 1, got %d", v.ComponentCount)
	}
	if v.AnnotationCount != 1 {
		t.Errorf("expected AnnotationCount 1, got %d", v.AnnotationCount)
	}
	if v.TypeCount != 1 {
		t.Errorf("expected TypeCount 1, got %d", v.TypeCount)
	}
	if v.FieldCount != 2 {
		t.Errorf("expected FieldCount 2, got %d", v.FieldCount)
	}
	if v.RelationCount != 1 {
		t.Errorf("expected RelationCount 1, got %d", v.RelationCount)
	}
	if v.InterfaceCount != 1 {
		t.Errorf("expected InterfaceCount 1, got %d", v.InterfaceCount)
	}
	if v.MethodCount != 2 {
		t.Errorf("expected MethodCount 2, got %d", v.MethodCount)
	}
	if v.FlowCount != 1 {
		t.Errorf("expected FlowCount 1, got %d", v.FlowCount)
	}
	// IfStepの中のステップも含めて4つ
	if v.StepCount != 4 {
		t.Errorf("expected StepCount 4, got %d", v.StepCount)
	}
	if v.StatesCount != 1 {
		t.Errorf("expected StatesCount 1, got %d", v.StatesCount)
	}
	if v.StateCount != 2 {
		t.Errorf("expected StateCount 2, got %d", v.StateCount)
	}
	if v.TransitionCount != 1 {
		t.Errorf("expected TransitionCount 1, got %d", v.TransitionCount)
	}
}

// OrderTrackingVisitor は訪問順を記録する
type OrderTrackingVisitor struct {
	BaseVisitor
	Order []string
}

func (v *OrderTrackingVisitor) VisitSpecFile(node *SpecFile) error {
	v.Order = append(v.Order, "SpecFile")
	return nil
}

func (v *OrderTrackingVisitor) VisitImportDecl(node *ImportDecl) error {
	v.Order = append(v.Order, "Import:"+node.Path)
	return nil
}

func (v *OrderTrackingVisitor) VisitComponentDecl(node *ComponentDecl) error {
	v.Order = append(v.Order, "Component:"+node.Name)
	return nil
}

func (v *OrderTrackingVisitor) VisitTypeDecl(node *TypeDecl) error {
	v.Order = append(v.Order, "Type:"+node.Name)
	return nil
}

func (v *OrderTrackingVisitor) VisitFieldDecl(node *FieldDecl) error {
	v.Order = append(v.Order, "Field:"+node.Name)
	return nil
}

// AV003: 走査順序（深さ優先）
func TestVisitor_Walk_Order(t *testing.T) {
	spec := &SpecFile{
		Path: "test.pact",
		Imports: []ImportDecl{
			{Path: "./a.pact"},
		},
		Component: &ComponentDecl{
			Name: "Foo",
			Body: ComponentBody{
				Types: []TypeDecl{
					{
						Name: "Bar",
						Kind: TypeKindStruct,
						Fields: []FieldDecl{
							{Name: "x"},
							{Name: "y"},
						},
					},
				},
			},
		},
	}

	v := &OrderTrackingVisitor{}
	if err := Walk(v, spec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"SpecFile",
		"Import:./a.pact",
		"Component:Foo",
		"Type:Bar",
		"Field:x",
		"Field:y",
	}

	if len(v.Order) != len(expected) {
		t.Fatalf("expected %d visits, got %d: %v", len(expected), len(v.Order), v.Order)
	}

	for i, exp := range expected {
		if v.Order[i] != exp {
			t.Errorf("expected order[%d]=%q, got %q", i, exp, v.Order[i])
		}
	}
}

// AV004: BaseVisitorデフォルト実装
func TestVisitor_BaseVisitor(t *testing.T) {
	spec := &SpecFile{
		Path: "test.pact",
		Imports: []ImportDecl{
			{Path: "./a.pact"},
		},
		Component: &ComponentDecl{
			Name: "Foo",
		},
	}

	v := &BaseVisitor{}
	if err := Walk(v, spec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// エラーなく完了することを確認
}

// ErrorVisitor は特定のノードでエラーを返す
type ErrorVisitor struct {
	BaseVisitor
	ErrorOn string
}

func (v *ErrorVisitor) VisitTypeDecl(node *TypeDecl) error {
	if node.Name == v.ErrorOn {
		return errors.New("test error")
	}
	return nil
}

// AV005: エラー伝播
func TestVisitor_Walk_ErrorPropagation(t *testing.T) {
	spec := &SpecFile{
		Path: "test.pact",
		Component: &ComponentDecl{
			Name: "Foo",
			Body: ComponentBody{
				Types: []TypeDecl{
					{Name: "First"},
					{Name: "Second"},
					{Name: "Third"},
				},
			},
		},
	}

	v := &ErrorVisitor{ErrorOn: "Second"}
	err := Walk(v, spec)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "test error" {
		t.Errorf("expected 'test error', got %q", err.Error())
	}
}
