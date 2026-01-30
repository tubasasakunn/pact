// Component with type definitions
component User {
	type Address {
		street: string
		city: string
		zipCode: string
	}

	type UserData {
		id: string
		email: string
		age: int
		address: Address
	}

	enum UserStatus {
		active
		inactive
		pending
	}
}
