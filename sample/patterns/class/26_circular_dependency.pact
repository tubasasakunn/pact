// Pattern 26: Circular dependency (A depends on B, B depends on A)
component ServiceA {
	type ServiceAData {
		id: string
		name: string
	}

	depends on ServiceB

	provides ServiceAAPI {
		GetFromA() -> ServiceAData
		ProcessWithB() -> string
	}
}

component ServiceB {
	type ServiceBData {
		id: string
		value: int
	}

	depends on ServiceA

	provides ServiceBAPI {
		GetFromB() -> ServiceBData
		ProcessWithA() -> string
	}
}
