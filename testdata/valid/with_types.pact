// Component with type definitions
type Email = string
type Age = int

type Address {
	street: string
	city: string
	zipCode: string
}

type UserStatus = "active" | "inactive" | "pending"

component User {
	id: string
	email: Email
	age: Age
	address: Address
	status: UserStatus
}
