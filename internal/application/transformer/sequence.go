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
		// 単一コンポーネント
		if file.Component != nil {
			for i := range file.Component.Body.Flows {
				if file.Component.Body.Flows[i].Name == opts.FlowName {
					targetFlow = &file.Component.Body.Flows[i]
					targetComponent = file.Component
					break
				}
			}
		}

		// 複数コンポーネント
		if targetFlow == nil {
			for j := range file.Components {
				comp := &file.Components[j]
				for i := range comp.Body.Flows {
					if comp.Body.Flows[i].Name == opts.FlowName {
						targetFlow = &comp.Body.Flows[i]
						targetComponent = comp
						break
					}
				}
				if targetFlow != nil {
					break
				}
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
	diagram.Events = t.transformSteps(targetFlow.Steps, targetComponent.Name, opts.IncludeReturn)

	return diagram, nil
}

func (t *SequenceTransformer) transformSteps(steps []ast.Step, from string, includeReturn bool) []sequence.Event {
	var events []sequence.Event

	for _, step := range steps {
		switch s := step.(type) {
		case *ast.CallStep:
			if call, ok := s.Expr.(*ast.CallExpr); ok {
				to := t.getCallTarget(call)
				msgType := sequence.MessageTypeSync
				if s.Await {
					msgType = sequence.MessageTypeAsync
				}
				events = append(events, &sequence.MessageEvent{
					From:        from,
					To:          to,
					Label:       call.Method,
					MessageType: msgType,
				})
			}

		case *ast.AssignStep:
			if call, ok := s.Value.(*ast.CallExpr); ok {
				to := t.getCallTarget(call)
				events = append(events, &sequence.MessageEvent{
					From:        from,
					To:          to,
					Label:       call.Method,
					MessageType: sequence.MessageTypeSync,
				})
				if includeReturn {
					events = append(events, &sequence.MessageEvent{
						From:        to,
						To:          from,
						Label:       s.Variable,
						MessageType: sequence.MessageTypeReturn,
					})
				}
			}

		case *ast.IfStep:
			// Then節のイベントを収集
			thenEvents := t.transformSteps(s.Then, from, includeReturn)

			fragment := &sequence.FragmentEvent{
				Type:   sequence.FragmentTypeAlt,
				Label:  "condition",
				Events: thenEvents,
			}

			// Else節がある場合
			if s.Else != nil {
				elseEvents := t.transformSteps(s.Else, from, includeReturn)
				fragment.AltLabel = "else"
				fragment.AltEvents = elseEvents
			}
			events = append(events, fragment)

		case *ast.ForStep:
			bodyEvents := t.transformSteps(s.Body, from, includeReturn)
			fragment := &sequence.FragmentEvent{
				Type:   sequence.FragmentTypeLoop,
				Label:  "for " + s.Variable,
				Events: bodyEvents,
			}
			events = append(events, fragment)

		case *ast.WhileStep:
			bodyEvents := t.transformSteps(s.Body, from, includeReturn)
			fragment := &sequence.FragmentEvent{
				Type:   sequence.FragmentTypeLoop,
				Label:  "while",
				Events: bodyEvents,
			}
			events = append(events, fragment)
		}
	}

	return events
}

func (t *SequenceTransformer) getCallTarget(call *ast.CallExpr) string {
	if v, ok := call.Object.(*ast.VariableExpr); ok {
		return v.Name
	}
	return "Unknown"
}
