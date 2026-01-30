// Pattern 27: Diamond dependency (A depends on B and C, B and C depend on D)
component TopService {
	type TopData {
		id: string
		status: string
	}

	depends on LeftService
	depends on RightService

	provides TopAPI {
		ExecuteAll() -> TopData
	}
}

component LeftService {
	type LeftData {
		leftId: string
		leftValue: int
	}

	depends on BaseService

	provides LeftAPI {
		ProcessLeft() -> LeftData
	}
}

component RightService {
	type RightData {
		rightId: string
		rightValue: int
	}

	depends on BaseService

	provides RightAPI {
		ProcessRight() -> RightData
	}
}

component BaseService {
	type BaseData {
		baseId: string
		timestamp: string
	}

	provides BaseAPI {
		GetBase() -> BaseData
		Initialize() -> bool
	}
}
