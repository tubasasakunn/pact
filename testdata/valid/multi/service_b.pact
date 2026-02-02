// Service B - uses shared definitions
import "./shared.pact"

component ServiceB {
	type ServiceBData {
		id: string
	}

	depends on ServiceA

	provides ExecuteAPI {
		Execute(entityId: string)
	}
}
