// Pattern 20: Component with contains relationship
component Car {
	type CarData {
		make: string
		model: string
		year: int
	}
	
	contains Engine
	contains Transmission
	contains Wheels
}

component Engine {
	type EngineData {
		horsepower: int
		cylinders: int
	}
}

component Transmission {
	type TransmissionData {
		gears: int
		automatic: bool
	}
}

component Wheels {
	type WheelData {
		size: int
		count: int
	}
}
