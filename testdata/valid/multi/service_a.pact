// Service A - uses shared definitions
import "./shared.pact"

component ServiceA {
	type ServiceAData {
		id: string
	}

	provides ProcessAPI {
		Process(entityId: string)
	}
}
