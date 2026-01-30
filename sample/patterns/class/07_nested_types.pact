// Pattern 7: Component with nested type references
component NestedTypes {
	type Address {
		street: string
		city: string
		zipCode: string
	}
	
	type Person {
		id: string
		name: string
		homeAddress: Address
		workAddress: Address
	}
}
