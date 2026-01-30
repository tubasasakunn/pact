// Service B - uses shared definitions
import "./shared.pact"

component ServiceB {
	id: ID

	relation ServiceA: uses

	method Execute(entityId: ID): void
}
