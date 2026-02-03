// Package transformer provides AST to diagram model transformations.
//
// Each diagram type has a dedicated transformer:
//   - ClassTransformer:    AST → class.Diagram
//   - SequenceTransformer: AST → sequence.Diagram
//   - StateTransformer:    AST → state.Diagram
//   - FlowTransformer:     AST → flow.Diagram
//
// All transformers follow the same method pattern:
//
//	Transform(files []*ast.SpecFile, opts *XxxOptions) (*xxx.Diagram, error)
//
// Transformers are stateless (except FlowTransformer which uses internal counters
// that are reset on each Transform call).
package transformer
