// Service A - uses shared definitions
import "./shared.pact"

component ServiceA {
	id: ID

	method Process(entityId: ID): void
}
