package transformer

import (
	"pact/internal/domain/ast"
	"pact/internal/domain/diagram/sequence"
	"pact/internal/domain/errors"
)

// SequenceTransformer はASTをシーケンス図に変換する
type SequenceTransformer struct{}

// NewSequenceTransformer は新しいSequenceTransformerを作成する
func NewSequenceTransformer() *SequenceTransformer {
	return &SequenceTransformer{}
}

// SequenceOptions はシーケンス図変換オプション
type SequenceOptions struct {
	FlowName       string
	IncludeReturn  bool
}

// Transform はASTをシーケンス図に変換する
func (t *SequenceTransformer) Transform(files []*ast.SpecFile, opts *SequenceOptions) (*sequence.Diagram, error) {
	if opts == nil || opts.FlowName == "" {
		return nil, &errors.TransformError{Source: "AST", Target: "SequenceDiagram", Message: "flow name is required"}
	}

	var targetFlow *ast.FlowDecl
	var targetComponent *ast.ComponentDecl

	for _, file := range files {
		if file.Component == nil {
			continue
		}
		for i := range file.Component.Body.Flows {
			if file.Component.Body.Flows[i].Name == opts.FlowName {
				targetFlow = &file.Component.Body.Flows[i]
				targetComponent = file.Component
				break
			}
		}
	}

	if targetFlow == nil {
		return nil, &errors.TransformError{
			Source:  "AST",
			Target:  "SequenceDiagram",
			Message: "flow not found: " + opts.FlowName,
		}
	}

	diagram := &sequence.Diagram{
		Participants: []sequence.Participant{},
		Events:       []sequence.Event{},
	}

	// 自身を参加者として追加
	diagram.Participants = append(diagram.Participants, sequence.Participant{
		ID:   targetComponent.Name,
		Name: targetComponent.Name,
		Type: sequence.ParticipantTypeDefault,
	})

	// 依存関係から参加者を収集
	for _, rel := range targetComponent.Body.Relations {
		if rel.Kind == ast.RelationDependsOn {
			pType := sequence.ParticipantTypeDefault
			if rel.TargetType != nil {
				switch *rel.TargetType {
				case "database":
					pType = sequence.ParticipantTypeDatabase
				case "external":
					pType = sequence.ParticipantTypeExternal
				case "queue":
					pType = sequence.ParticipantTypeQueue
				case "actor":
					pType = sequence.ParticipantTypeActor
				}
			}
			diagram.Participants = append(diagram.Participants, sequence.Participant{
				ID:   rel.Target,
				Name: rel.Target,
				Type: pType,
			})
		}
	}

	// ステップをイベントに変換
	t.transformSteps(targetFlow.Steps, diagram, targetComponent.Name, opts.IncludeReturn)

	return diagram, nil
}

func (t *SequenceTransformer) transformSteps(steps []ast.Step, diagram *sequence.Diagram, from string, includeReturn bool) {
	for _, step := range steps {
		switch s := step.(type) {
		case *ast.CallStep:
			if call, ok := s.Expr.(*ast.CallExpr); ok {
				to := t.getCallTarget(call)
				msgType := sequence.MessageTypeSync
				if s.Await {
					msgType = sequence.MessageTypeAsync
				}
				diagram.Events = append(diagram.Events, &sequence.MessageEvent{
					From:        from,
					To:          to,
					Label:       call.Method,
					MessageType: msgType,
				})
			}

		case *ast.AssignStep:
			if call, ok := s.Value.(*ast.CallExpr); ok {
				to := t.getCallTarget(call)
				diagram.Events = append(diagram.Events, &sequence.MessageEvent{
					From:        from,
					To:          to,
					Label:       call.Method,
					MessageType: sequence.MessageTypeSync,
				})
				if includeReturn {
					diagram.Events = append(diagram.Events, &sequence.MessageEvent{
						From:        to,
						To:          from,
						Label:       s.Variable,
						MessageType: sequence.MessageTypeReturn,
					})
				}
			}

		case *ast.IfStep:
			fragment := &sequence.FragmentEvent{
				Type:   sequence.FragmentTypeAlt,
				Label:  "condition",
				Events: []sequence.Event{},
			}
			t.transformSteps(s.Then, diagram, from, includeReturn)
			if len(s.Else) > 0 {
				fragment.AltLabel = "else"
				t.transformSteps(s.Else, diagram, from, includeReturn)
			}
			diagram.Events = append(diagram.Events, fragment)

		case *ast.ForStep:
			fragment := &sequence.FragmentEvent{
				Type:   sequence.FragmentTypeLoop,
				Label:  "for " + s.Variable,
				Events: []sequence.Event{},
			}
			t.transformSteps(s.Body, diagram, from, includeReturn)
			diagram.Events = append(diagram.Events, fragment)

		case *ast.WhileStep:
			fragment := &sequence.FragmentEvent{
				Type:   sequence.FragmentTypeLoop,
				Label:  "while",
				Events: []sequence.Event{},
			}
			t.transformSteps(s.Body, diagram, from, includeReturn)
			diagram.Events = append(diagram.Events, fragment)
		}
	}
}

func (t *SequenceTransformer) getCallTarget(call *ast.CallExpr) string {
	if v, ok := call.Object.(*ast.VariableExpr); ok {
		return v.Name
	}
	return "Unknown"
}
